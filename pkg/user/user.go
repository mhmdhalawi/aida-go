package user

import "time"

type User struct {
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Birthday    time.Time `json:"birthday"`
	Address     string   `json:"address"`
	PhoneNumber string   `json:"phoneNumber"`
}


