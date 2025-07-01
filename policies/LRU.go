package policies

// LRU implements a least-recently-used eviction policy. It uses a map for O(1)
// lookups and a doubly linked list to track usage. The head of the list is the
// most recently used item, and the tail is the least recently used.

// Node represents an element in the doubly linked list used by LRU.
type node[K comparable] struct {
	key  K
	prev *node[K]
	next *node[K]
}

// newNode creates a new list node.
func newNode[K comparable](key K) *node[K] {
	return &node[K]{key: key}
}

// LRU implements the EvictionPolicy interface with a least-recently-used strategy.
type LRU[K comparable] struct {
	nodes map[K]*node[K]
	head  *node[K]
	tail  *node[K]
}

// NewLRU creates a new LRU eviction policy.
func NewLRU[K comparable]() EvictionPolicy[K] {
	head := &node[K]{}
	tail := &node[K]{}

	head.next = tail
	tail.prev = head

	return &LRU[K]{
		nodes: make(map[K]*node[K]),
		head:  head,
		tail:  tail,
	}
}

// ------------------- helper functions -------------------

// addToFront adds a node to the front of the list, marking it as most recently used.
func (f *LRU[K]) placeAtFront(newNode *node[K]) {
	// currently MRU
	prevFront := f.head.next

	// place the new node at front
	newNode.next = prevFront
	prevFront.prev = newNode

	// connect it to head
	newNode.prev = f.head
	f.head.next = newNode
}

// removeNode removes a node from the list.
func (f *LRU[K]) removeNode(node *node[K]) {
	// get prev and next nodes
	nodePrev := node.prev
	nodeNext := node.next

	// remove the node from current position
	nodePrev.next = nodeNext
	nodeNext.prev = nodePrev
}

// --------------------------------------------------------

// TrackAddition adds a new key to the LRU policy.
// It is marked as the most recently used.
// If the key already exists in the cache,
// it is moved to the front.
func (f *LRU[K]) TrackAddition(key K) {
	if _, ok := f.nodes[key]; ok {
		f.TrackGet(key)
		return
	}

	newNode := newNode(key)
	// place it at the front
	f.placeAtFront(newNode)
	f.nodes[key] = newNode
}

// TrackGet marks a key as most recently used.
// If the key does not exist, it functions as a no-op.
func (f *LRU[K]) TrackGet(key K) {
	if _, ok := f.nodes[key]; ok {
		// the accessed key
		getNode := f.nodes[key]
		// remove it from current position...
		f.removeNode(getNode)
		// ...and place it at the front
		f.placeAtFront(getNode)
	}
}

// Evict removes and returns the least recently used key.
// If the cache is empty, it returns the zero value for type K.
func (f *LRU[K]) Evict() K {
	// The node to evict is the one before the tail sentinel.
	nodeToEvict := f.tail.prev

	// If the list is empty (head's next is tail), return zero value.
	if nodeToEvict == f.head {
		var zero K
		return zero
	}

	f.removeNode(nodeToEvict)
	delete(f.nodes, nodeToEvict.key)

	return nodeToEvict.key
}
