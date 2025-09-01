package presentations

import (
	"github.com/gin-gonic/gin"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/repositories"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/domain/usecase"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/presentations/handler"
	"gorm.io/gorm"
)

func RegisterPublicRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	authUC := usecase.NewAuthUsecase(userRepo)
	h := handler.NewAuthHandler(authUC)

	auth := rg.Group("/")
	{
		auth.GET("/ping", handler.PingHandler)
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}

func RegisterProtectedRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepository)
	h := handler.NewAuthHandler(authUsecase)

	userGroup := rg.Group("users")
	{
		userGroup.GET("/me", h.GetMe)
		userGroup.PUT("/change-password", h.UpdateMyPassword)
	}
}
