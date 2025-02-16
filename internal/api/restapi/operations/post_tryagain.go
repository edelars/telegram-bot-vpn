// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostTryagainHandlerFunc turns a function with the right signature into a post tryagain handler
type PostTryagainHandlerFunc func(PostTryagainParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostTryagainHandlerFunc) Handle(params PostTryagainParams) middleware.Responder {
	return fn(params)
}

// PostTryagainHandler interface for that can handle valid post tryagain params
type PostTryagainHandler interface {
	Handle(PostTryagainParams) middleware.Responder
}

// NewPostTryagain creates a new http.Handler for the post tryagain operation
func NewPostTryagain(ctx *middleware.Context, handler PostTryagainHandler) *PostTryagain {
	return &PostTryagain{Context: ctx, Handler: handler}
}

/* PostTryagain swagger:route POST /tryagain postTryagain

PostTryagain post tryagain API

*/
type PostTryagain struct {
	Context *middleware.Context
	Handler PostTryagainHandler
}

func (o *PostTryagain) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostTryagainParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
