package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"lotsoflovecindy/m/v2/postgres"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow necessary methods, including OPTIONS for preflight
		AllowedHeaders: []string{"Content-Type"},
	})

	log.Println("HEREEEEEEE")
	db, err := postgres.Connection()
	if err != nil {
		return
	}

	http.HandleFunc("/upload", uploadHandler(db))
	http.HandleFunc("/list-files", retrieveHandler(db))
	http.HandleFunc("/update-description", updateHandler(db))

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
