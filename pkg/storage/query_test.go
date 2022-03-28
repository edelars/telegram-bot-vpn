package storage

import (
	"testing"
)

func TestCreateStrongswanAccount_GetEncodedPassword(t *testing.T) {
	type fields struct {
		login    string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "qwerty",
			fields: struct {
				login    string
				password string
			}{login: "qwerty", password: "qwerty"},
			want: "X'717765727479'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CreateStrongswanAccount{
				login:    tt.fields.login,
				password: tt.fields.password,
			}
			if got := a.GetEncodedPassword(); got != tt.want {
				t.Errorf("GetEncodedPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
