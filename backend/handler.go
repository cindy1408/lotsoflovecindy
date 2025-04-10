package main

import (
	"encoding/json"
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

		// Respond with the post ID and file URL (send back to UI)
		response := map[string]string{
			"id":      post.ID.String(),
			"fileURL": fileURL,
		}
		responseData, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(responseData)
		if err != nil {
			http.Error(w, "Failed write response", http.StatusInternalServerError)
		}
	}
}

func updateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at /update-description")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Get the ID of the post to update from the form data
		postID := r.FormValue("id")
		if postID == "" {
			log.Println("HERE")
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		// Get the updated description from the form data
		description := r.FormValue("description")
		if description == "" {
			log.Println("HERE 1")
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		// Get the file URL from the form data
		fileURL := r.FormValue("fileURL")
		if fileURL == "" {
			log.Println("HERE 2")
			http.Error(w, "File URL is required", http.StatusBadRequest)
			return
		}

		// Find the post by ID
		var post models.Post
		if err := db.First(&post, "id = ?", postID).Error; err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Update the post fields
		post.ContentURL = fileURL
		post.Description = description
		post.Schedule = time.Now()

		// Save the updated post
		if err := db.Save(&post).Error; err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// Respond with success
		if _, err := fmt.Fprintf(w, "Post updated successfully! URL: %s", fileURL); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	}
}
