package policies

// EvictionPolicyType defines the type for eviction policy identifiers.
type EvictionPolicyType int

const (
	// FIFOPolicy represents the First-In, First-Out eviction policy.
	FIFOPolicy EvictionPolicyType = iota
	// LIFOPolicy represents the Last-In, First-Out eviction policy.
	LIFOPolicy
	// LRUPolicy represents the Least Recently Used eviction policy.
	LRUPolicy
)
