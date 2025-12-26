package main

import (
	"eventplanner/database"
	"eventplanner/models"
	"eventplanner/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	database.DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Invitation{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://eventplanner-frontend-khaledahmed04-dev.apps.rm2.thpm.p1.openshiftapps.com"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.AuthRoutes(r)
	routes.RegisterRoutes(r)

	r.Run(":8081") // Backend runs at localhost:8080
}
