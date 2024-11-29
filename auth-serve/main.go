// main.go

package main

import (
	"log"
	"os"

	"TrafiAuth/auth-serve/handlers"
	"TrafiAuth/auth-serve/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using environment variables")
    }

    app := fiber.New()

    // Initialize database and redis connections
    utils.InitMongoDB()
    defer utils.CloseMongoDB()

    utils.InitRedis()
    defer utils.CloseRedis()

    // Set up routes
    app.Post("/register", handlers.RegisterHandler)
    app.Post("/login", handlers.LoginHandler)
    app.Get("/validate", handlers.ValidateHandler)
    app.Post("/refresh", handlers.RefreshHandler)

    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("PORT environment variable not set")
    }

    log.Fatal(app.Listen(":" + port))
}
