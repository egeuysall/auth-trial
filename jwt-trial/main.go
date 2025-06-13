package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var jwtKey []byte

// Define your custom claims struct
type userClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
}

func createToken(userId string) (string, error) {
	claims := userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func verifyToken(tokenStr string) (*userClaims, error) {
	claims := &userClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables:", err)
	}

	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatal("JWT_KEY not set in environment")
	}
	jwtKey = []byte(key)

	token, err := createToken("56EC6CD9-5F20-42A9-90DD-B40AC6378CD7")
	if err != nil {
		fmt.Println("Error creating JWT token:", err)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		fmt.Println("Token could not be verified:", err)
		return
	}

	fmt.Println("Token is verified for user:", claims.UserId)
}
