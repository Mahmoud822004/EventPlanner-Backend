package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"eventplanner/models"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:Password@123@tcp(127.0.0.1:3306)/eventplanner?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(" Failed to connect to database!")
	}
	fmt.Println("Database connected successfully!")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Event{})	
	DB.AutoMigrate(&models.Invitation{})
}
