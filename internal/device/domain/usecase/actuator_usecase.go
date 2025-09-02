package usecase

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
)

type ActuatorControlUsecase struct {
	repo       devicerepo.ActuatorRepository
	mqttClient mqtt.Client
}

func NewActuatorControlUsecase(r devicerepo.ActuatorRepository, client mqtt.Client) *ActuatorControlUsecase {
	return &ActuatorControlUsecase{repo: r, mqttClient: client}
}

// Control toggles actuator status (on/off), updates DB, logs and publishes MQTT command
func (uc *ActuatorControlUsecase) Control(actuatorID uint64, desired string, triggeredBy string) error {
	_, err := uc.repo.FindByID(actuatorID)
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
		topic := fmt.Sprintf("actuators/%d/command", actuatorID)
		payload := fmt.Sprintf("{\"action\":\"%s\"}", desired)
		token := uc.mqttClient.Publish(topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Printf("failed publish actuator command: %v", token.Error())
		}
	}
	return nil
}
