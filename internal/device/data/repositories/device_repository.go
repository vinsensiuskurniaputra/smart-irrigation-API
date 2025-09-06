package repositories

import (
	devicemodels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	FindByUser(userID uint, limit, offset int) ([]devicemodels.Device, error)
	CountByUser(userID uint) (int64, error)
	FindDetail(id uint, userID uint) (*devicemodels.Device, error)
	UpdateStatusByCode(deviceCode string, status string) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) FindByUser(userID uint, limit, offset int) ([]devicemodels.Device, error) {
	var devices []devicemodels.Device
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	if err := query.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) CountByUser(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&devicemodels.Device{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *deviceRepository) FindDetail(id uint, userID uint) (*devicemodels.Device, error) {
	var device devicemodels.Device
	if err := r.db.Preload("User").Preload("Sensors").Preload("Actuators").Where("id = ? AND user_id = ?", id, userID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) UpdateStatusByCode(deviceCode string, status string) error {
	if status != "online" && status != "offline" {
		return nil // silently ignore invalid (could return error if preferred)
	}
	return r.db.Model(&devicemodels.Device{}).Where("device_code = ?", deviceCode).Update("status", status).Error
}
