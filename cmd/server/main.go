package main

import (
	"log"
	"net/http"

	"ticket-system/internal/handler"
)

func main() {
	mux := http.NewServeMux() // created router

	mux.HandleFunc("/health", handler.Health)

	log.Println("Server running on port 8080...")

	err:= http.ListenAndServe(":8080", mux)

	if err!=nil{
		log.Fatal(err)
	}
}
