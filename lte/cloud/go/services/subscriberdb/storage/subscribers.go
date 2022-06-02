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
	"strconv"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
)


type SubscriberStorage interface {
	// Initialize the backing store.
	Initialize() error

	// TODO Which further methods do we need?
	// + General design question: Do we want to have the same interface structure
	// as blobstore.go (one to initialize the table and start a transaction and
	// one interface for transactions) or like in ip_lookup.go where we only have
	// one interface that is capable of all that is needed?

	// Suggestions:
	SetSubscribers(networkID string, gatewayID string, blobs []Blob) error
	DeleteSubscribers(networkID string, gatewayID string) error

	// GetSubscriber(networkID string, gatewayID string, IMSI string)
	// GetAllSubscribersFromGateway(networkID string, gatewayID string)
	// GetAllSubscribersFromNetwork(network ID)
	// one or more of these methods are probably nice to have for automated tests
}


// Sketch of the proposed solution

// As before gateway uses ReportStates method (orc8r/cloud/go/services/state/servicers/southbound/servicer.go) of state service
// to send its state to the cloud.
// Instead of the current `ReportStatesRequest` proto message it might need to use a new one which encapsulates
// the list of `State` proto messages (defined in orc8r/protos/service303.proto) in a list again. Or can we actually still use the old message
// and deal with the empty list case? More details in [1]

// State service dumps the agw state in the `states` table. A reindexer that is registered somewhere will read the state object and write into the
// new subscriber storage table that is maintained by subscriberdb. This table is defined in https://confluence.tngtech.com/display/PFM/Subscriber+state+reindexing
//
// The table needs an interface which has at least Set and Delete methods. See the interface above for more details.
// For the Set method we need a type in which the data is passed to the method (analogous to Blobs / Blob from blobstore package in magma).
// Encoding and decoding needs to be taken care of in the process of writing the data.
// (!) Specify the output format of the Value entry and write a test for that.

// The reindexer: don't have much to say about this one yet.

// [1] proto message TODO How should the new message be defined? To answer that we probably have to look at the Sessiond side first?
// The requirement is, that the AGW needs to be able to send an empty list. Apparently this is not possible currently, an empty list is not sent at all.
// So we need the list to be an entry in another list. But does that mean that we need a new proto message at all or can we instead parse the new
// sessiond message into the old `ReportStatesRequest` and deal with an empty message case?
// If not, something like this might work (didn't check if that is valid syntax, just to illustrate the idea)
// message ReportStatesRequest {
// 	repeated repeated State states = 1;
// }

const (
	subscriberStorageTableName = "gateway_subscriber_states"

	subscriberStorageNidColumn = "network_id"
	subscriberStorageGwidCol = "gateway_id"
	subscriberStorageImsiCol = "IMSI"
	subscriberStorageTimestamp = "timestamp_in_milliseconds"
	subscriberStorageState = "reported_state"
)

type subscriberStorage struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewSubscriberStorage(db *sql.DB, builder sqorc.StatementBuilder) SubscriberStorage {
	return &subscriberStorage{
		db:      db,
		builder: builder,
	}
}

func (ss *subscriberStorage) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := ss.builder.CreateTable(subscriberStorageTableName).
			IfNotExists().
			Column(subscriberStorageNidColumn).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberStorageGwidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberStorageImsiCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberStorageTimestamp).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberStorageState).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(subscriberStorageNidColumn, subscriberStorageGwidCol, subscriberStorageImsiCol).
			Unique(subscriberStorageNidColumn, subscriberStorageGwidCol, subscriberStorageImsiCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("initialize subscriber storage table: %w", err)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}
// TODO How to use Unique and PrimaryKey?
// Primary Key implies uniqueness. Has to hold for the combination of all keys.
// Uniqueness: probably also has to hold for the combination of the values.

type Blob struct {
	Key     string
	Value   []byte
	Version uint64
}
// TODO do we need Version? There is a version in the value string as well. But they seem to differ.


func (ss *subscriberStorage) SetSubscribers(networkID string, gatewayID string, blobs []Blob) error {
	timeSec := int(clock.Now().UnixNano()) / int(time.Nanosecond)

	txFn := func(tx *sql.Tx) (interface{}, error) {
		// why do they use a statement cache here in ip_lookup.go SetIPs but not in the create table case?
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetSubscribers")

		for _, blob := range blobs {
			_, err := ss.builder.Insert(subscriberStorageTableName).
				Columns(subscriberStorageNidColumn, subscriberStorageGwidCol, subscriberStorageImsiCol, subscriberStorageTimestamp, subscriberStorageState).
				Values(networkID, gatewayID, blob.Key, strconv.Itoa(timeSec), blob.Value).
				// what should we do OnConflict? (currently we are doing nothing, i.e. not writing the value) (insert something like {Column: subscriberStorageImsiCol, Value: blob.Key} as UpsertValue?)
				OnConflict([]sqorc.UpsertValue{}, subscriberStorageNidColumn, subscriberStorageGwidCol, subscriberStorageImsiCol).
				RunWith(sc).
				Exec()
			if err != nil {
				return nil, fmt.Errorf("insert subscriber %+v: %w", blob, err)
			}
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

// Not tested yet.
func (ss *subscriberStorage) DeleteSubscribers(networkID string, gatewayID string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "DeleteSubscribers")

		_, err := ss.builder.Delete(subscriberStorageTableName).
			Where(map[string]string{subscriberStorageNidColumn: networkID, subscriberStorageGwidCol: gatewayID}).
			RunWith(sc).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("delete subscribers from network %v, gateway %v: %w", networkID, gatewayID, err)
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

// TODO Open questions:

// 1) Which write conflicts can appear and what to do then?
// Guess it could happen that we get the same IMSI twice for different gateways.
// but for different networks this shouldn't be possible? (A sim card can only be logged into one network, IMSI is bound to the sim card, not the phone (?))

// 3) What is the statement cache doing? is it necessary?

// 2) pass gateway and network ID as an argument or get them within the function from a passed context?
// assumption here is that the reindexer calls SetSubscribers and passes the IDs,
// the reindexer in turn could get them by calling a function with the code below

// 	// Get gateway information from context
//	gw := protos.GetClientGateway(ctx)
//	if gw == nil {
//		return nil, state.ErrMissingGateway
//	}
//	if !gw.Registered() {
//		return nil, state.ErrGatewayNotRegistered
//	}
//	hwID := gw.HardwareId
//	networkID := gw.NetworkId
