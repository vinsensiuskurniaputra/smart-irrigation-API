package seeders

import (
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

func seedSensors(db *gorm.DB) {
	var count int64
	db.Model(&models.Sensor{}).Count(&count)

	if count == 0 {
		sensors := []models.Sensor{
			{
				DeviceID:   1,
				SensorType: "soil_moisture",
			},
			{
				DeviceID:   1,
				SensorType: "temperature",
			},
			{
				DeviceID:   1,
				SensorType: "humidity",
			},
		}

		if err := db.Create(&sensors).Error; err != nil {
			log.Println("❌ Seeder: failed to insert sensors:", err)
			return
		}

		log.Println("✅ Seeder: sensors inserted")
	} else {
		log.Println("ℹ️ Seeder: sensors already exist, skipping")
	}
}
