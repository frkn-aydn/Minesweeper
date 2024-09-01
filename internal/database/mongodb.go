package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB encapsulates the MongoDB client and database name.
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// Connect establishes a connection to the MongoDB database.
func Connect() (*MongoDB, error) {
	// TODO: Currently hardcoding the connection string. This should be read from a configuration file.
	const mongodbConnectionUrl = "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(mongodbConnectionUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &MongoDB{
		client:   client,
		database: client.Database("minesweeper"),
	}, nil
}

// DB returns the MongoDB database instance.
func (db *MongoDB) DB() *mongo.Database {
	return db.database
}

// Disconnect closes the connection to MongoDB.
func (db *MongoDB) Disconnect() {
	if err := db.client.Disconnect(context.Background()); err != nil {
		fmt.Printf("failed to disconnect from MongoDB: %v\n", err)
	}

	fmt.Println("disconnected from MongoDB")
}

// Collection returns a MongoDB collection instance.
func (db *MongoDB) collection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

func (db *MongoDB) GameCollection() *mongo.Collection {
	return db.collection("game")
}
