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

package state

import (
	"encoding/json"

	"magma/orc8r/cloud/go/serde"
)

func NewStateSerde(stateType string, modelPtr serde.ValidateableBinaryConvertible) serde.Serde {
	return serde.NewBinarySerde(SerdeDomain, stateType, modelPtr)
}

// StringToStringMap is a generic map that holds key value pair both of type
// string. This is used on the gateway side in checkin_cli.py to simply test
// the connection between the cloud and the gateway.
type StringToStringMap map[string]string

func (m *StringToStringMap) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *StringToStringMap) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *StringToStringMap) ValidateModel() error {
	return nil
}

// ArbitraryJSON is a generic map[string]interface{} destination to unmarshal
// any JSON payload into which implements ValidateableBinaryConvertible.
// This is used for replicated gateway states which are serialized as
// protos on the gateway side to avoid double-defining structs in proto and
// swagger for the time being.
type ArbitraryJSON map[string]interface{}

func (j *ArbitraryJSON) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *ArbitraryJSON) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, j)
}

func (j *ArbitraryJSON) ValidateModel() error {
	return nil
}
