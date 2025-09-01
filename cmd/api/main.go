package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/config"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/database"
	"github.com/vinsensiuskurniaputra/smart-irrigation-API/internal/core/infrastructures/router"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Connect DB
	db := database.Connect(cfg)
	
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}


	r := gin.Default()

	// Register Routes
	router.RegisterRouter(r, db)

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
