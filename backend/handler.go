package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"lotsoflovecindy/m/v2/gcs"
	"lotsoflovecindy/m/v2/models"
	"lotsoflovecindy/m/v2/respositories"
)

// Upload handler which accepts the db connection
func uploadHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Upload file to GCS and get the file URL
		fileURL, err := gcs.UploadFileToGCS(w, file, header.Filename)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			return
		}

		// Create new post in the database with the file URL
		post := &models.Post{
			ID:          uuid.New(),
			ContentURL:  fileURL,
			Description: r.FormValue("description"), // optional: capture the description from form
			Schedule:    time.Now(),                 // or capture from form
			DateCreated: time.Now(),
		}

		if err := respositories.CreatePost(db, post); err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// Respond with success
		if _, err := fmt.Fprintf(w, "File uploaded and post created successfully! URL: %s", fileURL); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}
