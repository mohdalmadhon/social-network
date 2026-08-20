package models

import (
	"time"
)

type User struct {
	Id         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"createdAt"`
	DOB        time.Time `json:"dob"`
	Updated_at time.Time `json:"updatedAt"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DOB       string `json:"dob"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	About     string `json:"about"`
}

type Profile struct {
	UserID      int
	About       string
	Avatar_Path string
	Followers   int
	Following   int
	Posts       int
}

type Logger struct {
	Identifier string `json:"identifier"`
	Pass       string `json:"password"`
}
