package main

import (
	"log"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/database"
	mqttInfra "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/router"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	deviceusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/domain/usecase"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Connect DB
	db := database.Connect(cfg)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// === Tambahkan Middleware CORS ===
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // boleh semua origin (kalau mau restrict, ganti dengan domain tertentu)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Prepare MQTT (optional)
	var mq *mqttInfra.Client
	if cfg.MQTT.Broker != "" {
		mq = mqttInfra.NewClient(cfg)
	}

	// Register Routes (pass mqtt to device routes)
	router.RegisterRouter(r, db, mq, cfg)

	// Subscribe to sensor topics if connected
	if mq != nil && mq.IsConnected() {
		// Sensor reading consumer
		sensorRepo := devicerepo.NewSensorReadingRepository(db)
		sensorConsumer := deviceusecase.NewSensorReadingConsumer(sensorRepo)
		sensorTopic := "device/+/sensor/+" // narrow to only sensor readings
		if err := mq.Subscribe(sensorTopic, sensorConsumer.Handler()); err != nil {
			log.Printf("Failed subscribe MQTT topic %s: %v", sensorTopic, err)
		}

		// Actuator actual status consumer
		actuatorRepo := devicerepo.NewActuatorRepository(db)
		actuatorStatusConsumer := deviceusecase.NewActuatorStatusConsumer(actuatorRepo)
		actuatorStatusTopic := "device/+/actuator/+/actual-status"
		if err := mq.Subscribe(actuatorStatusTopic, actuatorStatusConsumer.Handler()); err != nil {
			log.Printf("Failed subscribe MQTT topic %s: %v", actuatorStatusTopic, err)
		}

		// Device online/offline status consumer
		deviceRepo := devicerepo.NewDeviceRepository(db)
		deviceStatusConsumer := deviceusecase.NewDeviceStatusConsumer(deviceRepo)
		deviceStatusTopic := "device/+/status"
		if err := mq.Subscribe(deviceStatusTopic, deviceStatusConsumer.Handler()); err != nil {
			log.Printf("Failed subscribe MQTT topic %s: %v", deviceStatusTopic, err)
		}
	}

	log.Printf("Server running on port %s", port)
	r.Run("0.0.0.0:" + port)
}
