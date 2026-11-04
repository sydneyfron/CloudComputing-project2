package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TempGauge tracks time-series temperature data per zipcode
	TempGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_temperature_fahrenheit",
		Help: "Forecasted temperature",
	}, []string{"zipcode"})

	// WindGauge tracks time-series wind speed data per zipcode
	WindGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "weather_wind_speed_mph",
		Help: "Forecasted wind speed",
	}, []string{"zipcode"})

	// AlertCounter tracks how many alerts are triggered per zipcode
	AlertCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "weather_alerts_total",
		Help: "Total number of weather alerts triggered",
	}, []string{"zipcode", "alert_type"})
)