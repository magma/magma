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
	"net/http/httptest"
	"testing"

	"fbc/lib/go/http/header"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type mockExporter struct {
	spans []*trace.SpanData
}

func (e *mockExporter) ExportSpan(s *trace.SpanData) {
	e.spans = append(e.spans, s)
}

func TestTracingMiddleware(t *testing.T) {
	exporter := &mockExporter{}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	handler = Tracing(TracingPublicEndpoint(false))(handler)
	handler.ServeHTTP(rec, req)

	spans := exporter.spans
	require.Len(t, spans, 1)

	assert.EqualValues(t, http.StatusTooManyRequests, spans[0].Attributes[ochttp.StatusCodeAttribute])
	assert.Equal(t, "HTTP GET /", spans[0].Name)
	assert.Equal(t, "/", spans[0].Attributes[ochttp.PathAttribute])
	assert.Equal(t, rec.Header().Get(header.XCorrelationID), spans[0].TraceID.String())
	assert.Equal(t, req.Method, spans[0].Attributes[ochttp.MethodAttribute])
}
