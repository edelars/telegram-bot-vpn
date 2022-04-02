package unchecked_payments_to_balance

import (
	"backend-vpn/pkg/storage"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type uncheckedPaymentsToBalanceHandler struct {
	db *sqlx.DB
}

func NewUncheckedPaymentsToBalanceHandler(db *sqlx.DB) *uncheckedPaymentsToBalanceHandler {
	return &uncheckedPaymentsToBalanceHandler{db: db}
}

func (h uncheckedPaymentsToBalanceHandler) Exec(ctx context.Context, args *storage.UncheckedPaymentsQuery) (err error) {

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	var rows *sqlx.Rows

	// 1 select pay
	sqlQuery := `select payments.id,payments.value,user1.invite_referal_payed,user2.tg_id as ref_user_id
		from payments
    	 LEFT JOIN
     	users as user1 ON payments.user_id = user1.id
     	LEFT JOIN
   		users as user2 ON user1.invite_referal_id = user2.referal_id
		where payments.user_id = (select id from users where tg_id =?)  and payments.checked = false`
	rows, err = tx.Queryx(sqlQuery, args.UserId)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to get payments id %d: %w", args.UserId, err)
	}

	var paymentsIdArr []int64
	var debtBal, invite_ref int64
	var ref_payed bool

	for rows.Next() {
		var r rowPay
		err = rows.StructScan(&r)
		if err != nil {

			_ = tx.Rollback()
			return fmt.Errorf("failed to scan user id %d: %w", args.UserId, err)
		}
		if r.Id.Valid {
			paymentsIdArr = append(paymentsIdArr, r.Id.Int64)
		}
		if r.Value.Valid {
			debtBal = debtBal + r.Value.Int64
		}
		if r.Ref_payed.Valid {
			ref_payed = r.Ref_payed.Bool
		}
		if r.Ref_user_id.Valid {
			invite_ref = r.Ref_user_id.Int64
		}

	}
	defer rows.Close()

	if len(paymentsIdArr) == 0 {
		return fmt.Errorf("there no payments for user: %d", args.UserId)
	}

	//2.1 insert
	sqlQuery2 := `insert into balance (user_id,debt) VALUES ((select id from users where tg_id =?),?);`
	_, err = tx.Exec(sqlQuery2, args.UserId, debtBal)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert balance %d: %w", args.UserId, err)
	}

	//2.1 insert if
	if !ref_payed && invite_ref > 0 {
		_, err = tx.Exec(sqlQuery2, invite_ref, debtBal)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert referal balance %d: %w", invite_ref, err)
		}
		//2.2 update user ref_payed
		sqlQuery22 := `update users set invite_referal_payed = true where tg_id = ?;`
		_, err = tx.Exec(sqlQuery22, args.UserId)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to update users invite_referal_payed %d: %w", args.UserId, err)
		}
	}

	//3 update payments
	sqlQuery3, argss, err := sqlx.In("update payments set checked = true where id in (?) and user_id = (select id from users where tg_id =?)", paymentsIdArr, args.UserId)
	sqlQuery3 = h.db.Rebind(sqlQuery3)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update payments %s: %w", args.UserId, err)
	}
	_, err = tx.Exec(sqlQuery3, argss...)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update payments %s: %w", args.UserId, err)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	args.Out.IncBalance = debtBal

	return nil
}

func (h *uncheckedPaymentsToBalanceHandler) Context() interface{} {
	return (*storage.UncheckedPaymentsQuery)(nil)
}

type rowPay struct {
	Id          sql.NullInt64 `db:"id"`
	Value       sql.NullInt64 `db:"value"`
	Ref_payed   sql.NullBool  `db:"invite_referal_payed"`
	Ref_user_id sql.NullInt64 `db:"ref_user_id"`
}
