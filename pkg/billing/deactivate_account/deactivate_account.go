package deactivate_account

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/storage/handlers/delete_strongswan_account"
	"context"
	"errors"
	"github.com/rs/zerolog"
	"time"
)

type DeactivateAccount struct {
	user *dto.User
}

func NewDeactivateAccount(user *dto.User) (a *DeactivateAccount, err error) {

	if user != nil {
		return &DeactivateAccount{user: user}, err
	}

	return a, errors.New("*user is empty")
}

func (a *DeactivateAccount) GetUser() *dto.User {
	return a.user
}

type deactivateAccountHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewActivateAccountHandler(ctrl controller.Controller, logger *zerolog.Logger) *deactivateAccountHandler {
	return &deactivateAccountHandler{ctrl: ctrl, logger: logger}
}

func (h *deactivateAccountHandler) Exec(ctx context.Context, args *DeactivateAccount) (err error) {

	var newUserDto dto.User
	newUserDto = *args.user

	newUserDto.ExpiredAt = time.Now()

	//0
	err, ssu := dto.NewStrongswanUser(newUserDto.Login, newUserDto.Password, false)
	if err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return err
	}

	//1
	if err := h.ctrl.Exec(context.Background(), &storage.DeleteStrongswanAccount{User: ssu}); err != nil {
		if err != delete_strongswan_account.ErrNotExist {
			h.logger.Debug().Err(err).Msg("fail")
			return err
		}
	}

	//2
	err, su := storage.NewSaveUserQuery(&newUserDto)
	if err != nil {
		return err
	}

	if err := h.ctrl.Exec(ctx, su); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return err
	}

	args.user = nil
	args.user = &newUserDto
	return err
}

func (h *deactivateAccountHandler) Context() interface{} {
	return (*DeactivateAccount)(nil)
}
