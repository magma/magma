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
)

type perSubDigestLookup struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

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

func (l *perSubDigestLookup) GetDigest(network string) (interface{}, error) {
	digestInfos, err := l.GetDigests([]string{network}, clock.Now().Unix())
	if err != nil {
		return nil, err
	}
	digestsBySubscriber := map[string]string{}
	for _, digestInfo := range digestInfos {
		digestsBySubscriber[digestInfo.Subscriber] = digestInfo.Digest
	}
	return digestsBySubscriber, nil
}

func (l *perSubDigestLookup) GetDigests(networks []string, lastUpdatedBefore int64) (DigestInfos, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		filters := squirrel.And{}
		if len(networks) > 0 {
			filters = append(filters, squirrel.Eq{flatDigestNidCol: networks})
		}
		rows, err := l.builder.
			Select(perSubDigestNidCol, perSubDigestSidCol, perSubDigestDigestCol).
			From(perSubDigestTableName).
			Where(filters).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "gets per sub digest for networks %+v", networks)
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetDigest")

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

func (l *perSubDigestLookup) SetDigest(network string, subscriber string, digest string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
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
			return nil, errors.Wrapf(err, "insert sub digest for network %+v and subscriber %+v", network, subscriber)
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
