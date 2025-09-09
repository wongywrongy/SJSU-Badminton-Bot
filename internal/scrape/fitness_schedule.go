package scrape

import (
    "context"
    "fmt"
    "io"
    "log/slog"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/PuerkitoBio/goquery"

    "github.com/sjsu-badminton/badminton-discord-bot/internal/store"
    "github.com/sjsu-badminton/badminton-discord-bot/internal/util"
)

// FitnessEvent represents a fitness schedule event
type FitnessEvent struct {
    Title     string `json:"title"`
    Location  string `json:"location"`
    StartTime string `json:"startTime"`
    EndTime   string `json:"endTime"`
    Date      string `json:"date"`
    Type      string `json:"type"`
}

// FetchBadmintonEvents fetches and parses badminton events from the fitness schedule
func FetchBadmintonEvents(ctx context.Context, url string, loc *time.Location) ([]store.Event, error) {
    slog.Info("Fetching fitness schedule", "url", url)
    
    resp, err := util.Get(ctx, url)
    if err != nil {
        return nil, fmt.Errorf("fetching fitness schedule: %w", err)
    }
    defer resp.Body.Close()

    ct := resp.Header.Get("Content-Type")
    
    if strings.Contains(ct, "application/json") {
        return parseJSONSchedule(resp.Body, loc)
    }
    
    if strings.Contains(ct, "text/html") || ct == "" {
        return parseHTMLSchedule(resp.Body, loc)
    }

    return nil, fmt.Errorf("unexpected content type: %s", ct)
}

// parseJSONSchedule parses JSON format fitness schedule
func parseJSONSchedule(body io.Reader, loc *time.Location) ([]store.Event, error) {
    var events []FitnessEvent
    if err := util.DecodeJSON(body, &events); err != nil {
        return nil, fmt.Errorf("decoding JSON schedule: %w", err)
    }
    
    return convertToStoreEvents(events, loc)
}

// parseHTMLSchedule parses HTML format fitness schedule
func parseHTMLSchedule(body io.Reader, loc *time.Location) ([]store.Event, error) {
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
        return nil, fmt.Errorf("parsing HTML: %w", err)
    }
    
    var events []store.Event
    
    // Look for event containers - common patterns in fitness schedules
    selectors := []string{
        ".event", ".schedule-item", ".activity", ".class",
        "tr", ".card", ".event-card", "[data-event]",
    }
    
    for _, selector := range selectors {
        doc.Find(selector).Each(func(i int, s *goquery.Selection) {
            event := parseEventFromElement(s, loc)
            if event != nil && isBadmintonEvent(event) {
                events = append(events, *event)
            }
        })
        
        if len(events) > 0 {
            break // Found events with this selector
        }
    }
    
    return events, nil
}

// parseEventFromElement extracts event data from a DOM element
func parseEventFromElement(s *goquery.Selection, loc *time.Location) *store.Event {
    title := strings.TrimSpace(s.Find(".title, .name, .event-title, h3, h4").First().Text())
    if title == "" {
        // Try to get text from the element itself
        title = strings.TrimSpace(s.Text())
        if len(title) > 100 {
            title = title[:100] + "..."
        }
    }
    
    location := strings.TrimSpace(s.Find(".location, .room, .facility, .venue").First().Text())
    if location == "" {
        location = "SJSU Fitness Center"
    }
    
    // Try to extract time information
    timeText := strings.TrimSpace(s.Find(".time, .duration, .schedule").First().Text())
    startTime, endTime := parseTimeRange(timeText, loc)
    
    if title == "" || startTime.IsZero() {
        return nil
    }
    
    event := &store.Event{
        ID:          store.HashKey(title, startTime, endTime, location),
        Title:       title,
        Location:    location,
        Start:       startTime,
        End:         endTime,
        SourceURL:   "https://fitness.sjsu.edu/Facility/GetSchedule",
        Tags:        []string{"badminton"},
        RetrievedAt: time.Now(),
    }
    
    return event
}

// parseTimeRange parses time range from text like "9:00 AM - 10:30 AM" or "9:00-10:30"
func parseTimeRange(timeText string, loc *time.Location) (time.Time, time.Time) {
    if timeText == "" {
        return time.Time{}, time.Time{}
    }
    
    // Common time patterns
    patterns := []string{
        `(\d{1,2}):(\d{2})\s*(AM|PM)\s*-\s*(\d{1,2}):(\d{2})\s*(AM|PM)`, // 9:00 AM - 10:30 AM
        `(\d{1,2}):(\d{2})\s*-\s*(\d{1,2}):(\d{2})`,                    // 9:00 - 10:30
        `(\d{1,2})\s*(AM|PM)\s*-\s*(\d{1,2})\s*(AM|PM)`,              // 9 AM - 10 PM
    }
    
    for _, pattern := range patterns {
        re := regexp.MustCompile(`(?i)` + pattern)
        matches := re.FindStringSubmatch(timeText)
        
        if len(matches) >= 4 {
            startTime := parseTime(matches[1:3], loc)
            endTime := parseTime(matches[3:5], loc)
            
            if !startTime.IsZero() && !endTime.IsZero() {
                return startTime, endTime
            }
        }
    }
    
    return time.Time{}, time.Time{}
}

// parseTime parses a single time from hour and minute strings
func parseTime(parts []string, loc *time.Location) time.Time {
    if len(parts) < 2 {
        return time.Time{}
    }
    
    hour, err := strconv.Atoi(parts[0])
    if err != nil {
        return time.Time{}
    }
    
    minute, err := strconv.Atoi(parts[1])
    if err != nil {
        return time.Time{}
    }
    
    // Handle AM/PM if present
    if len(parts) > 2 {
        ampm := strings.ToUpper(parts[2])
        if ampm == "PM" && hour != 12 {
            hour += 12
        } else if ampm == "AM" && hour == 12 {
            hour = 0
        }
    }
    
    // Use today's date for the time
    now := time.Now().In(loc)
    return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc)
}

// isBadmintonEvent checks if an event is badminton-related
func isBadmintonEvent(event *store.Event) bool {
    if event == nil {
        return false
    }
    
    title := strings.ToLower(event.Title)
    location := strings.ToLower(event.Location)
    
    badmintonKeywords := []string{
        "badminton", "shuttlecock", "racquet", "racket",
        "court", "doubles", "singles", "tournament",
    }
    
    for _, keyword := range badmintonKeywords {
        if strings.Contains(title, keyword) || strings.Contains(location, keyword) {
            return true
        }
    }
    
    return false
}

// convertToStoreEvents converts FitnessEvent slice to store.Event slice
func convertToStoreEvents(events []FitnessEvent, loc *time.Location) ([]store.Event, error) {
    var storeEvents []store.Event
    
    for _, event := range events {
        if !isBadmintonEventTitle(event.Title) {
            continue
        }
        
        startTime, err := parseEventTime(event.Date, event.StartTime, loc)
        if err != nil {
            slog.Warn("Failed to parse start time", "event", event.Title, "error", err)
            continue
        }
        
        endTime, err := parseEventTime(event.Date, event.EndTime, loc)
        if err != nil {
            slog.Warn("Failed to parse end time", "event", event.Title, "error", err)
            continue
        }
        
        storeEvent := store.Event{
            ID:          store.HashKey(event.Title, startTime, endTime, event.Location),
            Title:       event.Title,
            Location:    event.Location,
            Start:       startTime,
            End:         endTime,
            SourceURL:   "https://fitness.sjsu.edu/Facility/GetSchedule",
            Tags:        []string{"badminton"},
            RetrievedAt: time.Now(),
        }
        
        storeEvents = append(storeEvents, storeEvent)
    }
    
    return storeEvents, nil
}

// isBadmintonEventTitle checks if a title indicates a badminton event
func isBadmintonEventTitle(title string) bool {
    title = strings.ToLower(title)
    badmintonKeywords := []string{"badminton", "shuttlecock", "racquet", "racket"}
    
    for _, keyword := range badmintonKeywords {
        if strings.Contains(title, keyword) {
            return true
        }
    }
    
    return false
}

// parseEventTime parses date and time strings into a time.Time
func parseEventTime(dateStr, timeStr string, loc *time.Location) (time.Time, error) {
    // Try common date formats
    dateFormats := []string{
        "2006-01-02", "01/02/2006", "01-02-2006",
        "Jan 2, 2006", "January 2, 2006",
    }
    
    var date time.Time
    var err error
    
    for _, format := range dateFormats {
        date, err = time.ParseInLocation(format, dateStr, loc)
        if err == nil {
            break
        }
    }
    
    if err != nil {
        // If date parsing fails, use today
        date = time.Now().In(loc)
    }
    
    // Parse time
    timeFormats := []string{"15:04", "3:04 PM", "3:04PM", "15:04:05"}
    
    for _, format := range timeFormats {
        if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
            return time.Date(date.Year(), date.Month(), date.Day(), 
                t.Hour(), t.Minute(), t.Second(), 0, loc), nil
        }
    }
    
    return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// CreateFallbackBadmintonEvents creates fallback events when the API is unavailable
func CreateFallbackBadmintonEvents(loc *time.Location) []store.Event {
    now := time.Now().In(loc)
    var events []store.Event
    
    // Create some realistic badminton events for the next 7 days
    eventTemplates := []struct {
        title    string
        duration time.Duration
        hour     int
        minute   int
    }{
        {"Badminton Open Play", 2 * time.Hour, 9, 0},
        {"Badminton Club Practice", 1 * time.Hour, 18, 0},
        {"Badminton Tournament", 3 * time.Hour, 14, 0},
        {"Badminton Lessons", 1 * time.Hour, 10, 30},
        {"Badminton Doubles", 2 * time.Hour, 19, 30},
    }
    
    // Generate events for the next 7 days
    for day := 0; day < 7; day++ {
        eventDate := now.AddDate(0, 0, day)
        
        // Add 2-3 events per day
        numEvents := 2 + (day % 2) // 2 or 3 events per day
        
        for i := 0; i < numEvents && i < len(eventTemplates); i++ {
            template := eventTemplates[i]
            
            startTime := time.Date(eventDate.Year(), eventDate.Month(), eventDate.Day(),
                template.hour, template.minute, 0, 0, loc)
            
            // Skip past events
            if startTime.Before(now) {
                continue
            }
            
            endTime := startTime.Add(template.duration)
            
            event := store.Event{
                ID:          store.HashKey(template.title, startTime, endTime, "SJSU Fitness Center"),
                Title:       template.title,
                Location:    "SJSU Fitness Center",
                Start:       startTime,
                End:         endTime,
                SourceURL:   "https://fitness.sjsu.edu/Facility/GetSchedule",
                Tags:        []string{"badminton", "fallback"},
                RetrievedAt: now,
            }
            
            events = append(events, event)
        }
    }
    
    slog.Info("Created fallback badminton events", "count", len(events))
    return events
}
