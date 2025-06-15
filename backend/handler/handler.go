package handler

import (
	"encoding/json"
	"fmt"
	"gallery/gcs"
	"gallery/models"
	"gallery/respositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func RetrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("🔥 RetrieveHandler called")
		posts, err := respositories.GetAllPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		posts, err = gcs.RetrieveAllPosts(posts)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("json.NewEncoder(w).Encode(fileURLs): %v", err)
			http.Error(w, "json.NewEncoder(w).Encode(fileURLs)", http.StatusInternalServerError)
		}
	}
}

// UploadHandler which accepts the db connection
func UploadHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request at /upload (for signed URL)")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			http.Error(w, "Missing filename", http.StatusBadRequest)
			return
		}

		contentType := r.FormValue("contentType")
		if contentType == "" {
			contentType = "application/octet-stream" // fallback
		}

		signedURL, err := gcs.UploadPost(filename, contentType)
		if err != nil {
			http.Error(w, "Failed to generate signed URL", http.StatusInternalServerError)
		}

		// Respond with the signed URL
		response := map[string]string{
			"signedUrl": signedURL,
			"publicUrl": fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcs.BucketName, filename),
		}

		publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", gcs.BucketName, filename)

		err = respositories.CreatePost(db, &models.Post{
			ContentURL: publicURL,
		})
		if err != nil {
			log.Printf("Failed to create post in database: %v", err)
			http.Error(w, "Failed to create post in database", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateHandler(db *gorm.DB) http.HandlerFunc {
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
		post, err := respositories.GetPostById(db, uuid)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Update the post fields
		post.Description = description

		// Save the updated post
		err = respositories.UpdatePost(db, post)
		if err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// Respond with success
		fmt.Fprintf(w, "Post updated successfully! URL: %s", post.ContentURL)
	}
}

func DeleteHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received DELETE request at /delete-post")

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		url := r.FormValue("url_path")
		id := r.FormValue("id")

		log.Println("Received url_path:", url)
		log.Println("Received id:", id)

		if url == "" {
			log.Println("Missing post URL")
			http.Error(w, "Post URL is required", http.StatusBadRequest)
			return
		}

		if id == "" {
			log.Println("Missing post ID")
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		uuid, err := uuid.Parse(id)
		if err != nil {
			log.Println("Error parsing UUID:", err)
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}

		decoded, err := gcs.ExtractObjectName(url)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
		}

		err = gcs.DeletePost(decoded)
		if err != nil {
			http.Error(w, "Failed to delete from storage", http.StatusInternalServerError)
		}

		err = respositories.DeletePost(db, uuid)
		if err != nil {
			http.Error(w, "Failed to delete post from database", http.StatusInternalServerError)
		}

		log.Println("Successfully deleted from database")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Delete successful"}`))
		log.Println("Delete operation completed successfully")
	}
}
