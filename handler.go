package thinkgo

import (
	"reflect"

	"github.com/thinkoner/thinkgo/config"
	"github.com/thinkoner/thinkgo/context"
	"github.com/thinkoner/thinkgo/route"
	"github.com/thinkoner/thinkgo/session"
)

type RouteHandler struct {
	Route *route.Route
}

// NewRouteHandler The default RouteHandler
func NewRouteHandler(app *Application) Handler {
	return &RouteHandler{
		Route: app.Route,
	}
}

func (h *RouteHandler) Process(request *context.Request, next Closure) interface{} {

	rule, err := h.Route.Dispatch(request)

	if err != nil {
		return NotFoundResponse()
	}

	handler := rule.Handler

	if h, ok := handler.(Response); ok {
		return h
	}

	var res interface{}

	if handler != nil {

		t := reflect.TypeOf(handler)
		switch t.Kind() {
		case reflect.Func:

			v := reflect.ValueOf(handler).Call(
				parseParams(request, rule.Parameters),
			)
			res = v[0].Interface()
		default:
			res = handler
		}
	}
	return res
}

func parseParams(request *context.Request, parameters map[string]string) []reflect.Value {
	var params []reflect.Value
	params = append(params, reflect.ValueOf(request))

	for _, v := range parameters {
		params = append(params, reflect.ValueOf(v))
	}

	return params
}

type SessionHandler struct {
	manager *session.Manager
	app     *Application
}

//SessionHandler The default SessionHandler
func NewSessionHandler(app *Application) Handler {
	handler := &SessionHandler{}
	handler.manager = session.NewManager(&session.Config{
		Driver:     config.Session.Driver,
		CookieName: config.Session.CookieName,
		Lifetime:   config.Session.Lifetime,
		Encrypt:    config.Session.Encrypt,
		Files:      config.Session.Files,
	})

	handler.app = app

	return handler
}

func (h *SessionHandler) Process(req *context.Request, next Closure) interface{} {
	store := h.startSession(req)

	req.SetSession(store)

	result := next(req)

	if res, ok := result.(session.Response); ok {
		h.saveSession(res)
	}

	return result
}

func (h *SessionHandler) startSession(req *context.Request) *session.Store {
	return h.manager.SessionStart(req)
}

func (h *SessionHandler) saveSession(res session.Response) {
	h.manager.SessionSave(res)
}
