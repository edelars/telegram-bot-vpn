package configure

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage/handlers/create_strongswan_account"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"github.com/jmoiron/sqlx"
)

func MainController(ctrl *controller.ControllerImpl, mainDb *sqlx.DB, swDb *sqlx.DB) (e error) {

	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	propogateErr(ctrl.RegisterHandler(get_create_update_user.NewGetCreateUpdateUserHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(create_strongswan_account.NewCreateStrongswanAccountHandler(swDb)))

	return e

}
