package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	devicemodels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type LiveDataHandler struct {
	db         *gorm.DB
	sensorRepo *devicerepo.SensorReadingRepository
}

func NewLiveDataHandler(db *gorm.DB) *LiveDataHandler {
	return &LiveDataHandler{db: db, sensorRepo: devicerepo.NewSensorReadingRepository(db)}
}

type liveSensorResponse struct {
	SensorID uint64    `json:"sensor_id"`
	Type     string    `json:"type"`
	Readings []reading `json:"readings"`
}
type reading struct {
	Value      float64   `json:"value"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (h *LiveDataHandler) DeviceLive(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	_ = userIDVal.(uint) // could be used for auth check of device ownership

	deviceIDParam := c.Param("id")
	deviceID64, err := strconv.ParseUint(deviceIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Load sensors for device
	var sensors []devicemodels.Sensor
	if err := h.db.Where("device_id = ?", deviceID64).Find(&sensors).Error; err != nil {
		log.Printf("error load sensors: %v", err)
		return
	}
	
	sensorIDs := make([]uint64, 0, len(sensors))
	for _, s := range sensors {
		sensorIDs = append(sensorIDs, uint64(s.ID))
	}

	send := func() error {
		dataMap, err := h.sensorRepo.LastNSensorReadings(sensorIDs, 10)
		if err != nil {
			return err
		}
		resp := make([]liveSensorResponse, 0, len(sensors))
		for _, s := range sensors {
			readings := dataMap[uint64(s.ID)]
			rs := make([]reading, 0, len(readings))
			for i := len(readings) - 1; i >= 0; i-- { // chronological (oldest first)
				rs = append(rs, reading{Value: readings[i].Value, RecordedAt: readings[i].RecordedAt})
			}
			resp = append(resp, liveSensorResponse{SensorID: uint64(s.ID), Type: s.SensorType, Readings: rs})
		}
		return conn.WriteJSON(gin.H{"sensors": resp, "ts": time.Now()})
	}

	// send immediately
	if err := send(); err != nil {
		log.Printf("websocket initial send error: %v", err)
		return
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := send(); err != nil {
			log.Printf("websocket write error: %v", err)
			return
		}
	}
}
