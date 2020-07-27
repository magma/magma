/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	nodeTable = "testcontroller_nodes"

	idCol         = "pk"
	vpnIPCol      = "vpn_ip"
	tagCol        = "tag"
	availCol      = "available"
	lastLeasedCol = "last_leased_sec"
	leaseIdCol    = "lease_id"
)

// Based on current AWS workload runtime of 1hr+, a 2hr timeout should be
// pretty generous but not overly so
const leaseTimeout = 2 * time.Hour

const manualLeaseID = "manual"

var (
	selectedNextNode = func() {}
)

type sqlNodeLeasorStorage struct {
	db          *sql.DB
	idGenerator storage.IDGenerator
	builder     sqorc.StatementBuilder
}

func NewSQLNodeLeasorStorage(db *sql.DB, idGenerator storage.IDGenerator, builder sqorc.StatementBuilder) NodeLeasorStorage {
	return &sqlNodeLeasorStorage{db: db, idGenerator: idGenerator, builder: builder}
}

func (s *sqlNodeLeasorStorage) Init() (err error) {
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("%s; rollback error: %s", err, rollbackErr)
			}
		}
	}()

	_, err = s.builder.CreateTable(nodeTable).
		IfNotExists().
		Column(idCol).Type(sqorc.ColumnTypeText).PrimaryKey().EndColumn().
		Column(vpnIPCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(tagCol).Type(sqorc.ColumnTypeText).NotNull().Default("''").EndColumn().
		Column(availCol).Type(sqorc.ColumnTypeBool).NotNull().Default(true).EndColumn().
		Column(lastLeasedCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Column(leaseIdCol).Type(sqorc.ColumnTypeText).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create worker node lease table")
	}

	// Add tag column here. We can remove this code after it runs on prod once.
	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s text NOT NULL DEFAULT ''", nodeTable, tagCol))
	// sqlite doesn't support this alter table DDL but it doesn't matter since
	// we only use ephemeral in-memory sqlite for tests
	if err != nil && os.Getenv("SQL_DRIVER") != "sqlite3" {
		return errors.Wrap(err, "failed to add 'tag' column to nodes table")
	}
	return
}

func (s *sqlNodeLeasorStorage) GetNodes(ids []string, tag *string) (map[string]*CINode, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		builder := s.builder.Select(idCol, vpnIPCol, tagCol, availCol, lastLeasedCol).
			From(nodeTable).
			RunWith(tx)
		clauses := squirrel.And{}
		if !funk.IsEmpty(ids) {
			clauses = append(clauses, squirrel.Eq{idCol: ids})
		}
		if tag != nil {
			clauses = append(clauses, squirrel.Eq{tagCol: tag})
		}
		if !funk.IsEmpty(clauses) {
			builder = builder.Where(clauses)
		}

		rows, err := builder.Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve nodes")
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetNodes")
		return scanNodes(rows)
	}

	ret, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return map[string]*CINode{}, err
	}
	return ret.(map[string]*CINode), nil
}

func (s *sqlNodeLeasorStorage) CreateOrUpdateNode(node *MutableCINode) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.Insert(nodeTable).
			Columns(idCol, tagCol, vpnIPCol).
			Values(node.Id, node.Tag, node.VpnIP).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: tagCol, Value: node.Tag},
					{Column: vpnIPCol, Value: node.VpnIP},
				},
				idCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to write update to node")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlNodeLeasorStorage) DeleteNode(id string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.Delete(nodeTable).
			Where(squirrel.Eq{idCol: id}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete node")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlNodeLeasorStorage) LeaseNode(tag *string) (*NodeLease, error) {
	ret, err := sqorc.ExecInTx(s.db, nil, nil, s.getNodeLeaseTxFn("", tag, ""))
	switch {
	case err != nil:
		return nil, err
	case ret == nil:
		return nil, nil
	default:
		return ret.(*NodeLease), nil
	}
}

func (s *sqlNodeLeasorStorage) ReserveNode(id string) (*NodeLease, error) {
	ret, err := sqorc.ExecInTx(s.db, nil, nil, s.getNodeLeaseTxFn(id, nil, manualLeaseID))
	switch {
	case err != nil:
		return nil, err
	case ret == nil:
		return nil, nil
	default:
		return ret.(*NodeLease), nil
	}
}

func (s *sqlNodeLeasorStorage) ReleaseNode(id string, leaseID string) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		res, err := s.builder.Update(nodeTable).
			Set(availCol, true).
			Where(
				squirrel.And{
					squirrel.Eq{idCol: id},
					squirrel.Eq{leaseIdCol: leaseID},
				},
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to release node")
		}
		rowsAffected, raErr := res.RowsAffected()
		if raErr == nil {
			if rowsAffected <= 0 {
				return nil, ErrBadRelease
			}
		}
		// If the DB driver doesn't give us rows affected back, just carry on
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlNodeLeasorStorage) getNodeLeaseTxFn(id string, tag *string, specificLeaseID string) func(tx *sql.Tx) (interface{}, error) {
	return func(tx *sql.Tx) (interface{}, error) {
		now := clock.Now()
		timeoutThreashold := now.Add(-leaseTimeout)

		// SELECT id, vpn_ip, tag, available, last_leased_sec
		// FROM testcontroller_nodes
		// WHERE (available OR (NOT available AND last_leased_sec < 42)) AND tag = "foo"
		// LIMIT 1
		// FOR UPDATE SKIP LOCKED
		// If the ID is specified, the WHERE clause will only look for the ID, regardless of availability or tag.
		var whereClause interface{}
		if id != "" {
			whereClause = squirrel.Eq{idCol: id}
		} else {
			whereClause = squirrel.Or{
				squirrel.Eq{availCol: true},
				squirrel.And{
					squirrel.Eq{availCol: false},
					squirrel.Lt{lastLeasedCol: timeoutThreashold.Unix()},
				},
			}
			if tag != nil {
				whereClause = squirrel.And{
					whereClause.(squirrel.Or),
					squirrel.Eq{tagCol: *tag},
				}
			}
		}

		rows, err := s.builder.Select(idCol, vpnIPCol, tagCol, availCol, lastLeasedCol).
			From(nodeTable).
			Where(whereClause).
			Limit(1).
			Suffix("FOR UPDATE SKIP LOCKED").
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to acquire next available node")
		}
		defer sqorc.CloseRowsLogOnError(rows, "LeaseNode")

		nodesById, err := scanNodes(rows)
		if err != nil {
			return nil, err
		}
		if funk.IsEmpty(nodesById) {
			return nil, nil
		}
		selectedNode := nodesById[funk.Head(funk.Keys(nodesById)).(string)]

		// Call this callback in case we want to pause here during a test
		selectedNextNode()

		// Now mark the node as leased
		leaseID := specificLeaseID
		if leaseID == "" {
			leaseID = s.idGenerator.New()
		}
		_, err = s.builder.Update(nodeTable).
			Set(availCol, false).
			Set(lastLeasedCol, now.Unix()).
			Set(leaseIdCol, leaseID).
			Where(squirrel.Eq{idCol: selectedNode.Id}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "faield to mark node as leased")
		}
		return &NodeLease{
			Id:      selectedNode.Id,
			LeaseID: leaseID,
			VpnIP:   selectedNode.VpnIp,
		}, nil
	}
}

func scanNodes(rows *sql.Rows) (map[string]*CINode, error) {
	ret := map[string]*CINode{}
	for rows.Next() {
		var id, ip, tag string
		var avail bool
		var lastLeased int64

		err := rows.Scan(&id, &ip, &tag, &avail, &lastLeased)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan node row")
		}

		lastLeasedTs, err := ptypes.TimestampProto(time.Unix(lastLeased, 0))
		if err != nil {
			return nil, errors.Wrapf(err, "could not validate last leased time %d", lastLeasedTs)
		}
		ret[id] = &CINode{
			Id:            id,
			VpnIp:         ip,
			Tag:           tag,
			Available:     avail,
			LastLeaseTime: lastLeasedTs,
		}
	}
	return ret, nil
}
