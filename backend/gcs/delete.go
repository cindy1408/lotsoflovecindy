package gcs

import (
	"cloud.google.com/go/storage"
	"context"
)

func DeleteFileFromGCS(objectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)

	if object.Delete(ctx) != nil {
		return err
	}

	return nil
}
