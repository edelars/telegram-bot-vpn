package info

import (
	"backend-vpn/pkg/billing/account_info"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"context"
	"github.com/rs/zerolog"
)

type Info struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewInfo(ctrl controller.Controller, logger *zerolog.Logger) *Info {
	return &Info{ctrl: ctrl, logger: logger}
}

func (p *Info) Endpoint() interface{} {
	return `/info`
}

func (p *Info) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {
		ai := account_info.NewAccountInfo(data.Id)
		if err := p.ctrl.Exec(context.Background(), ai); err != nil {
			p.logger.Debug().Err(err).Msg("Info:AccountInfo fail")
			return "Ошибка, попробуйте позже"
		}

		return ai.Out.Message
	}
}

func (p *Info) Menu() []transport.MenuI {
	return nil
}
