package main

import (
	"api-undangan/config"
	"api-undangan/database"
	"api-undangan/middleware"
	"api-undangan/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	database.ConnectDB()
	
	r := gin.Default()
	
	r.Use(middleware.CORSMiddleware())
	
	routes.RegisterRoutes(r)
	
	log.Println("Server running on", config.Cfg.AppPort)
	if err := r.Run(config.Cfg.AppPort); err != nil {
		log.Fatal("Failed to run server", err)
	}
}