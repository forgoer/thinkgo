package session

type Handler interface {
	Read(id string) string
	Write(id string, data string)
}
