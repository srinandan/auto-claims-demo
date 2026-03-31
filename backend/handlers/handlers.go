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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"example.com/claims-app/database"
	"example.com/claims-app/models"
	"example.com/claims-app/pkg/gcs"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	BaseAIServiceURL string
	AIServiceURL     string
	FindShopsURL     string
	httpClient       *http.Client
)

func init() {
	BaseAIServiceURL = strings.TrimSpace(os.Getenv("AI_SERVICE_URL"))
	if BaseAIServiceURL == "" {
		log.Println("Warning: AI_SERVICE_URL environment variable is not set. AI features will be disabled or fail.")
	} else {
		log.Printf("Configured AI Service URL: %s", BaseAIServiceURL)
	}

	var err error
	// Construct service URLs from base
	AIServiceURL, err = url.JoinPath(BaseAIServiceURL, "process-claims")
	if err != nil {
		log.Printf("Error creating AIServiceURL: %v\n", err)
		// Fallback to simple concatenation if JoinPath fails
		AIServiceURL = fmt.Sprintf("%s/%s", strings.TrimRight(BaseAIServiceURL, "/"), "process-claims")
	}

	FindShopsURL, err = url.JoinPath(BaseAIServiceURL, "find-repair-shops")
	if err != nil {
		log.Printf("Error creating FindShopsURL: %v\n", err)
		FindShopsURL = fmt.Sprintf("%s/%s", strings.TrimRight(BaseAIServiceURL, "/"), "find-repair-shops")
	}

	// Initialize instrumented HTTP client
	httpClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   120 * time.Second, // 2 minute timeout for AI operations
	}
}

// ListClaims returns all claims
func ListClaims(c *gin.Context) {
	var claims []models.Claim
	query := database.DB.Preload("Photos").Preload("Estimates")

	policyNumber := c.Query("policy_number")
	if policyNumber != "" {
		query = query.Where("policy_number = ?", policyNumber)
	}

	if err := query.Find(&claims).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bucketName := os.Getenv("BUCKET_NAME")

	// Generate signed URLs for photos
	for i := range claims {
		for j := range claims[i].Photos {
			// If URL is a GCS path (doesn't start with http or uploads/), generate signed URL
			// We assume anything not starting with uploads/ is a GCS path
			if !strings.HasPrefix(claims[i].Photos[j].URL, "uploads/") && bucketName != "" {
				signedURL, err := gcs.GenerateSignedURL(c.Request.Context(), bucketName, claims[i].Photos[j].URL)
				if err == nil {
					claims[i].Photos[j].URL = signedURL
				} else {
					fmt.Printf("Error generating signed URL for %s: %v\n", claims[i].Photos[j].URL, err)
				}
			}
		}
	}

	c.JSON(http.StatusOK, claims)
}

// GetClaim returns a single claim by ID
func GetClaim(c *gin.Context) {
	id := c.Param("id")
	var claim models.Claim
	// Preload Photos and AnalysisResult. Note: AnalysisResult is singular for Photo
	if err := database.DB.Preload("Photos.AnalysisResult").Preload("Estimates").First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	bucketName := os.Getenv("BUCKET_NAME")

	// Generate signed URLs for photos
	for i := range claim.Photos {
		if !strings.HasPrefix(claim.Photos[i].URL, "uploads/") && bucketName != "" {
			signedURL, err := gcs.GenerateSignedURL(c.Request.Context(), bucketName, claim.Photos[i].URL)
			if err == nil {
				claim.Photos[i].URL = signedURL
			} else {
				fmt.Printf("Error generating signed URL for %s: %v\n", claim.Photos[i].URL, err)
			}
		}
	}

	c.JSON(http.StatusOK, claim)
}

// FindRepairShops triggers AI agent to find repair shops
func FindRepairShops(c *gin.Context) {
	id := c.Param("id")
	var claim models.Claim
	if err := database.DB.First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	var policy models.PolicyHolder
	if err := database.DB.Where("policy_number = ?", claim.PolicyNumber).First(&policy).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Policy holder not found"})
		return
	}

	// Determine damage type
	damageType := claim.Description
	if damageType == "" {
		damageType = "auto body repair"
	}

	reqBody := map[string]string{
		"zip_code":    fmt.Sprintf("%d", policy.InsuredZip),
		"state":       policy.PolicyState,
		"make":        policy.AutoMake,
		"model":       policy.AutoModel,
		"damage_type": damageType,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), "POST", FindShopsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to call AI service: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("AI service returned status %d", resp.StatusCode)})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode AI response"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// BookAppointment proxies request to AI service for booking
func BookAppointment(c *gin.Context) {
	id := c.Param("id")
	var claim models.Claim
	if err := database.DB.First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	var input struct {
		SessionID    string `json:"session_id"`
		Message      string `json:"message"`
		ShopName     string `json:"shop_name"`
		CustomerName string `json:"customer_name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use claim's customer name if not provided (though input struct defaults are empty string)
	customerName := input.CustomerName
	if customerName == "" {
		customerName = claim.CustomerName
	}

	// Construct context
	contextMap := map[string]interface{}{
		"shop_name":     input.ShopName,
		"customer_name": customerName,
	}

	reqBody := map[string]interface{}{
		"session_id": input.SessionID,
		"message":    input.Message,
		"context":    contextMap,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	u, _ := url.JoinPath(BaseAIServiceURL, "book-appointment")
	req, err := http.NewRequestWithContext(c.Request.Context(), "POST", u, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to call AI service: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("AI service returned status %d", resp.StatusCode)})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode AI response"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPolicy returns a policy holder by number
func GetPolicy(c *gin.Context) {
	number := c.Param("number")
	var policy models.PolicyHolder
	if err := database.DB.Where("policy_number = ?", number).First(&policy).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
		return
	}
	c.JSON(http.StatusOK, policy)
}

// UpdateClaim updates a claim status
func UpdateClaim(c *gin.Context) {
	id := c.Param("id")
	var claim models.Claim
	if err := database.DB.First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != "" {
		claim.Status = input.Status
	}

	database.DB.Save(&claim)
	c.JSON(http.StatusOK, claim)
}

// CreateClaim creates a claim with file uploads
func CreateClaim(c *gin.Context) {
	// Check if content type is multipart
	contentType := c.Request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be multipart/form-data"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data: " + err.Error()})
		return
	}

	policyNumber := c.PostForm("policy_number")
	customerName := c.PostForm("customer_name")
	description := c.PostForm("description")
	accidentDateStr := c.PostForm("accident_date")
	incidentCity := c.PostForm("incident_city")
	incidentState := c.PostForm("incident_state")
	incidentType := c.PostForm("incident_type")
	collisionType := c.PostForm("collision_type")
	severity := c.PostForm("severity")

	if policyNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Policy number is required"})
		return
	}

	var accidentDate time.Time
	if accidentDateStr != "" {
		// Try parsing date, handle errors gracefully or default
		accidentDate, err = time.Parse("2006-01-02", accidentDateStr)
		if err != nil {
			// Try RFC3339 just in case
			accidentDate, err = time.Parse(time.RFC3339, accidentDateStr)
			if err != nil {
				fmt.Printf("Error parsing accident date %s: %v\n", accidentDateStr, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid accident date format. Use YYYY-MM-DD or RFC3339"})
				return
			}
		}
	}

	claim := models.Claim{
		PolicyNumber:  policyNumber,
		CustomerName:  customerName,
		Description:   description,
		Status:        "New",
		AccidentDate:  accidentDate,
		IncidentCity:  incidentCity,
		IncidentState: incidentState,
		IncidentType:  incidentType,
		CollisionType: collisionType,
		Severity:      severity,
	}

	if err := database.DB.Create(&claim).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		// Fallback to local if bucket not set? Or error?
		// For now, let's log warning and error out or fallback.
		// User requirement implies bucket is mandatory.
		fmt.Println("Warning: BUCKET_NAME env var not set. Uploads will fail.")
	}

	files := form.File["files"]
	for _, file := range files {
		// Open the file
		f, err := file.Open()
		if err != nil {
			fmt.Printf("Failed to open uploaded file: %v\n", err)
			continue
		}

		// Construct object name: policyNumber/filename
		// Using policyNumber as folder
		objectName := fmt.Sprintf("%s/%s", policyNumber, file.Filename)

		// Upload to GCS
		if bucketName != "" {
			if err := gcs.UploadFile(context.Background(), bucketName, objectName, f); err != nil {
				fmt.Printf("Failed to upload file %s to GCS: %v\n", objectName, err)
				f.Close()
				continue
			}

			photo := models.Photo{
				ClaimID: claim.ID,
				URL:     objectName, // Store the object path
			}
			database.DB.Create(&photo)
		} else {
			f.Close()
			// Fallback to local storage if bucket not defined (legacy support or error?)
			// I'll return error to force configuration
			c.JSON(http.StatusInternalServerError, gin.H{"error": "BUCKET_NAME not configured"})
			return
		}
		f.Close()
	}

	// Reload to return full object
	database.DB.Preload("Photos").First(&claim, claim.ID)
	// We might want to sign the URLs in the response here too, but for creation it might be okay to return raw path or nothing?
	// Consistent behavior: return signed URLs.
	for i := range claim.Photos {
		if !strings.HasPrefix(claim.Photos[i].URL, "uploads/") && bucketName != "" {
			signedURL, err := gcs.GenerateSignedURL(c.Request.Context(), bucketName, claim.Photos[i].URL)
			if err == nil {
				claim.Photos[i].URL = signedURL
			}
		}
	}

	c.JSON(http.StatusCreated, claim)
}

// AnalyzeClaim triggers AI analysis for a claim
func AnalyzeClaim(c *gin.Context) {
	id := c.Param("id")
	var claim models.Claim
	if err := database.DB.Preload("Photos").First(&claim, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "BUCKET_NAME not configured"})
		return
	}

	// Initial status update
	claim.Status = "Analyzing"
	database.DB.Save(&claim)

	// Prepare JSON request to AI Service with GCS URIs
	fileURIs := []string{}
	for _, photo := range claim.Photos {
		if !strings.HasPrefix(photo.URL, "uploads/") {
			// Construct gs:// URI
			uri := fmt.Sprintf("gs://%s/%s", bucketName, photo.URL)
			fileURIs = append(fileURIs, uri)
		} else {
			// Skip local files or handle them?
			// The new AI service expects GCS URIs.
			fmt.Printf("Skipping local file %s for analysis\n", photo.URL)
		}
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"file_uris": fileURIs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body: " + err.Error()})
		return
	}

	// Check if AI Service is configured
	if BaseAIServiceURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI_SERVICE_URL environment variable is not set"})
		return
	}

	// Call AI Service
	req, err := http.NewRequestWithContext(c.Request.Context(), "POST", AIServiceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Use instrumented client
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to call AI service: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to call AI service: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("AI service returned status %d", resp.StatusCode)})
		return
	}

	// Parse Response
	// Structure matches ClaimsProcessResponse in Python
	var aiResponse struct {
		Findings    []string `json:"findings"`
		AgentResult struct {
			Decision  string `json:"decision"`
			Reasoning string `json:"reasoning"`
			Estimate  struct {
				Items []struct {
					Part string  `json:"part"`
					Cost float64 `json:"cost"`
				} `json:"items"`
				TotalCost float64 `json:"total_cost"`
			} `json:"estimate"`
		} `json:"agent_result"`
		PhotoAnalyses map[string][]struct {
			Label string    `json:"label"`
			Box   []float64 `json:"box"`
			Score float64   `json:"score"`
		} `json:"photo_analyses"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode AI response: " + err.Error()})
		return
	}

	// Update Analysis Results for Photos
	for filename, detections := range aiResponse.PhotoAnalyses {
		// Find the photo corresponding to this filename
		// The key in PhotoAnalyses from AI service might be the full GCS URI or just the filename?
		// We need to check how AI service returns the keys.
		// Usually it returns the filename or the identifier we sent.
		// If we sent gs://bucket/policy/filename, it might return that.
		// Let's assume AI service keys by the input string (the URI) or the filename.
		// I'll implement logic in AI service to return keys that we can match.
		// Safest is to match against the URI or the filename part.

		for _, photo := range claim.Photos {
			// Construct URI again to match
			uri := fmt.Sprintf("gs://%s/%s", bucketName, photo.URL)

			// Check if key matches URI or ends with filename (if AI service extracts filename)
			if filename == uri || strings.HasSuffix(uri, filename) || strings.HasSuffix(filename, photo.URL) {
				// Save Analysis Result
				detBytes, _ := json.Marshal(detections)

				// Extract unique labels for PartsDetected
				partsMap := make(map[string]bool)
				var partsList []string
				for _, d := range detections {
					if !partsMap[d.Label] {
						partsMap[d.Label] = true
						partsList = append(partsList, d.Label)
					}
				}
				partsDetected := strings.Join(partsList, ",")

				var analysis models.AnalysisResult
				if err := database.DB.Where("photo_id = ?", photo.ID).Limit(1).Find(&analysis).Error; err != nil {
					// Real DB error
					fmt.Printf("Error querying analysis: %v\n", err)
				}

				if analysis.ID == 0 {
					// Record not found (ID is 0), create new
					analysis = models.AnalysisResult{
						PhotoID:       photo.ID,
						QualityScore:  "Good", // Assuming Good for now
						Detections:    string(detBytes),
						PartsDetected: partsDetected,
						Severity:      "Unknown",
					}
					database.DB.Create(&analysis)
				} else {
					analysis.Detections = string(detBytes)
					analysis.PartsDetected = partsDetected
					database.DB.Save(&analysis)
				}
			}
		}
	}

	// Update Claim Status
	// Set to "Assessed" regardless of AI decision, to let user decide
	claim.Status = "Assessed"

	// Update Estimate
	estimateItems, _ := json.Marshal(aiResponse.AgentResult.Estimate.Items)
	totalCost := aiResponse.AgentResult.Estimate.TotalCost

	// Create/Update Estimate
	var estimate models.Estimate
	if err := database.DB.Where("claim_id = ?", claim.ID).Limit(1).Find(&estimate).Error; err != nil {
		fmt.Printf("Error querying estimate: %v\n", err)
	}

	if estimate.ID == 0 {
		estimate = models.Estimate{
			ClaimID:     claim.ID,
			TotalAmount: totalCost,
			Source:      "AI Agent",
			Items:       string(estimateItems),
		}
		database.DB.Create(&estimate)
	} else {
		estimate.TotalAmount = totalCost
		estimate.Source = "AI Agent"
		estimate.Items = string(estimateItems)
		database.DB.Save(&estimate)
	}

	database.DB.Save(&claim)

	// Return updated claim
	database.DB.Preload("Photos.AnalysisResult").Preload("Estimates").First(&claim, id)

	// Sign URLs for response
	for i := range claim.Photos {
		if !strings.HasPrefix(claim.Photos[i].URL, "uploads/") && bucketName != "" {
			signedURL, err := gcs.GenerateSignedURL(c.Request.Context(), bucketName, claim.Photos[i].URL)
			if err == nil {
				claim.Photos[i].URL = signedURL
			}
		}
	}

	c.JSON(http.StatusOK, claim)
}

// DeleteClaim deletes a claim and its associations
func DeleteClaim(c *gin.Context) {
	id := c.Param("id")

	// Delete associations first
	// Note: GORM might handle this if cascading delete is set up, but doing it explicitly ensures it works with Unscoped
	database.DB.Unscoped().Where("claim_id = ?", id).Delete(&models.Estimate{})

	// We also need to delete Photos and their AnalysisResults.
	// Since AnalysisResult depends on Photo, we should delete AnalysisResults for Photos of this Claim.
	// This is getting complicated to do manually.
	// But let's look at the models. Photo has gorm:"foreignKey:ClaimID".

	// Find photos to delete their analysis results
	var photos []models.Photo
	database.DB.Unscoped().Where("claim_id = ?", id).Find(&photos)
	for _, p := range photos {
		database.DB.Unscoped().Where("photo_id = ?", p.ID).Delete(&models.AnalysisResult{})
	}

	database.DB.Unscoped().Where("claim_id = ?", id).Delete(&models.Photo{})

	// Delete the claim
	result := database.DB.Unscoped().Delete(&models.Claim{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Claim not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
