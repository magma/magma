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

package object_store_test

import (
	"fmt"
	"testing"

	"magma/feg/gateway/object_store"

	"github.com/stretchr/testify/assert"
)

type mockRedisClient struct {
	dataMap map[string]string
}

func (client *mockRedisClient) HSet(hash string, field string, value string) error {
	client.dataMap[field] = value
	return nil
}

func (client *mockRedisClient) HGet(hash string, field string) (string, error) {
	str, ok := client.dataMap[field]
	if !ok {
		return "", fmt.Errorf("Not found: %s", field)
	}
	return str, nil
}

func (client *mockRedisClient) HGetAll(hash string) (map[string]string, error) {
	return client.dataMap, nil
}

func (client *mockRedisClient) HDel(hash string, field string) error {
	delete(client.dataMap, field)
	return nil
}

type testObject struct {
	foo string
}

func getSerializer() object_store.Serializer {
	return func(object interface{}) (string, error) {
		obj, ok := object.(*testObject)
		if !ok {
			return "", fmt.Errorf("Could not cast interface to object")
		}
		return obj.foo, nil
	}
}

func getDeserializer() object_store.Deserializer {
	return func(serialized string) (interface{}, error) {
		obj := &testObject{}
		obj.foo = serialized
		return obj, nil
	}
}

func TestRedisMap(t *testing.T) {
	redisClient := &mockRedisClient{dataMap: make(map[string]string)}
	redisMap := object_store.NewRedisMap(redisClient, "hash", getSerializer(), getDeserializer())

	var err error
	err = redisMap.Set("1", &testObject{foo: "first"})
	assert.NoError(t, err)
	err = redisMap.Set("2", &testObject{foo: "second"})
	assert.NoError(t, err)

	objRaw, err := redisMap.Get("1")
	assert.NoError(t, err)

	obj, ok := objRaw.(*testObject)
	assert.True(t, ok)
	assert.Equal(t, obj.foo, "first")

	objRaw2, err := redisMap.Get("2")
	assert.NoError(t, err)
	obj2, ok := objRaw2.(*testObject)
	assert.True(t, ok)
	assert.Equal(t, obj2.foo, "second")

	redisMap.Set("3", &testObject{foo: "last"})
	allVals, err := redisMap.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, len(allVals), 3)

	for _, valRaw := range allVals {
		_, ok := valRaw.(*testObject)
		assert.True(t, ok)
	}
}
