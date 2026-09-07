package models

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


