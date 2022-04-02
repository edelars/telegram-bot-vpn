package get_expired_users

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

type getExpiredUsersHandler struct {
	db *sqlx.DB
}

func NewGetExpiredUsersHandler(db *sqlx.DB) *getExpiredUsersHandler {
	return &getExpiredUsersHandler{db: db}
}

func (h *getExpiredUsersHandler) Exec(ctx context.Context, args *storage.GetExpiredUsers) (err error) {

	timeDisable := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)

	sqlQuery := `select tg_id, tg_login,created_at,referal_id,expired_at,invite_referal_id,password from users where expired_at > ? and expired_at <= ?`
	var rows *sqlx.Rows
	rows, err = h.db.QueryxContext(ctx, sqlQuery, timeDisable, time.Now())
	for rows.Next() {
		var p get_create_update_user.GetCreateUpdateUser
		var u dto.User
		if err = rows.StructScan(&p); err != nil {
			return err
		}

		if p.Login.Valid {
			u.Login = p.Login.String
		}
		if p.ExpiredAt.Valid {
			u.ExpiredAt = p.ExpiredAt.Time
		}
		if p.CreatedAt.Valid {
			u.CreatedAt = p.CreatedAt.Time
		}
		if p.ReferalId.Valid {
			u.ReferalId = p.ReferalId.String
		}
		if p.Id.Valid {
			u.Id = p.Id.Int64
		}
		if p.Password.Valid {
			u.Password = p.Password.String
		}

		args.Out = append(args.Out, &u)
	}
	defer rows.Close()

	return nil
}

func (h *getExpiredUsersHandler) Context() interface{} {
	return (*storage.GetExpiredUsers)(nil)
}
