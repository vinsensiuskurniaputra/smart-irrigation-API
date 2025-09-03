package routers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	coremqtt "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/mqtt"
	irrigationrepo "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/repositories"
	irrigationusecase "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/domain/usecase"
	irrigationhandler "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/presentations/handler"
	"gorm.io/gorm"
)

func Register(rg *gin.RouterGroup, predictionURL string, db *gorm.DB, mqttClient *coremqtt.Client) {
	repo := irrigationrepo.NewPlantRepository(db)
	var nativeClient mqtt.Client
	if mqttClient != nil {
		nativeClient = mqttClient.Native()
	}
	uc := irrigationusecase.NewIrrigationUsecase(predictionURL, repo, db, nativeClient)
	h := irrigationhandler.NewIrrigationHandler(uc)
	grp := rg.Group("/irrigation")
	{
		grp.POST("/predict", h.PredictPlant)
		grp.POST("/devices/:device_id/plant", h.SavePredicted)
		grp.GET("/plants/:plant_id", h.GetPlant)
		grp.PUT("/plants/:plant_id", h.UpdatePlant)
	}
}
