package cache

import (
	"errors"
	"fmt"
	"time"
)

type Store interface {
	// Get get cached value by key.
	Get(key string, val interface{}) error

	// Put set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error

	// Forget Remove an item from the cache.
	Forget(key string) error

	// Remove all items from the cache.
	Flush() error
}

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
func NewCache(adapter interface{}) (Store, error) {
	var store Store
	switch adapter.(type) {
	case string:
		store, ok := adapters[adapter.(string)]
		if !ok {
			err := fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapter.(string))
			return nil, err
		}
		return store, nil
	case Store:
		store = adapter.(Store)
	}

	return store, nil
}
