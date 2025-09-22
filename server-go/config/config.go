package config

import (
	"context"
	"fmt"
	"sync"

	self "github.com/selfxyz/self/sdk/sdk-go"
)

// InMemoryConfigStore provides in-memory storage for configurations and options
type InMemoryConfigStore struct {
	mu      sync.RWMutex
	configs map[string]self.VerificationConfig
}

// NewInMemoryConfigStore creates a new in-memory config store
func NewInMemoryConfigStore() *InMemoryConfigStore {
	return &InMemoryConfigStore{
		configs: make(map[string]self.VerificationConfig),
	}
}

// GetActionId implements the ConfigStore interface
func (store *InMemoryConfigStore) GetActionId(ctx context.Context, userIdentifier string, userDefinedData string) (string, error) {
	return userIdentifier, nil
}

// SetConfig implements the ConfigStore interface
func (store *InMemoryConfigStore) SetConfig(ctx context.Context, id string, config self.VerificationConfig) (bool, error) {
	// This method is required by the interface but not used in our implementation
	// The SDK will call GetConfig() to retrieve configurations
	store.mu.Lock()
	store.configs[id] = config
	store.mu.Unlock()
	return true, nil
}

// GetConfig implements the ConfigStore interface and returns self.VerificationConfig
func (store *InMemoryConfigStore) GetConfig(ctx context.Context, id string) (self.VerificationConfig, error) {
	// Check if we have a pre-stored config for this ID
	store.mu.RLock()
	config, exists := store.configs[id]
	store.mu.RUnlock()

	if exists {
		return config, nil
	}

	// Panic if no config found - matches expected behavior
	panic(fmt.Sprintf("Configuration not found for user: %s", id))
}

// Close cleans up resources (no-op for in-memory store)
func (store *InMemoryConfigStore) Close() error {
	return nil
}
