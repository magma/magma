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

	"magma/orc8r/cloud/go/blobstore"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type syncStoreReader struct {
	db         *sql.DB
	builder    sqorc.StatementBuilder
	resyncFact blobstore.BlobStorageFactory
}

func NewSyncStoreReader(db *sql.DB, builder sqorc.StatementBuilder, resyncFact blobstore.BlobStorageFactory) SyncStoreReader {
	return &syncStoreReader{db: db, builder: builder, resyncFact: resyncFact}
}

func (l *syncStoreReader) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(digestTableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(rootDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(leafDigestsCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			Column(lastUpdatedTimeCol).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			PrimaryKey(nidCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "initialize digest store table")
		}

		_, err = l.builder.CreateTable(cacheTableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(idCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(objCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(nidCol, idCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "initialize cached obj store table")
		}

		_, err = l.builder.CreateTable(cacheTmpTableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(idCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(objCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(nidCol, idCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize cached obj store tmp table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *syncStoreReader) GetDigests(networks []string, lastUpdatedBefore int64, loadLeaves bool) (map[string]*protos.DigestTree, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{squirrel.LtOrEq{lastUpdatedTimeCol: lastUpdatedBefore}}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{nidCol: networks})
		}
		rows, err := l.builder.
			Select(nidCol, rootDigestCol, leafDigestsCol, lastUpdatedTimeCol).
			From(digestTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "get digests for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigests")

		digestTrees := map[string]*protos.DigestTree{}
		for rows.Next() {
			network, rootDigest, leafDigestsMarshaled, lastUpdatedTime := "", "", []byte{}, int64(0)
			err = rows.Scan(&network, &rootDigest, &leafDigestsMarshaled, &lastUpdatedTime)
			if err != nil {
				return nil, errors.Wrapf(err, "get digests for network %+v, SQL row scan error", network)
			}

			digestTree := &protos.DigestTree{
				RootDigest: &protos.Digest{Md5Base64Digest: rootDigest},
			}
			if loadLeaves {
				leafDigests := &protos.LeafDigestsToSerialize{}
				err = proto.Unmarshal(leafDigestsMarshaled, leafDigests)
				if err != nil {
					return nil, errors.Wrapf(err, "unmarshal leaf digests for network %+v", network)
				}
				digestTree.LeafDigests = leafDigests.GetDigests()
			}
			digestTrees[network] = digestTree
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select digests for network, SQL rows error")
		}
		return digestTrees, nil
	}

	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.(map[string]*protos.DigestTree)
	return ret, nil
}

func (l *syncStoreReader) GetCachedByID(network string, ids []string) ([][]byte, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(idCol, objCol).
			From(cacheTableName).
			Where(squirrel.And{
				squirrel.Eq{nidCol: network},
				squirrel.Eq{idCol: ids},
			}).
			OrderBy(idCol).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "get cached objs by ID for network %+v", network)
		}
		objs, _, err := parseRows(rows)
		return objs, err
	}
	ret, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return ret.([][]byte), err
}

func (l *syncStoreReader) GetCachedByPage(network string, token string, pageSize uint64) ([][]byte, string, error) {
	lastIncludedId := ""
	if token != "" {
		decoded, err := configurator_storage.DeserializePageToken(token)
		if err != nil {
			return nil, "", err
		}
		lastIncludedId = decoded.LastIncludedEntity
	}
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(idCol, objCol).
			From(cacheTableName).
			Where(squirrel.And{
				squirrel.Eq{nidCol: network},
				squirrel.Gt{idCol: lastIncludedId},
			}).
			OrderBy(idCol).
			Limit(pageSize).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "get page for network %+v with token %+v", network, token)
		}
		return parsePage(rows, pageSize)
	}
	ret, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, "", err
	}
	info := ret.(*pageInfo)
	return info.objects, info.token, nil
}

func (l *syncStoreReader) GetLastResync(network string, gateway string) (uint32, error) {
	store, err := l.resyncFact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return uint32(0), errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(network, storage.TypeAndKey{Type: lastResyncBlobstoreType, Key: gateway})
	if err == merrors.ErrNotFound {
		// If this gw has never been resynced, return 0 to enforce first resync
		return uint32(0), nil
	}
	if err != nil {
		return uint32(0), errors.Wrapf(err, "get last resync time of network %+v, gateway %+v from blobstore", network, gateway)
	}

	lastResync := binary.LittleEndian.Uint32(blob.Value)
	return lastResync, store.Commit()
}

// parsePage parses the list of cached serialized objs, as well as returns the next
// page token based on the lastIncludedEntity in the current page.
// NOTE: The configurator_storage.EntityPageToken is used for simplicity & ease
// of transition between loading from configurator to loading from this cache.
// However, this generated token is unrelated to the configurator page tokens.
func parsePage(rows *sql.Rows, pageSize uint64) (*pageInfo, error) {
	objs, lastIncludedId, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	// Return empty token if we have definitely fetched all data in the db
	if uint64(len(objs)) < pageSize {
		return &pageInfo{token: "", objects: objs}, nil
	}
	nextToken, err := configurator_storage.SerializePageToken(&configurator_storage.EntityPageToken{
		LastIncludedEntity: lastIncludedId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get next page token for cached objs store")
	}
	return &pageInfo{token: nextToken, objects: objs}, nil
}

func parseRows(rows *sql.Rows) ([][]byte, string, error) {
	objs, lastIncludedId := [][]byte{}, ""
	for rows.Next() {
		id, obj := "", []byte{}
		err := rows.Scan(&id, &obj)
		if err != nil {
			return nil, "", errors.Wrap(err, "get serialized obj, SQL rows scan error")
		}
		objs = append(objs, obj)
		lastIncludedId = id
	}
	return objs, lastIncludedId, nil
}

type pageInfo struct {
	token   string
	objects [][]byte
}
