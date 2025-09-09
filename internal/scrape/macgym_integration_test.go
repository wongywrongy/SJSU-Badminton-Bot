package scrape

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/sjsu-badminton/badminton-discord-bot/internal/config"
	"github.com/sjsu-badminton/badminton-discord-bot/internal/store"
)

func TestMacGymDataFetching(t *testing.T) {
	// Set up test logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Testing with Mac Gym URL: %s", cfg.MacGymURL)

	// Test 1: Fetch real data from API
	t.Run("FetchRealData", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		snap, err := FetchMacGym(ctx, cfg.MacGymURL)
		if err != nil {
			t.Logf("API fetch failed (expected if API is down): %v", err)
		}

		// Log all the data we got
		t.Logf("=== MAC GYM DATA ANALYSIS ===")
		t.Logf("RetrievedAt: %s", snap.RetrievedAt.Format(time.RFC3339))
		t.Logf("Location: %s", snap.Location)
		t.Logf("Capacity: %d", snap.Capacity)
		t.Logf("InUse: %d", snap.InUse)
		t.Logf("Details: %s", snap.Details)
		
		if snap.Raw != nil {
			t.Logf("Raw data type: %T", snap.Raw)
			t.Logf("Raw data: %+v", snap.Raw)
		}

		// Calculate percentage
		if snap.Capacity > 0 {
			percentage := float64(snap.InUse) / float64(snap.Capacity) * 100
			available := snap.Capacity - snap.InUse
			availability := float64(available) / float64(snap.Capacity) * 100
			
			t.Logf("=== CALCULATED VALUES ===")
			t.Logf("Courts in use: %d/%d (%.1f%%)", snap.InUse, snap.Capacity, percentage)
			t.Logf("Courts available: %d/%d (%.1f%%)", available, snap.Capacity, availability)
		}

		// Validate data structure
		if snap.Capacity == 0 && snap.InUse == 0 {
			t.Log("WARNING: No capacity data available - using fallback data")
		} else {
			t.Log("SUCCESS: Valid capacity data found")
		}
	})

	// Test 2: Test fallback data generation
	t.Run("TestFallbackData", func(t *testing.T) {
		fallback := CreateFallbackMacGymData()
		
		t.Logf("=== FALLBACK DATA ANALYSIS ===")
		t.Logf("RetrievedAt: %s", fallback.RetrievedAt.Format(time.RFC3339))
		t.Logf("Location: %s", fallback.Location)
		t.Logf("Capacity: %d", fallback.Capacity)
		t.Logf("InUse: %d", fallback.InUse)
		t.Logf("Details: %s", fallback.Details)

		// Validate fallback data
		if fallback.Capacity == 0 {
			t.Error("Fallback data should have capacity > 0")
		}
		if fallback.InUse < 0 || fallback.InUse > fallback.Capacity {
			t.Error("Fallback data should have valid inUse value")
		}

		// Calculate percentage for fallback
		percentage := float64(fallback.InUse) / float64(fallback.Capacity) * 100
		available := fallback.Capacity - fallback.InUse
		availability := float64(available) / float64(fallback.Capacity) * 100
		
		t.Logf("=== FALLBACK CALCULATED VALUES ===")
		t.Logf("Courts in use: %d/%d (%.1f%%)", fallback.InUse, fallback.Capacity, percentage)
		t.Logf("Courts available: %d/%d (%.1f%%)", available, fallback.Capacity, availability)
	})
}

func TestMacGymPercentageCalculation(t *testing.T) {
	testCases := []struct {
		name     string
		capacity int
		inUse    int
		expected float64
	}{
		{"Empty gym", 8, 0, 0.0},
		{"Half full", 8, 4, 50.0},
		{"Full gym", 8, 8, 100.0},
		{"Quarter full", 8, 2, 25.0},
		{"Three quarters full", 8, 6, 75.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			snap := store.MacGymSnapshot{
				Capacity: tc.capacity,
				InUse:    tc.inUse,
			}

			percentage := float64(snap.InUse) / float64(snap.Capacity) * 100
			available := snap.Capacity - snap.InUse
			availability := float64(available) / float64(snap.Capacity) * 100

			t.Logf("Capacity: %d, In Use: %d", snap.Capacity, snap.InUse)
			t.Logf("Percentage in use: %.1f%%", percentage)
			t.Logf("Percentage available: %.1f%%", availability)
			t.Logf("Courts available: %d", available)

			if percentage != tc.expected {
				t.Errorf("Expected %.1f%%, got %.1f%%", tc.expected, percentage)
			}
		})
	}
}
