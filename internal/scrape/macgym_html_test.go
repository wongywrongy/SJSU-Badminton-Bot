package scrape

import (
	"context"
	"io"
	"log/slog"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/sjsu-badminton/badminton-discord-bot/internal/config"
	"github.com/sjsu-badminton/badminton-discord-bot/internal/util"
)

func TestMacGymHTMLParsing(t *testing.T) {
	// Set up test logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Testing HTML parsing with Mac Gym URL: %s", cfg.MacGymURL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fetch the raw HTML response
	r, err := util.Get(ctx, cfg.MacGymURL)
	if err != nil {
		t.Fatalf("Failed to fetch Mac Gym data: %v", err)
	}
	defer r.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(bodyBytes)
	t.Logf("=== RAW HTML RESPONSE ===")
	t.Logf("Content-Type: %s", r.Header.Get("Content-Type"))
	t.Logf("Response length: %d bytes", len(bodyStr))
	t.Logf("First 1000 characters: %s", bodyStr[:min(1000, len(bodyStr))])
	// Search for specific patterns in the HTML
	t.Logf("=== SEARCHING FOR PATTERNS ===")
	
	// Look for any numbers that might be counts
	countRegex := regexp.MustCompile(`\b(\d+)\b`)
	countMatches := countRegex.FindAllString(bodyStr, -1)
	t.Logf("All numbers found: %v", countMatches)
	
	// Look for time patterns
	timeRegex := regexp.MustCompile(`\d{2}/\d{2}/\d{4}\s+\d{2}:\d{2}\s+[AP]M`)
	timeMatches := timeRegex.FindAllString(bodyStr, -1)
	t.Logf("Time patterns found: %v", timeMatches)
	
	// Look for "MAC Gym" or similar
	macGymRegex := regexp.MustCompile(`(?i)(mac\s+gym|gym)`)
	macGymMatches := macGymRegex.FindAllString(bodyStr, -1)
	t.Logf("MAC Gym references found: %v", macGymMatches)
	
	// Look for JavaScript data
	jsDataRegex := regexp.MustCompile(`LocationId.*?(\d+)`)
	jsMatches := jsDataRegex.FindAllStringSubmatch(bodyStr, -1)
	t.Logf("Location IDs found: %v", jsMatches)
	t.Logf("Last 1000 characters: %s", bodyStr[max(0, len(bodyStr)-1000):])
	// Try to parse the HTML for the data we need
	status, count, updated := parseMacGymHTML(bodyStr)
	
	t.Logf("=== PARSED DATA ===")
	t.Logf("Status: %s", status)
	t.Logf("Count: %s", count)
	t.Logf("Updated: %s", updated)

	// Validate that we found the expected data
	if status == "" {
		t.Error("Could not find status in HTML")
	}
	if count == "" {
		t.Error("Could not find count in HTML")
	}
	if updated == "" {
		t.Error("Could not find updated time in HTML")
	}

	// Try to convert count to integer
	if count != "" {
		// Extract just the number from the count
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(count, -1)
		if len(matches) > 0 {
			t.Logf("Extracted count number: %s", matches[0])
		}
	}
}

func parseMacGymHTML(html string) (status, count, updated string) {
	// Look for patterns in the HTML that match the expected format:
	// MAC Gym (Open)
	// Last Count: 10
	// Updated: 09/09/2025 04:06 PM

	// Try to find the status (Open/Closed)
	statusRegex := regexp.MustCompile(`(?i)(open|closed)`)
	if matches := statusRegex.FindStringSubmatch(html); len(matches) > 0 {
		status = matches[1]
	}

	// Try to find the count
	countRegex := regexp.MustCompile(`(?i)(?:last\s+count|count)[:\s]*(\d+)`)
	if matches := countRegex.FindStringSubmatch(html); len(matches) > 0 {
		count = matches[1]
	}

	// Try to find the updated time
	updatedRegex := regexp.MustCompile(`(?i)(?:updated|last\s+updated)[:\s]*(\d{2}/\d{2}/\d{4}\s+\d{2}:\d{2}\s+[AP]M)`)
	if matches := updatedRegex.FindStringSubmatch(html); len(matches) > 0 {
		updated = matches[1]
	}

	return status, count, updated
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
