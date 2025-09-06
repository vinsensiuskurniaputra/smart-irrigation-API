package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	irrigationusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/domain/usecase"
)

var labelIndexMap = map[string]int{
	"cactus":   0,
	"chili":    1,
	"monstera": 2,
	"spinach":  3,
	"tomato":   4,
}

func labelToIndex(label string) (int, bool) {
	label = strings.ToLower(label)
	idx, ok := labelIndexMap[label]
	return idx, ok
}

// tries to find a "label" field or first string value
func extractLabelString(payload interface{}) string {
	if m, ok := payload.(map[string]interface{}); ok {
		if v, ok := m["label"]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		// Sometimes ML returns {"prediction": {"label":"tomato", ...}}
		if p, ok := m["prediction"]; ok {
			if pm, ok := p.(map[string]interface{}); ok {
				if v, ok := pm["label"]; ok {
					if s, ok := v.(string); ok {
						return s
					}
				}
			}
		}
		// fallback: scan values
		for _, v := range m {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

type IrrigationHandler struct {
	uc *irrigationusecase.IrrigationUsecase
}

func NewIrrigationHandler(uc *irrigationusecase.IrrigationUsecase) *IrrigationHandler {
	return &IrrigationHandler{uc: uc}
}

// PredictPlant forwards image file to ML service and streams back its JSON response.
// Expects form field name "file".
func (h *IrrigationHandler) PredictPlant(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer opened.Close()
	bytes, err := io.ReadAll(opened)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
		return
	}
	resp, err := h.uc.PredictPlant("file", file.Filename, bytes)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	proxyBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed read prediction response"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		var payload interface{}
		if err := json.Unmarshal(proxyBody, &payload); err == nil {
			// attempt extract label string
			labelStr := extractLabelString(payload)
			labelIndex := -1
			if labelStr != "" {
				if idx, ok := labelToIndex(labelStr); ok {
					labelIndex = idx
				}
			}
			if m, ok := payload.(map[string]interface{}); ok {
				m["label_index"] = labelIndex
				modified, _ := json.Marshal(m)
				c.Data(resp.StatusCode, "application/json", modified)
				return
			}
		}
	}
	// fallback raw proxy
	c.Data(resp.StatusCode, contentType, proxyBody)
}

type savePlantRequest struct {
	LabelIndex int    `json:"label_index"`
	ImageURL   string `json:"image_url"`
}

// SavePredicted stores plant with auto irrigation rule mapping based on AI label index.
func (h *IrrigationHandler) SavePredicted(c *gin.Context) {
	deviceIDParam := c.Param("device_id")
	deviceID, err := strconv.ParseUint(deviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}
	var req savePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plant, err := h.uc.SavePredictedPlant(deviceID, req.LabelIndex, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": plant})
}

// GetPlant returns a plant by id
func (h *IrrigationHandler) GetPlant(c *gin.Context) {
	idParam := c.Param("plant_id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	plant, err := h.uc.GetPlant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": plant})
}

type updatePlantRequest struct {
	LabelIndex *int    `json:"label_index"` // optional; when provided remaps plant
	ImageURL   *string `json:"image_url"`
}

// UpdatePlant updates label_index (thus plant + irrigation rule) and/or image url
func (h *IrrigationHandler) UpdatePlant(c *gin.Context) {
	idParam := c.Param("plant_id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updatePlantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plant, err := h.uc.UpdatePlantLabel(id, req.LabelIndex, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": plant})
}

// ListPlantsByDevice returns all plants for a given device id
func (h *IrrigationHandler) ListPlantsByDevice(c *gin.Context) {
	deviceIDParam := c.Param("device_id")
	deviceID, err := strconv.ParseUint(deviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}
	plants, err := h.uc.ListPlantsByDevice(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": plants})
}
