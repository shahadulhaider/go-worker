package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	var wg sync.WaitGroup

	collection := client.Database(dbName).Collection(collectionName)

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := collection.CountDocuments(context.Background(), bson.D{})
		if err != nil {
			log.Println("Error counting users:", err)
			return
		}

		log.Printf("Summary Statistics - Total Users: %d\n", count)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logChangesPeriodically(collection)
	}()

	wg.Wait()
}

func logChangesPeriodically(collection *mongo.Collection) {
	lastCount, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		log.Println("Error counting users:", err)
		return
	}

	for {
		time.Sleep(1 * time.Minute) // Adjust the interval as needed

		currentCount, err := collection.CountDocuments(context.Background(), bson.D{})
		if err != nil {
			log.Println("Error counting users:", err)
			continue
		}

		if currentCount != lastCount {
			newUsers, err := getNewUsers(collection, lastCount)
			if err != nil {
				log.Println("Error fetching new users:", err)
				continue
			}

			log.Printf("User count changed. Previous Count: %d, Current Count: %d\n", lastCount, currentCount)
			logNewUsers(newUsers)
			lastCount = currentCount
		}
	}
}

func getNewUsers(collection *mongo.Collection, lastCount int64) ([]bson.M, error) {
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var newUsers []bson.M
	for cursor.Next(context.Background()) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		newUsers = append(newUsers, user)
	}

	// Only return the new users since the last count
	return newUsers[int(lastCount):], nil
}

func logNewUsers(users []bson.M) {
	for _, user := range users {
		log.Println("New user added:", user)
	}
}
