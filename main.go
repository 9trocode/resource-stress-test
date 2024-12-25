package main

import (
    "fmt"
    "os"
    "os/signal"
    "runtime"
    "sync"
    "syscall"
    "time"
)

const (
    // Maximum number of CPUs to use
    MaxCPUs = 5

    // Maximum memory usage in bytes (6 GB)
    MaxMemoryBytes = 6 * 1024 * 1024 * 1024 // 6 GB

    // Interval for logging (e.g., every 5 seconds)
    LogInterval = 5 * time.Second

    // Number of CPU-bound worker goroutines
    NumWorkers = MaxCPUs
)

// AppState holds the application's state, including allocated memory
type AppState struct {
    Memory [][]byte
}

func main() {
    // Set the maximum number of CPUs
    runtime.GOMAXPROCS(MaxCPUs)
    fmt.Printf("Set GOMAXPROCS to %d\n", MaxCPUs)

    // Initialize application state
    state := &AppState{}

    // Allocate memory up to the limit
    var err error
    state.Memory, err = allocateMemory(MaxMemoryBytes)
    if err != nil {
        fmt.Printf("Memory allocation failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Memory allocation completed.")

    // Start CPU-bound worker goroutines
    var wg sync.WaitGroup
    wg.Add(NumWorkers)
    for i := 0; i < NumWorkers; i++ {
        go cpuBoundWorker(i, &wg)
    }
    fmt.Printf("Started %d CPU-bound workers.\n", NumWorkers)

    // Set up channel to listen for interrupt or terminate signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start a ticker for logging
    ticker := time.NewTicker(LogInterval)
    defer ticker.Stop()

    // Main loop
    for {
        select {
        case <-ticker.C:
            // Access memory to mark it as used
            if len(state.Memory) > 0 {
                _ = state.Memory[0][0] // Example access to prevent optimization
            }

            // Get memory stats
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            allocMB := float64(m.Alloc) / 1024 / 1024

            // Get number of CPUs
            numCPU := runtime.GOMAXPROCS(0)

            // Get number of goroutines
            numGoroutines := runtime.NumGoroutine()

            // Log the information
            fmt.Printf("Memory Usage: %.2f MB | vCPUs: %d | Goroutines: %d\n", allocMB, numCPU, numGoroutines)
        case sig := <-sigChan:
            fmt.Printf("\nReceived signal: %s. Shutting down gracefully...\n", sig)
            // Perform any necessary cleanup here
            os.Exit(0)
        }
    }
}

// allocateMemory allocates memory up to the specified limit.
// It returns a slice that holds the allocated memory to prevent garbage collection.
func allocateMemory(limit int64) ([][]byte, error) {
    var (
        mu       sync.Mutex
        memory   [][]byte
        alloced int64
    )

    chunkSize := int64(100 * 1024 * 1024) // 100 MB per chunk
    for alloced+chunkSize <= limit {
        chunk := make([]byte, chunkSize)
        // Initialize the chunk to ensure memory is committed
        for i := range chunk {
            chunk[i] = 1
        }

        mu.Lock()
        memory = append(memory, chunk)
        alloced += chunkSize
        mu.Unlock()

        fmt.Printf("Allocated %d MB / %d MB\n", alloced/(1024*1024), limit/(1024*1024))
        time.Sleep(100 * time.Millisecond) // Slight delay between allocations
    }

    // Final allocation to reach exactly the limit
    remaining := limit - alloced
    if remaining > 0 {
        chunk := make([]byte, remaining)
        for i := range chunk {
            chunk[i] = 1
        }

        mu.Lock()
        memory = append(memory, chunk)
        alloced += remaining
        mu.Unlock()

        fmt.Printf("Allocated %d MB / %d MB\n", alloced/(1024*1024), limit/(1024*1024))
    }

    return memory, nil
}

// cpuBoundWorker performs a CPU-intensive task continuously.
// For demonstration, it calculates prime numbers in an infinite loop.
func cpuBoundWorker(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Worker %d started.\n", id)
    num := 2
    for {
        if isPrime(num) {
            // Found a prime number; do something if needed
            // For now, we'll just ignore it to keep the loop busy
        }
        num++
        // To prevent integer overflow
        if num < 0 {
            num = 2
        }
    }
}

// isPrime checks if a number is prime.
// This is a naive implementation for demonstration purposes.
func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}
