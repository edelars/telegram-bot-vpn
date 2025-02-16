package pay_incoming_transaction

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"context"
	"github.com/rs/zerolog"
)

type PayTransaction struct {
	userId int64
	value  int
}

func NewPayTransaction(userId int64, value int) *PayTransaction {
	return &PayTransaction{
		userId: userId,
		value:  value,
	}

}

type PayIncomingTransactionHandler struct {
	ctrl          controller.Controller
	logger        *zerolog.Logger
	workerPayChan chan *storage.NewPayments
}

func NewPayIncomingTransactionHandler(ctrl controller.Controller, logger *zerolog.Logger, workerPayChan chan *storage.NewPayments) *PayIncomingTransactionHandler {
	return &PayIncomingTransactionHandler{ctrl: ctrl, logger: logger, workerPayChan: workerPayChan}
}

func (h *PayIncomingTransactionHandler) Exec(ctx context.Context, args *PayTransaction) (err error) {

	q := storage.NewPayments{
		UserId: args.userId,
		Value:  args.value,
	}

	if err := h.ctrl.Exec(ctx, &q); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return err
	}
	go func() {
		h.workerPayChan <- &q
	}()

	return nil
}

func (h *PayIncomingTransactionHandler) Context() interface{} {
	return (*PayTransaction)(nil)
}
