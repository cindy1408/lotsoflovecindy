package gcs

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

const bucketName = "lotsoflovecindy"

func UploadFileToGCS(w http.ResponseWriter, file io.Reader, fileName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create GCS client: %v", err)
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return "", err
	}
	defer client.Close() //nolint:errcheck

	bucket := client.Bucket(bucketName)
	object := bucket.Object(fileName)
	writer := object.NewWriter(ctx)

	_, err = io.Copy(writer, file)
	if err != nil {
		http.Error(w, "Failed to upload file to GCS", http.StatusInternalServerError)
		return "", err
	}

	if err = writer.Close(); err != nil {
		log.Printf("Failed to close GCS writer: %v", err)
		http.Error(w, fmt.Sprintf("Failed to finalize upload to GCS: %v", err), http.StatusInternalServerError)
		return "", err
	}

	// Return public URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return publicURL, nil
}
