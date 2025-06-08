package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"os"
	"time"
)

type serviceAccountJSON struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

const serviceAccount = "../../credentials.json"

func GenerateSignedURL(bucket, object string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	credsJSON, err := os.ReadFile(serviceAccount)
	if err != nil {
		return "", err
	}

	var sa serviceAccountJSON
	if err = json.Unmarshal(credsJSON, &sa); err != nil {
		return "", err
	}

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: sa.ClientEmail,
		PrivateKey:     []byte(sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return "", err
	}
	return url, nil
}
