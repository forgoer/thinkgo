package thinkgo

import (
	"fmt"
	"net/http"

	"time"

	"github.com/thinkoner/thinkgo/app"
	"github.com/thinkoner/thinkgo/config"
	"github.com/thinkoner/thinkgo/helper"
	"github.com/thinkoner/thinkgo/router"
	"github.com/thinkoner/thinkgo/view"
)

var application *app.Application

type registerRouteFunc func(route *router.Route)

type registerConfigFunc func()

type ThinkGo struct {
	App      *app.Application
	handlers []app.HandlerFunc
}

// BootStrap Create The Application
func BootStrap() *ThinkGo {
	application = app.NewApplication()
	think := &ThinkGo{
		App: application,
	}
	think.bootView()
	think.bootSession()
	think.bootRoute()
	return think
}

// RegisterRoute Register Route
func (think *ThinkGo) RegisterRoute(register registerRouteFunc) {
	route := think.App.GetRoute()
	defer route.Register()
	register(route)
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

func (think *ThinkGo) bootView() {
	v := view.NewView()
	v.SetPath(config.View.Path)
	think.App.RegisterView(v)
}

func (think *ThinkGo) bootSession() {
	think.handlers = append(think.handlers, app.NewSessionHandler)
}

func (think *ThinkGo) bootRoute() {
	r := router.NewRoute()
	r.Statics(config.Route.Static)
	think.App.RegisterRoute(r)
	think.handlers = append(think.handlers, app.NewRouteHandler)
}
