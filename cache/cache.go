package cache

import (
	"errors"
	"fmt"
)

var adapters = make(map[string]Store)

// Register Register a cache adapter available by the adapter name.
func Register(name string, adapter Store) error {
	if adapter == nil {
		return errors.New("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		return errors.New("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
	return nil
}

// NewCache Create a new cache by adapter name.
func Cache(adapter interface{}) (*Repository, error) {
	var store Store
	switch adapter.(type) {
	case string:
		var ok bool
		store, ok = adapters[adapter.(string)]
		if !ok {
			err := fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapter.(string))
			return nil, err
		}
	case Store:
		store = adapter.(Store)
	}

	return NewRepository(store), nil
}
