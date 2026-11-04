package api

type WeatherResponse struct {
    City struct {
        Name string `json:"name"`
    } `json:"city"`
    List []Forecast `json:"list"`
}

type Forecast struct {
    Dt   int64 `json:"dt"`
    Main struct {
        Temp     float64 `json:"temp"`
        TempMin  float64 `json:"temp_min"`
        TempMax  float64 `json:"temp_max"`
    } `json:"main"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
}