package main

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/agent/collector"
	"github.com/agatma/sprint1-http-server/internal/agent/handlers"
	"github.com/agatma/sprint1-http-server/internal/agent/storage"
	"log"
	"time"
)

func main() {
	parseFlags()
	collectMetricsTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	sendMetricsTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	metricStorage := &storage.MetricsStorage{
		Metrics: make(map[string]float64),
	}
	host := fmt.Sprintf("http://localhost%s", flagRunAddr)
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
