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

package policydb

import (
	"fmt"

	"magma/feg/gateway/object_store"
	lteProtos "magma/lte/cloud/go/protos"
	orc8rProtos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
)

func getRedisStateSerializer(serializer object_store.Serializer) object_store.Serializer {
	return func(object interface{}) (string, error) {
		serialized, err := serializer(object)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		redisState := &orc8rProtos.RedisState{SerializedMsg: []byte(serialized), Version: 0}
		bytes, err := proto.Marshal(redisState)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		return string(bytes[:]), nil
	}
}

func getRedisStateDeserializer(deserializer object_store.Deserializer) object_store.Deserializer {
	return func(serialized string) (interface{}, error) {
		redisStatePtr := &orc8rProtos.RedisState{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, redisStatePtr)
		if err != nil {
			return nil, err
		}
		return deserializer(string(redisStatePtr.SerializedMsg))
	}
}

func GetPolicySerializer() object_store.Serializer {
	serializer := func(object interface{}) (string, error) {
		policy, ok := object.(*lteProtos.PolicyRule)
		if !ok {
			return "", fmt.Errorf("Could not cast object to protobuf")
		}
		bytes, err := proto.Marshal(policy)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		return string(bytes[:]), nil
	}
	return getRedisStateSerializer(serializer)
}

func GetPolicyDeserializer() object_store.Deserializer {
	deserializer := func(serialized string) (interface{}, error) {
		policyPtr := &lteProtos.PolicyRule{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, policyPtr)
		if err != nil {
			return nil, err
		}
		return policyPtr, nil
	}
	return getRedisStateDeserializer(deserializer)
}

func GetBaseNameSerializer() object_store.Serializer {
	serializer := func(object interface{}) (string, error) {
		setPtr, ok := object.(*lteProtos.ChargingRuleNameSet)
		if !ok {
			return "", fmt.Errorf("Could not cast object to protobuf")
		}
		bytes, err := proto.Marshal(setPtr)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		return string(bytes[:]), nil
	}
	return getRedisStateSerializer(serializer)
}

func GetBaseNameDeserializer() object_store.Deserializer {
	deserializer := func(serialized string) (interface{}, error) {
		setPtr := &lteProtos.ChargingRuleNameSet{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, setPtr)
		if err != nil {
			return nil, err
		}
		return setPtr, nil
	}
	return getRedisStateDeserializer(deserializer)
}

func GetRuleMappingSerializer() object_store.Serializer {
	serializer := func(object interface{}) (string, error) {
		setPtr, ok := object.(*lteProtos.AssignedPolicies)
		if !ok {
			return "", fmt.Errorf("Could not cast object to protobuf")
		}
		bytes, err := proto.Marshal(setPtr)
		if err != nil {
			return "", fmt.Errorf("Could not marshal message")
		}
		return string(bytes[:]), nil
	}
	return getRedisStateSerializer(serializer)
}

func GetRuleMappingDeserializer() object_store.Deserializer {
	deserializer := func(serialized string) (interface{}, error) {
		setPtr := &lteProtos.AssignedPolicies{}
		bytes := []byte(serialized)
		err := proto.Unmarshal(bytes, setPtr)
		if err != nil {
			return nil, err
		}
		return setPtr, nil
	}
	return getRedisStateDeserializer(deserializer)
}
