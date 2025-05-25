package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"wadood/db"
	"wadood/routes"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load .env file")
	}

	err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	router := routes.Setup()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}

	fmt.Println("Server running on port 8080")
	log.Fatal(srv.ListenAndServe())
}
