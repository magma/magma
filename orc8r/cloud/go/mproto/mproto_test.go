/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mproto_test

import (
	"testing"

	"magma/orc8r/cloud/go/mproto"
	"magma/orc8r/cloud/go/mproto/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

// TestEncodeProtosDeterministic checks if EncodeProtosDeterministic truly enfores
// deterministic encoding by comparing encoded protobuf messages containing the same data
func TestEncodeProtosDeterministic(t *testing.T) {
	// Encoding simple proto messages (no compound fields)
	testDataObjs := []protos.TestData{
		protos.TestData{Key: "12345", Value: 10},
		protos.TestData{Key: "23456", Value: 15},
		protos.TestData{Key: "34567", Value: 20},
		protos.TestData{Key: "45678", Value: 25},
		protos.TestData{Key: "56789", Value: 30},
		protos.TestData{Key: "67890", Value: 35},
	}

	dataObjCount := len(testDataObjs)
	map1, map2 := map[string]proto.Message{}, map[string]proto.Message{}
	for i := 0; i < dataObjCount; i ++ {
		ind1, ind2 := i, (i + 3) % dataObjCount
		map1[testDataObjs[ind1].Key] = &testDataObjs[ind1]
		map2[testDataObjs[ind2].Key] = &testDataObjs[ind2]
	}

	encoded1, _ := mproto.EncodeProtosDeterministic(map1)
	encoded2, _ := mproto.EncodeProtosDeterministic(map2)
	assert.Equal(t, encoded1, encoded2)

	// Encoding compound proto messages (with repeated & map fields)
	submap1, submap2, submap3 := map[string]*protos.TestData{}, map[string]*protos.TestData{}, map[string]*protos.TestData{}
	for i := 0; i < dataObjCount; i ++ {
		submap1[testDataObjs[i].Key] = &testDataObjs[i]
		if i % 2 == 0 { submap2[testDataObjs[i].Key] = &testDataObjs[i] }
		if i % 3 == 0 { submap3[testDataObjs[i].Key] = &testDataObjs[i] }
	}

	subslice1, subslice2, subslice3 := []*protos.TestData{}, []*protos.TestData{}, []*protos.TestData{}
	for i := 0; i < dataObjCount; i ++ {
		subslice3 = append(subslice3, &testDataObjs[i])
		if i % 2 == 0 { subslice2 = append(subslice2, &testDataObjs[i]) }
		if i % 3 == 0 { subslice1 = append(subslice1, &testDataObjs[i]) }

	}

	testDataCollections := []protos.TestDataCollection{
		protos.TestDataCollection{
			ID: "c1",
			SingleData: &testDataObjs[0],
			DataMap: submap1,
			DataSlice: subslice1,
		},
		protos.TestDataCollection{
			ID: "c2",
			SingleData: &testDataObjs[1],
			DataMap: submap2,
			DataSlice: subslice2,
		},
		protos.TestDataCollection{
			ID: "c3",
			SingleData: &testDataObjs[2],
			DataMap: submap3,
			DataSlice: subslice3,
		},
	}

	map3, map4 := map[string]proto.Message{}, map[string]proto.Message{}
	for i := 0; i < 3; i ++ {
		ind1, ind2 := i, 3 - i - 1
		map3[testDataCollections[ind1].ID] = &testDataCollections[ind1]
		map4[testDataCollections[ind2].ID] = &testDataCollections[ind2]
	}

	encoded3, _ := mproto.EncodeProtosDeterministic(map3)
	encoded4, _ := mproto.EncodeProtosDeterministic(map4)
	assert.Equal(t, encoded3, encoded4)
}
