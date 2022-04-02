package auto_suggester_tariff_plan

import (
	"backend-vpn/pkg/config"
	"context"
	"errors"
)

type AutoSuggesterTariffPlan struct {
	HowMuchMoney int
	Out          struct {
		Selected   bool
		TariffDays int
		Cost       int
	}
}

type autoSuggesterTariffPlanHandler struct {
	env config.Environment
}

func NewAutoSuggesterTariffPlanHandler(env config.Environment) *autoSuggesterTariffPlanHandler {
	return &autoSuggesterTariffPlanHandler{env: env}
}

func (h *autoSuggesterTariffPlanHandler) Exec(ctx context.Context, args *AutoSuggesterTariffPlan) (err error) {

	switch {
	case args.HowMuchMoney == 0:
		return errors.New("no money - no honey")
	case args.HowMuchMoney >= h.env.Price12:
		args.Out.TariffDays = 365
		args.Out.Cost = h.env.Price12
	case args.HowMuchMoney >= h.env.Price06:
		args.Out.TariffDays = 180
		args.Out.Cost = h.env.Price06
	case args.HowMuchMoney >= h.env.Price01:
		args.Out.TariffDays = 30
		args.Out.Cost = h.env.Price01
	default:
		args.Out.TariffDays = 0
		args.Out.Cost = 0
	}

	args.Out.Selected = args.Out.Cost > 0

	return nil
}

func (h *autoSuggesterTariffPlanHandler) Context() interface{} {
	return (*AutoSuggesterTariffPlan)(nil)
}
