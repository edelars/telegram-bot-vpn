package storage

import (
	"backend-vpn/internal/dto"
	"errors"
	"math/rand"
	"strings"
	"time"
)

type UserQuery struct {
	login     string
	tg_id     int64
	referalId string
	inviteid  string
	Out       struct {
		User    *dto.User
		Created bool
	}
}

func NewUserQuery(login string, id int64, referalId string) *UserQuery {

	u := &dto.User{}

	uq := &UserQuery{
		login:     login,
		tg_id:     id,
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

func (u *UserQuery) GetId() int64 {
	return u.tg_id
}
func (u *UserQuery) GetReferalId() string {
	return u.inviteid
}

func (u *UserQuery) GetNewReferalId() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 6
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	u.inviteid = b.String()

	return u.inviteid
}
func (u *UserQuery) GetInviteReferalId() string {
	return u.referalId
}

type SaveUserQuery struct {
	user *dto.User
}

func NewSaveUserQuery(user *dto.User) (err error, u *SaveUserQuery) {
	if user == nil {
		return errors.New("user is nil"), u
	}
	if user.Login == "" {
		return errors.New("login is nil"), u
	}
	return err, &SaveUserQuery{user: user}
}

func (h *SaveUserQuery) GetUser() *dto.User {
	return h.user
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

type NewPayments struct {
	UserId int64
	Value  int
}

type UncheckedPaymentsQuery struct {
	UserId int64
	Out    struct {
		IncBalance int64
	}
}

type GetUserBalanceQuery struct {
	UserId int64
	Out    struct {
		TotalBalance int64
	}
}

type GetUser struct {
	UserId int64
	Out    struct {
		User *dto.User
	}
}

func NewGetUser(userId int64) *GetUser {
	u := &dto.User{}
	return &GetUser{
		UserId: userId,
		Out:    struct{ User *dto.User }{User: u},
	}
}

type WriteOffBalance struct {
	UserId int64
	Value  int
}

type GetExpiredUsers struct {
	Out []*dto.User
}
