package models

import "time"

type User struct {
	ID          int
	Username    string
	Password    string
	Name        *string
	Nationality *string
	BirthDate   *time.Time
}
