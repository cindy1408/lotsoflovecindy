package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
	"lotsoflovecindy/m/v2/respositories"
)

// RetrieveHandler HTTP handlers that use the repository functions
func RetrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := respositories.GetAllPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}

		//var fileURLs []string
		//for _, post := range posts {
		//	fileURLs = append(fileURLs, post.ContentURL)
		//}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("json.NewEncoder(w).Encode(fileURLs): %v", err)
			http.Error(w, "json.NewEncoder(w).Encode(fileURLs)", http.StatusInternalServerError)
		}
	}
}
