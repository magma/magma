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

package models_test

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/models"

	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
)

func TestJSONMapToProtobufStruct(t *testing.T) {
	jsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(jsonMap)
	assert.NoError(t, err)
	expectedProtobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), expectedProtobufStruct)
	assert.NoError(t, err)

	actualProtobufStruct, err := models.JSONMapToProtobufStruct(jsonMap)

	assert.NoError(t, err)
	assert.Equal(t, expectedProtobufStruct, actualProtobufStruct)
}

func TestProtobufStructToJSONMap(t *testing.T) {
	expectedJsonMap := map[string]interface{}{
		"nil":    nil,
		"number": 1.0,
		"string": "str",
		"struct": map[string]interface{}{
			"a": 2.0,
		},
		"list": []interface{}{1.0, "foo"},
	}
	marshaled, err := json.Marshal(expectedJsonMap)
	assert.NoError(t, err)
	protobufStruct := &structpb.Struct{}
	err = jsonpb.UnmarshalString(string(marshaled), protobufStruct)
	assert.NoError(t, err)

	actualJsonMap, err := models.ProtobufStructToJSONMap(protobufStruct)

	assert.NoError(t, err)
	assert.Equal(t, expectedJsonMap, actualJsonMap)
}
