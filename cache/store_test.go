package cache

import (
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func testCache(t *testing.T, cache Store) {
	var a int
	var b string
	err := cache.Get("a", &a)
	if err == nil {
		t.Error("Getting A found value that shouldn't exist:", a)
	}

	err = cache.Get("b", &b)
	if err == nil {
		t.Error("Getting B found value that shouldn't exist:", b)
	}

	cache.Put("a", 1, 10*time.Minute)
	cache.Put("b", "thinkgo", 10*time.Minute)

	err = cache.Get("a", &a)
	if err != nil {
		t.Error(err)
	}

	if a != 1 {
		t.Error("Expect: ", 1)
	}

	err = cache.Get("b", &b)
	if err != nil {
		t.Error(err)
	}

	if b != "thinkgo" {
		t.Error("Expect: ", "thinkgo")
	}
}

func TestMemoryCache(t *testing.T) {
	Register("memory", NewMemoryStore("thinkgo"))

	cache, err := NewCache("memory")

	if err != nil {
		t.Error(err)
	}
	testCache(t, cache)
}

func TestRedisCache(t *testing.T) {
	pool := &redis.Pool{
		MaxIdle:     5,
		MaxActive:   1000,
		IdleTimeout: 300 * time.Second,
		Wait:        true,
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	cache, err := NewCache(NewRedisStore(pool, "thinkgo"))
	if err != nil {
		t.Error(err)
	}
	testCache(t, cache)
}
