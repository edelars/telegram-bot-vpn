package price

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/pay"
	"fmt"
	"github.com/rs/zerolog"
)

type Price struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
	env    config.Environment
}

func NewPrice(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Price {
	return &Price{ctrl: ctrl, logger: logger, env: env}
}

func (p *Price) Endpoint() interface{} {
	return `/price`
}

func (p *Price) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		return fmt.Sprintf("Тестовый - бесплатно в течении 24ч\nМесяц (30 дней)- %dр\nПолгода (180 дней)- %dр\nГод (365 дней)- %d\nНажмите кнопку внизу для покупки",
			p.env.Price01, p.env.Price06, p.env.Price12)
	}
}

func (p *Price) Menu() (res []transport.MenuI) {
	res = append(res, pay.NewPayType00(p.ctrl, p.logger))
	res = append(res, pay.NewPayType01(p.ctrl, p.logger, p.env))
	res = append(res, pay.NewPayType06(p.ctrl, p.logger, p.env))
	res = append(res, pay.NewPayType12(p.ctrl, p.logger, p.env))
	return res
}
