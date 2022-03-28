package del

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/transport"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type Del struct {
	logger *zerolog.Logger
	ctrl   controller.Controller
	env    config.Environment
}

func NewDel(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Del {
	return &Del{
		ctrl:   ctrl,
		logger: logger,
		env:    env,
	}
}

func (p *Del) Endpoint() interface{} {
	return `/del`
}

func (p *Del) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		if err := p.ctrl.Exec(context.Background(), &storage.AccessRightQuery{Id: data.Id}); err != nil {
			p.logger.Debug().Msgf("add: no access %d", data.Id)
			return err.Error()
		}

		msgAr := data.GetMessageArray()
		if data.Message == "" || len(msgAr) < 2 {
			return "Укажите логин"
		}

		err, userSw := dto.NewStrongswanUser(msgAr[1], "", true)

		if err != nil {
			return err
		}
		q := storage.DeleteStrongswanAccount{
			User: userSw,
		}

		if err := p.ctrl.Exec(context.Background(), &q); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint(err)
		}
		p.logger.Debug().Msgf("del: delete account %s", msgAr[1])
		return fmt.Sprintf("deleted %s", msgAr[1])
	}
}

func (p *Del) Menu() (res []transport.MenuI) {
	return nil
}
