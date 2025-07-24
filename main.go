package main

import (
	"log"

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

	router := gin.Default()
	client, err := core.InitFirestore()
	if err != nil {
		return
	}
	core.FirestoreClient = client // Initialize the Firestore client globally

	defer client.Close()
	routes.SetupRoutes(router)
	router.Run(":8080")
}
