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

// Package health provides a thin client for using the health service from other cloud services.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package health

import (
	"context"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// getHealthClient is a utility function to get an RPC connection to the
// Health service
func getHealthClient() (protos.HealthClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewHealthClient(conn), nil
}

// GetActiveGateway returns the active federated gateway in the network specified by networkID
func GetActiveGateway(networkID string) (string, error) {
	client, err := getHealthClient()
	if err != nil {
		return "", err
	}

	// Currently, we use networkID as clusterID as we only support one cluster per network
	clusterState, err := client.GetClusterState(context.Background(), &protos.ClusterStateRequest{
		NetworkId: networkID,
		ClusterId: networkID,
	})
	if err != nil {
		return "", err
	}
	return clusterState.ActiveGatewayLogicalId, nil
}

// GetHealth fetches the health stats for a given gateway
// represented by a (networkID, logicalId)
func GetHealth(networkID string, logicalID string) (*protos.HealthStats, error) {
	if len(networkID) == 0 {
		return nil, fmt.Errorf("Empty networkId provided")
	}
	if len(logicalID) == 0 {
		return nil, fmt.Errorf("Empty logicalId provided")
	}
	client, err := getHealthClient()
	if err != nil {
		return nil, err
	}

	gatewayHealthReq := &protos.GatewayStatusRequest{
		NetworkId: networkID,
		LogicalId: logicalID,
	}
	return client.GetHealth(context.Background(), gatewayHealthReq)
}
