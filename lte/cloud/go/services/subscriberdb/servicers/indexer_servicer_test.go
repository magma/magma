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
	"encoding/base64"
	"net"
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/subscriberdb"
	subscriberdb_test_init "magma/lte/cloud/go/services/subscriberdb/test_init"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/stretchr/testify/assert"
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
		id00: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}})},
		id01: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}})},
		id10: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.1")}})},
		id11: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.2")}})},
	}

	// Index the imsi0->sid0 state, result is sid0->imsi0 reverse mapping
	errs, err := idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Empty(t, errs)
	gotA, err := subscriberdb.GetIMSIsForIP("nid0", "127.0.0.1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0", "IMSI1"}, gotA)
	gotB, err := subscriberdb.GetIMSIsForIP("nid0", "127.0.0.2")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI1"}, gotB)

	// Correctly handle per-state errs
	states = state_types.SerializedStatesByID{
		id00: {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": encodeIP("127.0.0.3")}})},
		id2:  {SerializedReportedState: serialize(t, &state.ArbitraryJSON{"ip": state.ArbitraryJSON{"address": "deadbeef"}})},
	}
	errs, err = idx.Index("nid0", states)
	assert.NoError(t, err)
	assert.Error(t, errs[id2])
	gotC, err := subscriberdb.GetIMSIsForIP("nid0", "127.0.0.3")
	assert.NoError(t, err)
	assert.Equal(t, []string{"IMSI0"}, gotC)
}

func serialize(t *testing.T, mobilitydState *state.ArbitraryJSON) []byte {
	bytes, err := serde.Serialize(mobilitydState, lte.MobilitydStateType, serdes.State)
	assert.NoError(t, err)
	return bytes
}

func encodeIP(ip string) string {
	ipBytes := net.ParseIP(ip)[12:16] // get just the IPv4 bytes
	return base64.StdEncoding.EncodeToString(ipBytes)
}
