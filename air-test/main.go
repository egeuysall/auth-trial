package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var notes = []Note{}
var idCounter = 1

func main() {
	http.HandleFunc("/notes", notesHandler)

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	case "POST":
		var n Note
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		n.ID = idCounter
		idCounter++
		notes = append(notes, n)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(n)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
