package main

import "fmt"

func main() {
	fmt.Println("--- Running Ganache Cache Examples ---")
	fmt.Println()

	fmt.Println("--- FIFO Cache Example ---")
	runFIFOExample()
	fmt.Println("--------------------------")
	fmt.Println()

	fmt.Println("--- LIFO Cache Example ---")
	runLIFOExample()
	fmt.Println("--------------------------")
	fmt.Println()

	fmt.Println("--- LRU Cache Example ---")
	runLRUExample()
	fmt.Println("-------------------------")
	fmt.Println()

	fmt.Println("--- All examples finished ---")
}