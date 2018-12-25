package config

import "time"

var App *AppConfig
var Route *RouteConfig
var View *ViewConfig

//var Database *DatabaseConfig
var Cookie *CookieConfig
var Session *SessionConfig

func init() {
	loadAppConfig()
	loadViewConfig()
	loadRouteConfig()
	loadCookieConfig()
	loadSessionConfig()
}

func loadAppConfig() {
	App = &AppConfig{
		Name:  "ThinkGo",
		Env:   "production",
		Debug: false,
	}
}

func loadRouteConfig() {
	Route = &RouteConfig{
		Static: map[string]string{
			"static": "public",
			"upload": "public",
		},
	}
}

func loadViewConfig() {
	View = &ViewConfig{
		Path: "view",
	}
}

func loadCookieConfig() {
	Cookie = &CookieConfig{
		Prefix:   "",
		Expires:  time.Hour * 4,
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}
}

func loadSessionConfig() {
	Session = &SessionConfig{
		Driver:     "file",
		Lifetime:   time.Hour * 4,
		Encrypt:    false,
		Files:      "temp/sessions",
		CookieName: "thinkgo_sessions",
	}
}
