package database

import (
	"fmt"
	"log"

	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) *gorm.DB {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.Database.Host,
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Name,
        cfg.Database.Port,
        cfg.Database.SSLMode,
    )

    var logLevel logger.LogLevel
    if cfg.Server.Env == "production" {
        logLevel = logger.Silent
    } else {
        logLevel = logger.Info
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })

    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db
    log.Println("Database connected successfully")
    return db
}

func GetDB() *gorm.DB {
    return DB
}