package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"angular-playground-backend/handler"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	response := HealthResponse{
		Status:    "OK",
		Message:   "Go backend server is running",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Go backend server starting on port %s\n", port)
	fmt.Printf("Health check available at http://localhost:%s/health\n", port)
	fmt.Printf("Verify endpoint available at http://localhost:%s/api/verify\n", port)

	// Five endpoints: health, verify, dual-verify, results, and saveOptions
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/verify", handler.VerifyHandler)
	http.HandleFunc("/api/saveOptions", handler.GoSaveOptions)

	// Default handler for root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Go Backend Server is Running!\nAvailable endpoints:\n- GET /health\n- POST /api/verify\n- POST /api/saveOptions\n- POST /api/dual-verify\n- GET /api/latest-results")
	})

	fmt.Printf("Server about to start listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
