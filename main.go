// main.go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("Starting resource-consuming application...")

	// Display the number of CPU cores available
	numCPU := runtime.NumCPU()
	fmt.Printf("Number of CPU cores available: %d\n", numCPU)

	// Allocate 4 GB of memory
	fmt.Println("Allocating 4 GB of memory...")
	var memoryUsage [][]byte
	const totalMemory = 4 * 1024 * 1024 * 1024 // 4 GB
	const chunkSize = 100 * 1024 * 1024       // 100 MB chunks

	for i := 0; i < totalMemory/chunkSize; i++ {
		chunk := make([]byte, chunkSize)
		// Initialize the chunk to ensure memory is actually allocated
		for j := range chunk {
			chunk[j] = byte(j % 256)
		}
		memoryUsage = append(memoryUsage, chunk)
		fmt.Printf("Allocated %d MB/%d MB\n", (i+1)*100, 4000)
	}

	fmt.Println("Memory allocation complete.")

	// Start CPU-intensive computations on up to 4 CPU cores
	fmt.Println("Starting CPU-intensive computations...")
	var wg sync.WaitGroup
	numGoroutines := 4 // Number of CPU-intensive goroutines

	// Start a goroutine to report memory usage periodically
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Memory Usage: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v\n",
				bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
			time.Sleep(1 * time.Minute)
		}
	}()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d started.\n", id+1)
			for {
				_ = fibonacci(35) // Adjust the input for desired CPU usage
			}
		}(i)
	}

	// Let the application run for a specified duration
	runDuration := 10 * time.Minute
	fmt.Printf("Application will run for %v...\n", runDuration)
	time.Sleep(runDuration)

	// Signal goroutines to stop (in this example, we simply exit the program)
	fmt.Println("Run duration complete. Exiting application.")
	// Note: In this simple example, goroutines are not gracefully stopped.
	// For a graceful shutdown, consider using context cancellation or other signaling mechanisms.
}

// fibonacci computes the nth Fibonacci number recursively.
// Note: This is intentionally inefficient to consume CPU resources.
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
