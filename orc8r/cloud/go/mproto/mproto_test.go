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
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/mproto"
	"magma/orc8r/cloud/go/mproto/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

var (
	testDataObjs = []protos.TestData{
		{Key: "12345", Value: 10},
		{Key: "23456", Value: 15},
		{Key: "34567", Value: 20},
		{Key: "45678", Value: 25},
		{Key: "56789", Value: 30},
		{Key: "67890", Value: 35},
	}
	dataObjCount = len(testDataObjs)

	testdataDir    = "testdata"
	goldenFilepath = filepath.Join(testdataDir, "determinisic-digest.golden")
)

// TestEncodeProtosDeterministic checks if EncodeProtosDeterministic truly enforces
// deterministic encoding by comparing encoded protobuf messages containing the same data
func TestEncodeProtosDeterministic(t *testing.T) {
	// Encode simple proto messages (no compound fields)
	protos1, protos2 := map[string]proto.Message{}, map[string]proto.Message{}
	for i := 0; i < dataObjCount; i++ {
		ind1, ind2 := i, (i+3)%dataObjCount
		protos1[testDataObjs[ind1].Key] = &testDataObjs[ind1]
		protos2[testDataObjs[ind2].Key] = &testDataObjs[ind2]
	}

	encoded1, err1 := mproto.EncodeProtosDeterministic(protos1)
	assert.NoError(t, err1)
	encoded2, err2 := mproto.EncodeProtosDeterministic(protos2)
	assert.NoError(t, err2)
	assert.Equal(t, encoded1, encoded2)

	// Encode compound proto messages (with repeated & map fields)
	testDataCollections := prepareTestDataCollections()
	dataCollectionCount := len(testDataCollections)
	protos3, protos4 := map[string]proto.Message{}, map[string]proto.Message{}
	for i := 0; i < dataCollectionCount; i++ {
		ind1, ind2 := i, dataCollectionCount-i-1
		protos3[testDataCollections[ind1].ID] = &testDataCollections[ind1]
		protos4[testDataCollections[ind2].ID] = &testDataCollections[ind2]
	}

	encoded3, err3 := mproto.EncodeProtosDeterministic(protos3)
	assert.NoError(t, err3)
	encoded4, err4 := mproto.EncodeProtosDeterministic(protos4)
	assert.NoError(t, err4)
	assert.Equal(t, encoded3, encoded4)
}

// TestEncodeProtosDeterministicGoldenFile checks if EncodeProtosDeterministic enforces
// deterministic encoding consistently over time by conducting a golden file test
func TestEncodeProtosDeterministicGoldenFile(t *testing.T) {
	testDataCollections := prepareTestDataCollections()
	dataCollectionCount := len(testDataCollections)
	protos := map[string]proto.Message{}
	for i := 0; i < dataCollectionCount; i++ {
		protos[testDataCollections[i].ID] = &testDataCollections[i]
	}
	encoded, err1 := mproto.EncodeProtosDeterministic(protos)
	assert.NoError(t, err1)

	// Compare resultant digest to encoded content stored in the golden file
	absGoldenFilepath, err2 := filepath.Abs(goldenFilepath)
	assert.NoError(t, err2)
	goldenFileContent, err3 := ioutil.ReadFile(absGoldenFilepath)
	assert.NoError(t, err3)

	assert.Equal(t, 0, bytes.Compare(goldenFileContent, encoded))
}

func prepareTestDataCollections() []protos.TestDataCollection {
	submap1, submap2, submap3 := map[string]*protos.TestData{}, map[string]*protos.TestData{}, map[string]*protos.TestData{}
	for i := 0; i < dataObjCount; i++ {
		submap1[testDataObjs[i].Key] = &testDataObjs[i]
		if i%2 == 0 {
			submap2[testDataObjs[i].Key] = &testDataObjs[i]
		}
		if i%3 == 0 {
			submap3[testDataObjs[i].Key] = &testDataObjs[i]
		}
	}

	subslice1, subslice2, subslice3 := []*protos.TestData{}, []*protos.TestData{}, []*protos.TestData{}
	for i := 0; i < dataObjCount; i++ {
		subslice3 = append(subslice3, &testDataObjs[i])
		if i%2 == 0 {
			subslice2 = append(subslice2, &testDataObjs[i])
		}
		if i%3 == 0 {
			subslice1 = append(subslice1, &testDataObjs[i])
		}
	}

	return []protos.TestDataCollection{
		{
			ID:         "c1",
			SingleData: &testDataObjs[0],
			DataMap:    submap1,
			DataSlice:  subslice1,
		},
		{
			ID:         "c2",
			SingleData: &testDataObjs[1],
			DataMap:    submap2,
			DataSlice:  subslice2,
		},
		{
			ID:         "c3",
			SingleData: &testDataObjs[2],
			DataMap:    submap3,
			DataSlice:  subslice3,
		},
	}
}
