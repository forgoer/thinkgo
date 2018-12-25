package session

type Request interface {
	Cookie(key string, value ...string) (string, error)
}

type Response interface {
	Cookie(name interface{}, params ...interface{}) error
}
