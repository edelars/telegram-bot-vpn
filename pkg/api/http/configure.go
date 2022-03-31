package http

import (
	"backend-vpn/internal/api/restapi"
	"backend-vpn/internal/api/restapi/operations"
	"backend-vpn/pkg/api/http/handlers"
	"backend-vpn/pkg/controller"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/rs/zerolog"
)

func NewServer(host string, port int, ctrl controller.Controller, logger *zerolog.Logger) (*restapi.Server, error) {
	spec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, err
	}

	api := operations.NewBackendAPI(spec)

	//	api.Logger = logger

	//api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.PostPayedHandler = handlers.NewPostPayedHandler(ctrl, logger)
	api.PostNotifyHandler = handlers.NewNotifyHandler(ctrl, logger)
	api.PostTryagainHandler = handlers.NewTryAgainHandler(ctrl, logger)

	api.ServerShutdown = func() {}

	//api.Middleware = func(builder middleware.Builder) http.Handler {
	//	docPath := "/docs"
	//
	//	middlewares := []func(h http.Handler) http.Handler{
	//		middlewarePprof,
	//		middlewareRedoc(docPath, spec),
	//		middlewareSpec(docPath, spec),
	//		middlewareHealthz(healthchecks...),
	//		middlewareMetrics,
	//		middlewareRecover,
	//		middlewareRequestID,
	//		middlewareLogging,
	//	}
	//
	//	return setupMiddleware(api.Context().RoutesHandler(builder), middlewares...)
	//}

	server := restapi.NewServer(api)
	server.Host = host
	server.Port = port

	return server, nil
}
