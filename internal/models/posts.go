package models

import "time"

// group id == 0 (PUBLIC)
// group id == -1 (PRIVATE)
type RegsiterPost struct {
	UserID        int
	GroupID       int    `json:"groupID"`
	Content       string `json:"content"`
	Image_path    string
	AllowComments int    `json:"allowComments"`
	Location      string `json:"location"`
	PeopleTagged  []int  `json:"taggedPeople"`
}

type TaggedPerson struct {
	Id         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	AvatarPath string `json:"avatarPath"`
}

type Post struct {
	Id             int            `json:"id"`
	UserId         int            `json:"userId"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	Username       *string        `json:"username"`
	AvatarPath     string         `json:"avatarPath"`
	Content        string         `json:"content"`
	ImagePath      *string        `json:"imagePath"`
	AllowComments  bool           `json:"allowComments"`
	Location       *string        `json:"location"`
	GroupId        *int           `json:"groupId"`
	CreatedAt      string         `json:"createdAt"`
	Relationship   string         `json:"relationship"`
	Visibility     string         `json:"visibility"`
	VisibilityUser string         `json:"visibilityUser"`
	TaggedPeople   []TaggedPerson `json:"taggedPeople"`
}

type Comment struct {
	ID                      int
	Content                 string
	LikeCount, DisLikeCount int
	User                    UserRegistration
	CreatedAt               time.Time
}
