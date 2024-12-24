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
	counter := 0
	startTime := time.Now()
	for {
		_ = rand.Float64() * rand.Float64() // Perform some floating-point calculations
		counter++
		if counter%100000 == 0 {
			elapsed := time.Since(startTime).Seconds()
			if elapsed < 1.0/targetCPU {
				time.Sleep(time.Duration((1.0/targetCPU - elapsed) * float64(time.Second)))
			}
			startTime = time.Now()
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
			if len(chunks) > 0 {
				allocatedMemory -= len(chunks[0])
				chunks = chunks[1:]
			}
		} else {
			chunk := make([]byte, chunkSize)
			for i := range chunk {
				chunk[i] = byte(rand.Intn(256))
			}
			chunks = append(chunks, chunk)
			allocatedMemory += len(chunk)
		}
		// Periodically print memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("[Memory Worker %d] Memory usage: %.2f MB allocated\n", id, float64(m.Alloc)/(1024*1024))
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	runtime.GOMAXPROCS(6) // Limit to 6 vCPUs

	var wg sync.WaitGroup

	cpuWorkers := 5 // Use 6 workers for CPU
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
