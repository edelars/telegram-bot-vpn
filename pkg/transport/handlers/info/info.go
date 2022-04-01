package info

import (
	"backend-vpn/pkg/transport"
	"fmt"
)

type Info struct {
}

func NewInfo() *Info {
	return &Info{}
}

func (p *Info) Endpoint() interface{} {
	return `/info`
}

func (p *Info) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		return fmt.Sprintf("Информация об аккаунте VPN:\n\n/login: %s\n/password: %s\npresharedkey(PSK): %s\n\n Аккаунт активен до(UTC +0):%s\nБаланс: %s руб.", "", "", "", "", "")
	}
	//TODO get data
}

func (p *Info) Menu() []transport.MenuI {
	return nil
}
