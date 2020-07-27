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

package blob

import (
	"magma/orc8r/cloud/go/blobstore/ent/schema"
)

const (
	// Label holds the string label denoting the blob type in the database.
	Label = "blob"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldNetworkID holds the string denoting the network_id vertex property in the database.
	FieldNetworkID = "network_id"
	// FieldType holds the string denoting the type vertex property in the database.
	FieldType = "type"
	// FieldKey holds the string denoting the key vertex property in the database.
	FieldKey = "key"
	// FieldValue holds the string denoting the value vertex property in the database.
	FieldValue = "value"
	// FieldVersion holds the string denoting the version vertex property in the database.
	FieldVersion = "version"

// Table declared below. We override the default constant definition
// of "ent", because we want to allow using the same blob API with
// different tables.

)

// Columns holds all SQL columns are blob fields.
var Columns = []string{
	FieldID,
	FieldNetworkID,
	FieldType,
	FieldKey,
	FieldValue,
	FieldVersion,
}

var (
	fields = schema.Blob{}.Fields()

	// descVersion is the schema descriptor for version field.
	descVersion = fields[4].Descriptor()
	// DefaultVersion holds the default value on creation for the version field.
	DefaultVersion = descVersion.Default.(uint64)
)

// Table holds the table name of the blob in the database.
var Table = "states"
