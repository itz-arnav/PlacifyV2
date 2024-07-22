package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"placify/backend/src/api"
	"placify/backend/src/storage"
	"placify/backend/src/validate"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	mongoDB := os.Getenv("MONGODB_DB")
	port := os.Getenv("PORT")

	if mongoURI == "" || mongoDB == "" {
		log.Fatalf("MONGODB_URI or MONGODB_DB not set in environment variables")
	}
	if port == "" {
		log.Println("PORT not set in environment variables, defaulting to 8080")
		port = "8080"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	database := client.Database(mongoDB)
	userStorage := storage.NewUserStorage(database)
	userValidator := validate.NewUserValidator()

	router := api.InitializeRouter(userStorage, userValidator)

	serverAddress := ":" + port
	log.Printf("Starting server on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
