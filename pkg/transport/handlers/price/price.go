package price

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/pay"
	"github.com/rs/zerolog"
)

type Price struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewPrice(ctrl controller.Controller, logger *zerolog.Logger) *Price {
	return &Price{ctrl: ctrl, logger: logger}
}

func (p *Price) Endpoint() interface{} {
	return `/price`
}

func (p *Price) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} { return "100500" }
}

func (p *Price) Menu() (res []transport.MenuI) {
	res = append(res, pay.NewPay(p.ctrl, p.logger))
	return res
}
