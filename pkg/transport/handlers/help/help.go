package help

import (
	"backend-vpn/pkg/transport"
	"fmt"
)

type Help struct {
}

func NewHelp() *Help {
	return &Help{}
}

func (p *Help) Endpoint() interface{} {
	return `/help`
}

func (p *Help) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		return fmt.Sprint("Здесь будет инфа о помощи...пока можно спросить у @edelars")
	}
	//TODO get data
}

func (p *Help) Menu() []transport.MenuI {
	return nil
}
