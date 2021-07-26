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

	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

type cacheWriter struct {
	network       string
	db            *sql.DB
	builder       sqorc.StatementBuilder
	invalidWriter bool
}

func NewCacheWriter(network string, db *sql.DB, builder sqorc.StatementBuilder) CacheWriter {
	return &cacheWriter{network: network, db: db, builder: builder, invalidWriter: false}
}

func (l *cacheWriter) InsertMany(objects map[string][]byte) error {
	if l.invalidWriter {
		return errors.Errorf("attempt to insert into network %+v with invalid cache writer", l.network)
	}
	if len(objects) == 0 {
		return nil
	}

	insertQuery := l.builder.
		Insert(cacheTmpTableName).
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
		return nil, errors.Wrapf(err, "insert objs into store for network %+v", l.network)
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	return err
}

func (l *cacheWriter) Apply() error {
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
			Delete(cacheTableName).
			Where(squirrel.And{
				squirrel.Eq{nidCol: l.network},
				squirrel.Expr(fmt.Sprintf(
					"(%s, %s) NOT IN (SELECT %s, %s FROM %s)",
					nidCol, idCol,
					nidCol, idCol, cacheTmpTableName,
				)),
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "clean up previous cached objs store table")
		}

		// The upsert query should look something like
		// INSERT INTO cached_objs
		//     SELECT network_id, id, obj FROM cached_objs_tmp
		// 	   WHERE network_id = ${network}
		// ON CONFLICT (network_id, id)
		// 	   DO UPDATE SET obj = cached_objs_tmp.obj
		conflictUpdateTarget := sqorc.FmtConflictUpdateTarget(cacheTmpTableName, objCol)
		_, err = l.builder.
			Insert(cacheTableName).
			Select(
				l.builder.
					Select(nidCol, idCol, objCol).
					From(cacheTmpTableName).
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
			return nil, errors.Wrap(err, "populate cached objs store table")
		}

		_, err = l.builder.
			Delete(cacheTmpTableName).
			Where(squirrel.Eq{nidCol: l.network}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "clean up tmp cached objs store table")
		}
		return nil, nil
	}
	_, err := sqorc.ExecInTx(l.db, nil, nil, txFn)
	if err != nil {
		return err
	}

	l.invalidWriter = true
	return nil
}
