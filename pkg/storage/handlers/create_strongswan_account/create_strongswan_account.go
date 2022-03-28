package create_strongswan_account

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CreateStrongswanAccountHandler struct {
	db *sqlx.DB
}

func NewCreateStrongswanAccountHandler(db *sqlx.DB) *CreateStrongswanAccountHandler {
	return &CreateStrongswanAccountHandler{db: db}
}

func (h CreateStrongswanAccountHandler) Exec(ctx context.Context, args *storage.CreateStrongswanAccount) (err error) {

	if args.User == nil {
		return errors.New("No *User")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	var last_id_in_identities, last_id_in_shared_secrets int64
	var res sql.Result

	// 1
	sqlQuery := `insert into identities (type, data) VALUES (?,?);`
	res, err = tx.Exec(sqlQuery, 2, args.User.GetLogin())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}
	last_id_in_identities, err = res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}

	//2
	sqlQuery2 := "insert into shared_secrets (type,data) VALUES (?,?);"
	res, err = tx.Exec(sqlQuery2, 2, args.User.GetPassword())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}
	last_id_in_shared_secrets, err = res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}

	//3
	sqlQuery3 := "insert into shared_secret_identity (shared_secret, identity) VALUES (?,?);"
	res, err = tx.Exec(sqlQuery3, last_id_in_shared_secrets, last_id_in_identities)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return err
}

func (h *CreateStrongswanAccountHandler) Context() interface{} {
	return (*storage.CreateStrongswanAccount)(nil)
}
