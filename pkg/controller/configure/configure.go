package configure

import (
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"github.com/jmoiron/sqlx"
)

func MainController(ctrl *controller.ControllerImpl, db *sqlx.DB) (e error) {

	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	propogateErr(ctrl.RegisterHandler(get_create_update_user.NewGetCreateUpdateUserHandler(db)))

	return e

}
