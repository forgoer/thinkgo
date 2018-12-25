package context

type Session interface {
	Get(name string, value ...interface{}) interface{}

	Set(name string, value interface{})

	All() map[string]interface{}

	Remove(name string) interface{}

	Forget(names ...string)

	Clear()

	Save()
}
