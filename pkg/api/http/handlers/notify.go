package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/pkg/controller"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
)

type NotifyHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewNotifyHandler(ctrl controller.Controller, logger *zerolog.Logger) *NotifyHandler {
	return &NotifyHandler{ctrl, logger}
}

func (h NotifyHandler) Handle(params operations.PostNotifyParams) middleware.Responder {

	h.logger.Debug().Msgf("NotifyHandler %s", params)
	return operations.NewPostNotifyOK()
}
