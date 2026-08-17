package models

import (
	"time"
)

type User struct {
	Id                                             int
	FirstName, LastName, Email, Password, Username string
	CreatedAt, DOB, Updated_at                     time.Time
}
