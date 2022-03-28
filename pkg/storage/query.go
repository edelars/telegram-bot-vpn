package storage

import (
	"backend-vpn/internal/dto"
	"strconv"
)

type UserQuery struct {
	login     string
	referalId string
	Out       struct {
		User    *dto.User
		Created bool
	}
}

func NewUserQuery(id int64, referalId string) *UserQuery {

	u := &dto.User{}

	uq := &UserQuery{
		login:     strconv.FormatInt(id, 10),
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

type CreateStrongswanAccount struct {
	User *dto.StrongswanUser
}

type DeleteStrongswanAccount struct {
	User *dto.StrongswanUser
}

type AccessRightQuery struct {
	Id int64
}
