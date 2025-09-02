package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/database"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
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

	// Register Routes
	router.RegisterRouter(r, db)

	// Init MQTT (if configured) and subscribe to sensor topics
	if cfg.MQTT.Broker != "" {
		mq := mqtt.NewClient(cfg)
		if mq.IsConnected() {
			sensorRepo := devicerepo.NewSensorReadingRepository(db)
			consumer := deviceusecase.NewSensorReadingConsumer(sensorRepo)
			topic := cfg.MQTT.Topic
			if topic == "" {
				topic = "sensors/#" // default wildcard
			}
			if err := mq.Subscribe(topic, consumer.Handler()); err != nil {
				log.Printf("Failed subscribe MQTT topic %s: %v", topic, err)
			}
		}
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
