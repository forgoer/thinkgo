package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func testCache(t *testing.T, cache *Repository) {
	var a int
	var b string
	var c Foo

	cache.Clear()

	assert.Error(t, cache.Get("a", &a))
	assert.Error(t, cache.Get("b", &b))

	assert.NoError(t, cache.Put("a", 1, 10*time.Minute))
	assert.NoError(t, cache.Put("b", "thinkgo", 10*time.Minute))

	assert.True(t, cache.Has("a"))
	assert.True(t, cache.Has("b"))

	assert.NoError(t, cache.Get("a", &a))
	assert.Equal(t, a, 1)
	assert.NoError(t, cache.Get("b", &b))
	assert.Equal(t, b, "thinkgo")

	assert.NoError(t, cache.Pull("b", &b))
	assert.Equal(t, b, "thinkgo")
	assert.False(t, cache.Has("b"))

	assert.NoError(t, cache.Set("b", "think go", 10*time.Minute))
	assert.Error(t, cache.Add("b", "think go", 10*time.Minute))

	assert.True(t, cache.Has("b"))
	assert.NoError(t, cache.Forget("b"))
	assert.False(t, cache.Has("b"))

	assert.NoError(t, cache.Put("c", Foo{
		Name: "thinkgo",
		Age:100,
	}, 10*time.Minute))
	assert.NoError(t,cache.Get("c", &c))
	fmt.Println(c)
	assert.Equal(t, c.Name , "thinkgo")
	assert.Equal(t, c.Age , 100)
	assert.NoError(t, cache.Delete("c"))
	assert.False(t, cache.Has("c"))

	_, ok := cache.GetStore().(Store)
	assert.True(t, ok)

	assert.NoError(t, cache.Clear())
	assert.False(t, cache.Has("a"))
	assert.False(t, cache.Has("b"))

	assert.NoError(t, cache.Remember("a", &a, 1*time.Minute, func() interface{} {
		return 1000
	}))

	assert.Equal(t, a, 1000)

	assert.NoError(t,cache.Remember("b", &b, 1*time.Minute, func() interface{} {
		return "hello thinkgo"
	}))

	assert.Equal(t, b, "hello thinkgo")
}

func TestMemoryCache(t *testing.T) {
	Register("memory", NewMemoryStore("thinkgo"))

	cache, err := Cache("memory")

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
			// if _, err := c.Do("AUTH", "123456"); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, nil
		},
	}

	cache, err := Cache(NewRedisStore(pool, "thinkgo"))
	if err != nil {
		t.Error(err)
	}
	testCache(t, cache)
}