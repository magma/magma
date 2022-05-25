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

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/protos"
)

const (
	// idCol contains the network-wide unique identifiers of the objects.
	idCol              = "id"
	nidCol             = "network_id"
	rootDigestCol      = "root_digest"
	leafDigestsCol     = "leaf_digests"
	objCol             = "obj"
	lastUpdatedTimeCol = "last_updated_at"

	lastResyncBlobstoreType  = "gateway_last_resync_time"
	cacheWriterBlobstoreType = "cache_writer_creation_time"
)

func NewSyncStoreReader(db *sql.DB, builder sqorc.StatementBuilder, fact blobstore.StoreFactory, config Config) (SyncStoreReader, error) {
	err := config.Validate(false)
	if err != nil {
		return nil, fmt.Errorf("invalid configs for syncstore reader: %w", err)
	}
	storeReader := &syncStore{
		db:                           db,
		builder:                      builder,
		fact:                         fact,
		cacheWriterValidIntervalSecs: config.CacheWriterValidIntervalSecs,
		tableNamePrefix:              config.TableNamePrefix,
		digestTableName:              fmt.Sprintf("%s_digest", config.TableNamePrefix),
		cacheTableName:               fmt.Sprintf("%s_cached_objs", config.TableNamePrefix),
	}
	return storeReader, nil
}

func (l *syncStore) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(l.digestTableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(rootDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(leafDigestsCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			Column(lastUpdatedTimeCol).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			PrimaryKey(nidCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("initialize digest store table: %w", err)
		}

		_, err = l.builder.CreateTable(l.cacheTableName).
			IfNotExists().
			Column(nidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(idCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(objCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(nidCol, idCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("initialize cached obj store table: %w", err)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *syncStore) GetDigests(networks []string, lastUpdatedBefore int64, loadLeaves bool) (DigestTrees, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{squirrel.LtOrEq{lastUpdatedTimeCol: lastUpdatedBefore}}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{nidCol: networks})
		}
		rows, err := l.builder.
			Select(nidCol, rootDigestCol, leafDigestsCol, lastUpdatedTimeCol).
			From(l.digestTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, fmt.Errorf("get digests for networks %+v: %w", networks, err)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigests")

		digestTrees := DigestTrees{}
		for rows.Next() {
			network, rootDigest, leafDigestsMarshaled, lastUpdatedTime := "", "", []byte{}, int64(0)
			err = rows.Scan(&network, &rootDigest, &leafDigestsMarshaled, &lastUpdatedTime)
			if err != nil {
				return nil, fmt.Errorf("get digests for network %+v, SQL row scan error: %w", network, err)
			}

			digestTree := &protos.DigestTree{RootDigest: &protos.Digest{Md5Base64Digest: rootDigest}}
			if loadLeaves {
				leafDigests := &protos.LeafDigests{}
				err = proto.Unmarshal(leafDigestsMarshaled, leafDigests)
				if err != nil {
					return nil, fmt.Errorf("unmarshal leaf digests for network %+v: %w", network, err)
				}
				digestTree.LeafDigests = leafDigests.Digests
			}
			digestTrees[network] = digestTree
		}
		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("select digests for network, SQL rows error: %w", err)
		}
		return digestTrees, nil
	}

	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.(DigestTrees)
	return ret, nil
}

func (l *syncStore) GetCachedByID(network string, ids []string) ([][]byte, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(idCol, objCol).
			From(l.cacheTableName).
			Where(
				squirrel.And{
					squirrel.Eq{nidCol: network},
					squirrel.Eq{idCol: ids},
				},
			).
			OrderBy(idCol).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, fmt.Errorf("get cached objs by ID for network %+v: %w", network, err)
		}
		objs, _, err := parseRows(rows)
		return objs, err
	}
	ret, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return ret.([][]byte), err
}

func (l *syncStore) GetCachedByPage(network string, token string, pageSize uint64) ([][]byte, string, error) {
	lastIncludedID := ""
	if token != "" {
		decoded, err := configurator_storage.DeserializePageToken(token)
		if err != nil {
			return nil, "", err
		}
		lastIncludedID = decoded.LastIncludedEntity
	}
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(idCol, objCol).
			From(l.cacheTableName).
			Where(
				squirrel.And{
					squirrel.Eq{nidCol: network},
					squirrel.Gt{idCol: lastIncludedID},
				},
			).
			OrderBy(idCol).
			Limit(pageSize).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, fmt.Errorf("get page for network %+v with token %+v: %w", network, token, err)
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

func (l *syncStore) GetLastResync(network string, gateway string) (int64, error) {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return int64(0), fmt.Errorf("error starting transaction: %w", err)
	}
	defer store.Rollback()

	blob, err := store.Get(network, storage.TK{Type: lastResyncBlobstoreType, Key: gateway})
	if err == merrors.ErrNotFound {
		// If this gw has never been resynced, return 0 to enforce first resync
		return int64(0), nil
	}
	if err != nil {
		return int64(0), fmt.Errorf("get last resync time of network %+v, gateway %+v from blobstore: %w", network, gateway, err)
	}

	lastResync := binary.LittleEndian.Uint64(blob.Value)
	return int64(lastResync), store.Commit()
}

// GetDigestTree returns the full digest tree of a single network.
func GetDigestTree(store SyncStoreReader, network string) (*protos.DigestTree, error) {
	digestTrees, err := store.GetDigests([]string{network}, clock.Now().Unix(), true)
	if err != nil {
		return nil, err
	}
	if _, ok := digestTrees[network]; !ok {
		emptyTree := &protos.DigestTree{
			RootDigest:  &protos.Digest{Md5Base64Digest: ""},
			LeafDigests: []*protos.LeafDigest{},
		}
		return emptyTree, nil
	}
	return digestTrees[network], nil
}

// parsePage parses the list of cached serialized objs, as well as returns the next
// page token based on the lastIncludedEntity in the current page.
//
// NOTE: The configurator_storage.EntityPageToken is used for simplicity & ease
// of transition between loading from configurator to loading from this cache.
// However, this generated token is unrelated to the configurator page tokens.
func parsePage(rows *sql.Rows, pageSize uint64) (*pageInfo, error) {
	objs, lastIncludedID, err := parseRows(rows)
	if err != nil {
		return nil, err
	}

	// Return empty token if we have definitely fetched all data in the db
	if uint64(len(objs)) < pageSize {
		return &pageInfo{token: "", objects: objs}, nil
	}
	nextToken := &configurator_storage.EntityPageToken{LastIncludedEntity: lastIncludedID}
	nextTokenSerialized, err := configurator_storage.SerializePageToken(nextToken)
	if err != nil {
		return nil, fmt.Errorf("get next page token for cached objs store: %w", err)
	}
	return &pageInfo{token: nextTokenSerialized, objects: objs}, nil
}

func parseRows(rows *sql.Rows) ([][]byte, string, error) {
	objs, lastIncludedID := [][]byte{}, ""
	for rows.Next() {
		id, obj := "", []byte{}
		err := rows.Scan(&id, &obj)
		if err != nil {
			return nil, "", fmt.Errorf("get serialized obj, SQL rows scan error: %w", err)
		}
		objs = append(objs, obj)
		lastIncludedID = id
	}
	err := rows.Err()
	if err != nil {
		return nil, "", fmt.Errorf("parse cached objs in store, SQL rows error: %w", err)
	}
	return objs, lastIncludedID, nil
}

type pageInfo struct {
	token   string
	objects [][]byte
}
