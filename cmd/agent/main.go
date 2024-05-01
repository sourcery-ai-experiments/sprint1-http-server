package main

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/agent/collector"
	"github.com/agatma/sprint1-http-server/internal/agent/handlers"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"log"
	"strings"
	"time"
)

func main() {
	parseFlags()
	collectMetricsTicker := time.NewTicker(options.pollInterval)
	sendMetricsTicker := time.NewTicker(options.reportInterval)
	metricStorage := &storage.MetricsStorage{
		Metrics: make(map[string]float64),
	}
	port := strings.Split(flagRunAddr, ":")[1]
	host := fmt.Sprintf("http://127.0.0.1:%s", port)
	var PollCount int64
	for {
		select {
		case <-collectMetricsTicker.C:
			metrics := collector.CollectMetrics()
			metricStorage.Metrics = metrics
			PollCount++
		case <-sendMetricsTicker.C:
			err := handlers.SendGaugeMetrics(host, metricStorage)
			if err != nil {
				log.Fatal(err)
				return
			}
			err = handlers.SendCounterMetrics(host, "PollCount", PollCount)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
}
