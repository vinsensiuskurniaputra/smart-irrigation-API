package presentation

import (
	"github.com/gin-gonic/gin"
	deviceRouters "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations/routers"
	"gorm.io/gorm"
)

func RegisterDeviceRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	deviceRouters.Register(rg, db)
}
