package database

import (
	"api-undangan/config"
	"api-undangan/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB
func ConnectDB(){
	c := config.Cfg
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}
	err = db.AutoMigrate(&models.Comment{}, &models.Reservation{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database", err)
	}
	seedUsers(db)
	DB = db
	log.Println("Database connect and migrate")
}