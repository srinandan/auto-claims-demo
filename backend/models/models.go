package models

import (
	"time"

	"gorm.io/gorm"
)

type Claim struct {
	gorm.Model
	PolicyNumber         string     `json:"policy_number"`
	CustomerName         string     `json:"customer_name"`
	Status               string     `json:"status"` // New, Simple, Complex, Total Loss
	Description          string     `json:"description"`
	AccidentDate         time.Time  `json:"accident_date"`
	IncidentCity         string     `json:"incident_city"`
	IncidentState        string     `json:"incident_state"`
	IncidentType         string     `json:"incident_type"`
	CollisionType        string     `json:"collision_type"`
	Severity             string     `json:"severity"`
	Photos               []Photo    `json:"photos" gorm:"foreignKey:ClaimID"`
	TotalLossProbability float64    `json:"total_loss_probability"`
	ActionRequired       bool       `json:"action_required"`
	Estimates            []Estimate `json:"estimates" gorm:"foreignKey:ClaimID"`
}

type Photo struct {
	gorm.Model
	ClaimID        uint           `json:"claim_id"`
	URL            string         `json:"url"`
	AnalysisResult AnalysisResult `json:"analysis_result" gorm:"foreignKey:PhotoID"`
}

type AnalysisResult struct {
	gorm.Model
	PhotoID       uint   `json:"photo_id"`
	QualityScore  string `json:"quality_score"`  // Good, Blurry, Dark
	Detections    string `json:"detections"`     // JSON string of bounding boxes
	PartsDetected string `json:"parts_detected"` // Comma separated list
	Severity      string `json:"severity"`       // Low, Medium, High
}

type Estimate struct {
	gorm.Model
	ClaimID     uint    `json:"claim_id"`
	TotalAmount float64 `json:"total_amount"`
	Items       string  `json:"items"` // JSON string of line items
	Source      string  `json:"source"` // "AI", "Shop"
}

type PolicyHolder struct {
	gorm.Model
	PolicyNumber          string    `json:"policy_number" gorm:"uniqueIndex"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	MonthsAsCustomer      int       `json:"months_as_customer"`
	Age                   int       `json:"age"`
	PolicyBindDate        time.Time `json:"policy_bind_date"`
	PolicyState           string    `json:"policy_state"`
	PolicyCSL             string    `json:"policy_csl"`
	PolicyDeductible      int       `json:"policy_deductible"`
	PolicyAnnualPremium   float64   `json:"policy_annual_premium"`
	UmbrellaLimit         int       `json:"umbrella_limit"`
	InsuredZip            int       `json:"insured_zip"`
	InsuredSex            string    `json:"insured_sex"`
	InsuredEducationLevel string    `json:"insured_education_level"`
	InsuredOccupation     string    `json:"insured_occupation"`
	InsuredHobbies        string    `json:"insured_hobbies"`
	InsuredRelationship   string    `json:"insured_relationship"`
	CapitalGains          int       `json:"capital_gains"`
	CapitalLoss           int       `json:"capital_loss"`
	AutoMake              string    `json:"auto_make"`
	AutoModel             string    `json:"auto_model"`
	AutoYear              int       `json:"auto_year"`
}
