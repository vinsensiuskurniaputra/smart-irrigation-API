package database

import (
	authSeeders "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/auth/data/seeders"
	devicesSeeders "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/device/data/seeders"
	irrigationSeeders "github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/irrigation/data/seeders"
	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) {
	authSeeders.RunSeeders(db)
	devicesSeeders.RunSeeders(db)
	irrigationSeeders.RunSeeders(db)
}
