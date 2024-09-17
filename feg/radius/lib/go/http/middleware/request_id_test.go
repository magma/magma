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
)

func TestRequestIDMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*http.Request)
		expect  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:    "without-id",
			prepare: func(req *http.Request) {},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.NotEmpty(t, rec.Header().Get(header.XRequestID))
			},
		},
		{
			name: "with-id",
			prepare: func(req *http.Request) {
				req.Header.Set(header.XRequestID, "f2314c55814a")
			},
			expect: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Equal(t, "f2314c55814a", rec.Header().Get(header.XRequestID))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			tt.prepare(req)
			handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
			RequestID(handler).ServeHTTP(rec, req)
			tt.expect(t, rec)
		})
	}
}
