package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hoarfrost32/ganache"
)

const (
	benchmarkOperations     = 5000
	benchmarkUniqueKeys     = 1000
	benchmarkCacheCap       = 200
	benchmarkHotKeyRatio    = 0.2 // 20% of keys are "hot"
	benchmarkHotAccessRatio = 0.8 // 80% of operations target "hot" keys
	benchmarkDataFile       = "data.csv"
)

// fetchFromFileStore simulates reading from a slow data source (a file on disk),
// representing the "cache miss penalty."
func fetchFromFileStore(key int, dataFilePath string) (int, bool) {
	file, err := os.Open(dataFilePath)
	if err != nil {
		fmt.Printf("Failed to open data file: %v\n", err)
		return 0, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			k, _ := strconv.Atoi(parts[0])
			if k == key {
				v, _ := strconv.Atoi(parts[1])
				return v, true
			}
		}
	}
	return 0, false
}

// runCacheBenchmark runs the benchmark with the ganache cache.
func runCacheBenchmark(dataFilePath string) {
	numHotKeys := int(float64(benchmarkUniqueKeys) * benchmarkHotKeyRatio)

	fmt.Println("--- Running Ganache Cache Benchmark (with File Store) ---")
	fmt.Printf("Configuration: Capacity=%d, Operations=%d, UniqueKeys=%d\n",
		benchmarkCacheCap, benchmarkOperations, benchmarkUniqueKeys)
	fmt.Printf("Workload: %.0f%% of accesses on %.0f%% of keys\n", benchmarkHotAccessRatio*100, benchmarkHotKeyRatio*100)

	cache, err := ganache.New[int, int](benchmarkCacheCap, "lru")
	if err != nil {
		fmt.Printf("Failed to create cache: %v\n", err)
		return
	}
	hits, misses := 0, 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startTime := time.Now()

	for i := 0; i < benchmarkOperations; i++ {
		var key int
		if r.Float64() < benchmarkHotAccessRatio {
			key = r.Intn(numHotKeys)
		} else {
			key = numHotKeys + r.Intn(benchmarkUniqueKeys-numHotKeys)
		}

		// This is a "read-through" cache pattern.
		if _, ok := cache.Get(key); ok {
			hits++
		} else {
			misses++
			if value, found := fetchFromFileStore(key, dataFilePath); found {
				cache.Put(key, value)
			}
		}
	}

	duration := time.Since(startTime)
	totalGets := hits + misses
	var hitRate float64
	if totalGets > 0 {
		hitRate = float64(hits) / float64(totalGets) * 100
	}

	fmt.Println("\n--- Benchmark Results ---")
	fmt.Printf("Total Duration:     %s\n", duration)
	fmt.Printf("Cache Hits:         %d\n", hits)
	fmt.Printf("Cache Misses:       %d\n", misses)
	fmt.Printf("Hit Rate:           %.2f%%\n", hitRate)
	fmt.Println("---------------------------------------------------------")
}

// runNoCacheBenchmark runs the benchmark without any caching, serving as a baseline.
func runNoCacheBenchmark(dataFilePath string) {
	numHotKeys := int(float64(benchmarkUniqueKeys) * benchmarkHotKeyRatio)

	fmt.Println("--- Running No-Cache Baseline Benchmark (with File Store) ---")
	fmt.Printf("Configuration: Operations=%d, UniqueKeys=%d\n",
		benchmarkOperations, benchmarkUniqueKeys)
	fmt.Printf("Workload: %.0f%% of accesses on %.0f%% of keys\n", benchmarkHotAccessRatio*100, benchmarkHotKeyRatio*100)

	gets := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startTime := time.Now()

	for i := 0; i < benchmarkOperations; i++ {
		var key int
		if r.Float64() < benchmarkHotAccessRatio {
			key = r.Intn(numHotKeys)
		} else {
			key = numHotKeys + r.Intn(benchmarkUniqueKeys-numHotKeys)
		}

		fetchFromFileStore(key, dataFilePath)
		gets++
	}

	duration := time.Since(startTime)

	fmt.Println("\n--- Benchmark Results ---")
	fmt.Printf("Total Duration:     %s\n", duration)
	fmt.Printf("Total Gets:         %d (all from file store)\n", gets)
	fmt.Println("---------------------------------------------------------")
}

// runBenchmarks executes both the cache and no-cache benchmarks for comparison.
func runBenchmarks() {
	dataFilePath := filepath.Join("examples", benchmarkDataFile)

	runCacheBenchmark(dataFilePath)
	fmt.Println()

	runNoCacheBenchmark(dataFilePath)
}
