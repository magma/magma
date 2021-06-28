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

type DigestLookup interface {
	// Initialize the backing store.
	Initialize() error

	// GetDigest returns the latest flat digest/per sub digests of a particular network.
	GetDigest(network string) (interface{}, error)

	// GetDigests returns a list of digests that satisfy the filtering criteria
	// specified by the arguments.
	// Caveats:
	// 1. If networks is empty, returns digests for all networks.
	// 2. lastUpdatedBefore is recorded in unix seconds. Filters for all digests that
	// were last updated earlier than this time. Only applied when querying for flat
	// digests (the source of truth for whether a network's digests are up-to-date).
	GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error)

	// SetDigest creates/updates the digest for a particular network/subscriber.
	SetDigest(network string, subscriber string, digest string) error

	// DeleteDigests removes digests by network IDs.
	DeleteDigests(networks []string) error
}

type flatDigestLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

const (
	flatDigestTableName = "subscriberdb_flat_digests"

	flatDigestNidCol             = "network_id"
	flatDigestDigestCol          = "digest"
	flatDigestLastUpdatedTimeCol = "last_updated_at"

	perSubDigestTableName = "subscriberdb_per_sub_digests"

	perSubDigestNidCol    = "network_id"
	perSubDigestSidCol    = "susbcriber_id"
	perSubDigestDigestCol = "digest"
)

func NewFlatDigestLookup(db *sql.DB, builder sqorc.StatementBuilder) DigestLookup {
	return &flatDigestLookup{db: db, builder: builder}
}

func (l *flatDigestLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(flatDigestTableName).
			IfNotExists().
			Column(flatDigestNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(flatDigestDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(flatDigestLastUpdatedTimeCol).Type(sqorc.ColumnTypeBigInt).NotNull().EndColumn().
			PrimaryKey(flatDigestNidCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize flat digest lookup table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *flatDigestLookup) GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{squirrel.LtOrEq{flatDigestLastUpdatedTimeCol: lastUpdatedBefore}}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{flatDigestNidCol: networks})
		}

		rows, err := l.builder.
			Select(flatDigestNidCol, flatDigestDigestCol, flatDigestLastUpdatedTimeCol).
			From(flatDigestTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "gets flat digest for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigest")

		digestInfos := DigestInfos{}
		for rows.Next() {
			network, digest, lastUpdatedTime := "", "", int64(0)
			err = rows.Scan(&network, &digest, &lastUpdatedTime)
			if err != nil {
				return nil, errors.Wrap(err, "select flat digests for networks, SQL row scan error")
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
			return nil, errors.Wrap(err, "select flat digests for network, SQL rows error")
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

func (l *flatDigestLookup) SetDigest(network string, subscriber string, digest string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		now := clock.Now().Unix()
		_, err := l.builder.
			Insert(flatDigestTableName).
			Columns(flatDigestNidCol, flatDigestDigestCol, flatDigestLastUpdatedTimeCol).
			Values(network, digest, now).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: flatDigestDigestCol, Value: digest},
					{Column: flatDigestLastUpdatedTimeCol, Value: now},
				},
				flatDigestNidCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "insert flat digest for network %+v", network)
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *flatDigestLookup) DeleteDigests(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(flatDigestTableName).
			Where(squirrel.Eq{flatDigestNidCol: networks}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "delete digests")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *flatDigestLookup) GetDigest(network string) (interface{}, error) {
	digestInfos, err := l.GetDigests([]string{network}, clock.Now().Unix())
	if err != nil {
		return nil, err
	}
	// There should be at most 1 digest for each network
	// if digest not found, return default value
	if len(digestInfos) == 0 {
		return DigestInfo{}, nil
	}
	digestInfo := digestInfos[0]
	return digestInfo, nil
}

// GetOutdatedNetworks returns all networks with digests last updated at a time
// earlier than the specified deadline.
func GetOutdatedNetworks(l DigestLookup, lastUpdatedBefore int64) ([]string, error) {
	digestInfos, err := l.GetDigests([]string{}, lastUpdatedBefore)
	if err != nil {
		return nil, err
	}
	networks := digestInfos.Networks()
	return networks, nil
}

// GetAllNetworks returns all unique networks currently stored.
func GetAllNetworks(l DigestLookup) ([]string, error) {
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
	Subscriber      string
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
