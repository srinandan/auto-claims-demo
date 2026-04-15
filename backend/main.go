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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"example.com/claims-app/database"
	"example.com/claims-app/handlers"
	"example.com/claims-app/pkg/seeder"
	"example.com/claims-app/pkg/telemetry"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	// Initialize Telemetry
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		// Fallback or skip
		// Since user requested GCP telemetry, we'll try to get it from metadata if library supports it automatically,
		// but providing an empty string to the exporter might cause it to auto-detect.
		// However, it's safer to warn if missing.
		fmt.Println("Warning: GOOGLE_CLOUD_PROJECT not set. Telemetry might fail or use default.")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	shutdown, err := telemetry.InitTelemetry(context.Background(), projectID, "auto-claims-backend")
	if err != nil {
		log.Printf("Failed to initialize telemetry: %v", err)
	} else {
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				log.Printf("Telemetry shutdown failed: %v", err)
			}
		}()
	}

	database.Connect()
	seeder.SeedPolicyHolders(database.DB)

	r := gin.Default()

	// Add OpenTelemetry middleware
	r.Use(otelgin.Middleware("auto-claims-backend"))

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "traceparent", "tracestate"}
	r.Use(cors.New(config))

	// Serve uploaded files
	r.Static("/uploads", "./uploads")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/claims", handlers.ListClaims)
		api.GET("/claims/:id", handlers.GetClaim)
		api.POST("/claims", handlers.CreateClaim)
		api.PUT("/claims/:id", handlers.UpdateClaim)
		api.DELETE("/claims/:id", handlers.DeleteClaim)
		api.POST("/claims/:id/analyze", handlers.AnalyzeClaim)
		api.POST("/claims/:id/repair-shops", handlers.FindRepairShops)
		api.POST("/claims/:id/book-appointment", handlers.BookAppointment)
		api.GET("/policies/:number", handlers.GetPolicy)
		api.POST("/resolve-address", handlers.ResolveAddress)
	}

	r.Run(fmt.Sprintf(":%s", PORT))
}
