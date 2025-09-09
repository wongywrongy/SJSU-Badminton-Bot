package store

import (
    "crypto/sha1"
    "encoding/hex"
    "fmt"
    "log/slog"
    "sort"
    "strings"
    "sync"
    "time"
)

type MacGymSnapshot struct {
    RetrievedAt time.Time
    Location    string
    Capacity    int
    InUse       int
    Details     string
    Raw         any
}

type Event struct {
    ID          string
    Title       string
    Location    string
    Start       time.Time
    End         time.Time
    SourceURL   string
    Tags        []string
    RetrievedAt time.Time
}

type MemoryStore struct {
    mu         sync.RWMutex
    mac        MacGymSnapshot
    events     map[string]Event
    subs       map[string]int // userID -> threshold
    lastAlert  time.Time      // for debouncing alerts
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        events:    make(map[string]Event),
        subs:      make(map[string]int),
        lastAlert: time.Time{},
    }
}

// HashKey creates a stable deduplication key for events
func HashKey(title string, start, end time.Time, loc string) string {
    key := strings.ToLower(fmt.Sprintf("%s|%s|%s|%s", 
        title, 
        start.UTC().Format(time.RFC3339), 
        end.UTC().Format(time.RFC3339), 
        loc))
    h := sha1.Sum([]byte(key))
    return hex.EncodeToString(h[:])
}

// SetMac updates the Mac Gym snapshot
func (m *MemoryStore) SetMac(s MacGymSnapshot) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    oldSnapshot := m.mac
    m.mac = s
    
    slog.Info("Updated Mac Gym snapshot", 
        "capacity", s.Capacity,
        "inUse", s.InUse,
        "details", s.Details)
    
    // Check for threshold alerts
    m.checkThresholdAlerts(oldSnapshot, s)
}

// GetMac returns a copy of the current Mac Gym snapshot
func (m *MemoryStore) GetMac() MacGymSnapshot {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    // Return a copy to avoid race conditions
    return MacGymSnapshot{
        RetrievedAt: m.mac.RetrievedAt,
        Location:    m.mac.Location,
        Capacity:    m.mac.Capacity,
        InUse:       m.mac.InUse,
        Details:     m.mac.Details,
        Raw:         m.mac.Raw,
    }
}

// UpsertEvents adds or updates events with deduplication
func (m *MemoryStore) UpsertEvents(es []Event) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    added := 0
    updated := 0
    
    for _, e := range es {
        if _, exists := m.events[e.ID]; exists {
            // Update existing event
            m.events[e.ID] = e
            updated++
            slog.Debug("Updated event", "id", e.ID, "title", e.Title)
        } else {
            // Add new event
            m.events[e.ID] = e
            added++
            slog.Info("Added new event", "id", e.ID, "title", e.Title, "start", e.Start)
        }
    }
    
    slog.Info("Upserted events", "added", added, "updated", updated, "total", len(m.events))
}

// ListUpcoming returns upcoming events sorted by start time
func (m *MemoryStore) ListUpcoming(now time.Time, days int) []Event {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    cutoff := now.AddDate(0, 0, days)
    var upcoming []Event
    
    for _, e := range m.events {
        if e.End.After(now) && e.Start.Before(cutoff) {
            upcoming = append(upcoming, e)
        }
    }
    
    // Sort by start time
    sort.Slice(upcoming, func(i, j int) bool {
        return upcoming[i].Start.Before(upcoming[j].Start)
    })
    
    return upcoming
}

// Subscribe adds a user to the alert subscription list
func (m *MemoryStore) Subscribe(userID string, threshold int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.subs[userID] = threshold
    slog.Info("User subscribed to alerts", "userID", userID, "threshold", threshold)
}

// Unsubscribe removes a user from the alert subscription list
func (m *MemoryStore) Unsubscribe(userID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    delete(m.subs, userID)
    slog.Info("User unsubscribed from alerts", "userID", userID)
}

// Subscribers returns a copy of the subscribers map
func (m *MemoryStore) Subscribers() map[string]int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    subs := make(map[string]int, len(m.subs))
    for k, v := range m.subs {
        subs[k] = v
    }
    return subs
}

// checkThresholdAlerts checks if occupancy thresholds have been crossed
func (m *MemoryStore) checkThresholdAlerts(old, new MacGymSnapshot) {
    if new.Capacity == 0 {
        return // No capacity data available
    }
    
    // Debounce alerts (max once per minute)
    if time.Since(m.lastAlert) < time.Minute {
        return
    }
    
    for userID, threshold := range m.subs {
        oldCrossed := old.InUse >= threshold
        newCrossed := new.InUse >= threshold
        
        // Alert if threshold was just crossed
        if !oldCrossed && newCrossed {
            m.lastAlert = time.Now()
            slog.Info("Threshold crossed", 
                "userID", userID, 
                "threshold", threshold, 
                "current", new.InUse, 
                "capacity", new.Capacity)
            // Note: Actual alert sending would be handled by the Discord client
        }
    }
}

// GetEventCount returns the total number of events in the store
func (m *MemoryStore) GetEventCount() int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return len(m.events)
}

// GetSubscriberCount returns the number of active subscribers
func (m *MemoryStore) GetSubscriberCount() int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return len(m.subs)
}
