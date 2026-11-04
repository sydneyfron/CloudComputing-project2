package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"weather-monitor/internal/concurrency"
	"weather-monitor/internal/evaluator"
	"weather-monitor/internal/kafka"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY environment variable is not set")
	}

	brokerAddress := "kafka:9092"
	kafkaTopic := "weather-alerts"

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Serving metrics at :8080/metrics")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go evaluator.Start(brokerAddress, kafkaTopic)

	log.Println("Waiting for Kafka to initialize...")
	time.Sleep(15 * time.Second)

	// --- NEW CSV PARSING LOGIC ---
	file, err := os.Open("testdata/concurrent_requests.csv")
	if err != nil {
		log.Fatalf("Failed to open testcase file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header row
	if _, err := reader.Read(); err != nil {
		log.Fatalf("Failed to read CSV header: %v", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV records: %v", err)
	}

	var requests []concurrency.WeatherRequest
	for _, record := range records {
		days, _ := strconv.Atoi(record[1])
		requests = append(requests, concurrency.WeatherRequest{
			Zipcode: record[0],
			Days:    days,
		})
	}
	// -----------------------------

	jobs := make(chan concurrency.WeatherRequest, len(requests))
	results := make(chan concurrency.WeatherResult, len(requests))
	var wg sync.WaitGroup

	// Spin up 5 concurrent workers to handle the larger file load
	for w := 1; w <= 5; w++ {
		wg.Add(1)
		go concurrency.Worker(w, apiKey, jobs, results, &wg)
	}

	for _, req := range requests {
		jobs <- req
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	kafka.PublishWeatherResults(results, brokerAddress, kafkaTopic)
	log.Println("Publisher finished. Keeping app alive for metrics scraping...")
	
	select {} 
}