package presentation

import (
	"github.com/gin-gonic/gin"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	routers "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations/routers"
	"gorm.io/gorm"
)

// RegisterDeviceRoutes registers device related HTTP + websocket + actuator routes.
// mqttClient may be nil; actuator control will skip publishing if nil or disconnected.
func RegisterDeviceRoutes(rg *gin.RouterGroup, db *gorm.DB, mqttClient *coremqtt.Client) {
	routers.Register(rg, db)
	routers.RegisterActuatorRoutes(rg, db, mqttClient)
}
