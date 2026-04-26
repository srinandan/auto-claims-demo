package handlers

import (
	"internal/database"
	"internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database")
	}

	db.AutoMigrate(&models.PolicyHolder{}, &models.Claim{}, &models.Photo{}, &models.AnalysisResult{}, &models.Estimate{})
	database.DB = db

	policy := models.PolicyHolder{PolicyNumber: "123", FirstName: "John", LastName: "Doe"}
	db.Create(&policy)
	return db
}
