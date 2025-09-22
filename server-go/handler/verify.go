package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"angular-playground-backend/config"

	self "github.com/selfxyz/self/sdk/sdk-go"
)

type VerifyRequest struct {
	AttestationID   int                     `json:"attestationId"`
	Proof           self.VcAndDiscloseProof `json:"proof"`
	PublicSignals   self.PublicSignals      `json:"publicSignals"`
	UserContextData string                  `json:"userContextData"`
}

type VerifyResponse struct {
	Status string `json:"status"`
	Result bool   `json:"result,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// Global config store instance - similar to TypeScript version
// Exported so save-option.go can use the same instance
var ConfigStoreInstance *config.InMemoryConfigStore

// VerifyHandler handles the verification endpoint
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if ConfigStoreInstance == nil {
		ConfigStoreInstance = config.NewInMemoryConfigStore()
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}

	ctx := context.Background()
	// Check if global config store is available
	if ConfigStoreInstance == nil {
		log.Printf("Config store not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Internal server error",
		})
		return
	}

	// For now, we'll create the verification config after we know the user identifier
	// The SDK will call GetConfig() during verification with the user's action ID

	// Define allowed attestation types
	allowedIds := map[self.AttestationId]bool{
		self.Passport: true,
		self.EUCard:   true,
	}

	// Use the same verifyEndpoint as TypeScript API to match scope calculation
	verifyEndpoint := "https://cc10778f114e.ngrok-free.app/api/verify"

	verifier, err := self.NewBackendVerifier(
		"self-playground",
		verifyEndpoint,
		true, // Use testnet for testing
		allowedIds,
		ConfigStoreInstance,
		self.UserIDTypeUUID, // Use UUID format for user IDs
	)
	if err != nil {
		log.Printf("Failed to initialize verifier: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Internal server error",
		})
		return
	}

	result, err := verifier.Verify(
		ctx,
		req.AttestationID,
		req.Proof,
		req.PublicSignals,
		req.UserContextData,
	)
	if err != nil {
		log.Printf("Verification failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: err.Error(),
		})
		return
	}

	if result == nil || !result.IsValidDetails.IsValid {
		log.Printf("Verification failed - invalid result")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Verification failed",
		})
		return
	}

	// Check if verification is valid
	if result.IsValidDetails.IsValid {
		// Create filtered subject - copy the struct to modify it

		// Return successful verification result with filtered data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "success",
			Result: result.IsValidDetails.IsValid,
		})
	} else {
		// Handle failed verification case
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Reason: "Verification failed",
		})
	}
}
