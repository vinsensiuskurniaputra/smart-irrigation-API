package presentations

import (
	"github.com/gin-gonic/gin"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	irrigationRouters "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/presentations/routers"
	"gorm.io/gorm"
)

// RegisterIrrigationRoutes registers irrigation related routes (currently prediction).
func RegisterIrrigationRoutes(rg *gin.RouterGroup, predictionURL string, db *gorm.DB, mqttClient *coremqtt.Client) {
	irrigationRouters.Register(rg, predictionURL, db, mqttClient)
}
