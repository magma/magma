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

// Package utils includes common helper functions used by FeG Rely components/services
package utils

import (
	"context"
	"fmt"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

// GetAllGatewaysIDs returns all Gateways served by calling FeG,
// Calling FeG ID is set by Identity framework & retrieved from ctx
func GetAllGatewayIDs(ctx context.Context) ([]string, error) {
	var err error
	res := []string{}
	fegId := protos.GetClientGateway(ctx)
	if fegId == nil {
		ctxMetadata, _ := metadata.FromIncomingContext(ctx)
		return res, fmt.Errorf(
			"Failed to get Identity of calling Federated Gateway from CTX Metadata: %+v", ctxMetadata)
	}
	networkID := fegId.GetNetworkId()
	logicalID := fegId.GetLogicalId()
	if len(networkID) == 0 || len(logicalID) == 0 {
		return res, fmt.Errorf("Unregistered Federated Gateway: %s", fegId.String())
	}
	cfg, err := getFegCfg(networkID, logicalID)
	if err != nil {
		return res, fmt.Errorf("Error getting Federated Gateway %s:%s configs: %v", networkID, logicalID, err)
	}

	// Find as many gateways as possible, swallowing intermediate errors
	for _, networkID := range cfg.ServedNetworkIds {
		gateways, _, err := configurator.LoadEntities(
			networkID, swag.String(orc8r.MagmadGatewayType), nil, nil, nil,
			configurator.EntityLoadCriteria{},
			serdes.Entity,
		)
		if err != nil {
			glog.Errorf("List Network '%s' Gateways error: %v", networkID, err)
			continue
		}
		for _, gatewayEntity := range gateways {
			record, err := device.GetDevice(networkID, orc8r.AccessGatewayRecordType, gatewayEntity.PhysicalID, serdes.Device)
			if err != nil {
				glog.Errorf("Find Gateway Record Error: %v for Gateway %s:%s", err, networkID, gatewayEntity.Key)
				continue
			}
			hardwareID := record.(*orc8rModels.GatewayDevice).HardwareID
			if len(hardwareID) > 0 {
				res = append(res, hardwareID)
			} else {
				glog.Errorf("Empty Harware ID for Gateway %s:%s", networkID, gatewayEntity.Key)
			}
		}
	}

	return res, nil
}

func getFegCfg(networkID, gatewayID string) (*models.GatewayFederationConfigs, error) {
	fegGateway, err := configurator.LoadEntity(
		networkID, feg.FegGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil && err != merrors.ErrNotFound {
		return nil, errors.WithStack(err)
	}
	if err == nil && fegGateway.Config != nil {
		return fegGateway.Config.(*models.GatewayFederationConfigs), nil
	}

	iNetworkConfig, err := configurator.LoadNetworkConfig(networkID, feg.FegNetworkType, serdes.Network)
	if err != nil {
		return nil, merrors.ErrNotFound
	}

	gwConfig := iNetworkConfig.(*models.GatewayFederationConfigs)
	return gwConfig, nil
}
