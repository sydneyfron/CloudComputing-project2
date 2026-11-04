package concurrency

import "weather-monitor/internal/api"

// WeatherRequest holds the specific parameters from different users
type WeatherRequest struct {
    Zipcode string
    Days    int
}

// WeatherResult bundles the original request with its retrieved API payload or error
type WeatherResult struct {
    Request WeatherRequest
    Data    *api.WeatherResponse
    Error   error
}