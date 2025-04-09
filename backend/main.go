package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"lotsoflovecindy/m/v2/handlers"
	"lotsoflovecindy/m/v2/postgres"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow necessary methods, including OPTIONS for preflight
		AllowedHeaders: []string{"Content-Type"},
	})

	db, err := postgres.Connection()
	if err != nil {
		return
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/upload", uploadHandler(db))
	http.HandleFunc("/list-files", handlers.RetrieveHandler(db))

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
