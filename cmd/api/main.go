package main

import (
	"log"
	"serviceCatalog/internal/database"
	"serviceCatalog/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Routes
	r.GET("/services", h.ListServices)
	r.GET("/services/:id", h.GetService)
	r.GET("/services/:id/versions", h.GetServiceVersions)

	// Start server
	log.Println("Server starting on :8080")
	r.Run(":8080")

}
