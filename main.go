package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func consumeCPU(wg *sync.WaitGroup, id int, targetCPU float64) {
	defer wg.Done()
	for {
		startTime := time.Now()
		for i := 0; i < int(1e6); i++ {
			_ = rand.Float64() * rand.Float64() // Perform calculations
		}
		elapsed := time.Since(startTime).Seconds()
		if elapsed < 1.0/targetCPU {
			time.Sleep(time.Duration((1.0/targetCPU - elapsed) * float64(time.Second)))
		}
	}
}

func consumeMemory(wg *sync.WaitGroup, memoryLimitMB int, id int) {
	defer wg.Done()
	var chunks [][]byte
	chunkSize := 10 * 1024 * 1024 // 10 MB per chunk
	allocatedMemory := 0
	for {
		if allocatedMemory+chunkSize > memoryLimitMB*1024*1024 {
			// Maintain memory limit by keeping it constant
			time.Sleep(100 * time.Millisecond)
			continue
		}
		chunk := make([]byte, chunkSize)
		for i := range chunk {
			chunk[i] = byte(rand.Intn(256))
		}
		chunks = append(chunks, chunk)
		allocatedMemory += len(chunk)
		// Periodically print memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("[Memory Worker %d] Memory usage: %.2f MB allocated (target: %d MB)\n", id, float64(m.Alloc)/(1024*1024), memoryLimitMB)
	}
}

func main() {
	runtime.GOMAXPROCS(6) // Limit to 6 vCPUs

	var wg sync.WaitGroup

	cpuWorkers := 5       // Use 6 workers for CPU
	memoryLimitMB := 5000 // Limit memory usage to 6 GB

	fmt.Println("Starting CPU and memory consumption...")

	// Start CPU consumers
	for i := 0; i < cpuWorkers; i++ {
		wg.Add(1)
		go consumeCPU(&wg, i, 1.0) // Each worker uses 1 vCPU
	}

	// Start memory consumers
	for i := 0; i < 1; i++ { // Use a single memory worker
		wg.Add(1)
		go consumeMemory(&wg, memoryLimitMB, i)
	}

	wg.Wait() // This will block forever since the goroutines never exit
}
