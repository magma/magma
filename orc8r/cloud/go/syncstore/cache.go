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
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-multierror"

	"magma/orc8r/cloud/go/sqorc"
)

type cacheWriter struct {
	network string
	// id is the db-wide unique identifier of each cacheWriter object, and is
	// also used to name the temporary table owned by the cacheWriter.
	id string
	// cacheTableName refers to the table with all the cached objects that the
	// cacheWriter will apply batch update into.
	cacheTableName string
	db             *sql.DB
	builder        sqorc.StatementBuilder
	invalid        bool
}

func (l *syncStore) NewCacheWriter(network string, id string) CacheWriter {
	writer := &cacheWriter{
		network:        network,
		id:             id,
		db:             l.db,
		builder:        l.builder,
		cacheTableName: l.cacheTableName,
		invalid:        false,
	}
	return writer
}

// InsertMany inserts a batch of objects into the temporary table of the
// CacheWriter object.
func (l *cacheWriter) InsertMany(objects map[string][]byte) error {
	if l.invalid {
		return fmt.Errorf("attempt to insert into network %+v with invalid cache writer", l.network)
	}
	if len(objects) == 0 {
		return nil
	}

	insertQuery := l.builder.
		Insert(l.id).
		Columns(nidCol, idCol, objCol)
	errs := &multierror.Error{}
	for id, obj := range objects {
		insertQuery = insertQuery.Values(l.network, id, obj)
	}
	if errs.ErrorOrNil() != nil {
		return errs
	}

	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := insertQuery.RunWith(tx).Exec()
		if err != nil {
			return nil, fmt.Errorf("insert objs into store for network %+v: %w", l.network, err)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *cacheWriter) Apply() error {
	if l.invalid {
		return fmt.Errorf("attempt to apply updates to network %+v with invalid cache writer", l.network)
	}
	txFn := func(tx *sql.Tx) (interface{}, error) {
		// HACK: hard coding part of this sql query because there currently doesn't exist good support
		// for "WHERE (row NOT IN other_table)" with squirrel
		//
		// The SQL query should look something like
		// DELETE FROM cached_objs WHERE
		//     network_id NOT IN ${networks}
		// AND
		//	   (network_id, id) NOT IN (SELECT network_id, id FROM cached_objs_tmp)
		_, err := l.builder.
			Delete(l.cacheTableName).
			Where(squirrel.And{
				squirrel.Eq{nidCol: l.network},
				squirrel.Expr(fmt.Sprintf(
					"(%s, %s) NOT IN (SELECT %s, %s FROM %s)",
					nidCol, idCol,
					nidCol, idCol, l.id,
				)),
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("clean up previous cached objs store table: %w", err)
		}

		// The upsert query should look something like
		// INSERT INTO cached_objs
		//     SELECT network_id, id, obj FROM cached_objs_tmp
		// 	   WHERE network_id = ${network}
		// ON CONFLICT (network_id, id)
		// 	   DO UPDATE SET obj = cached_objs_tmp.obj
		conflictUpdateTarget := sqorc.FmtConflictUpdateTarget(l.id, objCol)
		_, err = l.builder.
			Insert(l.cacheTableName).
			Select(
				l.builder.
					Select(nidCol, idCol, objCol).
					From(l.id).
					Where(squirrel.Eq{nidCol: l.network}),
			).
			OnConflict(
				[]sqorc.UpsertValue{{
					Column: objCol,
					Value:  squirrel.Expr(conflictUpdateTarget),
				}},
				nidCol, idCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("populate cached objs store table: %w", err)
		}

		_, err = l.builder.
			Delete(l.id).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, fmt.Errorf("clean up tmp cached objs store table: %w", err)
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return err
	}

	l.invalid = true
	return nil
}
