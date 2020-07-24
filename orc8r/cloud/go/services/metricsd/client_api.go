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

package metricsd

import (
	"context"

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	service_registry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// PushMetrics pushes a set of metrics to the metricsd service.
func PushMetrics(metrics protos.PushedMetricsContainer) error {
	client, err := getMetricsdClient()
	if err != nil {
		return err
	}
	_, err = client.Push(context.Background(), &metrics)
	return err
}

// getMetricsdClient is a utility function to get a RPC connection to the
// metricsd service
func getMetricsdClient() (protos.MetricsControllerClient, error) {
	conn, err := service_registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMetricsControllerClient(conn), err
}
