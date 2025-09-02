package routers

import (
	"github.com/gin-gonic/gin"
	devicerepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/repositories"
	deviceusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/domain/usecase"
	devicehandler "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/presentations/handler"
	"gorm.io/gorm"
)

func Register(rg *gin.RouterGroup, db *gorm.DB) {
	repo := devicerepo.NewDeviceRepository(db)
	uc := deviceusecase.NewDeviceUsecase(repo)
	h := devicehandler.NewDeviceHandler(uc)
	liveHandler := devicehandler.NewLiveDataHandler(db)

	devices := rg.Group("/devices")
	{
		devices.GET("/", h.List)
		devices.GET(":id", h.Detail)
		devices.GET(":id/live", liveHandler.DeviceLive)
	}
}
