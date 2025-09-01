package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"size:50;unique;not null"`
	Password     string `gorm:"size:255;not null"`
	Name         string `gorm:"size:100"`
	Email        string `gorm:"size:100;unique;not null"`
	Role         string `gorm:"type:varchar(20);default:'user';check:role IN ('admin','user')" json:"role"`
	PhotoProfile string `gorm:"size:255"`
}
