package cache

import (
	"time"

	"github.com/pkg/errors"
)

type Repository struct {
	store Store
}

func NewRepository(store Store) *Repository {
	s := &Repository{
		store: store,
	}
	return s
}

// Has Determine if an item exists in the cache.
func (r *Repository) Has(key string) bool {
	return r.store.Exist(key)
}

// Get Retrieve an item from the cache by key.
func (r *Repository) Get(key string, val interface{}) error {
	err := r.store.Get(key, val)
	return err
}

// Pull Retrieve an item from the cache and delete it.
func (r *Repository) Pull(key string, val interface{}) error {
	err := r.store.Get(key, val)
	if err != nil {
		return err
	}
	r.store.Forget(key)
	return nil
}

// Put Store an item in the cache.
func (r *Repository) Put(key string, val interface{}, timeout time.Duration) error {
	return r.store.Put(key, val, timeout)
}

// Set Store an item in the cache.
func (r *Repository) Set(key string, val interface{}, timeout time.Duration) error {
	return r.Put(key, val, timeout)
}

// Add Store an item in the cache if the key does not exist.
func (r *Repository) Add(key string, val interface{}, timeout time.Duration) error {
	if r.store.Exist(key) {
		return errors.New("the key already exists：" + key)
	}
	return r.store.Put(key, val, timeout)
}

// Remember Get an item from the cache, or store the default value.
func (r *Repository) Remember(key string, val interface{}, timeout time.Duration, callback func() interface{}) error {
	err := r.Get(key, val)
	if err == nil {
		return nil
	}

	value := callback()
	r.Put(key, value, timeout)

	return r.Get(key, val)
}

// Forget Remove an item from the cache.
func (r *Repository) Forget(key string) error {
	return r.store.Forget(key)
}

// Delete Alias for the "Delete" method.
func (r *Repository) Delete(key string) error {
	return r.Forget(key)
}

// Clear Remove all items from the cache.
func (r *Repository) Clear() error {
	return r.store.Flush()
}

// GetStore Get the cache store implementation.
func (r *Repository) GetStore() Store {
	return r.store
}
