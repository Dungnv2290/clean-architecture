package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoHandler define the MongoDb handler
type MongoHandler struct {
	db     *mongo.Database
	client *mongo.Client
}

// NewMongoHandler create new MongoHander
func NewMongoHandler() *MongoHandler {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoHandler{
		db:     client.Database(os.Getenv("MONGODB_DATABASE")),
		client: client,
	}
}

// Client return the client property
func (m *MongoHandler) Client() *mongo.Client {
	return m.client
}

// DB returns the database property
func (m *MongoHandler) Db() *mongo.Database {
	return m.db
}
