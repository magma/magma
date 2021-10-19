/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests obsidian receives",
		},
	)
	respStatuses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Number of obsidian response by status",
		},
		[]string{"code", "method"},
	)
	errorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "error_count",
			Help: "Number of http errors obsidian receives",
		},
	)
)

func init() {
	prometheus.MustRegister(requestCount, respStatuses, errorCount)
}


func isServerErrCode(code int) bool {
	return code >= http.StatusInternalServerError && code <= http.StatusNetworkAuthenticationRequired
}

// isHttpErrCode returns true for any non-2xx responses
func isHttpErrCode(code int) bool {
	return code < http.StatusOK || code > http.StatusIMUsed
}

// Logger is the middleware function
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			c.Error(err)
		}
		requestCount.Inc()
		status := c.Response().Status
		respStatuses.WithLabelValues(strconv.Itoa(status), c.Request().Method).Inc()
		if isServerErrCode(status) {
			glog.Infof("REST HTTP Error: %s, Status: %d", err, status)
			errorCount.Inc()
		} else if isHttpErrCode(status) {
			glog.V(1).Infof("REST HTTP Error: %s, Status: %d", err, status)
			errorCount.Inc()
		} else {
			glog.V(2).Infof(
				"REST API code: %v, method: %v, url: %v\n",
				status,
				c.Request().Method,
				c.Request().URL,
			)
		}
		return err
	}
}
