package main

import (
	"fmt"

	"github.com/hoarfrost32/ganache"
)

// runLRUExample demonstrates the cache with a Least Recently Used eviction policy.
func runLRUExample() {
	cache, err := ganache.New[string, int](3, "lru")
	if err != nil {
		fmt.Printf("Failed to create cache: %v\n", err)
		return
	}

	fmt.Println("Putting 'a': 1, 'b': 2, 'c': 3")
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Access 'a' to mark it as the most recently used item.
	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d. 'a' is now the most recently used item.\n", val)
	}

	// Add 'd', which evicts 'b' (the least recently used item).
	fmt.Println("Putting 'd': 4. This should evict 'b'.")
	cache.Put("d", 4)

	if _, ok := cache.Get("b"); !ok {
		fmt.Println("Get 'b': Not found (as expected).")
	}

	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Get 'a': %d\n", val)
	}
	if val, ok := cache.Get("c"); ok {
		fmt.Printf("Get 'c': %d\n", val)
	}
	if val, ok := cache.Get("d"); ok {
		fmt.Printf("Get 'd': %d\n", val)
	}
}
