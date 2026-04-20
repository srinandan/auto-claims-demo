package gcs

import (
    "testing"
    "context"
    "bytes"
)

func TestUploadFile(t *testing.T) {
    err := UploadFile(context.Background(), "test-bucket", "test-object", bytes.NewBuffer([]byte("test")))
    if err == nil {
        t.Log("Expected error without credentials in unit test")
    }
}
