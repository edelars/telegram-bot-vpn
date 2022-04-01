package pay_get_invoice

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/controller"
	"context"
	"errors"
	"github.com/rs/zerolog"
	"strconv"
)

type PayGetInvoice struct {
	valueToPay int
	user       dto.User
	Out        struct {
		Url string
	}
}

func NewPayInvoice(user dto.User, valueToPay int) *PayGetInvoice {
	return &PayGetInvoice{
		user:       user,
		valueToPay: valueToPay}
}

type PayGetInvoiceHandler struct {
	ctrl     controller.Controller
	logger   *zerolog.Logger
	provider ProviderI
}

func NewPayGetInvoiceHandler(ctrl controller.Controller, logger *zerolog.Logger, provider ProviderI) *PayGetInvoiceHandler {
	return &PayGetInvoiceHandler{ctrl: ctrl, logger: logger, provider: provider}
}

func (h *PayGetInvoiceHandler) Exec(ctx context.Context, args *PayGetInvoice) (err error) {

	if h.provider == nil {
		return errors.New("provider is nil")
	}

	pd := dto.PayData{
		Value: args.valueToPay,
		Order: strconv.FormatInt(args.user.Id, 10),
	}
	if err, args.Out.Url = h.provider.GenerateUrlToPay(pd); err != nil {
		return err
	}

	h.logger.Debug().Msgf("Got Url for user: %s to pay: %s", args.user.Login, args.Out.Url)

	return nil
}

func (h *PayGetInvoiceHandler) Context() interface{} {
	return (*PayGetInvoice)(nil)
}

type ProviderI interface {
	GenerateUrlToPay(dto.PayData) (error, string)
}
