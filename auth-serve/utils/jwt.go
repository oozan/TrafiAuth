package utils

import (
	"fmt"
	"os"
	"time"

	"TrafiAuth/auth-serve/common"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_KEY"))
	if len(jwtSecret) == 0 {
		err := fmt.Errorf("JWT_KEY environment variable not set")
		common.LogError(err, "Environment Variable Error")
		panic(err)
	}
}

func GenerateToken(email string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		common.LogError(err, "Error signing token")
		return "", err
	}
	return signedToken, nil
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("Unexpected signing method")
			common.LogError(err, "Token Parsing Error")
			return nil, err
		}
		return jwtSecret, nil
	})
	if err != nil {
		common.LogError(err, "Error parsing token")
		return nil, err
	}
	if !token.Valid {
		err := fmt.Errorf("Invalid token")
		common.LogError(err, "Token Validation Error")
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := fmt.Errorf("Invalid token claims")
		common.LogError(err, "Token Claims Error")
		return nil, err
	}
	return claims, nil
}
