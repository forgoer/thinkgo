package router

import (
	"github.com/forgoer/thinkgo/ctx"
	"net/http"
)

// Response an HTTP response interface
type Response interface {
	Send(w http.ResponseWriter)
}

// Closure Anonymous function, Used in Middleware Handler
type Closure func(req *ctx.Request) interface {
}

// MiddlewareFunc Handle an incoming request.
type Middleware func(request *ctx.Request, next Closure) interface{}
