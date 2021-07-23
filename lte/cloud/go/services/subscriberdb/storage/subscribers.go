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
	"fmt"

	lte_protos "magma/lte/cloud/go/protos"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	subscriberTableName    = "subscribers"
	subscriberTmpTableName = "subscribers_tmp"
	subscriberNidCol       = "network_id"
	subscriberSidCol       = "subscriber_id"
	subscriberProtoCol     = "subscriber_proto"
)

type SubStore struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewSubStore(db *sql.DB, builder sqorc.StatementBuilder) *SubStore {
	return &SubStore{db: db, builder: builder}
}

func (l *SubStore) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(subscriberTableName).
			IfNotExists().
			Column(subscriberNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberSidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberProtoCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(subscriberNidCol, subscriberSidCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "initialize sub proto store table")
		}
		_, err = l.builder.CreateTable(subscriberTmpTableName).
			IfNotExists().
			Column(subscriberNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberSidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subscriberProtoCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(subscriberNidCol, subscriberSidCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "initialize sub proto store tmp table")
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// InitializeUpdate prepares the db tables for a batch update.
func (l *SubStore) InitializeUpdate() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(subscriberTmpTableName).
			RunWith(tx).
			Exec()
		return nil, errors.Wrap(err, "clear sub protos tmp table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// InsertMany inserts a batch of subscriber data into the temporary table.
// NOTE: Caller of the function should enforce that the max size of the
// insertion aligns reasonably with the max page size of its corresponding
// load source.
func (l *SubStore) InsertMany(network string, subProtos []*lte_protos.SubscriberData) error {
	if len(subProtos) == 0 {
		return nil
	}
	insertQuery := l.builder.
		Insert(subscriberTmpTableName).
		Columns(subscriberNidCol, subscriberSidCol, subscriberProtoCol)
	errs := &multierror.Error{}
	for _, subProto := range subProtos {
		marshaled, err := proto.Marshal(subProto)
		if err != nil {
			multierror.Append(errs, errors.Wrapf(err, "serialize subproto of ID %+v", lte_protos.SidString(subProto.Sid)))
			continue
		}
		insertQuery = insertQuery.Values(network, lte_protos.SidString(subProto.Sid), marshaled)
	}
	if errs.ErrorOrNil() != nil {
		return errs
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := insertQuery.RunWith(tx).Exec()
		return nil, errors.Wrapf(err, "insert sub protos into store for network %+v", network)
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// ApplyUpdate applies all subscriber data updates onto the db table and
// completes the batch update.
func (l *SubStore) ApplyUpdate(network string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		// HACK: hard coding part of this sql query because there currently doesn't exist good support
		// for "WHERE (row NOT IN other_table)" with squirrel
		//
		// The SQL query should look something like
		// DELETE FROM subscribers WHERE
		//     network_id NOT IN ${networks}
		// AND
		//	   (network_id, subscriber_id) NOT IN (SELECT network_id, subscriber_id FROM subscribers_tmp)
		_, err := l.builder.
			Delete(subscriberTableName).
			Where(squirrel.And{
				squirrel.Eq{subscriberNidCol: network},
				squirrel.Expr(fmt.Sprintf(
					"(%s, %s) NOT IN (SELECT %s, %s FROM %s)",
					subscriberNidCol, subscriberSidCol,
					subscriberNidCol, subscriberSidCol, subscriberTmpTableName,
				)),
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "clean up previous sub proto store table")
		}

		// The upsert query should look something like
		// INSERT INTO subscribers
		//     SELECT network_id, subscriber_id, subscriber_proto FROM subscribers_tmp
		// 	   WHERE network_id = ${network}
		// ON CONFLICT (network_id, subscriber_id)
		// 	   DO UPDATE SET subscriber_proto = subscribers_tmp.subscriber_proto
		conflictUpdateTarget := sqorc.FmtConflictUpdateTarget(subscriberTmpTableName, subscriberProtoCol)
		_, err = l.builder.
			Insert(subscriberTableName).
			Select(
				l.builder.
					Select(subscriberNidCol, subscriberSidCol, subscriberProtoCol).
					From(subscriberTmpTableName).
					Where(squirrel.Eq{subscriberNidCol: network}),
			).
			OnConflict(
				[]sqorc.UpsertValue{{
					Column: subscriberProtoCol,
					Value:  squirrel.Expr(conflictUpdateTarget),
				}},
				subscriberNidCol, subscriberSidCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "populate sub proto store table")
		}

		_, err = l.builder.
			Delete(subscriberTmpTableName).
			Where(squirrel.Eq{subscriberNidCol: network}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "clean up tmp sub proto store table")
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// DeleteSubscribersForNetworks deletes the cached protos for a list of networks.
func (l *SubStore) DeleteSubscribersForNetworks(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(subscriberTableName).
			Where(squirrel.Eq{subscriberNidCol: networks}).
			RunWith(tx).
			Exec()
		return nil, errors.Wrapf(err, "delete sub protos for networks %+v", networks)
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// GetSubscribersPage gets a page of subscriber protos based on the page token
// and size, and also returns the token for the next page.
func (l *SubStore) GetSubscribersPage(network string, token string, pageSize uint64) ([]*lte_protos.SubscriberData, string, error) {
	lastIncludedSid := ""
	if token != "" {
		decoded, err := configurator_storage.DeserializePageToken(token)
		if err != nil {
			return nil, "", err
		}
		lastIncludedSid = decoded.LastIncludedEntity
	}
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(subscriberProtoCol).
			From(subscriberTableName).
			Where(squirrel.And{
				squirrel.Eq{subscriberNidCol: network},
				squirrel.Gt{subscriberSidCol: lastIncludedSid},
			}).
			OrderBy(subscriberSidCol).
			Limit(pageSize).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "get page for network %+v with token %+v", network, token)
		}
		return parseSubProtoRows(rows)
	}
	ret, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, "", err
	}
	subProtos := ret.([]*lte_protos.SubscriberData)
	nextPageToken, err := getNextPageToken(subProtos, pageSize)
	if err != nil {
		return nil, "", errors.Wrap(err, "get next page token")
	}
	return subProtos, nextPageToken, nil
}

// GetSubscribers returns an ordered list of subscriber protos with matching IDs.
// NOTE: Caller of the function should enforce that the max size of the requested
// subscribers aligns reasonably with its max page size.
func (l *SubStore) GetSubscribers(network string, sids []string) ([]*lte_protos.SubscriberData, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(subscriberProtoCol).
			From(subscriberTableName).
			Where(squirrel.And{
				squirrel.Eq{subscriberNidCol: network},
				squirrel.Eq{subscriberSidCol: sids},
			}).
			OrderBy(subscriberSidCol).
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrapf(err, "Get sub protos by ID for network %+v", network)
		}
		return parseSubProtoRows(rows)
	}
	ret, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return nil, err
	}
	subProtos := ret.([]*lte_protos.SubscriberData)
	return subProtos, nil
}

// parseSubProtoRows scans db rows of serialized subscriber data, and returns
// a deserialized list of subscriber protos.
func parseSubProtoRows(rows *sql.Rows) ([]*lte_protos.SubscriberData, error) {
	subProtos := []*lte_protos.SubscriberData{}
	for rows.Next() {
		subProtoMarshaled := []byte{}
		err := rows.Scan(&subProtoMarshaled)
		if err != nil {
			return nil, errors.Wrap(err, "get sub protos from store, SQL row scan error")
		}
		subProto := &lte_protos.SubscriberData{}
		err = proto.Unmarshal(subProtoMarshaled, subProto)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal sub protos from store")
		}
		subProtos = append(subProtos, subProto)
	}
	err := rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "get sub protos from store, SQL rows error")
	}
	return subProtos, nil
}

// getNextPageToken returns the next page token based on the lastIncludedEntity
// in the current page.
// NOTE: The configurator_storage.EntityPageToken is used for simplicity & ease
// of transition between loading from configurator to loading from this cache.
// However, this generated token is unrelated to the configurator page tokens.
func getNextPageToken(subProtos []*lte_protos.SubscriberData, pageSize uint64) (string, error) {
	// The next token is empty if we have definitely exhausted all protos in the db
	if uint64(len(subProtos)) < pageSize {
		return "", nil
	}
	lastSubProto := subProtos[len(subProtos)-1]
	nextToken := &configurator_storage.EntityPageToken{
		LastIncludedEntity: lte_protos.SidString(lastSubProto.Sid),
	}
	return configurator_storage.SerializePageToken(nextToken)
}
