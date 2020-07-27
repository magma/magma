/*
 * Copyright 2020 The Magma Authors
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"magma/orc8r/cloud/go/blobstore/ent/blob"
	"magma/orc8r/cloud/go/blobstore/ent/predicate"

	"github.com/facebookincubator/ent/dialect/sql"
)

// BlobDelete is the builder for deleting a Blob entity.
type BlobDelete struct {
	config
	predicates []predicate.Blob
}

// Where adds a new predicate to the delete builder.
func (bd *BlobDelete) Where(ps ...predicate.Blob) *BlobDelete {
	bd.predicates = append(bd.predicates, ps...)
	return bd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (bd *BlobDelete) Exec(ctx context.Context) (int, error) {
	return bd.sqlExec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (bd *BlobDelete) ExecX(ctx context.Context) int {
	n, err := bd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (bd *BlobDelete) sqlExec(ctx context.Context) (int, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(bd.driver.Dialect())
	)
	selector := builder.Select().From(sql.Table(blob.Table))
	for _, p := range bd.predicates {
		p(selector)
	}
	query, args := builder.Delete(blob.Table).FromSelect(selector).Query()
	if err := bd.driver.Exec(ctx, query, args, &res); err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affected), nil
}

// BlobDeleteOne is the builder for deleting a single Blob entity.
type BlobDeleteOne struct {
	bd *BlobDelete
}

// Exec executes the deletion query.
func (bdo *BlobDeleteOne) Exec(ctx context.Context) error {
	n, err := bdo.bd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &ErrNotFound{blob.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (bdo *BlobDeleteOne) ExecX(ctx context.Context) {
	bdo.bd.ExecX(ctx)
}
