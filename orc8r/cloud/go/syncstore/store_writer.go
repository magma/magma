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

package syncstore

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/Masterminds/squirrel"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type syncStore struct {
	cacheWriterValidIntervalSecs int64
	digestTableName              string
	cacheTableName               string
	tableNamePrefix              string
	db                           *sql.DB
	builder                      sqorc.StatementBuilder
	fact                         blobstore.BlobStorageFactory
}

func NewSyncStore(db *sql.DB, builder sqorc.StatementBuilder, fact blobstore.BlobStorageFactory, config Config) (SyncStore, error) {
	err := config.Validate(true)
	if err != nil {
		return nil, errors.Wrap(err, "invalid configs for syncstore")
	}
	store := &syncStore{
		db:                           db,
		builder:                      builder,
		fact:                         fact,
		cacheWriterValidIntervalSecs: config.CacheWriterValidIntervalSecs,
		tableNamePrefix:              config.TableNamePrefix,
		digestTableName:              fmt.Sprintf("%s_digest", config.TableNamePrefix),
		cacheTableName:               fmt.Sprintf("%s_cached_objs", config.TableNamePrefix),
	}
	return store, nil
}

func (l *syncStore) SetDigest(network string, digests *protos.DigestTree) error {
	rootDigest := digests.RootDigest.GetMd5Base64Digest()
	leafDigestsToSerialize := &protos.LeafDigests{Digests: digests.GetLeafDigests()}
	leafDigests, err := proto.Marshal(leafDigestsToSerialize)
	if err != nil {
		return errors.Wrapf(err, "marshal leaf digests for network %+v", network)
	}
	now := clock.Now().Unix()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Insert(fmt.Sprintf(l.digestTableName)).
			Columns(nidCol, rootDigestCol, leafDigestsCol, lastUpdatedTimeCol).
			Values(network, rootDigest, leafDigests, now).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: rootDigestCol, Value: rootDigest},
					{Column: leafDigestsCol, Value: leafDigests},
					{Column: lastUpdatedTimeCol, Value: now},
				},
				nidCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert digests for network %+v", network)
		}
		return nil, nil
	}
	_, err = sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *syncStore) UpdateCache(network string) (CacheWriter, error) {
	// The temporary table is namespaced with the unique id of the cacheWriter
	writerID := generateCacheWriterUUID(l.tableNamePrefix)
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(writerID).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(idCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(objCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(nidCol, idCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "create cached objs tmp table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}

	// The start time of cacheWriters is tracked by the store for garbage collection
	err = l.recordCacheWriterStartTime(network, writerID)
	if err != nil {
		return nil, errors.Wrapf(err, "record start time of cache writer %+v of network %+v", writerID, network)
	}

	return l.NewCacheWriter(network, writerID), nil
}

func (l *syncStore) RecordResync(network string, gateway string, t int64) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	err = store.CreateOrUpdate(network, blobstore.Blobs{{
		Type:  lastResyncBlobstoreType,
		Key:   gateway,
		Value: encodeInt64(t),
	}})
	if err != nil {
		return errors.Wrapf(err, "set last resync time of network %+v, gateway %+v in blobstore", network, gateway)
	}

	return store.Commit()
}

// recordCacheWriterStartTime registers the creation time of a cache writer in store.
func (l *syncStore) recordCacheWriterStartTime(network string, writerID string) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	err = store.CreateOrUpdate(network, blobstore.Blobs{{
		Type:  cacheWriterBlobstoreType,
		Key:   writerID,
		Value: encodeInt64(clock.Now().Unix()),
	}})
	if err != nil {
		return errors.Wrapf(err, "set start time of network %+v, cachewriter %+v in blobstore", network, writerID)
	}
	return store.Commit()
}

func (l *syncStore) CollectGarbage(trackedNetworks []string) {
	err := l.collectGarbageSQL(trackedNetworks)
	if err != nil {
		glog.Errorf("Collect syncstore garbage in sql tables: %+v", err)
	}

	err = l.collectGarbageLastResync(trackedNetworks)
	if err != nil {
		glog.Errorf("Collect syncstore garbage for last resync times: %+v", err)
	}

	err = l.collectGarbageCacheWriter(trackedNetworks)
	if err != nil {
		glog.Errorf("Collect syncstore garbage for cache writers: %+v", err)
	}
}

// collectGarbageSQL drops all contents in the digests and cached objects SQL storage
// that are unrelated to the tracked networks.
func (l *syncStore) collectGarbageSQL(tracked []string) error {
	errs := &multierror.Error{}
	tableNames := []string{l.cacheTableName, l.digestTableName}
	for _, tableName := range tableNames {
		txFn := func(tx *sql.Tx) (interface{}, error) {
			stored, err := l.getStoredNetworksSQL(tx, tableName)
			if err != nil {
				return nil, err
			}
			deleted, _ := funk.DifferenceString(stored, tracked)
			_, err = l.builder.
				Delete(tableName).
				Where(squirrel.Eq{nidCol: deleted}).
				RunWith(tx).
				Exec()
			return nil, err
		}
		_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "collect garbage for table %+v", tableName))
		}
	}
	return errs.ErrorOrNil()
}

func (l *syncStore) getStoredNetworksSQL(tx *sql.Tx, tableName string) ([]string, error) {
	rows, err := l.builder.Select(nidCol).From(tableName).RunWith(tx).Query()
	if err != nil {
		return nil, errors.Wrapf(err, "get all networks in store %+v", tableName)
	}
	var storedNetworks []string
	for rows.Next() {
		network := ""
		err = rows.Scan(&network)
		if err != nil {
			return nil, errors.Wrapf(err, "get all networks in store %+v, SQL rows scan error", tableName)
		}
		storedNetworks = append(storedNetworks, network)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "get all networks in store %+v, SQL rows error", tableName)
	}
	return storedNetworks, nil
}

func getStoredNetworksBlobstore(store blobstore.TransactionalBlobStorage) ([]string, error) {
	keysByNetwork, err := blobstore.ListKeysByNetwork(store)
	if err != nil {
		return nil, errors.Wrap(err, "list blobstore keys by network")
	}
	return funk.Keys(keysByNetwork).([]string), nil
}

// collectGarbageLastResync drops all lastResync type blobstore items unrelated
// to the tracked networks.
func (l *syncStore) collectGarbageLastResync(tracked []string) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer store.Rollback()

	stored, err := getStoredNetworksBlobstore(store)
	if err != nil {
		return errors.Wrap(err, "get all networks in blobstore")
	}
	deleted, _ := funk.DifferenceString(stored, tracked)

	errs := &multierror.Error{}
	for _, network := range deleted {
		keys, err := blobstore.ListKeys(store, network, lastResyncBlobstoreType)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		err = store.Delete(network, storage.MakeTKs(lastResyncBlobstoreType, keys))
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	err = store.Commit()
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs.ErrorOrNil()
}

func (l *syncStore) collectGarbageCacheWriter(tracked []string) error {
	errs := &multierror.Error{}

	invalidByNetwork, err := l.getInvalidCacheWriter(tracked, l.cacheWriterValidIntervalSecs)
	if err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "get invalid cache writers for tracked networks %+v", tracked))
	}

	// Attempt to drop the tmp tables of all invalid cacheWriters, and only delete the blobstore records of those
	// whose tables have been successfully dropped; the rest is left to be garbage collected in future runs
	deletedByNetwork, err := l.dropInvalidCaches(invalidByNetwork)
	if err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "drop invalid cache writer tables %+v", invalidByNetwork))
	}

	err = l.deleteCacheWriterBlobstoreRecords(deletedByNetwork)
	if err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "delete cache writer blobstore records %+v", deletedByNetwork))
	}

	return errs.ErrorOrNil()
}

// getInvalidCacheWriter returns a list of cache writer IDs from blobstore that
// either belong to already deleted networks or have expired.
func (l *syncStore) getInvalidCacheWriter(tracked []string, cacheWriterValidIntervalSecs int64) (map[string][]string, error) {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "error starting transaction")
	}
	defer store.Rollback()

	stored, err := getStoredNetworksBlobstore(store)
	if err != nil {
		return nil, errors.Wrap(err, "get all networks in blobstore")
	}

	deleted, _ := funk.DifferenceString(stored, tracked)

	invalidByNetwork := map[string][]string{}
	errs := &multierror.Error{}
	for _, network := range deleted {
		keys, err := blobstore.ListKeys(store, network, cacheWriterBlobstoreType)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "list cache writers of deleted network %+v", network))
			continue
		}
		invalidByNetwork[network] = keys
	}

	for _, network := range tracked {
		keys, err := blobstore.ListKeys(store, network, cacheWriterBlobstoreType)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "list all cache-writer-type blobstore keys of network %+v", network))
			continue
		}
		blobs, err := store.GetMany(network, storage.MakeTKs(cacheWriterBlobstoreType, keys))
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "get cache writer blobs of network %+v", network))
			continue
		}

		var invalid []string
		for _, blob := range blobs {
			creationTime := binary.LittleEndian.Uint64(blob.Value)
			if clock.Now().Unix()-int64(creationTime) > cacheWriterValidIntervalSecs {
				invalid = append(invalid, blob.Key)
			}
		}
		invalidByNetwork[network] = invalid
	}
	err = store.Commit()
	if err != nil {
		errs = multierror.Append(errs, err)
	}
	return invalidByNetwork, errs.ErrorOrNil()
}

// dropInvalidCaches drops the temporary caches held by invalid cache writers, and
// returns the IDs of those successfully dropped.
func (l *syncStore) dropInvalidCaches(invalidByNetwork map[string][]string) (map[string][]string, error) {
	errs := &multierror.Error{}
	deletedByNetwork := map[string][]string{}
	for network, invalid := range invalidByNetwork {
		var deleted []string
		for _, tableName := range invalid {
			txFn := func(tx *sql.Tx) (interface{}, error) {
				_, err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
				return nil, errors.Wrapf(err, "drop cache writer table %+v for network %+v", tableName, network)
			}
			_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}
			deleted = append(deleted, tableName)
		}
		deletedByNetwork[network] = deleted
	}
	return deletedByNetwork, errs.ErrorOrNil()
}

// deleteCacheWriterBlobstoreRecords removes the blobstore records of cache writers
// that are invalid, and whose temporary caches have been dropped.
func (l *syncStore) deleteCacheWriterBlobstoreRecords(deletedByNetwork map[string][]string) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer store.Rollback()

	errs := &multierror.Error{}
	for network, deleted := range deletedByNetwork {
		tks := storage.MakeTKs(cacheWriterBlobstoreType, deleted)
		err := store.Delete(network, tks)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "delete blobstore cache writer records %+v for network %+v", deleted, network))
		}
	}
	err = store.Commit()
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs.ErrorOrNil()
}

// generateCacheWriterUUID returns a universally unique ID for a cache writer.
func generateCacheWriterUUID(tableNamePrefix string) string {
	// Replace "-" symbols in the ID since they aren't supported by SQL variable names
	id := strings.Replace((&storage.UUIDGenerator{}).New(), "-", "_", -1)
	return fmt.Sprintf("%s_cache_writer_%s", tableNamePrefix, id)
}

func encodeInt64(n int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(n))
	return bytes
}
