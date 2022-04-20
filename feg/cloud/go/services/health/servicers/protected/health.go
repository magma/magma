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

package protected

import (
	"context"
	"fmt"

	fegprotos "magma/feg/cloud/go/protos"
	servicers "magma/feg/cloud/go/services/health/servicers/southbound"
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/blobstore"
)

type CloudHealthServer struct {
	Store storage.HealthBlobstore
}

func NewCloudHealthServer(factory blobstore.StoreFactory) (*CloudHealthServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	store, err := storage.NewHealthBlobstore(factory)
	return &CloudHealthServer{
		store,
	}, err
}

// GetHealth fetches the health stats for a given gateway
// represented by a (networkID, logicalId)
func (srv *CloudHealthServer) GetHealth(ctx context.Context, req *fegprotos.GatewayStatusRequest) (*fegprotos.HealthStats, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil GatewayHealthRequest")
	}
	if len(req.GetNetworkId()) == 0 || len(req.GetLogicalId()) == 0 {
		return nil, fmt.Errorf("Empty GatewayHealthRequest parameters provided")
	}
	gwHealthStats, err := srv.Store.GetHealth(req.NetworkId, req.LogicalId)
	if err != nil {
		return nil, fmt.Errorf("Get Health Error: '%s' for Gateway: %s", err, req.LogicalId)
	}
	// Update health status field with new HEALTHY/UNHEALTHY determination
	// as recency of an update is a factor in gateway health
	healthStatus, healthMessage, err := servicers.AnalyzeHealthStats(ctx, gwHealthStats, req.GetNetworkId())
	gwHealthStats.Health = &fegprotos.HealthStatus{
		Health:        healthStatus,
		HealthMessage: healthMessage,
	}
	return gwHealthStats, err
}

// GetClusterState takes a ClusterStateRequest containing a networkID and clusterID
// and returns the ClusterState or an error
func (srv *CloudHealthServer) GetClusterState(ctx context.Context, req *fegprotos.ClusterStateRequest) (*fegprotos.ClusterState, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil ClusterStateRequest")
	}
	if len(req.NetworkId) == 0 || len(req.ClusterId) == 0 {
		return nil, fmt.Errorf("Empty ClusterStateRequest parameters provided")
	}
	clusterState, err := srv.Store.GetClusterState(req.NetworkId, req.ClusterId)
	if err != nil {
		return nil, fmt.Errorf("Get Cluster State Error for networkID: %s, clusterID: %s; %s", req.NetworkId, req.ClusterId, err)
	}
	return clusterState, nil
}
