package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	"github.com/Starwish-Cottage/service/routes"
)

func initFirestore() (*firestore.Client, error) {
	ctx := context.Background()

	credsPath := os.Getenv("GOOGLE_FIREBASE_CREDENTIALS")
	projId := os.Getenv("PROJECT_ID")
	config := &firebase.Config{ProjectID: projId}
	opt := option.WithCredentialsFile(credsPath)

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v\n", err)
	}
	return client, nil
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	router := gin.Default()
	client, err := initFirestore()
	if err != nil {
		return
	}

	defer client.Close()
	routes.SetupRoutes(router)
	router.Run(":8080")
}
