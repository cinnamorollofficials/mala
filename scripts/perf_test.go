package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL    = "http://localhost:3000"
	virtualKey = "sk-gh-YOUR_VIRTUAL_KEY_HERE" // Replace with a real key from your DB
	numRequests = 100
	concurrency = 10
)

type Stats struct {
	mu           sync.Mutex
	latencies    []time.Duration
	statusCodes  map[int]int
	successCount int
	failCount    int
}

func main() {
	fmt.Printf("Starting performance test: %d requests with concurrency %d\n", numRequests, concurrency)
	
	stats := &Stats{
		statusCodes: make(map[int]int),
	}

	work := make(chan int, numRequests)
	for i := 0; i < numRequests; i++ {
		work <- i
	}
	close(work)

	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range work {
				sendRequest(stats)
			}
		}()
	}
	wg.Wait()
	duration := time.Since(startTime)

	printStats(stats, duration)
}

func sendRequest(stats *Stats) {
	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "Hello, my email is test@example.com and phone is 081234567890. What is 2+2?",
			},
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+virtualKey)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)

	stats.mu.Lock()
	defer stats.mu.Unlock()

	stats.latencies = append(stats.latencies, latency)
	if err != nil {
		stats.failCount++
		return
	}
	defer resp.Body.Close()

	stats.statusCodes[resp.StatusCode]++
	if resp.StatusCode == http.StatusOK {
		stats.successCount++
	} else {
		stats.failCount++
	}
}

func printStats(stats *Stats, totalDuration time.Duration) {
	var totalLatency time.Duration
	minLatency := stats.latencies[0]
	maxLatency := stats.latencies[0]

	for _, l := range stats.latencies {
		totalLatency += l
		if l < minLatency {
			minLatency = l
		}
		if l > maxLatency {
			maxLatency = l
		}
	}

	avgLatency := totalLatency / time.Duration(len(stats.latencies))
	rps := float64(stats.successCount+stats.failCount) / totalDuration.Seconds()

	fmt.Println("\n--- Performance Test Results ---")
	fmt.Printf("Total Duration:  %v\n", totalDuration)
	fmt.Printf("Requests/sec:    %.2f\n", rps)
	fmt.Printf("Total Requests:  %d\n", stats.successCount+stats.failCount)
	fmt.Printf("Success:         %d\n", stats.successCount)
	fmt.Printf("Failures:        %d\n", stats.failCount)
	fmt.Println("\n--- Latency ---")
	fmt.Printf("Min:             %v\n", minLatency)
	fmt.Printf("Max:             %v\n", maxLatency)
	fmt.Printf("Avg:             %v\n", avgLatency)
	fmt.Println("\n--- Status Codes ---")
	for code, count := range stats.statusCodes {
		fmt.Printf("%d: %d\n", code, count)
	}
}
