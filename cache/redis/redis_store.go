package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Store struct {
	pool   *redis.Pool // redis connection pool
	prefix string
}

func NewStore(pool *redis.Pool, prefix string) *Store {
	s := Store{}
	return s.SetPool(pool).SetPrefix(prefix)
}

// Get get cached value by key.
func (s *Store) Get(key string, val interface{}) error {
	c := s.pool.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", s.prefix+key))
	if err != nil {
		return err
	}

	return json.Unmarshal(b, val)
}

// Put set cached value with key and expire time.
func (s *Store) Put(key string, val interface{}, timeout time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c := s.pool.Get()
	defer c.Close()
	_, err = c.Do("SETEX", s.prefix+key, int64(timeout/time.Second), string(b))
	return err
}

// Exist check cache's existence in redis.
func (s *Store) Exist(key string) bool {
	c := s.pool.Get()
	defer c.Close()
	v, err := redis.Bool(c.Do("EXISTS", s.prefix+key))
	if err != nil {
		return false
	}
	return v
}

// Forget Remove an item from the cache.
func (s *Store) Forget(key string) error {
	c := s.pool.Get()
	defer c.Close()
	_, err := c.Do("DEL", s.prefix+key)
	return err
}

// Remove all items from the cache.
func (s *Store) Flush() error {
	c := s.pool.Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("KEYS", s.prefix+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		if _, err = c.Do("DEL", key); err != nil {
			return err
		}
	}
	return err
}

// SetPool Get the redis pool.
func (s *Store) SetPool(pool *redis.Pool) *Store {
	s.pool = pool
	return s
}

// GetPrefix Get the cache key prefix.
func (s *Store) GetPrefix() string {
	return s.prefix
}

// SetPrefix Set the cache key prefix.
func (s *Store) SetPrefix(prefix string) *Store {
	if len(prefix) != 0 {
		s.prefix = fmt.Sprintf("%s:", prefix)
	} else {
		s.prefix = ""
	}
	return s
}
