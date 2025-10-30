package main

import (
	"log"
	"net/http"

	"github.com/mhmdhalawi/aida-go/middleware"
	"github.com/mhmdhalawi/aida-go/routes"
	"github.com/mhmdhalawi/aida-go/utils"
)

func main() {
	users, err := utils.LoadFromFolder("./data")
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	http.Handle("GET /users", middleware.WithJSONHeaders(routes.Users(users)))

	log.Printf("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
