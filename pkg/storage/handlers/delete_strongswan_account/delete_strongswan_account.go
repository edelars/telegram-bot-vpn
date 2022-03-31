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

var (
	ErrNotExist = errors.New("no such user")
)

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

	var rows, rows2, rows3 *sqlx.Rows
	var res sql.Result

	// 1 select id identities
	sqlQuery := `select id from identities where data = ?;`
	rows, err = tx.Queryx(sqlQuery, args.User.GetLoginByte())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get user id %s: %w", args.User.GetLogin(), err)
	}

	var identitiesId int64
	if rows.Next() {
		err = rows.Scan(&identitiesId)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to scan user id %s: %w", args.User.GetLogin(), err)
		}
		rows.Close()
	} else {
		rows.Close()
		_ = tx.Rollback()
		return ErrNotExist
	}

	//2 delete identities
	sqlQuery2 := `delete from identities where id = ?`
	res, err = tx.ExecContext(ctx, sqlQuery2, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user identities %s: %w", args.User.GetLogin(), err)
	}

	print(res.RowsAffected())

	//3.1  get shared_secret id
	sqlQuery3 := "select shared_secret   from shared_secret_identity where identity= ?;"
	rows3, err = tx.Queryx(sqlQuery3, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get shared_secret user %s: %w", args.User.GetLogin(), err)
	}
	defer rows3.Close()
	var sharedSecretId int64
	for rows3.Next() {
		err = rows3.Scan(&sharedSecretId)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to get sharedSecretId %s: %w", args.User.GetLogin(), err)
		}
	}

	//3 get identity id
	sqlQuery4 := "select identity  from shared_secret_identity where shared_secret= ?;" //shared_secret
	rows2, err = tx.Queryx(sqlQuery4, sharedSecretId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get identity shared_secret user %s: %w", args.User.GetLogin(), err)
	}
	defer rows2.Close()
	var identityS int64
	var identitySrr []int64
	for rows2.Next() {
		err = rows2.Scan(&identityS)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to get identity sharedSecretId %s: %w", args.User.GetLogin(), err)
		}
		identitySrr = append(identitySrr, identityS)
	}

	//4 delete shared_secret_identity
	sqlQuery5, argss, err := sqlx.In("delete from shared_secret_identity  where identity in (?) AND shared_secret=?", identitySrr, sharedSecretId)
	sqlQuery5 = h.db.Rebind(sqlQuery5)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user shared_secret_identity %s: %w", args.User.GetLogin(), err)
	}
	_, err = tx.Exec(sqlQuery5, argss...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user shared_secret_identity %s: %w", args.User.GetLogin(), err)
	}

	//5 delete shared_secret
	sqlQuery6 := "delete from shared_secrets where id = ?;"
	_, err = tx.Exec(sqlQuery6, identitiesId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to delete user shared_secret by identity %s: %w", args.User.GetLogin(), err)
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
