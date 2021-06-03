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

package test

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"testing"

	"magma/orc8r/cloud/go/mproto"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

const (
	iterationCount = 1000
)

var (
	goldenFilepath = "../testdata/determinisic-digest.b64.golden"
)

// TestMarshalManyDeterministic checks if MarshalManyDeterministic truly enforces
// deterministic encoding by comparing encoded protobuf messages containing the same data.
//
// Ref: https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510
func TestMarshalManyDeterministic(t *testing.T) {
	// Encode basic proto messages (no compound fields).
	protos1 := map[string]proto.Message{
		"12345": &TestDataBasic{Key: "12345", Value: 10},
		"23456": &TestDataBasic{Key: "23456", Value: 15},
		"34567": &TestDataBasic{Key: "34567", Value: 20},
	}
	encodedCanon1, err1 := mproto.MarshalManyDeterministic(protos1)
	assert.NoError(t, err1)

	for i := 0; i < iterationCount; i++ {
		encoded, err := mproto.MarshalManyDeterministic(protos1)
		assert.NoError(t, err)
		assert.Equal(t, encoded, encodedCanon1)
	}

	// Encode compound proto messages (with repeated & map fields).
	testDataCollections := getTestDataCompound()
	protos2 := map[string]proto.Message{
		testDataCollections[0].Id: testDataCollections[0],
		testDataCollections[1].Id: testDataCollections[1],
		testDataCollections[2].Id: testDataCollections[2],
	}
	encodedCanon2, err2 := mproto.MarshalManyDeterministic(protos2)
	assert.NoError(t, err2)

	for i := 0; i < iterationCount; i++ {
		encoded, err := mproto.MarshalManyDeterministic(protos2)
		assert.NoError(t, err)
		assert.Equal(t, encoded, encodedCanon2)
	}
}

// TestMarshalManyDeterministicGoldenFile checks if MarshalManyDeterministic enforces
// deterministic encoding consistently over time by conducting a golden file test.
func TestMarshalManyDeterministicGoldenFile(t *testing.T) {
	absGoldenFilepath, err1 := filepath.Abs(goldenFilepath)
	assert.NoError(t, err1)
	goldenFileContent, err2 := ioutil.ReadFile(absGoldenFilepath)
	assert.NoError(t, err2)

	testDataCollections := getTestDataCompound()
	protos := map[string]proto.Message{
		testDataCollections[0].Id: testDataCollections[0],
		testDataCollections[1].Id: testDataCollections[1],
		testDataCollections[2].Id: testDataCollections[2],
	}

	// Compare encoded digest string to content stored in the golden file.
	encoded, err3 := mproto.MarshalManyDeterministic(protos)
	assert.NoError(t, err3)
	encodedB64 := base64.StdEncoding.EncodeToString(encoded)
	assert.Equal(t, string(goldenFileContent), encodedB64)
}

// getTestDataCompound generates proto messages with compound data fields
// (repeated fields and maps) to cover cases of protobuf nondeterminism.
func getTestDataCompound() []*TestDataCompound {
	return []*TestDataCompound{
		{
			Id:         "c1",
			SingleData: &TestDataBasic{Key: "12345", Value: 10},
			DataMap: map[string]*TestDataBasic{
				"12345": {Key: "12345", Value: 10},
				"23456": {Key: "23456", Value: 15},
				"34567": {Key: "34567", Value: 20},
				"45678": {Key: "45678", Value: 25},
			},
			DataSlice: []*TestDataBasic{
				{Key: "67890", Value: 35},
				{Key: "56789", Value: 30},
				{Key: "45678", Value: 25},
				{Key: "34567", Value: 20},
			},
		},
		{
			Id:         "c2",
			SingleData: &TestDataBasic{Key: "23456", Value: 15},
			DataMap: map[string]*TestDataBasic{
				"34567": {Key: "34567", Value: 20},
				"45678": {Key: "45678", Value: 25},
				"56789": {Key: "56789", Value: 30},
			},
			DataSlice: []*TestDataBasic{
				{Key: "45678", Value: 25},
				{Key: "34567", Value: 20},
				{Key: "23456", Value: 15},
			},
		},
		{
			Id:         "c3",
			SingleData: &TestDataBasic{Key: "34567", Value: 20},
			DataMap: map[string]*TestDataBasic{
				"56789": {Key: "56789", Value: 30},
				"67890": {Key: "67890", Value: 35},
			},
			DataSlice: []*TestDataBasic{
				{Key: "23456", Value: 15},
				{Key: "12345", Value: 10},
			},
		},
	}
}
