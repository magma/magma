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
	"encoding/base64"

	lte_protos "magma/lte/cloud/go/protos"
	configurator_storage "magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	subProtoTableName    = "subscriber_protos"
	subProtoTmpTableName = "subscriber_protos_tmp"

	subProtoNidCol   = "network_id"
	subProtoSidCol   = "subscriber_id"
	subProtoProtoCol = "subscriber_proto"
)

type SubProtoStore struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewSubProtoStore(db *sql.DB, builder sqorc.StatementBuilder) *SubProtoStore {
	return &SubProtoStore{db: db, builder: builder}
}

func (l *SubProtoStore) Initialize() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.CreateTable(subProtoTableName).
			IfNotExists().
			Column(subProtoNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subProtoSidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subProtoProtoCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(subProtoNidCol, subProtoSidCol).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "initialize sub proto store table")
		}

		_, err = l.builder.CreateTable(subProtoTmpTableName).
			IfNotExists().
			Column(subProtoNidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subProtoSidCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
			Column(subProtoProtoCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
			PrimaryKey(subProtoNidCol, subProtoSidCol).
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

func (l *SubProtoStore) InsertManyByNetwork(network string, subProtos []*lte_protos.SubscriberData) error {
	if len(subProtos) == 0 {
		return nil
	}

	insertQuery := l.builder.
		Insert(subProtoTmpTableName).
		Columns(subProtoNidCol, subProtoSidCol, subProtoProtoCol)
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

// CommitUpdateByNetwork cleans up and re-populates the subscriber proto store table with data from the temporary
// table for a particular network, and then truncates the temporary table.
func (l *SubProtoStore) CommitUpdateByNetwork(network string) error {
	// TODO(wangyyt1013): optimize this for large tables
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(subProtoTableName).
			Where(squirrel.Eq{subProtoNidCol: network}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "clean up previous sub proto store table")
		}

		_, err = l.builder.
			Insert(subProtoTableName).
			Select(
				l.builder.
					Select(subProtoNidCol, subProtoSidCol, subProtoProtoCol).
					From(subProtoTmpTableName).
					Where(squirrel.Eq{subProtoNidCol: network}),
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "populate sub proto store table")
		}

		_, err = l.builder.
			Delete(subProtoTmpTableName).
			Where(squirrel.Eq{subProtoNidCol: network}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrapf(err, "clean up tmp sub proto store table")
		}

		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *SubProtoStore) ClearTmpTable() error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(subProtoTmpTableName).
			RunWith(tx).
			Exec()
		return nil, errors.Wrapf(err, "clear sub protos tmp table")
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *SubProtoStore) DeleteSubProtos(networks []string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := l.builder.
			Delete(subProtoTableName).
			Where(squirrel.Eq{subProtoNidCol: networks}).
			RunWith(tx).
			Exec()
		return nil, errors.Wrapf(err, "delete sub protos for networks %+v", networks)
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

// GetPage gets a page of subscriber protos based on the page token and size, and also returns
// the token for the next page.
func (l *SubProtoStore) GetPage(network string, token string, pageSize int) ([]*lte_protos.SubscriberData, string, error) {
	lastIncludedSid := ""
	if token != "" {
		marshaledToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, "", errors.Wrapf(err, "decode page token")
		}

		buf := &configurator_storage.EntityPageToken{}
		err = proto.Unmarshal(marshaledToken, buf)
		if err != nil {
			return nil, "", errors.Wrapf(err, "unmarshal page token")
		}
		lastIncludedSid = buf.LastIncludedEntity
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(subProtoProtoCol).
			From(subProtoTableName).
			Where(squirrel.And{
				squirrel.Eq{subProtoNidCol: network},
				squirrel.Gt{subProtoSidCol: lastIncludedSid},
			}).
			OrderBy(subProtoSidCol).
			Limit(uint64(pageSize)).
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
	nextPageToken, err := getNextPageToken(subProtos)
	if err != nil {
		return nil, "", errors.Wrapf(err, "get next page token")
	}
	return subProtos, nextPageToken, nil
}

func (l *SubProtoStore) GetByIDs(network string, sids []string) ([]*lte_protos.SubscriberData, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		rows, err := l.builder.
			Select(subProtoProtoCol).
			From(subProtoTableName).
			Where(squirrel.And{
				squirrel.Eq{subProtoNidCol: network},
				squirrel.Eq{subProtoSidCol: sids},
			}).
			OrderBy(subProtoSidCol).
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

// parseSubProtoRows scans db rows of serialized subscriber data, and returns a deserialized list
// of subscriber protos.
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
			return nil, errors.Wrapf(err, "unmarshal sub protos from store")
		}
		subProtos = append(subProtos, subProto)
	}
	err := rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "get sub protos from store, SQL rows error")
	}
	return subProtos, nil
}

// getNextPageToken returns the next page token based on the lastIncludedEntity in the current page.
func getNextPageToken(subProtos []*lte_protos.SubscriberData) (string, error) {
	nextToken := ""
	if len(subProtos) > 0 {
		lastSubProto := subProtos[len(subProtos)-1]
		nextTokenUnmarshaled := &configurator_storage.EntityPageToken{
			LastIncludedEntity: lte_protos.SidString(lastSubProto.Sid),
		}
		nextTokenMarshaled, err := proto.Marshal(nextTokenUnmarshaled)
		if err != nil {
			return "", err
		}
		nextToken = base64.StdEncoding.EncodeToString(nextTokenMarshaled)
	}
	return nextToken, nil
}
