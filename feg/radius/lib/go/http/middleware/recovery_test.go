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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"fbc/lib/go/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestRecoveryMiddleware(t *testing.T) {
	core, o := observer.New(zap.ErrorLevel)
	logger := log.NewFactory(zap.New(core))
	errBadHandler := errors.New("bad handler")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(errBadHandler)
	})
	Recovery(RecoveryLogger(logger))(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	o = o.FilterMessage("panic recovery").FilterField(zap.Error(errBadHandler))
	require.Equal(t, 1, o.Len())
	assert.Condition(t, func() bool {
		for _, field := range o.TakeAll()[0].Context {
			if field.Key == "stacktrace" {
				return true
			}
		}
		return false
	})
}
