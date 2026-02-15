package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"cloud.google.com/go/storage"
)

var (
	iamClient *credentials.IamCredentialsClient
	initOnce  sync.Once
	initErr   error
)

// Init initializes the IAM client. It is safe to call multiple times.
func Init(ctx context.Context) error {
	initOnce.Do(func() {
		iamClient, initErr = credentials.NewIamCredentialsClient(ctx)
	})
	return initErr
}

// UploadFile uploads an object to GCS
func UploadFile(ctx context.Context, bucketName, objectName string, content io.Reader) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.Copy(wc, content); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}

// GenerateSignedURL generates a V4 signed URL for an object
func GenerateSignedURL(ctx context.Context, bucketName, objectName string) (string, error) {
	// Initialize if not already done (best effort, mainly for simpler usage)
	if iamClient == nil {
		if err := Init(ctx); err != nil {
			return "", fmt.Errorf("failed to initialize IAM client: %w", err)
		}
	}

	// 1. Determine Service Account Email
	email := os.Getenv("SERVICE_ACCOUNT_EMAIL")
	if email == "" {
		if metadata.OnGCE() {
			var err error
			email, err = metadata.Email("default")
			if err != nil {
				return "", fmt.Errorf("metadata.Email: %w", err)
			}
		} else {
			return "", fmt.Errorf("SERVICE_ACCOUNT_EMAIL env var not set and not on GCE")
		}
	}

	// 2. Define SignBytes function using IAM Credentials API
	signBytes := func(b []byte) ([]byte, error) {
		req := &credentialspb.SignBlobRequest{
			Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", email),
			Payload: b,
		}
		resp, err := iamClient.SignBlob(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("SignBlob: %w", err)
		}
		return resp.SignedBlob, nil
	}

	// 3. Generate URL
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: email,
		Expires:        time.Now().Add(15 * time.Minute),
		SignBytes:      signBytes,
	}

	return storage.SignedURL(bucketName, objectName, opts)
}
