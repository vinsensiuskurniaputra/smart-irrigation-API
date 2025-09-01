package database

import (
	"log"

	authModels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/models"
	devicesModels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/models"
	irrigationModels "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	// Daftar semua model
	models := []interface{}{
		&authModels.User{},
		&devicesModels.Device{},
		&devicesModels.Sensor{},
		&devicesModels.SensorReading{},
		&devicesModels.Actuator{},
		&devicesModels.ActuatorLog{},
		&irrigationModels.Plant{},
		&irrigationModels.IrrigationRule{},
	}

	// Drop semua tabel
	if err := db.Migrator().DropTable(models...); err != nil {
		log.Fatalf("❌ Failed to drop tables: %v", err)
	}
	log.Println("🗑️ Semua tabel berhasil dihapus")

	// Migrasi ulang
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	log.Println("✅ Database migration completed (reset & migrate)")
}
