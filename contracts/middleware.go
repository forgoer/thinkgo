package contracts

import "github.com/forgoer/thinkgo/ctx"

// Middleware interface
type Middleware interface {
	// New a Middleware with the application.
	New(app Application)

	// Handle an incoming request.
	Handle(request *ctx.Request, next Next) interface{}
}

// Next Anonymous function, Used in Middleware Handler
type Next func(req *ctx.Request) interface{}
