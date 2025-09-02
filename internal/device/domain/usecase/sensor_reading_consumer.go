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

// Expected topic pattern: sensors/{device_code}/{sensor_type}
// Payload JSON: {"value":123.45, "recorded_at":"optional RFC3339"}

type SensorReadingConsumer struct {
    sensorRepo *devicerepo.SensorReadingRepository
    // optionally we could cache sensorID mapping; kept simple now
}

func NewSensorReadingConsumer(sensorRepo *devicerepo.SensorReadingRepository) *SensorReadingConsumer {
    return &SensorReadingConsumer{sensorRepo: sensorRepo}
}

type sensorPayload struct {
    Value      float64   `json:"value"`
    RecordedAt time.Time `json:"recorded_at"`
    SensorID   uint64    `json:"sensor_id"` // alternative if provided
}

func (c *SensorReadingConsumer) Handler() mqtt.MessageHandler {
    return func(client mqtt.Client, msg mqtt.Message) {
        topic := msg.Topic()
        parts := strings.Split(topic, "/")
        if len(parts) < 3 { // sensors/{device_code}/{sensor_type}
            log.Printf("MQTT: invalid topic format: %s", topic)
            return
        }
        payload := msg.Payload()
        var sp sensorPayload
        if err := json.Unmarshal(payload, &sp); err != nil {
            log.Printf("MQTT: invalid json payload: %v", err)
            return
        }
        if sp.Value == 0 && !strings.Contains(string(payload), "0") { // naive check to allow zero
            log.Printf("MQTT: missing value in payload: %s", payload)
            return
        }
        // Prefer sensor_id from payload; else attempt derive from topic numeric segment
        sensorID := sp.SensorID
        if sensorID == 0 && len(parts) >= 4 { // sensors/{device_code}/{sensor_type}/{sensor_id}
            if id, err := strconv.ParseUint(parts[3], 10, 64); err == nil {
                sensorID = id
            }
        }
        if sensorID == 0 {
            log.Printf("MQTT: sensor id not provided (topic %s)", topic)
            return
        }
        if err := c.sensorRepo.Create(sensorID, sp.Value); err != nil {
            log.Printf("MQTT: failed inserting reading: %v", err)
            return
        }
        log.Printf("MQTT: inserted reading sensor=%d value=%f", sensorID, sp.Value)
    }
}
