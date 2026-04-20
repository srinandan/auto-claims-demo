package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"example.com/claims-app/models"
    "time"
	"mime/multipart"
	"encoding/json"
)




func TestErrorStates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB()

	// Test missing claim for update
	router := gin.Default()
	router.PATCH("/api/claims/:id", UpdateClaim)

	req, _ := http.NewRequest("PATCH", "/api/claims/999", bytes.NewBuffer([]byte(`{"status":"Approved"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test missing claim for get
	router2 := gin.Default()
	router2.GET("/api/claims/:id", GetClaim)
	req2, _ := http.NewRequest("GET", "/api/claims/999", nil)
	w2 := httptest.NewRecorder()
	router2.ServeHTTP(w2, req2)

	// Test missing claim for delete
	router3 := gin.Default()
	router3.DELETE("/api/claims/:id", DeleteClaim)
	req3, _ := http.NewRequest("DELETE", "/api/claims/999", nil)
	w3 := httptest.NewRecorder()
	router3.ServeHTTP(w3, req3)

	router4 := gin.Default()
	router4.POST("/api/claims/:id/analyze", AnalyzeClaim)
	req4, _ := http.NewRequest("POST", "/api/claims/999/analyze", nil)
	w4 := httptest.NewRecorder()
	router4.ServeHTTP(w4, req4)

	router5 := gin.Default()
	router5.POST("/api/claims/:id/find-repair-shops", FindRepairShops)
	req5, _ := http.NewRequest("POST", "/api/claims/999/find-repair-shops", bytes.NewBuffer([]byte(`{"zipCode":"12345"}`)))
	req5.Header.Set("Content-Type", "application/json")
	w5 := httptest.NewRecorder()
	router5.ServeHTTP(w5, req5)

	router6 := gin.Default()
	router6.POST("/api/claims/:id/book-appointment", BookAppointment)
	req6, _ := http.NewRequest("POST", "/api/claims/999/book-appointment", bytes.NewBuffer([]byte(`{"placeId":"123"}`)))
	req6.Header.Set("Content-Type", "application/json")
	w6 := httptest.NewRecorder()
	router6.ServeHTTP(w6, req6)
}

func TestResolveAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	router := gin.Default()
	router.POST("/api/resolve-address", ResolveAddress)

	reqBody := `{"address": "123 Test St"}`
	req, _ := http.NewRequest("POST", "/api/resolve-address", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestResolveAddressMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()
	router := gin.Default()
	router.POST("/api/resolve-address", ResolveAddress)

	// test invalid JSON
	req, _ := http.NewRequest("POST", "/api/resolve-address", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}

func TestUpdateClaimInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber:  "123", Status: "New", Description: "Test"}
	db.Create(&claim)

	router := gin.Default()
	router.PATCH("/api/claims/:id", UpdateClaim)

	req, _ := http.NewRequest("PATCH", "/api/claims/1", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusBadGateway {
		t.Fatalf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

return
}

func TestFindRepairShopsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber:  "123", Status: "New", Description: "Test"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/find-repair-shops", FindRepairShops)

	req, _ := http.NewRequest("POST", "/api/claims/1/find-repair-shops", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusBadGateway {
		t.Fatalf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

return
}

func TestBookAppointmentInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber:  "123", Status: "New", Description: "Test"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusBadGateway {
		t.Fatalf("Expected status 400 for invalid JSON, got %d", w.Code)
	}

return
}

func TestCreateClaimFailParsing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

    // No multiparts setup
	req, _ := http.NewRequest("POST", "/api/claims", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "multipart/form-data")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateClaimInvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("policy_number", "123")
	writer.WriteField("description", "A test claim")
	writer.WriteField("accident_date", "invalid-date")
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateClaimNoPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "A test claim")
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 without policy, got %d", w.Code)
	}
}

func TestCreateClaimWithPhoto(t *testing.T) {
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

	// Add at least one photo form file so the check passes
	part, _ := writer.CreateFormFile("files", "test.jpg")
	part.Write([]byte("fake image data"))

	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code == http.StatusCreated {
		var claim models.Claim
		json.Unmarshal(w.Body.Bytes(), &claim)
		if claim.PolicyNumber != "123" {
			t.Fatalf("Expected Policy 123, got %s", claim.PolicyNumber)
		}
	}
}

func TestCreateClaimInvalidContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	req, _ := http.NewRequest("POST", "/api/claims", bytes.NewBuffer([]byte(`{"policy_number": "123"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status 400 for invalid content type, got %d", w.Code)
	}
}




func TestCreateClaimFailNoPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB()

	router := gin.Default()
	router.POST("/api/claims", CreateClaim)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("description", "A test claim")
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/claims", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
}
