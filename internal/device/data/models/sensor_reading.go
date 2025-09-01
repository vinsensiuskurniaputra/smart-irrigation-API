package models

import (
	"time"

	"gorm.io/gorm"
)

type SensorReading struct {
	gorm.Model
	SensorID   uint64    `gorm:"not null" json:"sensor_id"`
	Value      float64   `gorm:"not null" json:"value"`
	RecordedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"recorded_at"`

	Sensor Sensor `gorm:"foreignKey:SensorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sensor"`
}
