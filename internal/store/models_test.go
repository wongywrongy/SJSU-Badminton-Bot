package store

import (
    "testing"
    "time"
)

func TestHashKey(t *testing.T) {
    title := "Test Event"
    start := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)
    end := time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC)
    location := "Test Location"

    key1 := HashKey(title, start, end, location)
    key2 := HashKey(title, start, end, location)

    if key1 != key2 {
        t.Errorf("HashKey should be deterministic, got different keys: %s vs %s", key1, key2)
    }

    if len(key1) != 40 { // SHA1 hex length
        t.Errorf("Expected hash key length 40, got %d", len(key1))
    }
}

func TestMemoryStore(t *testing.T) {
    store := NewMemoryStore()

    // Test initial state
    if store.GetEventCount() != 0 {
        t.Errorf("Expected 0 events initially, got %d", store.GetEventCount())
    }

    if store.GetSubscriberCount() != 0 {
        t.Errorf("Expected 0 subscribers initially, got %d", store.GetSubscriberCount())
    }

    // Test Mac Gym snapshot
    snap := MacGymSnapshot{
        RetrievedAt: time.Now(),
        Location:    "Test Gym",
        Capacity:    10,
        InUse:       5,
        Details:     "Test details",
    }

    store.SetMac(snap)
    retrieved := store.GetMac()

    if retrieved.Location != snap.Location {
        t.Errorf("Expected location %s, got %s", snap.Location, retrieved.Location)
    }

    if retrieved.Capacity != snap.Capacity {
        t.Errorf("Expected capacity %d, got %d", snap.Capacity, retrieved.Capacity)
    }

    // Test events
    event := Event{
        ID:        "test-event-1",
        Title:     "Test Badminton Event",
        Location:  "Test Court",
        Start:     time.Now().Add(time.Hour),
        End:       time.Now().Add(2 * time.Hour),
        SourceURL: "https://test.com",
        Tags:      []string{"badminton"},
    }

    store.UpsertEvents([]Event{event})

    if store.GetEventCount() != 1 {
        t.Errorf("Expected 1 event, got %d", store.GetEventCount())
    }

    upcoming := store.ListUpcoming(time.Now(), 1)
    if len(upcoming) != 1 {
        t.Errorf("Expected 1 upcoming event, got %d", len(upcoming))
    }

    // Test subscriptions
    store.Subscribe("user123", 5)
    if store.GetSubscriberCount() != 1 {
        t.Errorf("Expected 1 subscriber, got %d", store.GetSubscriberCount())
    }

    subs := store.Subscribers()
    if subs["user123"] != 5 {
        t.Errorf("Expected threshold 5 for user123, got %d", subs["user123"])
    }

    store.Unsubscribe("user123")
    if store.GetSubscriberCount() != 0 {
        t.Errorf("Expected 0 subscribers after unsubscribe, got %d", store.GetSubscriberCount())
    }
}

func TestEventDeduplication(t *testing.T) {
    store := NewMemoryStore()

    event1 := Event{
        ID:        "same-event",
        Title:     "Same Event",
        Location:  "Same Location",
        Start:     time.Now().Add(time.Hour),
        End:       time.Now().Add(2 * time.Hour),
        SourceURL: "https://test.com",
        Tags:      []string{"badminton"},
    }

    event2 := Event{
        ID:        "same-event", // Same ID
        Title:     "Same Event",
        Location:  "Same Location",
        Start:     time.Now().Add(time.Hour),
        End:       time.Now().Add(2 * time.Hour),
        SourceURL: "https://test.com",
        Tags:      []string{"badminton"},
    }

    // Add same event twice
    store.UpsertEvents([]Event{event1})
    store.UpsertEvents([]Event{event2})

    // Should still be only 1 event
    if store.GetEventCount() != 1 {
        t.Errorf("Expected 1 event after deduplication, got %d", store.GetEventCount())
    }
}
