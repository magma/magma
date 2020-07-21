/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"testing"
	"time"

	"magma/lte/cloud/go/lte"
	ltePlugin "magma/lte/cloud/go/plugin"
	lteModels "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb/obsidian/handlers"
	subscriberModels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers"
	handlers := handlers.GetHandlers()
	createSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.POST).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	// default sub profile should always succeed
	payload := &subscriberModels.Subscriber{
		ID:   "IMSI1234567890",
		Name: "Jane Doe",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "default",
		},
		ActiveApns: subscriberModels.ApnList{apn2, apn1},
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID:    "n1",
		Type:         lte.SubscriberEntityType,
		Key:          "IMSI1234567890",
		Name:         "Jane Doe",
		Config:       payload.Lte,
		GraphID:      "2",
		Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)

	// no cellular config on network and a non-default sub profile should be 500
	payload = &subscriberModels.Subscriber{
		ID: "IMSI0987654321",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "foo",
		},
		ActiveApns: subscriberModels.ApnList{apn2, apn1},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 500,
		ExpectedError:  "no cellular config found for network",
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI0987654321", configurator.FullEntityLoadCriteria())
	assert.EqualError(t, err, "Not found")

	// nonexistent sub profile should be 400
	err = configurator.UpdateNetworkConfig(
		"n1", lte.CellularNetworkType,
		&lteModels.NetworkCellularConfigs{
			Epc: &lteModels.NetworkEpcConfigs{
				SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{
					"blah": {
						MaxDlBitRate: 100,
						MaxUlBitRate: 100,
					},
				},
			},
		},
	)
	assert.NoError(t, err)
	payload = &subscriberModels.Subscriber{
		ID: "IMSI0987654321",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "foo",
		},
		ActiveApns: subscriberModels.ApnList{apn2, apn1},
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        payload,
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "subscriber profile foo does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)

	// other validation failure
	tc = tests.Test{
		Method: "POST",
		URL:    testURLRoot,
		Payload: &subscriberModels.Subscriber{
			ID: "IMSI1234567898",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
		},
		Handler:        createSubscriber,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "expected lte auth key to be 16 bytes but got 15 bytes",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestListSubscribers(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers"
	handlers := handlers.GetHandlers()
	listSubscribers := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{}),
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
				Config: &subscriberModels.LteSubscription{
					AuthAlgo: "MILENAGE",
					AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:    "ACTIVE",
				},
				Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
			},
			{
				Type: lte.SubscriberEntityType, Key: "IMSI0987654321",
				Config: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn1}},
			},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{
			"IMSI1234567890": {
				ID: "IMSI1234567890",
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				ActiveApns: subscriberModels.ApnList{apn2, apn1},
			},
			"IMSI0987654321": {
				ID: "IMSI0987654321",
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ActiveApns: subscriberModels.ApnList{apn1},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)

	// Now create some AGW-reported state for 1234567890
	// First we need to register a gateway which can report state
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw1"},
	)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)

	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listSubscribers,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]*subscriberModels.Subscriber{
			"IMSI1234567890": {
				ID: "IMSI1234567890",
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
					State:      "ACTIVE",
					SubProfile: "default",
				},
				ActiveApns: subscriberModels.ApnList{apn2, apn1},
				Monitoring: &subscriberModels.SubscriberStatus{
					Icmp: &subscriberModels.IcmpStatus{
						LastReportedTime: frozenClock,
						LatencyMs:        f32Ptr(12.34),
					},
				},
				State: &subscriberModels.SubscriberState{
					Mme:  mmeState,
					S1ap: s1apState,
					Spgw: spgwState,
					Mobility: []*subscriberModels.SubscriberIPAllocation{
						{
							Apn: "magma.apn",
							IP:  "192.168.128.134",
						},
						{
							Apn: "oai.ipv4",
							IP:  "192.168.128.174",
						},
					},
				},
			},
			"IMSI0987654321": {
				ID: "IMSI0987654321",
				Lte: &subscriberModels.LteSubscription{
					AuthAlgo:   "MILENAGE",
					AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
					State:      "ACTIVE",
					SubProfile: "foo",
				},
				ActiveApns: subscriberModels.ApnList{apn1},
			},
		}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	getSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.GET).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// No sub profile configured, we should return "default"
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Name: "Jane Doe",
			Config: &subscriberModels.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
			Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
		},
	)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.Subscriber{
			ID:   "IMSI1234567890",
			Name: "Jane Doe",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
		},
	}
	tests.RunUnitTest(t, e, tc)

	// Now create AGW
	// First we need to register a gateway which can report state
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw1"},
	)
	assert.NoError(t, err)
	frozenClock := int64(1000000)
	clock.SetAndFreezeClock(t, time.Unix(frozenClock, 0))
	defer clock.UnfreezeClock(t)
	icmpStatus := &subscriberModels.IcmpStatus{LatencyMs: f32Ptr(12.34)}
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportState(t, ctx, lte.ICMPStateType, "IMSI1234567890", icmpStatus)
	mmeState := state.ArbitraryJSON{"mme": "foo"}
	test_utils.ReportState(t, ctx, lte.MMEStateType, "IMSI1234567890", &mmeState)
	spgwState := state.ArbitraryJSON{"spgw": "foo"}
	test_utils.ReportState(t, ctx, lte.SPGWStateType, "IMSI1234567890", &spgwState)
	s1apState := state.ArbitraryJSON{"s1ap": "foo"}
	test_utils.ReportState(t, ctx, lte.S1APStateType, "IMSI1234567890", &s1apState)
	// Report 2 allocated IP addresses for the subscriber
	mobilitydState1 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiArg==",
		},
	}
	mobilitydState2 := state.ArbitraryJSON{
		"ip": map[string]interface{}{
			"address": "wKiAhg==",
		},
	}
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.oai.ipv4", &mobilitydState1)
	test_utils.ReportState(t, ctx, lte.MobilitydStateType, "IMSI1234567890.magma.apn", &mobilitydState2)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        getSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
		ExpectedResult: &subscriberModels.Subscriber{
			ID:   "IMSI1234567890",
			Name: "Jane Doe",
			Lte: &subscriberModels.LteSubscription{
				AuthAlgo:   "MILENAGE",
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			ActiveApns: subscriberModels.ApnList{apn2, apn1},
			Monitoring: &subscriberModels.SubscriberStatus{
				Icmp: &subscriberModels.IcmpStatus{
					LastReportedTime: frozenClock,
					LatencyMs:        f32Ptr(12.34),
				},
			},
			State: &subscriberModels.SubscriberState{
				Mme:  mmeState,
				S1ap: s1apState,
				Spgw: spgwState,
				Mobility: []*subscriberModels.SubscriberIPAllocation{
					{
						Apn: "magma.apn",
						IP:  "192.168.128.134",
					},
					{
						Apn: "oai.ipv4",
						IP:  "192.168.128.174",
					},
				},
			},
		},
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	updateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	// 404
	payload := &subscriberModels.Subscriber{
		ID: "IMSI1234567890",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "default",
		},
		ActiveApns: subscriberModels.ApnList{apn2, apn1},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateSubscriber,
		Payload:        payload,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// Happy path
	err = configurator.UpdateNetworkConfig(
		"n1", lte.CellularNetworkType,
		&lteModels.NetworkCellularConfigs{
			Epc: &lteModels.NetworkEpcConfigs{
				SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{
					"foo": {
						MaxUlBitRate: 100,
						MaxDlBitRate: 100,
					},
				},
			},
		},
	)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}},
		},
	)
	assert.NoError(t, err)

	payload = &subscriberModels.Subscriber{
		ID:   "IMSI1234567890",
		Name: "Jane Doe",
		Lte: &subscriberModels.LteSubscription{
			AuthAlgo:   "MILENAGE",
			AuthKey:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			AuthOpc:    []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			State:      "INACTIVE",
			SubProfile: "foo",
		},
		ActiveApns: subscriberModels.ApnList{apn2, apn1},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateSubscriber,
		Payload:        payload,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID:    "n1",
		Type:         lte.SubscriberEntityType,
		Key:          "IMSI1234567890",
		Name:         "Jane Doe",
		Config:       payload.Lte,
		GraphID:      "2",
		Version:      1,
		Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)

	// No profile matching
	payload.Lte.SubProfile = "bar"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateSubscriber,
		Payload:        payload,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 400,
		ExpectedError:  "subscriber profile bar does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	deleteSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.DELETE).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &subscriberModels.LteSubscription{
				AuthAlgo: "MILENAGE",
				AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:    "ACTIVE",
			},
			Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
		},
	)
	assert.NoError(t, err)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot,
		Handler:        deleteSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadAllEntitiesInNetwork("n1", lte.SubscriberEntityType, configurator.EntityLoadCriteria{})
	assert.Equal(t, 0, len(actual))
}

func TestActivateDeactivateSubscriber(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id"
	handlers := handlers.GetHandlers()
	activateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot+"/activate", obsidian.POST).HandlerFunc
	deactivateSubscriber := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot+"/deactivate", obsidian.POST).HandlerFunc

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	expected := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.LteSubscription{
			AuthAlgo: "MILENAGE",
			AuthKey:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:    "ACTIVE",
		},
		Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
	}
	_, err = configurator.CreateEntity("n1", expected)
	assert.NoError(t, err)
	expected.NetworkID = "n1"
	expected.GraphID = "2"
	expected.Version = 1

	// activate already activated subscriber
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot + "/activate",
		Handler:        activateSubscriber,
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 200,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// deactivate
	tc.URL = testURLRoot + "/deactivate"
	tc.Handler = deactivateSubscriber
	tests.RunUnitTest(t, e, tc)

	actual, err = configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.LteSubscription).State = "INACTIVE"
	expected.Version = 2
	assert.Equal(t, expected, actual)

	// deactivate deactivated sub
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.LteSubscription).State = "INACTIVE"
	expected.Version = 3
	assert.Equal(t, expected, actual)

	// activate
	tc.URL = testURLRoot + "/activate"
	tc.Handler = activateSubscriber
	tests.RunUnitTest(t, e, tc)
	actual, err = configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected.Config.(*subscriberModels.LteSubscription).State = "ACTIVE"
	expected.Version = 4
	assert.Equal(t, expected, actual)
}

func TestUpdateSubscriberProfile(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &ltePlugin.LteOrchestratorPlugin{})
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	err = configurator.UpdateNetworkConfig(
		"n1", lte.CellularNetworkType,
		&lteModels.NetworkCellularConfigs{
			Epc: &lteModels.NetworkEpcConfigs{
				SubProfiles: map[string]lteModels.NetworkEpcConfigsSubProfilesAnon{
					"foo": {
						MaxUlBitRate: 100,
						MaxDlBitRate: 100,
					},
				},
			},
		},
	)
	assert.NoError(t, err)

	//preseed 2 apns
	apn1, apn2 := "foo", "bar"
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.ApnEntityType, Key: apn1},
			{Type: lte.ApnEntityType, Key: apn2},
		},
	)
	assert.NoError(t, err)

	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
			Config: &subscriberModels.LteSubscription{
				AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				State:      "ACTIVE",
				SubProfile: "default",
			},
			Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
		},
	)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/lte/:network_id/subscribers/:subscriber_id/lte/sub_profile"
	handlers := handlers.GetHandlers()
	updateProfile := tests.GetHandlerByPathAndMethod(t, handlers, testURLRoot, obsidian.PUT).HandlerFunc

	// 404
	payload := "foo"
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI0987654321"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// bad profile
	payload = "bar"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 400,
		ExpectedError:  "subscriber profile bar does not exist for the network",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	payload = "foo"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err := configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	assert.NoError(t, err)
	expected := configurator.NetworkEntity{
		NetworkID: "n1", Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.LteSubscription{
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "foo",
		},
		GraphID:      "2",
		Version:      1,
		Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)

	// set to default
	payload = "default"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Handler:        updateProfile,
		Payload:        tests.JSONMarshaler(payload),
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", "IMSI1234567890"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actual, err = configurator.LoadEntity("n1", lte.SubscriberEntityType, "IMSI1234567890", configurator.FullEntityLoadCriteria())
	expected = configurator.NetworkEntity{
		NetworkID: "n1", Type: lte.SubscriberEntityType, Key: "IMSI1234567890",
		Config: &subscriberModels.LteSubscription{
			AuthKey:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			AuthOpc:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			State:      "ACTIVE",
			SubProfile: "default",
		},
		GraphID:      "2",
		Version:      2,
		Associations: []storage.TypeAndKey{{Type: lte.ApnEntityType, Key: apn2}, {Type: lte.ApnEntityType, Key: apn1}},
	}
	assert.Equal(t, expected, actual)
}

func f32Ptr(f float32) *float32 {
	return &f
}
