package models

import "time"

type User struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Birthday    time.Time `json:"birthday"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
}
