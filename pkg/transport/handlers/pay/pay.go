package pay

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/transport"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type Pay struct {
	ctrl     controller.Controller
	logger   *zerolog.Logger
	endpoint string
}

func NewPay(ctrl controller.Controller, logger *zerolog.Logger) *Pay {
	return &Pay{ctrl: ctrl, logger: logger, endpoint: "pay"}
}
func (p *Pay) Data() (text, unique string) {
	return "Оплатить VPN", p.endpoint
}
func (p *Pay) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		u := storage.NewUserQuery(data.Username, "")

		if err := p.ctrl.Exec(context.Background(), u); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint("Ошибка, попробуйте позже")
		}
		return u.Out.User.CreatedAt.String()
	}
}
