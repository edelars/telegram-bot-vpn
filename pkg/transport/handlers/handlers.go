package handlers

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/price"
	"backend-vpn/pkg/transport/handlers/start"
	"github.com/rs/zerolog"
)

func GetHandlers(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) (res []transport.HandlerI) {
	res = append(res, price.NewPrice(ctrl, logger, env))
	res = append(res, start.NewStart())
	return res
}
