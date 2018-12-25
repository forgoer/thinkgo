package thinkgo

type View interface {
	Render(name string, data interface{}) []byte
}
