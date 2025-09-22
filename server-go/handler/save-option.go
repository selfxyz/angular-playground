package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"angular-playground-backend/config"

	self "github.com/selfxyz/self/sdk/sdk-go"
)

type SaveOptionsRequest struct {
	UserID  string      `json:"userId"`
	Options interface{} `json:"options"`
}

type SaveOptionsResponse struct {
	Message string `json:"message"`
}

// ConfigStoreInstance is imported from verify.go - they share the same instance

func GoSaveOptions(w http.ResponseWriter, r *http.Request) {
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

	var req SaveOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic("Invalid JSON")
	}

	if req.UserID == "" {
		panic("User ID is required")
	}

	if req.Options == nil {
		panic("Options are required")
	}

	// Convert frontend options to VerificationConfig
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		panic(err)
	}

	// Parse as SelfAppDisclosureConfig first
	var frontendOptions self.VerificationConfig
	err = json.Unmarshal(optionsJSON, &frontendOptions)
	if err != nil {
		panic(fmt.Sprintf("Options do not match expected structure: %v", err))
	}

	ctx := context.Background()

	// Store the SelfAppDisclosureConfig in the configs map for GetDisclosureConfig
	if ConfigStoreInstance == nil {
		ConfigStoreInstance = config.NewInMemoryConfigStore()
	}
	_, err = ConfigStoreInstance.SetConfig(ctx, req.UserID, frontendOptions)

	if err != nil {
		panic(err)
	}

	log.Printf("Saved options for user: %s, options: %+v", req.UserID, req.Options)

	response := SaveOptionsResponse{
		Message: "Options saved successfully",
	}

	json.NewEncoder(w).Encode(response)
}
