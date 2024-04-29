package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

var log = logrus.New()

type MetricRepository interface {
	addGaugeValues(string, string) error
	addCounterValues(string, string) error
}

type MemStorage struct {
	gaugeValues   map[string]float64
	counterValues map[string]int64
}

var storage = &MemStorage{
	gaugeValues:   make(map[string]float64),
	counterValues: make(map[string]int64),
}

func (m *MemStorage) addGaugeValues(key string, value string) error {
	gaugeValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.gaugeValues[key] = gaugeValue
	return nil
}

func (m *MemStorage) addCounterValues(key string, value string) error {
	counterValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.counterValues[key] += int64(counterValue)
	return nil
}

func gaugeMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", req.Method)
		return
	}
	urlPath := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(urlPath) < 4 {
		http.Error(res, "incorrect input path", http.StatusBadRequest)
		return
	}
	metricName, metricValue := urlPath[2], urlPath[3]
	if metricName == "" {
		http.Error(res, "empty metric name", http.StatusNotFound)
		return
	}
	err := storage.addGaugeValues(metricName, metricValue)
	if err != nil {
		http.Error(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
	log.Info("Current storage state: ", storage)
}

func undefinedMetricType(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "incorrect metric type", http.StatusBadRequest)
}

func counterMetricHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method is not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", req.Method)
		return
	}
	urlPath := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(urlPath) < 4 {
		http.Error(res, "incorrect input path", http.StatusNotFound)
		return
	}
	metricName, metricValue := urlPath[2], urlPath[3]
	if metricName == "" {
		http.Error(res, "empty metric name", http.StatusNotFound)
		return
	}
	err := storage.addCounterValues(metricName, metricValue)
	if err != nil {
		http.Error(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
	log.Info("Current storage state: ", storage)
}

func init() {
	log.Formatter = &logrus.JSONFormatter{}
	log.Level = logrus.InfoLevel
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/gauge/`, gaugeMetricHandler)
	mux.HandleFunc(`/update/counter/`, counterMetricHandler)
	mux.HandleFunc(`/update/`, undefinedMetricType)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
