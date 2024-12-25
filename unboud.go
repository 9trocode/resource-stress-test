// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"runtime"
// 	"sync"
// 	"time"
// )

// func consumeCPU(wg *sync.WaitGroup, id int) {
// 	defer wg.Done()
// 	counter := 0
// 	startTime := time.Now()
// 	for {
// 		_ = rand.Float64() * rand.Float64() // Perform some floating-point calculations
// 		counter++
// 		if counter%2000000 == 0 { // Adjusted to reduce CPU intensity
// 			elapsed := time.Since(startTime).Seconds()
// 			fmt.Printf("[CPU Worker %d] CPU usage: %.2f vCPUs over %.2f seconds\n", id, float64(counter)/2000000.0, elapsed)
// 			startTime = time.Now()
// 		}
// 	}
// }

// func consumeMemory(wg *sync.WaitGroup, memoryChunkSize int, id int) {
// 	defer wg.Done()
// 	var chunks [][]byte
// 	counter := 0
// 	for {
// 		chunk := make([]byte, memoryChunkSize) // Allocate a chunk of memory
// 		for i := range chunk {
// 			chunk[i] = byte(rand.Intn(256)) // Fill memory with random data
// 		}
// 		chunks = append(chunks, chunk)
// 		counter++
// 		// Periodically print memory usage
// 		if counter%50 == 0 { // Increased memory logging frequency
// 			var m runtime.MemStats
// 			runtime.ReadMemStats(&m)
// 			fmt.Printf("[Memory Worker %d] Memory usage: %.2f MB allocated (%d iterations)\n", id, float64(m.Alloc)/(1024*1024), counter)
// 		}
// 	}
// }

// func main() {
// 	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all available CPUs

// 	var wg sync.WaitGroup

// 	cpuWorkers := runtime.NumCPU() // Match CPU workers to number of cores for balance
// 	memoryWorkers := 4             // Increase number of memory workers to prioritize memory usage
// 	memoryChunkSize := 20 * 1024 * 1024 // Allocate 20 MB per chunk for more memory usage

// 	fmt.Println("Starting CPU and memory consumption...")

// 	// Start CPU consumers
// 	for i := 0; i < cpuWorkers; i++ {
// 		wg.Add(1)
// 		go consumeCPU(&wg, i)
// 	}

// 	// Start memory consumers
// 	for i := 0; i < memoryWorkers; i++ {
// 		wg.Add(1)
// 		go consumeMemory(&wg, memoryChunkSize, i)
// 	}

// 	wg.Wait() // This will block forever since the goroutines never exit
// }
