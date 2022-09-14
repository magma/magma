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
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/sqorc"
)

// AmcManager is supposed to be a library that will replace radio controller
// it is not supposed to be a service but rather an interface to database
// could be implemented in this file as separate struct or combined with cbsd manager
// also its methods are supposed to be used in transaction (they should start a new one)
type AmcManager interface {
	// GetState is equivalent to GetState grpc method
	// it should return list of all feasible cbsd with grants
	// cbsd is considered feasible if and only if
	// - it has no pending requests
	// - one of the following conditions is satisfied
	//	 - it has all necessary parameters to perform sas requests (registration/grant)
	//   - it has some pending db action (e.g. it needs to be deleted)
	GetState(sq.BaseRunner) ([]*DetailedCbsd, error)
	CreateRequest(sq.BaseRunner, *MutableRequest) error
	DeleteCbsd(sq.BaseRunner, *DBCbsd) error
	// UpdateCbsd should replace AcknowledgeCbsdUpdate, AcknowledgeCbsdRelinquish
	// and StoreAvailableFrequencies
	// it should just update cbsd (no need to lock)
	UpdateCbsd(sq.BaseRunner, *DBCbsd, db.FieldMask) error
}

type MutableRequest struct {
	Request       *DBRequest
	DesiredTypeId *DBRequestType
}

// WithinTx is used to call AmcManager function inside single transaction.
func WithinTx[T any](db *sql.DB, f func(tx *sql.Tx) (T, error)) (res T, err error) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return res, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}

		switch err {
		case nil:
			err = tx.Commit()
		default:
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	res, err = f(tx)
	return res, err
}

func NewAmcManager(db *sql.DB, builder sqorc.StatementBuilder, errorChecker sqorc.ErrorChecker, locker sqorc.Locker) *amcManager {
	return &amcManager{
		&dpManager{
			db:           db,
			builder:      builder,
			cache:        &enumCache{cache: map[string]map[string]int64{}},
			errorChecker: errorChecker,
			locker:       locker,
		},
	}
}

type amcManager struct {
	*dpManager
}

// GetState TODO
func (m *amcManager) GetState(sq.BaseRunner) ([]*DetailedCbsd, error) {
	return []*DetailedCbsd{}, nil
}

// CreateRequest inserts given request into the DB.
func (m *amcManager) CreateRequest(tx sq.BaseRunner, data *MutableRequest) error {
	builder := m.builder.RunWith(tx)

	desiredTypeId, err := m.cache.getValue(builder, &DBRequestType{}, data.DesiredTypeId.Name.String)
	if err != nil {
		return err
	}
	data.Request.TypeId = db.MakeInt(desiredTypeId)

	mask := db.NewIncludeMask("type_id", "cbsd_id", "payload")
	_, err = db.NewQuery().WithBuilder(builder).From(data.Request).Insert(mask)
	return err
}

// DeleteCbsd removes given CBSD from the DB.
func (m *amcManager) DeleteCbsd(tx sq.BaseRunner, cbsd *DBCbsd) error {
	builder := m.builder.RunWith(tx)
	where := sq.Eq{"id": cbsd.Id}
	return db.NewQuery().WithBuilder(builder).From(cbsd).Where(where).Delete()
}

// UpdateCbsd TODO
func (m *amcManager) UpdateCbsd(sq.BaseRunner, *DBCbsd, db.FieldMask) error {
	return nil
}
