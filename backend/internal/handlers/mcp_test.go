package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
    "os"
)



func TestResolveAddressSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/resolve-address", ResolveAddress)

    os.Setenv("GOOGLE_MAPS_API_KEY", "test_key")
    defer os.Unsetenv("GOOGLE_MAPS_API_KEY")

	reqBody := `{"address": "123 Test St"}`
	req, _ := http.NewRequest("POST", "/api/resolve-address", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

}

func TestResolveAddressNoKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/resolve-address", ResolveAddress)

    os.Setenv("GOOGLE_MAPS_API_KEY", "")
    defer os.Unsetenv("GOOGLE_MAPS_API_KEY")

	reqBody := `{"address": "123 Test St"}`
	req, _ := http.NewRequest("POST", "/api/resolve-address", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("Expected 400 missing key, got %d", w.Code)
    }
}


func TestResolveAddressToolFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/resolve-address", ResolveAddress)

    os.Setenv("GOOGLE_MAPS_API_KEY", "test_key")

	reqBody := `{"address": "123 Test St"}`
	req, _ := http.NewRequest("POST", "/api/resolve-address", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("GOOGLE_MAPS_API_KEY")
}
