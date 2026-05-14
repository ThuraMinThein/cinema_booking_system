package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Endpoint struct {
	Name   string
	Method string
	URL    string
	BODY   string
}

var endpoints = []Endpoint{
	{
		Name:   "Health Check",
		Method: "GET",
		URL:    "http://localhost:3333/health",
	},
	{
		Name:   "Seats API",
		Method: "GET",
		URL:    "http://localhost:3333/seats",
		BODY:   `{"movie_id": 10}`,
	},
	// {
	// 	Name:   "booking API",
	// 	Method: "GET",
	// 	URL:    "http://localhost:3333/bookings",
	// },
}

const (
	totalRequests = 1000
	concurrency   = 50
	timeout       = 10 * time.Second
)

func main() {
	client := &http.Client{
		Timeout: timeout,
	}

	for _, endpoint := range endpoints {
		fmt.Printf("\n=============================\n")
		fmt.Printf("Testing: %s\n", endpoint.Name)
		fmt.Printf("Endpoint: %s %s\n", endpoint.Method, endpoint.URL)
		fmt.Printf("=============================\n")

		runLoadTest(client, endpoint)
	}
}

func runLoadTest(client *http.Client, endpoint Endpoint) {
	var successCount int64
	var failureCount int64

	var totalLatency int64

	start := time.Now()

	wg := sync.WaitGroup{}
	sem := make(chan struct{}, concurrency)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func(requestID int) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			reqStart := time.Now()

			req, err := http.NewRequest(endpoint.Method, endpoint.URL, strings.NewReader(endpoint.BODY))
			if err != nil {
				fmt.Printf("[ERROR] Request creation failed: %v\n", err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			resp, err := client.Do(req)
			latency := time.Since(reqStart).Milliseconds()

			atomic.AddInt64(&totalLatency, latency)

			if err != nil {
				fmt.Printf("[ERROR] Request failed: %v\n", err)
				atomic.AddInt64(&failureCount, 1)
				return
			}

			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failureCount, 1)
				fmt.Printf("[FAIL] Status Code: %d\n", resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()

	totalDuration := time.Since(start)

	avgLatency := float64(totalLatency) / float64(totalRequests)

	fmt.Printf("\nResults:\n")
	fmt.Printf("Total Requests : %d\n", totalRequests)
	fmt.Printf("Concurrency    : %d\n", concurrency)
	fmt.Printf("Success        : %d\n", successCount)
	fmt.Printf("Failures       : %d\n", failureCount)
	fmt.Printf("Total Time     : %s\n", totalDuration)
	fmt.Printf("Average Latency: %.2f ms\n", avgLatency)
	fmt.Printf("Requests/sec   : %.2f\n",
		float64(totalRequests)/totalDuration.Seconds())
}
