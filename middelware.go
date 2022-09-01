package httputils

import (
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/http-utils/tracing"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func BindStandardMiddleware(mux *chi.Mux) {
	mux.Use(middleware.Heartbeat("/healthz"))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(logging.RequestLogger())
	mux.Use(tracing.PrometheusHTTPRequestLogger())
	mux.Use(middleware.Recoverer)
}
