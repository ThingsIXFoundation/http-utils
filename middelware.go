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

package httputils

import (
	"github.com/ThingsIXFoundation/http-utils/logging"
	"github.com/ThingsIXFoundation/http-utils/tracing"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func BindStandardMiddleware(mux *chi.Mux) {
	mux.Use(middleware.Heartbeat("/healthz"))
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(logging.RequestLogger())
	mux.Use(tracing.PrometheusHTTPRequestLogger())
	mux.Use(middleware.Recoverer)
}
