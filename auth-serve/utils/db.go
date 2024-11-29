package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"TrafiAuth/auth-serve/common"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client


func InitMongoDB() {
    mongoURL := os.Getenv("MONGODB_URL")
    if mongoURL == "" {
        err := fmt.Errorf("MONGODB_URL environment variable not set")
        common.LogError(err, "Environment Variable Error")
        log.Fatal(err)
    }

    var err error
    mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
    if err != nil {
        common.LogError(err, "Failed to connect to MongoDB")
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    if err = mongoClient.Ping(context.Background(), nil); err != nil {
        common.LogError(err, "MongoDB ping failed")
        log.Fatal("MongoDB ping failed:", err)
    }

    fmt.Println("MongoDB connected successfully")
}

func CloseMongoDB() {
    if err := mongoClient.Disconnect(context.Background()); err != nil {
        common.LogError(err, "Error closing MongoDB connection")
        log.Fatal("Error closing MongoDB connection:", err)
    }
    fmt.Println("MongoDB connection closed")
}


func GetMongoCollection() *mongo.Collection {
    dbName := os.Getenv("MONGODB_DB")
    if dbName == "" {
        dbName = "auth"
    }
    return mongoClient.Database(dbName).Collection("users")
}
