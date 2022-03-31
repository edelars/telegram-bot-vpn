package dto

import "time"

type User struct {
	Login          string
	Password       string
	CreatedAt      time.Time
	ExpiredAt      time.Time
	ReferalId      string
	Psk            string
	UsedTestPeriod bool
	Id             int64
}
