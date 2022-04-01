package payments_worker

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/transport"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

const (
	timerDuration = 5
)

type PaymentsWorkerHandler struct {
	ctrl         controller.Controller
	logger       *zerolog.Logger
	cancelFunc   context.CancelFunc
	incomingChan chan *storage.NewPayments
	msgBot       transport.Transport
}

func NewPaymentsWorkerHandler(ctrl controller.Controller, logger *zerolog.Logger, incomingChan chan *storage.NewPayments, msgBot transport.Transport) *PaymentsWorkerHandler {
	return &PaymentsWorkerHandler{ctrl: ctrl, logger: logger, incomingChan: incomingChan, msgBot: msgBot}
}

func (h *PaymentsWorkerHandler) Run() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	h.cancelFunc = cancelFunc

	go h.start(ctx)
}

func (h *PaymentsWorkerHandler) Shutdown() error {
	h.logger.Debug().Msg("PaymentsWorkerHandler Got cancelFunc")
	h.cancelFunc()
	return nil
}

func (h *PaymentsWorkerHandler) start(ctx context.Context) {

	timer := time.NewTicker(time.Second * timerDuration)
	for {

		select {
		case <-ctx.Done():
			timer.Stop()
			close(h.incomingChan)
			h.logger.Debug().Msg("PaymentsWorkerHandler Got ctx.Done")
			return
		case <-timer.C:
			h.logger.Debug().Msg("PaymentsWorkerHandler Got timer")

		case pit, ok := <-h.incomingChan:
			h.logger.Debug().Msg("PaymentsWorkerHandler Got incomingChan")
			if ok {
				h.HandlePit(pit)
			}
		}

	}
}

func (h PaymentsWorkerHandler) HandlePit(pit *storage.NewPayments) {

	h.SendMsgToClient(pit.UserId, fmt.Sprintf("Получили оплату: %d руб.", pit.Value))

}

func (h PaymentsWorkerHandler) SendMsgToClient(userId int64, msg string) {
	if err := h.msgBot.Send(userId, msg); err != nil {
		h.logger.Err(err).Msgf("cant send msg to usr with id:%d", userId)
	}
}
