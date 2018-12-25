package config

import "time"

type SessionConfig struct {
	//Default Session Driver
	Driver string

	//Session Cookie Name
	CookieName string

	//Session Lifetime
	Lifetime time.Duration

	//Session Encryption
	Encrypt bool

	//Session File Location
	Files string
}
