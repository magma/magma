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
	"encoding/json"
	"fmt"
	"reflect"

	"magma/orc8r/lib/go/protos"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

type RedisStateClient struct {
	redisClient *redis.Client
	stateSerde  RedisStateSerde
}

// RedisStateSerde defines an interfaces for state serialization
type RedisStateSerde interface {
	// Serialize defines a method to serialize the state and state's version
	// into a byte slice. Version is expected to be incremented each time the
	// state is updated. This field is used to keep the gateway state in sync
	// with the orc8r.
	Serialize(in interface{}, version uint64) ([]byte, error)

	// Deserialize defines a method to deserialize the stored state into an
	// interface of the serde's internal data instance.
	Deserialize(in []byte) (interface{}, error)

	// GetStateType defines a method to fetch the serde's state type.
	GetStateType() string
}

type JsonStateSerde struct {
	stateType    string
	dataInstance interface{}
}

// NewJsonStateSerde creates a new serde for the given state type
func NewJsonStateSerde(stateType string, dataInstance interface{}) *JsonStateSerde {
	return &JsonStateSerde{
		stateType:    stateType,
		dataInstance: dataInstance,
	}
}

// NewDefaultRedisStateClient initializes a redis client with a state serde
func NewDefaultRedisStateClient(redisAddr string, serde RedisStateSerde) *RedisStateClient {
	return &RedisStateClient{
		redisClient: redis.NewClient(
			&redis.Options{
				Addr: redisAddr,
			},
		),
		stateSerde: serde,
	}
}

// Set creates or updates the state for the given key
func (r *RedisStateClient) Set(key string, value interface{}) error {
	version, err := r.GetVersion(key)
	if err != nil {
		return err
	}
	serializedValue, err := r.stateSerde.Serialize(value, version+1)
	if err != nil {
		return err
	}
	compositeKey := r.makeCompositeKey(key, r.stateSerde.GetStateType())
	return r.redisClient.Set(compositeKey, serializedValue, 0).Err()
}

// Get fetches the object for the provided key, using the registered serde to
// deserialize the stored data
func (r *RedisStateClient) Get(key string) (interface{}, error) {
	compositeKey := r.makeCompositeKey(key, r.stateSerde.GetStateType())
	serializedState, err := r.redisClient.Get(compositeKey).Result()
	if err != nil {
		return nil, err
	}
	redisState := &protos.RedisState{}
	err = proto.Unmarshal([]byte(serializedState), redisState)
	if err != nil {
		return nil, err
	}
	if redisState.GetIsGarbage() {
		return nil, fmt.Errorf("object found for key %s is garbage", key)
	}
	return r.stateSerde.Deserialize([]byte(serializedState))
}

// Delete removes the key from redis. The function MarkAsGarbage should be
// preferred as the Delete function won't allow for any garbage collection.
func (r *RedisStateClient) Delete(key string) (bool, error) {
	compositeKey := r.makeCompositeKey(key, r.stateSerde.GetStateType())
	deleted, err := r.redisClient.Del(compositeKey).Result()
	return deleted != 0, err
}

// GetVersion returns the state's version for the provided key
func (r *RedisStateClient) GetVersion(key string) (uint64, error) {
	compositeKey := r.makeCompositeKey(key, r.stateSerde.GetStateType())
	exists, err := r.redisClient.Exists(compositeKey).Result()
	if err != nil {
		return 0, err
	}
	if exists == 0 {
		return 0, nil
	}
	serializedState, err := r.redisClient.Get(compositeKey).Result()
	if err != nil {
		return 0, err
	}
	redisState := &protos.RedisState{}
	err = proto.Unmarshal([]byte(serializedState), redisState)
	if err != nil {
		return 0, err
	}
	if redisState.GetIsGarbage() {
		return 0, fmt.Errorf("object found for key %s is garbage", key)
	}
	return redisState.GetVersion(), err
}

// MarkAsGarbage sets the object to be viewed as garbage. Once set, the object
// will no longer be returned on subsequent fetches. This object will be
// cleaned up via async garbage collection
func (r *RedisStateClient) MarkAsGarbage(key string) error {
	compositeKey := r.makeCompositeKey(key, r.stateSerde.GetStateType())
	serializedState, err := r.redisClient.Get(compositeKey).Result()
	if err != nil {
		return err
	}
	redisState := &protos.RedisState{}
	err = proto.Unmarshal([]byte(serializedState), redisState)
	if err != nil {
		return err
	}
	redisState.IsGarbage = true
	serializedProto, err := proto.Marshal(redisState)
	if err != nil {
		return nil
	}
	return r.redisClient.Set(compositeKey, serializedProto, 0).Err()
}

func (r *RedisStateClient) makeCompositeKey(key string, stateType string) string {
	return fmt.Sprintf("%s:%s", key, stateType)
}

// Serialize serializes the given interface, setting the state's version to the
// provided value
func (j *JsonStateSerde) Serialize(in interface{}, version uint64) ([]byte, error) {
	serializedJson, err := json.Marshal(in)
	if err != nil {
		return []byte{}, err
	}
	redisState := &protos.RedisState{
		SerializedMsg: serializedJson,
		Version:       version,
		IsGarbage:     false,
	}
	return proto.Marshal(redisState)
}

// Deserialize deserializes the provided byte slice into the an instance of the
// serde's data model
func (j *JsonStateSerde) Deserialize(in []byte) (interface{}, error) {
	redisState := &protos.RedisState{}
	err := proto.Unmarshal(in, redisState)
	if err != nil {
		return nil, err
	}
	model := reflect.New(reflect.TypeOf(j.dataInstance).Elem()).Interface()
	err = json.Unmarshal(redisState.GetSerializedMsg(), model)
	return model, err
}

// GetStateType returns the serde's state type
func (j *JsonStateSerde) GetStateType() string {
	return j.stateType
}
