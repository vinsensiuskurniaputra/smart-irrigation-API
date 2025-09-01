package models

import (
	"gorm.io/gorm"
)

type Actuator struct {
	gorm.Model
	DeviceID     uint64    `gorm:"not null" json:"device_id"`
	Device       Device    `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"device"`
	ActuatorName string    `gorm:"type:varchar(100);not null" json:"actuator_name"`
	Type         string    `gorm:"type:varchar(20);not null;check:type IN ('pump','lamp','fan','sprayer','other')" json:"type"`
	PinNumber    string    `gorm:"type:varchar(10);not null" json:"pin_number"`
	Status       string    `gorm:"type:varchar(10);default:'off';check:status IN ('on','off')" json:"status"`
}
