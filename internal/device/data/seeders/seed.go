package seeders

import "gorm.io/gorm"


func RunSeeders(db *gorm.DB) {
	seedDevice(db)
	seedActuators(db)
	seedSensors(db)
}
