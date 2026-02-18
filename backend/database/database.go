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

package database

import (
	"os"
	"log"
	"example.com/claims-app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "claims.db"
	}
	DB, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Claim{}, &models.Photo{}, &models.AnalysisResult{}, &models.Estimate{}, &models.PolicyHolder{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database connected and migrated successfully.")
}
