package usecase

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
)

// Topic pattern expected: device/{device_code}/actuator/{actuator_id}/actual-status
// Payload JSON: {"value":"on"} or {"value":"off"}

type ActuatorStatusConsumer struct {
	repo devicerepo.ActuatorRepository
}

func NewActuatorStatusConsumer(repo devicerepo.ActuatorRepository) *ActuatorStatusConsumer {
	return &ActuatorStatusConsumer{repo: repo}
}

type actuatorStatusPayload struct {
	Value string `json:"value"`
}

func (c *ActuatorStatusConsumer) Handler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		parts := strings.Split(topic, "/")
		// device/{device_code}/actuator/{actuator_id}/actual-status
		if len(parts) != 5 || parts[0] != "device" || parts[2] != "actuator" || parts[4] != "actual-status" {
			// not our concern
			return
		}

		actuatorIDStr := parts[3]
		actuatorID, err := strconv.ParseUint(actuatorIDStr, 10, 64)
		if err != nil || actuatorID == 0 {
			log.Printf("MQTT: invalid actuator id in topic: %s", actuatorIDStr)
			return
		}

		var p actuatorStatusPayload
		if err := json.Unmarshal(msg.Payload(), &p); err != nil {
			log.Printf("MQTT: invalid actuator status payload: %v", err)
			return
		}
		val := strings.ToLower(strings.TrimSpace(p.Value))
		if val != "on" && val != "off" {
			log.Printf("MQTT: invalid actuator status value: %s", p.Value)
			return
		}
		if err := c.repo.UpdateStatus(actuatorID, val); err != nil {
			log.Printf("MQTT: failed updating actuator status (id=%d): %v", actuatorID, err)
			return
		}
		log.Printf("MQTT: actuator actual status updated id=%d value=%s", actuatorID, val)
	}
}
