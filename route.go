package thinkgo

import (
	"github.com/forgoer/thinkgo/contracts"
	"github.com/forgoer/thinkgo/ctx"
	"github.com/forgoer/thinkgo/router"
)

// routeDispatcher The default route Dispatcher
type routeDispatcher struct {
	app contracts.Application
}

func (h *routeDispatcher) New(app contracts.Application) {
	h.app = app
}

// Handle the request to a router and return the response.
func (h *routeDispatcher) Handle(request *ctx.Request, next contracts.Next) interface{} {
	rule, err := h.Route().Dispatch(request)

	if err != nil {
		return ctx.NotFoundResponse()
	}

	return router.RunRoute(request, rule)
}

// Route return a router from app
func (h *routeDispatcher) Route() *router.Route {
	return h.app.Route()
}
