package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"example.com/claims-app/internal/models"
	"github.com/gin-gonic/gin"
)

func TestBookAppointmentFailJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	claim := models.Claim{PolicyNumber: "123"}
	db.Create(&claim)

	router := gin.Default()
	router.POST("/api/claims/:id/book-appointment", BookAppointment)

	req, _ := http.NewRequest("POST", "/api/claims/1/book-appointment", bytes.NewBuffer([]byte(`{invalid}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

    // Test that the handler handles bind error correctly
    if w.Code != http.StatusBadRequest {
        t.Fatalf("Expected 400, got %d", w.Code)
    }
}
