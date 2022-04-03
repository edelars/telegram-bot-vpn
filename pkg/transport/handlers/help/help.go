package help

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/iphone"
	"fmt"
	"github.com/rs/zerolog"
)

type Help struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewHelp(ctrl controller.Controller, logger *zerolog.Logger) *Help {
	return &Help{ctrl: ctrl, logger: logger}
}

func (p *Help) Endpoint() interface{} {
	return `/help`
}

func (p *Help) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		return fmt.Sprint("Мы используем протоколы IKEv2 (PSK) и IPSec Xauth PSK.\nЗдесь будет инфа о помощи...пока можно спросить у @edelars")
	}
	//TODO get data
}

func (p *Help) Menu() (res []transport.MenuI) {
	res = append(res, iphone.NewIphone(p.ctrl, p.logger, "Файл настроек Iphone", "iphone"))
	return res
}
