package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// SimulateCPUUsage creates a CPU-bound goroutine to simulate load.
func SimulateCPUUsage(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	for {
		select {
		case <-stopChan:
			return
		default:
			// Simulate CPU load by performing calculations.
			_ = rand.Float64() * rand.Float64()
		}
	}
}

// SimulateMemoryUsage allocates and maintains a fixed amount of memory.
func SimulateMemoryUsage(sizeInBytes int, stopChan <-chan struct{}) {
	memory := make([]byte, sizeInBytes)
	for i := range memory {
		memory[i] = byte(rand.Intn(256))
	}
	<-stopChan // Wait until the application is stopped.
}

// LogResourceUsage logs memory and CPU usage periodically.
func LogResourceUsage(interval time.Duration, stopChan <-chan struct{}) {
	var m runtime.MemStats
	for {
		select {
		case <-stopChan:
			return
		case <-time.After(interval):
			runtime.ReadMemStats(&m)
			cpuCount := runtime.NumCPU()
			goroutines := runtime.NumGoroutine()
			fmt.Printf("[Resource Usage] Alloc = %.2f MB, TotalAlloc = %.2f MB, Sys = %.2f MB, NumGC = %d, Goroutines = %d, CPUs = %d\n",
				float64(m.Alloc)/1024/1024,
				float64(m.TotalAlloc)/1024/1024,
				float64(m.Sys)/1024/1024,
				m.NumGC,
				goroutines,
				cpuCount)
		}
	}
}

func main() {
	const (
		cpuCount   = 5                // Number of CPU cores to simulate.
		ramLimit   = 5 * 1024 * 1024 * 1024 // Memory limit in bytes (6 GB).
		logInterval = 5 * time.Second  // Interval for logging resource usage.
	)

	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	// Simulate CPU usage.
	for i := 0; i < cpuCount; i++ {
		wg.Add(1)
		go SimulateCPUUsage(&wg, stopChan)
	}

	// Simulate memory usage.
	go SimulateMemoryUsage(ramLimit, stopChan)

	// Start logging resource usage.
	go LogResourceUsage(logInterval, stopChan)

	fmt.Println("Application is running. Press Ctrl+C to stop.")

	// Block until the program is terminated.
	stopSignal := make(chan struct{})
	<-stopSignal

	// Clean up resources.
	close(stopChan)
	wg.Wait()
	fmt.Println("Application stopped.")
}
