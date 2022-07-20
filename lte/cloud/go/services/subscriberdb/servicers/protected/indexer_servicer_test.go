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
	"encoding/base64"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	subscriberdb_test_init "magma/lte/cloud/go/services/subscriberdb/test_init"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"
)

func TestIndexerIP(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go
	)
	var (
		types = []string{lte.MobilitydStateType} // copied from indexer_servicer.go
	)

	subscriberdb_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(subscriberdb.ServiceName, version, types...)

	id00 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI0.apn0"}
	id01 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI0.apn1"}
	id10 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI1.apn0"}
	id11 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI1.apn1"}
	id2 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI2.apn2"}

	// Reported state takes the form
	// {
	//   "state": "ALLOCATED",
	//   "sid": {"id": "IMSI001010000000001.magma.ipv4"},
	//   "ipBlock": {"netAddress": "wKiAAA==", "prefixLen": 24},
	//   "ip": {"address": "wKiArg=="}
	//  }
	states := state_types.SerializedStatesByID{
		id00: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}}, lte.MobilitydStateType)},
		id01: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}}, lte.MobilitydStateType)},
		id10: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}}, lte.MobilitydStateType)},
		id11: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.2")}}, lte.MobilitydStateType)},
	}

	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	gotA, err := subscriberdb.GetIMSIsForIP(context.Background(), "nid0", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0", "IMSI1"}, gotA)
	gotB, err := subscriberdb.GetIMSIsForIP(context.Background(), "nid0", "127.0.0.2")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI1"}, gotB)

	// Correctly handle per-state errs
	states = state_types.SerializedStatesByID{
		id00: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.3")}}, lte.MobilitydStateType)},
		id2:  {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": "deadbeef"}}, lte.MobilitydStateType)},
	}
	errs, err = idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Error(t, errs[id2])
	gotC, err := subscriberdb.GetIMSIsForIP(context.Background(), "nid0", "127.0.0.3")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0"}, gotC)
}

func serialize(t *testing.T, state serde.ValidateableBinaryConvertible, stateType string) []byte {
	bytes, err := serde.Serialize(state, stateType, serdes.State)
	assert.NoError(t, err)
	return bytes
}

func encodeIP(ip string) string {
	ipBytes := net.ParseIP(ip)[12:16] // get just the IPv4 bytes
	return base64.StdEncoding.EncodeToString(ipBytes)
}

const gwid1 = "snowflake_ID_1"
const gwid2 = "snowflake_ID_2"

func TestIndexerSubscriber(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go
	)

	var (
		types = []string{lte.GatewaySubscriberStateType}
	)

	subscriberStorage := subscriberdb_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(subscriberdb.ServiceName, version, types...)

	// define IDs
	id0 := state_types.ID{Type: lte.GatewaySubscriberStateType, DeviceID: gwid1}
	id1 := state_types.ID{Type: lte.GatewaySubscriberStateType, DeviceID: gwid2}
	// define gateway states
	gatewayState0 := storage.CreateTestGatewaySubscriberState("IMSI001010000000123", "IMSI001010000000456", "IMSI001010000000789", "IMSI001010000001011")
	gatewayState1 := storage.CreateTestGatewaySubscriberState("IMSI002020000000123", "IMSI002020000000456", "IMSI002020000000789")

	// empty gateway state
	emptyState := storage.GatewaySubscriberState{Subscribers: map[string]state.ArbitraryJSON{}}
	emptySerializedState := state_types.SerializedStatesByID{
		id0: {SerializedReportedState: serialize(t, &emptyState, lte.GatewaySubscriberStateType)},
	}
	errs, err := idx.Index("nid0", emptySerializedState)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err := subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, emptyState))

	// one gateway
	serializedGatewayState := state_types.SerializedStatesByID{
		id0: {SerializedReportedState: serialize(t, &gatewayState0, lte.GatewaySubscriberStateType)},
	}
	errs, err = idx.Index("nid0", serializedGatewayState)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, gatewayState0))

	// two gatewayIDs
	serializedGatewayState = state_types.SerializedStatesByID{
		id0: {SerializedReportedState: serialize(t, &gatewayState1, lte.GatewaySubscriberStateType)},
		id1: {SerializedReportedState: serialize(t, &gatewayState0, lte.GatewaySubscriberStateType)},
	}
	errs, err = idx.Index("nid0", serializedGatewayState)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, gatewayState1))
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid2)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, gatewayState0))

	// empty gateways
	emptySerializedState = state_types.SerializedStatesByID{
		id0: {SerializedReportedState: serialize(t, &emptyState, lte.GatewaySubscriberStateType)},
		id1: {SerializedReportedState: serialize(t, &emptyState, lte.GatewaySubscriberStateType)},
	}
	errs, err = idx.Index("nid0", emptySerializedState)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, emptyState))
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid2)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, emptyState))
}

func TestIndexerSubscriberAndIP(t *testing.T) {
	const (
		version indexer.Version = 1 // copied from indexer_servicer.go
	)

	var (
		types = []string{lte.GatewaySubscriberStateType, lte.MobilitydStateType}
	)

	subscriberStorage := subscriberdb_test_init.StartTestService(t)
	idx := indexer.NewRemoteIndexer(subscriberdb.ServiceName, version, types...)

	// define IDs
	id0 := state_types.ID{Type: lte.GatewaySubscriberStateType, DeviceID: gwid1}
	id1 := state_types.ID{Type: lte.GatewaySubscriberStateType, DeviceID: gwid2}
	// define gateway states
	gatewayState0 := storage.CreateTestGatewaySubscriberState("IMSI001010000000123", "IMSI001010000000456", "IMSI001010000000789", "IMSI001010000001011")
	gatewayState1 := storage.CreateTestGatewaySubscriberState("IMSI002020000000123", "IMSI002020000000456", "IMSI002020000000789")

	// define mobilityd state ids
	idm0 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI0.apn0"}
	idm1 := state_types.ID{Type: lte.MobilitydStateType, DeviceID: "IMSI0.apn1"}

	// mixed state
	states := state_types.SerializedStatesByID{
		id0:  {SerializedReportedState: serialize(t, &gatewayState0, lte.GatewaySubscriberStateType)},
		idm0: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}}, lte.MobilitydStateType)},
		idm1: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}}, lte.MobilitydStateType)},
		id1:  {SerializedReportedState: serialize(t, &gatewayState1, lte.GatewaySubscriberStateType)},
	}
	errs, err := idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err := subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, gatewayState0))
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid2)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, gatewayState1))
	foundIPs, err := subscriberdb.GetIMSIsForIP(context.Background(), "nid0", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0"}, foundIPs)

	// create empty gateway states
	emptyState := storage.GatewaySubscriberState{Subscribers: map[string]state.ArbitraryJSON{}}
	emptySerializedState := state_types.SerializedStatesByID{
		id0: {SerializedReportedState: serialize(t, &emptyState, lte.GatewaySubscriberStateType)},
		id1: {SerializedReportedState: serialize(t, &emptyState, lte.GatewaySubscriberStateType)},
	}
	// empty gateway states
	errs, err = idx.Index("nid0", emptySerializedState)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid1)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, emptyState))
	foundStates, err = subscriberStorage.GetSubscribersForGateway("nid0", gwid2)
	assert.NoError(t, err)
	assert.True(t, cmp.Equal(*foundStates, emptyState))
	foundIPs, err = subscriberdb.GetIMSIsForIP(context.Background(), "nid0", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0"}, foundIPs)
}
