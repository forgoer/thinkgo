package session

import (
	"time"
)

var customHandlers map[string]Handler

type Config struct {
	//Default Session Driver
	Driver string

	CookieName string

	//Session Lifetime
	Lifetime time.Duration

	//Session Encryption
	Encrypt bool

	//Session File Location
	Files string
}

type Manager struct {
	store          *Store
	Config         *Config
}

func NewManager(config *Config) *Manager {
	m := &Manager{
		Config: config,
	}

	return m
}

func (m *Manager) SessionStart(req Request) *Store {
	m.parseStore()

	if handler, ok := m.usingCookieSessions(); ok {
		handler.SetRequest(req)
	}
	name, _ := req.Cookie(m.store.GetName())
	m.store.SetId(name)
	m.store.Start()
	return m.store
}

func (m *Manager) SessionSave(res Response) *Store {
	if handler, ok := m.usingCookieSessions(); ok {
		handler.SetResponse(res)
	}
	res.Cookie(m.store.GetName(), m.store.GetId())
	m.store.Save()
	return m.store
}

func Extend(driver string, handler Handler) {
	if customHandlers == nil {
		customHandlers = make(map[string]Handler)
	}
	customHandlers[driver] = handler
}

func (m *Manager) buildSession(handler Handler) *Store {
	store := NewStore(m.Config.CookieName, handler)
	return store
}

func (m *Manager) usingCookieSessions() (handler *CookieHandler, ok bool) {
	handler, ok = m.store.GetHandler().(*CookieHandler)
	return
}

func (m *Manager) parseStore() {
	if m.store != nil {
		return
	}

	m.store = m.buildSession(
		m.parseStoreHandler(),
	)
}

func (m *Manager) parseStoreHandler() Handler {
	var storeHandler Handler
	switch m.Config.Driver {
	case "cookie":
		storeHandler = &CookieHandler{}
	case "file":
		storeHandler = &FileHandler{
			Path:     m.Config.Files,
			Lifetime: m.Config.Lifetime,
		}
	default:
		var ok bool
		storeHandler, ok = customHandlers[m.Config.Driver]
		if !ok {
			panic("Unsupported session driver: " + m.Config.Driver)
		}
	}

	return storeHandler
}
