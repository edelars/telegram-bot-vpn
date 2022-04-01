package configure

import (
	"backend-vpn/pkg/access_right"
	"backend-vpn/pkg/billing/activate_account"
	"backend-vpn/pkg/billing/pay_get_invoice"
	"backend-vpn/pkg/billing/pay_incoming_transaction"
	"backend-vpn/pkg/billing/pay_prepare"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/storage/handlers/create_strongswan_account"
	"backend-vpn/pkg/storage/handlers/delete_strongswan_account"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"backend-vpn/pkg/storage/handlers/new_payments"
	"backend-vpn/pkg/storage/handlers/save_user"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func MainController(ctrl *controller.ControllerImpl,
	mainDb *sqlx.DB,
	swDb *sqlx.DB,
	env config.Environment,
	logger *zerolog.Logger,
	provider pay_get_invoice.ProviderI,
	workerPayChan chan *storage.NewPayments) (e error) {

	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	//main db
	propogateErr(ctrl.RegisterHandler(get_create_update_user.NewGetCreateUpdateUserHandler(mainDb, env)))
	propogateErr(ctrl.RegisterHandler(save_user.NewSaveUserHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(new_payments.NewNewPaymentsHandler(mainDb)))

	//swan
	propogateErr(ctrl.RegisterHandler(create_strongswan_account.NewCreateStrongswanAccountHandler(swDb, env)))
	propogateErr(ctrl.RegisterHandler(delete_strongswan_account.NewDeleteStrongswanAccountHandler(swDb)))

	propogateErr(ctrl.RegisterHandler(access_right.NewAccessRightHandlerHandler(env)))

	//Billing
	propogateErr(ctrl.RegisterHandler(pay_prepare.NewPayPrepareHandler(ctrl, logger)))
	propogateErr(ctrl.RegisterHandler(activate_account.NewActivateAccountHandler(ctrl, logger)))
	propogateErr(ctrl.RegisterHandler(pay_get_invoice.NewPayGetInvoiceHandler(ctrl, logger, provider)))
	propogateErr(ctrl.RegisterHandler(pay_incoming_transaction.NewPayIncomingTransactionHandler(ctrl, logger, workerPayChan)))

	return e

}
