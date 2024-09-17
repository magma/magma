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

package health_client

import (
	"context"

	"magma/cwf/k8s/cwf_operator/pkg/registry"
	"magma/feg/cloud/go/protos"
	orc8rprotos "magma/orc8r/lib/go/protos"
)

type HealthClient struct {
	registry.ConnectionRegistry
}

// NewHealthClient creates a new health client with an embedded connection
// registry, to allow for reuse of existing connections.
func NewHealthClient() *HealthClient {
	connReg := registry.NewK8sConnectionRegistry()
	return &HealthClient{
		ConnectionRegistry: connReg,
	}
}

// getGatewayClient is a utility function to get an RPC connection to
// the health service at the provided service address.
func (h *HealthClient) getGatewayClient(addr string, port int) (protos.ServiceHealthClient, error) {
	conn, err := h.GetConnection(addr, port)
	if err != nil {
		return nil, err
	}
	client := protos.NewServiceHealthClient(conn)
	return client, nil
}

// GetHealthStatus calls the provided service address to obtain health
// status from a gateway.
func (h *HealthClient) GetHealthStatus(address string, port int) (*protos.HealthStatus, error) {
	client, err := h.getGatewayClient(address, port)
	if err != nil {
		return nil, err
	}
	return client.GetHealthStatus(context.Background(), &orc8rprotos.Void{})
}

// Enable calls the provided service address to enable gateway functionality
// after a standby gateway is promoted.
func (h *HealthClient) Enable(address string, port int) error {
	client, err := h.getGatewayClient(address, port)
	if err != nil {
		return err
	}
	_, err = client.Enable(context.Background(), &orc8rprotos.Void{})
	return err
}

// Disable calls the provided service address to disable gateway functionality
// after an active gateway is demoted.
func (h *HealthClient) Disable(address string, port int) error {
	req := &protos.DisableMessage{}
	client, err := h.getGatewayClient(address, port)
	if err != nil {
		return err
	}
	_, err = client.Disable(context.Background(), req)
	return err
}
