package think

import (
	"github.com/forgoer/thinkgo/context"
)

// HandlerFunc Handle the application.
type HandlerFunc func(app *Application) Handler

// Closure Anonymous function, Used in Middleware Handler
type Closure func(req *context.Request) interface{}

// Handler Middleware Handler interface
type Handler interface {
	Process(request *context.Request, next Closure) interface{}
}
