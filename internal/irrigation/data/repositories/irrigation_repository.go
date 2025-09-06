package repositories

import (
	"errors"
	"strings"

	models "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/models"
	"gorm.io/gorm"
)

type PlantRepository interface {
	CreatePlant(deviceID uint64, plantName string, imageURL string) (*models.Plant, error)
	GetPlant(id uint64) (*models.Plant, error)
	UpdatePlant(id uint64, plantName string, imageURL *string) (*models.Plant, error)
	ListPlantsByDevice(deviceID uint64) ([]*models.Plant, error)
	FindByDevice(deviceID uint64) (*models.Plant, error)
}

type plantRepository struct{ db *gorm.DB }

func NewPlantRepository(db *gorm.DB) PlantRepository { return &plantRepository{db: db} }

func normalize(name string) string { return strings.ToLower(name) }

func (r *plantRepository) findRuleID(plantName string) (uint64, error) {
	var rule models.IrrigationRule
	if err := r.db.Where("LOWER(plant_name) = ?", normalize(plantName)).First(&rule).Error; err != nil {
		return 0, err
	}
	return uint64(rule.ID), nil
}

func (r *plantRepository) CreatePlant(deviceID uint64, plantName string, imageURL string) (*models.Plant, error) {
	ruleID, err := r.findRuleID(plantName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("irrigation rule not found for plant")
		}
		return nil, err
	}
	plant := &models.Plant{DeviceID: deviceID, IrrigationRuleID: ruleID, PlantName: plantName, ImageURL: imageURL}
	if err := r.db.Create(plant).Error; err != nil {
		return nil, err
	}
	return plant, nil
}

func (r *plantRepository) GetPlant(id uint64) (*models.Plant, error) {
	var p models.Plant
	if err := r.db.Preload("IrrigationRule").First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *plantRepository) UpdatePlant(id uint64, plantName string, imageURL *string) (*models.Plant, error) {
	var p models.Plant
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	if plantName != "" && plantName != p.PlantName {
		ruleID, err := r.findRuleID(plantName)
		if err != nil {
			return nil, err
		}
		p.PlantName = plantName
		p.IrrigationRuleID = ruleID
	}
	if imageURL != nil {
		p.ImageURL = *imageURL
	}
	if err := r.db.Save(&p).Error; err != nil {
		return nil, err
	}
	r.db.Preload("IrrigationRule").First(&p, id)
	return &p, nil
}

func (r *plantRepository) ListPlantsByDevice(deviceID uint64) ([]*models.Plant, error) {
	var plants []*models.Plant
	if err := r.db.Preload("IrrigationRule").Where("device_id = ?", deviceID).Find(&plants).Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (r *plantRepository) FindByDevice(deviceID uint64) (*models.Plant, error) {
	var p models.Plant
	if err := r.db.Preload("IrrigationRule").Where("device_id = ?", deviceID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}
