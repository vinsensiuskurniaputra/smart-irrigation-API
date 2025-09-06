package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicemodels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	models "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/models"
	irrigationrepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/repositories"
	irrigationdto "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/domain/dto"
	"gorm.io/gorm"
)

type PredictionResult struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

type IrrigationUsecase struct {
	PredictionURL string
	HTTPClient    *http.Client
	PlantRepo     irrigationrepo.PlantRepository
	DB            *gorm.DB
	MQTTClient    mqtt.Client
}

func NewIrrigationUsecase(predictionURL string, plantRepo irrigationrepo.PlantRepository, db *gorm.DB, mqttClient mqtt.Client) *IrrigationUsecase {
	return &IrrigationUsecase{PredictionURL: predictionURL, HTTPClient: &http.Client{}, PlantRepo: plantRepo, DB: db, MQTTClient: mqttClient}
}

func (uc *IrrigationUsecase) PredictPlant(fileField string, filename string, fileBytes []byte) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileField, filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, uc.PredictionURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return uc.HTTPClient.Do(req)
}

var aiLabelMap = map[int]string{
	0: "cactus",
	1: "chili",
	2: "monstera",
	3: "spinach",
	4: "tomato",
}

func (uc *IrrigationUsecase) SavePredictedPlant(deviceID uint64, labelIndex int, imageURL string) (interface{}, error) {
	name, ok := aiLabelMap[labelIndex]
	if !ok {
		return nil, fmt.Errorf("invalid label index")
	}
	// check if plant already exists for device
	existing, err := uc.PlantRepo.FindByDevice(deviceID)
	if err == nil && existing != nil {
		// update existing plant
		p, err2 := uc.PlantRepo.UpdatePlant(uint64(existing.ID), name, &imageURL)
		if err2 != nil {
			return nil, err2
		}
		uc.publishRule(p)
		return toPlantDTO(p), nil
	}
	// create new plant
	plant, err := uc.PlantRepo.CreatePlant(deviceID, name, imageURL)
	if err != nil {
		return nil, err
	}
	loaded, err := uc.PlantRepo.GetPlant(uint64(plant.ID))
	if err == nil {
		uc.publishRule(loaded)
	}
	return toPlantDTO(loaded), nil
}

func (uc *IrrigationUsecase) GetPlant(id uint64) (*irrigationdto.PlantDTO, error) {
	p, err := uc.PlantRepo.GetPlant(id)
	if err != nil {
		return nil, err
	}
	return toPlantDTO(p), nil
}

func (uc *IrrigationUsecase) ListPlantsByDevice(deviceID uint64) ([]*irrigationdto.PlantDTO, error) {
	plants, err := uc.PlantRepo.ListPlantsByDevice(deviceID)
	if err != nil {
		return nil, err
	}
	res := make([]*irrigationdto.PlantDTO, 0, len(plants))
	for _, p := range plants {
		res = append(res, toPlantDTO(p))
	}
	return res, nil
}

func (uc *IrrigationUsecase) UpdatePlant(id uint64, plantName string, imageURL *string) (*irrigationdto.PlantDTO, error) {
	p, err := uc.PlantRepo.UpdatePlant(id, plantName, imageURL)
	if err != nil {
		return nil, err
	}
	return toPlantDTO(p), nil
}

// UpdatePlantLabel remaps plant by label index (if provided) and/or updates image
func (uc *IrrigationUsecase) UpdatePlantLabel(id uint64, labelIndex *int, imageURL *string) (*irrigationdto.PlantDTO, error) {
	var plantName string
	if labelIndex != nil {
		n, ok := aiLabelMap[*labelIndex]
		if !ok {
			return nil, fmt.Errorf("invalid label index")
		}
		plantName = n
	}
	p, err := uc.PlantRepo.UpdatePlant(id, plantName, imageURL)
	if err != nil {
		return nil, err
	}
	uc.publishRule(p)
	return toPlantDTO(p), nil
}

// publishRule publishes irrigation rule to topic:
// device/{device_code}/rule/{plant_name}/{min}/{max}/{preferred_temp}/{preferred_humidity}
// values formatted with no trailing zeros when possible.
func (uc *IrrigationUsecase) publishRule(p *models.Plant) {
	if uc.MQTTClient == nil || !uc.MQTTClient.IsConnected() || p == nil {
		return
	}
	// ensure rule loaded
	if p.IrrigationRule.ID == 0 {
		var reload models.Plant
		if err := uc.DB.Preload("IrrigationRule").First(&reload, p.ID).Error; err == nil {
			p = &reload
		}
	}
	if p.IrrigationRule.ID == 0 {
		return
	}
	var device devicemodels.Device
	if err := uc.DB.Select("device_code").First(&device, p.DeviceID).Error; err != nil {
		return
	}

	// Build JSON payload
	rule := p.IrrigationRule
	payload := map[string]interface{}{
		"id":                 rule.ID,
		"plant_name":         rule.PlantName,
		"min_moisture":       rule.MinMoisture,
		"max_moisture":       rule.MaxMoisture,
		"preferred_temp":     rule.PreferredTemp,
		"preferred_humidity": rule.PreferredHumidity,
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	topic := fmt.Sprintf("device/%s/rule", device.DeviceCode)
	token := uc.MQTTClient.Publish(topic, 0, true, bytes)
	token.Wait()
	if token.Error() != nil {
		log.Printf("failed publish rule topic: %v", token.Error())
	}
}

func toPlantDTO(p *models.Plant) *irrigationdto.PlantDTO {
	dto := &irrigationdto.PlantDTO{
		ID:               uint(p.ID),
		DeviceID:         p.DeviceID,
		IrrigationRuleID: p.IrrigationRuleID,
		PlantName:        p.PlantName,
		ImageURL:         p.ImageURL,
	}
	if p.IrrigationRule.ID != 0 {
		dto.Rule = &irrigationdto.IrrigationRuleDTO{
			ID:                uint(p.IrrigationRule.ID),
			PlantName:         p.IrrigationRule.PlantName,
			MinMoisture:       p.IrrigationRule.MinMoisture,
			MaxMoisture:       p.IrrigationRule.MaxMoisture,
			PreferredTemp:     p.IrrigationRule.PreferredTemp,
			PreferredHumidity: p.IrrigationRule.PreferredHumidity,
		}
	}
	return dto
}
