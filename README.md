# Ganache: An In-Memory Caching Library

Ganache is a lightweight, thread-safe, and extensible in-memory caching library for Go. It is designed for general-purpose use and provides multiple standard eviction policies out of the box, while also allowing for the implementation of custom policies.

## Features

- **Multiple Eviction Policies**: Supports `LRU` (Least Recently Used), `FIFO` (First-In, First-Out), and `LIFO` (Last-In, First-Out) policies.
- **Extensible by Design**: Easy to add your own custom eviction policies by implementing a simple interface.
- **Thread-Safe**: All cache operations (`Get` and `Put`) are safe for concurrent use.
- **Generics**: Fully utilizes Go generics for compile-time type safety.
- **Simple API**: A clean and minimal API for straightforward integration into any project.

## Getting Started

To use Ganache in your Go project, you can add it using `go get`:

```sh
go get github.com/hoarfrost32/ganache
```

## Usage

Creating and using a new cache is simple. You need to specify the cache's capacity and the desired eviction policy.

```go
package main

import (
	"fmt"
	"github.com/hoarfrost32/ganache"
	"github.com/hoarfrost32/ganache/policies"
)

func main() {
	// Create a new cache with a capacity of 3 and an LRU eviction policy
	cache, err := ganache.New[string, int](3, "lru")
	if err != nil {
		fmt.Printf("Failed to create cache: %v\n", err)
		return
	}

	// Add items to the cache
	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	// Retrieve an item. This marks 'a' as the most recently used.
	if val, ok := cache.Get("a"); ok {
		fmt.Printf("Got 'a': %d\n", val)
	}

	// Add another item. Because the cache is full, this will evict the
	// least recently used item ('b' in this case).
	cache.Put("d", 4)

	if _, ok := cache.Get("b"); !ok {
		fmt.Println("Cache miss for key 'b' (as expected, it was evicted).")
	}
}
```

## Eviction Policies

Ganache comes with three built-in eviction policies, which are specified by a string identifier when creating a cache:

- `"lru"`: **Least Recently Used**. Evicts the item that has not been accessed for the longest time.
- `"fifo"`: **First-In, First-Out**. Evicts the item that was added first, regardless of access patterns.
- `"lifo"`: **Last-In, First-Out**. Evicts the item that was added most recently.

## Running Examples & Benchmarks

The `/examples` directory contains demonstrations for each policy and a benchmark to showcase the performance benefits of caching.

To run all examples from the root of the directory:

```sh
go mod tidy
go run ./examples
```

This will run through three basic sanity check examples on `LIFO`, `FIFO`, and `LRU`, and then execute a benchmark that compares the performance of the `ganache` cache against a no-cache baseline where every read hits a simulated slow file store (with LRU policy).

## Extensibility: Adding a Custom Eviction Policy

Ganache is designed to be extensible, allowing you to add your own custom eviction policies without modifying the library's source code. This is achieved through a registration system.

The `EvictionPolicy` interface is defined as:

```Enterpret/policies/policy.go#L5-12
type EvictionPolicy[K comparable] interface {
	// TrackAddition is called when a new item is added to the cache.
	TrackAddition(key K)
	// TrackGet is called when an item is accessed from the cache.
	TrackGet(key K)
	// Evict determines and returns the key of the item to be evicted.
	Evict() K
}
```

To add and use a custom policy, follow these steps:

1.  **Implement the Policy**: Create a new struct and its constructor, ensuring it implements the `EvictionPolicy` interface.

2.  **Register Your Policy**: Use the `policies.Register` function to add your policy to the global registry. This is typically done in an `init()` function in your own code.

    ```go
    package main

    import "github.com/hoarfrost32/ganache/policies"

    // Define your custom policy
    type CustomPolicy[K comparable] struct{ /* ... internal state ... */ }
    func NewCustomPolicy[K comparable]() policies.EvictionPolicy[K] {
        return &CustomPolicy[K]{}
    }
    
    // ... implement interface methods for MRUPolicy ...

    // register it with a unique name in an init function
    func init() {
        policies.Register[string]("custom", NewCustomPolicy[string])
    }
    ```

3.  **Use It**: Once registered, you can create a cache using the string name you provided.

    ```go
    // Now you can create a cache with your custom policy
    mruCache, err := ganache.New[string, int](50, "mru")
    if err != nil {
        // handle error
    }
    ```

This approach keeps the API simple for standard use cases while offering powerful, decoupled extensibility for custom requirements.

## Development Setup

This project uses [Nix Flakes](https://nixos.wiki/wiki/Flakes) to provide a consistent and reproducible development environment. To use it, install Nix with Flake support and run:

```sh
nix develop
```

This command will automatically provision a shell with Go, `gopls` (the Go language server), `delve` (the debugger), and `golangci-lint`. All Go environment variables like `GOPATH` are automatically configured to be local to the project directory.

## Evaluation Criteria Checklist

This implementation was designed to meet the following criteria:

-   **Quality of Code**:
    -   Clean, commented, and idiomatic Go.
    -   Use of modern Go generics for type-safety and reusability.
    -   Consistent formatting and linting best practices.
-   **Extensibility of Low Level Design**:
    -   Eviction policies are decoupled from the core cache logic via the `EvictionPolicy` interface.
    -   The `Cache` struct operates independently of any specific policy's implementation details.
-   **Must Haves**:
    -   [x] **Support for multiple Standard Eviction Policies (FIFO, LRU, LIFO)**: Implemented and demonstrated in examples.
    -   [x] **Support to add custom eviction policies**: Implemented via the `EvictionPolicy` interface.
-   **Good To Have**:
    -   [x] **Thread safety**: Implemented using `sync.Mutex` to protect shared state in the `Cache`.

## Considerations

There is a notable tradeoff between thread safety and concurrency made in the implementation of this library. In particular, Go's `sync.RWMutex` could not be put to use at all, since even the `Get` method has the side effect of `TrackGet` in its eviction policy. Working solely with a mutex lock means that multi-threaded applications looking to make use of this library will be forced to proceed sequentially when interfacing with the cache. 

One way to retain concurrency would be to introduce Sharding. The cache can internally be split up into shards of equal size, each with their own storage and `EvictionPolicy`. This would lead to substantial improvements in concurrency, largely controlled by the number of shards generated.