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
	UpdateCbsd(sq.BaseRunner, *DBCbsd, db.FieldMask) error
}

type MutableRequest struct {
	Request     *DBRequest
	RequestType *DBRequestType
}

// WithinTx is used to call AmcManager function inside single transaction.
func WithinTx[T any](db *sql.DB, f func(tx *sql.Tx) (T, error)) (res T, err error) {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return res, err
	}

	defer func() {
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

// CreateRequest inserts given request into the DB.
func (m *amcManager) CreateRequest(tx sq.BaseRunner, data *MutableRequest) error {
	builder := m.builder.RunWith(tx)

	desiredTypeId, err := m.cache.getValue(builder, &DBRequestType{}, data.RequestType.Name.String)
	if err != nil {
		return err
	}
	data.Request.TypeId = db.MakeInt(desiredTypeId)

	_, err = db.NewQuery().
		WithBuilder(builder).
		From(data.Request).
		Insert(db.NewIncludeMask("type_id", "cbsd_id", "payload"))
	return err
}

// DeleteCbsd removes given CBSD from the DB.
func (m *amcManager) DeleteCbsd(tx sq.BaseRunner, cbsd *DBCbsd) error {
	builder := m.builder.RunWith(tx)
	where := sq.Eq{"id": cbsd.Id}
	return db.NewQuery().
		WithBuilder(builder).
		From(cbsd).
		Where(where).
		Delete()
}

// UpdateCbsd update CBSD in the DB with given mask.
func (m *amcManager) UpdateCbsd(tx sq.BaseRunner, cbsd *DBCbsd, mask db.FieldMask) error {
	builder := m.builder.RunWith(tx)
	_, err := db.NewQuery().
		WithBuilder(builder).
		From(cbsd).
		Select(db.NewIncludeMask()).
		Where(sq.Eq{"id": cbsd.Id}).
		Update(mask)
	return err
}

func (m *amcManager) GetState(tx sq.BaseRunner) ([]*DetailedCbsd, error) {
	runner := m.getQueryRunner(tx)
	return runner.getState()
}

func notNull(fields ...string) sq.Sqlizer {
	filters := make(sq.And, len(fields))
	for i, f := range fields {
		filters[i] = sq.NotEq{f: nil}
	}
	return filters
}

func (r *queryRunner) getState() ([]*DetailedCbsd, error) {
	multiStepFields := []string{"fcc_id", "user_id", "number_of_ports", "min_power", "max_power", "antenna_gain"}
	singleStepFields := append(multiStepFields, "latitude_deg", "longitude_deg", "height_m", "height_type")
	res, err := db.NewQuery().
		WithBuilder(r.builder).
		From(&DBCbsd{}).
		Select(db.NewExcludeMask("network_id", "state_id", "desired_state_id")).
		Join(db.NewQuery().
			From(&DBCbsdState{}).
			As("t1").
			On(db.On(CbsdTable, "state_id", "t1", "id")).
			Select(db.NewIncludeMask("name"))).
		Join(db.NewQuery().
			From(&DBCbsdState{}).
			As("t2").
			On(db.On(CbsdTable, "desired_state_id", "t2", "id")).
			Select(db.NewIncludeMask("name"))).
		Join(db.NewQuery().
			From(&DBGrant{}).
			On(db.On(CbsdTable, "id", GrantTable, "cbsd_id")).
			Select(db.NewIncludeMask("grant_id", "heartbeat_interval", "last_heartbeat_request_time", "low_frequency", "high_frequency")).
			Join(db.NewQuery().
				From(&DBGrantState{}).
				On(db.On(GrantTable, "state_id", GrantStateTable, "id")).
				Select(db.NewIncludeMask("name"))).
			Nullable()).
		Nullable().
		Join(db.NewQuery().
			From(&DBRequest{}).
			On(db.On(CbsdTable, "id", RequestTable, "cbsd_id")).
			Select(db.NewIncludeMask()).
			Join(db.NewQuery().
				From(&DBRequestType{}).
				On(db.On(RequestTable, "type_id", RequestTypeTable, "id")).
				Select(db.NewIncludeMask())).
			Nullable()).
		Nullable().
		Where(sq.And{
			sq.Eq{RequestTable + ".id": nil},
			sq.Or{
				sq.Eq{CbsdTable + ".should_deregister": true},
				sq.Eq{"should_relinquish": true},
				sq.Eq{"is_deleted": true},
				sq.And{
					sq.Eq{"single_step_enabled": false},
					notNull(multiStepFields...),
				},
				sq.And{
					sq.Eq{"single_step_enabled": true},
					sq.Eq{"cbsd_category": "a"},
					sq.Eq{"indoor_deployment": true},
					notNull(singleStepFields...),
				},
			},
		}).
		OrderBy(CbsdTable+".id", db.OrderAsc).
		List()

	if err != nil {
		return nil, err
	}

	cbsds := make([]*DetailedCbsd, 0, len(res))

	cbsdId := int64(-1)
	cbsdIndex := -1

	for i, models := range res {
		cbsd := models[0].(*DBCbsd)
		if cbsd.Id.Int64 != cbsdId {
			cbsds = append(cbsds, convertCbsdToDetails(models))
			cbsdIndex += 1
		}
		grant := res[i][3]
		if grant.(*DBGrant).LowFrequencyHz.Valid {
			grantState := res[i][4]
			cbsds[cbsdIndex].Grants = append(cbsds[cbsdIndex].Grants, &DetailedGrant{
				Grant:      grant.(*DBGrant),
				GrantState: grantState.(*DBGrantState),
			})
		}
		cbsdId = cbsd.Id.Int64
	}

	return cbsds, nil
}
