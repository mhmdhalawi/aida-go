package user

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// LoadFromFolder reads all JSON files in a folder
func LoadFromFolder(folder string) ([]User, error) {
	var people []User

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read %s: %v", path, err)
			return nil
		}

		var filePeople []User
		if err := json.Unmarshal(fileContent, &filePeople); err != nil {
			log.Printf("Failed to parse %s: %v", path, err)
			return nil
		}

		for _, p := range filePeople {
			if p.FirstName == "" || p.LastName == "" || p.Address == "" || p.PhoneNumber == "" || p.Birthday.IsZero() {
				log.Printf("Skipping incomplete entry in %s: %+v", path, p)
				continue
			}
			people = append(people, p)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking folder: %w", err)
	}

	return people, nil
}
