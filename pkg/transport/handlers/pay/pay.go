package pay

import (
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
}

func newPay(ctrl controller.Controller, logger *zerolog.Logger, text string, endpoint string, price int) *Pay {
	return &Pay{ctrl: ctrl, logger: logger, endpoint: endpoint, text: text, price: price}
}
func (p *Pay) Data() (text, unique string) {
	return p.text, p.endpoint
}
func (p *Pay) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		u := storage.NewUserQuery(data.Id, "")

		if err := p.ctrl.Exec(context.Background(), u); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint("Ошибка, попробуйте позже")
		}
		p.logger.Debug().Msgf("user %s press %s button", u.Out.User, p.endpoint)
		return u.Out.User.CreatedAt.String()
	}
}
func NewPayType00(ctrl controller.Controller, logger *zerolog.Logger) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("1 день - Бесплатно"), "pay00", 0)
}

func NewPayType01(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("30 дней %dр", env.Price01), "pay01", env.Price01)
}

func NewPayType06(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("180 дней %dр", env.Price06), "pay06", env.Price06)
}

func NewPayType12(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Pay {
	return newPay(ctrl, logger, fmt.Sprintf("365 дней %dр", env.Price12), "pay12", env.Price12)
}
