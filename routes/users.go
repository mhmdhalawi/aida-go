package routes

import (
	"encoding/json"
	"net/http"

	"github.com/mhmdhalawi/aida-go/models"
)

// Users returns an http.Handler that writes the users as JSON.
func Users(users []models.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
