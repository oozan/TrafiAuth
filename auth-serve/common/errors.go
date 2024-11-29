package common

import (
	"log"

	"github.com/gofiber/fiber/v2"
)


func HandleError(c *fiber.Ctx, statusCode int, err error, message string) error {
    if err != nil {
        log.Printf("Error: %v", err)
    }
    return c.Status(statusCode).JSON(fiber.Map{"error": message})
}

func LogError(err error, message string) {
    if err != nil {
        log.Printf("%s: %v", message, err)
    }
}
