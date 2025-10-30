package main

import (
	"fmt"
	"log"

	"github.com/mhmdhalawi/aida-go/pkg/user"
)

func main() {
	users, err := user.LoadFromFolder("./data")
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	fmt.Printf("Loaded %d valid entries:\n", len(users))
	for _, u := range users {
		fmt.Printf("%s %s, Birthday: %s, Phone: %s\n",
			u.FirstName, u.LastName, u.Birthday.Format("2006-01-02"), u.PhoneNumber)
	}
}
