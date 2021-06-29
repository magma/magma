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

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type perSubDigestLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

const (
	perSubDigestTableName = "subscriberdb_per_sub_digests"

	perSubDigestNidCol    = "network_id"
	perSubDigestSidCol    = "susbcriber_id"
	perSubDigestDigestCol = "digest"
)

func NewPerSubDigestLookup(db *sql.DB, builder sqorc.StatementBuilder) DigestLookup {
	return &perSubDigestLookup{db: db, builder: builder}
}

func (l *perSubDigestLookup) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(perSubDigestTableName).
			IfNotExists().
			Column(perSubDigestNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(perSubDigestSidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(perSubDigestDigestCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			PrimaryKey(perSubDigestNidCol, perSubDigestSidCol).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "initialize per sub digest lookup table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// GetDigests for PerSubDigestLookup returns a list of digests ordered by their network ID and subscriber ID.
func (l *perSubDigestLookup) GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{perSubDigestNidCol: networks})
		}
		rows, err := l.builder.
			Select(perSubDigestNidCol, perSubDigestSidCol, perSubDigestDigestCol).
			From(perSubDigestTableName).
			Where(filters).
			OrderBy(perSubDigestNidCol, perSubDigestSidCol).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "gets per sub digest for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigests")

		digestInfos := DigestInfos{}
		for rows.Next() {
			network, digest, subscriber := "", "", ""
			err = rows.Scan(&network, &subscriber, &digest)
			if err != nil {
				return nil, errors.Wrap(err, "select per sub digests for networks, SQL row scan error")
			}
			digestInfo := DigestInfo{
				Network:    network,
				Subscriber: subscriber,
				Digest:     digest,
			}
			digestInfos = append(digestInfos, digestInfo)
		}
		err = rows.Err()
		if err != nil {
			return nil, errors.Wrap(err, "select per sub digests for network, SQL rows error")
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

func (l *perSubDigestLookup) SetDigest(network string, args interface{}) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		perSubDigestUpsertArgs, ok := args.(PerSubDigestUpsertArgs)
		if !ok {
			return nil, errors.Errorf("invalid args for setting flat digest of network %+v", network)
		}
		toRenew := perSubDigestUpsertArgs.ToRenew
		deleted := perSubDigestUpsertArgs.Deleted

		for _, subscriber := range deleted {
			_, err := l.builder.
				Delete(perSubDigestTableName).
				Where(squirrel.And{
					squirrel.Eq{perSubDigestNidCol: network},
					squirrel.Eq{perSubDigestSidCol: subscriber},
				}).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, errors.Wrapf(err, "delete digest of subscriber %+v of network %+v", subscriber, network)
			}
		}
		for subscriber, digest := range toRenew {
			_, err := l.builder.
				Insert(perSubDigestTableName).
				Columns(perSubDigestNidCol, perSubDigestSidCol, perSubDigestDigestCol).
				Values(network, subscriber, digest).
				OnConflict(
					[]sqorc.UpsertValue{
						{Column: perSubDigestDigestCol, Value: digest},
					},
					perSubDigestNidCol, perSubDigestSidCol,
				).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, errors.Wrapf(err, "insert sub digest for subscriber %+v of network %+v", subscriber, network)
			}
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *perSubDigestLookup) DeleteDigests(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(perSubDigestTableName).
			Where(squirrel.Eq{perSubDigestNidCol: networks}).
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

// PerSubDigestUpsertArgs specifies for a certain network:
// 1. subscribers whose digests need to be added/updated in store
// 2. subscribers who need to be removed from the store.
type PerSubDigestUpsertArgs struct {
	ToRenew map[string]string
	Deleted []string
}
