/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package JsonStore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/thoas/go-funk"
)

const (
	nidCol = "network_id"
	typeCol= "type"
	keyCol = "\"key\""
	valCol = "value"
	verCol = "version"
)

func NewSQLStoreFactory(tableName string, db *sql.DB, sqlBuilder sqorc.StatementBuilder) StoreFactory {
	return &sqlStoreFactory{tableName: tableName, db: db, builder: sqlBuilder}
}

type sqlStoreFactory struct {
	tableName string
	db 		  *sql.DB
	builder   sqorc.StatementBuilder
}

func (fact *sqlStoreFactory) StartTransaction(opts *storage.TxOptions) (Store, error) {
	tx, err := fact.db.BeginTx(context.Background(), getSqlOpts(opts))
	if err != nil {
		return nil, err
	}
	return &sqlStore{tableName: fact.tableName, tx: tx, builder: fact.builder}, nil
}

func getSqlOpts(opts *storage.TxOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}
	if opts.Isolation == 0 {
		return &sql.TxOptions{ReadOnly: opts.ReadOnly}
	}
	return &sql.TxOptions{ReadOnly: opts.ReadOnly, Isolation: sql.IsolationLevel(opts.Isolation)}
}

func (fact *sqlStoreFactory) InitializeFactory() error {
	tx, err := fact.db.Begin()
	if err != nil {
		return err
	}
	err = fact.initTable(tx, fact.tableName)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			glog.Errorf("error rolling back transaction initializing JsonStore factory: %s", rollbackErr)
		}

		return err
	}
	return tx.Commit()
}

func (fact *sqlStoreFactory) initTable(tx *sql.Tx, tableName string) error {
	_, err := fact.builder.CreateTable(tableName).
		IfNotExists().
		Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(typeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(keyCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(valCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(verCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		PrimaryKey(nidCol, typeCol, keyCol).
		RunWith(tx).
		Exec()
	return err
}

type sqlStore struct {
	tableName string
	tx        *sql.Tx
	builder   sqorc.StatementBuilder
}

func (store *sqlStore) Commit() error {
	if store.tx == nil {
		return errors.New("there is no current transaction to commit")
	}

	err := store.tx.Commit()
	store.tx = nil
	return err
}

func (store *sqlStore) Rollback() error {
	if store.tx == nil {
		return errors.New("there is no current transaction to rollback")
	}

	err := store.tx.Rollback()
	store.tx = nil
	return err
}

func (store *sqlStore) Get(networkID string, id storage.TK) (Json, error){
	multiRet, err := store.GetMany(networkID, storage.TKs{id})
	if err != nil {
		return Json{}, err
	}
	if len(multiRet) == 0 {
		return Json{}, merrors.ErrNotFound
	}
	return multiRet[0], nil
}

func (store *sqlStore) GetMany(networkID string, ids storage.TKs) (Jsons, error) {
	if err := store.validateTx(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	whereCondition := getWhereCondition(networkID, ids)
	rows, err := store.builder.Select(typeCol, keyCol, valCol, verCol).From(store.tableName).
		Where(whereCondition).
		RunWith(store.tx).
		Query()
	
	if err != nil {
		return nil ,err
	}

	defer sqorc.CloseRowsLogOnError(rows, "GetMany")

	var jsons Jsons
	for rows.Next() {
		var json Json

		err = rows.Scan(&json.Type, &json.Key, &json.Value, &json.Version)
		if err != nil {
			return nil ,err
		}
		jsons = append(jsons, json)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("sql rows err: %w", err)
	}
	return jsons, nil
}

func (store *sqlStore) Search(filter SearchFilter, criteria LoadCriteria) (map[string]Jsons, error) {
	ret := map[string]Jsons{}
	if err := store.validateTx(); err!=nil {
		return ret , err
	}

	selectCols := []string{nidCol,typeCol, keyCol, verCol}
	if criteria.LoadValue {
		selectCols = append(selectCols, valCol)
	}

	whereCondition := sq.And{}
	if filter.NetworkID != nil {
		whereCondition =  append(whereCondition, sq.Eq{nidCol : *filter.NetworkID})

	}

	if !funk.IsEmpty(filter.Types) {
		whereCondition = append(whereCondition, sq.Eq{typeCol: filter.GetKeys()})
	}
	// Apply only one of prefix or match predicates; prefix takes precedence
	if !funk.IsEmpty(filter.KeyPrefix) {
		whereCondition = append(whereCondition, sq.Like{keyCol: fmt.Sprintf("%s%%", *filter.KeyPrefix)})
	} else {
		if !funk.IsEmpty(filter.Keys) {
			whereCondition = append(whereCondition, sq.Eq{keyCol: filter.GetKeys()})
		}
	}

	rows, err := store.builder.Select(selectCols...).From(store.tableName).
		Where(whereCondition).
		RunWith(store.tx).
		Query()
	if err != nil {
		return ret, fmt.Errorf("failed to query DB: %w", err)
	}
	defer sqorc.CloseRowsLogOnError(rows, "GetMany")

	for rows.Next() {
		var nid, t, k string
		var version uint64
		var val string
		scanArgs := []interface{}{&nid, &t, &k, &version}
		if criteria.LoadValue {
			scanArgs = append(scanArgs, &val)
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return ret, fmt.Errorf("failed to scan blob row: %w", err)
		}

		nidCol := ret[nid]
		nidCol = append(nidCol, Json{Type: t, Key: k, Value: val, Version: version})
		ret[nid] = nidCol
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("sql rows err: %w", err)
	}
	return ret, nil
}

func (store *sqlStore) Write(networkID string, jsons Jsons) error {
	existingJsons, err := store.GetMany(networkID, getJsonIDs(jsons))

	if err != nil {
		return fmt.Errorf("error reading existing jsons: %s", err)
	}
	jsonsToCreateAndChange := partitionJsonsToCreateAndChange(jsons, existingJsons)

	if len(jsonsToCreateAndChange.jsonsToChange) > 0 {
		err := store.updateExistingJsons(networkID, jsonsToCreateAndChange.jsonsToChange)
		if err != nil {
			return err
		}
	}
	if len(jsonsToCreateAndChange.jsonsToCreate) > 0 {
		err := store.insertNewJsons(networkID, jsonsToCreateAndChange.jsonsToCreate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *sqlStore) Delete(networkID string, ids storage.TKs) error {
	if err := store.validateTx(); err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	whereCondition := getWhereCondition(networkID, ids)
	_, err := store.builder.Delete(store.tableName).
		Where(whereCondition).
		RunWith(store.tx).
		Exec()
	return err
}

func (store *sqlStore) GetExistingKeys(keys []string, filter SearchFilter) ([]string, error) {
	if err := store.validateTx(); err != nil {
		return nil, err
	}

	whereConditions := make(sq.Or, 0, len(keys))
	for _, key := range keys {
		and := sq.And{sq.Eq{keyCol: key}}
		if funk.NotEmpty(filter.NetworkID) {
			and = append(and, sq.Eq{nidCol: filter.NetworkID})
		}

		whereConditions = append(whereConditions, and)
	}
	rows, err := store.builder.Select(keyCol).Distinct().From(store.tableName).
		Where(whereConditions).
		RunWith(store.tx).
		Query()
	if err != nil {
		return nil, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "GetExistingKeys")
	var scannedKeys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, err
		}
		scannedKeys = append(scannedKeys, key)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("sql rows err: %w", err)
	}
	return scannedKeys, nil
}

func (store *sqlStore) IncrementVersion(networkID string, id storage.TK) error {
	if err := store.validateTx(); err != nil {
		return err
	}

	_, err := store.builder.Insert(store.tableName).
		Columns(nidCol, typeCol, keyCol, verCol).
		Values(networkID, id.Type, id.Key, 1).
		OnConflict(
			[]sqorc.UpsertValue{{Column: verCol, Value: sq.Expr(fmt.Sprintf("%s.%s+1", store.tableName, verCol))}},
			nidCol, typeCol, keyCol,
		).
		RunWith(store.tx).
		Exec()
	if err != nil {
		return fmt.Errorf("error incrementing version on network %s with type %s and key %s: %w", networkID, id.Type, id.Key, err)
	}
	return nil
}

func (store *sqlStore) validateTx() error {
	if store.tx == nil {
		return errors.New("no transaction is available")
	}
	return nil
}

func (store *sqlStore) updateExistingJsons(networkID string, JsonsToChange map[storage.TK]JsonChange) error {
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "updateExistingJsons")

	for _, jsonID := range getSortedTKs(JsonsToChange) {
		change := JsonsToChange[jsonID]
		updatedVersion := change.old.Version + 1
		if change.new.Version != 0 {
			updatedVersion = change.new.Version
		}
		_, err := store.builder.Update(store.tableName).
			Set(valCol, change.new.Value).
			Set(verCol, updatedVersion).
			Where(
				// Use explicit sq.And to preserve ordering of WHERE clause items
				sq.And{
					sq.Eq{nidCol: networkID},
					sq.Eq{typeCol: jsonID.Type},
					sq.Eq{keyCol: jsonID.Key},
				},
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return fmt.Errorf("error updating blob (%s, %s, %s): %s", networkID, jsonID.Type, jsonID.Key, err)
		}
	}
	return nil
}

func (store *sqlStore) insertNewJsons(networkID string, jsons Jsons) error {
	insertBuilder := store.builder.Insert(store.tableName).
		Columns(nidCol, typeCol, keyCol, valCol, verCol)
	for _, json := range jsons {
		insertBuilder = insertBuilder.Values(networkID, json.Type, json.Key, json.Value, json.Version)
	}
	_, err := insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return fmt.Errorf("error creating jsons: %w", err)
	}
	return nil
}

func getWhereCondition(networkID string, ids storage.TKs) sq.Or {
	whereConditions := make(sq.Or, 0, len(ids))
	for _, id := range ids {
		// Use explicit sq.And to preserve ordering of clauses for testing
		whereConditions = append(whereConditions, sq.And{
			sq.Eq{nidCol: networkID},
			sq.Eq{typeCol: id.Type},
			sq.Eq{keyCol: id.Key},
		})
	}
	return whereConditions
}

func getJsonIDs(jsons Jsons) storage.TKs {
	ret := make(storage.TKs, 0, len(jsons))
	for _, json := range jsons {
		ret = append(ret, storage.TK{Type: json.Type, Key: json.Key})
	}
	return ret
}

type JsonChange struct {
	old Json
	new Json
}

type JsonsToCreateAndChange struct {
	jsonsToCreate Jsons
	jsonsToChange map[storage.TK]JsonChange
}

func partitionJsonsToCreateAndChange(jsonsToUpdate Jsons, existingJsons Jsons) JsonsToCreateAndChange {
	ret := JsonsToCreateAndChange{
		jsonsToCreate: Jsons{},
		jsonsToChange: map[storage.TK]JsonChange{},
	}
	existingJsonsByID := existingJsons.ByTK()

	for _, json := range jsonsToUpdate {
		jsonID := storage.TK{Type: json.Type, Key: json.Key}
		oldjson, exists := existingJsonsByID[jsonID]
		if exists {
			ret.jsonsToChange[jsonID] = JsonChange{old: oldjson, new: json}
		} else {
			ret.jsonsToCreate = append(ret.jsonsToCreate, json)
		}
	}
	return ret
}

func getSortedTKs(JsonToChange map[storage.TK]JsonChange) storage.TKs {
	ret := make(storage.TKs, 0, len(JsonToChange))
	for k := range JsonToChange {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].String() < ret[j].String() })
	return ret
}

