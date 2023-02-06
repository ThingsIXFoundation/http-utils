// Copyright 2023 Stichting ThingsIX Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package logging

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// RequestLogger writes HTTP request/response information to the log.
// Making it possible to trace each individual request from request, to action, to reply.
func RequestLogger() func(next http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			reqSize := r.ContentLength
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				fields := logrus.Fields{
					"status_code": ww.Status(),
					"req_size":    reqSize,
					"resp_size":   ww.BytesWritten(),
					"duration_ms": int64(time.Since(t1)) / 1000000,
					"remote_ip":   remoteIP,
					"method":      r.Method,
					"path":        r.RequestURI,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				logrus.WithFields(fields).Info("HTTP request")
			}()

			h.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
