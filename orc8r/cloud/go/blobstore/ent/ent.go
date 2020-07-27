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
	"fmt"
	"magma/orc8r/cloud/go/blobstore/ent/blob"
	"strings"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
)

// Order applies an ordering on either graph traversal or sql selector.
type Order func(*sql.Selector)

// Asc applies the given fields in ASC order.
func Asc(fields ...string) Order {
	return Order(
		func(s *sql.Selector) {
			for _, f := range fields {
				s.OrderBy(sql.Asc(f))
			}
		},
	)
}

// Desc applies the given fields in DESC order.
func Desc(fields ...string) Order {
	return Order(
		func(s *sql.Selector) {
			for _, f := range fields {
				s.OrderBy(sql.Desc(f))
			}
		},
	)
}

// Aggregate applies an aggregation step on the group-by traversal/selector.
type Aggregate struct {
	// SQL the column wrapped with the aggregation function.
	SQL func(*sql.Selector) string
}

// As is a pseudo aggregation function for renaming another other functions with custom names. For example:
//
//	GroupBy(field1, field2).
//	Aggregate(ent.As(ent.Sum(field1), "sum_field1"), (ent.As(ent.Sum(field2), "sum_field2")).
//	Scan(ctx, &v)
//
func As(fn Aggregate, end string) Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.As(fn.SQL(s), end)
		},
	}
}

// Count applies the "count" aggregation function on each group.
func Count() Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.Count("*")
		},
	}
}

// Max applies the "max" aggregation function on the given field of each group.
func Max(field string) Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.Max(s.C(field))
		},
	}
}

// Mean applies the "mean" aggregation function on the given field of each group.
func Mean(field string) Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.Avg(s.C(field))
		},
	}
}

// Min applies the "min" aggregation function on the given field of each group.
func Min(field string) Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.Min(s.C(field))
		},
	}
}

// Sum applies the "sum" aggregation function on the given field of each group.
func Sum(field string) Aggregate {
	return Aggregate{
		SQL: func(s *sql.Selector) string {
			return sql.Sum(s.C(field))
		},
	}
}

// ErrNotFound returns when trying to fetch a specific entity and it was not found in the database.
type ErrNotFound struct {
	label string
}

// Error implements the error interface.
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("ent: %s not found", e.label)
}

// IsNotFound returns a boolean indicating whether the error is a not found error.
func IsNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

// MaskNotFound masks nor found error.
func MaskNotFound(err error) error {
	if IsNotFound(err) {
		return nil
	}
	return err
}

// ErrNotSingular returns when trying to fetch a singular entity and more then one was found in the database.
type ErrNotSingular struct {
	label string
}

// Error implements the error interface.
func (e *ErrNotSingular) Error() string {
	return fmt.Sprintf("ent: %s not singular", e.label)
}

// IsNotSingular returns a boolean indicating whether the error is a not singular error.
func IsNotSingular(err error) bool {
	_, ok := err.(*ErrNotSingular)
	return ok
}

// ErrConstraintFailed returns when trying to create/update one or more entities and
// one or more of their constraints failed. For example, violation of edge or field uniqueness.
type ErrConstraintFailed struct {
	msg  string
	wrap error
}

// Error implements the error interface.
func (e ErrConstraintFailed) Error() string {
	return fmt.Sprintf("ent: unique constraint failed: %s", e.msg)
}

// Unwrap implements the errors.Wrapper interface.
func (e *ErrConstraintFailed) Unwrap() error {
	return e.wrap
}

// IsConstraintFailure returns a boolean indicating whether the error is a constraint failure.
func IsConstraintFailure(err error) bool {
	_, ok := err.(*ErrConstraintFailed)
	return ok
}

func isSQLConstraintError(err error) (*ErrConstraintFailed, bool) {
	var (
		msg = err.Error()
		// error format per dialect.
		errors = [...]string{
			"Error 1062",               // MySQL 1062 error (ER_DUP_ENTRY).
			"UNIQUE constraint failed", // SQLite.
			"duplicate key value violates unique constraint", // PostgreSQL.
		}
	)
	for i := range errors {
		if strings.Contains(msg, errors[i]) {
			return &ErrConstraintFailed{msg, err}, true
		}
	}
	return nil, false
}

// rollback calls to tx.Rollback and wraps the given error with the rollback error if occurred.
func rollback(tx dialect.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%s: %v", err.Error(), rerr)
	}
	if err, ok := isSQLConstraintError(err); ok {
		return err
	}
	return err
}

// insertLastID invokes the insert query on the transaction and returns the LastInsertID.
func insertLastID(ctx context.Context, tx dialect.Tx, insert *sql.InsertBuilder) (int64, error) {
	var (
		res         sql.Result
		query, args = insert.Returning().Query()
	)
	// return zero value for id, because it doesn't
	// exist on the database.
	return 0, tx.Exec(ctx, query, args, &res)
}

func init() {
	// remove the id field from blob columns.
	blob.Columns = blob.Columns[1:]
}

// BeginTx returns a transactional client with options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(*sql.Driver).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: &txDriver{tx: tx, drv: c.driver}, log: c.log, debug: c.debug}
	return &Tx{
		config: cfg,
		Blob:   NewBlobClient(cfg),
	}, nil
}

// keys returns the keys/ids from the edge map.
func keys(m map[int]struct{}) []int {
	s := make([]int, 0, len(m))
	for id := range m {
		s = append(s, id)
	}
	return s
}
