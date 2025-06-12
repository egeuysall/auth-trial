package api

import (
	"encoding/json"
	"github.com/egeuysall/auth-trial/supabase"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Welcome to Auth Trial v1. Available routes: GET /, GET /ping, POST /v1/login, POST /v1/users, GET /v1/users/{id}, DELETE /v1/users/{id}, PATCH /v1/users/{id}",
	}
	SendJson(w, response, http.StatusOK)
}

func CheckPing(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "pong",
	}
	SendJson(w, response, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input supabase.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	foundUser, err := supabase.GetUser(ctx, supabase.User{Email: input.Email})
	if err != nil {
		if err == pgx.ErrNoRows {
			SendError(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			SendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if !CheckPassword(input.Password, foundUser.Password) {
		SendError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"message": "Login successful",
	}
	SendJson(w, response, http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input supabase.User

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		SendError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	input.Password = hashedPassword
	err = supabase.CreateUser(ctx, input)

	if err != nil {
		SendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "User created",
	}

	SendJson(w, response, http.StatusCreated)
}
