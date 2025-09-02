package repositories

import (
	"time"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

type SensorReadingRepository struct {
	db *gorm.DB
}

func NewSensorReadingRepository(db *gorm.DB) *SensorReadingRepository {
	return &SensorReadingRepository{db}
}

func (r *SensorReadingRepository) Create(sensorID uint64, value float64) error {
	reading := models.SensorReading{
		SensorID:   sensorID,
		Value:      value,
		RecordedAt: time.Now(),
	}
	return r.db.Create(&reading).Error
}

func (r *SensorReadingRepository) LastNSensorReadings(sensorIDs []uint64, n int) (map[uint64][]models.SensorReading, error) {
	if n <= 0 {
		n = 10
	}
	result := make(map[uint64][]models.SensorReading)
	for _, id := range sensorIDs {
		var readings []models.SensorReading
		if err := r.db.Where("sensor_id = ?", id).Order("recorded_at DESC").Limit(n).Find(&readings).Error; err != nil {
			return nil, err
		}
		result[id] = readings
	}
	return result, nil
}
