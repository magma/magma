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
	"context"
	"fmt"
	"time"

	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state/wrappers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	validSecsSinceStateReported = 180
)

type HAServicer struct{}

// NewHAServicer creates a new service implementing the HA proto file
func NewHAServicer() lte_protos.HaServer {
	return &HAServicer{}
}

// GetEnodebOffloadState fetches all primary gateways that the calling gateway
// is in a gateway pool with. For each of these gateways, it then fetches the
// offload state for each of the gateway's ENBs.
func (s *HAServicer) GetEnodebOffloadState(ctx context.Context, req *lte_protos.GetEnodebOffloadStateRequest) (*lte_protos.GetEnodebOffloadStateResponse, error) {
	ret := &lte_protos.GetEnodebOffloadStateResponse{}
	secondaryGw := protos.GetClientGateway(ctx)
	if secondaryGw == nil {
		return ret, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !secondaryGw.Registered() {
		return ret, status.Errorf(codes.PermissionDenied, "gateway is not registered")
	}
	cfg, err := configurator.LoadEntityConfig(secondaryGw.GetNetworkId(), lte.CellularGatewayEntityType, secondaryGw.LogicalId, lte_models.EntitySerdes)
	if err != nil {
		errors.Wrap(err, "unable to load cellular gateway configs to find primary gateway's in its pool")
		return ret, err
	}
	cellularCfg, ok := cfg.(*lte_models.GatewayCellularConfigs)
	if !ok {
		return ret, status.Errorf(codes.Internal, "could not convert stored config to type GatewayCellularConfigs for gw %s", secondaryGw.LogicalId)
	}
	if cellularCfg.Pooling == nil || len(cellularCfg.Pooling) == 0 {
		return ret, fmt.Errorf("gateway '%s' is not configured in a gateway pool", secondaryGw.LogicalId)
	}

	// All gateway pool records must have the same capacity, so use the first
	// entry
	callingRelativeCapacity := cellularCfg.Pooling[0].MmeRelativeCapacity
	gwIDsToEnbs := map[string][]string{}
	for _, record := range cellularCfg.Pooling {
		gwIDsToEnbsInPool, err := s.getPrimaryGatewaysToEnodebs(secondaryGw.GetNetworkId(), string(record.GatewayPoolID), callingRelativeCapacity)
		if err != nil {
			return &lte_protos.GetEnodebOffloadStateResponse{}, err
		}
		// Since a gateway can be in multiple pools, it is possible
		// there could be key collisions. Since the ENBs values will be the
		// same each time, this is okay.
		for k, v := range gwIDsToEnbsInPool {
			gwIDsToEnbs[k] = v
		}
	}
	glog.V(2).Infof("Found the following primary gatewayIDs to ENB SNs: %v", gwIDsToEnbs)

	enbSNsToOffloadState := map[uint32]lte_protos.GetEnodebOffloadStateResponse_EnodebOffloadState{}
	for primaryGwID, enbs := range gwIDsToEnbs {
		isCheckinValid, err := s.isGatewayCheckinValid(secondaryGw.NetworkId, primaryGwID)
		// Since a secondary gateway can serve multiple primary gateways, if
		// we are unable to fetch checkin state for a primary, we should
		// still continue gathering the offload state to return, rather
		// than returning the error.
		if err != nil {
			glog.Error(err)
			continue
		} else if !isCheckinValid {
			continue
		}
		for _, enb := range enbs {
			offloadState, err := s.getOffloadStateForEnb(secondaryGw.NetworkId, primaryGwID, enb)
			// Since a secondary gateway can offload multiple ENBs, if we are
			// unable to fetch offload state for an ENB, we should continue
			// gathering offload state for other ENBs, rather than returning
			// the error.
			if err != nil {
				glog.Error(err)
				continue
			}
			enbID, err := s.getEnodebID(secondaryGw.NetworkId, secondaryGw.LogicalId, enb)
			if err != nil {
				glog.Error(err)
				continue
			}
			enbSNsToOffloadState[enbID] = offloadState
		}
	}
	return &lte_protos.GetEnodebOffloadStateResponse{EnodebOffloadStates: enbSNsToOffloadState}, nil
}

func (s *HAServicer) getPrimaryGatewaysToEnodebs(networkID string, gatewayPoolID string, callingRelativeCapacity uint32) (map[string][]string, error) {
	poolEnt, err := configurator.LoadEntity(
		networkID, lte.CellularGatewayPoolEntityType, gatewayPoolID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		lte_models.EntitySerdes,
	)
	if err != nil {
		return map[string][]string{}, err
	}
	primaryGatewaysToEnbs := map[string][]string{}
	for _, gw := range poolEnt.Associations.Filter(lte.CellularGatewayEntityType) {
		ent, err := configurator.LoadEntity(
			networkID, lte.CellularGatewayEntityType, gw.Key,
			configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
			lte_models.EntitySerdes,
		)
		if err != nil {
			return map[string][]string{}, err
		}
		cellularCfg, ok := ent.Config.(*lte_models.GatewayCellularConfigs)
		if !ok {
			return map[string][]string{}, fmt.Errorf("could not convert stored config to type GatewayCellularConfigs for gw %s", gw.Key)
		}
		if cellularCfg.Pooling == nil || len(cellularCfg.Pooling) == 0 {
			return map[string][]string{}, fmt.Errorf("gateway '%s' is not configured in a gateway pool", gw.Key)
		}
		// primary gateways are those with strictly higher relative capacity
		// than the calling gateway
		currentRelativeCapacity := cellularCfg.Pooling[0].MmeRelativeCapacity
		if currentRelativeCapacity > callingRelativeCapacity {
			enbs := ent.Associations.Filter(lte.CellularEnodebEntityType).Keys()
			primaryGatewaysToEnbs[gw.Key] = enbs
		}
	}
	return primaryGatewaysToEnbs, nil
}

func (s *HAServicer) isGatewayCheckinValid(networkID string, gatewayID string) (bool, error) {
	hwID, err := s.getHardwareIDFromGatewayID(networkID, gatewayID)
	if err != nil {
		return false, err
	}
	status, err := wrappers.GetGatewayStatus(networkID, hwID)
	if err != nil {
		return false, err
	}
	timeSinceCheckin := time.Now().Unix() - int64(status.CheckinTime)/1000
	return timeSinceCheckin < validSecsSinceStateReported, nil
}

func (s *HAServicer) getOffloadStateForEnb(networkID string, primaryGwID string, enbSN string) (lte_protos.GetEnodebOffloadStateResponse_EnodebOffloadState, error) {
	enodebState, err := lte_service.GetEnodebState(networkID, primaryGwID, enbSN)
	if err != nil {
		return lte_protos.GetEnodebOffloadStateResponse_NO_OP, err
	}
	timeSinceReported := time.Now().Unix() - int64(enodebState.TimeReported)/1000
	if timeSinceReported > validSecsSinceStateReported {
		glog.V(2).Infof("Returning NO_OP offload state for ENB %s; Time is %d secs too stale", enbSN, timeSinceReported)
		return lte_protos.GetEnodebOffloadStateResponse_NO_OP, nil
	}
	if !*enodebState.EnodebConnected || !*enodebState.MmeConnected {
		glog.V(2).Infof("Returning NO_OP offload state for ENB %s; Enodeb state does not have Enodeb connected or MME connected", enbSN)
		return lte_protos.GetEnodebOffloadStateResponse_NO_OP, nil
	}
	if enodebState.UesConnected == 0 {
		glog.V(2).Infof("Returning PRIMARY_CONNECTED offload state for ENB %s; no UEs connected", enbSN)
		return lte_protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED, nil
	}
	glog.V(2).Infof("Returning PRIMARY_CONNECTED_AND_SERVING_UES offload state for ENB %s", enbSN)
	return lte_protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED_AND_SERVING_UES, nil
}

func (s *HAServicer) getEnodebID(networkID string, hwID string, enodebSn string) (uint32, error) {
	cfg, err := configurator.LoadEntityConfig(networkID, lte.CellularEnodebEntityType, enodebSn, serdes.Entity)
	if err != nil {
		return 0, err
	}
	enodebCfg, ok := cfg.(*lte_models.EnodebConfig)
	if !ok {
		return 0, fmt.Errorf("could not convert Enodeb config to proper type for ENB '%s'", enodebSn)
	}
	switch enodebCfg.ConfigType {
	case "MANAGED":
		if enodebCfg == nil || enodebCfg.ManagedConfig == nil {
			return 0, fmt.Errorf("could not extract ENB ID from config for ENB '%s'; config was nil", enodebSn)
		}
		return *enodebCfg.ManagedConfig.CellID, nil
	case "UNMANAGED":
		if enodebCfg == nil || enodebCfg.UnmanagedConfig == nil {
			return 0, fmt.Errorf("could not extract ENB ID from config for ENB '%s'; config was nil", enodebSn)
		}
		return *enodebCfg.UnmanagedConfig.CellID, nil
	default:
		return 0, fmt.Errorf("invalid enodeb config type '%s' for ENB '%s'", enodebCfg.ConfigType, enodebSn)
	}
}

func (s *HAServicer) getHardwareIDFromGatewayID(networkID string, gatewayID string) (string, error) {
	ent, err := configurator.LoadEntity(
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadMetadata: true}, serdes.Entity,
	)
	if err != nil {
		return "", err
	}
	return ent.PhysicalID, nil
}
