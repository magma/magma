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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerMiddleware(t *testing.T) {
	assert.Panics(t, func() { Logger(nil) })
	core, o := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	ctx, span := trace.StartSpan(context.Background(), "test")
	defer span.End()

	req := httptest.NewRequest(http.MethodPost, "/foo/bar", nil).WithContext(ctx)
	rec := httptest.NewRecorder()
	Logger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "request error", http.StatusInsufficientStorage)
	})).ServeHTTP(rec, req)

	fields := o.TakeAll()
	require.Len(t, fields, 1)
	assert.Equal(t, "HTTP request", fields[0].Message)
	assert.Equal(t, zap.InfoLevel, fields[0].Level)
	m := fields[0].ContextMap()
	assert.Equal(t, http.MethodPost, m["method"])
	assert.EqualValues(t, http.StatusInsufficientStorage, m["status"])
	assert.Equal(t, "/foo/bar", m["url"])
	assert.Equal(t, span.SpanContext().TraceID.String(), m["trace_id"])
	assert.Equal(t, span.SpanContext().SpanID.String(), m["span_id"])
}
