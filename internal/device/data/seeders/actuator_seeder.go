package seeders

import (
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	"gorm.io/gorm"
)

func seedActuators(db *gorm.DB) {
	var count int64
	db.Model(&models.Actuator{}).Count(&count)

	if count == 0 {
		actuators := []models.Actuator{
			{
				DeviceID:     1,
				ActuatorName: "Water Pump",
				Type:         "pump",
				PinNumber:   "A1",
				Status:      "off",
			},
		}

		if err := db.Create(&actuators).Error; err != nil {
			log.Println("❌ Seeder: failed to insert actuators:", err)
			return
		}

		log.Println("✅ Seeder: actuators inserted")
	} else {
		log.Println("ℹ️ Seeder: actuators already exist, skipping")
	}
}