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

package servicers

import (
	"fmt"
	"time"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/health/metrics"
	"magma/feg/cloud/go/services/health/storage"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthStatus int

type HealthServer struct {
	store storage.HealthBlobstore
}

func NewHealthServer(factory blobstore.BlobStorageFactory) (*HealthServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	store, err := storage.NewHealthBlobstore(factory)
	return &HealthServer{
		store,
	}, err
}

type healthConfig struct {
	services              []string
	cpuUtilThreshold      float32
	memAvailableThreshold float32
	staleUpdateThreshold  uint32
}

// GetHealth fetches the health stats for a given gateway
// represented by a (networkID, logicalId)
func (srv *HealthServer) GetHealth(ctx context.Context, req *fegprotos.GatewayStatusRequest) (*fegprotos.HealthStats, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil GatewayHealthRequest")
	}
	if len(req.GetNetworkId()) == 0 || len(req.GetLogicalId()) == 0 {
		return nil, fmt.Errorf("Empty GatewayHealthRequest parameters provided")
	}
	gwHealthStats, err := srv.store.GetHealth(req.NetworkId, req.LogicalId)
	if err != nil {
		return nil, fmt.Errorf("Get Health Error: '%s' for Gateway: %s", err, req.LogicalId)
	}
	// Update health status field with new HEALTHY/UNHEALTHY determination
	// as recency of an update is a factor in gateway health
	healthStatus, healthMessage, err := AnalyzeHealthStats(gwHealthStats, req.GetNetworkId())
	gwHealthStats.Health = &fegprotos.HealthStatus{
		Health:        healthStatus,
		HealthMessage: healthMessage,
	}
	return gwHealthStats, err
}

func (srv *HealthServer) UpdateHealth(ctx context.Context, req *fegprotos.HealthRequest) (*fegprotos.HealthResponse, error) {
	healthResponse := &fegprotos.HealthResponse{
		Action: fegprotos.HealthResponse_SYSTEM_DOWN,
		Time:   uint64(clock.Now().UnixNano()) / uint64(time.Millisecond),
	}
	if req == nil {
		return healthResponse, fmt.Errorf("Nil HealthRequest")
	}
	// Get gateway id from context
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return healthResponse, status.Errorf(
			codes.PermissionDenied, "Missing Gateway Identity")
	}
	if !gw.Registered() {
		return healthResponse, status.Errorf(
			codes.PermissionDenied, "Gateway is not registered")
	}

	networkID := gw.GetNetworkId()
	logicalID := gw.GetLogicalId()

	req.HealthStats.Time = healthResponse.Time

	// Override gateway's view of it's health with cloud's view
	healthState, healthMsg, _ := AnalyzeHealthStats(req.HealthStats, networkID)
	req.HealthStats.Health = &fegprotos.HealthStatus{
		Health:        healthState,
		HealthMessage: healthMsg,
	}
	err := srv.store.UpdateHealth(networkID, logicalID, req.HealthStats)
	if err != nil {
		healthResponse.Action = fegprotos.HealthResponse_NONE
		errMsg := fmt.Errorf("Update Health Error: '%s' for Gateway: %s", err, gw)
		glog.Error(errMsg)
		return healthResponse, errMsg
	}
	// Get FeGs registered in configurator, then make a health decision based off of the
	// the number of gateways, which gateway is active, and gateway health
	magmadGatewayTypeVal := orc8r.MagmadGatewayType
	gateways, _, err := configurator.LoadEntities(
		networkID, &magmadGatewayTypeVal, nil, nil, nil,
		configurator.EntityLoadCriteria{},
		serdes.Entity,
	)
	if err != nil {
		errMsg := fmt.Errorf(
			"Update Health Error: Could not retrieve gateways registered in network: %s",
			networkID,
		)
		glog.Error(errMsg)
		return healthResponse, errMsg
	}
	var requestedAction fegprotos.HealthResponse_RequestedAction
	switch len(gateways) {
	case 0:
		err = fmt.Errorf("Zero gateways found registered in NetworkID: %s of Gateway: %s", networkID, logicalID)
	case 1:
		requestedAction, err = srv.analyzeSingleFegState(networkID, logicalID)
	case 2:
		requestedAction, err = srv.analyzeDualFegState(networkID, logicalID, req.HealthStats, gateways)
	default:
		err = fmt.Errorf("Unsupported number of gateways registered in %s", networkID)
	}
	if err != nil {
		glog.Error(err)
		return healthResponse, fmt.Errorf("Update Health Error: '%s' for Gateway: %s", err, gw)
	}
	healthResponse.Action = requestedAction
	return healthResponse, nil
}

// GetClusterState takes a ClusterStateRequest containing a networkID and clusterID
// and returns the ClusterState or an error
func (srv *HealthServer) GetClusterState(ctx context.Context, req *fegprotos.ClusterStateRequest) (*fegprotos.ClusterState, error) {
	if req == nil {
		return nil, fmt.Errorf("Nil ClusterStateRequest")
	}
	if len(req.NetworkId) == 0 || len(req.ClusterId) == 0 {
		return nil, fmt.Errorf("Empty ClusterStateRequest parameters provided")
	}
	clusterState, err := srv.store.GetClusterState(req.NetworkId, req.ClusterId)
	if err != nil {
		return nil, fmt.Errorf("Get Cluster State Error for networkID: %s, clusterID: %s; %s", req.NetworkId, req.ClusterId, err)
	}
	return clusterState, nil
}

// analyzeDualFegState finds the current active gateway for the provided networkID.
// If the current active is unhealthy and the standby is healthy, a failover occurs.
// Otherwise, the state is left as is. The action returned is dependent on whether
// the request is from the active or standby
func (srv *HealthServer) analyzeDualFegState(
	networkID string,
	gatewayID string,
	gatewayHealth *fegprotos.HealthStats,
	clusterGateways []configurator.NetworkEntity,
) (fegprotos.HealthResponse_RequestedAction, error) {
	if gatewayHealth == nil {
		return fegprotos.HealthResponse_NONE, fmt.Errorf("Nil GatewayHealth provided")
	}
	// Get cluster state, initializing the active to the current gateway if the clusterState doesn't already exist
	clusterState, err := srv.store.GetClusterState(networkID, gatewayID)
	if err != nil {
		return fegprotos.HealthResponse_NONE, fmt.Errorf(
			"Error while trying to get clusterState for network: %s, gateway: %s; %s",
			networkID,
			gatewayID,
			err,
		)
	}
	activeID := clusterState.ActiveGatewayLogicalId

	// Sanity check to ensure that the active gateway is registered in magmad
	if !isActiveGatewayRegistered(activeID, clusterGateways) {
		reason := "active is not registered"
		return srv.failover(networkID, gatewayID, activeID, gatewayID, reason)
	}

	// We need to get the GatewayID and health for the other FeG in the cluster
	otherGatewayID := getOtherGatewayID(gatewayID, clusterGateways)
	otherGatewayHealth, err := srv.store.GetHealth(networkID, otherGatewayID)
	if err != nil {
		glog.Errorf("Unable to retrieve health data for gateway: %s; %s", otherGatewayID, err)

		// If we can't get the health data for the active, failover to standby
		if otherGatewayID == activeID {
			reason := "unable to get health of active"
			return srv.failover(networkID, gatewayID, activeID, gatewayID, reason)
		}
		// If we can't get the health data for the standby, leave the active as is
		return fegprotos.HealthResponse_SYSTEM_UP, nil
	}
	currentHealth, currentHealthDescription, err := AnalyzeHealthStats(gatewayHealth, networkID)
	if err != nil {
		return fegprotos.HealthResponse_NONE, err
	}
	otherHealth, otherHealthDescription, err := AnalyzeHealthStats(otherGatewayHealth, networkID)
	if err != nil {
		return fegprotos.HealthResponse_NONE, err
	}
	// Update gauge metric for how many gateways are healthy
	metrics.SetHealthyGatewayMetric(networkID, currentHealth, otherHealth)

	// Determine what to send back based off of health of active and standby, as well as where the request is from
	if gatewayID == activeID {
		return srv.determineAction(networkID, gatewayID, gatewayID, currentHealth, otherGatewayID, otherHealth, currentHealthDescription)
	}
	return srv.determineAction(networkID, gatewayID, otherGatewayID, otherHealth, gatewayID, currentHealth, otherHealthDescription)
}

// determineAction compares the health status of the two FeGs and determines
// if a failover should occur. The action returned is dependent on health
// status as well as which FeG the request is from.
func (srv *HealthServer) determineAction(
	networkID string,
	currentID string,
	activeID string,
	activeHealth fegprotos.HealthStatus_HealthState,
	standbyID string,
	standbyHealth fegprotos.HealthStatus_HealthState,
	activeHealthDesc string,
) (fegprotos.HealthResponse_RequestedAction, error) {
	// Only failover if active is unhealthy and standby is healthy
	if activeHealth == fegprotos.HealthStatus_UNHEALTHY && standbyHealth == fegprotos.HealthStatus_HEALTHY {
		return srv.failover(networkID, standbyID, activeID, currentID, activeHealthDesc)
	}

	// Otherwise, active stays UP and standby stays DOWN
	if currentID == activeID {
		return fegprotos.HealthResponse_SYSTEM_UP, nil
	}
	return fegprotos.HealthResponse_SYSTEM_DOWN, nil
}

// failover updates the active gateway to a new active and returns the appropriate
// action depending on which gateway the request is from (Active or Standby)
func (srv *HealthServer) failover(
	networkID string,
	newActive string,
	oldActive string,
	currentID string,
	reason string,
) (fegprotos.HealthResponse_RequestedAction, error) {
	glog.Infof("Failing over for network: %s from: %s to: %s; Reason: %s", networkID, oldActive, newActive, reason)
	metrics.ActiveGatewayChanged.WithLabelValues(networkID).Inc()
	err := srv.store.UpdateClusterState(networkID, networkID, newActive)
	if err != nil {
		errMsg := fmt.Errorf(
			"Unable to store updated cluster state for networkID %s from: %s to: %s ; %s",
			networkID,
			oldActive,
			newActive,
			err,
		)
		return fegprotos.HealthResponse_NONE, errMsg
	}
	if currentID == newActive {
		return fegprotos.HealthResponse_SYSTEM_UP, nil
	}
	return fegprotos.HealthResponse_SYSTEM_DOWN, nil
}

// analyzeSingleFegState ensures that the active ID in clusterState is set correctly.
// It then returns SYSTEM_UP to ensure that a solo feg will always remain ACTIVE
func (srv *HealthServer) analyzeSingleFegState(
	networkID string,
	gatewayID string,
) (fegprotos.HealthResponse_RequestedAction, error) {
	clusterState, err := srv.store.GetClusterState(networkID, gatewayID)
	if err != nil {
		return fegprotos.HealthResponse_SYSTEM_UP, err
	}
	// If current gatewayID is listed as active, then stay ACTIVE regardless of health
	if gatewayID == clusterState.ActiveGatewayLogicalId {
		return fegprotos.HealthResponse_SYSTEM_UP, err
	}
	// Otherwise there is a mismatch, and active needs to be updated
	glog.V(2).Infof("Updating active for networkID: %s to: %s", networkID, gatewayID)

	err = srv.store.UpdateClusterState(networkID, networkID, gatewayID)
	if err != nil {
		return fegprotos.HealthResponse_SYSTEM_UP, err
	}
	return fegprotos.HealthResponse_SYSTEM_UP, nil
}

// AnalyzeHealthStats takes a HealthStats proto and determines if it is
// HEALTHY or UNHEALTHY based on the configuration for the provided networkID
func AnalyzeHealthStats(
	healthData *fegprotos.HealthStats,
	networkID string,
) (fegprotos.HealthStatus_HealthState, string, error) {
	config := GetHealthConfigForNetwork(networkID)
	if healthData == nil {
		return fegprotos.HealthStatus_UNHEALTHY, "", fmt.Errorf("Nil HealthStats provided")
	}
	updateDelta := clock.Now().Unix() - int64(healthData.Time)/1000
	if updateDelta > int64(config.staleUpdateThreshold) {
		return fegprotos.HealthStatus_UNHEALTHY, "Health update is stale", nil
	}
	if !isSystemHealthy(healthData.GetSystemStatus(), config) {
		return fegprotos.HealthStatus_UNHEALTHY, "System unhealthy", nil
	}
	for _, service := range config.services {
		if !isServiceHealthy(healthData.ServiceStatus, service) {
			return fegprotos.HealthStatus_UNHEALTHY, fmt.Sprintf("Service: %s unhealthy", service), nil
		}
	}
	return fegprotos.HealthStatus_HEALTHY, "OK", nil
}

func isSystemHealthy(status *fegprotos.SystemHealthStats, config *healthConfig) bool {
	if status.CpuUtilPct >= config.cpuUtilThreshold {
		return false
	}
	usedMemoryBytes := status.MemTotalBytes - status.MemAvailableBytes
	exceedsMemThreshold := status.MemTotalBytes != 0 &&
		float64(usedMemoryBytes)/float64(status.MemTotalBytes) >= float64(config.memAvailableThreshold)
	return !exceedsMemThreshold
}

func isServiceHealthy(
	serviceData map[string]*fegprotos.ServiceHealthStats,
	serviceName string,
) bool {
	srvStatus, statusFound := serviceData[serviceName]
	if !statusFound {
		return false
	}
	if srvStatus.ServiceState == fegprotos.ServiceHealthStats_UNAVAILABLE || srvStatus.ServiceHealthStatus == nil {
		return false
	}
	if srvStatus.ServiceHealthStatus.Health != fegprotos.HealthStatus_HEALTHY {
		return false
	}
	return true
}

// getOtherGatewayID gets the gatewayID of the gateway in 'gateways' that is not 'gatewayID'
// If more than two gateways exist in the slice, an empty string is returned
func getOtherGatewayID(gatewayID string, gateways []configurator.NetworkEntity) string {
	if len(gateways) > 2 {
		return ""
	}
	for _, gw := range gateways {
		if gatewayID != gw.Key {
			return gw.Key
		}
	}
	return ""
}

func isActiveGatewayRegistered(activeID string, gateways []configurator.NetworkEntity) bool {
	for _, gateway := range gateways {
		if gateway.Key == activeID {
			return true
		}
	}
	return false
}
