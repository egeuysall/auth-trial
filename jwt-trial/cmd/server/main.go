package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var jwtKey = []byte("secret") // never hardcode in real apps

func createToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := createToken("guldeny")
	w.Write([]byte(token))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "No token", http.StatusUnauthorized)
		return
	}
	tokenStr := auth[len("Bearer "):]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"]
	w.Write([]byte(fmt.Sprintf("Hello, %v!", userID)))
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", protectedHandler)
	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
