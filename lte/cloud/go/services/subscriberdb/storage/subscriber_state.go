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
	"context"
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

	GetSubscribersForGateway(networkID string, gatewayID string) (*GatewaySubscriberState, error)

	DeleteSubscribersForGateway(networkID string, gatewayID string) error

	SetAllSubscribersForGateway(networkID string, gatewayID string, subscriberStates *GatewaySubscriberState) error

	GetSubscribersForIMSIs(networkID string, imsis []string) (ImsiStateMap, error)
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

// key: IMSI, value: sessiond subscriber state
type ImsiStateMap = map[string]state.ArbitraryJSON

type GatewaySubscriberState struct {
	Subscribers ImsiStateMap `json:"subscribers"`
}

func (j *GatewaySubscriberState) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *GatewaySubscriberState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, j)
}

func (j *GatewaySubscriberState) ValidateModel(context.Context) error {
	return nil
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
		_, err = ss.builder.CreateIndex("network_id_imsi_idx").
			IfNotExists().
			On(tableName).
			Columns(nidColumn, imsiColumn).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("failed to create nid,imsi index: %w", err)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	return err
}

func (ss *subscriberStorage) GetSubscribersForGateway(networkID string, gatewayID string) (*GatewaySubscriberState, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := ss.getSubscribers(tx, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("Error getting subscribers for nid / gwid:  %v / %v", networkID, gatewayID)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetSubscribersForGateway")

		mappings := GatewaySubscriberState{Subscribers: map[string]state.ArbitraryJSON{}}
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
			mappings.Subscribers[imsi] = reportedState
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
	ret := txRet.(GatewaySubscriberState)
	return &ret, nil

}

func (ss *subscriberStorage) SetAllSubscribersForGateway(networkID string, gatewayID string, subscriberStates *GatewaySubscriberState) error {
	timeSec := clock.Now().Unix()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetAllSubscribersForGateway")

		err := ss.deleteAllSubscribers(sc, networkID, gatewayID)
		if err != nil {
			return nil, fmt.Errorf("error deleting subscribers before update: %w", err)
		}

		for imsi, blob := range subscriberStates.Subscribers {
			value, err := blob.MarshalBinary()
			if err != nil {
				return nil, fmt.Errorf("convert subscriber value %+v: %w", blob, err)
			}

			err = ss.insertSubscriber(sc, networkID, gatewayID, imsi, timeSec, value)
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

// GetSubscribersForIMSIs takes a network ID and a list of IMSIs and returns the
// subscriber states of those IMSIs in the network.
// If imsis == nil, states for all IMSIs are returned.
func (ss *subscriberStorage) GetSubscribersForIMSIs(networkID string, imsis []string) (ImsiStateMap, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		var rows *sql.Rows
		var err error
		if imsis == nil {
			rows, err = ss.getSubscribersForNetwork(tx, networkID)
		} else {
			rows, err = ss.getSubscribersForIMSIs(tx, networkID, imsis)
		}
		if err != nil {
			return nil, fmt.Errorf("Error getting subscribers for nid / imsis:  %v / %v", networkID, imsis)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetSubscribersForIMSIs")

		mappings, err := ss.convertRowsToMap(rows)
		if err != nil {
			return nil, fmt.Errorf("GetSubscribersForIMSIs, SQL rows error: %w", err)
		}

		return mappings, nil
	}
	txRet, err := sqorc.ExecInTx(ss.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.(ImsiStateMap)
	return ret, nil
}

func (ss *subscriberStorage) convertRowsToMap(rows *sql.Rows) (ImsiStateMap, error) {
	mappings := ImsiStateMap{}
	lastUpdates := map[string]int64{}
	for rows.Next() {
		var rawState []byte
		imsi := ""
		var lastUpdated int64
		err := rows.Scan(&imsi, &rawState, &lastUpdated)
		if err != nil {
			return nil, fmt.Errorf("GetSubscribersForIMSIs, SQL row scan error: %w", err)
		}
		reportedState := make(state.ArbitraryJSON)
		err = json.Unmarshal(rawState, &reportedState)
		if err != nil {
			return nil, fmt.Errorf("GetSubscribersForIMSIs, error unmarshaling state: %w", err)
		}
		if lastUpdated > lastUpdates[imsi] {
			lastUpdates[imsi] = lastUpdated
			mappings[imsi] = reportedState
		}
	}
	err := rows.Err()
	return mappings, err
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

func (ss *subscriberStorage) getSubscribersForIMSIs(tx *sql.Tx, networkID string, imsis []string) (*sql.Rows, error) {
	rows, err := ss.builder.
		Select(imsiColumn, stateColumn, lastUpdatedColumn).
		From(tableName).
		Where(squirrel.Eq{nidColumn: networkID, imsiColumn: imsis}).
		RunWith(tx).
		Query()
	return rows, err
}

func (ss *subscriberStorage) getSubscribersForNetwork(tx *sql.Tx, networkID string) (*sql.Rows, error) {
	rows, err := ss.builder.
		Select(imsiColumn, stateColumn, lastUpdatedColumn).
		From(tableName).
		Where(squirrel.Eq{nidColumn: networkID}).
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
