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

	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"
)

type CbsdManager interface {
	CreateCbsd(networkId string, data *DBCbsd) error
	UpdateCbsd(networkId string, id int64, data *DBCbsd) error
	DeleteCbsd(networkId string, id int64) error
	FetchCbsd(networkId string, id int64) (*DetailedCbsd, error)
	ListCbsd(networkId string, pagination *Pagination) ([]*DetailedCbsd, error)
}

type DetailedCbsd struct {
	Cbsd       *DBCbsd
	CbsdState  *DBCbsdState
	Channel    *DBChannel
	Grant      *DBGrant
	GrantState *DBGrantState
}

func NewCbsdManager(db *sql.DB, builder sqorc.StatementBuilder) *cbsdManager {
	return &cbsdManager{
		db:      db,
		builder: builder,
		cache:   &enumCache{cache: map[string]map[string]int64{}},
	}
}

type cbsdManager struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
	cache   *enumCache
}

type enumCache struct {
	cache map[string]map[string]int64
}

func (c *cbsdManager) CreateCbsd(networkId string, data *DBCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.createCbsdWithActiveModeConfig(networkId, data)
		return nil, err
	})
	return makeError(err)
}

func (c *cbsdManager) UpdateCbsd(networkId string, id int64, data *DBCbsd) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.updateCbsd(networkId, id, data)
		return nil, err
	})
	return makeError(err)
}

func (c *cbsdManager) DeleteCbsd(networkId string, id int64) error {
	_, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		err := runner.deleteAndDeregisterCbsd(networkId, id)
		return nil, err
	})
	return makeError(err)
}

func (c *cbsdManager) FetchCbsd(networkId string, id int64) (*DetailedCbsd, error) {
	cbsd, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		return runner.fetchDetailedCbsd(networkId, id)
	})
	if err != nil {
		return nil, makeError(err)
	}
	return cbsd.(*DetailedCbsd), nil
}

func (c *cbsdManager) ListCbsd(networkId string, pagination *Pagination) ([]*DetailedCbsd, error) {
	cbsds, err := sqorc.ExecInTx(c.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		runner := c.getInTransactionManager(tx)
		return runner.listDetailedCbsd(networkId, pagination)
	})
	if err != nil {
		return nil, makeError(err)
	}
	return cbsds.([]*DetailedCbsd), nil
}

func (c *cbsdManager) getInTransactionManager(tx sq.BaseRunner) *cbsdManagerInTransaction {
	return &cbsdManagerInTransaction{
		builder: c.builder.RunWith(tx),
		cache:   c.cache,
	}
}

type cbsdManagerInTransaction struct {
	builder sq.StatementBuilderType
	cache   *enumCache
}

func (c *cbsdManagerInTransaction) createCbsdWithActiveModeConfig(networkId string, data *DBCbsd) error {
	unregisteredState, err := c.cache.getValue(c.builder, &DBCbsdState{}, "unregistered")
	if err != nil {
		return err
	}
	data.StateId = db.MakeInt(unregisteredState)
	data.NetworkId = db.MakeString(networkId)
	columns := append(getCbsdWriteFields(), "state_id", "network_id")
	id, err := db.NewQuery().
		WithBuilder(c.builder).
		From(data).
		Select(db.NewIncludeMask(columns...)).
		Insert()
	if err != nil {
		return err
	}
	registeredState, err := c.cache.getValue(c.builder, &DBCbsdState{}, "registered")
	if err != nil {
		return err
	}
	_, err = db.NewQuery().
		WithBuilder(c.builder).
		From(&DBActiveModeConfig{
			CbsdId:         db.MakeInt(id),
			DesiredStateId: db.MakeInt(registeredState),
		}).
		Select(db.NewIncludeMask("cbsd_id", "desired_state_id")).
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
	e.cache[meta.Table][name] = r[0].Fields()["id"].Value().(int64)
	return e.cache[meta.Table][name], nil
}

func getCbsdWriteFields() []string {
	return []string{
		"fcc_id", "cbsd_serial_number", "user_id",
		"min_power", "max_power", "antenna_gain", "number_of_ports",
	}
}

func (c *cbsdManagerInTransaction) updateCbsd(networkId string, id int64, data *DBCbsd) error {
	_, err := c.fetchCbsd(networkId, id, "id")
	if err != nil {
		return err
	}
	return db.NewQuery().
		WithBuilder(c.builder).
		From(data).
		Select(db.NewIncludeMask(getCbsdWriteFields()...)).
		Where(sq.Eq{"id": id}).
		Update()
}

func (c *cbsdManagerInTransaction) fetchCbsd(networkId string, id int64, columns ...string) (*DBCbsd, error) {
	cbsd, err := db.NewQuery().
		WithBuilder(c.builder).
		From(&DBCbsd{}).
		Select(db.NewIncludeMask(columns...)).
		Where(sq.Eq{"network_id": networkId, "id": id}).
		Fetch()
	if err != nil {
		return nil, err
	}
	return cbsd[0].(*DBCbsd), nil
}

func (c *cbsdManagerInTransaction) deleteAndDeregisterCbsd(networkId string, id int64) error {
	cbsd, err := c.fetchCbsd(networkId, id, "cbsd_id")
	if err != nil {
		return err
	}
	if cbsd.CbsdId.Valid {
		if err := c.deregisterCbsd(cbsd.CbsdId.String); err != nil {
			return err
		}
	}
	return db.NewQuery().
		WithBuilder(c.builder).
		From(&DBCbsd{}).
		Where(sq.Eq{"id": id}).
		Delete()
}

func (c *cbsdManagerInTransaction) deregisterCbsd(cbsdId string) error {
	deregistration, err := c.cache.getValue(c.builder, &DBRequestType{}, "deregistrationRequest")
	if err != nil {
		return err
	}
	pending, err := c.cache.getValue(c.builder, &DBRequestState{}, "pending")
	if err != nil {
		return err
	}
	payload := buildDeregistrationPayload(cbsdId)
	request := &DBRequest{
		TypeId:  db.MakeInt(deregistration),
		StateId: db.MakeInt(pending),
		Payload: db.MakeString(payload),
	}
	_, err = db.NewQuery().
		WithBuilder(c.builder).
		From(request).
		Select(db.NewIncludeMask("type_id", "state_id", "payload")).
		Insert()
	return err
}

func buildDeregistrationPayload(cbsdId string) string {
	r := &deregistrationRequestMessage{
		DeregistrationRequest: []*deregistrationRequest{{
			CbsdId: cbsdId,
		}},
	}
	b, _ := json.Marshal(r)
	return string(b)
}

type deregistrationRequestMessage struct {
	DeregistrationRequest []*deregistrationRequest `json:"deregistrationRequest"`
}

type deregistrationRequest struct {
	CbsdId string `json:"cbsdId"`
}

func (c *cbsdManagerInTransaction) fetchDetailedCbsd(networkId string, id int64) (*DetailedCbsd, error) {
	res, err := buildDetailedCbsdQuery(c.builder).
		Where(sq.Eq{
			CbsdTable + ".id":         id,
			CbsdTable + ".network_id": networkId,
		}).
		Fetch()
	if err != nil {
		return nil, err
	}
	return convertToDetails(res), nil
}

func convertToDetails(models []db.Model) *DetailedCbsd {
	return &DetailedCbsd{
		Cbsd:       models[0].(*DBCbsd),
		CbsdState:  models[1].(*DBCbsdState),
		Channel:    models[2].(*DBChannel),
		Grant:      models[3].(*DBGrant),
		GrantState: models[4].(*DBGrantState),
	}
}

func buildDetailedCbsdQuery(builder sq.StatementBuilderType) *db.Query {
	return db.NewQuery().
		WithBuilder(builder).
		From(&DBCbsd{}).
		Select(db.NewExcludeMask("network_id", "state_id")).
		Join(db.NewQuery().
			From(&DBCbsdState{}).
			Select(db.NewIncludeMask("name"))).
		Join(db.NewQuery().
			From(&DBChannel{}).
			Select(db.NewIncludeMask("low_frequency", "high_frequency", "last_used_max_eirp")).
			Join(db.NewQuery().
				From(&DBGrant{}).
				Select(db.NewIncludeMask("grant_expire_time", "transmit_expire_time")).
				Join(db.NewQuery().
					From(&DBGrantState{}).
					Select(db.NewIncludeMask("name")).
					Where(sq.NotEq{GrantStateTable + ".name": "idle"}))).
			Nullable())
}

func (c *cbsdManagerInTransaction) listDetailedCbsd(networkId string, pagination *Pagination) ([]*DetailedCbsd, error) {
	query := buildDetailedCbsdQuery(c.builder)
	res, err := buildPagination(query, pagination).
		Where(sq.Eq{CbsdTable + ".network_id": networkId}).
		OrderBy(CbsdTable+".id", db.OrderAsc).
		List()
	if err != nil {
		return nil, err
	}
	cbsds := make([]*DetailedCbsd, len(res))
	for i, models := range res {
		cbsds[i] = convertToDetails(models)
	}
	return cbsds, nil
}

func makeError(err error) error {
	if err == sql.ErrNoRows {
		return merrors.ErrNotFound
	}
	return err
}
