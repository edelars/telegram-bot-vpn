// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostPayedHandlerFunc turns a function with the right signature into a post payed handler
type PostPayedHandlerFunc func(PostPayedParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostPayedHandlerFunc) Handle(params PostPayedParams) middleware.Responder {
	return fn(params)
}

// PostPayedHandler interface for that can handle valid post payed params
type PostPayedHandler interface {
	Handle(PostPayedParams) middleware.Responder
}

// NewPostPayed creates a new http.Handler for the post payed operation
func NewPostPayed(ctx *middleware.Context, handler PostPayedHandler) *PostPayed {
	return &PostPayed{Context: ctx, Handler: handler}
}

/* PostPayed swagger:route POST /payed postPayed

PostPayed post payed API

*/
type PostPayed struct {
	Context *middleware.Context
	Handler PostPayedHandler
}

func (o *PostPayed) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostPayedParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
