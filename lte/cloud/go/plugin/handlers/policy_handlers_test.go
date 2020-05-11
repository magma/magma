/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"fmt"
	"testing"

	"magma/lte/cloud/go/lte"
	lteplugin "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/plugin/handlers"
	"magma/lte/cloud/go/plugin/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// Basic API workflow tests
func TestPolicyDBHandlersBasic(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.LteNetworkType})
	assert.NoError(t, err)

	listPolicies := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.GET).HandlerFunc
	createPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	getPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.GET).HandlerFunc
	updatePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.PUT).HandlerFunc
	deletePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.DELETE).HandlerFunc
	listNames := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.GET).HandlerFunc
	createName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.POST).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.GET).HandlerFunc
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.PUT).HandlerFunc
	deleteName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.DELETE).HandlerFunc

	// Test empty response
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.PolicyRule{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test add policy rule
	testRule := &models.PolicyRule{
		ID: "PolicyRule1",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPV4Dst:   "42.42.42.42",
					IPV4Src:   "192.168.0.1/24",
					TCPDst:    2,
					TCPSrc:    1,
					UDPDst:    4,
					UDPSrc:    3,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that policy rule was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.PolicyRule{
			"PolicyRule1": testRule,
		}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{"PolicyRule1"})
	tests.RunUnitTest(t, e, tc)

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update Rule Using URL based ID
	testRule.FlowList = []*models.FlowDescription{
		{Action: swag.String("PERMIT"), Match: &models.FlowMatch{IPProto: swag.String("IPPROTO_ICMP"), Direction: swag.String("DOWNLINK")}},
	}
	testRule.Priority, testRule.RatingGroup, testRule.TrackingType = swag.Uint32(10), *swag.Uint32(3), "ONLY_OCS"
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        testRule,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        updatePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update results
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Delete a rule
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/policies/rules/PolicyRule1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "PolicyRule1"},
		Handler:        deletePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listPolicies,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.PolicyRule{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/rules"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test Multi Match Add Rule
	testRule = &models.PolicyRule{
		ID: "Test_mult",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_TCP"),
					TCPDst:    2,
					TCPSrc:    1,
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPV4Dst:   "42.42.42.42",
					IPV4Src:   "192.168.0.1/24",
					TCPDst:    2,
					TCPSrc:    1,
					UDPDst:    4,
					UDPSrc:    3,
				},
			},
		},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Read Rule Using URL based ID
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_mult",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_mult"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test adding a rule with QoS
	testRule = &models.PolicyRule{
		ID:           "Test_qos",
		FlowList:     []*models.FlowDescription{},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
		Qos: &models.FlowQos{
			MaxReqBwUl: swag.Uint32(2000),
			MaxReqBwDl: swag.Uint32(1000),
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that QoS rule was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_qos",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_qos"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)

	// Test adding rule with redirect information
	testRule = &models.PolicyRule{
		ID:           "Test_redirect",
		FlowList:     []*models.FlowDescription{},
		Priority:     swag.Uint32(5),
		RatingGroup:  *swag.Uint32(2),
		TrackingType: "ONLY_OCS",
		Redirect: &models.RedirectInformation{
			Support:       swag.String("ENABLED"),
			AddressType:   swag.String("URL"),
			ServerAddress: swag.String("127.0.0.1"),
		},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that rule with redirect was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/Test_redirect",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "Test_redirect"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)


	// Test add rule with app name match
	testRule = &models.PolicyRule{
		ID: "test_app_policy",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_ICMP"),
					IPV4Dst:   "42.42.42.42",
					IPV4Src:   "192.168.0.1/24",
					TCPDst:    2,
					TCPSrc:    1,
					UDPDst:    4,
					UDPSrc:    3,
				},
			},
		},
		Priority:       swag.Uint32(5),
		RatingGroup:    *swag.Uint32(2),
		TrackingType:   "ONLY_OCS",
		AppName:	    "INSTAGRAM"
		AppServiceType: "VIDEO"
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        testRule,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Check that rule with app name was added
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/rules/test_app_policy",
		Payload:        nil,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "test_app_policy"},
		Handler:        getPolicy,
		ExpectedStatus: 200,
		ExpectedResult: testRule,
	}
	tests.RunUnitTest(t, e, tc)


	// Now run base name test cases using the rules created above

	// Test Listing All Base Names
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.BaseNameRecord{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)

	// Test Add BaseName
	baseNameRecord := &models.BaseNameRecord{
		Name:      "Test",
		RuleNames: models.RuleNames{"Test_qos", "Test_redirect"},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/base_names",
		Payload:        baseNameRecord,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createName,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test Read BaseName Using URL based name
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: &models.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos", "Test_redirect"},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Test Update BaseName Using URL based name
	tc = tests.Test{
		Method: "PUT",
		URL:    "/magma/v1/networks/n1/policies/base_names/Test",
		Payload: &models.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos"},
		},
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Verify update BaseName
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: &models.BaseNameRecord{
			Name:      "Test",
			RuleNames: []string{"Test_qos"},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Get all BaseNames
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.BaseNameRecord{
			"Test": {
				Name:      "Test",
				RuleNames: []string{"Test_qos"},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{"Test"})
	tests.RunUnitTest(t, e, tc)

	// Delete a BaseName
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/networks/n1/policies/base_names/Test",
		Payload:        nil,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "Test"},
		Handler:        deleteName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// Confirm delete
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/networks/n1/policies/base_names?view=full",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        listNames,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*models.BaseNameRecord{}),
	}
	tests.RunUnitTest(t, e, tc)
	tc.URL = "/magma/v1/networks/n1/policies/base_names"
	tc.ExpectedResult = tests.JSONMarshaler([]string{})
	tests.RunUnitTest(t, e, tc)
}

// Associate base names and policies to subscribers
func TestPolicyHandlersAssociations(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &lteplugin.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	err := configurator.CreateNetwork(configurator.Network{ID: "n1", Type: lte.LteNetworkType})
	assert.NoError(t, err)

	createPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules", obsidian.POST).HandlerFunc
	getPolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.GET).HandlerFunc
	updatePolicy := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/rules/:rule_id", obsidian.PUT).HandlerFunc

	createName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names", obsidian.POST).HandlerFunc
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.GET).HandlerFunc
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/policies/base_names/:base_name", obsidian.PUT).HandlerFunc

	// preseed 3 subscribers
	imsi1, imsi2, imsi3 := "IMSI1234567890", "IMSI0987654321", "IMSI1111111111"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.SubscriberEntityType, Key: imsi1},
			{Type: lte.SubscriberEntityType, Key: imsi2},
			{Type: lte.SubscriberEntityType, Key: imsi3},
		},
	)
	assert.NoError(t, err)

	// Create rule assigned to s1, s2
	expectedP1 := &models.PolicyRule{
		AssignedSubscribers: []models.SubscriberID{models.SubscriberID(imsi2), models.SubscriberID(imsi1)},
		ID:                  "p1",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        expectedP1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	validatePolicy(
		t, e, getPolicy,
		expectedP1,
		configurator.NetworkEntity{
			NetworkID:    "n1",
			Type:         lte.PolicyRuleEntityType,
			Key:          "p1",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: imsi2}, {Type: lte.SubscriberEntityType, Key: imsi1}},
		},
	)

	// Update rule to assign to s3
	expectedP1.AssignedSubscribers = []models.SubscriberID{models.SubscriberID(imsi3)}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/rules/p1",
		Payload:        expectedP1,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", "p1"},
		Handler:        updatePolicy,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validatePolicy(
		t, e, getPolicy,
		expectedP1,
		configurator.NetworkEntity{
			NetworkID:    "n1",
			Type:         lte.PolicyRuleEntityType,
			Key:          "p1",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: imsi3}},
			Version:      1,
		},
	)

	// Create another policy p2 unbound to subs
	expectedP2 := &models.PolicyRule{
		ID: "p2",
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
				},
			},
		},
		Priority: swag.Uint32(1),
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/rules",
		Payload:        expectedP2,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createPolicy,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Create base name bound to p1 and imsi1, imsi2
	expectedBN := &models.BaseNameRecord{
		Name:                "b1",
		RuleNames:           models.RuleNames{"p1"},
		AssignedSubscribers: []models.SubscriberID{models.SubscriberID(imsi2), models.SubscriberID(imsi1)},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/networks/n1/policies/base_names",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createName,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "p1"},
				{Type: lte.SubscriberEntityType, Key: imsi2},
				{Type: lte.SubscriberEntityType, Key: imsi1},
			},
		},
	)

	// Update base name to bind to p2 and s3
	expectedBN = &models.BaseNameRecord{
		Name:                "b1",
		RuleNames:           models.RuleNames{"p2"},
		AssignedSubscribers: []models.SubscriberID{models.SubscriberID(imsi3)},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/base_names/b1",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "b1"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "p2"},
				{Type: lte.SubscriberEntityType, Key: imsi3},
			},
			Version: 1,
		},
	)

	// Update base name to bind to p1, p2, s1, s2, s3
	expectedBN = &models.BaseNameRecord{
		Name:                "b1",
		RuleNames:           models.RuleNames{"p1", "p2"},
		AssignedSubscribers: []models.SubscriberID{models.SubscriberID(imsi2), models.SubscriberID(imsi3), models.SubscriberID(imsi1)},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/networks/n1/policies/base_names/b1",
		Payload:        expectedBN,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", "b1"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	validateBaseName(
		t, e, getName,
		expectedBN,
		configurator.NetworkEntity{
			NetworkID: "n1",
			Type:      lte.BaseNameEntityType,
			Key:       "b1",
			GraphID:   "10",
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "p1"},
				{Type: lte.PolicyRuleEntityType, Key: "p2"},
				{Type: lte.SubscriberEntityType, Key: imsi2},
				{Type: lte.SubscriberEntityType, Key: imsi3},
				{Type: lte.SubscriberEntityType, Key: imsi1},
			},
			Version: 2,
		},
	)
}

// config will be filled from the expected model
func validatePolicy(t *testing.T, e *echo.Echo, getRule echo.HandlerFunc, expectedModel *models.PolicyRule, expectedEnt configurator.NetworkEntity) {
	expectedEnt.Config = getExpectedRuleConfig(expectedModel)

	actual, err := configurator.LoadEntity("n1", lte.PolicyRuleEntityType, string(expectedModel.ID), configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	assert.Equal(t, expectedEnt, actual)
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("/magma/v1/networks/n1/policies/rules/%s", expectedModel.ID),
		Payload:        expectedModel,
		ParamNames:     []string{"network_id", "rule_id"},
		ParamValues:    []string{"n1", string(expectedModel.ID)},
		Handler:        getRule,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}

func getExpectedRuleConfig(m *models.PolicyRule) *models.PolicyRuleConfig {
	return &models.PolicyRuleConfig{
		FlowList:      m.FlowList,
		MonitoringKey: m.MonitoringKey,
		Priority:      m.Priority,
		Qos:           m.Qos,
		RatingGroup:   m.RatingGroup,
		Redirect:      m.Redirect,
		TrackingType:  m.TrackingType,
	}
}

func validateBaseName(t *testing.T, e *echo.Echo, getName echo.HandlerFunc, expectedModel *models.BaseNameRecord, expectedEnt configurator.NetworkEntity) {
	actual, err := configurator.LoadEntity("n1", lte.BaseNameEntityType, string(expectedModel.Name), configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	assert.Equal(t, expectedEnt, actual)
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("/magma/v1/networks/n1/policies/base_names/%s", expectedModel.Name),
		Payload:        expectedModel,
		ParamNames:     []string{"network_id", "base_name"},
		ParamValues:    []string{"n1", string(expectedModel.Name)},
		Handler:        getName,
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)
}
