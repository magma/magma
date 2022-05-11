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

package blobstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
)

const (
	nidCol  = "network_id"
	typeCol = "type"
	keyCol  = "\"key\""
	valCol  = "value"
	verCol  = "version"
)

// NewSQLStoreFactory returns a StoreFactory implementation which
// will return storage APIs backed by SQL.
func NewSQLStoreFactory(tableName string, db *sql.DB, sqlBuilder sqorc.StatementBuilder) StoreFactory {
	return &sqlStoreFactory{tableName: tableName, db: db, builder: sqlBuilder}
}

type sqlStoreFactory struct {
	tableName string
	db        *sql.DB
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
			glog.Errorf("error rolling back transaction initializing blobstore factory: %s", rollbackErr)
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
		Column(valCol).Type(sqorc.ColumnTypeBytes).EndColumn().
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
		return errors.New("There is no current transaction to commit")
	}

	err := store.tx.Commit()
	store.tx = nil
	return err
}

func (store *sqlStore) Rollback() error {
	if store.tx == nil {
		return errors.New("There is no current transaction to rollback")
	}

	err := store.tx.Rollback()
	store.tx = nil
	return err
}

func (store *sqlStore) Get(networkID string, id storage.TK) (Blob, error) {
	multiRet, err := store.GetMany(networkID, storage.TKs{id})
	if err != nil {
		return Blob{}, err
	}
	if len(multiRet) == 0 {
		return Blob{}, merrors.ErrNotFound
	}
	return multiRet[0], nil
}

func (store *sqlStore) GetMany(networkID string, ids storage.TKs) (Blobs, error) {
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
		return nil, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "GetMany")

	var blobs Blobs
	for rows.Next() {
		var t, k string
		var val []byte
		var version uint64

		err = rows.Scan(&t, &k, &val, &version)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, Blob{Type: t, Key: k, Value: val, Version: version})
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("sql rows err: %w", err)
	}
	return blobs, nil
}

func (store *sqlStore) Search(filter SearchFilter, criteria LoadCriteria) (map[string]Blobs, error) {
	ret := map[string]Blobs{}
	if err := store.validateTx(); err != nil {
		return ret, err
	}

	// Get select columns from load criteria
	selectCols := []string{nidCol, typeCol, keyCol, verCol}
	if criteria.LoadValue {
		selectCols = append(selectCols, valCol)
	}

	// Get where condition from search filter
	// Use and condition to deterministically order clauses for testing
	whereCondition := sq.And{}
	if filter.NetworkID != nil {
		whereCondition = append(whereCondition, sq.Eq{nidCol: *filter.NetworkID})
	}
	if !funk.IsEmpty(filter.Types) {
		whereCondition = append(whereCondition, sq.Eq{typeCol: filter.GetTypes()})
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
		var val []byte
		scanArgs := []interface{}{&nid, &t, &k, &version}
		if criteria.LoadValue {
			scanArgs = append(scanArgs, &val)
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return ret, fmt.Errorf("failed to scan blob row: %w", err)
		}

		nidCol := ret[nid]
		nidCol = append(nidCol, Blob{Type: t, Key: k, Value: val, Version: version})
		ret[nid] = nidCol
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("sql rows err: %w", err)
	}
	return ret, nil
}

func (store *sqlStore) Write(networkID string, blobs Blobs) error {
	// defer tx validation to GetMany
	existingBlobs, err := store.GetMany(networkID, getBlobIDs(blobs))
	if err != nil {
		return fmt.Errorf("Error reading existing blobs: %s", err)
	}
	blobsToCreateAndChange := partitionBlobsToCreateAndChange(blobs, existingBlobs)

	if len(blobsToCreateAndChange.blobsToChange) > 0 {
		err := store.updateExistingBlobs(networkID, blobsToCreateAndChange.blobsToChange)
		if err != nil {
			return err
		}
	}
	if len(blobsToCreateAndChange.blobsToCreate) > 0 {
		err := store.insertNewBlobs(networkID, blobsToCreateAndChange.blobsToCreate)
		if err != nil {
			return err
		}
	}

	return nil
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
		return fmt.Errorf("Error incrementing version on network %s with type %s and key %s: %w", networkID, id.Type, id.Key, err)
	}
	return nil
}

func (store *sqlStore) validateTx() error {
	if store.tx == nil {
		return errors.New("no transaction is available")
	}
	return nil
}

func (store *sqlStore) updateExistingBlobs(networkID string, blobsToChange map[storage.TK]blobChange) error {
	// Let squirrel cache prepared statements for us (there should only be 1)
	sc := sq.NewStmtCache(store.tx)
	defer sqorc.ClearStatementCacheLogOnError(sc, "updateExistingBlobs")

	// Sort keys for deterministic behavior in tests
	for _, blobID := range getSortedTKs(blobsToChange) {
		change := blobsToChange[blobID]
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
					sq.Eq{typeCol: blobID.Type},
					sq.Eq{keyCol: blobID.Key},
				},
			).
			RunWith(sc).
			Exec()
		if err != nil {
			return fmt.Errorf("error updating blob (%s, %s, %s): %s", networkID, blobID.Type, blobID.Key, err)
		}
	}
	return nil
}

func (store *sqlStore) insertNewBlobs(networkID string, blobs Blobs) error {
	insertBuilder := store.builder.Insert(store.tableName).
		Columns(nidCol, typeCol, keyCol, valCol, verCol)
	for _, blob := range blobs {
		insertBuilder = insertBuilder.Values(networkID, blob.Type, blob.Key, blob.Value, blob.Version)
	}
	_, err := insertBuilder.RunWith(store.tx).Exec()
	if err != nil {
		return fmt.Errorf("error creating blobs: %w", err)
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

func getBlobIDs(blobs Blobs) storage.TKs {
	ret := make(storage.TKs, 0, len(blobs))
	for _, blob := range blobs {
		ret = append(ret, storage.TK{Type: blob.Type, Key: blob.Key})
	}
	return ret
}

type blobChange struct {
	old Blob
	new Blob
}

type blobsToCreateAndChange struct {
	blobsToCreate Blobs
	blobsToChange map[storage.TK]blobChange
}

func partitionBlobsToCreateAndChange(blobsToUpdate Blobs, existingBlobs Blobs) blobsToCreateAndChange {
	ret := blobsToCreateAndChange{
		blobsToCreate: Blobs{},
		blobsToChange: map[storage.TK]blobChange{},
	}
	existingBlobsByID := existingBlobs.ByTK()

	for _, blob := range blobsToUpdate {
		blobID := storage.TK{Type: blob.Type, Key: blob.Key}
		oldBlob, exists := existingBlobsByID[blobID]
		if exists {
			ret.blobsToChange[blobID] = blobChange{old: oldBlob, new: blob}
		} else {
			ret.blobsToCreate = append(ret.blobsToCreate, blob)
		}
	}
	return ret
}

func getSortedTKs(blobsToChange map[storage.TK]blobChange) storage.TKs {
	ret := make(storage.TKs, 0, len(blobsToChange))
	for k := range blobsToChange {
		ret = append(ret, k)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].String() < ret[j].String() })
	return ret
}
