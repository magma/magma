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

package collection

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
)

// MetricCollector provides an API to query for metrics.
type MetricCollector interface {
	// GetMetrics returns a collection of prometheus MetricFamily structures
	// which contain collected metrics.
	GetMetrics() ([]*io_prometheus_client.MetricFamily, error)
}
