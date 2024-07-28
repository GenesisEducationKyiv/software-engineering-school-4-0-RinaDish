package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	httpRequestsTotal *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	dbQueryDuration *prometheus.HistogramVec
	errorsTotal *prometheus.CounterVec
	Registry *prometheus.Registry
}

func NewMetrics() Metrics {
	m := Metrics{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "handler"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Duration of database queries",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query"},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors",
			},
			[]string{"type"},
		),
		Registry: prometheus.NewRegistry(),
	}
	m.registerMetrics()
	return m
}

func (m Metrics) registerMetrics() {
	m.Registry.MustRegister(m.httpRequestsTotal)
	m.Registry.MustRegister(m.httpRequestDuration)
	m.Registry.MustRegister(m.dbQueryDuration)
	m.Registry.MustRegister(m.errorsTotal)
}

func (m Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		timer := prometheus.NewTimer(m.httpRequestDuration.WithLabelValues(r.URL.Path, r.Method))
		defer timer.ObserveDuration()

		next.ServeHTTP(w, r)
	})
}

func (m Metrics) MonitorDBQuery(queryName string, queryFunc func() error) error {
	timer := prometheus.NewTimer(m.dbQueryDuration.WithLabelValues(queryName))
	defer timer.ObserveDuration()

	err := queryFunc()
	if err != nil {
		m.errorsTotal.WithLabelValues("db_query").Inc()
	}
	return err
}

func (m Metrics) MetricsHandler() http.Handler {
	return promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{})
}
