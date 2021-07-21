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

package storage

import (
	"database/sql"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type digestStore struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

type DigestStore interface {
	// Initialize the backing store.
	Initialize() error

	// GetDigests returns a list of digests that satisfy the filtering criteria
	// specified by the arguments.
	// Caveats:
	// 1. If networks is empty, returns digests for all networks.
	// 2. lastUpdatedBefore is recorded in unix seconds. Filters for all digests that
	// were last updated earlier than this time.
	GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error)

	// SetDigest creates/updates the subscribers digest for a particular network.
	SetDigest(network string, digest string) error

	// DeleteDigests removes digests by network IDs.
	DeleteDigests(networks []string) error
}

const (
	digestTableName = "subscriberdb_flat_digests"

	digestNidCol             = "network_id"
	digestDigestCol          = "digest"
	digestLastUpdatedTimeCol = "last_updated_at"
)

func NewDigestStore(db *sql.DB, builder sqorc.StatementBuilder) DigestStore {
	return &digestStore{db: db, builder: builder}
}

func (l *digestStore) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(digestTableName).
			IfNotExists().
			Column(digestNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(digestDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(digestLastUpdatedTimeCol).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			PrimaryKey(digestNidCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize digest store table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *digestStore) GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{squirrel.LtOrEq{digestLastUpdatedTimeCol: lastUpdatedBefore}}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{digestNidCol: networks})
		}

		rows, err := l.builder.
			Select(digestNidCol, digestDigestCol, digestLastUpdatedTimeCol).
			From(digestTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "gets digest for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigest")

		digestInfos := DigestInfos{}
		for rows.Next() {
			network, digest, lastUpdatedTime := "", "", int64(0)
			err = rows.Scan(&network, &digest, &lastUpdatedTime)
			if err != nil {
				return nil, errors.Wrap(err, "select digests for networks, SQL row scan error")
			}
			digestInfo := DigestInfo{
				Network:         network,
				Digest:          digest,
				LastUpdatedTime: lastUpdatedTime,
			}
			digestInfos = append(digestInfos, digestInfo)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select digests for network, SQL rows error")
		}
		return digestInfos, nil
	}

	txRet, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	ret := txRet.(DigestInfos)
	return ret, nil
}

func (l *digestStore) SetDigest(network string, digest string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		now := clock.Now().Unix()
		_, err := l.builder.
			Insert(digestTableName).
			Columns(digestNidCol, digestDigestCol, digestLastUpdatedTimeCol).
			Values(network, digest, now).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: digestDigestCol, Value: digest},
					{Column: digestLastUpdatedTimeCol, Value: now},
				},
				digestNidCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert digest for network %+v", network)
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *digestStore) DeleteDigests(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(digestTableName).
			Where(squirrel.Eq{digestNidCol: networks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "delete digests")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// GetDigest returns the digest information for a particular network.
func GetDigest(l DigestStore, network string) (string, error) {
	digestInfos, err := l.GetDigests([]string{network}, clock.Now().Unix())
	if err != nil {
		return "", err
	}
	// There should be at most 1 digest for each network
	// if digest not found, return default value
	if len(digestInfos) == 0 {
		return "", nil
	}
	return digestInfos[0].Digest, nil
}

// GetOutdatedNetworks returns all networks with digests last updated at a time
// earlier than the specified deadline.
func GetOutdatedNetworks(l DigestStore, lastUpdatedBefore int64) ([]string, error) {
	digestInfos, err := l.GetDigests([]string{}, lastUpdatedBefore)
	if err != nil {
		return nil, err
	}
	networks := digestInfos.Networks()
	return networks, nil
}

// GetAllNetworks returns all unique networks currently stored.
func GetAllNetworks(l DigestStore) ([]string, error) {
	digestInfos, err := l.GetDigests([]string{}, clock.Now().Unix())
	if err != nil {
		return nil, err
	}
	networks := digestInfos.Networks()
	networksUniq := funk.UniqString(networks)
	return networksUniq, nil
}

type DigestInfo struct {
	Network         string
	Digest          string
	LastUpdatedTime int64
}

type DigestInfos []DigestInfo

// Networks returns a list of network IDs for all digests in DigestInfos.
func (digestInfos DigestInfos) Networks() []string {
	ret := []string{}
	for _, digestInfo := range digestInfos {
		ret = append(ret, digestInfo.Network)
	}
	return ret
}
