package gcs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func RetrieveAllFilesFromGCS(w http.ResponseWriter) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return err
	}
	defer client.Close() //nolint:errcheck

	bucket := client.Bucket(bucketName)
	query := &storage.Query{}
	it := bucket.Objects(ctx, query)

	var fileURLs []string

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to list objects: %v", err)
			http.Error(w, "Failed to list objects", http.StatusInternalServerError)
			return err
		}
		fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, attrs.Name)
		fileURLs = append(fileURLs, fileURL)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(fileURLs)
	if err != nil {
		log.Printf("json.NewEncoder(w).Encode(fileURLs): %v", err)
		http.Error(w, "json.NewEncoder(w).Encode(fileURLs)", http.StatusInternalServerError)
	}

	return nil
}
