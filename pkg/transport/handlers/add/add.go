package add

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

type Add struct {
	logger *zerolog.Logger
	ctrl   controller.Controller
	env    config.Environment
}

func NewAdd(ctrl controller.Controller, logger *zerolog.Logger, env config.Environment) *Add {
	return &Add{
		ctrl:   ctrl,
		logger: logger,
		env:    env,
	}
}

func (p *Add) Endpoint() interface{} {
	return `/add`
}

func (p *Add) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		if err := p.ctrl.Exec(context.Background(), &storage.AccessRightQuery{Id: data.Id}); err != nil {
			p.logger.Debug().Msgf("add: no access %d", data.Id)
			return err.Error()
		}

		msgAr := data.GetMessageArray()
		if data.Message == "" || len(msgAr) < 2 {
			return "Укажите логин"
		}

		err, userS := dto.NewStrongswanUser(msgAr[1], "", true)

		if err != nil {
			return err
		}

		q := storage.CreateStrongswanAccount{
			User: userS,
		}

		if err := p.ctrl.Exec(context.Background(), &q); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint(err)
		}
		p.logger.Debug().Msgf("add: create new account %s", msgAr[1])
		return fmt.Sprintf("password %s", userS.GetPassword())
	}
}

func (p *Add) Menu() (res []transport.MenuI) {
	return nil
}
