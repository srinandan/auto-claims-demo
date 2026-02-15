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
		api.POST("/claims/:id/analyze", handlers.AnalyzeClaim)
		api.GET("/policies/:number", handlers.GetPolicy)
	}

	r.Run(":8080")
}
