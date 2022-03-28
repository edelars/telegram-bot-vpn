package configure

import (
	"backend-vpn/pkg/access_right"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage/handlers/create_strongswan_account"
	"backend-vpn/pkg/storage/handlers/delete_strongswan_account"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"github.com/jmoiron/sqlx"
)

func MainController(ctrl *controller.ControllerImpl, mainDb *sqlx.DB, swDb *sqlx.DB, env config.Environment) (e error) {

	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	propogateErr(ctrl.RegisterHandler(get_create_update_user.NewGetCreateUpdateUserHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(create_strongswan_account.NewCreateStrongswanAccountHandler(swDb)))
	propogateErr(ctrl.RegisterHandler(delete_strongswan_account.NewDeleteStrongswanAccountHandler(swDb)))
	propogateErr(ctrl.RegisterHandler(access_right.NewAccessRightHandlerHandler(env)))

	return e

}
