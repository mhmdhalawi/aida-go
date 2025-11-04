package loader

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mhmdhalawi/aida-go/models"
)

// LoadFromFolder reads all JSON files in a folder
func LoadFromFolder(folder string) ([]models.User, error) {
	var users []models.User

	// Ensure the folder exists and is a directory
	info, statErr := os.Stat(folder)
	if statErr != nil {
		return nil, fmt.Errorf("folder does not exist: %w", statErr)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", folder)
	}

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}

		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		log.Printf("Checking file: %s", path)

		fileContent, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read %s: %v", path, err)
			return nil
		}

		var filePeople []models.User
		if err := json.Unmarshal(fileContent, &filePeople); err != nil {
			log.Printf("Failed to parse %s: %v", path, err)
			return nil
		}

		for _, u := range filePeople {
			if u.ID == 0 || u.FirstName == "" || u.LastName == "" || u.Address == "" || u.PhoneNumber == "" || u.Birthday.IsZero() {
				log.Printf("Skipping incomplete entry in %s: %+v", path, u)
				continue
			}
			users = append(users, u)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking folder: %w", err)
	}

	return users, nil
}
