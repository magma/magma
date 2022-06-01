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

	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/merrors"
)

type CbsdManager interface {
	CreateCbsd(networkId string, data *MutableCbsd) error
	UpdateCbsd(networkId string, id int64, data *MutableCbsd) error
	EnodebdUpdateCbsd(data *DBCbsd) error
	DeleteCbsd(networkId string, id int64) error
	FetchCbsd(networkId string, id int64) (*DetailedCbsd, error)
	ListCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error)
	DeregisterCbsd(networkId string, id int64) error
}

type CbsdFilter struct {
	SerialNumber string
}

type DetailedCbsdList struct {
	Cbsds []*DetailedCbsd
	Count int64
}

type MutableCbsd struct {
	Cbsd         *DBCbsd
	DesiredState *DBCbsdState
}

type DetailedCbsd struct {
	Cbsd         *DBCbsd
	CbsdState    *DBCbsdState
	DesiredState *DBCbsdState
	Grant        *DBGrant
	GrantState   *DBGrantState
}

func NewCbsdManager(db *sql.DB, builder sqorc.StatementBuilder, errorChecker sqorc.ErrorChecker, locker sqorc.Locker) *cbsdManager {
	return &cbsdManager{
		db:           db,
		builder:      builder,
		cache:        &enumCache{cache: map[string]map[string]int64{}},
		errorChecker: errorChecker,
		locker:       locker,
	}
}

type cbsdManager struct {
	db           *sql.DB
	builder      sqorc.StatementBuilder
	cache        *enumCache
	errorChecker sqorc.ErrorChecker
	locker       sqorc.Locker
}

type enumCache struct {
	cache map[string]map[string]int64
}

func (c *cbsdManager) CreateCbsd(networkId string, data *MutableCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.createCbsd(networkId, data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) UpdateCbsd(networkId string, id int64, data *MutableCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.updateCbsd(networkId, id, data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) EnodebdUpdateCbsd(cbsd *DBCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.enodebdUpdateCbsd(cbsd)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) DeleteCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.markCbsdAsDeleted(networkId, id)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) FetchCbsd(networkId string, id int64) (*DetailedCbsd, error) {
	cbsd, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		return runner.fetchDetailedCbsd(networkId, id)
	})
	if err != nil {
		return nil, makeError(err, c.errorChecker)
	}
	return cbsd.(*DetailedCbsd), nil
}

func (c *cbsdManager) ListCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error) {
	cbsds, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		return runner.listDetailedCbsd(networkId, pagination, filter)
	})
	if err != nil {
		return nil, makeError(err, c.errorChecker)
	}
	return cbsds.(*DetailedCbsdList), nil
}

func (c *cbsdManager) DeregisterCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.markCbsdAsUpdated(networkId, id)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) getInTransactionManager(tx sq.BaseRunner) *cbsdManagerInTransaction {
	return &cbsdManagerInTransaction{
		builder: c.builder.RunWith(tx),
		cache:   c.cache,
		locker:  c.locker,
	}
}

type cbsdManagerInTransaction struct {
	builder sq.StatementBuilderType
	cache   *enumCache
	locker  sqorc.Locker
}

func (c *cbsdManagerInTransaction) createCbsd(networkId string, data *MutableCbsd) error {
	unregisteredState, err := c.cache.getValue(c.builder, &DBCbsdState{}, "unregistered")
	if err != nil {
		return err
	}
	desiredState, err := c.cache.getValue(c.builder, &DBCbsdState{}, data.DesiredState.Name.String)
	if err != nil {
		return err
	}
	data.Cbsd.StateId = db.MakeInt(unregisteredState)
	data.Cbsd.DesiredStateId = db.MakeInt(desiredState)
	data.Cbsd.NetworkId = db.MakeString(networkId)
	columns := append(getCbsdWriteFields(), "state_id", "network_id")
	_, err = db.NewQuery().
		WithBuilder(c.builder).
		From(data.Cbsd).
		Select(db.NewIncludeMask(columns...)).
		Insert()
	return err
}

func (e *enumCache) getValue(builder sq.StatementBuilderType, model db.Model, name string) (int64, error) {
	meta := model.GetMetadata()
	_, ok := e.cache[meta.Table]
	if !ok {
		e.cache[meta.Table] = map[string]int64{}
	}
	if value, ok := e.cache[meta.Table][name]; ok {
		return value, nil
	}
	r, err := db.NewQuery().
		WithBuilder(builder).
		From(model).
		Select(db.NewIncludeMask("id")).
		Where(sq.Eq{"name": name}).
		Fetch()
	if err != nil {
		return 0, err
	}
	e.cache[meta.Table][name] = r[0].(EnumModel).GetId()
	return e.cache[meta.Table][name], nil
}

func getCbsdWriteFields() []string {
	return []string{
		"fcc_id", "cbsd_serial_number", "user_id", "desired_state_id",
		"min_power", "max_power", "antenna_gain", "number_of_ports",
		"preferred_bandwidth_mhz", "preferred_frequencies_mhz", "single_step_enabled",
		"cbsd_category", "latitude_deg", "longitude_deg", "height_m", "height_type", "horizontal_accuracy_m",
		"antenna_azimuth_deg", "antenna_downtilt_deg", "antenna_beamwidth_deg", "antenna_model", "eirp_capability_dbm_mhz",
		"indoor_deployment", "cpi_digital_signature",
	}
}

func getEnodebdWritableFields() []string {
	return []string{
		"antenna_gain", "cbsd_category", "latitude_deg", "longitude_deg",
		"height_m", "height_type", "indoor_deployment", "cpi_digital_signature",
	}
}

func (c *cbsdManagerInTransaction) updateCbsd(networkId string, id int64, data *MutableCbsd) error {
	mask := db.NewIncludeMask("id")
	if _, err := c.selectForUpdateIfCbsdExists(mask, getCbsdFiltersWithId(networkId, id)); err != nil {
		return err
	}
	desiredState, err := c.cache.getValue(c.builder, &DBCbsdState{}, data.DesiredState.Name.String)
	if err != nil {
		return err
	}
	data.Cbsd.DesiredStateId = db.MakeInt(desiredState)
	data.Cbsd.ShouldDeregister = db.MakeBool(true)
	columns := append(getCbsdWriteFields(), "should_deregister")
	return db.NewQuery().
		WithBuilder(c.builder).
		From(data.Cbsd).
		Select(db.NewIncludeMask(columns...)).
		Where(sq.Eq{"id": id}).
		Update()
}

func (c *cbsdManagerInTransaction) enodebdUpdateCbsd(data *DBCbsd) error {
	mask := db.NewIncludeMask("id")
	filters := sq.Eq{"cbsd_serial_number": data.CbsdSerialNumber}
	if _, err := c.selectForUpdateIfCbsdExists(mask, filters); err != nil {
		return err
	}
	data.ShouldDeregister = db.MakeBool(true)
	columns := append(getEnodebdWritableFields(), "should_deregister")
	return db.NewQuery().
		WithBuilder(c.builder).
		From(data).
		Select(db.NewIncludeMask(columns...)).
		Where(filters).
		Update()
}

func (c *cbsdManagerInTransaction) selectForUpdateIfCbsdExists(mask db.FieldMask, filters sq.Eq) (*DBCbsd, error) {
	res, err := db.NewQuery().
		WithBuilder(c.builder).
		From(&DBCbsd{}).
		Select(mask).
		Where(filters).
		Lock(c.locker.WithLock()).
		Fetch()
	if err != nil {
		return nil, err
	}
	return res[0].(*DBCbsd), nil
}

func (c *cbsdManagerInTransaction) markCbsdAsDeleted(networkId string, id int64) error {
	mask := db.NewIncludeMask("id")
	if _, err := c.selectForUpdateIfCbsdExists(mask, getCbsdFiltersWithId(networkId, id)); err != nil {
		return err
	}
	return db.NewQuery().
		WithBuilder(c.builder).
		From(&DBCbsd{IsDeleted: db.MakeBool(true)}).
		Select(db.NewIncludeMask("is_deleted")).
		Where(sq.Eq{"id": id}).
		Update()
}

func (c *cbsdManagerInTransaction) markCbsdAsUpdated(networkId string, id int64) error {
	mask := db.NewIncludeMask("id")
	if _, err := c.selectForUpdateIfCbsdExists(mask, getCbsdFiltersWithId(networkId, id)); err != nil {
		return err
	}
	return db.NewQuery().
		WithBuilder(c.builder).
		From(&DBCbsd{ShouldDeregister: db.MakeBool(true)}).
		Select(db.NewIncludeMask("should_deregister")).
		Where(sq.Eq{"id": id}).
		Update()
}

func (c *cbsdManagerInTransaction) fetchDetailedCbsd(networkId string, id int64) (*DetailedCbsd, error) {
	res, err := buildDetailedCbsdQuery(c.builder).
		Where(getCbsdFiltersWithId(networkId, id)).
		Fetch()
	if err != nil {
		return nil, err
	}
	return convertToDetails(res), nil
}

func convertToDetails(models []db.Model) *DetailedCbsd {
	return &DetailedCbsd{
		Cbsd:         models[0].(*DBCbsd),
		CbsdState:    models[1].(*DBCbsdState),
		DesiredState: models[2].(*DBCbsdState),
		Grant:        models[3].(*DBGrant),
		GrantState:   models[4].(*DBGrantState),
	}
}

func buildDetailedCbsdQuery(builder sq.StatementBuilderType) *db.Query {
	return db.NewQuery().
		WithBuilder(builder).
		From(&DBCbsd{}).
		Select(db.NewExcludeMask("network_id", "state_id", "desired_state_id",
			"is_deleted", "should_deregister", "grant_attempts")).
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
			Select(db.NewIncludeMask(
				"grant_expire_time", "transmit_expire_time",
				"low_frequency", "high_frequency", "max_eirp")).
			Join(db.NewQuery().
				From(&DBGrantState{}).
				On(sq.And{
					db.On(GrantTable, "state_id", GrantStateTable, "id"),
					sq.NotEq{GrantStateTable + ".name": "idle"},
				}).
				Select(db.NewIncludeMask("name"))).
			Nullable())
}

func (c *cbsdManagerInTransaction) listDetailedCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error) {
	count, err := countCbsds(networkId, filter, c.builder)
	if err != nil {
		return nil, err
	}
	query := buildDetailedCbsdQuery(c.builder)
	res, err := buildPagination(query, pagination).
		Where(getCbsdFilters(networkId, filter)).
		OrderBy(CbsdTable+".id", db.OrderAsc).
		List()
	if err != nil {
		return nil, err
	}
	cbsds := make([]*DetailedCbsd, len(res))
	for i, models := range res {
		cbsds[i] = convertToDetails(models)
	}
	return &DetailedCbsdList{
		Cbsds: cbsds,
		Count: count,
	}, nil
}

func countCbsds(networkId string, filter *CbsdFilter, builder sq.StatementBuilderType) (int64, error) {
	return db.NewQuery().
		WithBuilder(builder).
		From(&DBCbsd{}).
		Where(getCbsdFilters(networkId, filter)).
		Count()
}

func makeError(err error, checker sqorc.ErrorChecker) error {
	if err == sql.ErrNoRows {
		return merrors.ErrNotFound
	}
	return checker.GetError(err)
}

func getCbsdFiltersWithId(networkId string, id int64) sq.Eq {
	filters := getCbsdFilters(networkId, nil)
	filters[CbsdTable+".id"] = id
	return filters
}

func getCbsdFilters(networkId string, filter *CbsdFilter) sq.Eq {
	filters := sq.Eq{
		CbsdTable + ".network_id": networkId,
		CbsdTable + ".is_deleted": false,
	}
	if filter != nil {
		if filter.SerialNumber != "" {
			filters[CbsdTable+".cbsd_serial_number"] = filter.SerialNumber
		}
	}
	return filters
}
