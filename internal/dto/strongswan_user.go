package dto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type StrongswanUser struct {
	login    string
	password string
}

func NewStrongswanUser(login string, pass string, generateNewPass bool) (err error, p *StrongswanUser) {

	p = &StrongswanUser{
		login: login,
	}

	if login == "" {
		return errors.New("login cant be empty"), nil
	}
	//TODO check @ at login

	if generateNewPass {
		p.password = p.generateNewPassword()
	} else {
		p.password = pass
	}

	return err, p
}

// GetEncodedPassword  return X'717765727479'  == qwerty hex
func (a *StrongswanUser) GetEncodedPassword() string {
	return a.encode(a.password)
}

func (a *StrongswanUser) GetEncodedLogin() string {
	return a.encode(a.login)
}

func (a *StrongswanUser) GetLogin() string {
	return a.login
}

func (a *StrongswanUser) GetPassword() string {
	return a.password
}

func (a *StrongswanUser) encode(s string) string {
	hx := hex.EncodeToString([]byte(s))
	return fmt.Sprintf("%s", hx)
}

func (a *StrongswanUser) generateNewPassword() string {
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
