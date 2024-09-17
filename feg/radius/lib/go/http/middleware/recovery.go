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

	"fbc/lib/go/log"

	"go.uber.org/zap"
)

type (
	// RecoveryOption controls the behavior of the recovery middleware.
	RecoveryOption func(*recoveryOptions)

	recoveryOptions struct {
		logger log.Factory
	}
)

// RecoveryLogger returns a RecoveryOption that sets the logger
// for panic recovery logging.
func RecoveryLogger(logger log.Factory) RecoveryOption {
	return func(opts *recoveryOptions) {
		opts.logger = logger
	}
}

// Recovery returns an http request panic recovery middleware.
func Recovery(options ...RecoveryOption) func(http.Handler) http.Handler {
	opts := recoveryOptions{}
	for _, option := range options {
		option(&opts)
	}

	logger := opts.logger
	if logger == nil {
		logger = log.NewNopFactory()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					err, _ := err.(error)
					logger.For(r.Context()).Error("panic recovery", zap.Error(err), zap.Stack("stacktrace"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
