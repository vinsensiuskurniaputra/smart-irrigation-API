package usecase

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicemodels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	"gorm.io/gorm"
)

type ActuatorControlUsecase struct {
	repo       devicerepo.ActuatorRepository
	db         *gorm.DB
	mqttClient mqtt.Client
}

func NewActuatorControlUsecase(r devicerepo.ActuatorRepository, db *gorm.DB, client mqtt.Client) *ActuatorControlUsecase {
	return &ActuatorControlUsecase{repo: r, db: db, mqttClient: client}
}

// Control toggles actuator status (on/off), updates DB, logs and publishes MQTT command
func (uc *ActuatorControlUsecase) Control(actuatorID uint64, desired string, triggeredBy string) error {
	act, err := uc.repo.FindByID(actuatorID)
	if err != nil {
		return err
	}
	if desired != "on" && desired != "off" {
		return fmt.Errorf("invalid desired state")
	}
	if err := uc.repo.UpdateStatus(actuatorID, desired); err != nil {
		return err
	}
	if err := uc.repo.LogAction(actuatorID, desired, triggeredBy); err != nil {
		log.Printf("warn: failed to log actuator action: %v", err)
	}
	// Publish MQTT command. Topic pattern: actuators/{actuator_id}/command
	if uc.mqttClient != nil && uc.mqttClient.IsConnected() {
		// fetch device to compose topic with device_code
		var device devicemodels.Device
		if err := uc.db.Select("device_code").First(&device, act.DeviceID).Error; err != nil {
			log.Printf("warn: cannot fetch device for actuator topic: %v", err)
		}
		topic := fmt.Sprintf("device/%s/actuator/%d", device.DeviceCode, actuatorID)
		payload := fmt.Sprintf("{\"value\":\"%s\"}", desired)
		token := uc.mqttClient.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("failed publish actuator command: %v", token.Error())
		}
	}
	return nil
}
