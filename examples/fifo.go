package main

import (
	"fmt"
	"time"

	"github.com/hoarfrost32/ganache"
	"github.com/hoarfrost32/ganache/policies"
)

// runFIFOExample demonstrates the cache with a First-In, First-Out eviction policy.
func runFIFOExample() {
	fifoPolicy := policies.NewFIFO[string]()
	cache := ganache.New[string, int](3, 0*time.Second, fifoPolicy)

	fmt.Println("Putting 'a': 1, 'b': 2, 'c': 3")
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d. Access order doesn't matter for FIFO.\n", val)
	}

	fmt.Println("Putting 'd': 4. This should evict 'a'.")
	cache.Put("d", 4)

	if _, ok := cache.Get("a"); !ok {
		fmt.Println("Get 'a': Not found (as expected).")
	}

	if val, ok := cache.Get("b"); ok {
		fmt.Printf("Get 'b': %d\n", val)
	}
	if val, ok := cache.Get("c"); ok {
		fmt.Printf("Get 'c': %d\n", val)
	}
	if val, ok := cache.Get("d"); ok {
		fmt.Printf("Get 'd': %d\n", val)
	}
}
