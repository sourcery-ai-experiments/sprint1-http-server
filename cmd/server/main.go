package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

var availableMetrics = map[string]struct{}{gauge: {}, counter: {}}

type MetricRepository interface {
	addGaugeValues(string, string) error
	addCounterValues(string, string) error
	getGaugeValues(string) (float64, bool)
	getCounterValues(string) (int64, bool)
}

type MemStorage struct {
	gaugeValues   map[string]float64
	counterValues map[string]int64
}

func (m *MemStorage) addGaugeValues(key string, value string) error {
	gaugeValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.gaugeValues[key] = gaugeValue
	return nil
}

func (m *MemStorage) getGaugeValues(key string) (float64, bool) {
	value, found := m.gaugeValues[key]
	return value, found
}

func (m *MemStorage) addCounterValues(key string, value string) error {
	counterValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.counterValues[key] += int64(counterValue)
	return nil
}

func (m *MemStorage) getCounterValues(key string) (int64, bool) {
	value, found := m.counterValues[key]
	return value, found
}

func gaugeMetricHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}
	urlPath := strings.Split(req.URL.Path, "/")[2:]
	metricType := urlPath[0]
	if _, found := availableMetrics[metricType]; !found {
		http.Error(res, "incorrect metric type", http.StatusBadRequest)
		return
	}
	metricName := urlPath[1]
	if metricName == "" {
		http.Error(res, "incorrect metric name", http.StatusNotFound)
		return
	}
	metricValue := urlPath[2]
	storage := &MemStorage{
		gaugeValues:   make(map[string]float64),
		counterValues: make(map[string]int64),
	}
	var err error
	if metricType == gauge {
		err = storage.addGaugeValues(metricName, metricValue)
	} else {
		err = storage.addCounterValues(metricName, metricValue)
	}
	if err != nil {
		http.Error(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
	fmt.Println(storage)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, gaugeMetricHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
