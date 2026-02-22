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
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	//"time"

	"example.com/claims-app/database"
	"example.com/claims-app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Claim{}, &models.Photo{}, &models.Estimate{}, &models.AnalysisResult{})
	return db
}

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
