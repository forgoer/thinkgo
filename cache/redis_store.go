package cache

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// ReferenceKeyForever Forever reference key.
const ReferenceKeyForever = "forever_ref"

// ReferenceKeyStandard Standard reference key.
const ReferenceKeyStandard = "standard_ref"

type RedisStore struct {
	pool   *redis.Pool // redis connection pool
	tagSet *TagSet
	prefix string
}

// NewRedisStore Create a redis cache store
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

	err = s.pushStandardKeys(key)
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
	err := s.pushStandardKeys(key)
	if err != nil {
		return 0, err
	}

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
	err := s.pushStandardKeys(key)
	if err != nil {
		return 0, err
	}

	c := s.pool.Get()
	defer c.Close()

	var by = 1
	if len(value) > 0 {
		by = value[0]
	}

	return redis.Int(c.Do("DECRBY", s.prefix+key, by))
}

// Forever Store an item in the cache indefinitely.
func (s *RedisStore) Forever(key string, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = s.pushForeverKeys(key)
	if err != nil {
		return err
	}

	c := s.pool.Get()
	defer c.Close()
	_, err = c.Do("SET", s.prefix+key, string(b))
	return err
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
	if s.tagSet != nil {
		err := s.deleteForeverKeys()
		if err != nil {
			return err
		}
		err = s.deleteStandardKeys()
		if err != nil {
			return err
		}
		err = s.tagSet.Reset()
		if err != nil {
			return err
		}
		return nil
	}

	return s.FlushByPrefix("")
}

func (s *RedisStore) FlushByPrefix(prefix string) error {
	c := s.pool.Get()
	defer c.Close()

	var err error
	iter := 0
	keys := []string{}
	pattern := s.prefix + prefix + "*"
	for {
		arr, err := redis.Values(c.Do("SCAN", iter, "MATCH", pattern))
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

	length := len(keys)
	if length == 0 {
		return nil
	}

	var keysChunk []interface{}
	for i, key := range keys {
		keysChunk = append(keysChunk, key)
		if i == length-1 || len(keysChunk) == 1000 {
			_, err = c.Do("DEL", keysChunk...)
			if err != nil {
				return err
			}
			keysChunk = nil
		}
	}

	return nil
}

func (s *RedisStore) Tags(names ...string) Store {
	if len(names) == 0 {
		return s
	}
	ss := s.clone()
	ss.tagSet = NewTagSet(s, names)

	return ss
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

func (s *RedisStore) clone() *RedisStore {
	return &RedisStore{
		pool:   s.pool,
		prefix: s.prefix,
	}
}

func (s *RedisStore) pushStandardKeys(key string) error {
	return s.pushKeys(key, ReferenceKeyStandard)
}

func (s *RedisStore) pushForeverKeys(key string) error {
	return s.pushKeys(key, ReferenceKeyForever)
}

func (s *RedisStore) pushKeys(key, reference string) error {
	if s.tagSet == nil {
		return nil
	}

	namespace, err := s.tagSet.GetNamespace()
	if err != nil {
		return err
	}

	fullKey := s.prefix + key
	segments := strings.Split(namespace, "|")

	c := s.pool.Get()
	defer c.Close()
	for _, segment := range segments {
		_, err = c.Do("SADD", s.referenceKey(segment, reference), fullKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *RedisStore) deleteStandardKeys() error {
	return s.deleteKeysByReference(ReferenceKeyStandard)
}

func (s *RedisStore) deleteForeverKeys() error {
	return s.deleteKeysByReference(ReferenceKeyForever)
}

func (s *RedisStore) deleteKeysByReference(reference string) error {
	if s.tagSet == nil {
		return nil
	}

	namespace, err := s.tagSet.GetNamespace()
	if err != nil {
		return err
	}
	segments := strings.Split(namespace, "|")
	c := s.pool.Get()
	defer c.Close()

	for _, segment := range segments {
		segment = s.referenceKey(segment, reference)
		err = s.deleteKeys(segment)
		if err != nil {
			return err
		}
		_, err = c.Do("DEL", segment)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RedisStore) deleteKeys(referenceKey string) error {
	c := s.pool.Get()
	defer c.Close()
	keys, err := redis.Strings(c.Do("SMEMBERS", referenceKey))
	if err != nil {
		return err
	}
	var length = len(keys)
	if length == 0 {
		return nil
	}

	var keysChunk []interface{}
	for i, key := range keys {
		keysChunk = append(keysChunk, key)
		if i == length-1 || len(keysChunk) == 1000 {
			_, err = c.Do("DEL", keysChunk...)
			if err != nil {
				return err
			}
			keysChunk = nil
		}
	}

	return nil
}

func (s *RedisStore) referenceKey(segment, suffix string) string {
	return s.prefix + segment + ":" + suffix
}
