package handlers

import (
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/add"
	"backend-vpn/pkg/transport/handlers/del"
	"backend-vpn/pkg/transport/handlers/help"
	"backend-vpn/pkg/transport/handlers/info"
	"backend-vpn/pkg/transport/handlers/price"
	"backend-vpn/pkg/transport/handlers/ref"
	"backend-vpn/pkg/transport/handlers/start"
	"github.com/rs/zerolog"
)

func GetHandlers(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) (res []transport.HandlerI) {
	res = append(res, price.NewPrice(ctrl, logger, env))
	res = append(res, start.NewStart())
	res = append(res, add.NewAdd(ctrl, logger, env))
	res = append(res, del.NewDel(ctrl, logger, env))
	res = append(res, help.NewHelp())
	res = append(res, info.NewInfo(ctrl, logger))
	res = append(res, ref.NewRef())
	return res
}
