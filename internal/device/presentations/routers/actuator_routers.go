package routers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	deviceusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/domain/usecase"
	devicehandler "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations/handler"
	"gorm.io/gorm"
)

// RegisterActuatorRoutes sets up actuator endpoints
func RegisterActuatorRoutes(rg *gin.RouterGroup, db *gorm.DB, mqttClient *coremqtt.Client) {
	repo := devicerepo.NewActuatorRepository(db)
	var rawClient mqtt.Client
	if mqttClient != nil && mqttClient.IsConnected() {
		rawClient = mqttClient.Native()
	}
	uc := deviceusecase.NewActuatorControlUsecase(repo, rawClient)
	h := devicehandler.NewActuatorHandler(uc)
	a := rg.Group("/actuators")
	{
		a.POST(":id/control", h.Control)
	}
}
