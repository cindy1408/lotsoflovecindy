package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"lotsoflovecindy/m/v2/gcs"
	"lotsoflovecindy/m/v2/models"
	"lotsoflovecindy/m/v2/respositories"
	"net/http"
)

func RetrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🔥 RetrieveHandler called")
		posts, err := respositories.GetAllPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		const bucketName = "lotsoflovecindy"
		gcsBaseURL := "https://storage.googleapis.com/" + bucketName + "/"

		for i, post := range posts {
			if len(post.ContentURL) > len(gcsBaseURL) && post.ContentURL[:len(gcsBaseURL)] == gcsBaseURL {
				objectName := post.ContentURL[len(gcsBaseURL):]
				signedURL, err := gcs.GenerateSignedURL(bucketName, objectName)
				if err != nil {
					log.Printf("Failed to sign URL for object %s: %v", objectName, err)
					continue // fallback: leave original URL
				}
				posts[i].ContentURL = signedURL
			}
		}

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

		const bucketName = "lotsoflovecindy"

		// Generate signed URL for PUT
		signedURL, err := gcs.GenerateUploadSignedUploadURL(bucketName, filename, contentType)
		if err != nil {
			log.Printf("Failed to generate signed URL: %v", err)
			http.Error(w, "Failed to generate signed URL", http.StatusInternalServerError)
			return
		}

		// Respond with the signed URL
		response := map[string]string{
			"signedUrl": signedURL,
			"publicUrl": fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename),
		}

		publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, filename)

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
		var post models.Post
		if err := db.First(&post, "id = ?", uuid).Error; err != nil {
			log.Println("Post not found:", err)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		// Update the post fields
		post.Description = description

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
