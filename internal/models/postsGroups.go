package models

type Group struct {
	ID     int
	UserID int
	Name   string
	Users  []UserRegistration
}

type NewGroup struct {
	GroupID int    `json:"groupId"`
	UserID  int    `json:"userId"`
	Name    string `json:"name"`
	Users   []int  `json:"users"`
}
