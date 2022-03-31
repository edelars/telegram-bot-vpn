package save_user

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SaveUserHandler struct {
	db *sqlx.DB
}

func NewSaveUserHandler(db *sqlx.DB) *SaveUserHandler {
	return &SaveUserHandler{db: db}
}

func (h SaveUserHandler) Exec(ctx context.Context, args *storage.SaveUserQuery) (err error) {

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	sqlQuery := `insert into users (tg_id, used_test_period, tg_login, referal_id, expired_at, invite_referal_id,password) VALUES (?,?,?,?,?,null,?)
 					ON DUPLICATE KEY UPDATE referal_id=?, expired_at = ?, password=?,used_test_period = ?;`

	user := args.GetUser()
	var rows *sqlx.Rows
	rows, err = tx.Queryx(sqlQuery,
		user.Id, user.UsedTestPeriod, user.Login, user.ReferalId, user.ExpiredAt, user.Password,
		user.ReferalId, user.ExpiredAt, user.Password, user.UsedTestPeriod)
	if err != nil {
		return fmt.Errorf("failed to update user info %s: %w", user.Login, err)
	}
	defer rows.Close()

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (h *SaveUserHandler) Context() interface{} {
	return (*storage.SaveUserQuery)(nil)
}
