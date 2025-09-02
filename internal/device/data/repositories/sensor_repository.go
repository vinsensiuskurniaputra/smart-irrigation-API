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
