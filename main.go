package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"wadood/routes"
	"wadood/handlers"
)

func main() {
	err := handlers.InitDB()
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