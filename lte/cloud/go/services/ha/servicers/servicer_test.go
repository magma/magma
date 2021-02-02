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

package servicers_test

import (
	"context"
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/ha/servicers"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	orc8r_protos "magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestHAServicer_GetEnodebOffloadState(t *testing.T) {
	configurator_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	lte_test_init.StartTestService(t)
	servicer := servicers.NewHAServicer()

	testNetworkId := "n1"
	testGwHwId1 := "hw1"
	testGwId1 := "g1"
	testGwHwId2 := "hw2"
	testGwId2 := "g2"
	testGwPool := "pool1"
	enbSn := "enb1"
	err := configurator.CreateNetwork(configurator.Network{ID: testNetworkId}, serdes.Network)
	assert.NoError(t, err)

	// Initialize HA network topology
	_, err = configurator.CreateEntity(
		testNetworkId,
		configurator.NetworkEntity{
			Type:   lte.CellularEnodebEntityType,
			Key:    enbSn,
			Config: newDefaultUnmanagedEnodebConfig(),
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		testNetworkId,
		[]configurator.NetworkEntity{
			{
				Type: lte.CellularGatewayEntityType, Key: testGwId1,
				Config:       newDefaultGatewayConfig(1, 255),
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enbSn}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwId1,
				Name: "foobar", Description: "foo bar",
				PhysicalID:   testGwHwId1,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: testGwId1}},
			},
			{
				Type: lte.CellularGatewayEntityType, Key: testGwId2,
				Config:       newDefaultGatewayConfig(2, 1),
				Associations: []storage.TypeAndKey{{Type: lte.CellularEnodebEntityType, Key: enbSn}},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwId2,
				Name: "foobar2", Description: "foo bar",
				PhysicalID:   testGwHwId2,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: testGwId2}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		testNetworkId,
		[]configurator.NetworkEntity{
			{
				Type: lte.CellularGatewayPoolEntityType, Key: testGwPool,
				Config: &lte_models.CellularGatewayPoolConfigs{
					MmeGroupID: 1,
				},
				Associations: []storage.TypeAndKey{
					{Type: lte.CellularGatewayEntityType, Key: testGwId1},
					{Type: lte.CellularGatewayEntityType, Key: testGwId2},
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Initialize network state for given devices
	gwStatus := &models.GatewayStatus{
		CheckinTime: uint64(time.Now().Unix()),
		HardwareID:  testGwHwId1,
	}
	ctx1 := test_utils.GetContextWithCertificate(t, testGwHwId1)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)
	gwStatus2 := &models.GatewayStatus{
		CheckinTime: uint64(time.Now().Unix()),
		HardwareID:  testGwHwId2,
	}
	ctx2 := test_utils.GetContextWithCertificate(t, testGwHwId2)
	test_utils.ReportGatewayStatus(t, ctx2, gwStatus2)

	enbState := getDefaultEnodebState(testGwId1)
	reportEnodebState(t, testNetworkId, testGwId1, enbSn, enbState)

	// First test primary healthy
	ctx := orc8r_protos.NewGatewayIdentity(testGwHwId2, testNetworkId, testGwId2).NewContextWithIdentity(context.Background())
	res, err := servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes := &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED_AND_SERVING_UES,
		},
	}
	assert.Equal(t, expectedRes, res)

	// Now simulate failed primary not checking in
	stateTooOld := time.Now().Add(-time.Second * 600)
	clock.SetAndFreezeClock(t, stateTooOld)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{},
	}
	assert.Equal(t, expectedRes, res)
	clock.UnfreezeClock(t)
	test_utils.ReportGatewayStatus(t, ctx1, gwStatus)

	// Simulate too old of ENB state
	clock.SetAndFreezeClock(t, stateTooOld)
	reportEnodebState(t, testNetworkId, testGwId1, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_NO_OP,
		},
	}
	assert.Equal(t, expectedRes, res)
	clock.UnfreezeClock(t)

	// ENB not connected
	enbState.EnodebConnected = swag.Bool(false)
	reportEnodebState(t, testNetworkId, testGwId1, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_NO_OP,
		},
	}
	assert.Equal(t, expectedRes, res)

	// Connected but no UEs connected
	enbState.EnodebConnected = swag.Bool(true)
	enbState.UesConnected = 0
	reportEnodebState(t, testNetworkId, testGwId1, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED,
		},
	}
	assert.Equal(t, expectedRes, res)

	// Back to connected with users
	enbState.UesConnected = 10
	reportEnodebState(t, testNetworkId, testGwId1, enbSn, enbState)

	// Simulate secondary gateway updating state to ensure primary state
	// can still be fetched
	reportEnodebState(t, testNetworkId, testGwId2, enbSn, enbState)

	res, err = servicer.GetEnodebOffloadState(ctx, &protos.GetEnodebOffloadStateRequest{})
	assert.NoError(t, err)
	expectedRes = &protos.GetEnodebOffloadStateResponse{
		EnodebOffloadStates: map[uint32]protos.GetEnodebOffloadStateResponse_EnodebOffloadState{
			138: protos.GetEnodebOffloadStateResponse_PRIMARY_CONNECTED_AND_SERVING_UES,
		},
	}
	assert.Equal(t, expectedRes, res)
}

func reportEnodebState(t *testing.T, networkID string, gatewayID string, enodebSerial string, req *lte_models.EnodebState) {
	req.TimeReported = uint64(clock.Now().UnixNano()) / uint64(time.Millisecond)
	serializedEnodebState, err := serde.Serialize(req, lte.EnodebStateType, serdes.State)
	assert.NoError(t, err)
	err = lte_service.SetEnodebState(networkID, gatewayID, enodebSerial, serializedEnodebState)
	assert.NoError(t, err)
}

func newDefaultGatewayConfig(mmeCode uint32, mmeRelCap uint32) *lte_models.GatewayCellularConfigs {
	return &lte_models.GatewayCellularConfigs{
		Ran: &lte_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lte_models.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &lte_models.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              nil,
			NonEpsServiceControl: swag.Uint32(0),
		},
		DNS: &lte_models.GatewayDNSConfigs{
			DhcpServerEnabled: swag.Bool(true),
			EnableCaching:     swag.Bool(false),
			LocalTTL:          swag.Int32(0),
		},
		HeConfig: &lte_models.GatewayHeConfig{},
		Pooling: lte_models.CellularGatewayPoolRecords{
			{
				GatewayPoolID:       "pool1",
				MmeCode:             mmeCode,
				MmeRelativeCapacity: mmeRelCap,
			},
		},
	}
}

func newDefaultUnmanagedEnodebConfig() *lte_models.EnodebConfig {
	ip := strfmt.IPv4("192.168.0.124")
	return &lte_models.EnodebConfig{
		ConfigType: "UNMANAGED",
		UnmanagedConfig: &lte_models.UnmanagedEnodebConfiguration{
			CellID:    swag.Uint32(138),
			Tac:       swag.Uint32(1),
			IPAddress: &ip,
		},
	}
}

func getDefaultEnodebState(gwID string) *lte_models.EnodebState {
	return &lte_models.EnodebState{
		MmeConnected:       swag.Bool(true),
		EnodebConnected:    swag.Bool(true),
		IPAddress:          "10.0.0.1",
		ReportingGatewayID: gwID,
		EnodebConfigured:   swag.Bool(true),
		GpsConnected:       swag.Bool(true),
		GpsLatitude:        swag.String("foo"),
		GpsLongitude:       swag.String("bar"),
		OpstateEnabled:     swag.Bool(true),
		PtpConnected:       swag.Bool(true),
		RfTxOn:             swag.Bool(true),
		RfTxDesired:        swag.Bool(true),
		FsmState:           swag.String("abc"),
		UesConnected:       5,
	}
}
