package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func sendJson(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	err := json.NewEncoder(w).Encode(errorResponse)
	if err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func getPassword() string {
	password := os.Getenv("PASSWORD")

	if password == "" {
		log.Fatal("PASSWORD is not set in environment")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}
	return hashedPassword
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		sendError(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	hashedPassword := getPassword()
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))

	if err != nil {
		sendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "Password is correct",
	}
	sendJson(w, response, http.StatusOK)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	r := chi.NewRouter()
	r.Post("/login", handleLogin)

	log.Println("Server starting on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal("Error starting server")
	}
}
