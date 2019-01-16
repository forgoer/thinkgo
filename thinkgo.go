package thinkgo

import (
	"fmt"
	"net/http"

	"time"

	"github.com/thinkoner/thinkgo/config"
	"github.com/thinkoner/thinkgo/helper"
	"github.com/thinkoner/thinkgo/route"
	"github.com/thinkoner/thinkgo/view"
)

var application *Application

type handlerFunc func(app *Application) Handler

type registerRouteFunc func(route *route.Route)

type registerConfigFunc func()

type ThinkGo struct {
	App      *Application
	handlers []handlerFunc
}

// BootStrap Create The Application
func BootStrap() *ThinkGo {
	application = NewApplication()
	think := &ThinkGo{
		App: application,
	}
	think.bootView()
	think.bootSession()
	think.bootRoute()
	return think
}

func (think *ThinkGo) bootView() {
	v := view.NewView()
	v.SetPath(config.View.Path)
	think.App.RegisterView(v)

}

func (think *ThinkGo) bootSession() {
	think.handlers = append(think.handlers, NewSessionHandler)
}

func (think *ThinkGo) bootRoute() {
	r := route.NewRoute()

	r.Statics(config.Route.Static)

	think.App.RegisterRoute(r)
	think.handlers = append(think.handlers, NewRouteHandler)
}

// RegisterRoute Register Route
func (think *ThinkGo) RegisterRoute(register registerRouteFunc) {
	register(think.App.Route)
}

// RegisterConfig Register Config
func (think *ThinkGo) RegisterConfig(register registerConfigFunc) {
	register()
}

// Run thinkgo application.
// Run() default run on HttpPort
// Run("localhost")
// Run(":9011")
// Run("127.0.0.1:9011")
func (think *ThinkGo) Run(params ...string) {
	var err error
	var endRunning = make(chan bool, 1)

	var addrs = helper.ParseAddr(params...)

	pipeline := NewPipeline()

	for _, h := range think.handlers {
		pipeline.Pipe(h(think.App))
	}

	go func() {

		fmt.Printf("context server Running on http://%s", addrs)

		err = http.ListenAndServe(addrs, pipeline)

		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(100 * time.Microsecond)
			endRunning <- true
		}
	}()

	<-endRunning
}
