package models

import (
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/models"
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	UserID     uint64 `gorm:"not null" json:"user_id"`
	DeviceName string `gorm:"type:varchar(100);not null" json:"device_name"`
	DeviceCode string `gorm:"type:varchar(100);uniqueIndex;not null" json:"device_code"`
	Status     string `gorm:"type:varchar(20);default:'offline';check:status IN ('online','offline')" json:"status"`

	User      models.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Sensors   []Sensor    `gorm:"foreignKey:DeviceID" json:"sensors"`
	Actuators []Actuator  `gorm:"foreignKey:DeviceID" json:"actuators"`
}
