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

// Package servicesrs implements various relay RPCs to relay messages from FeG to Gateways via Controller
package servicers

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/feg/cloud/go/services/feg_relay/utils"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"
)

// FegToGwRelayServer is a server serving requests from FeG to Access Gateway
type FegToGwRelayServer struct {
}

// NewFegToGwRelayServer creates a new FegToGwRelayServer
func NewFegToGwRelayServer() (*FegToGwRelayServer, error) {
	return &FegToGwRelayServer{}, nil
}

func getHwIDFromIMSI(ctx context.Context, imsi string) (string, error) {
	gw := protos.GetClientGateway(ctx)
	// directoryd prefixes imsi with "IMSI" when updating the location
	if !strings.HasPrefix(imsi, "IMSI") {
		imsi = fmt.Sprintf("IMSI%s", imsi)
	}
	servedIds, err := getFegServedIds(gw.GetNetworkId())
	if err != nil {
		return "", err
	}
	for _, nid := range servedIds {
		hwID, err := directoryd.GetHWIDForIMSI(nid, imsi)
		if err == nil && len(hwID) != 0 {
			glog.V(2).Infof("IMSI to send is %v\n", imsi)
			return hwID, nil
		}
	}
	return "", fmt.Errorf("could not find gateway location for IMSI: %s", imsi)
}

func getHwIDFromTeid(ctx context.Context, teid string) (string, error) {
	gw := protos.GetClientGateway(ctx)
	servedIds, err := getFegServedIds(gw.GetNetworkId())
	if err != nil {
		return "", err
	}
	for _, nid := range servedIds {
		hwID, err := directoryd.GetHWIDForSgwCTeid(nid, teid)
		if err == nil && len(hwID) != 0 {
			glog.V(2).Infof("TEID to send is %s", teid)
			return hwID, nil
		}
		glog.V(2).Infof("hwid for teid %s not found at network %s: %s", teid, nid, err)
	}
	return "", fmt.Errorf("could not find gateway location for teid: %s", teid)
}

func getGWSGSServiceConnCtx(ctx context.Context, imsi string) (*grpc.ClientConn, context.Context, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, nil, err
	}
	hwID, err := getHwIDFromIMSI(ctx, imsi)
	if err != nil {
		errorStr := fmt.Sprintf(
			"unable to get HwID from IMSI %v. err: %v\n",
			imsi,
			err,
		)
		glog.Error(errorStr)
		return nil, nil, fmt.Errorf(errorStr)
	}
	conn, ctx, err := gateway_registry.GetGatewayConnection(
		gateway_registry.GwSgsService, hwID)
	if err != nil {
		errorStr := fmt.Sprintf(
			"unable to get connection to the gateway: %v",
			err,
		)
		return nil, nil, fmt.Errorf(errorStr)
	}
	return conn, ctx, nil
}

func getAllGWSGSServiceConnCtx(ctx context.Context) ([]*grpc.ClientConn, []context.Context, error) {
	var connList []*grpc.ClientConn
	var ctxList []context.Context

	hwIDs, err := utils.GetAllGatewayIDs(ctx)
	if err != nil {
		return connList, ctxList, err
	}
	for _, hwID := range hwIDs {
		conn, ctx, err := gateway_registry.GetGatewayConnection(
			gateway_registry.GwSgsService,
			hwID,
		)
		if err != nil {
			return connList, ctxList, err
		}
		connList = append(connList, conn)
		ctxList = append(ctxList, ctx)
	}

	return connList, ctxList, nil
}

// getFegServedIds returns ServedNetworkIds of the given FeG networkId and appends to the list all ServedNetworkIds of
// the network's Neutral Host Network if any
func getFegServedIds(networkId string) ([]string, error) {
	if len(networkId) == 0 {
		return []string{}, fmt.Errorf("Empty networkID provided.")
	}
	fegCfg, err := configurator.LoadNetworkConfig(networkId, feg.FegNetworkType, serdes.Network)
	if err != nil || fegCfg == nil {
		return []string{}, fmt.Errorf("unable to retrieve config for federation network: %s", networkId)
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		return []string{}, fmt.Errorf("invalid federation network config found for network: %s", networkId)
	}
	// getFegServedIds will always return ServedNetworkIds of the given FeG networkId, but if the given FeG Network
	// is also serving Neutral Host Network, getFegServedIds would append all ServedNetworkIds of this
	// Neutral Host Network to the result
	if len(networkFegConfigs.ServedNhIds) == 0 {
		return networkFegConfigs.ServedNetworkIds, nil
	}
	// If this is a NH FeG Network, add served networks from ServedNhIds FeG networks
	nids := networkFegConfigs.ServedNetworkIds // prepend 'local' ServedNhIds to the combined result
	glog.V(2).Infof("getFegServedIds: nonempty Served NH Networks list for network: %s", networkId)
	for _, nhNetworkId := range networkFegConfigs.ServedNhIds {
		if len(nhNetworkId) > 0 {
			nhFegCfg, err := configurator.LoadNetworkConfig(nhNetworkId, feg.FegNetworkType, serdes.Network)
			if err != nil || nhFegCfg == nil {
				glog.Errorf("unable to retrieve config for NH federation network '%s': %v", nhNetworkId, err)
				continue
			}
			nhNetworkFegConfigs, ok := nhFegCfg.(*models.NetworkFederationConfigs)
			if !ok || nhNetworkFegConfigs == nil {
				glog.Errorf("invalid FeG network config found for NH network '%s': %T", nhNetworkId, nhFegCfg)
				continue
			}
			nids = append(nids, nhNetworkFegConfigs.ServedNetworkIds...)
		}
	}
	return nids, nil
}

func validateFegContext(ctx context.Context) error {
	fegID := protos.GetClientGateway(ctx)
	if fegID == nil {
		ctxMetadata, _ := metadata.FromIncomingContext(ctx)
		errorStr := fmt.Sprintf(
			"Failed to get Identity of calling Federated Gateway from CTX Metadata: %+v",
			ctxMetadata,
		)
		glog.Error(errorStr)
		return fmt.Errorf(errorStr)
	}
	if !fegID.Registered() {
		return fmt.Errorf("federated gateway not registered")
	}
	return nil
}
