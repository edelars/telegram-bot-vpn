package handlers

import (
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/pkg/controller"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
)

type PostPayedHandler struct {
	ctrl   controller.Controller
	logger *zerolog.Logger
}

func NewPostPayedHandler(ctrl controller.Controller, logger *zerolog.Logger) *PostPayedHandler {
	return &PostPayedHandler{ctrl, logger}
}

func (h PostPayedHandler) Handle(params operations.PostPayedParams) middleware.Responder {
	//ctx := params.HTTPRequest.Context()

	//cmd := &deployProject.DeployProjectAsync{
	//	ID:               params.ProjectID,
	//	StopBeforeDeploy: params.BodyParameters.StopBeforeDeploy,
	//}
	//if err := h.ctrl.Exec(ctx, cmd); err != nil {
	//	logger.Warn().Err(err).Msg("failed to deploy project")
	//
	//	if errors.Is(err, model.ErrLocked) {
	//		return deploy.NewPutProjectsProjectIDConflict()
	//	}
	//
	//	return deploy.NewPutProjectsProjectIDInternalServerError()
	//}

	h.logger.Debug().Msgf("project dep %s", params)

	return operations.NewPostPayedOK()
}
