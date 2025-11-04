package main

import (
	"log"
	"net/http"

	"github.com/mhmdhalawi/aida-go/middleware"
	routes "github.com/mhmdhalawi/aida-go/routes/users"
)

func main() {

	http.Handle("GET /users", middleware.WithHeaders(routes.Users()))

	log.Printf("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
