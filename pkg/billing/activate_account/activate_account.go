package activate_account

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

type ActivateAccount struct {
	user     *dto.User
	itsTrial bool
	toDays   int
}

func NewActivateAccount(user *dto.User, toDays int, itsTrial bool) (a *ActivateAccount, err error) {

	if user != nil {
		return &ActivateAccount{user: user, toDays: toDays, itsTrial: itsTrial}, err
	}

	return a, errors.New("*user is empty")
}

func (a *ActivateAccount) GetUser() *dto.User {
	return a.user
}

type activateAccountHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewActivateAccountHandler(ctrl controller.Controller, logger *zerolog.Logger) *activateAccountHandler {
	return &activateAccountHandler{ctrl: ctrl, logger: logger}
}

func (h *activateAccountHandler) Exec(ctx context.Context, args *ActivateAccount) (err error) {

	if args.toDays < 0 {
		return errors.New("days too low")
	}

	var newUserDto dto.User
	newUserDto = *args.user

	addTime := time.Duration(args.toDays) * 24 * time.Hour

	if newUserDto.ExpiredAt.Before(time.Now()) {
		newUserDto.ExpiredAt = time.Now().Add(addTime)
	} else {
		newUserDto.ExpiredAt = newUserDto.ExpiredAt.Add(addTime)
	}
	newUserDto.UsedTestPeriod = args.itsTrial

	var genNewPass bool
	if newUserDto.Password == "" {
		genNewPass = true
	}

	//0
	err, ssu := dto.NewStrongswanUser(newUserDto.Login, newUserDto.Password, genNewPass)
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
	if err := h.ctrl.Exec(context.Background(), &storage.CreateStrongswanAccount{User: ssu}); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return err
	}
	newUserDto.Password = ssu.GetPassword()

	//3
	err, su := storage.NewSaveUserQuery(&newUserDto)
	if err != nil {
		return err
	}

	if err := h.ctrl.Exec(context.Background(), su); err != nil {
		h.logger.Debug().Err(err).Msg("fail")
		return err
	}

	args.user = nil
	args.user = &newUserDto
	return err
}

func (h *activateAccountHandler) Context() interface{} {
	return (*ActivateAccount)(nil)
}
