package usecase

import (
	"errors"

	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/domain/dto"
)

type DeviceUsecase struct {
	repo devicerepo.DeviceRepository
}

func NewDeviceUsecase(r devicerepo.DeviceRepository) *DeviceUsecase {
	return &DeviceUsecase{repo: r}
}

func (uc *DeviceUsecase) List(userID uint, page, pageSize int) ([]dto.DeviceListItemDTO, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	devices, err := uc.repo.FindByUser(userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := uc.repo.CountByUser(userID)
	if err != nil {
		return nil, 0, err
	}
	list := make([]dto.DeviceListItemDTO, 0, len(devices))
	for _, d := range devices {
		list = append(list, dto.DeviceListItemDTO{
			ID:         d.ID,
			DeviceName: d.DeviceName,
			DeviceCode: d.DeviceCode,
			Status:     d.Status,
		})
	}
	return list, total, nil
}

func (uc *DeviceUsecase) Detail(id uint, userID uint) (*dto.DeviceDetailDTO, error) {
	device, err := uc.repo.FindDetail(id, userID)
	if err != nil {
		return nil, err
	}
	if device.ID == 0 {
		return nil, errors.New("device not found")
	}
	detail := &dto.DeviceDetailDTO{
		ID:         device.ID,
		DeviceName: device.DeviceName,
		DeviceCode: device.DeviceCode,
		Status:     device.Status,
	}
	for _, s := range device.Sensors {
		detail.Sensors = append(detail.Sensors, dto.SensorDTO{ID: s.ID, SensorType: s.SensorType})
	}
	for _, a := range device.Actuators {
		detail.Actuators = append(detail.Actuators, dto.ActuatorDTO{ID: a.ID, ActuatorName: a.ActuatorName, Type: a.Type, PinNumber: a.PinNumber, Status: a.Status, Mode: a.Mode})
	}
	return detail, nil
}
