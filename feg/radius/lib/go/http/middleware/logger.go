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

	"fbc/lib/go/http/requestlog"

	"go.uber.org/zap"
)

type logger struct{ *zap.Logger }

func (l logger) Log(ent *requestlog.Entry) {
	l.Info("HTTP request",
		zap.String("method", ent.RequestMethod),
		zap.String("url", ent.RequestURL),
		zap.Int("status", ent.Status),
		zap.String("user_agent", ent.UserAgent),
		zap.String("remote_ip", ent.RemoteIP),
		zap.String("server_ip", ent.ServerIP),
		zap.String("referer", ent.Referer),
		zap.Stringer("trace_id", ent.TraceID),
		zap.Stringer("span_id", ent.SpanID),
		zap.Duration("latency", ent.Latency),
		zap.Int64("bytes_in", ent.RequestHeaderSize+ent.RequestBodySize),
		zap.Int64("bytes_out", ent.ResponseHeaderSize+ent.ResponseBodySize),
	)
}

// Logger returns an http request logging middleware.
func Logger(l *zap.Logger) func(http.Handler) http.Handler {
	if l == nil {
		panic("logger middleware requires a logger")
	}
	return func(next http.Handler) http.Handler {
		return requestlog.NewHandler(logger{l}, next)
	}
}
