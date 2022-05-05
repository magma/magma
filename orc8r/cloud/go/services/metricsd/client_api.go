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

	"github.com/golang/glog"

	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/protos"
	service_registry "magma/orc8r/lib/go/registry"
)

// PushMetrics pushes a set of metrics to the metricsd service.
func PushMetrics(ctx context.Context, metrics *protos.PushedMetricsContainer) error {
	client, err := getCloudMetricsdClient()
	if err != nil {
		return err
	}
	_, err = client.Push(ctx, metrics)
	return err
}

// getCloudMetricsdClient is a utility function to get a RPC connection to the
// metricsd service
func getCloudMetricsdClient() (protos.CloudMetricsControllerClient, error) {
	conn, err := service_registry.GetConnection(ServiceName, protos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewCloudMetricsControllerClient(conn), err
}
