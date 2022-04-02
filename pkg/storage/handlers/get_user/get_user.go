package get_user

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"context"
	"github.com/jmoiron/sqlx"
)

type getUserHandler struct {
	db  *sqlx.DB
	env config.Environment
}

func NewGetUserHandler(db *sqlx.DB, env config.Environment) *getUserHandler {
	return &getUserHandler{db: db, env: env}
}

func (h getUserHandler) Exec(ctx context.Context, args *storage.GetUser) (err error) {

	var p get_create_update_user.GetCreateUpdateUser

	sqlQuery := `select tg_id, tg_login,created_at,referal_id,expired_at,invite_referal_id,password from users where tg_id = ? limit 1`
	var rows *sqlx.Rows
	rows, err = h.db.QueryxContext(ctx, sqlQuery, args.UserId)
	for rows.Next() {
		if err = rows.StructScan(&p); err != nil {
			return err
		}
	}

	defer rows.Close()

	if p.Login.Valid {
		args.Out.User.Login = p.Login.String
	}
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
	args.Out.User.Psk = h.env.Psk

	return nil
}

func (h *getUserHandler) Context() interface{} {
	return (*storage.GetUser)(nil)
}
