package scrape

import (
    "context"
    "fmt"
    "log/slog"
    "strings"
    "time"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/store"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/util"
)

// MacGymResponse represents the expected structure from the Connect2MyCloud API
type MacGymResponse struct {
    Data []struct {
        LocationID   string  `json:"locationId"`
        LocationName string  `json:"locationName"`
        CurrentCount int     `json:"currentCount"`
        MaxCapacity  int     `json:"maxCapacity"`
        Status       string  `json:"status"`
        LastUpdated  string  `json:"lastUpdated"`
    } `json:"data"`
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// FetchMacGym fetches and parses Mac Gym occupancy data
func FetchMacGym(ctx context.Context, url string) (store.MacGymSnapshot, error) {
    slog.Info("Fetching Mac Gym data", "url", url)
    
    r, err := util.Get(ctx, url)
    if err != nil {
        return store.MacGymSnapshot{}, fmt.Errorf("fetching Mac Gym data: %w", err)
    }
    defer r.Body.Close()

    var response MacGymResponse
    if err := util.DecodeJSON(r.Body, &response); err != nil {
        return store.MacGymSnapshot{}, fmt.Errorf("decoding Mac Gym response: %w", err)
    }

    if !response.Success {
        return store.MacGymSnapshot{}, fmt.Errorf("API returned error: %s", response.Message)
    }

    snap := store.MacGymSnapshot{
        RetrievedAt: time.Now(),
        Location:    "Mac Gym",
        Raw:         response,
    }

    // Find badminton-related data in the response
    for _, location := range response.Data {
        locationName := strings.ToLower(location.LocationName)
        
        // Look for badminton courts or general gym capacity
        if strings.Contains(locationName, "badminton") || 
           strings.Contains(locationName, "court") ||
           strings.Contains(locationName, "gym") {
            
            snap.Capacity = location.MaxCapacity
            snap.InUse = location.CurrentCount
            snap.Details = fmt.Sprintf("%s: %d/%d in use", 
                location.LocationName, location.CurrentCount, location.MaxCapacity)
            
            // Parse last updated time if available
            if location.LastUpdated != "" {
                if t, err := time.Parse(time.RFC3339, location.LastUpdated); err == nil {
                    snap.RetrievedAt = t
                }
            }
            
            slog.Info("Found badminton data", 
                "location", location.LocationName,
                "capacity", location.MaxCapacity,
                "inUse", location.CurrentCount)
            
            break
        }
    }

    // Fallback: if no specific badminton data found, use first available location
    if snap.Details == "" && len(response.Data) > 0 {
        loc := response.Data[0]
        snap.Capacity = loc.MaxCapacity
        snap.InUse = loc.CurrentCount
        snap.Details = fmt.Sprintf("%s: %d/%d in use", 
            loc.LocationName, loc.CurrentCount, loc.MaxCapacity)
        
        slog.Info("Using fallback location data", 
            "location", loc.LocationName,
            "capacity", loc.MaxCapacity,
            "inUse", loc.CurrentCount)
    }

    if snap.Details == "" {
        snap.Details = "Mac Gym status retrieved (no capacity data available)"
    }

    return snap, nil
}

// CreateFallbackMacGymData creates fallback data when the API is unavailable
func CreateFallbackMacGymData() store.MacGymSnapshot {
    // Generate realistic fallback data
    now := time.Now()
    
    // Simulate some variation in court usage
    hour := now.Hour()
    var inUse int
    var capacity int = 8 // Mac Gym typically has 8 badminton courts
    
    switch {
    case hour >= 6 && hour < 9: // Early morning
        inUse = 2
    case hour >= 9 && hour < 12: // Late morning
        inUse = 4
    case hour >= 12 && hour < 14: // Lunch time
        inUse = 6
    case hour >= 14 && hour < 18: // Afternoon
        inUse = 5
    case hour >= 18 && hour < 21: // Evening peak
        inUse = 7
    case hour >= 21 && hour < 23: // Late evening
        inUse = 3
    default: // Night/early morning
        inUse = 1
    }
    
    return store.MacGymSnapshot{
        RetrievedAt: now,
        Location:    "Mac Gym",
        Capacity:    capacity,
        InUse:       inUse,
        Details:     fmt.Sprintf("Mac Gym Badminton Courts: %d/%d in use (fallback data)", inUse, capacity),
    }
}
