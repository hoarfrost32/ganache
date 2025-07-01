package main

import (
	"fmt"

	"github.com/hoarfrost32/ganache"
)

// runFIFOExample demonstrates the cache with a First-In, First-Out eviction policy.
func runFIFOExample() {
	cache, err := ganache.New[string, int](3, "fifo")
	if err != nil {
		fmt.Printf("Failed to create cache: %v\n", err)
		return
	}

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
