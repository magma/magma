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
	"fmt"
	"strings"

	"github.com/facebookincubator/ent/dialect/sql"
)

// Blob is the model entity for the Blob schema.
type Blob struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// NetworkID holds the value of the "network_id" field.
	NetworkID string `json:"network_id,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Key holds the value of the "key" field.
	Key string `json:"key,omitempty"`
	// Value holds the value of the "value" field.
	Value []byte `json:"value,omitempty"`
	// Version holds the value of the "version" field.
	Version uint64 `json:"version,omitempty"`
}

// FromRows scans the sql response data into Blob.
func (b *Blob) FromRows(rows *sql.Rows) error {
	var scanb struct {
		ID        int
		NetworkID sql.NullString
		Type      sql.NullString
		Key       sql.NullString
		Value     []byte
		Version   sql.NullInt64
	}
	// the order here should be the same as in the `blob.Columns`.
	if err := rows.Scan(
		&scanb.ID,
		&scanb.NetworkID,
		&scanb.Type,
		&scanb.Key,
		&scanb.Value,
		&scanb.Version,
	); err != nil {
		return err
	}
	b.ID = scanb.ID
	b.NetworkID = scanb.NetworkID.String
	b.Type = scanb.Type.String
	b.Key = scanb.Key.String
	b.Value = scanb.Value
	b.Version = uint64(scanb.Version.Int64)
	return nil
}

// Update returns a builder for updating this Blob.
// Note that, you need to call Blob.Unwrap() before calling this method, if this Blob
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Blob) Update() *BlobUpdateOne {
	return (&BlobClient{b.config}).UpdateOne(b)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (b *Blob) Unwrap() *Blob {
	tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Blob is not a transactional entity")
	}
	b.config.driver = tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Blob) String() string {
	var builder strings.Builder
	builder.WriteString("Blob(")
	builder.WriteString(fmt.Sprintf("id=%v", b.ID))
	builder.WriteString(", network_id=")
	builder.WriteString(b.NetworkID)
	builder.WriteString(", type=")
	builder.WriteString(b.Type)
	builder.WriteString(", key=")
	builder.WriteString(b.Key)
	builder.WriteString(", value=")
	builder.WriteString(fmt.Sprintf("%v", b.Value))
	builder.WriteString(", version=")
	builder.WriteString(fmt.Sprintf("%v", b.Version))
	builder.WriteByte(')')
	return builder.String()
}

// Blobs is a parsable slice of Blob.
type Blobs []*Blob

// FromRows scans the sql response data into Blobs.
func (b *Blobs) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanb := &Blob{}
		if err := scanb.FromRows(rows); err != nil {
			return err
		}
		*b = append(*b, scanb)
	}
	return nil
}

func (b Blobs) config(cfg config) {
	for _i := range b {
		b[_i].config = cfg
	}
}
