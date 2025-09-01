package seeders

import "gorm.io/gorm"

func RunSeeders(db *gorm.DB) {
	seedIrrigationRules(db)
}
