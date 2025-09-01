package models

import (
	"gorm.io/gorm"
)

type ActuatorLog struct {
	gorm.Model
	ActuatorID  uint64   `gorm:"not null" json:"actuator_id"`
	Actuator    Actuator `gorm:"foreignKey:ActuatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"actuator"`
	Action      string   `gorm:"type:varchar(10);not null;check:action IN ('on','off')" json:"action"`
	TriggeredBy string   `gorm:"type:varchar(20);not null;check:triggered_by IN ('manual','auto','schedule')" json:"triggered_by"`
}
