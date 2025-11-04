package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mhmdhalawi/aida-go/models"
)

// Helper to create temp JSON test data
func writeJSONFile(t *testing.T, dir, name string, data any) {
	t.Helper()
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed writing %s: %v", path, err)
	}
}

// Mock models.User slice for testing
func mockUsers() []models.User {
	return []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Address: "123 Street", PhoneNumber: "+111", Birthday: time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC)},
		{ID: 2, FirstName: "Emma", LastName: "Smith", Address: "456 Road", PhoneNumber: "+222", Birthday: time.Date(1988, 7, 5, 0, 0, 0, 0, time.UTC)},
		{ID: 3, FirstName: "Ali", LastName: "Khan", Address: "789 Ave", PhoneNumber: "+333", Birthday: time.Date(2001, 2, 15, 0, 0, 0, 0, time.UTC)},
		{ID: 4, FirstName: "Nora", LastName: "Lee", Address: "321 Blvd", PhoneNumber: "+444", Birthday: time.Date(1995, 12, 30, 0, 0, 0, 0, time.UTC)},
		{ID: 5, FirstName: "Liam", LastName: "Brown", Address: "654 Drive", PhoneNumber: "+555", Birthday: time.Date(2005, 3, 20, 0, 0, 0, 0, time.UTC)},
		{ID: 6, FirstName: "Olivia", LastName: "Miller", Address: "111 Lane", PhoneNumber: "+666", Birthday: time.Date(1999, 8, 1, 0, 0, 0, 0, time.UTC)},
	}
}

func TestParseDateRange(t *testing.T) {
	tests := []struct {
		input        string
		expectedYear int
		expectOk     bool
	}{
		{"2006", 2006, true},
		{"2006-05", 2006, true},
		{"2006-05-01", 2006, true},
		{"2006-01-01 to 2008-01-01", 2006, true},
		{"", 0, false},
		{"invalid", 0, false},
	}

	for _, tt := range tests {
		start, end, ok := ParseDateRange(tt.input)
		if ok != tt.expectOk {
			t.Errorf("expected ok=%v for %q, got %v", tt.expectOk, tt.input, ok)
		}
		if ok && start.Year() != tt.expectedYear {
			t.Errorf("expected start year %d, got %d", tt.expectedYear, start.Year())
		}
		if ok && end.Before(start) {
			t.Errorf("end %v should not be before start %v", end, start)
		}
	}
}

func TestFilterUsers_ByName(t *testing.T) {
	users := mockUsers()

	filtered := filterUsers(users, "emma", nil, nil)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 result for 'emma', got %d", len(filtered))
	}
	if filtered[0].FirstName != "Emma" {
		t.Errorf("expected 'Emma', got %s", filtered[0].FirstName)
	}
}

func TestFilterUsers_ByDateRange(t *testing.T) {
	users := mockUsers()
	start := time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	filtered := filterUsers(users, "", &start, &end)
	for _, u := range filtered {
		if u.Birthday.Before(start) || u.Birthday.After(end) {
			t.Errorf("user %v has birthday outside range", u)
		}
	}
}

func TestUsersHandler_Pagination(t *testing.T) {
	// Create a temporary root directory
	root := t.TempDir()

	// Create a "data" subfolder so Users() finds "./data"
	dataDir := filepath.Join(root, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatalf("failed to create data dir: %v", err)
	}

	// Write JSON test data to ./data/people.json
	writeJSONFile(t, dataDir, "people.json", mockUsers())

	// Change working directory so "./data" is found
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	if err := os.Chdir(root); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	// Create test request and response
	req := httptest.NewRequest("GET", "/users?cursor=0", nil)
	w := httptest.NewRecorder()

	handler := Users()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Users      []models.User `json:"users"`
		NextCursor int           `json:"nextCursor"`
		Total      int           `json:"total"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Users) != 5 {
		t.Errorf("expected 5 users on first page, got %d", len(resp.Users))
	}
	if resp.NextCursor != 5 {
		t.Errorf("expected nextCursor=5, got %d", resp.NextCursor)
	}
	if resp.Total != len(mockUsers()) {
		t.Errorf("expected total=%d, got %d", len(mockUsers()), resp.Total)
	}
}

func TestUsersHandler_LastPage(t *testing.T) {
	// Create a temporary data folder with mock users
	dir := t.TempDir()
	dataPath := filepath.Join(dir, "data")
	if err := os.Mkdir(dataPath, 0755); err != nil {
		t.Fatalf("failed to create data folder: %v", err)
	}

	writeJSONFile(t, dataPath, "users.json", mockUsers())

	// Temporarily change working directory so Users() loads from ./data
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(dir)

	req := httptest.NewRequest("GET", "/users?cursor=5", nil)
	w := httptest.NewRecorder()
	Users().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		NextCursor int `json:"nextCursor"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.NextCursor != -1 {
		t.Errorf("expected nextCursor=-1 on last page, got %v", resp.NextCursor)
	}
}
