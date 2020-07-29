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
)

func init() {
	prometheus.MustRegister(requestCount, respStatuses)
}

// CollectStats is the middleware function
func CollectStats(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		requestCount.Inc()
		status := strconv.Itoa(c.Response().Status)
		respStatuses.WithLabelValues(status, c.Request().Method).Inc()
		glog.V(2).Infof(
			"REST API code: %v, method: %v, url: %v\n",
			status,
			c.Request().Method,
			c.Request().URL,
		)
		return nil
	}
}
