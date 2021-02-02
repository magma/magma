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
package tests

import (
	"context"
	"sync"
	"testing"

	healthTestInit "magma/feg/cloud/go/services/health/test_init"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	models2 "magma/feg/cloud/go/services/feg/obsidian/models"
	health_servicers "magma/feg/cloud/go/services/health/servicers"
	healthTestUtils "magma/feg/cloud/go/services/health/test_utils"
	"magma/lte/cloud/go/lte"
	models3 "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

var (
	nhNetworkID           = "nh"
	servingFegNetworkID   = "serving_feg"
	federatedLteNetworkID = "federated_lte"
	nhImsi                = "123456000000101"
	nhPlmnId              = nhImsi[:6]
	agwHwId               = "lte_gw_hw_id"
	agwId                 = "lte_gw_id"
	fegHwId               = "feg_hw_id"
	fegId                 = "feg_id"
	nonNhImsi             = "654321000000102"
	nonNhPlmnId           = nonNhImsi[:6]
	nhFegHwId             = "nh_feg_hw_id"
	nhFegId               = "nh_feg_id"

	once sync.Once
)

func setupNeutralHostNetworks(t *testing.T) *health_servicers.TestHealthServer {
	once.Do(func() {
		configuratorTestInit.StartTestService(t)
		deviceTestInit.StartTestService(t)
	})
	testHealthServicer, err := healthTestInit.StartTestService(t)
	assert.NoError(t, err)

	nhFegCfg := models2.NewDefaultNetworkFederationConfigs()
	nhFegCfg.NhRoutes = models2.NhRoutes{nhPlmnId: servingFegNetworkID}
	nhFegCfg.ServedNetworkIds = models2.ServedNetworkIds{federatedLteNetworkID}

	servingFegCfg := models2.NewDefaultNetworkFederationConfigs()
	servingFegCfg.ServedNhIds = models2.ServedNhIds{nhNetworkID}

	lteNetCfg := models2.NewDefaultFederatedNetworkConfigs()
	lteNetCfg.FegNetworkID = &nhNetworkID

	// Neutral Host Network
	nhNetworkConfig := configurator.Network{
		ID:          nhNetworkID,
		Type:        feg.FederationNetworkType,
		Name:        "TestNeutralHost",
		Description: "Test Neutral Host",
		Configs: map[string]interface{}{
			feg.FegNetworkType:          nhFegCfg,
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
		},
	}
	// Serving FeG Network
	servingFegNetworkCfg := configurator.Network{
		ID:          servingFegNetworkID,
		Type:        feg.FederationNetworkType,
		Name:        "serving_feg_network",
		Description: "Serving FeG Network",
		Configs: map[string]interface{}{
			feg.FegNetworkType:          servingFegCfg,
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
		},
	}
	// Federated LTE Network
	federatedLteNetCfg := configurator.Network{
		ID:          federatedLteNetworkID,
		Type:        feg.FederatedLteNetworkType,
		Name:        "Federated_FeG_Network",
		Description: "Federated FeG Network",
		Configs: map[string]interface{}{
			feg.FederatedNetworkType:      lteNetCfg,
			lte.CellularNetworkConfigType: models3.NewDefaultTDDNetworkConfig(),
			orc8r.NetworkFeaturesConfig:   models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:         models.NewDefaultDNSConfig(),
		},
	}
	networkConfigs := []configurator.Network{
		nhNetworkConfig,
		servingFegNetworkCfg,
		federatedLteNetCfg,
	}
	_, err = configurator.CreateNetworks(networkConfigs, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		federatedLteNetworkID,
		[]configurator.NetworkEntity{
			{Type: lte.CellularEnodebEntityType, Key: "enb1"},
			{Type: lte.CellularEnodebEntityType, Key: "enb2"},
			{
				Type: lte.CellularGatewayEntityType, Key: agwId,
				Config: &models3.GatewayCellularConfigs{
					Epc: &models3.GatewayEpcConfigs{NatEnabled: swag.Bool(true), IPBlock: "192.168.0.0/24"},
					Ran: &models3.GatewayRanConfigs{Pci: 260, TransmitEnabled: swag.Bool(true)},
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularEnodebEntityType, Key: "enb1"},
					{Type: lte.CellularEnodebEntityType, Key: "enb2"},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: agwId,
				Name: "lte_gateway", Description: "federated lte gateway",
				PhysicalID: agwHwId,
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: agwId}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: agwId},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		federatedLteNetworkID, orc8r.AccessGatewayRecordType, agwHwId,
		&models.GatewayDevice{HardwareID: agwHwId, Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		servingFegNetworkID,
		[]configurator.NetworkEntity{
			{
				Type: feg.FegGatewayType, Key: fegId,
			},
			{
				Type: orc8r.MagmadGatewayType, Key: fegId,
				Name: "feg_gateway", Description: "federation gateway",
				PhysicalID: fegHwId,
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
				Associations: []storage.TypeAndKey{{Type: feg.FegGatewayType, Key: fegId}},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: fegId},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		servingFegNetworkID, orc8r.AccessGatewayRecordType, fegHwId,
		&models.GatewayDevice{HardwareID: fegHwId, Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	actualNHNet, err := configurator.LoadNetwork(nhNetworkID, true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, nhNetworkConfig, actualNHNet)

	actualFeGNet, err := configurator.LoadNetwork(servingFegNetworkID, true, true, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, servingFegNetworkCfg, actualFeGNet)

	// Update Serving FeG Health status
	ctx := protos.NewGatewayIdentity(fegHwId, servingFegNetworkID, fegId).NewContextWithIdentity(context.Background())
	req := healthTestUtils.GetHealthyRequest()
	_, err = testHealthServicer.UpdateHealth(ctx, req)
	assert.NoError(t, err)

	return testHealthServicer
}
