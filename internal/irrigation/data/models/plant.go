package models

import (
	deviceModel "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

type Plant struct {
	gorm.Model
	DeviceID         uint64 `gorm:"not null" json:"device_id"`
	IrrigationRuleID uint64 `gorm:"not null" json:"irrigation_rule_id"`
	PlantName        string `gorm:"type:varchar(100);not null" json:"plant_name"`
	ImageURL         string `gorm:"type:text" json:"image_url"`

	Device         deviceModel.Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"device"`
	IrrigationRule IrrigationRule     `gorm:"foreignKey:IrrigationRuleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"irrigation_rule"`
}
