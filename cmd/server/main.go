package main

import (
	"log"
	"net/http"
	"os"
	"ticket-system/internal/middleware"
	"ticket-system/internal/repository"
	"ticket-system/internal/service"

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

	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService)

	mux := http.NewServeMux() //create router

	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)

	mux.Handle(
		"/tickets",
		middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.Tickets)),
	)

	mux.Handle(
		"/tickets/",
		middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.Tickets)),
	)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Ticket System API"))
	})

	//log.Println("Server running on port 8080....")
	//err = http.ListenAndServe(":8080", mux)
	//if err != nil {
	//	log.Fatal(err)
	//}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port " + port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
