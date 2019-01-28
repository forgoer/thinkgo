package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetPool() *redis.Pool {
	return &redis.Pool{
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
}

func getRedisStore() *RedisStore {
	pool := GetPool()
	return NewRedisStore(pool, "cache")
}

func TestRedisStoreInt(t *testing.T) {
	s := getRedisStore()
	err := s.Put("int", 9811, 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	var v int

	err = s.Get("int", &v)
	if err != nil {
		t.Error(err)
	}

	t.Logf("int:%d", v)
}

func TestRedisStoreString(t *testing.T) {
	s := getRedisStore()
	err := s.Put("str", "this is a string", 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	var str string

	err = s.Get("str", &str)
	if err != nil {
		t.Error(err)
	}

	t.Logf("str:%s", str)
}

func TestStoreStruct(t *testing.T) {
	s := getRedisStore()
	err := s.Put(
		"user", CacheUser{
			Name: "alice",
			Age:  16,
		},
		10*time.Minute,
	)
	if err != nil {
		t.Error(err)
	}

	user := &CacheUser{}

	err = s.Get("user", user)
	if err != nil {
		t.Error(err)
	}

	t.Logf("user:name=%s,age=%d", user.Name, user.Age)
}

func TestRedisStoreForgetAndExist(t *testing.T) {
	s := getRedisStore()
	err := s.Put("forget", "Forget me", 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	exist := s.Exist("forget")
	if exist != true {
		t.Error(errors.New("Expect true"))
	}

	err = s.Forget("forget")
	if err != nil {
		t.Error(err)
	}

	exist = s.Exist("forget")
	if exist == true {
		t.Error(errors.New("Expect false"))
	}
}

func TestRedisStoreFlush(t *testing.T) {
	s := getRedisStore()
	err := s.Put("Flush", "Flush all", 10*time.Minute)
	if err != nil {
		t.Error(err)
	}

	exist := s.Exist("Flush")
	if exist != true {
		t.Error(errors.New("Expect true"))
	}

	err = s.Flush()
	if err != nil {
		t.Error(err)
	}

	exist = s.Exist("Flush")
	if exist == true {
		t.Error(errors.New("Expect false"))
	}
}
