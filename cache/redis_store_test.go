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
			c, err := redis.Dial("tcp", "10.0.41.242:6379")
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", "abc-123"); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 0); err != nil {
				c.Close()
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

func TestRedisStore_Tags(t *testing.T) {
	cache := getRedisStore()
	key := "john"
	value := "LosAngeles"
	err := cache.Tags("people", "artists").Put(key, value, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	var val1 string
	cache.Get(key, &val1)
	if value != val1 {
		t.Errorf("%s != %s", value, val1)
	}

	var val2 string
	cache.Tags("people").Get(key, &val2)
	if value != val2 {
		t.Errorf("%s != %s", value, val2)
	}

	var val3 string
	cache.Tags("artists").Get(key, &val3)
	if value != val3 {
		t.Errorf("%s != %s", value, val3)
	}

	cache.Tags("people").Put("bob", "NewYork", time.Hour)

	err = cache.Tags("artists").Flush()
	if err != nil {
		t.Fatal(err)
	}

	err = cache.Tags("artists").Get(key, &val1)
	if err == nil {
		t.Fatal("err should not be nil")
	}

	cache.Tags("people").Get("bob", &val2)
	if "NewYork" != val2 {
		t.Errorf("%s != %s", "NewYork", val2)
	}
}
