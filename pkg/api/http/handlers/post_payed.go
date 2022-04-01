package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/pkg/billing/pay_incoming_transaction"
	"backend-vpn/pkg/controller"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
	"strconv"
)

type PostPayedHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewPostPayedHandler(ctrl controller.Controller, logger *zerolog.Logger) *PostPayedHandler {
	return &PostPayedHandler{ctrl, logger}
}

func (h PostPayedHandler) Handle(params operations.PostPayedParams) middleware.Responder {

	h.logger.Debug().Msgf("Payed %s", params)

	ctx := params.HTTPRequest.Context()

	if params.MERCHANTORDERID == nil || params.AMOUNT == nil {
		h.logger.Err(errors.New("nil data in query")).Msg("")
		return operations.NewPostPayedInternalServerError()
	}

	userId, err := strconv.ParseInt(*params.MERCHANTORDERID, 10, 64)
	if err != nil {
		h.logger.Err(err).Msgf(" err conver MERCHANTORDERID %s ", params.MERCHANTORDERID)
		return operations.NewPostPayedInternalServerError()
	}

	pit := pay_incoming_transaction.NewPayTransaction(userId, int(*params.AMOUNT))

	if err := h.ctrl.Exec(ctx, pit); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return operations.NewPostPayedInternalServerError()
	}

	return operations.NewPostPayedOK().WithPayload("YES")
}
