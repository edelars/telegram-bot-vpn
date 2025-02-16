package account_info

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
)

type AccountInfo struct {
	user    *dto.User
	balance int
	userId  int64
	Out     struct {
		Message string
	}
}

func NewAccountInfoWithData(user *dto.User, balance int) *AccountInfo {
	return &AccountInfo{user: user, balance: balance, Out: struct{ Message string }{}}
}

func NewAccountInfo(userId int64) *AccountInfo {
	return &AccountInfo{userId: userId, Out: struct{ Message string }{}}
}

type accountInfoHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
	env    config.Environment
}

func NewAccountInfoHandler(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *accountInfoHandler {
	return &accountInfoHandler{logger: logger, ctrl: ctrl, env: env}
}

func (h accountInfoHandler) Exec(ctx context.Context, args *AccountInfo) (err error) {

	if args.user == nil {
		if args.userId == 0 {
			return errors.New("userid is zero")
		}

		ngu := storage.NewGetUser(args.userId)
		if err := h.ctrl.Exec(ctx, ngu); err != nil {
			h.logger.Debug().Err(err).Msg("accountInfoHandler:NewGetUser fail")
			return err
		}
		args.user = ngu.Out.User

		bal := storage.GetUserBalanceQuery{UserId: args.userId}
		if err := h.ctrl.Exec(ctx, &bal); err != nil {
			h.logger.Debug().Err(err).Msg("accountInfoHandler:GetUserBalanceQuery fail")
		}
		args.balance = int(bal.Out.TotalBalance)
	}
	ipadr := h.env.OurServersIP
	args.Out.Message = fmt.Sprintf("Информация об аккаунте VPN:\n\nАдрес IP: %s\nlogin: %s\npassword: %s\npresharedkey(PSK): %s\n\nАктивен до: %s\nБаланс: %d руб.\n\nId: %d\nRef id: %s",
		ipadr, args.user.Login, args.user.Password, args.user.Psk, args.user.ExpiredAt.Format("02 Jan 06 15:04 MST"), args.balance, args.user.Id, args.user.ReferalId)

	return nil
}

func (h *accountInfoHandler) Context() interface{} {
	return (*AccountInfo)(nil)
}
