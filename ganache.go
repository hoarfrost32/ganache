// Ganache is a general-purpose, simple and extensible
// in-memory caching library in go.
package ganache

import (
	"github.com/hoarfrost32/ganache/policies"
	"time"
)

// Cache is a generic, in-memory key-value store with a configurable
// eviction policy and size limit.
type Cache[K comparable, V any] struct {
	storage        map[K]V
	capacity       int
	backupInterval time.Duration
	policy         policies.EvictionPolicy[K]
}

// New creates a new Cache.
// The capacity must be a positive integer. If capacity is zero or negative,
// the cache will not store any new items.
func New[K comparable, V any](capacity int, backupInterval time.Duration, policy policies.EvictionPolicy[K]) *Cache[K, V] {
	return &Cache[K, V]{
		storage:        make(map[K]V),
		capacity:       capacity,
		backupInterval: backupInterval,
		policy:         policy,
	}
}

// Get retrieves a value from the cache for a given key.
// It returns the value and a boolean that is true if the key was found.
// Accessing a key with Get marks it as recently used by the eviction policy.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, found := c.storage[key]
	if found {
		c.policy.TrackGet(key)
	}
	return value, found
}

// Put adds a key-value pair to the cache.
// If the key already exists, its value is updated. If the cache is at capacity,
// an item is evicted based on the configured eviction policy to make room for
// the new item.
// If the cache was initialized with a capacity of 0 or less, Put is a no-op
// for new keys.
func (c *Cache[K, V]) Put(key K, value V) {
	if c.capacity <= 0 {
		return
	}

	if _, found := c.storage[key]; found {
		c.storage[key] = value
		c.policy.TrackGet(key)
		return
	}

	if len(c.storage) >= c.capacity {
		evictedKey := c.policy.Evict()
		delete(c.storage, evictedKey)
	}

	c.storage[key] = value
	c.policy.TrackAddition(key)
}
