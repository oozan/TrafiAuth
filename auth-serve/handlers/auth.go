package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"TrafiAuth/auth-serve/common"
	"TrafiAuth/auth-serve/models"
	"TrafiAuth/auth-serve/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)


func RegisterHandler(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return common.HandleError(c, http.StatusBadRequest, err, "Invalid request body")
    }

    collection := utils.GetMongoCollection()
    ctx := context.Background()

    count, err := collection.CountDocuments(ctx, bson.M{"email": user.Email})
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Database error")
    }
    if count > 0 {
        return common.HandleError(c, http.StatusConflict, nil, "User already exists")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error hashing password")
    }
    user.Password = string(hashedPassword)

    _, err = collection.InsertOne(ctx, user)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error creating user")
    }

    return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Registration successful"})
}


func LoginHandler(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return common.HandleError(c, http.StatusBadRequest, err, "Invalid request body")
    }

    collection := utils.GetMongoCollection()
    ctx := context.Background()

    var storedUser models.User
    err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&storedUser)
    if err != nil {
        return common.HandleError(c, http.StatusUnauthorized, err, "User not found")
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
    if err != nil {
        return common.HandleError(c, http.StatusUnauthorized, err, "Incorrect password")
    }

    accessToken, err := utils.GenerateToken(user.Email, 15*time.Minute)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error generating access token")
    }

    refreshToken, err := utils.GenerateToken(user.Email, 24*time.Hour)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error generating refresh token")
    }

    err = utils.StoreRefreshToken(user.Email, refreshToken)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error storing refresh token")
    }

    utils.SetAuthCookies(c, accessToken, refreshToken)

    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}


func ValidateHandler(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return common.HandleError(c, http.StatusUnauthorized, nil, "Missing Authorization header")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return common.HandleError(c, http.StatusUnauthorized, nil, "Invalid Authorization header")
    }

    claims, err := utils.ParseToken(parts[1])
    if err != nil {
        return common.HandleError(c, http.StatusUnauthorized, err, "Invalid token")
    }

    email, ok := claims["email"].(string)
    if !ok {
        return common.HandleError(c, http.StatusUnauthorized, nil, "Invalid token claims")
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{"email": email, "message": "Token is valid"})
}


func RefreshHandler(c *fiber.Ctx) error {
    refreshToken := c.Cookies("refresh_token")
    if refreshToken == "" {
        return common.HandleError(c, http.StatusUnauthorized, nil, "Missing refresh token")
    }

    claims, err := utils.ParseToken(refreshToken)
    if err != nil {
        return common.HandleError(c, http.StatusUnauthorized, err, "Invalid refresh token")
    }

    email, ok := claims["email"].(string)
    if !ok {
        return common.HandleError(c, http.StatusUnauthorized, nil, "Invalid token claims")
    }

    storedToken, err := utils.GetStoredRefreshToken(email)
    if err != nil || storedToken != refreshToken {
        return common.HandleError(c, http.StatusUnauthorized, err, "Refresh token mismatch")
    }

    newAccessToken, err := utils.GenerateToken(email, 15*time.Minute)
    if err != nil {
        return common.HandleError(c, http.StatusInternalServerError, err, "Error generating new access token")
    }

    utils.SetAccessTokenCookie(c, newAccessToken)

    return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Access token refreshed"})
}
