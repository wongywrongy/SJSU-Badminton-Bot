package scrape

import (
    "strings"
    "testing"
    "time"
)

func TestCreateFallbackMacGymData(t *testing.T) {
    snap := CreateFallbackMacGymData()
    
    if snap.Location != "Mac Gym" {
        t.Errorf("Expected location 'Mac Gym', got '%s'", snap.Location)
    }
    
    if snap.Capacity != 8 {
        t.Errorf("Expected capacity 8, got %d", snap.Capacity)
    }
    
    if snap.InUse < 0 || snap.InUse > snap.Capacity {
        t.Errorf("InUse %d should be between 0 and %d", snap.InUse, snap.Capacity)
    }
    
    if snap.Details == "" {
        t.Error("Details should not be empty")
    }
    
    if !strings.Contains(snap.Details, "fallback data") {
        t.Error("Details should indicate this is fallback data")
    }
}

func TestCreateFallbackBadmintonEvents(t *testing.T) {
    loc, err := time.LoadLocation("America/Los_Angeles")
    if err != nil {
        t.Fatalf("Failed to load timezone: %v", err)
    }
    
    events := CreateFallbackBadmintonEvents(loc)
    
    if len(events) == 0 {
        t.Error("Should create at least some fallback events")
    }
    
    now := time.Now().In(loc)
    
    for i, event := range events {
        if event.Title == "" {
            t.Errorf("Event[%d] title should not be empty", i)
        }
        
        if !strings.Contains(strings.ToLower(event.Title), "badminton") {
            t.Errorf("Event[%d] title should contain 'badminton', got '%s'", i, event.Title)
        }
        
        if event.Location != "SJSU Fitness Center" {
            t.Errorf("Event[%d] location should be 'SJSU Fitness Center', got '%s'", i, event.Location)
        }
        
        if event.Start.Before(now) {
            t.Errorf("Event[%d] start time %v should be in the future", i, event.Start)
        }
        
        if event.End.Before(event.Start) {
            t.Errorf("Event[%d] end time %v should be after start time %v", i, event.End, event.Start)
        }
        
        if event.ID == "" {
            t.Errorf("Event[%d] ID should not be empty", i)
        }
        
        // Check for fallback tag
        hasFallbackTag := false
        for _, tag := range event.Tags {
            if tag == "fallback" {
                hasFallbackTag = true
                break
            }
        }
        if !hasFallbackTag {
            t.Errorf("Event[%d] should have 'fallback' tag", i)
        }
    }
    
    // Check that events are spread across multiple days
    uniqueDays := make(map[string]bool)
    for _, event := range events {
        day := event.Start.Format("2006-01-02")
        uniqueDays[day] = true
    }
    
    if len(uniqueDays) < 3 {
        t.Errorf("Should have events across at least 3 days, got %d", len(uniqueDays))
    }
}
