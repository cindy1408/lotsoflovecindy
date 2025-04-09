package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
	"lotsoflovecindy/m/v2/respositories"
)

// HTTP handlers that use the repository functions
func RetrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HERRRREEEE!!!!")
		posts, err := respositories.GetAllPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		log.Println("HERRRREEEE")
		var fileURLs []string
		for _, post := range posts {
			fileURLs = append(fileURLs, post.ContentURL)
		}

		log.Println("file url: ", fileURLs)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(fileURLs)
		if err != nil {
			log.Printf("json.NewEncoder(w).Encode(fileURLs): %v", err)
			http.Error(w, "json.NewEncoder(w).Encode(fileURLs)", http.StatusInternalServerError)
		}
	}
}
