package seeder

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/claims-app/models"
	"gorm.io/gorm"
)

const DataURL = "https://raw.githubusercontent.com/hongdnn/claimwise_ai/refs/heads/main/data/insurance_claims.csv"

func SeedPolicyHolders(db *gorm.DB) {
	// Check if table is empty
	var count int64
	db.Model(&models.PolicyHolder{}).Count(&count)
	if count > 0 {
		log.Println("PolicyHolders table already seeded.")
		return
	}

	log.Println("Seeding PolicyHolders table from remote CSV...")

	resp, err := http.Get(DataURL)
	if err != nil {
		log.Printf("Failed to fetch CSV: %v", err)
		return
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Failed to read CSV: %v", err)
		return
	}

	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		// Helper to parse int safely
		parseInt := func(s string) int {
			val, _ := strconv.Atoi(s)
			return val
		}
		// Helper to parse float safely
		parseFloat := func(s string) float64 {
			val, _ := strconv.ParseFloat(s, 64)
			return val
		}
		// Helper to parse time safely (2014-10-17)
		parseTime := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		policy := models.PolicyHolder{
			MonthsAsCustomer:      parseInt(record[0]),
			Age:                   parseInt(record[1]),
			PolicyNumber:          record[2],
			PolicyBindDate:        parseTime(record[3]),
			PolicyState:           record[4],
			PolicyCSL:             record[5],
			PolicyDeductible:      parseInt(record[6]),
			PolicyAnnualPremium:   parseFloat(record[7]),
			UmbrellaLimit:         parseInt(record[8]),
			InsuredZip:            parseInt(record[9]),
			InsuredSex:            record[10],
			InsuredEducationLevel: record[11],
			InsuredOccupation:     record[12],
			InsuredHobbies:        record[13],
			InsuredRelationship:   record[14],
			CapitalGains:          parseInt(record[15]),
			CapitalLoss:           parseInt(record[16]),
			AutoMake:              record[35],
			AutoModel:             record[36],
			AutoYear:              parseInt(record[37]),
		}

		if err := db.Create(&policy).Error; err != nil {
			log.Printf("Failed to insert policy %s: %v", policy.PolicyNumber, err)
		}
	}

	log.Println("PolicyHolders table seeded successfully.")
}
