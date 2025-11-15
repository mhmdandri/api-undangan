package main

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	database.ConnectDB()
	
	r := gin.Default()
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://wedding.mohaproject.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	routes.RegisterRoutes(r)
	
	log.Println("Server running on", config.Cfg.AppPort)
	if err := r.Run(config.Cfg.AppPort); err != nil {
		log.Fatal("Failed to run server", err)
	}
}