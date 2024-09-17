/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package middleware

import (
	"net/http"

	"fbc/lib/go/http/header"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

// TracingOption controls the behavior of the tracing middleware.
type TracingOption func(*ochttp.Handler)

// TracingPublicEndpoint should be set to true for public accessible endpoints.
func TracingPublicEndpoint(public bool) TracingOption {
	return func(h *ochttp.Handler) {
		h.IsPublicEndpoint = public
	}
}

// Tracing returns an http request tracing middleware.
func Tracing(options ...TracingOption) func(http.Handler) http.Handler {
	handler := &ochttp.Handler{
		FormatSpanName: func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		},
	}
	for _, option := range options {
		option(handler)
	}

	return func(next http.Handler) http.Handler {
		handler.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if span := trace.FromContext(r.Context()); span != nil {
				if sc := span.SpanContext(); sc.IsSampled() {
					w.Header().Set(header.XCorrelationID, sc.TraceID.String())
				}
			}
			next.ServeHTTP(w, r)
		})
		return handler
	}
}
