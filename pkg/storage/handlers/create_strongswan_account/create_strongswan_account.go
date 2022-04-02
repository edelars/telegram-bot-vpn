package create_strongswan_account

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/storage"
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math/big"
	"net"
	"strings"
)

type createStrongswanAccountHandler struct {
	db  *sqlx.DB
	env config.Environment
}

func NewCreateStrongswanAccountHandler(db *sqlx.DB, env config.Environment) *createStrongswanAccountHandler {
	return &createStrongswanAccountHandler{db: db, env: env}
}

func (h createStrongswanAccountHandler) Exec(ctx context.Context, args *storage.CreateStrongswanAccount) (err error) {

	if args.User == nil {
		return errors.New("No *User")
	}

	tx, err := h.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}
	var last_id_in_identities, last_id_in_shared_secrets int64
	var res sql.Result

	// 1.1
	sqlQuery := `insert into identities (type, data) VALUES (?,?);`
	res, err = tx.Exec(sqlQuery, 2, args.User.GetLoginByte())
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}
	last_id_in_identities, err = res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert user %s: %w", args.User.GetLogin(), err)
	}

	//2
	sqlQuery2 := "insert into shared_secrets (type,data) VALUES (?,?);"
	res, err = tx.Exec(sqlQuery2, 2, args.User.GetPasswordByte())
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
	var insertedIdSrv []int64
	err, insertedIdSrv = h.GetOrCreateServersIdent(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

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

func (h *createStrongswanAccountHandler) GetOrCreateServersIdent(tx *sqlx.Tx) (err error, insertedIdSrv []int64) {

	srvIp := strings.Fields(h.env.OurServersIP)
	sqlQueryIns := `insert into identities (type, data) VALUES (?,?);`
	sqlQueryGet := `select id from identities where data = ?;`
	var last_id int64
	var res sql.Result
	var rows *sqlx.Rows

	for _, s := range srvIp {
		ss := Pack32BinaryIP4(s)

		rows, err = tx.Queryx(sqlQueryGet, ss)
		if err != nil {
			return fmt.Errorf("failed to get srv user %s: %w", s, err), insertedIdSrv
		}
		defer rows.Close()
		if rows.Next() {
			err = rows.Scan(&last_id)
			if err != nil {
				return fmt.Errorf("failed to scan srv user %s: %w", s, err), insertedIdSrv
			}
			insertedIdSrv = append(insertedIdSrv, last_id)
			continue
		}

		res, err = tx.Exec(sqlQueryIns, 1, ss)
		if err != nil {
			return fmt.Errorf("failed to insert srv user %s: %w", s, err), insertedIdSrv
		}
		last_id, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to insert srv user %s: %w", s, err), insertedIdSrv
		}
		insertedIdSrv = append(insertedIdSrv, last_id)
	}

	return err, insertedIdSrv
}

func Pack32BinaryIP4(ip4Address string) []byte {
	ipv4Decimal := IP4toInt(net.ParseIP(ip4Address))

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, uint32(ipv4Decimal))

	if err != nil {
		fmt.Println("Unable to write to buffer:", err) //TODO fix shit
	}
	return buf.Bytes()
}

func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

func (h *createStrongswanAccountHandler) Context() interface{} {
	return (*storage.CreateStrongswanAccount)(nil)
}
