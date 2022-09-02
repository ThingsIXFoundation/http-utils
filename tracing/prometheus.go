package tracing

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// prometheusmHTTPLog is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type prometheusmHTTPLog struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

func (c prometheusmHTTPLog) handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		rp := routePattern(r)
		c.reqs.WithLabelValues(fmt.Sprintf("%d", ww.Status()), r.Method, rp).Inc()
		c.latency.WithLabelValues(fmt.Sprintf("%d", ww.Status()), r.Method, rp).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}

// PrometheusHTTPRequestLogger collects and emits HTTP requests prometheus metrics under the keys:
// - http_requests_total: vector with the number of times an endpoint is called
// - http_request_duration_ms: histogram of HTTP response times
func PrometheusHTTPRequestLogger() func(http.Handler) http.Handler {
	m := prometheusmHTTPLog{
		reqs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "how many http requests processed, partitioned by status code, method and path",
		}, []string{"code", "method", "path"}),
		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_ms",
			Help:    "how long it took to process the request, partitioned by status code, method and HTTP path",
			Buckets: []float64{100, 250, 500, 1000, 2500},
		}, []string{"code", "method", "path"}),
	}

	prometheus.MustRegister(m.reqs)
	prometheus.MustRegister(m.latency)

	return m.handler
}

func routePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		return pattern
	}
	// empty which means one of the middlewares returned an error before
	// the final handler is located based on the path. Find it ourselves.
	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// path isn't mapped in router
		logrus.WithField("path", routePath).Warn("unknown endpoint requested")
		return "<unknown>"
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}
