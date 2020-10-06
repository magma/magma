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

package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

type mockJsonObject struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestRedisStateClientTest(t *testing.T) {
	mockRedis, err := miniredis.Run()
	assert.NoError(t, err)

	client := NewDefaultRedisStateClient(mockRedis.Addr(), &JsonStateSerde{
		stateType:    "mock_state",
		dataInstance: &mockJsonObject{},
	})

	expectedObj := mockJsonObject{
		Foo: "foo1",
		Bar: 9,
	}

	// Test Happy Path
	err = client.Set("mock1", expectedObj)
	assert.NoError(t, err)

	actualObjI, err := client.Get("mock1")
	assert.NoError(t, err)
	actualObj, ok := actualObjI.(*mockJsonObject)
	assert.True(t, ok)
	assert.Equal(t, expectedObj, *actualObj)

	version, err := client.GetVersion("mock1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), version)

	expectedObj.Bar = 12
	err = client.Set("mock1", expectedObj)
	assert.NoError(t, err)

	actualObjI, err = client.Get("mock1")
	assert.NoError(t, err)
	actualObj, ok = actualObjI.(*mockJsonObject)
	assert.True(t, ok)
	assert.Equal(t, expectedObj, *actualObj)

	version, err = client.GetVersion("mock1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), version)

	err = client.MarkAsGarbage("mock1")
	assert.NoError(t, err)

	actualObjI, err = client.Get("mock1")
	assert.EqualError(t, err, "object found for key mock1 is garbage")

	deleted, err := client.Delete("mock1")
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Test error handling
	_, err = client.Get("nonexistent")
	assert.Error(t, err)

	version, err = client.GetVersion("nonexistent")
	assert.Equal(t, uint64(0), version)

	err = client.MarkAsGarbage("nonexistent")
	assert.Error(t, err)

	deleted, err = client.Delete("nonexistent")
	assert.NoError(t, err)
	assert.False(t, deleted)
}
