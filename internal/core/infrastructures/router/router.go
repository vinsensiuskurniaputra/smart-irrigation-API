package router

import (
	"github.com/gin-gonic/gin"
	authRoute "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/presentations"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	middleware "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/middlewares"
	devicePresentation "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations"
	"gorm.io/gorm"
)

func RegisterRouter(r *gin.Engine, db *gorm.DB, mqttClient *coremqtt.Client) {

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
		}
	}
}
