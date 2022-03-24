package configure

import "backend-vpn/pkg/controller"

func MainController(ctrl *controller.ControllerImpl) (e error) {

	propogateErr := func(err error) {
		if err != nil {
			e = err
		}
	}

	//	propogateErr(ctrl.RegisterHandler(command.NewPostEventHandler(publisher)))

}
