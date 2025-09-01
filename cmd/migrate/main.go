package main

import (
	"log"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/database"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Connect DB
	db := database.Connect(cfg)

	// Run AutoMigrate / Seeder
	database.Migrate(db)
	database.RunSeeders(db)

	log.Printf("Migration and seeding completed successfully")
}
