package configure

import (
	"backend-vpn/pkg/access_right"
	"backend-vpn/pkg/billing/account_info"
	"backend-vpn/pkg/billing/activate_account"
	"backend-vpn/pkg/billing/auto_suggester_tariff_plan"
	"backend-vpn/pkg/billing/deactivate_account"
	"backend-vpn/pkg/billing/pay_get_invoice"
	"backend-vpn/pkg/billing/pay_incoming_transaction"
	"backend-vpn/pkg/billing/pay_prepare"
	"backend-vpn/pkg/config"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/storage/handlers/create_strongswan_account"
	"backend-vpn/pkg/storage/handlers/delete_strongswan_account"
	"backend-vpn/pkg/storage/handlers/get_create_update_user"
	"backend-vpn/pkg/storage/handlers/get_expired_users"
	"backend-vpn/pkg/storage/handlers/get_user"
	"backend-vpn/pkg/storage/handlers/get_user_balance"
	"backend-vpn/pkg/storage/handlers/new_payments"
	"backend-vpn/pkg/storage/handlers/save_user"
	"backend-vpn/pkg/storage/handlers/unchecked_payments_to_balance"
	"backend-vpn/pkg/storage/handlers/writeoff_balance"
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
	propogateErr(ctrl.RegisterHandler(unchecked_payments_to_balance.NewUncheckedPaymentsToBalanceHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(get_user_balance.NewGetUserBalanceHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(get_user.NewGetUserHandler(mainDb, env)))
	propogateErr(ctrl.RegisterHandler(writeoff_balance.NewWriteoffBalanceHandler(mainDb)))
	propogateErr(ctrl.RegisterHandler(get_expired_users.NewGetExpiredUsersHandler(mainDb)))

	//swan
	propogateErr(ctrl.RegisterHandler(create_strongswan_account.NewCreateStrongswanAccountHandler(swDb, env)))
	propogateErr(ctrl.RegisterHandler(delete_strongswan_account.NewDeleteStrongswanAccountHandler(swDb)))

	propogateErr(ctrl.RegisterHandler(access_right.NewAccessRightHandlerHandler(env)))

	//Billing
	propogateErr(ctrl.RegisterHandler(pay_prepare.NewPayPrepareHandler(ctrl, logger)))
	propogateErr(ctrl.RegisterHandler(activate_account.NewActivateAccountHandler(ctrl, logger)))
	propogateErr(ctrl.RegisterHandler(pay_get_invoice.NewPayGetInvoiceHandler(ctrl, logger, provider)))
	propogateErr(ctrl.RegisterHandler(pay_incoming_transaction.NewPayIncomingTransactionHandler(ctrl, logger, workerPayChan)))
	propogateErr(ctrl.RegisterHandler(auto_suggester_tariff_plan.NewAutoSuggesterTariffPlanHandler(env)))
	propogateErr(ctrl.RegisterHandler(deactivate_account.NewActivateAccountHandler(ctrl, logger)))
	propogateErr(ctrl.RegisterHandler(account_info.NewAccountInfoHandler(ctrl, logger, env)))

	return e

}
