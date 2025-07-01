package policies

// EvictionPolicy defines the interface for cache eviction strategies.
// It tracks item usage and determines which item to evict when the cache is full.
type EvictionPolicy[K comparable] interface {
	// TrackAddition is called when a new item is added to the cache.
	TrackAddition(key K)
	// TrackGet is called when an item is accessed from the cache.
	TrackGet(key K)
	// Evict determines and returns the key of the item to be evicted.
	Evict() K
}
