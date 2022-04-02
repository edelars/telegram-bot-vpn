package pay_prepare

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/billing/activate_account"
	"backend-vpn/pkg/billing/pay_get_invoice"
	"backend-vpn/pkg/controller"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

const (
	messageSuccess = "\nДанные для подключения:\nLogin: %s\n Password: %s\nPreSharedKey (PSK): %s\n\nДля получения помощи по настройке VPN наберите команду /help или выберете её в меню"
)

type PayPrepare struct {
	valueToPay int
	dayCount   int
	user       *dto.User
	Out        struct {
		Message string
	}
}

func NewPayPrepare(user *dto.User, valueToPay int, dayCount int) *PayPrepare {
	return &PayPrepare{
		user:       user,
		valueToPay: valueToPay,
		dayCount:   dayCount}
}

type PayPrepareHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewPayPrepareHandler(ctrl controller.Controller, logger *zerolog.Logger) *PayPrepareHandler {
	return &PayPrepareHandler{ctrl: ctrl, logger: logger}
}

func (h *PayPrepareHandler) Exec(ctx context.Context, args *PayPrepare) (err error) {

	var a *activate_account.ActivateAccount
	if args.valueToPay > 0 {

		pgi := pay_get_invoice.NewPayInvoice(*args.user, args.valueToPay)
		if err := h.ctrl.Exec(ctx, pgi); err != nil {
			h.logger.Debug().Err(err).Msg("fail")
			return err
		}

		args.Out.Message = fmt.Sprintf("После оплаты, ваш аккаунт будет активирован\n\n Ссылка для оплаты:\n%s", pgi.Out.Url)

	} else {

		if args.user.UsedTestPeriod {
			args.Out.Message = "Вы уже использовали тестовый период"
			return nil
		}

		a, err = activate_account.NewActivateAccount(args.user, args.dayCount, true)
		if err != nil {
			return err
		}
		if err := h.ctrl.Exec(ctx, a); err != nil {
			h.logger.Debug().Err(err).Msg("fail")
			return err
		}
		usr := a.GetUser()
		h.logger.Debug().Msgf("New trial account for user %s ", usr.Login)
		args.Out.Message = fmt.Sprintf("Вам активирован пробный период на %d день.\n"+messageSuccess, usr.Login, usr.Password, usr.Psk, args.dayCount)
	}

	return err
}

func (h *PayPrepareHandler) Context() interface{} {
	return (*PayPrepare)(nil)
}
