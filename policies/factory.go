package policies

import "fmt"

// CreatePolicy instantiates and returns an eviction policy based on its registered name,
// with a fall back to built-in policies if the name is not found. If the name does not 
// correspond to any policy at all, it returns an error.
func CreatePolicy[K comparable](name string) (EvictionPolicy[K], error) {
	// check if a custom policy with this name has been registered.
	if constructor, ok := lookup[K](name); ok {
		return constructor(), nil
	}

	// If not, fall back to built-in policies.
	switch name {
	case "fifo":
		return NewFIFO[K](), nil
	case "lifo":
		return NewLIFO[K](), nil
	case "lru":
		return NewLRU[K](), nil
	default:
		return nil, fmt.Errorf("unknown eviction policy: %s", name)
	}
}
