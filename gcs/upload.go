package gcs

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	"cloud.google.com/go/storage"
// )

// const (
// 	bucketName = "lotsoflovecindy" // Change this to your bucket
// )

// func uploadFile(w http.ResponseWriter, r *http.Request) {
// 	ctx := context.Background()

// 	// Initialize a GCS client
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		http.Error(w, "Failed to create client", http.StatusInternalServerError)
// 		return
// 	}
// 	defer client.Close()

// 	// Parse the form to get the file
// 	r.ParseMultipartForm(10 << 20)
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		http.Error(w, "File not received", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Upload file to GCS
// 	object := client.Bucket(bucketName).Object(handler.Filename)
// 	wc := object.NewWriter(ctx)
// 	if _, err := io.Copy(wc, file); err != nil {
// 		http.Error(w, "Failed to upload", http.StatusInternalServerError)
// 		return
// 	}
// 	if err := wc.Close(); err != nil {
// 		http.Error(w, "Failed to close writer", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
// }
