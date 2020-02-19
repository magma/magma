/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	plugin2 "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
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
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)

	testNetworkId := "n1"
	testGwHwId := "hw1"
	testGwLogicalId := "g1"
	testSubscriberId := "s1"
	testPolicyId := "p1"
	testBaseName := "b1"

	// Initialize network
	err := configurator.CreateNetwork(configurator.Network{ID: testNetworkId})
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
				Type: lte.CellularGatewayType, Key: testGwLogicalId,
				Config: newDefaultGatewayConfig(),
				Associations: []storage.TypeAndKey{
					{Type: lte.SubscriberEntityType, Key: testSubscriberId},
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: testGwLogicalId,
				Name: "foobar", Description: "foo bar",
				PhysicalID:   testGwHwId,
				Associations: []storage.TypeAndKey{{Type: lte.CellularGatewayType, Key: testGwLogicalId}},
			},
		},
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
		testNetworkId,
		lte.PolicyRuleEntityType,
		testPolicyId,
		configurator.FullEntityLoadCriteria(),
	)
	testPolicy := (&models.PolicyRule{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(testPolicy.AssignedSubscribers))

	// Verify that the base name is associated to the subscriber
	ent, err = configurator.LoadEntity(
		testNetworkId,
		lte.BaseNameEntityType,
		testBaseName,
		configurator.FullEntityLoadCriteria(),
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
		testNetworkId,
		lte.PolicyRuleEntityType,
		testPolicyId,
		configurator.EntityLoadCriteria{LoadConfig: true},
	)
	testPolicy = (&models.PolicyRule{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(testPolicy.AssignedSubscribers))

	// Verify that the base name is disassociated from the subscriber
	ent, err = configurator.LoadEntity(
		testNetworkId,
		lte.BaseNameEntityType,
		testBaseName,
		configurator.FullEntityLoadCriteria(),
	)
	baseName = (&models.BaseNameRecord{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(baseName.AssignedSubscribers))
}

func newDefaultGatewayConfig() *models.GatewayCellularConfigs {
	return &models.GatewayCellularConfigs{
		Ran: &models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &models.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &models.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              []uint32{},
			NonEpsServiceControl: swag.Uint32(0),
		},
	}
}
