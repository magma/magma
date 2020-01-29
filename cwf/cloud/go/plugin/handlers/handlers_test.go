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

	"magma/cwf/cloud/go/cwf"
	plugin2 "magma/cwf/cloud/go/plugin"
	"magma/cwf/cloud/go/plugin/handlers"
	models2 "magma/cwf/cloud/go/plugin/models"
	"magma/feg/cloud/go/feg"
	plugin3 "magma/feg/cloud/go/plugin"
	models3 "magma/feg/cloud/go/plugin/models"
	plugin4 "magma/lte/cloud/go/plugin"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCwfNetworks(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin4.LteOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.CwfOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin3.FegOrchestratorPlugin{})
	test_init.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listNetworks := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf", obsidian.GET).HandlerFunc
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf", obsidian.POST).HandlerFunc
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id", obsidian.GET).HandlerFunc
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id", obsidian.PUT).HandlerFunc
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id", obsidian.DELETE).HandlerFunc
	getNetworkFederationConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/federation", obsidian.GET).HandlerFunc
	getCarrierWifiConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/carrier_wifi", obsidian.GET).HandlerFunc
	getSubscriberDirectory := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/subscribers/:subscriber_id/directory_record", obsidian.GET).HandlerFunc

	// Test ListNetworks
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	seedCwfNetworks(t)
	fegNetworkID := "n5"

	// Test CreateNetwork
	tc = tests.Test{
		Method: "POST",
		URL:    "/magma/v1/cwf",
		Payload: tests.JSONMarshaler(
			&models2.CwfNetwork{
				CarrierWifi: models2.NewDefaultNetworkCarrierWifiConfigs(),
				Federation: &models3.FederatedNetworkConfigs{
					FegNetworkID: &fegNetworkID,
				},
				Description: "Foo Bar",
				DNS:         models.NewDefaultDNSConfig(),
				Features:    models.NewDefaultFeaturesConfig(),
				ID:          "n4",
				Name:        "foobar",
			},
		),
		Handler:        createNetwork,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test ListNetworks
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", "n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetNetwork
	expectedN1 := &models2.CwfNetwork{
		CarrierWifi: models2.NewDefaultNetworkCarrierWifiConfigs(),
		Federation: &models3.FederatedNetworkConfigs{
			FegNetworkID: &fegNetworkID,
		},
		Description: "Foo Bar",
		DNS:         models.NewDefaultDNSConfig(),
		Features:    models.NewDefaultFeaturesConfig(),
		ID:          "n1",
		Name:        "foobar",
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedN1),
	}
	tests.RunUnitTest(t, e, tc)

	// Test UpdateNetwork
	payloadN1 := &models2.CwfNetwork{
		ID:          "n1",
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		CarrierWifi: models2.NewDefaultModifiedNetworkCarrierWifiConfigs(),
		Federation: &models3.FederatedNetworkConfigs{
			FegNetworkID: &fegNetworkID,
		},
		Features: &models.NetworkFeatures{
			Features: map[string]string{
				"bar": "baz",
				"baz": "quz",
			},
		},
		DNS: &models.NetworkDNSConfig{
			EnableCaching: swag.Bool(true),
			LocalTTL:      swag.Uint32(120),
			Records: []*models.DNSConfigRecord{
				{
					Domain:  "foobar.com",
					ARecord: []strfmt.IPv4{"127.0.0.1", "127.0.0.2"},
					AaaaRecord: []strfmt.IPv6{
						"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
						"1234:0db8:85a3:0000:0000:8a2e:0370:1234",
					},
				},
				{
					Domain:  "facebook.com",
					ARecord: []strfmt.IPv4{"127.0.0.3"},
				},
			},
		},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/cwf/n1",
		Payload:        payloadN1,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualN1, err := configurator.LoadNetwork("n1", true, true)
	assert.NoError(t, err)
	expected := configurator.Network{
		ID:          "n1",
		Type:        cwf.CwfNetworkType,
		Name:        "updated foobar",
		Description: "Updated Foo Bar",
		Configs: map[string]interface{}{
			cwf.CwfNetworkType: models2.NewDefaultModifiedNetworkCarrierWifiConfigs(),
			feg.FederatedNetworkType: &models3.FederatedNetworkConfigs{
				FegNetworkID: &fegNetworkID,
			},
			orc8r.DnsdNetworkType:       payloadN1.DNS,
			orc8r.NetworkFeaturesConfig: payloadN1.Features,
		},
		Version: 1,
	}
	assert.Equal(t, expected, actualN1)

	// Test GetFederationPartialGet
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/federation",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkFederationConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(&models3.FederatedNetworkConfigs{
			FegNetworkID: &fegNetworkID,
		}),
		ExpectedError: "",
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetCarrierWifiPartialGet
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/carrier_wifi",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getCarrierWifiConfig,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models2.NewDefaultModifiedNetworkCarrierWifiConfigs()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// setup fixtures in backend
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
	)
	assert.NoError(t, err)

	seedCwfGateway(t)

	reqRecord := &directoryd.DirectoryRecord{
		LocationHistory: []string{"hw1"},
		Identifiers: map[string]interface{}{
			"mac_addr":  "aa:aa:aa:aa:aa:aa",
			"ipv4_addr": "192.168.127.1",
		},
	}
	expectedRecord := &models2.CwfSubscriberDirectoryRecord{
		LocationHistory: []string{"hw1"},
		MacAddr:         "aa:aa:aa:aa:aa:aa",
		IPV4Addr:        "192.168.127.1",
	}
	subID := "IMSI123456"
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	reportSubscriberDirectoryRecord(t, ctx, subID, reqRecord)

	// Test GetSubscriberDirectoryRecord
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/subscribers/IMSI123456/directory_record",
		ParamNames:     []string{"network_id", "subscriber_id"},
		ParamValues:    []string{"n1", subID},
		Handler:        getSubscriberDirectory,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedRecord),
	}
	tests.RunUnitTest(t, e, tc)

	// Test DeleteNetwork
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/cwf/n1",
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf",
		Payload:        nil,
		Handler:        listNetworks,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n3", "n4"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestCwfGateways(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = plugin.RegisterPluginForTests(t, &plugin2.CwfOrchestratorPlugin{})
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.GetUnfreezeClockDeferFunc(t)()
	test_init.StartTestService(t)
	stateTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)

	e := echo.New()

	obsidianHandlers := handlers.GetHandlers()
	listGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways", obsidian.GET).HandlerFunc
	getGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways/:gateway_id", obsidian.GET).HandlerFunc
	getCarrierWifiGatewayConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways/:gateway_id/carrier_wifi", obsidian.GET).HandlerFunc
	updateCarrierWifiGatewayConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways/:gateway_id/carrier_wifi", obsidian.PUT).HandlerFunc
	updateGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways/:gateway_id", obsidian.PUT).HandlerFunc
	deleteGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways/:gateway_id", obsidian.DELETE).HandlerFunc
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways", obsidian.POST).HandlerFunc
	seedCwfNetworks(t)

	// setup fixtures in backend
	_, err := configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
	)
	assert.NoError(t, err)

	// Test CreateGateway
	payload := &models2.MutableCwfGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		CarrierWifi: &models2.GatewayCwfConfigs{
			AllowedGrePeers: models2.AllowedGrePeers{
				{IP: "1.1.1.1"},
			},
		},
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          5,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/cwf/n1/gateways",
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// Test ListGateways
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	expected := map[string]*models2.CwfGateway{
		"g1": {
			ID: "g1",
			Device: &models.GatewayDevice{
				HardwareID: "hw1",
				Key:        &models.ChallengeKey{KeyType: "ECHO"},
			},
			Name: "foobar", Description: "foo bar",
			CarrierWifi: &models2.GatewayCwfConfigs{
				AllowedGrePeers: models2.AllowedGrePeers{
					{IP: "1.1.1.1"},
				},
			},
			Tier: "t1",
			Magmad: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         15,
				CheckinTimeout:          5,
			},
			Status: models.NewDefaultGatewayStatus("hw1"),
		},
	}
	expected["g1"].Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expected["g1"].Status.CertExpirationTime = time.Unix(1000000, 0).Add(time.Hour * 4).Unix()
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways",
		Handler:        listGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expected),
	}
	tests.RunUnitTest(t, e, tc)

	// Test GetGateway
	expectedGet := &models2.CwfGateway{
		ID: "g1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Name: "foobar", Description: "foo bar",
		CarrierWifi: &models2.GatewayCwfConfigs{
			AllowedGrePeers: models2.AllowedGrePeers{
				{IP: "1.1.1.1"},
			},
		},
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Status: models.NewDefaultGatewayStatus("hw1"),
	}
	expectedGet.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expectedGet.Status.CertExpirationTime = time.Unix(1000000, 0).Add(time.Hour * 4).Unix()
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)

	// Test UpdateGateway
	payload = &models2.MutableCwfGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "newname",
		Description: "bar baz",
		CarrierWifi: &models2.GatewayCwfConfigs{
			AllowedGrePeers: models2.AllowedGrePeers{{IP: "1.1.1.1"}},
		},
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Tier: "t1",
	}

	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        updateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	expectedGet.Name = "newname"
	expectedGet.Description = "bar baz"
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGet,
	}
	tests.RunUnitTest(t, e, tc)

	// Test get gateway CarrierWifi config
	expectedGwConfGet := &models2.GatewayCwfConfigs{
		AllowedGrePeers: models2.AllowedGrePeers{
			{IP: "1.1.1.1"},
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        getCarrierWifiGatewayConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedGwConfGet,
	}
	tests.RunUnitTest(t, e, tc)

	// Test update gateway CarrierWifi config (invalid config)
	badPayloadConf := &models2.GatewayCwfConfigs{
		AllowedGrePeers: models2.AllowedGrePeers{
			{IP: "2.2.2.2/24", Key: swag.Uint32(444)},
			{IP: "2.2.2.2/24", Key: swag.Uint32(444)},
		},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        updateCarrierWifiGatewayConfig,
		Payload:        badPayloadConf,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 400,
		ExpectedError:  "Found duplicate peer 2.2.2.2/24 with key 444",
	}
	tests.RunUnitTest(t, e, tc)

	// Test update gateway CarrierWifi config
	payloadConf := &models2.GatewayCwfConfigs{
		AllowedGrePeers: models2.AllowedGrePeers{
			{IP: "2.2.2.2/24", Key: swag.Uint32(321)},
			{IP: "2.2.2.3/24", Key: swag.Uint32(321)},
		},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        updateCarrierWifiGatewayConfig,
		Payload:        payloadConf,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        getCarrierWifiGatewayConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: payloadConf,
	}
	tests.RunUnitTest(t, e, tc)

	// Test DeleteGateway
	tc = tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        deleteGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	tc = tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/cwf/n1/gateways/g1",
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
}

// n1, n3 are cwf networks, n2, n5 are not
func seedCwfNetworks(t *testing.T) {
	fegNetworkID := "n5"
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          fegNetworkID,
				Type:        feg.FederationNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					feg.FegNetworkType:          models3.NewDefaultNetworkFederationConfigs(),
					orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
					orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
				},
			},
		},
	)
	assert.NoError(t, err)
	_, err = configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        cwf.CwfNetworkType,
				Name:        "foobar",
				Description: "Foo Bar",
				Configs: map[string]interface{}{
					cwf.CwfNetworkType: models2.NewDefaultNetworkCarrierWifiConfigs(),
					feg.FederatedNetworkType: &models3.FederatedNetworkConfigs{
						FegNetworkID: &fegNetworkID,
					},
					orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
					orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
				},
			},
			{
				ID:          "n2",
				Type:        "blah",
				Name:        "foobar",
				Description: "Foo Bar",
				Configs:     map[string]interface{}{},
			},
			{
				ID:          "n3",
				Type:        cwf.CwfNetworkType,
				Name:        "barfoo",
				Description: "Bar Foo",
				Configs:     map[string]interface{}{},
			},
		},
	)
	assert.NoError(t, err)
}

func reportSubscriberDirectoryRecord(t *testing.T, ctx context.Context, id string, req *directoryd.DirectoryRecord) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedRecord, err := serde.Serialize(state.SerdeDomain, orc8r.DirectoryRecordType, req)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     orc8r.DirectoryRecordType,
			DeviceID: id,
			Value:    serializedRecord,
			Version:  1,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

func seedCwfGateway(t *testing.T) {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/cwf/:network_id/gateways", obsidian.POST).HandlerFunc

	payload := &models2.MutableCwfGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		CarrierWifi: &models2.GatewayCwfConfigs{
			AllowedGrePeers: models2.AllowedGrePeers{
				{IP: "1.1.1.1/24"},
			},
		},
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          5,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/cwf/n1/gateways",
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
}
