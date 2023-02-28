package thinkgo

import (
	"fmt"
	"net/http"

	"time"

	"github.com/forgoer/thinkgo/contracts"

	"github.com/forgoer/thinkgo/helper"
	"github.com/forgoer/thinkgo/router"
)

type RouteFunc func(route *router.Route)

type Think struct {
	App      *Application
	handlers []contracts.Middleware
}

// New Create The Application
func New() *Think {
	application := NewApplication()

	t := &Think{
		App: application,
	}
	return t
}

// Route register Route
func (th *Think) Route(register RouteFunc) {
	route := th.App.Route()
	defer route.Register()
	register(route)
}

// RegisterConfig Register Config
func (th *Think) RegisterHandler(handler contracts.Middleware) {
	th.handlers = append(th.handlers, handler)
}

// Run thinkgo application.
// Run() default run on HttpPort
// Run("localhost")
// Run(":9011")
// Run("127.0.0.1:9011")
func (th *Think) Run(params ...string) {
	var err error
	var endRunning = make(chan bool, 1)

	var addrs = helper.ParseAddr(params...)

	pl := NewPipeline()
	for _, h := range th.handlers {
		h.New(th.App)
		pl.Pipe(h)
	}

	// register route handler
	pl.Then(th.dispatchToRouter())

	th.App.Logger().Debug("\r\nLoaded routes:\r\n%s", string(th.App.Route().Dump()))

	go func() {
		th.App.Logger().Debug("ThinkGo server running on http://%s", addrs)

		err = http.ListenAndServe(addrs, pl)

		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(100 * time.Microsecond)
			endRunning <- true
		}
	}()

	<-endRunning
}

func (th *Think) dispatchToRouter() contracts.Middleware {
	h := &routeDispatcher{}
	h.New(th.App)

	return h
}

//func (th *Think) bootView() {
//	v := view.NewView()
//	v.SetPath(config.View.Path)
//	th.App.RegisterView(v)
//}
