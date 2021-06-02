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
	b64 "encoding/base64"
	"io/ioutil"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/mproto"
	mocks "magma/orc8r/cloud/go/mproto/mocks"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

const (
	iterationCount = 1000
)

var (
	testdataDir    = "testdata"
	goldenFilepath = filepath.Join(testdataDir, "determinisic-digest.b64.golden")
)

// TestMarshalManyDeterministic checks if MarshalManyDeterministic truly enforces
// deterministic encoding by comparing encoded protobuf messages containing the same data
func TestMarshalManyDeterministic(t *testing.T) {
	// Encode basic proto messages (no compound fields)
	testDataObjs := getTestDataBasic()
	protos1 := map[string]proto.Message{
		testDataObjs[0].Key: &testDataObjs[0],
		testDataObjs[1].Key: &testDataObjs[1],
		testDataObjs[2].Key: &testDataObjs[2],
	}
	protos2 := map[string]proto.Message{
		testDataObjs[0].Key: &testDataObjs[0],
		testDataObjs[1].Key: &testDataObjs[1],
		testDataObjs[2].Key: &testDataObjs[2],
	}

	for i := 0; i < iterationCount; i++ {
		encoded1, err1 := mproto.MarshalManyDeterministic(protos1)
		assert.NoError(t, err1)
		encoded2, err2 := mproto.MarshalManyDeterministic(protos2)
		assert.NoError(t, err2)
		assert.Equal(t, encoded1, encoded2)
	}

	// Encode compound proto messages (with repeated & map fields)
	testDataCollections := getTestDataCompound()
	protos3 := map[string]proto.Message{
		testDataCollections[0].Id: &testDataCollections[0],
		testDataCollections[1].Id: &testDataCollections[1],
		testDataCollections[2].Id: &testDataCollections[2],
	}
	protos4 := map[string]proto.Message{
		testDataCollections[0].Id: &testDataCollections[0],
		testDataCollections[1].Id: &testDataCollections[1],
		testDataCollections[2].Id: &testDataCollections[2],
	}

	for i := 0; i < iterationCount; i++ {
		encoded3, err3 := mproto.MarshalManyDeterministic(protos3)
		assert.NoError(t, err3)
		encoded4, err4 := mproto.MarshalManyDeterministic(protos4)
		assert.NoError(t, err4)
		assert.Equal(t, encoded3, encoded4)
	}
}

// TestMarshalManyDeterministicGoldenFile checks if MarshalManyDeterministic enforces
// deterministic encoding consistently over time by conducting a golden file test
func TestMarshalManyDeterministicGoldenFile(t *testing.T) {
	absGoldenFilepath, err1 := filepath.Abs(goldenFilepath)
	assert.NoError(t, err1)
	goldenFileContent, err2 := ioutil.ReadFile(absGoldenFilepath)
	assert.NoError(t, err2)

	testDataCollections := getTestDataCompound()
	protos := map[string]proto.Message{
		testDataCollections[0].Id: &testDataCollections[0],
		testDataCollections[1].Id: &testDataCollections[1],
		testDataCollections[2].Id: &testDataCollections[2],
	}

	// Compare encoded digest string to content stored in the golden file
	encoded, err3 := mproto.MarshalManyDeterministic(protos)
	assert.NoError(t, err3)
	encodedB64 := b64.StdEncoding.EncodeToString(encoded)
	assert.Equal(t, string(goldenFileContent), encodedB64)
}

// getTestDataCompound generates proto messages with compound data fields
// (repeated fields and maps) to cover cases of protobuf nondeterminism
//
// https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510
func getTestDataCompound() []mocks.TestDataCompound {
	testDataObjs := getTestDataBasic()

	return []mocks.TestDataCompound{
		{
			Id:         "c1",
			SingleData: &testDataObjs[0],
			DataMap: map[string]*mocks.TestDataBasic{
				testDataObjs[0].Key: &testDataObjs[0],
				testDataObjs[1].Key: &testDataObjs[1],
				testDataObjs[2].Key: &testDataObjs[2],
				testDataObjs[3].Key: &testDataObjs[3],
			},
			DataSlice: []*mocks.TestDataBasic{&testDataObjs[5], &testDataObjs[4], &testDataObjs[3], &testDataObjs[2]},
		},
		{
			Id:         "c2",
			SingleData: &testDataObjs[1],
			DataMap: map[string]*mocks.TestDataBasic{
				testDataObjs[2].Key: &testDataObjs[2],
				testDataObjs[3].Key: &testDataObjs[3],
				testDataObjs[4].Key: &testDataObjs[4],
			},
			DataSlice: []*mocks.TestDataBasic{&testDataObjs[3], &testDataObjs[2], &testDataObjs[1]},
		},
		{
			Id:         "c3",
			SingleData: &testDataObjs[2],
			DataMap: map[string]*mocks.TestDataBasic{
				testDataObjs[4].Key: &testDataObjs[4],
				testDataObjs[5].Key: &testDataObjs[5],
			},
			DataSlice: []*mocks.TestDataBasic{&testDataObjs[1], &testDataObjs[0]},
		},
	}
}

// getTestDataCompound generates proto messages with basic data fields
// (including simple integers) to cover cases of protobuf nondeterminism
//
// https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510
func getTestDataBasic() []mocks.TestDataBasic {
	return []mocks.TestDataBasic{
		{Key: "12345", Value: 10},
		{Key: "23456", Value: 15},
		{Key: "34567", Value: 20},
		{Key: "45678", Value: 25},
		{Key: "56789", Value: 30},
		{Key: "67890", Value: 35},
	}
}
