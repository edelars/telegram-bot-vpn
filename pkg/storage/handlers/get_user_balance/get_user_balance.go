package get_user_balance

import (
	"backend-vpn/pkg/storage"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type getUserBalanceHandler struct {
	db *sqlx.DB
}

func NewGetUserBalanceHandler(db *sqlx.DB) *getUserBalanceHandler {
	return &getUserBalanceHandler{db: db}
}

func (h getUserBalanceHandler) Exec(ctx context.Context, args *storage.GetUserBalanceQuery) (err error) {

	var rows *sqlx.Rows

	// 1 select balance
	sqlQuery := `select (sum(debt)-sum(credit)) as balance from balance where user_id = (select id from users where tg_id =?)`
	rows, err = h.db.Queryx(sqlQuery, args.UserId)
	if err != nil {
		return fmt.Errorf("failed to get balance id %d: %w", args.UserId, err)
	}

	var balance int64

	for rows.Next() {
		err = rows.Scan(&balance)
		if err != nil {
			return fmt.Errorf("failed to scan balance id %d: %w", args.UserId, err)
		}
	}
	defer rows.Close()

	args.Out.TotalBalance = balance

	return nil
}

func (h *getUserBalanceHandler) Context() interface{} {
	return (*storage.GetUserBalanceQuery)(nil)
}
