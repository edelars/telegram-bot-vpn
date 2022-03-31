package create_strongswan_account

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type CreateStrongswanAccountHandler struct {
	db  *sqlx.DB
	env config.Environment
}

func NewCreateStrongswanAccountHandler(db *sqlx.DB, env config.Environment) *CreateStrongswanAccountHandler {
	return &CreateStrongswanAccountHandler{db: db, env: env}
}

func (h CreateStrongswanAccountHandler) Exec(ctx context.Context, args *storage.CreateStrongswanAccount) (err error) {

	if args.User == nil {
		return errors.New("No *User")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	var last_id_in_identities, last_id_in_shared_secrets, last_id int64
	var res sql.Result

	// 1.1
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

	//1.2 srv ip
	var insertedIdSrv []int64
	srvIp := strings.Fields(h.env.OurServersIP)
	for _, s := range srvIp {
		res, err = tx.Exec(sqlQuery, 1, s)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert srv user %s: %w", s, err)
		}
		last_id, err = res.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert srv user %s: %w", s, err)
		}
		insertedIdSrv = append(insertedIdSrv, last_id)
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

	//3.1
	sqlQuery3 := "insert into shared_secret_identity (shared_secret, identity) VALUES (?,?);"
	res, err = tx.Exec(sqlQuery3, last_id_in_shared_secrets, last_id_in_identities)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}

	//3.2
	for _, i2 := range insertedIdSrv {
		res, err = tx.Exec(sqlQuery3, last_id_in_shared_secrets, i2)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert srv shared_secret user %s: %w", i2, err)
		}
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
