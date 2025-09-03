package seeders

import (
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/models"
	"gorm.io/gorm"
)

func seedIrrigationRules(db *gorm.DB) {
	var count int64
	db.Model(&models.IrrigationRule{}).Count(&count)

	if count == 0 {
		rules := []models.IrrigationRule{
			{
				PlantName:         "Tomato",
				MinMoisture:       60,
				MaxMoisture:       80,
				PreferredTemp:     20, // °C
				PreferredHumidity: 60, // %
			},
			{
				PlantName:         "Cactus",
				MinMoisture:       10,
				MaxMoisture:       30,
				PreferredTemp:     25, // °C
				PreferredHumidity: 30, // %
			},
			{
				PlantName:         "Spinach",
				MinMoisture:       50,
				MaxMoisture:       70,
				PreferredTemp:     18, // °C
				PreferredHumidity: 65, // %
			},
			{
				PlantName:         "Monstera",
				MinMoisture:       50,
				MaxMoisture:       70,
				PreferredTemp:     24, // °C
				PreferredHumidity: 60, // %
			},
			{
				PlantName:         "Chili",
				MinMoisture:       60,
				MaxMoisture:       80,
				PreferredTemp:     17, // °C
				PreferredHumidity: 70, // %
			},
		}

		if err := db.Create(&rules).Error; err != nil {
			log.Println("❌ Seeder: failed to insert irrigation rules:", err)
			return
		}

		log.Println("✅ Seeder: irrigation rules inserted")
	} else {
		log.Println("ℹ️ Seeder: irrigation rules already exist, skipping")
	}
}
