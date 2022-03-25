package dto

import "time"

type User struct {
	Login     string
	Password  string
	CreatedAt time.Time
	ExpiredAt time.Time
	ReferalId string
}

func NewUser(login string) *User {
	return &User{Login: login}
}
