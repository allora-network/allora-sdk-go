package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	rpcMetricsOnce       sync.Once
	rpcRequestDuration   *prometheus.HistogramVec
	rpcRequestTotal      *prometheus.CounterVec
	rpcAttemptsHistogram *prometheus.HistogramVec
	rpcRetryTotal        *prometheus.CounterVec
)

func initRPCMetrics() {
	rpcMetricsOnce.Do(func() {
		rpcRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    metricName("rpc_request_duration_seconds"),
				Help:    "Duration of RPC requests grouped by protocol, endpoint, service, method, and outcome.",
				Buckets: []float64{0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30},
			},
			[]string{"protocol", "endpoint", "service", "method", "outcome"},
		)

		rpcRequestTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName("rpc_requests_total"),
				Help: "Total number of RPC requests grouped by protocol, endpoint, service, method, and outcome.",
			},
			[]string{"protocol", "endpoint", "service", "method", "outcome"},
		)

		rpcAttemptsHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    metricName("rpc_request_attempts"),
				Help:    "Number of attempts performed per RPC call grouped by protocol, service, method, and outcome.",
				Buckets: []float64{1, 2, 3, 4, 5, 7, 10, 15},
			},
			[]string{"protocol", "service", "method", "outcome"},
		)

		rpcRetryTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName("rpc_request_retries_total"),
				Help: "Total retry attempts beyond the first try grouped by protocol, service, method, and outcome.",
			},
			[]string{"protocol", "service", "method", "outcome"},
		)

		prometheus.MustRegister(rpcRequestDuration, rpcRequestTotal, rpcAttemptsHistogram, rpcRetryTotal)
	})
}

// ObserveRPCRequest records latency and total counters for a single RPC attempt.
func ObserveRPCRequest(protocol, endpoint, service, method, outcome string, attempt int, duration time.Duration) {
	initRPCMetrics()
	labels := []string{protocol, endpoint, service, method, outcome}
	rpcRequestDuration.WithLabelValues(labels...).Observe(duration.Seconds())
	rpcRequestTotal.WithLabelValues(labels...).Inc()
	if attempt > 1 {
		rpcRetryTotal.WithLabelValues(protocol, service, method, outcome).Inc()
	}
}

// ObserveRPCAttempts records the number of attempts required for an RPC call.
func ObserveRPCAttempts(protocol, service, method, outcome string, attempts int) {
	initRPCMetrics()
	if attempts < 0 {
		attempts = 0
	}
	rpcAttemptsHistogram.WithLabelValues(protocol, service, method, outcome).Observe(float64(attempts))
}
