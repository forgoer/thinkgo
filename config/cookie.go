package config

import "time"

type CookieConfig struct {
	Prefix   string
	Expires  time.Duration
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}
