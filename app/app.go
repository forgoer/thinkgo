package app

import (
	"github.com/thinkoner/thinkgo/router"
	"github.com/thinkoner/thinkgo/view"
)

// Application the ThinkGo Application
type Application struct {
	view  *view.View
	route *router.Route
}

// NewApplication returns a new ThinkGo Application
func NewApplication() *Application {
	return &Application{}
}

// RegisterRoute Register Route for Application
func (a *Application) RegisterRoute(r *router.Route) {
	a.route = r
}

// RegisterView Register View for Application
func (a *Application) RegisterView(v *view.View) {
	a.view = v
}

// GetRoute Get the router of the application
func (a *Application) GetRoute() *router.Route {
	return a.route
}

// GetView Get the view of the application
func (a *Application) GetView() *view.View {
	return a.view
}
