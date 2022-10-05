/*
 Copyright 2022 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package storage

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/sqorc"
)

const nid1 = "network_1"
const nid2 = "network_2"

const gwid1 = "snowflake_ID_1"
const gwid2 = "snowflake_ID_2"

func TestSubscriberStorage(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	s := NewSubscriberStorage(db, sqorc.GetSqlBuilder())
	assert.NoError(t, s.Initialize())

	t.Run("Test inserting, querying and deleting subscriber states for a gateway", func(t *testing.T) {
		subscriberStates := CreateTestGatewaySubscriberState("IMSI001010000000123", "IMSI001010000000456")
		err = s.SetAllSubscribersForGateway(nid1, gwid1, &subscriberStates)
		assert.NoError(t, err)

		foundStates, err := s.GetSubscribersForGateway(nid2, gwid1)
		assert.NoError(t, err)
		assert.Empty(t, foundStates.Subscribers)
		foundStates, err = s.GetSubscribersForGateway(nid1, gwid2)
		assert.NoError(t, err)
		assert.Empty(t, foundStates.Subscribers)

		foundStates, err = s.GetSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(*foundStates, subscriberStates))

		err = s.DeleteSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		foundStates, err = s.GetSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		assert.Empty(t, foundStates.Subscribers)
	})

	t.Run("Test updating subscriber states for a gateway", func(t *testing.T) {
		oldSubscriberStatesGw1 := CreateTestGatewaySubscriberState("IMSI001010000000123", "IMSI001010000000456")
		subscriberStatesGw2 := CreateTestGatewaySubscriberState("IMSI001010000000789")

		err = s.SetAllSubscribersForGateway(nid1, gwid1, &oldSubscriberStatesGw1)
		assert.NoError(t, err)
		err = s.SetAllSubscribersForGateway(nid1, gwid2, &subscriberStatesGw2)
		assert.NoError(t, err)

		foundStates, err := s.GetSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(*foundStates, oldSubscriberStatesGw1))

		// Test updating states for gateway 1

		newSubscriberStatesGw1 := CreateTestGatewaySubscriberState("IMSI001010000000012", "IMSI001010000000123")
		err = s.SetAllSubscribersForGateway(nid1, gwid1, &newSubscriberStatesGw1)
		assert.NoError(t, err)

		foundStates, err = s.GetSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(*foundStates, newSubscriberStatesGw1))

		// Test emptying states for gateway 1

		err = s.SetAllSubscribersForGateway(nid1, gwid1, &GatewaySubscriberState{Subscribers: map[string]state.ArbitraryJSON{}})
		assert.NoError(t, err)

		foundStates, err = s.GetSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		assert.Empty(t, foundStates.Subscribers)

		// Gateway 2 should be unaffected

		foundStates, err = s.GetSubscribersForGateway(nid1, gwid2)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(*foundStates, subscriberStatesGw2))

		err = s.DeleteSubscribersForGateway(nid1, gwid1)
		assert.NoError(t, err)
		err = s.DeleteSubscribersForGateway(nid1, gwid2)
		assert.NoError(t, err)
	})

	t.Run("Test getting subscriber states by imsi and getting all subscribers", func(t *testing.T) {
		imsisGw1 := []string{"IMSI001010000000123", "IMSI001010000000456"}
		imsisGw2 := []string{"IMSI002020000000123", "IMSI002020000000456"}
		subscriberStatesGw1 := CreateTestGatewaySubscriberState(imsisGw1...)
		subscriberStatesGw2 := CreateTestGatewaySubscriberState(imsisGw2...)

		var frozenTime int64 = 100000
		err = s.(*subscriberStorage).setTestGatewaySubscriberState(t, nid1, gwid1, &subscriberStatesGw1, frozenTime)
		assert.NoError(t, err)
		err = s.(*subscriberStorage).setTestGatewaySubscriberState(t, nid1, gwid2, &subscriberStatesGw2, frozenTime)
		assert.NoError(t, err)
		err = s.(*subscriberStorage).setTestGatewaySubscriberState(t, nid2, gwid1, &subscriberStatesGw1, frozenTime)
		assert.NoError(t, err)

		// get all IMSIs for network

		foundStates, err := s.GetSubscribersForIMSIs(nid1, nil)
		assert.NoError(t, err)
		var expectedStates = map[string]state.ArbitraryJSON{}
		copyMap(expectedStates, subscriberStatesGw1.Subscribers)
		copyMap(expectedStates, subscriberStatesGw2.Subscribers)
		assert.True(t, cmp.Equal(foundStates, expectedStates))

		// get selected IMSIs for network (from two gateways)

		selectedIMSIs := []string{"IMSI001010000000123", "IMSI002020000000456"}
		expectedStates = map[string]state.ArbitraryJSON{
			"IMSI001010000000123": subscriberStatesGw1.Subscribers["IMSI001010000000123"],
			"IMSI002020000000456": subscriberStatesGw2.Subscribers["IMSI002020000000456"],
		}
		foundStates, err = s.GetSubscribersForIMSIs(nid1, selectedIMSIs)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(foundStates, expectedStates))

		// get latest entry for duplicate IMSI, appearing in two gateways

		imsisGw2 = []string{"IMSI001010000000123", "IMSI002020000000456"}
		subscriberStatesGw2 = CreateTestGatewaySubscriberState(imsisGw2...)
		subscriberStatesGw2.Subscribers["IMSI001010000000123"]["magma.ipv4"].([]interface{})[0].(map[string]interface{})["active_duration_sec"] = float64(40)
		subscriberStatesGw2.Subscribers["IMSI001010000000123"]["magma.ipv4"].([]interface{})[0].(map[string]interface{})["ipv4"] = "10.2.1.12"

		frozenTime = 200000
		err = s.(*subscriberStorage).setTestGatewaySubscriberState(t, nid1, gwid2, &subscriberStatesGw2, frozenTime)
		assert.NoError(t, err)
		foundStates, err = s.GetSubscribersForIMSIs(nid1, selectedIMSIs)
		assert.NoError(t, err)
		assert.True(t, cmp.Equal(foundStates, subscriberStatesGw2.Subscribers))
	})
}

func (ss *subscriberStorage) setTestGatewaySubscriberState(t *testing.T, networkID string, gatewayID string, subscriberStates *GatewaySubscriberState, frozenTime int64) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "setTestGatewaySubscriberState")

		err := ss.deleteAllSubscribers(sc, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("error deleting subscribers before update: %w", err)
		}

		for imsi, blob := range subscriberStates.Subscribers {
			value, err := blob.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("convert subscriber value %+v: %w", blob, err)
			}

			err = ss.insertSubscriber(sc, networkID, gatewayID, imsi, frozenTime, value)
			if err != nil {
				return nil, fmt.Errorf("insert subscriber %+v: %w", blob, err)
			}
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

// copyMap copies contents of src to dst, overwriting
// values to already existing keys with the value from src
func copyMap(dst, src map[string]state.ArbitraryJSON) {
	for k, v := range src {
		dst[k] = v
	}
}
