package models

import (
	"gorm.io/gorm"
)

type Sensor struct {
	gorm.Model
	DeviceID   uint64 `gorm:"not null" json:"device_id"`
	SensorType string `gorm:"type:varchar(50);not null;check:sensor_type IN ('soil_moisture','temperature','humidity')" json:"sensor_type"`
	Device     Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"device"`
}
