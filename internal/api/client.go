package api

import (
    "encoding/json"
    "fmt"
    "net/http"
)

const baseURL = "https://api.openweathermap.org/data/2.5/forecast"

func FetchWeather(zipcode string, days int, apiKey string) (*WeatherResponse, error) {
    // Note: OpenWeatherMap standard free tier uses 'cnt' for timestamps. 
    // Adjust calculation based on the specific API version you select.
    url := fmt.Sprintf("%s?zip=%s,US&cnt=%d&appid=%s&units=imperial", baseURL, zipcode, days, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
    }

    var weatherData WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
        return nil, err
    }

    return &weatherData, nil
}