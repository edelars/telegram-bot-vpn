package get_create_update_user

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

var (
	errNotExist = errors.New("Not exist")
)

type GetCreateUpdateUser struct {
	Id              sql.NullInt64  `db:"tg_id"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	ReferalId       sql.NullString `db:"referal_id"`
	ExpiredAt       sql.NullTime   `db:"expired_at"`
	InviteReferalId sql.NullString `db:"invite_referal_id"`
	Login           sql.NullString `db:"tg_login"`
	Password        sql.NullString `db:"password"`
	UsedTestPeriod  sql.NullBool   `db:"used_test_period"`
}

type GetCreateUpdateUserHandler struct {
	db  *sqlx.DB
	env config.Environment
}

func NewGetCreateUpdateUserHandler(db *sqlx.DB, env config.Environment) *GetCreateUpdateUserHandler {
	return &GetCreateUpdateUserHandler{db: db, env: env}
}

func (h GetCreateUpdateUserHandler) Exec(ctx context.Context, args *storage.UserQuery) (err error) {

	userLogin := args.GetLogin()
	if userLogin == "" {
		userLogin = strconv.FormatInt(args.GetId(), 10)
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	var p GetCreateUpdateUser

	err = get(tx, &p, args.GetId())
	switch err {
	case errNotExist:
		if err = create(tx, &p, userLogin, args.GetId(), args.GetNewReferalId(), args.GetInviteReferalId()); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = get(tx, &p, args.GetId()); err != nil {
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
	if p.Id.Valid {
		args.Out.User.Id = p.Id.Int64
	}
	if p.Password.Valid {
		args.Out.User.Password = p.Password.String
	}
	if p.UsedTestPeriod.Valid {
		args.Out.User.UsedTestPeriod = p.UsedTestPeriod.Bool
	}

	args.Out.User.Psk = h.env.Psk

	return nil
}

func get(tx *sqlx.Tx, p *GetCreateUpdateUser, tgId int64) (err error) {
	var rows *sqlx.Rows

	sqlQuery := `select tg_id, tg_login,created_at,referal_id,expired_at,invite_referal_id,password,used_test_period from users where tg_id = ?`

	rows, err = tx.QueryxContext(context.Background(), sqlQuery, tgId)
	if err != nil {
		return fmt.Errorf("failed to query user %s: %w", tgId, err)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.StructScan(&p)
	} else {
		return errNotExist
	}
	return nil
}

func create(tx *sqlx.Tx, p *GetCreateUpdateUser, userLogin string, tgId int64, referalId, invite_referal_id string) (err error) {
	var rows *sqlx.Rows

	sqlQuery := `insert into users (created_at, tg_login, referal_id, expired_at, invite_referal_id,tg_id) VALUES (?,?,?,?,?,?);`
	expired := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	rows, err = tx.Queryx(sqlQuery, time.Now(), userLogin, referalId, expired, invite_referal_id, tgId)
	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", userLogin, err)
	}
	defer rows.Close()

	return nil
}
func (h *GetCreateUpdateUserHandler) Context() interface{} {
	return (*storage.UserQuery)(nil)
}
