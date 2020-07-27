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

package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Blob holds the schema definition for the Blob entity.
type Blob struct {
	ent.Schema
}

func (Blob) Config() ent.Config {
	return ent.Config{
		Table: "states",
	}
}

// Fields of the Blob.
func (Blob) Fields() []ent.Field {
	return []ent.Field{
		field.Text("network_id"),
		field.Text("type"),
		field.Text("key"),
		field.Bytes("value").
			Optional(),
		field.Uint64("version").
			Default(0),
	}
}
