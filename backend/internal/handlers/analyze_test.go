package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"example.com/claims-app/internal/models"
	"github.com/gin-gonic/gin"
)

func TestAnalyzeClaimWithoutBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestAnalyzeClaimWithBucketNoPhotos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{}}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestAnalyzeClaimWithBucketPhotos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"findings": ["dent"], "agent_result": {"decision": "Approved", "reasoning": "Looks good", "estimate": {"items": [{"name": "part", "cost": 10}], "total_parts": 100, "total_cost": 200}}, "photo_analyses": {"gs://test-bucket/test.jpg": [{"label": "dent"}]}}`))
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	originalBaseURL := BaseAIServiceURL
	BaseAIServiceURL = ts.URL
	defer func() { BaseAIServiceURL = originalBaseURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestAnalyzeClaimAiFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	originalBaseURL := BaseAIServiceURL
	BaseAIServiceURL = ts.URL
	defer func() { BaseAIServiceURL = originalBaseURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestAnalyzeClaimAiFailDecode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{invalid-json}`))
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	originalBaseURL := BaseAIServiceURL
	BaseAIServiceURL = ts.URL
	defer func() { BaseAIServiceURL = originalBaseURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestAnalyzeClaimNoBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

	originalBaseURL := BaseAIServiceURL
	BaseAIServiceURL = ""
	defer func() { BaseAIServiceURL = originalBaseURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}
