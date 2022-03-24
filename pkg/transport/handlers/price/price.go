package price

import (
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/pay"
)

type Price struct {
}

func NewPrice() *Price {
	return &Price{}
}

func (p *Price) Endpoint() interface{} {
	return `/price`
}

func (p *Price) Handler() func() interface{} {
	return func() interface{} { return "100500" }
}

func (p *Price) Menu() (res []transport.MenuI) {
	res = append(res, pay.NewPay())
	return res
}
