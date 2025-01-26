package main

import (
	"fmt"
	"log"
	"serviceCatalog/config"
	"serviceCatalog/internal/database"
	"serviceCatalog/internal/handlers"
	"serviceCatalog/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize handlers
	handler := handlers.NewHandler(db)
	router := setupRouter(handler)

	// Start server
	log.Printf("Server starting on :%d \n", cfg.Server.Port)
	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func setupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/services", h.ListServices)
	r.GET("/services/:id", h.GetService)
	r.GET("/services/:id/versions", h.GetServiceVersions)
	return r
}
