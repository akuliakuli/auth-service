package main

import (
	"github.com/akuliakuli/auth-service/internal/handlers"
	"github.com/akuliakuli/auth-service/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDB()

	http.HandleFunc("/auth/token", handlers.GenerateTokenHandler)
	http.HandleFunc("/auth/refresh", handlers.RefreshTokenHandler)

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
