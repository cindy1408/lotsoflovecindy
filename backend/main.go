package main

import (
	"fmt"
	"gallery/handler"
	"log"
	"net/http"

	"gallery/postgres"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	db, err := postgres.Connection()
	if err != nil {
		return
	}

	http.HandleFunc("/list-files", handler.RetrieveHandler(db))
	http.HandleFunc("/upload", handler.UploadHandler(db))
	http.HandleFunc("/update-description", handler.UpdateHandler(db))
	http.HandleFunc("/delete-post", handler.DeleteHandler(db))

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}
