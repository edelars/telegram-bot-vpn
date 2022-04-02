package writeoff_balance

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type writeoffBalanceHandler struct {
	db *sqlx.DB
}

func NewWriteoffBalanceHandler(db *sqlx.DB) *writeoffBalanceHandler {
	return &writeoffBalanceHandler{db: db}
}

func (h writeoffBalanceHandler) Exec(ctx context.Context, args *storage.WriteOffBalance) (err error) {

	if args.UserId == 0 {
		return errors.New("userId is zero")
	}
	if args.Value <= 0 {
		return errors.New("value is too low")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	sqlQuery := `insert into balance (user_id,credit) VALUES ((select id from users where tg_id =?),?);`

	_, err = tx.ExecContext(ctx, sqlQuery, args.UserId, args.Value)

	if err != nil {
		return fmt.Errorf("failed to writeoff balance user %d: %w", args.UserId, err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (h *writeoffBalanceHandler) Context() interface{} {
	return (*storage.WriteOffBalance)(nil)
}
