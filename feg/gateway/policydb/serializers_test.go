/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package policydb_test

import (
	"testing"

	"magma/feg/gateway/policydb"
	lteProtos "magma/lte/cloud/go/protos"
	orc8rProtos "magma/orc8r/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestPolicySerializer(t *testing.T) {
	policy := getDefaultPolicy()
	policySerializer := policydb.GetPolicySerializer()

	// should return a serialized RedisState{SerializedMsg: <serialized policy>}
	serializedRedisState, err := policySerializer(policy)
	assert.NoError(t, err)

	expectedSerializedPolicyRedisState := createSerializedPolicyRedisState(t, policy)

	assert.Equal(t, expectedSerializedPolicyRedisState, serializedRedisState)
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
	clearOutMetaFields(policy)
	clearOutMetaFields(deserializedPolicy)

	assert.Equal(t, policy, deserializedPolicy)
}

func getDefaultPolicy() *lteProtos.PolicyRule {
	return &lteProtos.PolicyRule{
		Id:            "static1",
		Priority:      1,
		RatingGroup:   2,
		MonitoringKey: "mkey1",
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

func clearOutMetaFields(policy *lteProtos.PolicyRule) {
	policy.XXX_NoUnkeyedLiteral = struct{}{}
	policy.XXX_unrecognized = nil
	policy.XXX_sizecache = 0
}
