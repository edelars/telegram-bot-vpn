package handlers

import (
	"backend-vpn/pkg/transport"
	"backend-vpn/pkg/transport/handlers/price"
	"backend-vpn/pkg/transport/handlers/start"
)

func GetHandlers() (res []transport.HandlerI) {
	res = append(res, price.NewPrice())
	res = append(res, start.NewStart())
	return res
}
