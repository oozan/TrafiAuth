package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"TrafiAuth/auth-serve/common"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		err := fmt.Errorf("REDIS_ADDR environment variable not set")
		common.LogError(err, "Environment Variable Error")
		log.Fatal(err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		common.LogError(err, "Failed to connect to Redis")
		log.Fatal("Failed to connect to Redis:", err)
	}

	fmt.Println("Redis connected successfully")
}

func CloseRedis() {
	if err := redisClient.Close(); err != nil {
		common.LogError(err, "Error closing Redis connection")
		log.Fatal("Error closing Redis connection:", err)
	}
	fmt.Println("Redis connection closed")
}

func StoreRefreshToken(email, token string) error {
	err := redisClient.Set(context.Background(), email+"_refresh", token, 24*time.Hour).Err()
	if err != nil {
		common.LogError(err, "Error storing refresh token")
		return err
	}
	return nil
}

func GetStoredRefreshToken(email string) (string, error) {
	result, err := redisClient.Get(context.Background(), email+"_refresh").Result()
	if err != nil {
		common.LogError(err, "Error getting stored refresh token")
		return "", err
	}
	return result, nil
}
