package router

import (
	"github.com/gin-gonic/gin"
	authRoute "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/presentations"
	middleware "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/middlewares"
	"gorm.io/gorm"
)

func RegisterRouter(r *gin.Engine, db *gorm.DB) {

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
		}
	}
}
