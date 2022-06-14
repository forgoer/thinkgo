package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisStore struct {
	pool   *redis.Pool // redis connection pool
	prefix string
}

// NewStore Create a redis cache store
func NewRedisStore(pool *redis.Pool, prefix string) *RedisStore {
	s := RedisStore{}
	return s.SetPool(pool).SetPrefix(prefix)
}

// Get get cached value by key.
func (s *RedisStore) Get(key string, val interface{}) error {
	c := s.pool.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", s.prefix+key))
	if err != nil {
		return err
	}

	return json.Unmarshal(b, val)
}

// Put set cached value with key and expire time.
func (s *RedisStore) Put(key string, val interface{}, timeout time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c := s.pool.Get()
	defer c.Close()
	_, err = c.Do("SETEX", s.prefix+key, int64(timeout/time.Second), string(b))
	return err
}

// Increment the value of an item in the cache.
func (s *RedisStore) Increment(key string, value ...int) (int, error) {
	c := s.pool.Get()
	defer c.Close()

	var by = 1
	if len(value) > 0 {
		by = value[0]
	}

	return redis.Int(c.Do("INCRBY", s.prefix+key, by))
}

// Decrement the value of an item in the cache.
func (s *RedisStore) Decrement(key string, value ...int) (int, error) {
	c := s.pool.Get()
	defer c.Close()

	var by = 1
	if len(value) > 0 {
		by = value[0]
	}

	return redis.Int(c.Do("DECRBY", s.prefix+key, by))
}

// Exist check cache's existence in redis.
func (s *RedisStore) Exist(key string) bool {
	c := s.pool.Get()
	defer c.Close()
	v, err := redis.Bool(c.Do("EXISTS", s.prefix+key))
	if err != nil {
		return false
	}
	return v
}

// Expire set value expire time.
func (s *RedisStore) Expire(key string, timeout time.Duration) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", s.prefix+key, int64(timeout/time.Second))

	return err
}

// Forget Remove an item from the cache.
func (s *RedisStore) Forget(key string) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", s.prefix+key)
	return err
}

// Remove all items from the cache.
func (s *RedisStore) Flush() error {
	c := s.pool.Get()
	defer c.Close()

	var err error
	iter := 0
	keys := []string{}

	for {
		arr, err := redis.Values(c.Do("SCAN", iter, "MATCH", s.prefix+"*"))
		if err != nil {
			return err
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}
	for _, key := range keys {
		if _, err = c.Do("DEL", key); err != nil {
			return err
		}
	}
	return err
}

func (s *RedisStore) TTL(key string) (int64, error) {
	c := s.pool.Get()
	defer c.Close()

	return redis.Int64(c.Do("TTL", s.prefix+key))
}

// SetPool Get the redis pool.
func (s *RedisStore) SetPool(pool *redis.Pool) *RedisStore {
	s.pool = pool
	return s
}

// GetPrefix Get the cache key prefix.
func (s *RedisStore) GetPrefix() string {
	return s.prefix
}

// SetPrefix Set the cache key prefix.
func (s *RedisStore) SetPrefix(prefix string) *RedisStore {
	if len(prefix) != 0 {
		s.prefix = fmt.Sprintf("%s:", prefix)
	} else {
		s.prefix = ""
	}
	return s
}
