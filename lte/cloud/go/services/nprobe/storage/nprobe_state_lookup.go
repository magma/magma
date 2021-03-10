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

package storage

import (
	"database/sql"

	"magma/lte/cloud/go/services/nprobe/protos"
	"magma/orc8r/cloud/go/sqorc"
	orc8r_protos "magma/orc8r/lib/go/protos"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type NProbeStateService interface {
	// Initialize the backing store.
	Initialize() error

	// GetNProbeState returns the NProbeState keyed by networkID and taskID
	GetNProbeState(networkID, taskID string) (*protos.NProbeState, error)

	// SetNProbeState sets current NProbeState for a given networkID and taskID
	SetNProbeState(networkID, taskID, targetID string, state *protos.NProbeState) error

	// DeleteNProbeState deletes NProbeState for a given networkID and taskID
	DeleteNProbeState(networkID, taskID string) error
}

const (
	tableName      = "lte_nprobe_state"
	nidCol         = "network_id"
	taskIdCol      = "task_id"
	targetIdCol    = "target_id"
	nprobeStateCol = "nprobe_state"
)

type nprobeStateManager struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewNProbeStateService(db *sql.DB, builder sqorc.StatementBuilder) NProbeStateService {
	return &nprobeStateManager{db: db, builder: builder}
}

func (np *nprobeStateManager) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := np.builder.CreateTable(tableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(taskIdCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(targetIdCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(nprobeStateCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(nidCol, taskIdCol).
			Unique(nidCol, targetIdCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize nprobe state lookup table")
	}
	_, err := sqorc.ExecInTx(np.db, nil, nil, txFn)
	return err
}

func (np *nprobeStateManager) GetNProbeState(networkID, taskID string) (*protos.NProbeState, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := np.builder.
			Select(nprobeStateCol).
			From(tableName).
			Where(squirrel.Eq{nidCol: networkID, taskIdCol: taskID}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "select NProbeState for targetID %v", taskID)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetNProbeState")

		var marshaledState []byte
		for rows.Next() {
			err = rows.Scan(&marshaledState)
			if err != nil {
				return nil, errors.Wrap(err, "select NProbeState for (taskID), SQL row scan error")
			}
			// always return the first record as all cols are unique
			break
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select NProbeState for (taskID), SQL rows error")
		}
		return marshaledState, nil
	}
	txtRet, err := sqorc.ExecInTx(np.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txtRet.([]byte)
	if len(ret) == 0 {
		return nil, errors.Wrap(err, "Not found")
	}
	state := &protos.NProbeState{}
	err = orc8r_protos.Unmarshal(ret, state)
	return state, err
}

func (np *nprobeStateManager) SetNProbeState(networkID, taskID, targetID string, state *protos.NProbeState) error {
	marshaledState, err := orc8r_protos.Marshal(state)
	if err != nil {
		return err
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetNProbeState")

		_, err := np.builder.
			Insert(tableName).
			Columns(nidCol, taskIdCol, targetIdCol, nprobeStateCol).
			Values(networkID, taskID, targetID, marshaledState).
			OnConflict(
				[]sqorc.UpsertValue{{Column: nprobeStateCol, Value: marshaledState}},
				nidCol, taskIdCol, targetIdCol,
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert NProbeState %+v", state)
		}
		return nil, nil
	}
	_, err = sqorc.ExecInTx(np.db, nil, nil, txFn)
	return err
}

func (np *nprobeStateManager) DeleteNProbeState(networkID, taskID string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "DeleteNProbeState")

		_, err := np.builder.
			Delete(tableName).
			Where("network_id = ? AND task_id = ?", networkID, taskID).
			RunWith(sc).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "delete NProbeState %+v", taskID)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(np.db, nil, nil, txFn)
	return err
}
