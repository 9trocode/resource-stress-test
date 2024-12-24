package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func consumeCPU(wg *sync.WaitGroup) {
	defer wg.Done()
	counter := 0
	startTime := time.Now()
	for {
		_ = rand.Float64() * rand.Float64() // Perform some floating-point calculations
		counter++
		if counter%1000000 == 0 {
			elapsed := time.Since(startTime).Seconds()
			fmt.Printf("CPU usage: %.2f vCPUs over %.2f seconds\n", float64(counter)/1000000.0, elapsed)
			startTime = time.Now()
		}
	}
}

func consumeMemory(wg *sync.WaitGroup, memoryChunkSize int) {
	defer wg.Done()
	var chunks [][]byte
	counter := 0
	for {
		chunk := make([]byte, memoryChunkSize) // Allocate a chunk of memory
		for i := range chunk {
			chunk[i] = byte(rand.Intn(256)) // Fill memory with random data
		}
		chunks = append(chunks, chunk)
		counter++
		// Periodically print memory usage
		if counter%100 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Memory usage: %.2f MB allocated (%d iterations)\n", float64(m.Alloc)/(1024*1024), counter)
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all available CPUs

	var wg sync.WaitGroup

	cpuWorkers := runtime.NumCPU() * 2 // Double the number of CPU cores
	memoryChunkSize := 10 * 1024 * 1024 // Allocate 10 MB per chunk

	fmt.Println("Starting CPU and memory consumption...")

	// Start CPU consumers
	for i := 0; i < cpuWorkers; i++ {
		wg.Add(1)
		go consumeCPU(&wg)
	}

	// Start memory consumers
	wg.Add(1)
	go consumeMemory(&wg, memoryChunkSize)

	wg.Wait() // This will block forever since the goroutines never exit
}
