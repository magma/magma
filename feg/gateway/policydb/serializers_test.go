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

package policydb_test

import (
	"testing"

	"magma/feg/gateway/object_store"
	"magma/feg/gateway/policydb"
	lteProtos "magma/lte/cloud/go/protos"
	orc8rProtos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestPolicySerializer(t *testing.T) {
	policy := getDefaultPolicy()
	policySerializer := policydb.GetPolicySerializer()
	expectedString := createSerializedPolicyRedisState(t, policy)
	testRedisStateSerializer(t, policy, policySerializer, expectedString)
}

func TestPolicyDeserializer(t *testing.T) {
	policy := getDefaultPolicy()
	policyDeserializer := policydb.GetPolicyDeserializer()

	serializedPolicyRedisState := createSerializedPolicyRedisState(t, policy)
	iPolicy, err := policyDeserializer(serializedPolicyRedisState)
	assert.NoError(t, err)
	deserializedPolicy, ok := iPolicy.(*lteProtos.PolicyRule)
	assert.True(t, ok)

	// clear out meta fields
	clearOutMetaFieldsFromPolicy(policy)
	clearOutMetaFieldsFromPolicy(deserializedPolicy)

	assert.Equal(t, policy, deserializedPolicy)
}

func TestNameSetSerializer(t *testing.T) {
	nameSet := getDefaultNameSet()
	nameSetSerializer := policydb.GetBaseNameSerializer()
	expectedString := createSerializedNameSetRedisState(t, nameSet)

	testRedisStateSerializer(t, nameSet, nameSetSerializer, expectedString)
}

func TestNameSetDeserializer(t *testing.T) {
	nameSet := getDefaultNameSet()
	nameSetDeserializer := policydb.GetBaseNameDeserializer()

	serializedPolicyRedisState := createSerializedNameSetRedisState(t, nameSet)
	iPolicy, err := nameSetDeserializer(serializedPolicyRedisState)
	assert.NoError(t, err)
	deserializedPolicy, ok := iPolicy.(*lteProtos.ChargingRuleNameSet)
	assert.True(t, ok)

	// clear out meta fields
	clearOutMetaFieldsFromNameSet(nameSet)
	clearOutMetaFieldsFromNameSet(deserializedPolicy)

	assert.Equal(t, nameSet, deserializedPolicy)
}

func testRedisStateSerializer(t *testing.T, msg proto.Message, serializer object_store.Serializer, expectedString string) {
	// should return a serialized RedisState{SerializedMsg: <serialized policy>}
	serializedRedisState, err := serializer(msg)
	assert.NoError(t, err)
	assert.Equal(t, expectedString, serializedRedisState)
}

func getDefaultPolicy() *lteProtos.PolicyRule {
	return &lteProtos.PolicyRule{
		Id:            "static1",
		Priority:      1,
		RatingGroup:   2,
		MonitoringKey: []byte("mkey1"),
		Redirect:      nil,
		FlowList:      []*lteProtos.FlowDescription{{Action: lteProtos.FlowDescription_PERMIT}},
		Qos:           nil,
		TrackingType:  lteProtos.PolicyRule_OCS_AND_PCRF,
		HardTimeout:   0,
	}
}

func createSerializedPolicyRedisState(t *testing.T, policy *lteProtos.PolicyRule) string {
	serializedPolicy, err := proto.Marshal(policy)
	assert.NoError(t, err)
	expectedRedisState := &orc8rProtos.RedisState{
		SerializedMsg: serializedPolicy,
		Version:       0,
	}
	serializedRedisState, err := proto.Marshal(expectedRedisState)
	assert.NoError(t, err)
	return string(serializedRedisState)
}

func clearOutMetaFieldsFromPolicy(policy *lteProtos.PolicyRule) {
	policy.XXX_NoUnkeyedLiteral = struct{}{}
	policy.XXX_unrecognized = nil
	policy.XXX_sizecache = 0
}

func getDefaultNameSet() *lteProtos.ChargingRuleNameSet {
	return &lteProtos.ChargingRuleNameSet{
		RuleNames: []string{"static1"},
	}
}

func createSerializedNameSetRedisState(t *testing.T, nameSet *lteProtos.ChargingRuleNameSet) string {
	serializedNameSet, err := proto.Marshal(nameSet)
	assert.NoError(t, err)
	expectedRedisState := &orc8rProtos.RedisState{
		SerializedMsg: serializedNameSet,
		Version:       0,
	}
	serializedRedisState, err := proto.Marshal(expectedRedisState)
	assert.NoError(t, err)
	return string(serializedRedisState)
}

func clearOutMetaFieldsFromNameSet(nameSet *lteProtos.ChargingRuleNameSet) {
	nameSet.XXX_NoUnkeyedLiteral = struct{}{}
	nameSet.XXX_unrecognized = nil
	nameSet.XXX_sizecache = 0
}
