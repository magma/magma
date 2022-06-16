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
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/sqorc"
)

type SubscriberStorage interface {
	// Initialize the backing store.
	Initialize() error

	GetSubscribersForGateway(networkID string, gatewayID string) ([]SubscriberState, error)

	DeleteSubscribersForGateway(networkID string, gatewayID string) error

	SetAllSubscribersForGateway(networkID string, gatewayID string, subscriberStates []SubscriberState) error
}

const (
	tableName = "gateway_subscriber_states"

	nidColumn         = "network_id"
	gwidColumn        = "gateway_id"
	imsiColumn        = "imsi"
	lastUpdatedColumn = "last_updated_at"
	stateColumn       = "reported_state"
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
		_, err := ss.builder.CreateTable(tableName).
			IfNotExists().
			Column(nidColumn).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(gwidColumn).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(imsiColumn).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(lastUpdatedColumn).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			Column(stateColumn).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(nidColumn, gwidColumn, imsiColumn).
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

type SubscriberState struct {
	Imsi  string
	Value state.ArbitraryJSON
}

func (ss *subscriberStorage) GetSubscribersForGateway(networkID string, gatewayID string) ([]SubscriberState, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := ss.getSubscribers(tx, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("Error getting subscribers for nid / gwid:  %v / %v", networkID, gatewayID)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetSubscribersForGateway")

		var mappings []SubscriberState
		for rows.Next() {
			var rawState []byte
			imsi := ""
			err = rows.Scan(&imsi, &rawState)
			if err != nil {
				return nil, fmt.Errorf("GetSubscribersForGateway, SQL row scan error: %w", err)
			}
			reportedState := make(state.ArbitraryJSON)
			err = json.Unmarshal(rawState, &reportedState)
			if err != nil {
				return nil, fmt.Errorf("GetSubscribersForGateway, error unmarshaling state: %w", err)
			}
			mappings = append(mappings, SubscriberState{Imsi: imsi, Value: reportedState})
		}
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("GetSubscribersForGateway, SQL rows error: %w", err)
		}

		return mappings, nil
	}
	txRet, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.([]SubscriberState)
	return ret, nil

}

func (ss *subscriberStorage) SetAllSubscribersForGateway(networkID string, gatewayID string, subscriberStates []SubscriberState) error {
	timeSec := clock.Now().Unix()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetAllSubscribersForGateway")

		err := ss.deleteAllSubscribers(sc, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("error deleting subscribers before update: %w", err)
		}

		for _, blob := range subscriberStates {
			value, err := blob.Value.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("convert subscriber value %+v: %w", blob, err)
			}

			err = ss.insertSubscriber(sc, networkID, gatewayID, blob.Imsi, timeSec, value)
			if err != nil {
				return nil, fmt.Errorf("insert subscriber %+v: %w", blob, err)
			}
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

func (ss *subscriberStorage) DeleteSubscribersForGateway(networkID string, gatewayID string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "DeleteSubscribersForGateway")

		err := ss.deleteAllSubscribers(sc, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("delete subscribers from network %v, gateway %v: %w", networkID, gatewayID, err)
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

func (ss *subscriberStorage) getSubscribers(tx *sql.Tx, networkID string, gatewayID string) (*sql.Rows, error) {
	rows, err := ss.builder.
		Select(imsiColumn, stateColumn).
		From(tableName).
		Where(squirrel.Eq{nidColumn: networkID, gwidColumn: gatewayID}).
		RunWith(tx).
		Query()
	return rows, err
}

func (ss *subscriberStorage) deleteAllSubscribers(sc *squirrel.StmtCache, networkID string, gatewayID string) error {
	_, err := ss.builder.Delete(tableName).
		Where(squirrel.Eq{nidColumn: networkID, gwidColumn: gatewayID}).
		RunWith(sc).
		Exec()
	return err
}

func (ss *subscriberStorage) insertSubscriber(sc *squirrel.StmtCache, networkID string, gatewayID string, imsi string, timeSec int64, value []byte) error {
	_, err := ss.builder.Insert(tableName).
		Columns(nidColumn, gwidColumn, imsiColumn, lastUpdatedColumn, stateColumn).
		Values(networkID, gatewayID, imsi, timeSec, value).
		RunWith(sc).
		Exec()
	return err
}
