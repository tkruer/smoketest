package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	numRequests = 100
	apiURL      = "https://api.example.com/endpoint"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(numRequests)

	latencies := make([]time.Duration, numRequests)
	client := &http.Client{}

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(apiURL)
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", i, err)
				return
			}
			resp.Body.Close()
			latencies[i] = time.Since(start)
		}(i)
	}

	wg.Wait()
	saveLatenciesToCSV(latencies)
	fmt.Println("Latency data saved to latencies.csv")
}

func saveLatenciesToCSV(latencies []time.Duration) {
	file, err := os.Create("latencies.csv")
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i, latency := range latencies {
		err := writer.Write([]string{fmt.Sprintf("Request %d", i+1), latency.String()})
		if err != nil {
			fmt.Printf("Failed to write record to CSV: %v\n", err)
			return
		}
	}
}
