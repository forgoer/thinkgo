package middleware

import (
	"github.com/forgoer/thinkgo/config"
	. "github.com/forgoer/thinkgo/contracts"
	"github.com/forgoer/thinkgo/ctx"
	"github.com/forgoer/thinkgo/session"
)

type Session struct {
	Manager *session.Manager
	app     Application
}

// New the default Session handler
func (h *Session) New(app Application) {
	h.Manager = session.NewManager(&session.Config{
		Driver:     config.Session.Driver,
		CookieName: config.Session.CookieName,
		Lifetime:   config.Session.Lifetime,
		Encrypt:    config.Session.Encrypt,
		Files:      config.Session.Files,
	})

	h.app = app
}

// Handle an incoming request.
func (h *Session) Handle(req *ctx.Request, next Next) interface{} {
	store := h.startSession(req)

	req.SetSession(store)

	result := next(req)

	if res, ok := result.(session.Response); ok {
		h.saveSession(res)
	}

	return result
}

func (h *Session) startSession(req *ctx.Request) *session.Store {
	return h.Manager.SessionStart(req)
}

func (h *Session) saveSession(res session.Response) {
	h.Manager.SessionSave(res)
}
