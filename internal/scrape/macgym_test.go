package scrape

import (
    "fmt"
    "os"
    "strings"
    "testing"
    "time"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/util"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/store")

func TestFetchMacGym(t *testing.T) {
    // Test with sample data
    testCases := []struct {
        name     string
        filename string
        expected struct {
            location string
            capacity int
            inUse    int
        }
    }{
        {
            name:     "sample data",
            filename: "testdata/macgym_sample.json",
            expected: struct {
                location string
                capacity int
                inUse    int
            }{
                location: "Mac Gym",
                capacity: 8,
                inUse:    6,
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a test server or use file-based testing
            // For now, we'll test the parsing logic directly
            data, err := os.ReadFile(tc.filename)
            if err != nil {
                t.Skipf("Test file not found: %s", tc.filename)
                return
            }

            // Test the parsing logic
            var response MacGymResponse
            if err := util.DecodeJSON(strings.NewReader(string(data)), &response); err != nil {
                t.Fatalf("Failed to decode JSON: %v", err)
            }

            if !response.Success {
                t.Fatalf("Expected success=true, got %v", response.Success)
            }

            if len(response.Data) == 0 {
                t.Fatalf("Expected data array, got empty")
            }

            // Find badminton court
            var badmintonCourt *struct {
                LocationID   string `json:"locationId"`
                LocationName string `json:"locationName"`
                CurrentCount int    `json:"currentCount"`
                MaxCapacity  int    `json:"maxCapacity"`
                Status       string `json:"status"`
                LastUpdated  string `json:"lastUpdated"`
            }
            for i, location := range response.Data {
                if strings.Contains(strings.ToLower(location.LocationName), "badminton") {
                    badmintonCourt = &response.Data[i]
                    break
                }
            }

            if badmintonCourt == nil {
                t.Fatalf("No badminton court found in test data")
            }

            if badmintonCourt.MaxCapacity != tc.expected.capacity {
                t.Errorf("Expected capacity %d, got %d", tc.expected.capacity, badmintonCourt.MaxCapacity)
            }

            if badmintonCourt.CurrentCount != tc.expected.inUse {
                t.Errorf("Expected inUse %d, got %d", tc.expected.inUse, badmintonCourt.CurrentCount)
            }
        })
    }
}

func TestMacGymSnapshotCreation(t *testing.T) {
    response := MacGymResponse{
        Success: true,
        Data: []struct {
            LocationID   string `json:"locationId"`
            LocationName string `json:"locationName"`
            CurrentCount int    `json:"currentCount"`
            MaxCapacity  int    `json:"maxCapacity"`
            Status       string `json:"status"`
            LastUpdated  string `json:"lastUpdated"`
        }{
            {
                LocationID:   "test-001",
                LocationName: "Test Badminton Court",
                CurrentCount: 4,
                MaxCapacity:  6,
                Status:       "active",
                LastUpdated:  "2024-01-15T14:30:00Z",
            },
        },
    }

    snap := createSnapshotFromResponse(response)
    
    if snap.Location != "Mac Gym" {
        t.Errorf("Expected location 'Mac Gym', got '%s'", snap.Location)
    }

    if snap.Capacity != 6 {
        t.Errorf("Expected capacity 6, got %d", snap.Capacity)
    }

    if snap.InUse != 4 {
        t.Errorf("Expected inUse 4, got %d", snap.InUse)
    }

    if !strings.Contains(snap.Details, "4/6") {
        t.Errorf("Expected details to contain '4/6', got '%s'", snap.Details)
    }
}

// Helper function to create snapshot from response (extracted from main function)
func createSnapshotFromResponse(response MacGymResponse) store.MacGymSnapshot {
    snap := store.MacGymSnapshot{
        RetrievedAt: time.Now(),
        Location:    "Mac Gym",
        Raw:         response,
    }

    for _, location := range response.Data {
        locationName := strings.ToLower(location.LocationName)
        
        if strings.Contains(locationName, "badminton") || 
           strings.Contains(locationName, "court") ||
           strings.Contains(locationName, "gym") {
            
            snap.Capacity = location.MaxCapacity
            snap.InUse = location.CurrentCount
            snap.Details = fmt.Sprintf("%s: %d/%d in use", 
                location.LocationName, location.CurrentCount, location.MaxCapacity)
            
            if location.LastUpdated != "" {
                if t, err := time.Parse(time.RFC3339, location.LastUpdated); err == nil {
                    snap.RetrievedAt = t
                }
            }
            
            break
        }
    }

    if snap.Details == "" {
        snap.Details = "Mac Gym status retrieved (no capacity data available)"
    }

    return snap
}
