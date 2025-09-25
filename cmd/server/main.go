package main

import (
	"log"
	"net/http"

	"calculator/internal/calculator"
	"calculator/internal/handlers"
	"calculator/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize calculator service
	calcService := calculator.NewService(cfg.PackSizes)

	// Initialize handlers
	handler := handlers.NewHandler(calcService)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Test edge case on startup
	log.Println("Testing edge case: Pack Sizes: [23, 31, 53], Amount: 500000")
	calcService.SetPackSizes([]int{23, 31, 53})
	result := calcService.CalculatePacks(500000)
	log.Printf("Edge case result: %+v", result)

	// Reset to configured pack sizes
	calcService.SetPackSizes(cfg.PackSizes)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Initial pack sizes: %v", calcService.GetPackSizes())

	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
