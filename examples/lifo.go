package main

import (
	"fmt"

	"github.com/hoarfrost32/ganache"
)

// runLIFOExample demonstrates the cache with a Last-In, First-Out eviction policy.
func runLIFOExample() {
	cache, err := ganache.New[string, int](3, "lifo")
	if err != nil {
		fmt.Printf("Failed to create cache: %v\n", err)
		return
	}

	fmt.Println("Putting 'a': 1, 'b': 2, 'c': 3")
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d. Access order doesn't matter for LIFO.\n", val)
	}

	fmt.Println("Putting 'd': 4. This should evict 'c'.")
	cache.Put("d", 4)

	if _, ok := cache.Get("c"); !ok {
		fmt.Println("Get 'c': Not found (as expected).")
	}

	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d\n", val)
	}
	if val, ok := cache.Get("b"); ok {
		fmt.Printf("Get 'b': %d\n", val)
	}
	if val, ok := cache.Get("d"); ok {
		fmt.Printf("Get 'd': %d\n", val)
	}
}
