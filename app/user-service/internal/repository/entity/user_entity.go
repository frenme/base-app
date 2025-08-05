package entity

import (
	"time"
)

type UserEntity struct {
	ID          int
	Username    string
	Password    string
	Name        *string
	Nationality *string
	BirthDate   *time.Time
}

func (UserEntity) TableName() string {
	return "users"
}
