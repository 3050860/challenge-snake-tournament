package tenv

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client

func TestMongoDB(t *testing.T) {
	ctx := context.Background()

	// Run the MongoDB container
	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container after the test finishes
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Get the connection string (replica set URL) from the container
	connectionString, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// Connect to MongoDB using the Go driver
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %s", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("failed to disconnect client: %s", err)
		}
	}()

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Errorf("failed to ping database: %s", err)
	}
	MongoClient = client
	// Now you can use 'client' to perform database operations in your tests
}
