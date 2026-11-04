package kafka

import (
    "context"
    "encoding/json"
    "log"
    "weather-monitor/internal/concurrency"

    "github.com/segmentio/kafka-go"
)

func PublishWeatherResults(results <-chan concurrency.WeatherResult, brokerAddress string, topic string) {
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{brokerAddress},
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    })
    defer writer.Close()

    for res := range results {
        if res.Error != nil {
            log.Printf("Skipping failed request for %s: %v", res.Request.Zipcode, res.Error)
            continue
        }

        payload, err := json.Marshal(res)
        if err != nil {
            log.Printf("Failed to marshal data for %s: %v", res.Request.Zipcode, err)
            continue
        }

        err = writer.WriteMessages(context.Background(),
            kafka.Message{
                Key:   []byte(res.Request.Zipcode),
                Value: payload,
            },
        )
        if err != nil {
            log.Printf("Failed to write to Kafka for %s: %v", res.Request.Zipcode, err)
        } else {
            log.Printf("Successfully published %s weather data to Kafka topic: %s", res.Request.Zipcode, topic)
        }
    }
}