package usecase

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
)

// Expected topic pattern: device/{device_code}/sensor/{sensor_id}
// Payload JSON: {"value":"88.5"} or {"value":88.5}

type SensorReadingConsumer struct {
	sensorRepo *devicerepo.SensorReadingRepository
	// optionally we could cache sensorID mapping; kept simple now
}

func NewSensorReadingConsumer(sensorRepo *devicerepo.SensorReadingRepository) *SensorReadingConsumer {
	return &SensorReadingConsumer{sensorRepo: sensorRepo}
}

type sensorPayload struct {
	Value      interface{} `json:"value"`
	RecordedAt time.Time   `json:"recorded_at"`
}

func (c *SensorReadingConsumer) Handler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		parts := strings.Split(topic, "/")
		// device/{device_code}/sensor/{sensor_id}
		if len(parts) != 4 || parts[0] != "device" || parts[2] != "sensor" {
			log.Printf("MQTT: invalid topic format (expected device/{device_code}/sensor/{sensor_id}): %s", topic)
			return
		}
		payload := msg.Payload()
		var sp sensorPayload
		if err := json.Unmarshal(payload, &sp); err != nil {
			log.Printf("MQTT: invalid json payload: %v", err)
			return
		}
		sensorIDStr := parts[3]
		sensorID, err := strconv.ParseUint(sensorIDStr, 10, 64)
		if err != nil {
			log.Printf("MQTT: invalid sensor id in topic: %s", sensorIDStr)
			return
		}
		// Parse value which may be string or number
		var numericValue float64
		switch v := sp.Value.(type) {
		case float64:
			numericValue = v
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				numericValue = f
			} else {
				log.Printf("MQTT: invalid numeric string value: %v", v)
				return
			}
		default:
			log.Printf("MQTT: unsupported value type %T", v)
			return
		}
		if sensorID == 0 {
			log.Printf("MQTT: sensor id not provided (topic %s)", topic)
			return
		}
		if err := c.sensorRepo.Create(sensorID, numericValue); err != nil {
			log.Printf("MQTT: failed inserting reading: %v", err)
			return
		}
		log.Printf("MQTT: inserted reading sensor=%d value=%f", sensorID, numericValue)
	}
}
