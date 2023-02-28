package contracts

import "github.com/forgoer/thinkgo/router"

type Application interface {
	// Debug determine if the application is running with debug mode enabled.
	Debug() bool
	Route() *router.Route
}
