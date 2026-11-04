package concurrency

import (
    "sync"
    "weather-monitor/internal/api"
)

// Worker processes incoming requests concurrently and sends outcomes to the results channel
func Worker(id int, apiKey string, jobs <-chan WeatherRequest, results chan<- WeatherResult, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        data, err := api.FetchWeather(j.Zipcode, j.Days, apiKey)
        results <- WeatherResult{
            Request: j,
            Data:    data,
            Error:   err,
        }
    }
}