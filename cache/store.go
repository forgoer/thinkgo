package cache

import (
	"time"
)

type Store interface {
	// Get get cached value by key.
	Get(key string, val interface{}) error

	// Put set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error

	// Exist check cache's existence in redis.
	Exist(key string) bool

	// Forget Remove an item from the cache.
	Forget(key string) error

	// Flush Remove all items from the cache.
	Flush() error
}

