/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dbtest

import (
	"database/sql"
	"testing"

	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/sqorc"
)

type ResourceManager interface {
	InTransaction(func()) error
	GetBuilder() sq.StatementBuilderType
	CreateTables(...db.Model) error
	InsertResources(db.FieldMask, ...db.Model) error
	DropResources(...db.Model) error
}

type resourceManager struct {
	tx      sq.BaseRunner
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func NewResourceManager(t *testing.T, db *sql.DB, builder sqorc.StatementBuilder) *resourceManager {
	if t == nil {
		panic("for tests only")
	}
	return &resourceManager{
		db:      db,
		builder: builder,
	}
}

func (r *resourceManager) InTransaction(fn func()) error {
	return r.inTransactionWithError(func() error {
		fn()
		return nil
	})
}

func (r *resourceManager) GetBuilder() sq.StatementBuilderType {
	return r.builder.RunWith(r.tx)
}

func (r *resourceManager) CreateTables(models ...db.Model) error {
	return r.inTransactionWithError(func() error {
		for _, model := range models {
			if err := db.CreateTable(r.tx, r.builder, model.GetMetadata()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *resourceManager) DropResources(models ...db.Model) error {
	return r.inTransactionWithError(func() error {
		for _, model := range models {
			err := db.NewQuery().
				WithBuilder(r.GetBuilder()).
				From(model).
				Where(sq.Eq{}).
				Delete()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *resourceManager) InsertResources(mask db.FieldMask, models ...db.Model) error {
	return r.inTransactionWithError(func() error {
		for _, model := range models {
			_, err := db.NewQuery().
				WithBuilder(r.GetBuilder()).
				From(model).
				Insert(mask)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *resourceManager) inTransactionWithError(fn func() error) error {
	_, err := sqorc.ExecInTx(r.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		r.tx = tx
		return nil, fn()
	})
	return err
}
