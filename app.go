package thinkgo

import (
	"github.com/thinkoner/thinkgo/route"
)

// Application the ThinkGo Application
type Application struct {
	View  View
	Route *route.Route
}

// NewApplication returns a new ThinkGo Application
func NewApplication() *Application {
	return &Application{}
}

// RegisterRoute Register Route for Application
func (a *Application) RegisterRoute(r *route.Route) {
	a.Route = r
}

// RegisterView Register View for Application
func (a *Application) RegisterView(v View) {
	a.View = v
}
