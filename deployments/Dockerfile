FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /weather-app cmd/weather-monitor/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /weather-app /weather-app
# Bring the testdata folder into the final scratch container
COPY --from=builder /app/testdata /testdata
EXPOSE 8080
ENTRYPOINT ["/weather-app"]