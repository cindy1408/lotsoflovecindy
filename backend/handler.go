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

func retrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🔥 RetrieveHandler called")
		posts, err := respositories.GetAllPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("json.NewEncoder(w).Encode(fileURLs): %v", err)
			http.Error(w, "json.NewEncoder(w).Encode(fileURLs)", http.StatusInternalServerError)
		}
	}
}

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
		fmt.Println("Update description endpoint hit!")

		// Ensure method is POST
		if r.Method != http.MethodPost {
			log.Println("Invalid request method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Get the URL of the post to update from the form data
		id := r.FormValue("id")
		if id == "" {
			log.Println("Missing post ID")
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		uuid, err := uuid.Parse(id)
		if err != nil {
			log.Print("uuid.Parse")
		}

		// Get the URL of the post to update from the form data
		url := r.FormValue("url_path")
		if url == "" {
			log.Println("Missing post URL")
			http.Error(w, "Post URL is required", http.StatusBadRequest)
			return
		}

		// Get the updated description from the form data
		description := r.FormValue("description")
		if description == "" {
			log.Println("Missing updated description")
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		// Find the post by URL (assuming postUrl is a unique identifier)
		var post models.Post
		if err := db.First(&post, "id = ?", uuid).Error; err != nil {
			log.Println("Post not found:", err)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Update the post fields
		post.Description = description
		post.Schedule = time.Now()

		// Save the updated post
		if err := db.Save(&post).Error; err != nil {
			log.Println("Failed to update post:", err)
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// Respond with success
		fmt.Fprintf(w, "Post updated successfully! URL: %s", post.ContentURL)
	}
}
