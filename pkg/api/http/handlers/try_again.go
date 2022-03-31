package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/pkg/controller"

	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
)

type TryAgainHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewTryAgainHandler(ctrl controller.Controller, logger *zerolog.Logger) *TryAgainHandler {
	return &TryAgainHandler{ctrl, logger}
}

func (h TryAgainHandler) Handle(params operations.PostTryagainParams) middleware.Responder {

	h.logger.Debug().Msgf("TryAgainHandler %s", params)
	return operations.NewPostTryagainOK()
}
