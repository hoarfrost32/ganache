package policies

// FIFO implements a first-in, first-out eviction policy.
type FIFO[K comparable] struct {
	keyQueue []K
}

// NewFIFO creates a new FIFO eviction policy.
func NewFIFO[K comparable]() EvictionPolicy[K] {
	return &FIFO[K]{
		keyQueue: make([]K, 0),
	}
}

// TrackAddition adds a key to the end of the queue.
func (f *FIFO[K]) TrackAddition(key K) {
	f.keyQueue = append(f.keyQueue, key)
}

// TrackGet does nothing in a FIFO policy, as access order doesn't affect eviction.
func (f *FIFO[K]) TrackGet(key K) {}

// Evict removes and returns the oldest key from the queue.
// If the queue is empty, it returns the zero value for type K.
func (f *FIFO[K]) Evict() K {
	if len(f.keyQueue) == 0 {
		var zero K
		return zero
	}
	evictedKey := f.keyQueue[0]
	f.keyQueue = f.keyQueue[1:]
	return evictedKey
}
