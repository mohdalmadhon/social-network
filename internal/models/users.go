package models

import "time"

type UserRegistration struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	About     string
	DOB       time.Time
	Password  string
	Avatar    string
}
