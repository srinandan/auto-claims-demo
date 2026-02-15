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
	"example.com/claims-app/database"
	"example.com/claims-app/handlers"
	"example.com/claims-app/pkg/seeder"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	seeder.SeedPolicyHolders(database.DB)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
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
		api.POST("/claims/:id/analyze", handlers.AnalyzeClaim)
		api.GET("/policies/:number", handlers.GetPolicy)
	}

	r.Run(":8080")
}
