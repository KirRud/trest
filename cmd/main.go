package main

import (
	"log"
	"trest/internal/config"
	"trest/internal/controller"
	"trest/internal/database"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
		return
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Error initializing database:", err)
		return
	}

	handler := controller.NewHandler(db, cfg.SecretKey)
	router := controller.InitRoutes(handler)

	router.Run(":8001")
}
