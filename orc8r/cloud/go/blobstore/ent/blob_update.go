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
	"magma/orc8r/cloud/go/blobstore/ent/predicate"

	"github.com/facebookincubator/ent/dialect/sql"
)

// BlobUpdate is the builder for updating Blob entities.
type BlobUpdate struct {
	config
	network_id *string
	_type      *string
	key        *string
	value      *[]byte
	clearvalue bool
	version    *uint64
	addversion *uint64
	predicates []predicate.Blob
}

// Where adds a new predicate for the builder.
func (bu *BlobUpdate) Where(ps ...predicate.Blob) *BlobUpdate {
	bu.predicates = append(bu.predicates, ps...)
	return bu
}

// SetNetworkID sets the network_id field.
func (bu *BlobUpdate) SetNetworkID(s string) *BlobUpdate {
	bu.network_id = &s
	return bu
}

// SetType sets the type field.
func (bu *BlobUpdate) SetType(s string) *BlobUpdate {
	bu._type = &s
	return bu
}

// SetKey sets the key field.
func (bu *BlobUpdate) SetKey(s string) *BlobUpdate {
	bu.key = &s
	return bu
}

// SetValue sets the value field.
func (bu *BlobUpdate) SetValue(b []byte) *BlobUpdate {
	bu.value = &b
	return bu
}

// ClearValue clears the value of value.
func (bu *BlobUpdate) ClearValue() *BlobUpdate {
	bu.value = nil
	bu.clearvalue = true
	return bu
}

// SetVersion sets the version field.
func (bu *BlobUpdate) SetVersion(u uint64) *BlobUpdate {
	bu.version = &u
	bu.addversion = nil
	return bu
}

// SetNillableVersion sets the version field if the given value is not nil.
func (bu *BlobUpdate) SetNillableVersion(u *uint64) *BlobUpdate {
	if u != nil {
		bu.SetVersion(*u)
	}
	return bu
}

// AddVersion adds u to version.
func (bu *BlobUpdate) AddVersion(u uint64) *BlobUpdate {
	if bu.addversion == nil {
		bu.addversion = &u
	} else {
		*bu.addversion += u
	}
	return bu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (bu *BlobUpdate) Save(ctx context.Context) (int, error) {
	return bu.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BlobUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BlobUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BlobUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (bu *BlobUpdate) sqlSave(ctx context.Context) (int, error) {
	var (
		builder  = sql.Dialect(bu.driver.Dialect())
		selector = builder.Select().From(builder.Table(blob.Table))
	)
	for _, p := range bu.predicates {
		p(selector)
	}
	var (
		res     sql.Result
		updater = builder.Update(blob.Table).Where(selector.P())
	)
	if value := bu.network_id; value != nil {
		updater.Set(blob.FieldNetworkID, *value)
	}
	if value := bu._type; value != nil {
		updater.Set(blob.FieldType, *value)
	}
	if value := bu.key; value != nil {
		updater.Set(blob.FieldKey, *value)
	}
	if value := bu.value; value != nil {
		updater.Set(blob.FieldValue, *value)
	}
	if bu.clearvalue {
		updater.SetNull(blob.FieldValue)
	}
	if value := bu.version; value != nil {
		updater.Set(blob.FieldVersion, *value)
	}
	if value := bu.addversion; value != nil {
		updater.Add(blob.FieldVersion, *value)
	}
	if updater.Empty() {
		return 0, nil
	}
	tx, err := bu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	query, args := updater.Query()
	if err := tx.Exec(ctx, query, args, &res); err != nil {
		return 0, rollback(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	// returning the amount of nodes is not necessary.
	return 0, nil

}

// BlobUpdateOne is the builder for updating a single Blob entity.
type BlobUpdateOne struct {
	config
	id         int
	network_id *string
	_type      *string
	key        *string
	value      *[]byte
	clearvalue bool
	version    *uint64
	addversion *uint64
}

// SetNetworkID sets the network_id field.
func (buo *BlobUpdateOne) SetNetworkID(s string) *BlobUpdateOne {
	buo.network_id = &s
	return buo
}

// SetType sets the type field.
func (buo *BlobUpdateOne) SetType(s string) *BlobUpdateOne {
	buo._type = &s
	return buo
}

// SetKey sets the key field.
func (buo *BlobUpdateOne) SetKey(s string) *BlobUpdateOne {
	buo.key = &s
	return buo
}

// SetValue sets the value field.
func (buo *BlobUpdateOne) SetValue(b []byte) *BlobUpdateOne {
	buo.value = &b
	return buo
}

// ClearValue clears the value of value.
func (buo *BlobUpdateOne) ClearValue() *BlobUpdateOne {
	buo.value = nil
	buo.clearvalue = true
	return buo
}

// SetVersion sets the version field.
func (buo *BlobUpdateOne) SetVersion(u uint64) *BlobUpdateOne {
	buo.version = &u
	buo.addversion = nil
	return buo
}

// SetNillableVersion sets the version field if the given value is not nil.
func (buo *BlobUpdateOne) SetNillableVersion(u *uint64) *BlobUpdateOne {
	if u != nil {
		buo.SetVersion(*u)
	}
	return buo
}

// AddVersion adds u to version.
func (buo *BlobUpdateOne) AddVersion(u uint64) *BlobUpdateOne {
	if buo.addversion == nil {
		buo.addversion = &u
	} else {
		*buo.addversion += u
	}
	return buo
}

// Save executes the query and returns the updated entity.
func (buo *BlobUpdateOne) Save(ctx context.Context) (*Blob, error) {
	return buo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BlobUpdateOne) SaveX(ctx context.Context) *Blob {
	b, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return b
}

// Exec executes the query on the entity.
func (buo *BlobUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BlobUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (*BlobUpdateOne) sqlSave(context.Context) (*Blob, error) {
	return nil, fmt.Errorf("cannot perform update-one on models with complex PK")
}
