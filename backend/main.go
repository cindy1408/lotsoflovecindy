package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"lotsoflovecindy/m/v2/backend/gcs"
)

func main() {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"}, // Allow necessary methods, including OPTIONS for preflight
		AllowedHeaders: []string{"Content-Type"},
	})

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/list-files", func(w http.ResponseWriter, r *http.Request) {
		err := gcs.RetrieveAllFilesFromGCS(w)
		if err != nil {
			log.Printf("Error retrieving files: %v", err)
		}
	})

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux)))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /upload")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close() //nolint:errcheck

	log.Printf("File received: %s", header.Filename)

	err = gcs.UploadFileToGCS(w, file, header.Filename)
	if err != nil {
		return
	}

	if _, err := fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
