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

package object_store

import (
	"github.com/golang/glog"
)

// ObjectMap is an interface for getting objects from an arbitrary data store
type ObjectMap interface {
	Set(key string, object interface{}) error
	Delete(key string) error
	Get(key string) (interface{}, error)
	GetAll() (map[string]interface{}, error)
	DeleteAll() error
}

// Serializer turns an object into a string
type Serializer func(object interface{}) (string, error)

// Deserializer turns a string into an object
type Deserializer func(serialized string) (interface{}, error)

// RedisMap is an ObjectMap that stores objects in Redis using hash functions
type RedisMap struct {
	client       RedisClient
	hash         string
	serializer   Serializer
	deserializer Deserializer
}

// NewRedisMap creates a new redis map
func NewRedisMap(
	client RedisClient,
	hash string,
	serializer Serializer,
	deserializer Deserializer,
) *RedisMap {
	return &RedisMap{
		client:       client,
		hash:         hash,
		serializer:   serializer,
		deserializer: deserializer,
	}
}

// Set sets an object in the map
func (rm *RedisMap) Set(key string, object interface{}) error {
	str, err := rm.serializer(object)
	if err != nil {
		return err
	}
	return rm.client.HSet(rm.hash, key, str)
}

// Get retrieves an object from the map
func (rm *RedisMap) Get(key string) (interface{}, error) {
	val, err := rm.client.HGet(rm.hash, key)
	if err != nil {
		return nil, err
	}
	return rm.deserializer(val)
}

func (rm *RedisMap) Delete(key string) error {
	return rm.client.HDel(rm.hash, key)
}

// GetAll returns all objects in the map
func (rm *RedisMap) GetAll() (map[string]interface{}, error) {
	valMap, err := rm.client.HGetAll(rm.hash)
	if err != nil {
		return nil, err
	}
	returnVals := make(map[string]interface{})
	for key, val := range valMap {
		obj, err := rm.deserializer(val)
		if err != nil {
			glog.Errorf("Unable to parse key %s because: %s", key, err.Error())
		} else {
			returnVals[key] = obj
		}
	}
	return returnVals, nil
}

func (rm *RedisMap) DeleteAll() error {
	valMap, err := rm.client.HGetAll(rm.hash)
	if err != nil {
		return err
	}
	for key := range valMap {
		err = rm.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
