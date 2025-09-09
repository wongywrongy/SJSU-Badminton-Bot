package config

import (
    "os"
    "testing"
)

func TestLoad(t *testing.T) {
    // Save original environment
    originalToken := os.Getenv("DISCORD_BOT_TOKEN")
    defer func() {
        if originalToken != "" {
            os.Setenv("DISCORD_BOT_TOKEN", originalToken)
        } else {
            os.Unsetenv("DISCORD_BOT_TOKEN")
        }
    }()
    
    // Test minimal valid config
    os.Setenv("DISCORD_BOT_TOKEN", "test-token")
    cfg, err := Load()
    
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
        return
    }
    
    if cfg.Token != "test-token" {
        t.Errorf("Expected token 'test-token', got '%s'", cfg.Token)
    }
    
    if cfg.MacGymURL == "" {
        t.Error("Expected default MacGymURL to be set")
    }
    
    if cfg.FitnessURL == "" {
        t.Error("Expected default FitnessURL to be set")
    }
}

func TestLoadMissingToken(t *testing.T) {
    // Clear environment
    os.Unsetenv("DISCORD_BOT_TOKEN")
    
    _, err := Load()
    if err == nil {
        t.Error("Expected error for missing token but got none")
    }
}
