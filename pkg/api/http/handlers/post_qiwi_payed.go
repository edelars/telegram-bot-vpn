package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/internal/dto"
	"backend-vpn/internal/provider/qiwi"
	"backend-vpn/pkg/billing/pay_incoming_transaction"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
	"strconv"
)

type PostQiwiPayedHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
	env    config.Environment
}

func NewPostQiwiPayedHandler(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *PostQiwiPayedHandler {
	return &PostQiwiPayedHandler{ctrl, logger, env}
}

func (h PostQiwiPayedHandler) Handle(params operations.PostQiwiPayedParams) middleware.Responder {

	h.logger.Debug().Msgf("Payed %s", params.QiwiPayOpts)

	ctx := params.HTTPRequest.Context()

	if params.QiwiPayOpts == nil || params.QiwiPayOpts.Bill == nil {
		h.logger.Err(errors.New("nil data in query")).Msg("")
		return operations.NewPostPayedInternalServerError()
	}

	payValue, err := strconv.Atoi(*params.QiwiPayOpts.Bill.Amount.Value)

	qw := qiwi.NewQiwi(h.env)
	pd := dto.PayData{
		Value: payValue,
		Order: *params.QiwiPayOpts.Bill.BillID,
	}
	r := qw.SignVerify(pd, params.XAPISignatureSHA256, *params.QiwiPayOpts.Bill.Status.Value, params.QiwiPayOpts.Bill.SiteID)
	h.logger.Debug().Msgf("SignVerify result: %v", r)
	if r == false {
		return operations.NewPostQiwiPayedInternalServerError()
	}

	userId, err := strconv.ParseInt(*params.QiwiPayOpts.Bill.BillID, 10, 64)
	if err != nil {
		h.logger.Err(err).Msgf(" err convert BillID %s ", params.QiwiPayOpts.Bill.BillID)
		return operations.NewPostQiwiPayedInternalServerError()
	}

	pit := pay_incoming_transaction.NewPayTransaction(userId, payValue)

	if err := h.ctrl.Exec(ctx, pit); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return operations.NewPostQiwiPayedInternalServerError()
	}

	return operations.NewPostQiwiPayedOK()
}
