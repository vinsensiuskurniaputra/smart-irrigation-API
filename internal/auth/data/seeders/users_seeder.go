package seeders

import (
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		// bikin hash password
		adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		userPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

		users := []models.User{
			{
				Username: "admin",
				Password: string(adminPass),
				Name:     "Administrator",
				Email:    "admin@example.com",
				Role:     "admin",
			},
			{
				Username: "johndoe",
				Password: string(userPass),
				Name:     "John Doe",
				Email:    "johndoe@example.com",
				Role:     "user",
			},
		}

		if err := db.Create(&users).Error; err != nil {
			log.Println("❌ Seeder: failed to insert users:", err)
			return
		}

		log.Println("✅ Seeder: users inserted")
	} else {
		log.Println("ℹ️ Seeder: users already exist, skipping")
	}
}
