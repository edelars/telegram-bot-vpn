package payments_worker

import (
	"backend-vpn/pkg/billing/account_info"
	"backend-vpn/pkg/billing/activate_account"
	"backend-vpn/pkg/billing/auto_suggester_tariff_plan"
	"backend-vpn/pkg/billing/deactivate_account"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/transport"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

const (
	timerDuration = 60
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
				h.HandlePit(ctx, pit)
			}
		}

	}
}

func (h PaymentsWorkerHandler) HandlePit(ctx context.Context, pit *storage.NewPayments) {

	h.SendMsgToClient(pit.UserId, fmt.Sprintf("Получили оплату: %d руб.", pit.Value))

	upq := storage.UncheckedPaymentsQuery{UserId: pit.UserId}
	if err := h.ctrl.Exec(ctx, &upq); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:UncheckedPaymentsQuery fail")
	}

	bal := storage.GetUserBalanceQuery{UserId: pit.UserId}
	if err := h.ctrl.Exec(ctx, &bal); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:GetUserBalanceQuery fail")
	}

	pln := auto_suggester_tariff_plan.AutoSuggesterTariffPlan{HowMuchMoney: int(bal.Out.TotalBalance)}
	if err := h.ctrl.Exec(ctx, &pln); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:AutoSuggesterTariffPlan fail")
	}

	if !pln.Out.Selected {
		h.SendMsgToClient(pit.UserId, fmt.Sprintf("Недостаточно денег на балансе (%d руб.) для активации тарифа", bal.Out.TotalBalance))
		h.logger.Debug().Msgf("PaymentsWorkerHandler not selected tariff plan, return")
		return
	}

	ngu := storage.NewGetUser(pit.UserId)
	if err := h.ctrl.Exec(ctx, ngu); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:NewGetUser fail")
		return
	}

	aa, err := activate_account.NewActivateAccount(ngu.Out.User, pln.Out.TariffDays, false)
	if err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:NewActivateAccount fail")
		return
	}
	if err := h.ctrl.Exec(ctx, aa); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:ActivateAccount fail")
		return
	}

	wob := storage.WriteOffBalance{UserId: pit.UserId, Value: pln.Out.Cost}
	if err := h.ctrl.Exec(ctx, &wob); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:WriteOffBalance fail")

		da, err := deactivate_account.NewDeactivateAccount(ngu.Out.User)
		if err != nil {
			h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:NewDeactivateAccount fail")
			return
		}
		if err := h.ctrl.Exec(ctx, da); err != nil {
			h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:DeactivateAccount fail")
		}
		return
	}
	h.SendMsgToClient(pit.UserId, fmt.Sprintf("Вам активировано дней: %d\nБаланс: %d", pln.Out.TariffDays, int(bal.Out.TotalBalance)-pln.Out.Cost))

	ai := account_info.NewAccountInfoWithData(ngu.Out.User, int(bal.Out.TotalBalance)-pln.Out.Cost)
	if err := h.ctrl.Exec(ctx, ai); err != nil {
		h.logger.Debug().Err(err).Msg("PaymentsWorkerHandler:AccountInfo fail")
		return
	}

	h.SendMsgToClient(pit.UserId, ai.Out.Message)

}

func (h PaymentsWorkerHandler) SendMsgToClient(userId int64, msg string) {
	if err := h.msgBot.Send(userId, msg); err != nil {
		h.logger.Err(err).Msgf("cant send msg to usr with id:%d", userId)
	}
}
