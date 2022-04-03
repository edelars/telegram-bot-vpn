package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/internal/dto"
	"backend-vpn/internal/provider/freekassa"
	"backend-vpn/pkg/billing/pay_incoming_transaction"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
	"strconv"
)

type PostPayedHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
	env    config.Environment
}

func NewPostPayedHandler(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *PostPayedHandler {
	return &PostPayedHandler{ctrl, logger, env}
}

func (h PostPayedHandler) Handle(params operations.PostPayedParams) middleware.Responder {

	h.logger.Debug().Msgf("Payed %s", params)

	ctx := params.HTTPRequest.Context()

	if params.MERCHANTORDERID == nil || params.AMOUNT == nil {
		h.logger.Err(errors.New("nil data in query")).Msg("")
		return operations.NewPostPayedInternalServerError()
	}

	if params.SIGN != nil {
		fk := freekassa.NewFreeKassa(h.env)
		pd := dto.PayData{
			Value: int(*params.AMOUNT),
			Order: *params.MERCHANTORDERID,
		}
		r := fk.SignVerify(pd, *params.SIGN)
		h.logger.Debug().Msgf("SignVerify result: %v", r)
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
