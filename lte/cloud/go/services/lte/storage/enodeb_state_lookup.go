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

	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type EnodebStateLookup interface {
	// Initialize the backing store.
	Initialize() error

	// GetEnodebState returns the EnodebState keyed by networkID, gatewayID,
	// enodebSN.
	GetEnodebState(networkID string, gatewayID string, enodebSN string) ([]byte, error)

	// SetEnodebState sets current EnodebState for a given networkID,
	// gatewayID, enodebSN.
	SetEnodebState(networkID string, gatewayID string, enodebSN string, serializedEnodebState []byte) error
}

const (
	tableName = "lte_multi_gateway_enodeb_state"

	nidCol      = "network_id"
	gidCol      = "gateway_id"
	enbSnCol    = "enodeb_sn"
	enbStateCol = "enodeb_state"
)

type enodebStateLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewEnodebStateLookup(db *sql.DB, builder sqorc.StatementBuilder) EnodebStateLookup {
	return &enodebStateLookup{db: db, builder: builder}
}

func (l *enodebStateLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(tableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(gidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(enbSnCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(enbStateCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(nidCol, gidCol, enbSnCol).
			Unique(nidCol, gidCol, enbSnCol, enbStateCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize enodeb state lookup table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *enodebStateLookup) GetEnodebState(networkID string, gatewayID string, enodebSN string) ([]byte, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(enbStateCol).
			From(tableName).
			Where(squirrel.Eq{nidCol: networkID, gidCol: gatewayID, enbSnCol: enodebSN}).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "select EnodebState for gatewayID, enodebSN %v, %v", gatewayID, enodebSN)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetEnodebState")

		var serializedState []byte
		for rows.Next() {
			err = rows.Scan(&serializedState)
			if err != nil {
				return nil, errors.Wrap(err, "select EnodebState for (gatewayID, enodebSN), SQL row scan error")
			}
			// always return the first record as all cols are unique
			break
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select EnodebState for (gatewayID, enodebSN), SQL rows error")
		}
		return serializedState, nil
	}
	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.([]byte)
	if len(ret) == 0 {
		return nil, merrors.ErrNotFound
	}
	return ret, nil
}

func (l *enodebStateLookup) SetEnodebState(networkID string, gatewayID string, enodebSN string, enodebState []byte) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		sc := squirrel.NewStmtCache(tx)
		defer sqorc.ClearStatementCacheLogOnError(sc, "SetEnodebState")

		_, err := l.builder.
			Insert(tableName).
			Columns(nidCol, gidCol, enbSnCol, enbStateCol).
			Values(networkID, gatewayID, enodebSN, enodebState).
			OnConflict(
				[]sqorc.UpsertValue{{Column: enbStateCol, Value: enodebState}},
				nidCol, gidCol, enbSnCol,
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert EnodbeState %+v", enodebState)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}
