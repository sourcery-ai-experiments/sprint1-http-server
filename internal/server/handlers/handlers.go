package handlers

import (
	"fmt"
	"github.com/agatma/sprint1-http-server/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// handleError is a helper function to handle HTTP errors.
func handleError(res http.ResponseWriter, errMsg string, statusCode int) {
	http.Error(res, errMsg, statusCode)
}

// checkMetricName checks if the metric name is provided and not empty.
func checkMetricName(metricName string) bool {
	return metricName != ""
}

func AddGaugeMetric(res http.ResponseWriter, req *http.Request) {
	metricName, metricValue := chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue")
	if !checkMetricName(metricName) {
		handleError(res, "empty metric name", http.StatusNotFound)
		return
	}
	if err := storage.Storage.AddGaugeValues(metricName, metricValue); err != nil {
		handleError(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
}

func AddCounterMetric(res http.ResponseWriter, req *http.Request) {
	metricName, metricValue := chi.URLParam(req, "metricName"), chi.URLParam(req, "metricValue")
	if !checkMetricName(metricName) {
		handleError(res, "empty metric name", http.StatusNotFound)
		return
	}
	if err := storage.Storage.AddCounterValues(metricName, metricValue); err != nil {
		handleError(res, "incorrect metric value", http.StatusBadRequest)
		return
	}
}

func GetMetric(res http.ResponseWriter, req *http.Request) {
	metricType, metricName := chi.URLParam(req, "metricType"), chi.URLParam(req, "metricName")
	var v interface{}
	var found bool
	fmt.Println(metricName)
	switch metricType {
	case "gauge":
		v, found = storage.Storage.GetGaugeValues(metricName)
	case "counter":
		v, found = storage.Storage.GetCounterValues(metricName)
	default:
		handleError(res, "incorrect metric type", http.StatusNotFound)
		return
	}
	if !checkMetricName(metricName) {
		handleError(res, "empty metric name", http.StatusNotFound)
		return
	}
	if !found {
		handleError(res, "metric is not found", http.StatusNotFound)
		return
	}
	res.Write([]byte(fmt.Sprintf("%s=%v", metricName, v)))
}

func GetAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	html := "<html><body><ul>"
	for key, value := range storage.Storage.GetAllGaugeValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	for key, value := range storage.Storage.GetAllCounterValues() {
		html += fmt.Sprintf("<li>%s: %v</li>", key, value)
	}
	html += "</ul></body></html>"
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(html))
}
