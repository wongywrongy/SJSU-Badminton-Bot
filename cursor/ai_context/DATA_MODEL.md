# Data Model

```go
// internal/store/models.go
package store

import "time"

// Mac Gym snapshot focused on badminton courts/capacity (best‑effort mapping).
type MacGymSnapshot struct {
    RetrievedAt time.Time
    Location    string   // e.g., "Mac Gym"
    Capacity    int      // if available
    InUse       int      // if available
    Details     string   // free text summary
    Raw         any      // original parsed payload for debugging
}

type Event struct {
    ID         string    // stable dedupe key: hash(title+start+end+location)
    Title      string
    Location   string
    Start      time.Time // in America/Los_Angeles
    End        time.Time
    SourceURL  string
    Tags       []string  // include "badminton" when matched
    RetrievedAt time.Time
}
```

**Dedupe key:** `sha1(strings.ToLower(title+startUTC+endUTC+location))`.
