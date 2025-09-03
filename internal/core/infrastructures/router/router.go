package router

import (
	"github.com/gin-gonic/gin"
	authRoute "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/presentations"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	middleware "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/middlewares"
	devicePresentation "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations"
	irrigationPresentation "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/presentations"
	"gorm.io/gorm"
)

func RegisterRouter(r *gin.Engine, db *gorm.DB, mqttClient *coremqtt.Client, cfg *config.Config) {

	authMiddleware := middleware.AuthMiddleware()

	router := r.Group("/api/v1")
	{
		public := router.Group("/")
		{
			authRoute.RegisterPublicRoutes(public, db)
		}

		auth := router.Group("/")
		auth.Use(authMiddleware)
		{
			authRoute.RegisterProtectedRoutes(auth, db)
			devicePresentation.RegisterDeviceRoutes(auth, db, mqttClient)
			irrigationPresentation.RegisterIrrigationRoutes(auth, cfg.ML.PredictionURL, db, mqttClient)
		}
	}
}
