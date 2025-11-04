package loader

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed writing %s: %v", path, err)
	}
	return path
}

func TestLoadFromFolder_BasicAndSkips(t *testing.T) {
	dir := t.TempDir()

	// Valid JSON with two valid and one incomplete user
	writeFile(t, dir, "valid-1.json", `[
		{"id":1,"first_name":"John","last_name":"Doe","birthday":"1990-05-14T00:00:00Z","address":"123","phone_number":"+1"},
		{"id":2,"first_name":"Emma","last_name":"Smith","birthday":"1988-11-02T00:00:00Z","address":"456","phone_number":"+2"},
		{"id":0,"first_name":"","last_name":"","birthday":"0001-01-01T00:00:00Z","address":"","phone_number":""}
	]`)

	// Malformed JSON should be skipped (logged), not fail the whole run
	writeFile(t, dir, "broken.json", `[{"id":3,`) // invalid JSON

	// Non-JSON file should be ignored
	writeFile(t, dir, "note.txt", "hello")

	// JSON in subdirectory
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatalf("failed creating subdir: %v", err)
	}
	writeFile(t, sub, "valid-2.json", `[{"id":4,"first_name":"Ava","last_name":"Jones","birthday":"1993-12-25T00:00:00Z","address":"789","phone_number":"+3"}]`)

	users, err := LoadFromFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users (2 + 1 from subdir), got %d", len(users))
	}

	// Spot-check fields and time parsing
	// birthdays must be parsed to non-zero times
	for i, u := range users {
		if u.Birthday.IsZero() {
			t.Fatalf("user %d has zero birthday: %+v", i, u)
		}
	}

	// Ensure a specific expected user exists
	found := false
	for _, u := range users {
		if u.ID == 2 && u.FirstName == "Emma" && u.LastName == "Smith" {
			parsed, _ := time.Parse(time.RFC3339, "1988-11-02T00:00:00Z")
			if !u.Birthday.Equal(parsed) {
				t.Fatalf("unexpected birthday for Emma: %v", u.Birthday)
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected Emma (id=2) in results, got %#v", users)
	}
}

func TestLoadFromFolder_EmptyFolder(t *testing.T) {
	dir := t.TempDir()
	users, err := LoadFromFolder(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestLoadFromFolder_MissingFolderReturnsError(t *testing.T) {
	// Use a path that does not exist
	nonexistent := filepath.Join(t.TempDir(), "does-not-exist")
	_, err := LoadFromFolder(nonexistent)
	if err == nil {
		t.Fatalf("expected error for missing folder, got nil")
	}
}
