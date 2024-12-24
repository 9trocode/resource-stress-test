// main.go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Set GOMAXPROCS to 1 to limit the application to 1 CPU core
	runtime.GOMAXPROCS(1)
	fmt.Println("Starting resource-consuming application...")

	// Allocate 2 GB of memory
	fmt.Println("Allocating 2 GB of memory...")
	var memoryUsage [][]byte
	// 2 GB = 2 * 1024 * 1024 * 1024 bytes
	const totalMemory = 2 * 1024 * 1024 * 1024 // 2 GB
	const chunkSize = 100 * 1024 * 1024       // 100 MB chunks

	for i := 0; i < totalMemory/chunkSize; i++ {
		chunk := make([]byte, chunkSize)
		// Initialize the chunk to ensure memory is actually allocated
		for j := range chunk {
			chunk[j] = byte(j % 256)
		}
		memoryUsage = append(memoryUsage, chunk)
		fmt.Printf("Allocated %d MB/%d MB\n", (i+1)*100, 2000)
	}

	fmt.Println("Memory allocation complete.")

	// Start CPU-intensive computation
	fmt.Println("Starting CPU-intensive computation...")
	done := make(chan bool)

	go func() {
		// Perform CPU-intensive calculations indefinitely
		for {
			_ = fibonacci(35) // Adjust the input for desired CPU usage
		}
	}()

	// Let the application run for a specified duration
	runDuration := 10 * time.Minute
	fmt.Printf("Application will run for %v...\n", runDuration)
	time.Sleep(runDuration)

	// Signal the CPU goroutine to stop (in this case, we exit the program)
	fmt.Println("Run duration complete. Exiting application.")
	done <- true
}

// fibonacci computes the nth Fibonacci number recursively.
// Note: This is intentionally inefficient to consume CPU resources.
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
