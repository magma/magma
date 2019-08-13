/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers_test

import (
	"crypto/x509"
	"os"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/handlers"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/security/key"
	checkind_test_utils "magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestListGateways(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	// empty case
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        handlers.ListGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.MagmadGatewayType, Key: "gz"},
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.MagmadGatewayType, Key: "g3"},
		},
	)
	assert.NoError(t, err)
	tc.ExpectedResult = tests.JSONMarshaler([]string{"g1", "g2", "g3", "gz"})
	tests.RunUnitTest(t, e, tc)
}

func TestCreateGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	// create 2 tiers
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t2"})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	// Register device with gateway
	payload := &models.MagmadGateway{
		Device: &models.GatewayDevice{
			HardwareID: "foo-bar-baz-123-42",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g1",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t1",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        handlers.CreateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)

	// verify a few things
	// gateway should have been created
	// device should have been created
	// tier should have an updated assoc
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "foo-bar-baz-123-42")
	assert.NoError(t, err)

	expectedEnts := []configurator.NetworkEntity{
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "foo-bar-baz-123-42",
			Config:             payload.Magmad,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t1"}},
			GraphID:            "2",
		},
		{
			NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:      "2",
			Version:      1,
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, payload.Device, actualDevice)

	// test registering gateway with existing device
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hello-world-42",
		&models.GatewayDevice{
			HardwareID: "hello-world-42",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
	)
	assert.NoError(t, err)

	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)
	pubkeyB64 := strfmt.Base64(marshaledPubKey)
	payload = &models.MagmadGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hello-world-42",
			Key: &models.ChallengeKey{
				KeyType: "SOFTWARE_ECDSA_SHA256",
				Key:     &pubkeyB64,
			},
		},
		ID:          "g2",
		Name:        "barfoo",
		Description: "bar foo",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t2",
	}
	tc.Payload = payload
	tests.RunUnitTest(t, e, tc)

	// verify results - device key should have changed
	actualEnts, _, err = configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g2"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t2"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err = device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hello-world-42")
	assert.NoError(t, err)

	expectedEnts = []configurator.NetworkEntity{
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g2",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "hello-world-42",
			Config:             payload.Magmad,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t2"}},
			GraphID:            "4",
		},
		{
			NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g2"}},
			GraphID:      "4",
			Version:      1,
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, payload.Device, actualDevice)

	// bad tier ID
	payload = &models.MagmadGateway{
		Device: &models.GatewayDevice{
			HardwareID: "doesnt-matter",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g3",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t3",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        handlers.CreateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "requested tier does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	// device already registered
	payload = &models.MagmadGateway{
		Device: &models.GatewayDevice{
			HardwareID: "foo-bar-baz-123-42",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		ID:          "g4",
		Name:        "foobar",
		Description: "foo bar",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		Tier: "t1",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Handler:        handlers.CreateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "device foo-bar-baz-123-42 is already mapped to gateway g1",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.GetUnfreezeClockDeferFunc(t)()

	test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	state_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g2",
				Name: "barfoo", Description: "bar foo",
				PhysicalID: "hw2",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{
					{Type: orc8r.MagmadGatewayType, Key: "g1"},
					{Type: orc8r.MagmadGatewayType, Key: "g2"},
				},
			},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, checkind_test_utils.GetGatewayStatusSwaggerFixture("hw1"))

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	// happy path
	expected := &models.MagmadGateway{
		ID: "g1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Name: "foobar", Description: "foo bar",
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Status: checkind_test_utils.GetGatewayStatusSwaggerFixture("hw1"),
	}
	expected.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
	expected.Status.CertExpirationTime = time.Unix(1000000, 0).Add(time.Hour * 4).Unix()

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1",
		Handler:        handlers.GetGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expected,
	}
	tests.RunUnitTest(t, e, tc)

	// get a gateway without a device or status
	expected = &models.MagmadGateway{
		ID:   "g2",
		Name: "barfoo", Description: "bar foo",
		Tier: "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
	}
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g2",
		Handler:        handlers.GetGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 200,
		ExpectedResult: expected,
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g3",
		Handler:        handlers.GetGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g3"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")

	test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			},
			{Type: orc8r.UpgradeTierEntityType, Key: "t2"},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	// update everything
	privateKey, err := key.GenerateKey("P256", 0)
	assert.NoError(t, err)
	marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(privateKey))
	assert.NoError(t, err)
	pubkeyB64 := strfmt.Base64(marshaledPubKey)
	payload := &models.MagmadGateway{
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "SOFTWARE_ECDSA_SHA256", Key: &pubkeyB64},
		},
		ID:          "g1",
		Name:        "barbaz",
		Description: "bar baz",
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         25,
			CheckinTimeout:          15,
			AutoupgradePollInterval: 200,
			AutoupgradeEnabled:      swag.Bool(false),
			FeatureFlags:            map[string]bool{"foo": false},
			DynamicServices:         []string{"d1", "d2"},
		},
		Tier: "t2",
	}

	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1",
		Handler:        handlers.UpdateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// load and validate
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t2"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1")
	assert.NoError(t, err)

	expectedEnts := []configurator.NetworkEntity{
		{
			NetworkID: "n1", Type: orc8r.MagmadGatewayType, Key: "g1",
			Name: string(payload.Name), Description: string(payload.Description),
			PhysicalID:         "hw1",
			Config:             payload.Magmad,
			ParentAssociations: []storage.TypeAndKey{{Type: orc8r.UpgradeTierEntityType, Key: "t2"}},
			GraphID:            "6",
			Version:            1,
		},
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "2", Version: 1},
		{
			NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t2",
			Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			GraphID:      "6",
			Version:      1,
		},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, payload.Device, actualDevice)

	// 404
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g3",
		Handler:        handlers.UpdateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g3"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteGateway(t *testing.T) {
	_ = plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")

	test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config: &models.MagmadGatewayConfigs{
					AutoupgradeEnabled:      swag.Bool(true),
					AutoupgradePollInterval: 300,
					CheckinInterval:         15,
					CheckinTimeout:          5,
				},
			},
			{
				Type: orc8r.UpgradeTierEntityType, Key: "t1",
				Associations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "g1"}},
			},
		},
	)
	assert.NoError(t, err)
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}})
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot + "/g1",
		Handler:        handlers.DeleteGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	// load, verify results
	actualEnts, _, err := configurator.LoadEntities(
		"n1", nil, nil, nil,
		[]storage.TypeAndKey{
			{Type: orc8r.MagmadGatewayType, Key: "g1"},
			{Type: orc8r.UpgradeTierEntityType, Key: "t1"},
		},
		configurator.FullEntityLoadCriteria(),
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1")
	assert.NoError(t, err)

	expectedEnts := []configurator.NetworkEntity{
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "2"},
	}
	assert.Equal(t, expectedEnts, actualEnts)
	assert.Equal(t, &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}}, actualDevice)
}
