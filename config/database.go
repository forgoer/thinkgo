package config

type DatabaseConfig struct {
	Connection  string
	Connections map[string]Connection
}

type Connection struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Charset  string
	Prefix   string
	Engine   string
}
