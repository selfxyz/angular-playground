package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	Result bool   `json:"result"`
	Reason string `json:"reason"`
}

// Global config store instance - similar to TypeScript version
// Exported so save-option.go can use the same instance
var ConfigStoreInstance *self.InMemoryConfigStore

// VerifyHandler handles the verification endpoint
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	fmt.Println("here1")
	if ConfigStoreInstance == nil {
		ConfigStoreInstance = self.NewInMemoryConfigStore(func(ctx context.Context, userIdentifier string, userDefinedData string) (string, error) {
			return userIdentifier, nil
		})
	}

	fmt.Println("here2")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("here3")
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	fmt.Println("here4")
	w.Header().Set("Content-Type", "application/json")

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("here5")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}

	fmt.Println("here6")

	ctx := context.Background()

	// For now, we'll create the verification config after we know the user identifier
	// The SDK will call GetConfig() during verification with the user's action ID

	// Define allowed attestation types
	allowedIds := map[self.AttestationId]bool{
		self.Passport: true,
		self.EUCard:   true,
		self.Aadhaar:  true,
	}

	// Use the same verifyEndpoint as TypeScript API to match scope calculation
	verifyEndpoint := "https://ceaf1286c8f7.ngrok-free.app/api/verify"

	verifier, err := self.NewBackendVerifier(
		"self-playground",
		verifyEndpoint,
		true, // Use testnet for testing
		allowedIds,
		ConfigStoreInstance,
		self.UserIDTypeUUID, // Use UUID format for user IDs
	)
	if err != nil {
		fmt.Println("here7")
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
		fmt.Println("here8")
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
		fmt.Println("here9")
		log.Printf("Verification failed - invalid result")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Verification failed",
		})
		return
	}

	if !result.IsValidDetails.IsMinimumAgeValid {
		fmt.Println("here10")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Minimum age check failed",
		})
		return
	}

	if result.IsValidDetails.IsOfacValid {
		fmt.Println("here11")
		response := VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "OFAC check failed",
		}
		fmt.Printf("Sending response: %+v\n", response)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Printf("JSON encoding error: %v\n", err)
		}
		return
	}

	fmt.Println("here12")
	// Check if verification is valid
	if result.IsValidDetails.IsValid {
		// Create filtered subject - copy the struct to modify it

		// Return successful verification result with filtered data
		fmt.Println("here13")
		response := VerifyResponse{
			Status: "success",
			Result: result.IsValidDetails.IsValid,
		}
		fmt.Printf("Sending success response: %+v\n", response)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Printf("JSON encoding error: %v\n", err)
		}
	} else {
		// Handle failed verification case
		fmt.Println("here14")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(VerifyResponse{
			Status: "error",
			Result: false,
			Reason: "Verification failed",
		})
	}
}
