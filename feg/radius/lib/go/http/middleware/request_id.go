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

	"github.com/google/uuid"
)

// RequestID returns a X-Request-ID middleware.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(header.XRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		w.Header().Set(header.XRequestID, rid)
		next.ServeHTTP(w, r)
	})
}
