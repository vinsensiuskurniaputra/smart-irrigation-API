package models

import (
	"gorm.io/gorm"
)

type IrrigationRule struct {
	gorm.Model
	PlantName         string  `gorm:"type:varchar(100);not null" json:"plant_name"`
	MinMoisture       float64 `gorm:"not null" json:"min_moisture"`
	MaxMoisture       float64 `gorm:"not null" json:"max_moisture"`
	PreferredTemp     float64 `gorm:"not null" json:"preferred_temp"`
	PreferredHumidity float64 `gorm:"not null" json:"preferred_humidity"`
}
