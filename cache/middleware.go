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

package cache

import (
	"net/http"
)

// DisableCacheOnGetRequests sets Cache-Control: no-store on GET responses by
// default. Handlers can overwrite this when required.
func DisableCacheOnGetRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// set it before callback is executed allowing handlers to
			// overwrite it
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
