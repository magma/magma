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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type syncStore struct {
	SyncStoreReader
	db      *sql.DB
	builder sqorc.StatementBuilder
	fact    blobstore.BlobStorageFactory
}

func NewSyncStore(db *sql.DB, builder sqorc.StatementBuilder, fact blobstore.BlobStorageFactory) SyncStore {
	reader := NewSyncStoreReader(db, builder, fact)
	return &syncStore{SyncStoreReader: reader, db: db, builder: builder, fact: fact}
}

func (l *syncStore) CollectGarbage(trackedNetworks []string) error {
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
		// Need to remove all networks that are in store but no longer tracked
		deletedNetworks, _ := funk.DifferenceString(inStoreNetworks, trackedNetworks)

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
		// TODO(wangyyt1013): collect garbage in last resync time store.
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
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
	// prepare the rows associated with the network for a batch update.
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(cacheTmpTableName).
			Where(squirrel.Eq{nidCol: network}).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "clear cached objs tmp table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}

	// TODO(wangyyt1013): add concurrency protection
	return NewCacheWriter(network, l.db, l.builder), nil
}

func (l *syncStore) RecordResync(network string, gatewayID string, t int64) error {
	// TODO(wangyyt1013)
	return nil
}
