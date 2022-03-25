package storage

import (
	"backend-vpn/internal/dto"
)

type UserQuery struct {
	login     string
	referalId string
	Out       struct {
		User    *dto.User
		Created bool
	}
}

func NewUserQuery(login, referalId string) *UserQuery {

	u := &dto.User{}

	uq := &UserQuery{
		login:     login,
		referalId: referalId,
		Out: struct {
			User    *dto.User
			Created bool
		}{User: u, Created: false},
	}

	return uq
}

func (u *UserQuery) GetLogin() string {
	return u.login
}

func (u *UserQuery) GetReferalId() string {
	return u.referalId
}
