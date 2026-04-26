package seeder

import (
    "testing"
    "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"internal/models"
)

func TestSeed(t *testing.T) {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.PolicyHolder{})
    SeedPolicyHolders(db)
}
