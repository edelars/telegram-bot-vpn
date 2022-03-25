package get_create_update_user

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	errNotExist = errors.New("Not exist")
)

type GetCreateUpdateUser struct {
	Id              sql.NullInt64  `db:"id"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	ReferalId       sql.NullString `db:"referal_id"`
	ExpiredAt       sql.NullTime   `db:"expired_at"`
	InviteReferalId sql.NullString `db:"invite_referal_id"`
}

type GetCreateUpdateUserHandler struct {
	db *sqlx.DB
}

func NewGetCreateUpdateUserHandler(db *sqlx.DB) *GetCreateUpdateUserHandler {
	return &GetCreateUpdateUserHandler{db: db}
}

func (h GetCreateUpdateUserHandler) Exec(ctx context.Context, args *storage.UserQuery) (err error) {

	userLogin := args.GetLogin()
	if userLogin == "" {
		return fmt.Errorf("user login in query not set")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	var p GetCreateUpdateUser

	err = get(tx, &p, userLogin)
	switch err {
	case errNotExist:
		if err = create(tx, &p, userLogin, args.GetReferalId()); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = get(tx, &p, userLogin); err != nil {
			_ = tx.Rollback()
			return err
		}
	case nil:

	default:
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	args.Out.User.Login = userLogin
	if p.ExpiredAt.Valid {
		args.Out.User.ExpiredAt = p.ExpiredAt.Time
	}
	if p.CreatedAt.Valid {
		args.Out.User.CreatedAt = p.CreatedAt.Time
	}
	if p.ReferalId.Valid {
		args.Out.User.ReferalId = p.ReferalId.String
	}
	return nil
}

func get(tx *sqlx.Tx, p *GetCreateUpdateUser, userLogin string) (err error) {
	var rows *sqlx.Rows

	sqlQuery := `select id,created_at,referal_id,expired_at,invite_referal_id from users where tg_login = ?`

	rows, err = tx.QueryxContext(context.Background(), sqlQuery, userLogin)
	if err != nil {
		return fmt.Errorf("failed to query user %s: %w", userLogin, err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&p)
	} else {
		return errNotExist
	}
	return nil
}

func create(tx *sqlx.Tx, p *GetCreateUpdateUser, userLogin, referalId string) (err error) {
	var rows *sqlx.Rows

	sqlQuery := `insert into users (created_at, tg_login, referal_id, expired_at, invite_referal_id) VALUES (default,?,?,?,null);`

	rows, err = tx.Queryx(sqlQuery, userLogin, referalId, time.Now())
	if err != nil {
		return fmt.Errorf("failed to query user %s: %w", userLogin, err)
	}
	defer rows.Close()

	return nil
}
func (h *GetCreateUpdateUserHandler) Context() interface{} {
	return (*storage.UserQuery)(nil)
}
