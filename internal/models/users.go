package models

import (
	"time"
)

type UserRegistration struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	About     string
	DOB       time.Time
	Password  string
	Avatar    string
	IsPrivate int
}

type UserLogger struct {
	Identifier string
	Pass       string
}

type UserData struct {
	UserInfo       UserRegistration
	NumOfFollowers int
	NumOfFollowing int
	NumOfPosts     int
	About          UserAbout
	Followers      map[int]UserRegistration
	Following      map[int]UserRegistration
	IsPrivate      int
}

type UserAbout struct {
	Bio       string
	Work      string `json:"Work"`
	Education string `json:"Education"`
	Travel    string `json:"Travel"`
	Intrests  string `json:"Intrests"`
	Hobbies   string `json:"Hobbies"`
	Website   string `json:"Website"`
	Linkedin  string `json:"Linkedin"`
	Instgram  string `json:"Instgram"`
	Twitter   string `json:"Twitter"`
}

type UserProfile struct {
	About    UserAbout
	UserInfo UserRegistration
}
