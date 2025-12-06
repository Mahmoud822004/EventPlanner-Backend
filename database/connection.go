package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"eventplanner/models"
)

var DB *gorm.DB

func Connect() {
	// Use environment variables from docker-compose
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	fmt.Println("Database connected successfully!")

	// Auto-migrate tables
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Event{})
	DB.AutoMigrate(&models.Invitation{})
}
