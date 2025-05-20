package main

import (
	"fmt"
	"log"
	"net/http"

	"lotsoflovecindy/m/v2/handlers"
	"lotsoflovecindy/m/v2/postgres"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
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

	http.HandleFunc("/upload", handlers.UploadHandler(db))
	http.HandleFunc("/list-files", handlers.RetrieveHandler(db))
	http.HandleFunc("/update-description", handlers.UpdateHandler(db))
	http.HandleFunc("/delete-post", handlers.DeleteHandler(db))

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
