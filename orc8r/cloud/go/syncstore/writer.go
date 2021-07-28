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
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type syncStore struct {
	SyncStoreReader
	cacheWriters map[string]CacheWriter
	db           *sql.DB
	builder      sqorc.StatementBuilder
	fact         blobstore.BlobStorageFactory
}

func NewSyncStore(db *sql.DB, builder sqorc.StatementBuilder, fact blobstore.BlobStorageFactory) SyncStore {
	reader := NewSyncStoreReader(db, builder, fact)
	return &syncStore{SyncStoreReader: reader, cacheWriters: map[string]CacheWriter{}, db: db, builder: builder, fact: fact}
}

func (l *syncStore) CollectGarbage(trackedNetworks []string) error {
	var deletedNetworks []string
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.Select(nidCol).From(digestTableName).RunWith(tx).Query()
		if err != nil {
			return nil, errors.Wrap(err, "get all networks in store")
		}
		inStoreNetworks := []string{}
		for rows.Next() {
			network := ""
			err = rows.Scan(&network)
			if err != nil {
				return nil, errors.Wrap(err, "get all networks in store, SQL rows scan error")
			}
			inStoreNetworks = append(inStoreNetworks, network)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "get all networks in store, SQL rows error")
		}
		// Need to remove all networks that are in store but no longer tracked
		deletedNetworks, _ = funk.DifferenceString(inStoreNetworks, trackedNetworks)

		_, err = l.builder.
			Delete(digestTableName).
			Where(squirrel.Eq{nidCol: deletedNetworks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "delete digests for networks %+v", deletedNetworks)
		}
		_, err = l.builder.
			Delete(cacheTableName).
			Where(squirrel.Eq{nidCol: deletedNetworks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "delete cached objs for networks %+v", deletedNetworks)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return err
	}

	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer store.Rollback()

	errs := &multierror.Error{}
	for _, network := range deletedNetworks {
		filter := blobstore.CreateSearchFilter(&network, []string{lastResyncBlobstoreType, cacheWriterBlobstoreType}, nil, nil)
		criteria := blobstore.LoadCriteria{LoadValue: false}
		networkBlobs, err := store.Search(filter, criteria)
		if err != nil {
			multierror.Append(errs, err)
			continue
		}
		err = store.Delete(network, networkBlobs[network].TKs())
		if err != nil {
			multierror.Append(errs, err)
		}
		for _, tk := range networkBlobs[network].TKs() {
			if tk.Type == cacheWriterBlobstoreType {
				l.cacheWriters[tk.Key].SetInvalid()
				delete(l.cacheWriters, tk.Key)
			}
		}
	}
	if errs.ErrorOrNil() != nil {
		return errors.Wrapf(errs.ErrorOrNil(), "delete info of networks %+v from blobstore", deletedNetworks)
	}

	// A cacheWriter expires after a configurable interval, and is removed from the tracked list
	for _, network := range trackedNetworks {
		filter := blobstore.CreateSearchFilter(&network, []string{cacheWriterBlobstoreType}, nil, nil)
		criteria := blobstore.LoadCriteria{LoadValue: true}
		networkBlobs, err := store.Search(filter, criteria)
		if err != nil {
			multierror.Append(errs, err)
		}
		for _, blob := range networkBlobs[network] {
			creationTime := binary.LittleEndian.Uint64(blob.Value)
			if clock.Now().Unix()-int64(creationTime) > cacheWriterValidIntervalSecs {
				l.cacheWriters[blob.Key].SetInvalid()
				delete(l.cacheWriters, blob.Key)
			}
		}
	}
	return store.Commit()
}

func (l *syncStore) SetDigest(network string, digests *protos.DigestTree) error {
	rootDigest := digests.GetRootDigest().GetMd5Base64Digest()
	leafDigestsToSerialize := &protos.LeafDigestsToSerialize{Digests: digests.GetLeafDigests()}
	leafDigests, err := proto.Marshal(leafDigestsToSerialize)
	if err != nil {
		return errors.Wrapf(err, "marshal leaf digests for network %+v", network)
	}
	now := clock.Now().Unix()

	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Insert(digestTableName).
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
	// the temporary table is namespaced with the unique id of the cacheWriter
	writerId := fmt.Sprintf("cache_writer_%s", strings.Replace((&storage.UUIDGenerator{}).New(), "-", "_", -1))
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(writerId).
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

	// The creation time of cacheWriters is tracked by the store for garbage collection
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	creationTime := make([]byte, 8)
	binary.LittleEndian.PutUint64(creationTime, uint64(clock.Now().Unix()))
	err = store.CreateOrUpdate(network, blobstore.Blobs{{
		Type:  cacheWriterBlobstoreType,
		Key:   writerId,
		Value: creationTime,
	}})
	if err != nil {
		return nil, errors.Wrapf(err, "set start time of network %+v, cachewriter %+v in blobstore", network, writerId)
	}
	err = store.Commit()
	if err != nil {
		return nil, err
	}

	l.cacheWriters[writerId] = NewCacheWriter(network, writerId, l.db, l.builder)
	return l.cacheWriters[writerId], nil
}

func (l *syncStore) RecordResync(network string, gateway string, t uint64) error {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	lastResyncBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastResyncBytes, t)
	err = store.CreateOrUpdate(network, blobstore.Blobs{{
		Type:  lastResyncBlobstoreType,
		Key:   gateway,
		Value: lastResyncBytes,
	}})
	if err != nil {
		return errors.Wrapf(err, "set last resync time of network %+v, gateway %+v in blobstore", network, gateway)
	}

	return store.Commit()
}
