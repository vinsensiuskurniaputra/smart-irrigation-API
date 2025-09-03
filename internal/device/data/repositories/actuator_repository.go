package repositories

import (
	"errors"

	devicemodels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

type ActuatorRepository interface {
	FindByID(id uint64) (*devicemodels.Actuator, error)
	UpdateStatus(id uint64, status string) error
	UpdateMode(id uint64, mode string) error
	LogAction(actuatorID uint64, action, triggeredBy string) error
}

type actuatorRepository struct {
	db *gorm.DB
}

func NewActuatorRepository(db *gorm.DB) ActuatorRepository {
	return &actuatorRepository{db: db}
}

func (r *actuatorRepository) FindByID(id uint64) (*devicemodels.Actuator, error) {
	var a devicemodels.Actuator
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *actuatorRepository) UpdateStatus(id uint64, status string) error {
	if status != "on" && status != "off" {
		return errors.New("invalid status")
	}
	return r.db.Model(&devicemodels.Actuator{}).Where("id = ?", id).Update("status", status).Error
}

func (r *actuatorRepository) UpdateMode(id uint64, mode string) error {
	if mode != "auto" && mode != "manual" {
		return errors.New("invalid mode")
	}
	return r.db.Model(&devicemodels.Actuator{}).Where("id = ?", id).Update("mode", mode).Error
}

func (r *actuatorRepository) LogAction(actuatorID uint64, action, triggeredBy string) error {
	log := devicemodels.ActuatorLog{ActuatorID: actuatorID, Action: action, TriggeredBy: triggeredBy}
	return r.db.Create(&log).Error
}
