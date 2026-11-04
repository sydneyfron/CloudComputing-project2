package evaluator

import (
	"context"
	"encoding/json"
	"log"
	"weather-monitor/internal/concurrency"
	"weather-monitor/internal/metrics"

	"github.com/segmentio/kafka-go"
)

const (
	HighTempThreshold = 85.0
	LowTempThreshold  = 32.0
	HighWindThreshold = 20.0
)

// Start Evaluator listens to the Kafka channel and triggers alerts[cite: 1]
func Start(brokerAddress string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  "weather-evaluator-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("Evaluator started, listening for Kafka messages...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading from Kafka: %v", err)
			continue
		}

		var res concurrency.WeatherResult
		if err := json.Unmarshal(m.Value, &res); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			continue
		}

		evaluateThresholds(res)
	}
}

func evaluateThresholds(res concurrency.WeatherResult) {
	zip := res.Request.Zipcode

	// For demonstration, we evaluate the first forecast period in the list[cite: 1]
	if len(res.Data.List) > 0 {
		forecast := res.Data.List[0]
		temp := forecast.Main.Temp
		wind := forecast.Wind.Speed

		// Update Prometheus Time-Series Metrics[cite: 1]
		metrics.TempGauge.WithLabelValues(zip).Set(temp)
		metrics.WindGauge.WithLabelValues(zip).Set(wind)

		// Trigger Alerts[cite: 1]
		if temp > HighTempThreshold {
			log.Printf("ALERT: High temperature (%.2fF) for %s!", temp, zip)
			metrics.AlertCounter.WithLabelValues(zip, "high_temp").Inc()
		}
		if temp < LowTempThreshold {
			log.Printf("ALERT: Low temperature (%.2fF) for %s!", temp, zip)
			metrics.AlertCounter.WithLabelValues(zip, "low_temp").Inc()
		}
		if wind > HighWindThreshold {
			log.Printf("ALERT: High wind (%.2fmph) for %s!", wind, zip)
			metrics.AlertCounter.WithLabelValues(zip, "high_wind").Inc()
		}
	}
}