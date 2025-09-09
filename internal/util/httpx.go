package util

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "time"
)

var client = &http.Client{ Timeout: 10 * time.Second }

type Doer interface{ Do(*http.Request) (*http.Response, error) }

func Get(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, fmt.Errorf("creating request: %w", err)
    }
    
    req.Header.Set("User-Agent", "sjsu-badminton-bot/1.0")
    req.Header.Set("Accept", "application/json, text/html, */*")

    var resp *http.Response
    backoff := 250 * time.Millisecond
    
    for i := 0; i < 3; i++ {
        resp, err = client.Do(req)
        if err == nil && resp.StatusCode < 500 { 
            break 
        }
        
        if resp != nil {
            resp.Body.Close()
        }
        
        slog.Warn("HTTP request failed, retrying", 
            "attempt", i+1, 
            "url", url, 
            "error", err,
            "status", func() int {
                if resp != nil { return resp.StatusCode }
                return 0
            }())
        
        time.Sleep(backoff)
        backoff *= 2
    }
    
    if err != nil { 
        return nil, fmt.Errorf("request failed after retries: %w", err) 
    }
    
    if resp.StatusCode >= 400 { 
        resp.Body.Close()
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status) 
    }
    
    return resp, nil
}

func DecodeJSON(r io.Reader, v any) error {
    if err := json.NewDecoder(r).Decode(v); err != nil {
        return fmt.Errorf("JSON decode: %w", err)
    }
    return nil
}
