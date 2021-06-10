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

package policydb_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"magma/feg/gateway/policydb"
	fegstreamer "magma/gateway/streamer"
	"magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/lib/go/definitions"
	orcprotos "magma/orc8r/lib/go/protos"
	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// Mock Cloud Streamer
type mockStreamProvider struct {
	name string
}

var (
	firstUpdateChan = make(chan struct{}, 100)
	onceTestsInit   sync.Once
)

func (m *mockStreamProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*orcprotos.DataUpdate, error) {
	// Data for stream name "base_names"
	rs1, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule11", "rule12"}})
	rs2, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule21", "rule22"}})

	// Data for stream name "policydb"
	pr, _ := proto.Marshal(&protos.PolicyRule{
		Id: "simple_match",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	})

	var updates []*orcprotos.DataUpdate

	// Determine what streamprovider it is before sending the updates
	switch m.name {
	case "base_names":
		updates = []*orcprotos.DataUpdate{
			{Key: "base_1", Value: rs1},
			{Key: "base_2", Value: rs2},
		}
	case "policydb":
		updates = []*orcprotos.DataUpdate{
			{Key: "simple_match", Value: pr},
		}
	default:
		updates = nil
	}

	go func() {
		time.Sleep(time.Millisecond * 100)
		firstUpdateChan <- struct{}{}
	}()
	return updates, nil
}

func (m *mockStreamProvider) GetExtraArgs() *any.Any {
	return nil
}

// Mock GW Cloud Service registry
type mockCloudRegistry struct {
	*platform_registry.ServiceRegistry
}

func (cr mockCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	if service != definitions.StreamerServiceName {
		return nil, fmt.Errorf("Not Implemented")
	}
	return platform_registry.GetConnection(streamer.ServiceName)
}

func (cr mockCloudRegistry) GetCloudConnectionFromServiceConfig(serviceConfig *config.ConfigMap, service string) (*grpc.ClientConn, error) {
	return nil, fmt.Errorf("Not Implemented")
}

type mockObjectStore struct {
	objMap sync.Map
}

func (os *mockObjectStore) Set(key string, object interface{}) error {
	os.objMap.Store(key, object)
	return nil
}

func (os *mockObjectStore) Get(key string) (interface{}, error) {
	val, ok := os.objMap.Load(key)
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return val, nil
}

func (os *mockObjectStore) Delete(key string) error {
	os.objMap.Delete(key)
	return nil
}

func (os *mockObjectStore) GetAll() (map[string]interface{}, error) {
	returnVals := make(map[string]interface{})
	addToMap := func(key, value interface{}) bool {
		keyStr, _ := key.(string)
		returnVals[keyStr] = value
		return true
	}
	os.objMap.Range(addToMap)
	return returnVals, nil
}

func (os *mockObjectStore) DeleteAll() error {
	valsByKey, err := os.GetAll()
	if err != nil {
		return err
	}
	for key := range valsByKey {
		os.Delete(key)
	}
	return nil
}

func initOnce(t *testing.T) {
	streamer_test_init.StartTestService(t)
}

func TestPolicyDBBaseNamesWithGRPC(t *testing.T) {
	onceTestsInit.Do(func() { initOnce(t) })
	streamer_test_init.StartNewTestProvider(t, &mockStreamProvider{name: "base_names"}, "base_names")
	dbClient := &policydb.RedisPolicyDBClient{
		PolicyMap:      &mockObjectStore{},
		BaseNameMap:    &mockObjectStore{},
		StreamerClient: fegstreamer.NewStreamerClient(mockCloudRegistry{}),
	}
	l := policydb.NewBaseNameStreamListener(dbClient.BaseNameMap)
	assert.NoError(t, dbClient.StreamerClient.AddListener(l))
	go dbClient.StreamerClient.Stream(l)

	select {
	case <-firstUpdateChan:
	case <-time.After(10 * time.Second):
		t.Fatal("PolicyDB base name update test timed out")
	}

	ruleIDs := dbClient.GetRuleIDsForBaseNames([]string{"base_1", "base_2"})
	assert.ElementsMatch(t, ruleIDs, []string{"rule11", "rule12", "rule21", "rule22"})
}

func TestPolicyDBRulesWithGRPC(t *testing.T) {
	onceTestsInit.Do(func() { initOnce(t) })
	streamer_test_init.StartNewTestProvider(t, &mockStreamProvider{name: "policydb"}, "policydb")
	dbClient := &policydb.RedisPolicyDBClient{
		PolicyMap:      &mockObjectStore{},
		BaseNameMap:    &mockObjectStore{},
		StreamerClient: fegstreamer.NewStreamerClient(mockCloudRegistry{}),
	}
	go dbClient.StreamerClient.Stream(policydb.NewPolicyDBStreamListener(dbClient.PolicyMap))

	select {
	case <-firstUpdateChan:
	case <-time.After(10 * time.Second):
		t.Fatal("PolicyDB rules update test timed out")
	}

	prExpected := &protos.PolicyRule{
		Id: "simple_match",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}

	prExpectedBytes, _ := proto.Marshal(prExpected)

	policyRule, _ := dbClient.GetPolicyRuleByID("simple_match")
	policyRuleBytes, _ := proto.Marshal(policyRule)

	assert.Equal(t, policyRuleBytes, prExpectedBytes)
}

func TestPolicyDBBaseNamesWithMockUpdates(t *testing.T) {
	dbClient := &policydb.RedisPolicyDBClient{
		PolicyMap:      &mockObjectStore{},
		BaseNameMap:    &mockObjectStore{},
		StreamerClient: fegstreamer.NewStreamerClient(mockCloudRegistry{}),
	}
	listener := policydb.NewBaseNameStreamListener(dbClient.BaseNameMap)
	go dbClient.StreamerClient.Stream(listener)

	rs1, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule11", "rule12"}})
	rs2, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule21", "rule22"}})
	updates := []*orcprotos.DataUpdate{
		{Key: "base_1", Value: rs1},
		{Key: "base_2", Value: rs2},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	ruleIDs := dbClient.GetRuleIDsForBaseNames([]string{"base_1", "base_2"})
	assert.ElementsMatch(t, ruleIDs, []string{"rule11", "rule12", "rule21", "rule22"})

	updates = []*orcprotos.DataUpdate{
		{Key: "base_1", Value: rs1},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	ruleIDs = dbClient.GetRuleIDsForBaseNames([]string{"base_1", "base_2"})
	assert.ElementsMatch(t, ruleIDs, []string{"rule11", "rule12"})

	rs3, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule31"}})
	updates = []*orcprotos.DataUpdate{
		{Key: "base_1", Value: rs1},
		{Key: "base_2", Value: rs3},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	ruleIDs = dbClient.GetRuleIDsForBaseNames([]string{"base_1", "base_2"})
	assert.ElementsMatch(t, ruleIDs, []string{"rule11", "rule12", "rule31"})
}

func TestPolicyDBRulesWithMockUpdates(t *testing.T) {
	dbClient := &policydb.RedisPolicyDBClient{
		PolicyMap:      &mockObjectStore{},
		BaseNameMap:    &mockObjectStore{},
		StreamerClient: fegstreamer.NewStreamerClient(mockCloudRegistry{}),
	}
	listener := policydb.NewPolicyDBStreamListener(dbClient.PolicyMap)
	go dbClient.StreamerClient.Stream(listener)

	// PolicyRules for the test
	prObject1 := &protos.PolicyRule{
		Id: "simple_match1",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}

	prObject2 := &protos.PolicyRule{
		Id: "simple_match2",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}

	prObject3 := &protos.PolicyRule{
		Id: "simple_match21",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}

	pr1, _ := proto.Marshal(prObject1)
	pr2, _ := proto.Marshal(prObject2)
	pr3, _ := proto.Marshal(prObject3)
	updates := []*orcprotos.DataUpdate{
		{Key: "simple_match1", Value: pr1},
		{Key: "simple_match2", Value: pr2},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	policyRule1, _ := dbClient.GetPolicyRuleByID("simple_match1")
	policyRule2, _ := dbClient.GetPolicyRuleByID("simple_match2")
	policyRule1Bytes, _ := proto.Marshal(policyRule1)
	policyRule2Bytes, _ := proto.Marshal(policyRule2)
	assert.Equal(t, policyRule1Bytes, pr1)
	assert.Equal(t, policyRule2Bytes, pr2)

	// Check that simple_match2 doesn't exist after this update
	updates = []*orcprotos.DataUpdate{
		{Key: "simple_match1", Value: pr1},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	policyRule11, _ := dbClient.GetPolicyRuleByID("simple_match1")
	policyRule21, err := dbClient.GetPolicyRuleByID("simple_match2")
	policyRule11Bytes, _ := proto.Marshal(policyRule11)

	assert.Equal(t, policyRule11Bytes, pr1)
	assert.NotNil(t, err)
	assert.Empty(t, policyRule21)

	// Check that simple_match1 updates to pr2 and simple_match2 adds pr3
	updates = []*orcprotos.DataUpdate{
		{Key: "simple_match1", Value: pr2},
		{Key: "simple_match2", Value: pr3},
	}
	listener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})
	policyRule12, _ := dbClient.GetPolicyRuleByID("simple_match1")
	policyRule22, _ := dbClient.GetPolicyRuleByID("simple_match2")
	policyRule12Bytes, _ := proto.Marshal(policyRule12)
	policyRule22Bytes, _ := proto.Marshal(policyRule22)
	assert.Equal(t, policyRule12Bytes, pr2)
	assert.Equal(t, policyRule22Bytes, pr3)
}

func TestOmnipresentRulesWithMockUpdates(t *testing.T) {
	dbClient := &policydb.RedisPolicyDBClient{
		PolicyMap:        &mockObjectStore{},
		BaseNameMap:      &mockObjectStore{},
		OmnipresentRules: &mockObjectStore{},
		StreamerClient:   fegstreamer.NewStreamerClient(mockCloudRegistry{}),
	}
	baseNameListener := policydb.NewBaseNameStreamListener(dbClient.BaseNameMap)
	omnipresentRulesListener := policydb.NewOmnipresentRulesListener(dbClient.OmnipresentRules)

	go dbClient.StreamerClient.Stream(baseNameListener)
	go dbClient.StreamerClient.Stream(omnipresentRulesListener)

	// base case
	ruleIDs, baseNames := dbClient.GetOmnipresentRules()
	assert.ElementsMatch(t, []string{}, ruleIDs)
	assert.ElementsMatch(t, []string{}, baseNames)

	// with update
	ruleSet1, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule11", "rule12"}})
	ruleSet2, _ := proto.Marshal(&protos.ChargingRuleNameSet{RuleNames: []string{"rule21", "rule22"}})
	updates := []*orcprotos.DataUpdate{
		{Key: "base_1", Value: ruleSet1},
		{Key: "base_2", Value: ruleSet2},
	}
	baseNameListener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})

	omnipresentRules, _ := proto.Marshal(&protos.AssignedPolicies{AssignedPolicies: []string{"rule1"}, AssignedBaseNames: []string{"base_1"}})
	updates = []*orcprotos.DataUpdate{
		{Key: "", Value: omnipresentRules},
	}
	omnipresentRulesListener.Update(&orcprotos.DataUpdateBatch{Updates: updates, Resync: true})

	ruleIDs, baseNames = dbClient.GetOmnipresentRules()
	assert.ElementsMatch(t, []string{"rule1"}, ruleIDs)
	assert.ElementsMatch(t, []string{"base_1"}, baseNames)
}
