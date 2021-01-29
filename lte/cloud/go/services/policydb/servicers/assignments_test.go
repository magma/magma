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

package servicers_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	lteModels "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/policydb/servicers"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/storage"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestAssignmentsServicer(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)

	testNetworkId := "n1"
	testGwHwId := "hw1"
	testGwLogicalId := "g1"
	testSubscriberId := "s1"
	testPolicyId := "p1"
	testBaseName := "b1"

	// Initialize network
	err := configurator.CreateNetwork(configurator.Network{ID: testNetworkId}, serdes.Network)
	assert.NoError(t, err)

	// Initialize gateway -> subscriber, and create a policy rule
	_, err = configurator.CreateEntities(
		testNetworkId,
		[]configurator.NetworkEntity{
			{Type: lte.SubscriberEntityType, Key: testSubscriberId},
			{
				Type: lte.PolicyRuleEntityType,
				Key:  testPolicyId,
				Config: &models.PolicyRule{
					ID:                  models.PolicyID(testPolicyId),
					FlowList:            []*models.FlowDescription{},
					Priority:            swag.Uint32(5),
					RatingGroup:         *swag.Uint32(2),
					TrackingType:        "ONLY_OCS",
					AppName:             "INSTAGRAM",
					AssignedSubscribers: []models.SubscriberID{},
				},
			},
			{
				Type: lte.BaseNameEntityType,
				Key:  testBaseName,
				Config: &models.BaseNameRecord{
					AssignedSubscribers: []models.SubscriberID{},
					Name:                models.BaseName(testBaseName),
					RuleNames:           models.RuleNames([]string{}),
				},
			},
			{
				Type: lte.CellularGatewayEntityType, Key: testGwLogicalId,
				Config: newDefaultGatewayConfig(),
				Associations: []storage.TypeAndKey{
					{Type: lte.SubscriberEntityType, Key: testSubscriberId},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwLogicalId,
				Name: "foobar", Description: "foo bar",
				PhysicalID:   testGwHwId,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayEntityType, Key: testGwLogicalId}},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Create an identity and context for sending requests as gateway
	id := orcprotos.Identity{}
	idgw := orcprotos.Identity_Gateway{HardwareId: testGwHwId, NetworkId: testNetworkId, LogicalId: testGwLogicalId}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	srv := servicers.NewPolicyAssignmentServer()

	// Associate the rule to the subscriber, missing subscriber ID
	req := &protos.EnableStaticRuleRequest{Imsi: "s0", RuleIds: []string{testPolicyId}}
	_, err = srv.EnableStaticRules(ctx, req)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Either a subscriber or one more rules/basenames are not found")

	// Associate the rule to the subscriber, missing policy rule ID
	req = &protos.EnableStaticRuleRequest{Imsi: testSubscriberId, RuleIds: []string{"p0"}}
	_, err = srv.EnableStaticRules(ctx, req)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Either a subscriber or one more rules/basenames are not found")

	// Associate the basename to the subscriber, missing basename
	req = &protos.EnableStaticRuleRequest{Imsi: testSubscriberId, BaseNames: []string{"b0"}}
	_, err = srv.EnableStaticRules(ctx, req)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Either a subscriber or one more rules/basenames are not found")

	// Associate the rule to the subscriber successfully
	req = &protos.EnableStaticRuleRequest{Imsi: testSubscriberId, RuleIds: []string{testPolicyId}, BaseNames: []string{testBaseName}}
	_, err = srv.EnableStaticRules(ctx, req)
	assert.NoError(t, err)

	// Verify that the rule is associated to the subscriber
	ent, err := configurator.LoadEntity(
		testNetworkId, lte.PolicyRuleEntityType, testPolicyId,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	testPolicy := (&models.PolicyRule{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(testPolicy.AssignedSubscribers))

	// Verify that the base name is associated to the subscriber
	ent, err = configurator.LoadEntity(
		testNetworkId, lte.BaseNameEntityType, testBaseName,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	baseName := (&models.BaseNameRecord{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(baseName.AssignedSubscribers))

	// Disassociate the rule from the subscriber
	req2 := &protos.DisableStaticRuleRequest{Imsi: testSubscriberId, RuleIds: []string{testPolicyId}, BaseNames: []string{testBaseName}}
	_, err = srv.DisableStaticRules(ctx, req2)
	assert.NoError(t, err)

	// Verify that the rule is disassociated from the subscriber
	ent, err = configurator.LoadEntity(
		testNetworkId, lte.PolicyRuleEntityType, testPolicyId,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	testPolicy = (&models.PolicyRule{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(testPolicy.AssignedSubscribers))

	// Verify that the base name is disassociated from the subscriber
	ent, err = configurator.LoadEntity(
		testNetworkId, lte.BaseNameEntityType, testBaseName,
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	baseName = (&models.BaseNameRecord{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(baseName.AssignedSubscribers))
}

func newDefaultGatewayConfig() *lteModels.GatewayCellularConfigs {
	return &lteModels.GatewayCellularConfigs{
		Ran: &lteModels.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &lteModels.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &lteModels.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              []uint32{},
			NonEpsServiceControl: swag.Uint32(0),
		},
	}
}
