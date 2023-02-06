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
	"context"

	"github.com/go-chi/chi/v5/middleware"
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
