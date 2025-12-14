package database

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
    Client *mongo.Client
    DB     *mongo.Database
}

func ConnectMongoDB(uri, dbName string) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout for Atlas
    defer cancel()
    
    // Create client options with Server API for MongoDB Atlas
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    clientOptions := options.Client().
        ApplyURI(uri).
        SetServerAPIOptions(serverAPI).
        SetMaxPoolSize(10).
        SetMinPoolSize(5)
    
    // Connect to MongoDB
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }
    
    // Ping the database with Server API for MongoDB Atlas
    var result bson.M
    if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
    }
    
    log.Println("✅ Connected to MongoDB Atlas successfully!")
    
    return &MongoDB{
        Client: client,
        DB:     client.Database(dbName),
    }, nil
}

func (m *MongoDB) Disconnect() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    return m.Client.Disconnect(ctx)
}