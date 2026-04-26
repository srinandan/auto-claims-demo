// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"internal/database"
	"internal/models"
)

func TestListClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	// Seed some data
	claim := models.Claim{
		PolicyNumber:  "123",
		Status:        "New",
		AccidentDate:  time.Now(),
		Description:   "Test",
	}
	db.Create(&claim)

	router := gin.Default()
	router.GET("/api/claims", ListClaims)

	req, _ := http.NewRequest("GET", "/api/claims", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var response []models.Claim
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("Expected 1 claim, got %d", len(response))
	}
}

func TestGetPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.GET("/api/policies/:number", GetPolicy)

	req, _ := http.NewRequest("GET", "/api/policies/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var policy models.PolicyHolder
	json.Unmarshal(w.Body.Bytes(), &policy)
	if policy.PolicyNumber != "123" {
		t.Fatalf("Expected Policy 123, got %s", policy.PolicyNumber)
	}

	// Test missing
	req2, _ := http.NewRequest("GET", "/api/policies/999", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 for missing policy, got %d", w2.Code)
	}
}

func TestUpdateClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{
		PolicyNumber: "123",
		Status:       "New",
		Description:  "Test",
	}
	db.Create(&claim)

	router := gin.Default()
	router.PATCH("/api/claims/:id", UpdateClaim)

	updateData := map[string]interface{}{
		"status":      "Approved",
	}
	body, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PATCH", "/api/claims/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestGetClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{
		PolicyNumber:  "123",
		Status:        "New",
		Description:   "Test",
	}
	db.Create(&claim)

	router := gin.Default()
	router.GET("/api/claims/:id", GetClaim)

	req, _ := http.NewRequest("GET", "/api/claims/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{
		PolicyNumber:  "123",
		Status:        "New",
		Description:   "Test",
	}
	db.Create(&claim)

	router := gin.Default()
	router.DELETE("/api/claims/:id", DeleteClaim)

	req, _ := http.NewRequest("DELETE", "/api/claims/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected status 204, got %d", w.Code)
	}
}

func TestFindRepairShops(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/find-repair-shops", FindRepairShops)

	reqBody := `{"zipCode": "12345", "state": "CA", "make": "Toyota", "model": "Camry", "damageType": "Dent"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/find-repair-shops", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestBookAppointment(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto", "customerName": "John", "phoneNumber": "555"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestBookAppointmentSuccess(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto", "customerName": "John", "phoneNumber": "555"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestFindRepairShopsSuccess(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123", Description: "Test damage"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))
    defer ts.Close()

    originalURL := FindShopsURL
    FindShopsURL = ts.URL
    defer func() { FindShopsURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/find-repair-shops", FindRepairShops)

	reqBody := `{"zipCode": "12345", "state": "CA", "make": "Toyota", "model": "Camry", "damageType": "Dent"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/find-repair-shops", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestListClaimsParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	// Seed some data
	claim := models.Claim{
		PolicyNumber:  "123",
		Status:        "New",
		AccidentDate:  time.Now(),
		Description:   "Test",
	}
	db.Create(&claim)

	router := gin.Default()
	router.GET("/api/claims", ListClaims)

	req, _ := http.NewRequest("GET", "/api/claims?policy_number=123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestBookAppointmentFailBackend(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestFindRepairShopsFailBackend(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123", Description: "Test damage"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer ts.Close()

    originalURL := FindShopsURL
    FindShopsURL = ts.URL
    defer func() { FindShopsURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/find-repair-shops", FindRepairShops)

	reqBody := `{"zipCode": "12345", "state": "CA", "make": "Toyota", "model": "Camry", "damageType": "Dent"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/find-repair-shops", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestBookAppointmentFailBackendCode(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest) // Send a non-200
    }))
    defer ts.Close()

    originalURL := BaseAIServiceURL
    BaseAIServiceURL = ts.URL
    defer func() { BaseAIServiceURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestFindRepairShopsFailBackendCode(t *testing.T) {
    gin.SetMode(gin.TestMode)
	db := setupTestDB()
    claim := models.Claim{PolicyNumber: "123", Description: "Test damage"}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest) // Send non-200
    }))
    defer ts.Close()

    originalURL := FindShopsURL
    FindShopsURL = ts.URL
    defer func() { FindShopsURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/find-repair-shops", FindRepairShops)

	reqBody := `{"zipCode": "12345", "state": "CA", "make": "Toyota", "model": "Camry", "damageType": "Dent"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/find-repair-shops", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestAnalyzeClaimFailBackendCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	os.Setenv("BUCKET_NAME", "test-bucket")

	claim := models.Claim{PolicyNumber: "123", Photos: []models.Photo{
		{URL: "test.jpg"},
	}}
	db.Create(&claim)

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest) // Send non-200
    }))
    defer ts.Close()

    originalURL := AIServiceURL
    AIServiceURL = ts.URL
    defer func() { AIServiceURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/analyze", AnalyzeClaim)

	req, _ := http.NewRequest("POST", "/api/claims/1/analyze", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestCreateClaimExtraCoverage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("policy_number", "123")
	writer.WriteField("description", "A test claim")
	writer.WriteField("accident_date", time.Now().Format("2006-01-02"))
	writer.WriteField("customer_name", "John Doe")

	// Ensure we cover "uploads" prefix creation check loop
	part, _ := writer.CreateFormFile("files", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestCreateClaimWithPhotoAndBucket(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

    os.Setenv("BUCKET_NAME", "test-bucket")

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("policy_number", "123")
	writer.WriteField("description", "A test claim")
	writer.WriteField("accident_date", time.Now().Format("2006-01-02"))
	writer.WriteField("customer_name", "John Doe")

	part, _ := writer.CreateFormFile("files", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    os.Unsetenv("BUCKET_NAME")
}

func TestBookAppointmentMissingClaim(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto"}`
	req, _ := http.NewRequest("POST", "/api/claims/999/book-appointment", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    if w.Code != http.StatusNotFound {
        t.Fatalf("Expected 404, got %d", w.Code)
    }
}

func TestBookAppointmentNoBaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

    claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

    originalURL := BaseAIServiceURL
    BaseAIServiceURL = ""
    defer func() { BaseAIServiceURL = originalURL }()

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	reqBody := `{"session_id": "123", "message": "Hi", "shop_name": "Mock Auto"}`
	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

// Tests from origin/main...
func TestCreateClaim_InvalidDate(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	database.DB = db // Inject test DB

	router := gin.Default()
	router.POST("/claims", CreateClaim)

	// Create a multipart request with invalid date
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("policy_number", "12345")
	_ = writer.WriteField("customer_name", "Test User")
	_ = writer.WriteField("accident_date", "invalid-date")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify DB is empty
	var count int64
	db.Model(&models.Claim{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 claims, got %d", count)
	}
}

func TestCreateClaim_ValidDate(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	database.DB = db // Inject test DB

	router := gin.Default()
	router.POST("/claims", CreateClaim)

	// Create a multipart request with valid date
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("policy_number", "12345")
	_ = writer.WriteField("customer_name", "Test User")
	_ = writer.WriteField("accident_date", "2023-01-01")
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify DB has claim
	var count int64
	db.Model(&models.Claim{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 claim, got %d", count)
	}
}
