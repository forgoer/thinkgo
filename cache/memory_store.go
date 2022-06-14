package cache

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type item struct {
	Object     interface{}
	Expiration int64
}

// Expired Returns true if the item has expired.
func (item item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type MemoryStore struct {
	prefix       string
	items        map[string]item
	mu           sync.RWMutex
	cleanupTimer *time.Timer
}

// NewStore Create a memory cache store
func NewMemoryStore(prefix string) *MemoryStore {
	s := &MemoryStore{
		items: make(map[string]item),
	}
	return s.SetPrefix(prefix)
}

// Get get cached value by key.
// func (s *Store) Get(key string) (interface{}, error) {
func (s *MemoryStore) Get(key string, val interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[s.prefix+key]
	if !ok {
		return errors.New("not found")
	}

	if item.Expired() {
		return errors.New("expired")
	}

	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid unmarshal")
	}

	rv = rv.Elem()

	rv.Set(reflect.ValueOf(item.Object))

	return nil
}

// Put set cached value with key and expire time.
func (s *MemoryStore) Put(key string, val interface{}, timeout time.Duration) error {
	var e int64
	if timeout > 0 {
		e = time.Now().Add(timeout).UnixNano()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	s.items[s.prefix+key] = item{
		Object:     val,
		Expiration: e,
	}

	if e > 0 {
		s.DeleteExpired()
	}

	return nil
}

// Increment the value of an item in the cache.
func (s *MemoryStore) Increment(key string, value ...int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var by = 1
	if len(value) > 0 {
		by = value[0]
	}

	exist, ok := s.items[s.prefix+key]
	if !ok {
		s.items[s.prefix+key] = item{
			Object: 1 + by,
		}
	} else {
		by = exist.Object.(int) + by
		exist.Object = by
		s.items[s.prefix+key] = exist
	}

	return by, nil
}

// Decrement the value of an item in the cache.
func (s *MemoryStore) Decrement(key string, value ...int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var by = 1
	if len(value) > 0 {
		by = value[0]
	}

	exist, ok := s.items[s.prefix+key]
	if !ok {
		s.items[s.prefix+key] = item{
			Object: 0 - by,
		}
	} else {
		by = exist.Object.(int) - by
		exist.Object = by
		s.items[s.prefix+key] = exist
	}

	return by, nil
}

// Exist check cache's existence in memory.
func (s *MemoryStore) Exist(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[s.prefix+key]

	if item.Expired() {
		return false
	}

	return ok
}

// Expire set value expire time.
func (s *MemoryStore) Expire(key string, timeout time.Duration) error {
	var e int64
	if timeout > 0 {
		e = time.Now().Add(timeout).UnixNano()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.Exist(key) {
		return errors.New("key not exist")
	}

	item := s.items[s.prefix+key]
	item.Expiration = e
	s.items[s.prefix+key] = item

	if e > 0 {
		s.DeleteExpired()
	}

	return nil
}

// Forget Remove an item from the cache.
func (s *MemoryStore) Forget(key string) error {
	delete(s.items, s.prefix+key)
	return nil
}

// Remove all items from the cache.
func (s *MemoryStore) Flush() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.items = map[string]item{}

	return nil
}

func (s *MemoryStore) TTL(key string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[s.prefix+key]
	if !ok {
		return 0, errors.New("not found")
	}

	if item.Expired() {
		return 0, errors.New("not found")
	}

	return item.Expiration - time.Now().UnixNano(), nil
}

// GetPrefix Get the cache key prefix.
func (s *MemoryStore) GetPrefix() string {
	return s.prefix
}

// SetPrefix Set the cache key prefix.
func (s *MemoryStore) SetPrefix(prefix string) *MemoryStore {
	if len(prefix) != 0 {
		s.prefix = fmt.Sprintf("%s:", prefix)
	} else {
		s.prefix = ""
	}
	return s
}

// Delete all expired items from the cache.
func (s *MemoryStore) DeleteExpired() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cleanupTimer != nil {
		s.cleanupTimer.Stop()
	}

	smallestDuration := 0 * time.Nanosecond
	for key, item := range s.items {
		if item.Expiration == 0 {
			continue
		}
		// "Inlining" of expired
		if item.Expired() {
			delete(s.items, key)
		} else {
			// Find the item chronologically closest to its end-of-lifespan.
			sub := item.Expiration - time.Now().UnixNano()

			if smallestDuration == 0 {
				smallestDuration = time.Duration(sub) * time.Nanosecond
			} else {
				if time.Duration(sub)*time.Nanosecond < smallestDuration {
					smallestDuration = time.Duration(sub) * time.Nanosecond
				}
			}
		}
	}

	if smallestDuration > 0 {
		s.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go s.DeleteExpired()
		})
	}
}
