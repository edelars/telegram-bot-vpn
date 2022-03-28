package delete_strongswan_account

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type DeleteStrongswanAccountHandler struct {
	db *sqlx.DB
}

func NewDeleteStrongswanAccountHandler(db *sqlx.DB) *DeleteStrongswanAccountHandler {
	return &DeleteStrongswanAccountHandler{db: db}
}

func (h DeleteStrongswanAccountHandler) Exec(ctx context.Context, args *storage.DeleteStrongswanAccount) (err error) {

	if args.User == nil {
		return errors.New("No *User")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	var rows, rows2 *sqlx.Rows

	// 1 select id identities
	sqlQuery := `select id from identities where data = ?;`
	rows, err = tx.Queryx(sqlQuery, args.User.GetEncodedLogin())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get user id %s: %w", args.User.GetLogin(), err)
	}
	defer rows.Close()
	var identitiesId int64
	if rows.Next() {
		err = rows.Scan(&identitiesId)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to get user id %s: %w", args.User.GetLogin(), err)
		}
	} else {
		return errors.New("no such user")
	}

	//2 delete identities
	sqlQuery2 := "delete from identities where id = ?"
	_, err = tx.Exec(sqlQuery2, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user identities %s: %w", args.User.GetLogin(), err)
	}

	//3 get shared_secret id
	sqlQuery3 := "select shared_secret from shared_secret_identity where identity = ?;"
	rows2, err = tx.Queryx(sqlQuery3, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get shared_secret user %s: %w", args.User.GetLogin(), err)
	}
	defer rows2.Close()
	var sharedSecretId int64
	for rows2.Next() {
		err = rows2.Scan(&sharedSecretId)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to get sharedSecretId %s: %w", args.User.GetLogin(), err)
		}
	}

	//4 delete shared_secret_identity
	sqlQuery4 := "delete from shared_secret_identity where identity = ?"
	_, err = tx.Exec(sqlQuery4, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user shared_secret_identity %s: %w", args.User.GetLogin(), err)
	}

	//5 delete shared_secret
	sqlQuery5 := "delete from shared_secrets where id = ?"
	_, err = tx.Exec(sqlQuery5, sharedSecretId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user shared_secret %s: %w", args.User.GetLogin(), err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return err
}

func (h *DeleteStrongswanAccountHandler) Context() interface{} {
	return (*storage.DeleteStrongswanAccount)(nil)
}
