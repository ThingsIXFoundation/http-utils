package logging

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// WithContext logs the request ID if present in the context.
// If there is not request id it will not do anything.
func WithContext(ctx context.Context) *logrus.Entry {
	if rid := middleware.GetReqID(ctx); rid != "" {
		return logrus.WithField("request_id", rid)
	}
	return logrus.NewEntry(logrus.StandardLogger())
}
