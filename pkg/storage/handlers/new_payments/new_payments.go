package new_payments

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type newPaymentsHandler struct {
	db *sqlx.DB
}

var (
	ErrNotExist   = errors.New("no such user")
	ErrUserIdZero = errors.New("userid is 0")
)

func NewNewPaymentsHandler(db *sqlx.DB) *newPaymentsHandler {
	return &newPaymentsHandler{db: db}
}

func (h newPaymentsHandler) Exec(ctx context.Context, args *storage.NewPayments) (err error) {

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	var rows *sqlx.Rows

	// 1 select id
	sqlQuery := `select id from users where tg_id = ? limit 1;`
	rows, err = tx.Queryx(sqlQuery, args.UserId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get user id %d: %w", args.UserId, err)
	}

	var userId int
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			rows.Close()
			_ = tx.Rollback()
			return fmt.Errorf("failed to scan user id %d: %w", args.UserId, err)
		}
		rows.Close()
	} else {
		rows.Close()
		_ = tx.Rollback()
		return ErrNotExist
	}

	if userId == 0 {
		return ErrUserIdZero
	}

	//2 insert
	sqlQuery2 := `insert into payments (user_id, created_at,value) VALUES (?,?,?);`
	_, err = tx.Exec(sqlQuery2, userId, time.Now(), args.Value)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert payments %d: %w", args.UserId, err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (h *newPaymentsHandler) Context() interface{} {
	return (*storage.NewPayments)(nil)
}
