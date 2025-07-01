package policies

// LIFO implements a last-in, first-out eviction policy.
type LIFO[K comparable] struct {
	keyStack []K
}

// NewLIFO creates a new LIFO eviction policy.
func NewLIFO[K comparable]() EvictionPolicy[K] {
	return &LIFO[K]{
		keyStack: make([]K, 0),
	}
}

// TrackAddition adds a key to the top of the stack.
func (f *LIFO[K]) TrackAddition(key K) {
	f.keyStack = append(f.keyStack, key)
}

// TrackGet does nothing in a LIFO policy, as access order doesn't affect eviction.
func (f *LIFO[K]) TrackGet(key K) {}

// Evict removes and returns the most recently added key from the stack.
// If the stack is empty, it returns the zero value for type K.
func (f *LIFO[K]) Evict() K {
	if len(f.keyStack) == 0 {
		var zero K
		return zero
	}
	evictedKey := f.keyStack[len(f.keyStack)-1]
	f.keyStack = f.keyStack[:len(f.keyStack)-1]
	return evictedKey
}
