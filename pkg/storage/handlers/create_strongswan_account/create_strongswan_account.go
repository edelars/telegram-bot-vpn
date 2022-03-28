package create_strongswan_account

import (
	"backend-vpn/pkg/storage"
	"context"
	"github.com/jmoiron/sqlx"
)

type CreateStrongswanAccountHandler struct {
	db *sqlx.DB
}

func NewCreateStrongswanAccountHandler(db *sqlx.DB) *CreateStrongswanAccountHandler {
	return &CreateStrongswanAccountHandler{db: db}
}

func (h CreateStrongswanAccountHandler) Exec(ctx context.Context, args *storage.UserQuery) (err error) {

	return nil
}

func (h *CreateStrongswanAccountHandler) Context() interface{} {
	return (*storage.CreateStrongswanAccount)(nil)
}
