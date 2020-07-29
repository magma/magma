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
	"errors"
	"magma/orc8r/cloud/go/blobstore/ent/blob"

	"github.com/facebookincubator/ent/dialect/sql"
)

// BlobCreate is the builder for creating a Blob entity.
type BlobCreate struct {
	config
	network_id *string
	_type      *string
	key        *string
	value      *[]byte
	version    *uint64
}

// SetNetworkID sets the network_id field.
func (bc *BlobCreate) SetNetworkID(s string) *BlobCreate {
	bc.network_id = &s
	return bc
}

// SetType sets the type field.
func (bc *BlobCreate) SetType(s string) *BlobCreate {
	bc._type = &s
	return bc
}

// SetKey sets the key field.
func (bc *BlobCreate) SetKey(s string) *BlobCreate {
	bc.key = &s
	return bc
}

// SetValue sets the value field.
func (bc *BlobCreate) SetValue(b []byte) *BlobCreate {
	bc.value = &b
	return bc
}

// SetVersion sets the version field.
func (bc *BlobCreate) SetVersion(u uint64) *BlobCreate {
	bc.version = &u
	return bc
}

// SetNillableVersion sets the version field if the given value is not nil.
func (bc *BlobCreate) SetNillableVersion(u *uint64) *BlobCreate {
	if u != nil {
		bc.SetVersion(*u)
	}
	return bc
}

// Save creates the Blob in the database.
func (bc *BlobCreate) Save(ctx context.Context) (*Blob, error) {
	if bc.network_id == nil {
		return nil, errors.New("ent: missing required field \"network_id\"")
	}
	if bc._type == nil {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if bc.key == nil {
		return nil, errors.New("ent: missing required field \"key\"")
	}
	if bc.version == nil {
		v := blob.DefaultVersion
		bc.version = &v
	}
	return bc.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (bc *BlobCreate) SaveX(ctx context.Context) *Blob {
	v, err := bc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (bc *BlobCreate) sqlSave(ctx context.Context) (*Blob, error) {
	var (
		builder = sql.Dialect(bc.driver.Dialect())
		b       = &Blob{config: bc.config}
	)
	tx, err := bc.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(blob.Table).Default()
	if value := bc.network_id; value != nil {
		insert.Set(blob.FieldNetworkID, *value)
		b.NetworkID = *value
	}
	if value := bc._type; value != nil {
		insert.Set(blob.FieldType, *value)
		b.Type = *value
	}
	if value := bc.key; value != nil {
		insert.Set(blob.FieldKey, *value)
		b.Key = *value
	}
	if value := bc.value; value != nil {
		insert.Set(blob.FieldValue, *value)
		b.Value = *value
	}
	if value := bc.version; value != nil {
		insert.Set(blob.FieldVersion, *value)
		b.Version = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(blob.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	b.ID = int(id)
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return b, nil
}
