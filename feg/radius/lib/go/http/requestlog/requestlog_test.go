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

package requestlog

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
)

func TestHandler(t *testing.T) {
	const (
		requestMsg  = "Hello, World!"
		responseMsg = "I see you."
		userAgent   = "Request Log Test UA"
		referer     = "http://www.example.com/"
	)
	r, err := http.NewRequest(http.MethodPost, "http://localhost/foo", strings.NewReader(requestMsg))
	require.NoError(t, err)
	r.Header.Set("User-Agent", userAgent)
	r.Header.Set("Referer", referer)
	requestHdrSize := len(fmt.Sprintf("User-Agent: %s\r\nReferer: %s\r\nContent-Length: %v\r\n", userAgent, referer, len(requestMsg)))
	responseHdrSize := len(fmt.Sprintf("Content-Length: %v\r\n", len(responseMsg)))
	ent, spanCtx, err := roundTrip(r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(responseMsg)))
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, responseMsg)
	}))
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, ent.RequestMethod)
	assert.Equal(t, "/foo", ent.RequestURL)
	assert.True(t, ent.RequestHeaderSize >= int64(requestHdrSize))
	assert.Len(t, requestMsg, int(ent.RequestBodySize))
	assert.Equal(t, userAgent, ent.UserAgent)
	assert.Equal(t, referer, ent.Referer)
	assert.Equal(t, "HTTP/1.1", ent.Proto)
	assert.Equal(t, http.StatusOK, ent.Status)
	assert.True(t, ent.ResponseHeaderSize >= int64(responseHdrSize))
	assert.Len(t, responseMsg, int(ent.ResponseBodySize))
	assert.Equal(t, spanCtx.TraceID, ent.TraceID)
	assert.Equal(t, spanCtx.SpanID, ent.SpanID)
}

type testSpanHandler struct {
	h       http.Handler
	spanCtx *trace.SpanContext
}

func (sh *testSpanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "test")
	defer span.End()
	r = r.WithContext(ctx)
	sc := trace.FromContext(ctx).SpanContext()
	sh.spanCtx = &sc
	sh.h.ServeHTTP(w, r)
}

func roundTrip(r *http.Request, h http.Handler) (*Entry, *trace.SpanContext, error) {
	capture := new(captureLogger)
	hh := NewHandler(capture, h)
	handler := &testSpanHandler{h: hh}
	s := httptest.NewServer(handler)
	defer s.Close()
	r.URL.Host = s.URL[len("http://"):]
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, nil, err
	}
	resp.Body.Close()
	return &capture.ent, handler.spanCtx, nil
}

type captureLogger struct {
	ent Entry
}

func (cl *captureLogger) Log(ent *Entry) {
	cl.ent = *ent
}
