package usecase

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
)

// Topic pattern: device/{device_code}/status
// Payload: {"value":"online"} or {"value":"offline"}

type DeviceStatusConsumer struct {
	repo devicerepo.DeviceRepository
}

func NewDeviceStatusConsumer(r devicerepo.DeviceRepository) *DeviceStatusConsumer {
	return &DeviceStatusConsumer{repo: r}
}

type deviceStatusPayload struct {
	Value string `json:"value"`
}

func (c *DeviceStatusConsumer) Handler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic() // device/{device_code}/status
		parts := strings.Split(topic, "/")
		if len(parts) != 3 || parts[0] != "device" || parts[2] != "status" {
			return
		}
		deviceCode := parts[1]
		if deviceCode == "" {
			return
		}
		var p deviceStatusPayload
		if err := json.Unmarshal(msg.Payload(), &p); err != nil {
			log.Printf("MQTT: invalid device status payload: %v", err)
			return
		}
		val := strings.ToLower(strings.TrimSpace(p.Value))
		if val != "online" && val != "offline" {
			log.Printf("MQTT: invalid device status value: %s", p.Value)
			return
		}
		if err := c.repo.UpdateStatusByCode(deviceCode, val); err != nil {
			log.Printf("MQTT: failed update device status code=%s: %v", deviceCode, err)
			return
		}
		log.Printf("MQTT: device status updated code=%s value=%s", deviceCode, val)
	}
}
