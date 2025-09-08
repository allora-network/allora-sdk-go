package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	prefixMu     sync.RWMutex
	metricPrefix string
)

// SetPrefix updates the prefix applied to newly created metric names. Call this
// before using any metrics helpers to ensure the prefix is applied
// consistently.
func SetPrefix(p string) {
	prefixMu.Lock()
	metricPrefix = p
	prefixMu.Unlock()
}

func metricName(name string) string {
	prefixMu.RLock()
	defer prefixMu.RUnlock()
	return metricPrefix + name
}

var (
	httpMetricsOnce    sync.Once
	requestDuration    *prometheus.HistogramVec
	responseStatus     *prometheus.CounterVec
	responseStatusCode *prometheus.CounterVec
)

func initHTTPMetrics() {
	httpMetricsOnce.Do(func() {
		requestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    metricName("http_request_duration_seconds"),
				Help:    "Duration of HTTP requests by host, path, and method.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"host", "path", "method"},
		)

		responseStatus = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName("http_response_status_total"),
				Help: "Count of HTTP response status texts by host and path.",
			},
			[]string{"host", "path", "method", "status"},
		)

		responseStatusCode = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName("http_response_status_code_total"),
				Help: "Count of HTTP response status codes by host and path.",
			},
			[]string{"host", "path", "method", "status_code"},
		)

		prometheus.MustRegister(requestDuration, responseStatus, responseStatusCode)
	})
}

// MetricsRoundTripper wraps a base RoundTripper and records Prometheus metrics.
type MetricsRoundTripper struct {
	base http.RoundTripper
}

func NewMetricsRoundTripper(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	initHTTPMetrics()
	return &MetricsRoundTripper{base: base}
}

func (m *MetricsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := m.base.RoundTrip(req)
	duration := time.Since(start).Seconds()

	host := req.URL.Host
	path := req.URL.Path
	method := req.Method

	requestDuration.WithLabelValues(host, path, method).Observe(duration)

	if err != nil {
		responseStatus.WithLabelValues(host, path, method, "error").Inc()
		responseStatusCode.WithLabelValues(host, path, method, "error").Inc()
		return nil, err
	}

	statusCode := strconv.Itoa(resp.StatusCode)
	responseStatus.WithLabelValues(host, path, method, http.StatusText(resp.StatusCode)).Inc()
	responseStatusCode.WithLabelValues(host, path, method, statusCode).Inc()

	return resp, nil
}
