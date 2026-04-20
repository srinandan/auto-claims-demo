package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"example.com/claims-app/internal/models"
	"github.com/gin-gonic/gin"
    "os"
)

func TestGetClaimBucketURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    os.Setenv("BUCKET_NAME", "test-bucket")
	router := gin.Default()
	router.GET("/api/claims/:id", GetClaim)

	req, _ := http.NewRequest("GET", "/api/claims/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}
