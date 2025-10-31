package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mhmdhalawi/aida-go/models"
	"github.com/mhmdhalawi/aida-go/utils"
)

// Users returns an http.Handler that writes the users as JSON.
func Users() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := utils.LoadFromFolder("./data")
		if err != nil {
			log.Fatalf("Failed to load users: %v", err)
		}

		// Query params: q (free text for first/last name or date range), from, to (date range)
		query := strings.TrimSpace(r.URL.Query().Get("q"))

		var fromDatePtr, toDatePtr *time.Time

		// If q contains a recognizable range, use it to set from/to and ignore name search.

		if qFrom, qTo, ok := ParseDateRange(query); ok {
			fromDatePtr = &qFrom
			toDatePtr = &qTo

		}

		filteredUsers := filterUsers(users, query, fromDatePtr, toDatePtr)

		cursorStr := strings.TrimSpace(r.URL.Query().Get("cursor"))
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil {
			cursor = 0
		}

		max := 5

		start := min(cursor, len(filteredUsers))

		end := min(start+max, len(filteredUsers))

		page := filteredUsers[start:end]

		// Next cursor
		nextCursor := end
		if nextCursor >= len(filteredUsers) {
			nextCursor = -1 // indicate no more pages
		}

		// Response
		resp := map[string]any{
			"users":      page,
			"nextCursor": nextCursor,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

	})

}

func filterUsers(users []models.User, query string, fromDatePtr, toDatePtr *time.Time) []models.User {
	q := strings.ToLower(query)

	filteredUsers := make([]models.User, 0)

	for _, u := range users {
		name := strings.ToLower(u.FirstName + " " + u.LastName)
		birthday := u.Birthday.Truncate(24 * time.Hour)

		if strings.Contains(name, q) {
			filteredUsers = append(filteredUsers, u)
			continue
		}

		if fromDatePtr != nil && toDatePtr != nil {
			if birthday.After(fromDatePtr.Truncate(24*time.Hour)) && birthday.Before(toDatePtr.Truncate(24*time.Hour)) {
				filteredUsers = append(filteredUsers, u)
			}
			continue
		}

	}

	return filteredUsers
}

// ParseDateRange parses a query like "2006" or "2006 to 2010" or "2006-10-01 to 2010-10-01"
// It returns two time.Time values (start, end)
func ParseDateRange(q string) (time.Time, time.Time, bool) {
	q = strings.TrimSpace(q)
	if q == "" {
		return time.Time{}, time.Time{}, false
	}

	// Split by "to" if it exists
	parts := strings.Split(strings.ToLower(q), "to")
	layouts := []string{
		"2006",
		"2006-01",
		"2006-1",
		"2006-01-02",
		"2006-1-2",
	}

	parse := func(s string) (time.Time, error) {
		s = strings.TrimSpace(s)
		var t time.Time
		var err error
		for _, layout := range layouts {
			t, err = time.Parse(layout, s)
			if err == nil {
				return t, nil
			}
		}
		return time.Time{}, err
	}

	if len(parts) == 1 {
		// Single date or year
		start, err := parse(parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, false
		}
		// Assume full year if only year provided
		end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
		return start, end, true
	}

	if len(parts) == 2 {
		start, err1 := parse(parts[0])
		end, err2 := parse(parts[1])
		if err1 != nil {
			return time.Time{}, time.Time{}, false
		}
		if err2 != nil {
			return time.Time{}, time.Time{}, false
		}
		return start, end, true
	}

	return time.Time{}, time.Time{}, false
}
