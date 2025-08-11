package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Starwish-Cottage/service/core"
	"github.com/Starwish-Cottage/service/v1/routes"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	client, err := core.InitFirestore()
	if err != nil {
		return
	}
	core.FirestoreClient = client // Initialize the Firestore client globally

	defer client.Close()

	router := gin.Default()
	// router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	routes.SetupRoutes(router)
	router.Run("0.0.0.0:8080")
}
