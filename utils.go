package thinkgo

import (
	"github.com/thinkoner/thinkgo/context"
	"github.com/thinkoner/thinkgo/config"
	"time"
)

func parseCookieHandler() *context.Cookie {
	return &context.Cookie{
		Config: &context.CookieConfig{
			Prefix:   config.Cookie.Prefix,
			Path:     config.Cookie.Path,
			Domain:   config.Cookie.Domain,
			Expires:  time.Now().Add(config.Cookie.Expires),
			MaxAge:   config.Cookie.MaxAge,
			Secure:   config.Cookie.Secure,
			HttpOnly: config.Cookie.HttpOnly,
		},
	}
}
