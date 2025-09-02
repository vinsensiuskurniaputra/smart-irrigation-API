package seeders

import (
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

func seedDevice(db *gorm.DB) {
	var count int64
	db.Model(&models.Device{}).Count(&count)

	if count == 0 {
		devices := []models.Device{
			{
				UserID:     1,
				DeviceName: "Greenhouse",
				DeviceCode: "GH-001",
				Status:     "online",
			},
		}

		if err := db.Create(&devices).Error; err != nil {
			log.Println("❌ Seeder: failed to insert devices:", err)
			return
		}

		log.Println("✅ Seeder: devices inserted")
	} else {
		log.Println("ℹ️ Seeder: devices already exist, skipping")
	}
}
