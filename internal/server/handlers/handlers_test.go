package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGaugeMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		url         string
		metricName  string
		metricValue string
		want        want
		method      string
	}{
		{
			name:        "statusOkGauge",
			url:         "/update/gauge/{metricName}/{metricValue}",
			metricName:  "someMetric",
			metricValue: "13.0",
			method:      http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "statusIncorrectMetricValue",
			url:         "/update/gauge/{metricName}/{metricValue}",
			metricName:  "someMetric",
			metricValue: "string",
			method:      http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.metricName)
			rctx.URLParams.Add("metricValue", tt.metricValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			AddGaugeMetric(w, r)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}

func TestAddCounterMetricHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name        string
		url         string
		metricName  string
		metricValue string
		want        want
		method      string
	}{
		{
			name:        "statusOkCounter",
			url:         "/update/counter/{metricName}/{metricValue}",
			metricName:  "someMetric",
			metricValue: "13",
			method:      http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:        "statusIncorrectMetricValue",
			url:         "/update/counter/{metricName}/{metricValue}",
			metricName:  "someMetric",
			metricValue: "string",
			method:      http.MethodPost,
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("metricName", tt.metricName)
			rctx.URLParams.Add("metricValue", tt.metricValue)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			AddCounterMetric(w, r)
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
