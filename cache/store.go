package cache

import (
	"time"
)

type Store interface {
	// Get get cached value by key.
	Get(key string, val interface{}) error

	// Put set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error

	// Increment the value of an item in the cache.
	Increment(key string, value ...int) (int, error)

	// Decrement the value of an item in the cache.
	Decrement(key string, value ...int) (int, error)

	// Exist check cache's existence in redis.
	Exist(key string) bool

	// Expire set value expire time.
	Expire(key string, timeout time.Duration) error

	// Forget Remove an item from the cache.
	Forget(key string) error

	// Flush Remove all items from the cache.
	Flush() error

	// TTL get the ttl of the key.
	TTL(key string) (int64, error)
}
