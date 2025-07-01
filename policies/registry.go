package policies

import "sync"

// PolicyConstructor defines the function signature for creating a new eviction policy.
type PolicyConstructor[K comparable] func() EvictionPolicy[K]

// customPolicies holds the registry for user-defined eviction policies.
var (
	customPolicies = make(map[string]any)
	lock           = sync.RWMutex{}
)

// Register adds a new custom eviction policy constructor to the global registry.
// The name must be unique. If a policy with the same name is already registered,
// this function will panic.
//
// Example of registering a custom policy:
//
//	func init() {
//	    policies.Register[string]("custom", NewCustomPolicy[string])
//	}
func Register[K comparable](name string, constructor PolicyConstructor[K]) {
	lock.Lock()
	defer lock.Unlock()

	if _, exists := customPolicies[name]; exists {
		panic("policy with name " + name + " is already registered")
	}
	customPolicies[name] = constructor
}

// lookup attempts to find and return a custom policy constructor from the registry.
func lookup[K comparable](name string) (PolicyConstructor[K], bool) {
	lock.RLock()
	defer lock.RUnlock()

	constructorAny, found := customPolicies[name]
	if !found {
		return nil, false
	}

	constructor, ok := constructorAny.(PolicyConstructor[K])
	return constructor, ok
}
