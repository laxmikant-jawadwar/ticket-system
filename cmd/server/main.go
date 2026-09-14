package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"ticket-system/internal/database"
	"ticket-system/internal/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close() //when main finishes go close db handle

	log.Println("Database connected successfully")

	mux := http.NewServeMux() //create router

	mux.HandleFunc("/health", handler.Health)

	log.Println("Server running on port 8080....")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
