package storage

import (
	"backend-vpn/internal/dto"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
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

type CreateStrongswanAccount struct {
	login    string
	password string
}

func NewCreateStrongswanAccount(login string, pass string, generateNewPass bool) *CreateStrongswanAccount {

	p := &CreateStrongswanAccount{
		login: login,
	}

	if generateNewPass {
		p.password = p.generateNewPassword()
	} else {
		p.password = pass
	}

	return p
}

// GetEncodedPassword  return X'717765727479'  == qwerty hex
func (a *CreateStrongswanAccount) GetEncodedPassword() string {
	hx := hex.EncodeToString([]byte(a.password))
	return fmt.Sprintf("X'%s'", hx)
}

func (a *CreateStrongswanAccount) generateNewPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
