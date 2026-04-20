package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"example.com/claims-app/models"
	"github.com/gin-gonic/gin"
    "os"
)

func TestListClaimsBucketURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    os.Setenv("BUCKET_NAME", "test-bucket")
	router := gin.Default()
	router.GET("/api/claims", ListClaims)

	req, _ := http.NewRequest("GET", "/api/claims", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}
