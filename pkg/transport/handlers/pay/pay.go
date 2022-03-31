package pay

import (
	"backend-vpn/pkg/billing/pay_prepare"
	"backend-vpn/pkg/config"
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
	text     string
	price    int
	dayCount int
}

func newPay(ctrl controller.Controller, logger *zerolog.Logger, text string, endpoint string, price int, dayCount int) *Pay {
	return &Pay{ctrl: ctrl, logger: logger, endpoint: endpoint, text: text, price: price, dayCount: dayCount}
}
func (p *Pay) Data() (text, unique string) {
	return p.text, p.endpoint
}
func (p *Pay) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		u := storage.NewUserQuery(data.Username, data.Id, "")

		if err := p.ctrl.Exec(context.Background(), u); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint("Ошибка, попробуйте позже")
		}
		p.logger.Debug().Msgf("user %s press %s button", u.Out.User, p.endpoint)

		pp := pay_prepare.NewPayPrepare(u.Out.User, p.price, p.dayCount)
		if err := p.ctrl.Exec(context.Background(), pp); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint("Ошибка, попробуйте позже")
		}
		return pp.Out.Message
	}
}
func NewPayType00(ctrl controller.Controller, logger *zerolog.Logger) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("1 день - Бесплатно"), "pay00", 0, 1)
}

func NewPayType01(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("30 дней %dр", env.Price01), "pay01", env.Price01, 30)
}

func NewPayType06(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("180 дней %dр", env.Price06), "pay06", env.Price06, 180)
}

func NewPayType12(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("365 дней %dр", env.Price12), "pay12", env.Price12, 365)
}
