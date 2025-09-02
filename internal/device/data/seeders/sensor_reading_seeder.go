package seeders

import (
	"log"
	"time"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

func seedSensorReading(db *gorm.DB) {
	var count int64
	db.Model(&models.SensorReading{}).Count(&count)

	if count == 0 {
		readings := []models.SensorReading{
			{
				SensorID:  1,
				Value:    45.5,
				RecordedAt: time.Now(),
			},
			{
				SensorID:  1,
				Value:    22.3,
				RecordedAt: time.Now(),
			},
			{
				SensorID:  1,
				Value:    60.0,
				RecordedAt: time.Now(),
			},
		}

		if err := db.Create(&readings).Error; err != nil {
			log.Println("❌ Seeder: failed to insert sensor readings:", err)
			return
		}

		log.Println("✅ Seeder: sensor readings inserted")
	} else {
		log.Println("ℹ️ Seeder: sensor readings already exist, skipping")
	}
}
