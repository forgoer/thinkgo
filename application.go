package thinkgo

import (
	"github.com/forgoer/thinkgo/log"
	"github.com/forgoer/thinkgo/log/record"
	"github.com/forgoer/thinkgo/router"
	"github.com/forgoer/thinkgo/view"
)

// Application the ThinkGo Application
type Application struct {
	Env    string
	logger *log.Logger
	view   *view.View
	route  *router.Route
}

// NewApplication returns a new ThinkGo Application
func NewApplication() *Application {
	app := Application{Env: "production"}

	app.logger = log.NewLogger("develop", record.DEBUG)

	app.route = router.New()
	//app.route.Statics(config.Route.Static)

	app.view = view.New()
	//app.view.ParseGlob(config.View.Path)

	return &app
}

// Route return the router of the application
func (a *Application) Logger() *log.Logger {
	return a.logger
}

// Route return the router of the application
func (a *Application) Route() *router.Route {
	return a.route
}

// View Get the view of the application
func (a *Application) View() *view.View {
	return a.view
}

// Debug Determine if the application is running with debug mode enabled.
func (a *Application) Debug() bool {
	return true
}
