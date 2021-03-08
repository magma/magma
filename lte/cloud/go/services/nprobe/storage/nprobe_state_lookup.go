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

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type NProbeStateLookup interface {
	// Initialize the backing store.
	Initialize() error

	// GetNProbeState returns the NProbeState keyed by networkID and targetID
	GetNProbeState(networkID string, targetID string) (*models.NetworkProbeState, error)

	// SetNProbeState sets current NProbeState for a given networkID and targetID
	SetNProbeState(networkID string, targetID string, state *models.NetworkProbeState) error
}

const (
	tableName      = "lte_nprobe_state"
	nidCol         = "network_id"
	targetIdCol    = "target_id"
	nprobeStateCol = "nprobe_state"
)

type nprobeStateLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewNProbeStateLookup(db *sql.DB, builder sqorc.StatementBuilder) NProbeStateLookup {
	return &nprobeStateLookup{db: db, builder: builder}
}

func (np *nprobeStateLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := np.builder.CreateTable(tableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(targetIdCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(nprobeStateCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(nidCol, targetIdCol).
			Unique(nidCol, targetIdCol, nprobeStateCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize nprobe state lookup table")
	}
	_, err := sqorc.ExecInTx(np.db, nil, nil, txFn)
	return err
}

func (np *nprobeStateLookup) GetNProbeState(networkID string, targetID string) (*models.NetworkProbeState, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := np.builder.
			Select(nprobeStateCol).
			From(tableName).
			Where(squirrel.Eq{nidCol: networkID, targetIdCol: targetID}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "select NProbeState for targetID %v", targetID)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetNProbeState")

		var marshaledState []byte
		for rows.Next() {
			err = rows.Scan(&marshaledState)
			if err != nil {
				return nil, errors.Wrap(err, "select NProbeState for (targetID), SQL row scan error")
			}
			// always return the first record as all cols are unique
			break
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select NProbeState for (targetID), SQL rows error")
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
	state := models.NetworkProbeState{}
	err = state.UnmarshalBinary(ret)
	return &state, err
}

func (np *nprobeStateLookup) SetNProbeState(networkID string, targetID string, state *models.NetworkProbeState) error {

	marshaledState, err := state.MarshalBinary()
	if err != nil {
		return err
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetNProbeState")

		_, err := np.builder.
			Insert(tableName).
			Columns(nidCol, targetIdCol, nprobeStateCol).
			Values(networkID, targetID, marshaledState).
			OnConflict(
				[]sqorc.UpsertValue{{Column: nprobeStateCol, Value: marshaledState}},
				nidCol, targetIdCol,
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
