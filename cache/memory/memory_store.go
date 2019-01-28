package memory

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Item struct {
	Object     interface{}
	Expiration int64
}

// Expired Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type Store struct {
	prefix       string
	items        map[string]Item
	mu           sync.RWMutex
	cleanupTimer *time.Timer
}

// NewStore Create a memory cache store
func NewStore(prefix string) *Store {
	s := &Store{
		items: make(map[string]Item),
	}
	return s.SetPrefix(prefix)
}

// Get get cached value by key.
// func (s *Store) Get(key string) (interface{}, error) {
func (s *Store) Get(key string, val interface{}) error {
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
func (s *Store) Put(key string, val interface{}, timeout time.Duration) error {
	var e int64
	if timeout > 0 {
		e = time.Now().Add(timeout).UnixNano()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	s.items[s.prefix+key] = Item{
		Object:     val,
		Expiration: e,
	}

	if e > 0 {
		s.DeleteExpired()
	}

	return nil
}

// Exist check cache's existence in memory.
func (s *Store) Exist(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[s.prefix+key]

	if item.Expired() {
		return false
	}

	return ok
}

// Forget Remove an item from the cache.
func (s *Store) Forget(key string) error {
	delete(s.items, s.prefix+key)
	return nil
}

// Remove all items from the cache.
func (s *Store) Flush() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.items = map[string]Item{}

	return nil
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

// Delete all expired items from the cache.
func (s *Store) DeleteExpired() {
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
			if time.Duration(sub)*time.Nanosecond < smallestDuration {
				smallestDuration = time.Duration(sub) * time.Nanosecond
			}
		}
	}

	if smallestDuration > 0 {
		s.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go s.DeleteExpired()
		})
	}
}
