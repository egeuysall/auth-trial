package main

import (
	"github.com/egeuysall/auth-trial/api"
	"github.com/egeuysall/auth-trial/supabase"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	dbConn := supabase.Connect()
	defer dbConn.Close()

	// Assign db connection
	supabase.Db = dbConn

	// Define the router
	router := api.NewRouter()

	// Start the server
	log.Printf("Server starting on http://localhost:8080")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal(err)
	}
}
