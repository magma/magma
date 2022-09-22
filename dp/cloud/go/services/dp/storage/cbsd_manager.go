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
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/merrors"
)

type CbsdManager interface {
	CreateCbsd(networkId string, data *MutableCbsd) error
	UpdateCbsd(networkId string, id int64, data *MutableCbsd) error
	EnodebdUpdateCbsd(data *DBCbsd) (*DetailedCbsd, error)
	DeleteCbsd(networkId string, id int64) error
	FetchCbsd(networkId string, id int64) (*DetailedCbsd, error)
	ListCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error)
	DeregisterCbsd(networkId string, id int64) error
	RelinquishCbsd(networkId string, id int64) error
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
	Grants       []*DetailedGrant
}

type DetailedGrant struct {
	Grant      *DBGrant
	GrantState *DBGrantState
}

func NewCbsdManager(db *sql.DB, builder sqorc.StatementBuilder, errorChecker sqorc.ErrorChecker, locker sqorc.Locker) *cbsdManager {
	return &cbsdManager{
		&dpManager{
			db:           db,
			builder:      builder,
			cache:        &enumCache{cache: map[string]map[string]int64{}},
			errorChecker: errorChecker,
			locker:       locker,
		},
	}
}

type cbsdManager struct {
	*dpManager
}

type enumCache struct {
	cache map[string]map[string]int64
}

func (c *cbsdManager) CreateCbsd(networkId string, data *MutableCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		err := runner.createCbsd(networkId, data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) UpdateCbsd(networkId string, id int64, data *MutableCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		err := runner.updateCbsd(networkId, id, data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) EnodebdUpdateCbsd(data *DBCbsd) (*DetailedCbsd, error) {
	grantJoinClause := getGrantJoinClauseForEnodebdUpdate()
	result, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		err := runner.enodebdUpdateCbsd(data)
		if err != nil {
			return nil, err
		}
		return runner.fetchDetailedCbsd(sq.Eq{"cbsd_serial_number": data.CbsdSerialNumber}, grantJoinClause)
	})
	if err != nil {
		return nil, makeError(err, c.errorChecker)
	}
	return result.(*DetailedCbsd), nil
}

func (c *cbsdManager) DeleteCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		data := &DBCbsd{IsDeleted: db.MakeBool(true)}
		err := runner.updateField(networkId, id, "is_deleted", data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) FetchCbsd(networkId string, id int64) (*DetailedCbsd, error) {
	cbsd, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		grantJoinClause := db.On(GrantTable, "state_id", GrantStateTable, "id")
		return runner.fetchDetailedCbsd(getCbsdFiltersWithId(networkId, id), grantJoinClause)
	})
	if err != nil {
		return nil, makeError(err, c.errorChecker)
	}
	return cbsd.(*DetailedCbsd), nil
}

func (c *cbsdManager) ListCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error) {
	cbsds, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		return runner.listDetailedCbsd(networkId, pagination, filter)
	})
	if err != nil {
		return nil, makeError(err, c.errorChecker)
	}
	return cbsds.(*DetailedCbsdList), nil
}

func (c *cbsdManager) DeregisterCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		data := &DBCbsd{ShouldDeregister: db.MakeBool(true)}
		err := runner.updateField(networkId, id, "should_deregister", data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (c *cbsdManager) RelinquishCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getQueryRunner(tx)
		data := &DBCbsd{ShouldRelinquish: db.MakeBool(true)}
		err := runner.updateField(networkId, id, "should_relinquish", data)
		return nil, err
	})
	return makeError(err, c.errorChecker)
}

func (r *queryRunner) createCbsd(networkId string, data *MutableCbsd) error {
	unregisteredState, err := r.cache.getValue(r.builder, &DBCbsdState{}, "unregistered")
	if err != nil {
		return err
	}
	desiredState, err := r.cache.getValue(r.builder, &DBCbsdState{}, data.DesiredState.Name.String)
	if err != nil {
		return err
	}
	data.Cbsd.StateId = db.MakeInt(unregisteredState)
	data.Cbsd.DesiredStateId = db.MakeInt(desiredState)
	data.Cbsd.NetworkId = db.MakeString(networkId)
	columns := append(getCbsdWriteFields(), "state_id", "network_id")
	mask := db.NewIncludeMask(columns...)
	_, err = db.NewQuery().
		WithBuilder(r.builder).
		From(data.Cbsd).
		Insert(mask)
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
		"indoor_deployment", "cpi_digital_signature", "carrier_aggregation_enabled", "max_ibw_mhz", "grant_redundancy",
	}
}

func getEnodebdWritableFields() []string {
	return []string{
		"cbsd_category", "latitude_deg", "longitude_deg",
		"height_m", "height_type", "indoor_deployment", "cpi_digital_signature",
	}
}

func (r *queryRunner) updateCbsd(networkId string, id int64, data *MutableCbsd) error {
	mask := db.NewIncludeMask("id")
	if _, err := r.selectForUpdateIfCbsdExists(mask, getCbsdFiltersWithId(networkId, id)); err != nil {
		return err
	}
	desiredState, err := r.cache.getValue(r.builder, &DBCbsdState{}, data.DesiredState.Name.String)
	if err != nil {
		return err
	}
	data.Cbsd.DesiredStateId = db.MakeInt(desiredState)
	data.Cbsd.ShouldDeregister = db.MakeBool(true)
	columns := append(getCbsdWriteFields(), "should_deregister")
	mask = db.NewIncludeMask(columns...)
	_, err = db.NewQuery().
		WithBuilder(r.builder).
		From(data.Cbsd).
		Select(db.NewIncludeMask()).
		Where(sq.Eq{"id": id}).
		Update(mask)
	return err
}

func (r *queryRunner) enodebdUpdateCbsd(data *DBCbsd) error {
	identifiers := []string{"cbsd_serial_number", "network_id"}
	maskFields := append(identifiers, getEnodebdWritableFields()...)
	mask := db.NewIncludeMask(maskFields...)
	filters := sq.Eq{"cbsd_serial_number": data.CbsdSerialNumber}
	cbsd, err := r.selectForUpdateIfCbsdExists(mask, filters)

	if err != nil {
		return err
	}

	columns := []string{"last_seen"}

	if ShouldEnodebdUpdateInstallationParams(cbsd, data) {
		cols := append(getEnodebdWritableFields(), "should_deregister")
		columns = append(columns, cols...)
		data.ShouldDeregister = db.MakeBool(true)
	}

	_, err = db.NewQuery().
		WithBuilder(r.builder).
		From(data).
		Select(db.NewIncludeMask("id")).
		Where(filters).
		Update(db.NewIncludeMask(columns...))

	return err
}

func (r *queryRunner) selectForUpdateIfCbsdExists(mask db.FieldMask, filters sq.Eq) (*DBCbsd, error) {
	res, err := db.NewQuery().
		WithBuilder(r.builder).
		From(&DBCbsd{}).
		Select(mask).
		Where(filters).
		Lock(r.locker.WithLock()).
		Fetch()
	if err != nil {
		return nil, err
	}
	return res[0].(*DBCbsd), nil
}

func (r *queryRunner) updateField(networkId string, id int64, field string, data *DBCbsd) error {
	mask := db.NewIncludeMask("id")
	if _, err := r.selectForUpdateIfCbsdExists(mask, getCbsdFiltersWithId(networkId, id)); err != nil {
		return err
	}
	mask = db.NewIncludeMask(field)
	_, err := db.NewQuery().
		WithBuilder(r.builder).
		From(data).
		Select(db.NewIncludeMask()).
		Where(sq.Eq{"id": id}).
		Update(mask)
	if err != nil {
		return err
	}
	return nil
}

func (r *queryRunner) fetchDetailedCbsd(filter sq.Eq, grantJoinClause sq.Sqlizer) (*DetailedCbsd, error) {
	rawCbsd, err := buildDetailedCbsdQuery(r.builder).
		Where(filter).
		Fetch()
	if err != nil {
		return nil, err
	}
	cbsd := convertCbsdToDetails(rawCbsd)
	if err := getGrantsForCbsds(r.builder, grantJoinClause, cbsd); err != nil {
		return nil, err
	}
	return cbsd, nil
}

func convertCbsdToDetails(models []db.Model) *DetailedCbsd {
	return &DetailedCbsd{
		Cbsd:         models[0].(*DBCbsd),
		CbsdState:    models[1].(*DBCbsdState),
		DesiredState: models[2].(*DBCbsdState),
	}
}

func buildDetailedCbsdQuery(builder sq.StatementBuilderType) *db.Query {
	return db.NewQuery().
		WithBuilder(builder).
		From(&DBCbsd{}).
		Select(db.NewExcludeMask(
			"state_id", "desired_state_id",
			"is_deleted", "should_deregister", "should_relinquish",
			"available_frequencies", "channels")).
		Join(db.NewQuery().
			From(&DBCbsdState{}).
			As("t1").
			On(db.On(CbsdTable, "state_id", "t1", "id")).
			Select(db.NewIncludeMask("name"))).
		Join(db.NewQuery().
			From(&DBCbsdState{}).
			As("t2").
			On(db.On(CbsdTable, "desired_state_id", "t2", "id")).
			Select(db.NewIncludeMask("name")))
}

func getGrantsForCbsds(builder sq.StatementBuilderType, grantJoinClause sq.Sqlizer, cbsds ...*DetailedCbsd) error {
	idList, idMap := make([]int64, len(cbsds)), make(map[int64]*DetailedCbsd, len(cbsds))
	for i, c := range cbsds {
		idList[i] = c.Cbsd.Id.Int64
		idMap[c.Cbsd.Id.Int64] = c
	}
	rawGrants, err := buildDetailedGrantQuery(builder, grantJoinClause).
		Where(sq.Eq{"cbsd_id": idList}).
		OrderBy(GrantTable+".low_frequency", db.OrderAsc).
		List()
	if err != nil {
		return err
	}
	for _, models := range rawGrants {
		g := &DetailedGrant{
			Grant:      models[0].(*DBGrant),
			GrantState: models[1].(*DBGrantState),
		}
		c := idMap[g.Grant.CbsdId.Int64]
		g.Grant.CbsdId = sql.NullInt64{}
		c.Grants = append(c.Grants, g)
	}
	return nil
}

func buildDetailedGrantQuery(builder sq.StatementBuilderType, on sq.Sqlizer) *db.Query {
	return db.NewQuery().
		WithBuilder(builder).
		From(&DBGrant{}).
		Select(db.NewIncludeMask(
			"cbsd_id", "grant_expire_time", "transmit_expire_time",
			"low_frequency", "high_frequency", "max_eirp")).
		Join(db.NewQuery().
			From(&DBGrantState{}).
			On(on).
			Select(db.NewIncludeMask("name")))
}

func (r *queryRunner) listDetailedCbsd(networkId string, pagination *Pagination, filter *CbsdFilter) (*DetailedCbsdList, error) {
	count, err := countCbsds(networkId, filter, r.builder)
	if err != nil {
		return nil, err
	}
	query := buildDetailedCbsdQuery(r.builder)
	res, err := buildPagination(query, pagination).
		Where(getCbsdFilters(networkId, filter)).
		OrderBy(CbsdTable+".id", db.OrderAsc).
		List()
	if err != nil {
		return nil, err
	}
	cbsds := make([]*DetailedCbsd, len(res))
	for i, models := range res {
		cbsds[i] = convertCbsdToDetails(models)
	}
	on := db.On(GrantTable, "state_id", GrantStateTable, "id")
	if err := getGrantsForCbsds(r.builder, on, cbsds...); err != nil {
		return nil, err
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

func getGrantJoinClauseForEnodebdUpdate() sq.Sqlizer {
	return sq.And{
		db.On(GrantTable, "state_id", GrantStateTable, "id"),
		sq.Eq{GrantStateTable + ".name": "authorized"},
		sq.Or{
			sq.Eq{"transmit_expire_time": nil},
			sq.Gt{"transmit_expire_time": db.MakeTime(clock.Now().UTC())},
		},
		sq.Or{
			sq.Eq{"grant_expire_time": nil},
			sq.Gt{"grant_expire_time": db.MakeTime(clock.Now().UTC())},
		},
	}
}
