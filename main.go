package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set.")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set.")
	}

	collectionName := os.Getenv("COLLECTION_NAME")
	if collectionName == "" {
		log.Fatal("COLLECTION_NAME environment variable is not set.")
	}

	cronInterval := os.Getenv("CRON_INTERVAL")
	if cronInterval == "" {
		cronInterval = "@every 1m" 
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Create a cron scheduler
	c := cron.New()

	// Add a cron job to run the worker function at regular intervals
	_, err = c.AddFunc(cronInterval, func() {
		processAndLogSummary(client, dbName, collectionName)
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start the cron scheduler in a separate goroutine
	go c.Start()

	// Wait for a signal to gracefully shut down the cron scheduler
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the cron scheduler gracefully
	c.Stop()
}

func processAndLogSummary(client *mongo.Client, dbName, collectionName string) {
	collection := client.Database(dbName).Collection(collectionName)

	// Count the number of users in the database
	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error counting users:", err)
		return
	}

	log.Printf("Summary Statistics - Total Users: %d\n", count)
}
