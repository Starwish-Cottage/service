package core

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirestoreClient *firestore.Client

func InitFirestore() (*firestore.Client, error) {
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
