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

package handlers_test

import (
	"crypto/x509"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/device"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/security/key"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestListGateways(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	listGateways := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways", obsidian.GET).HandlerFunc

	// empty case
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Handler:        listGateways,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(map[string]models.MagmadGateway{}),
	}
	tests.RunUnitTest(t, e, tc)

	// happy path
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: orc8r.MagmadGatewayType, Key: "g1", Config: &models.MagmadGatewayConfigs{}, PhysicalID: "hw1"},
			{Type: orc8r.MagmadGatewayType, Key: "g2", Config: &models.MagmadGatewayConfigs{CheckinInterval: 15}},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	expectedResult := map[string]models.MagmadGateway{
		"g1": {ID: "g1", Magmad: &models.MagmadGatewayConfigs{}},
		"g2": {ID: "g2", Magmad: &models.MagmadGatewayConfigs{CheckinInterval: 15}},
	}
	tc.ExpectedResult = tests.JSONMarshaler(expectedResult)
	tests.RunUnitTest(t, e, tc)

	// add device and state to g1
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)
	gatewayRecord := &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}}
	err = device.RegisterDevice("n1", orc8r.AccessGatewayRecordType, "hw1", gatewayRecord, serdes.Device)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	expectedState := models.NewDefaultGatewayStatus("hw1")
	expectedState.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))

	expectedResult = map[string]models.MagmadGateway{
		"g1": {ID: "g1", Magmad: &models.MagmadGatewayConfigs{}, Device: gatewayRecord, Status: expectedState},
		"g2": {ID: "g2", Magmad: &models.MagmadGatewayConfigs{CheckinInterval: 15}},
	}
	tc.ExpectedResult = tests.JSONMarshaler(expectedResult)
	tests.RunUnitTest(t, e, tc)
}

func TestCreateGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	// create 2 tiers
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t1"}, serdes.Entity)
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.UpgradeTierEntityType, Key: "t2"}, serdes.Entity)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	createGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways", obsidian.POST).HandlerFunc

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
		Handler:        createGateway,
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "foo-bar-baz-123-42", serdes.Device)
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
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
		serdes.Device,
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err = device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hello-world-42", serdes.Device)
	assert.NoError(t, err)

	expectedEnts = configurator.NetworkEntities{
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
		Handler:        createGateway,
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
		Handler:        createGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		ExpectedStatus: 400,
		ExpectedError:  "device foo-bar-baz-123-42 is already mapped to gateway g1",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetGateway(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	getGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id", obsidian.GET).HandlerFunc

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
		Status: models.NewDefaultGatewayStatus("hw1"),
	}
	expected.Status.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))

	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1",
		Handler:        getGateway,
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
		Handler:        getGateway,
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
		Handler:        getGateway,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g3"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateGateway(t *testing.T) {
	stateTestInit.StartTestService(t)
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id", obsidian.PUT).HandlerFunc

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
		Handler:        updateGateway,
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	actualDevice, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.NoError(t, err)

	expectedEnts := configurator.NetworkEntities{
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

	// 400 mismatch gateway_id in parameter vs. payload
	tc = tests.Test{
		Method:                 "PUT",
		URL:                    testURLRoot + "/g3",
		Handler:                updateGateway,
		Payload:                payload,
		ParamNames:             []string{"network_id", "gateway_id"},
		ParamValues:            []string{"n1", "g3"},
		ExpectedStatus:         400,
		ExpectedErrorSubstring: "gateway ID from parameter (g3) and payload (g1) must match",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	payload.ID = "g3"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g3",
		Handler:        updateGateway,
		Payload:        payload,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g3"},
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestDeleteGateway(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	deleteGateway := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id", obsidian.DELETE).HandlerFunc

	tc := tests.Test{
		Method:         "DELETE",
		URL:            testURLRoot + "/g1",
		Handler:        deleteGateway,
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
		serdes.Entity,
	)
	assert.NoError(t, err)
	_, err = device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw1", serdes.Device)
	assert.EqualError(t, err, "Not found")

	expectedEnts := configurator.NetworkEntities{
		{NetworkID: "n1", Type: orc8r.UpgradeTierEntityType, Key: "t1", GraphID: "2"},
	}
	assert.Equal(t, expectedEnts, actualEnts)
}

func TestGetPartialReadHandlers(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(1000000, 0))
	defer clock.UnfreezeClock(t)

	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	gwConfig := &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 300,
		CheckinInterval:         15,
		CheckinTimeout:          5,
	}
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config:     gwConfig,
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g2",
				Name: "barfoo", Description: "bar foo",
				PhysicalID: "hw2",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, "hw1")
	test_utils.ReportGatewayStatus(t, ctx, models.NewDefaultGatewayStatus("hw1"))

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	getGatewayName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/name", obsidian.GET).HandlerFunc
	getGatewayDescription := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/description", obsidian.GET).HandlerFunc
	getGatewayState := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/status", obsidian.GET).HandlerFunc
	getGatewayDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/device", obsidian.GET).HandlerFunc
	getGatewayConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/magmad", obsidian.GET).HandlerFunc

	// happy path name
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/name",
		Handler:        getGatewayName,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("foobar"),
	}
	tests.RunUnitTest(t, e, tc)

	// happy path desc
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/description",
		Handler:        getGatewayDescription,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("foo bar"),
	}
	tests.RunUnitTest(t, e, tc)

	expectedState := models.NewDefaultGatewayStatus("hw1")
	expectedState.CheckinTime = uint64(time.Unix(1000000, 0).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))

	// happy path state
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/state",
		Handler:        getGatewayState,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: expectedState,
	}
	tests.RunUnitTest(t, e, tc)

	// 404 state
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g2/state",
		Handler:        getGatewayState,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path device
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/device",
		Handler:        getGatewayDevice,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: &models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
	}
	tests.RunUnitTest(t, e, tc)

	// 404 device
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g2/device",
		Handler:        getGatewayDevice,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case magmad config
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/magmad",
		Handler:        getGatewayConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: gwConfig,
	}
	tests.RunUnitTest(t, e, tc)

	// 404 magmad config
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g2/magmad",
		Handler:        getGatewayConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)
}

func TestGetGatewayTierHandler(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"
	obsidianHandlers := handlers.GetObsidianHandlers()
	getGatewayTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/tier", obsidian.GET).HandlerFunc

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	// 404 tier
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/tier",
		Handler:        getGatewayTier,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(""),
	}
	tests.RunUnitTest(t, e, tc)

	// add a tier and tier -> gateway association
	_, err = configurator.CreateEntity(
		"n1",
		configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType,
			Key:  "t1",
			Associations: []storage.TypeAndKey{
				{
					Type: orc8r.MagmadGatewayType,
					Key:  "g1",
				},
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// happy
	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot + "/g1/tier",
		Handler:        getGatewayTier,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("t1"),
	}
	tests.RunUnitTest(t, e, tc)
}

func TestUpdateGatewayTierHandler(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)
	stateTestInit.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"
	obsidianHandlers := handlers.GetObsidianHandlers()
	updateGatewayTier := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/tier", obsidian.PUT).HandlerFunc

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	// 404 tier
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/tier",
		Handler:        updateGatewayTier,
		Payload:        tests.JSONMarshaler(models.TierID("t1")),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 400,
		ExpectedError:  "Tier t1 does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	// add 2 tiers
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.UpgradeTierEntityType,
				Key:  "t1",
			},
			{
				Type: orc8r.UpgradeTierEntityType,
				Key:  "t2",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// happy add a tier
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/tier",
		Handler:        updateGatewayTier,
		Payload:        tests.JSONMarshaler(models.TierID("t1")),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entities, _, err := configurator.LoadEntities(
		"n1",
		swag.String(orc8r.UpgradeTierEntityType),
		nil,
		nil,
		nil,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	expectedTiers := configurator.NetworkEntities{
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType,
			Key:       "t1",
			Associations: []storage.TypeAndKey{
				{
					Type: orc8r.MagmadGatewayType,
					Key:  "g1",
				},
			},
			GraphID: "2",
			Version: 1,
		},
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType,
			Key:       "t2",
			GraphID:   "6",
			Version:   0,
		},
	}
	assert.Equal(t, expectedTiers, entities)

	// happy switch to a different tier
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/tier",
		Handler:        updateGatewayTier,
		Payload:        tests.JSONMarshaler(models.TierID("t2")),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	entities, _, err = configurator.LoadEntities(
		"n1",
		swag.String(orc8r.UpgradeTierEntityType),
		nil,
		nil,
		nil,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	expectedTiers = configurator.NetworkEntities{
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType,
			Key:       "t1",
			GraphID:   "2",
			Version:   2,
		},
		{
			NetworkID: "n1",
			Type:      orc8r.UpgradeTierEntityType,
			Key:       "t2",
			Associations: []storage.TypeAndKey{
				{
					Type: orc8r.MagmadGatewayType,
					Key:  "g1",
				},
			},
			GraphID: "6",
			Version: 1,
		},
	}
	assert.Equal(t, expectedTiers, entities)
}

func TestGetPartialUpdateHandlers(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"}, serdes.Network)
	assert.NoError(t, err)

	gwConfig := &models.MagmadGatewayConfigs{
		AutoupgradeEnabled:      swag.Bool(true),
		AutoupgradePollInterval: 300,
		CheckinInterval:         15,
		CheckinTimeout:          5,
	}
	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{
				Type: orc8r.MagmadGatewayType, Key: "g1",
				Name: "foobar", Description: "foo bar",
				PhysicalID: "hw1",
				Config:     gwConfig,
			},
			{
				Type: orc8r.MagmadGatewayType, Key: "g2",
				Name: "barfoo", Description: "bar foo",
				PhysicalID: "hw2",
			},
		},
		serdes.Entity,
	)
	assert.NoError(t, err)
	err = device.RegisterDevice(
		"n1", orc8r.AccessGatewayRecordType, "hw1",
		&models.GatewayDevice{HardwareID: "hw1", Key: &models.ChallengeKey{KeyType: "ECHO"}},
		serdes.Device,
	)
	assert.NoError(t, err)

	e := echo.New()
	testURLRoot := "/magma/v1/networks/n1/gateways"

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateGatewayName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/name", obsidian.PUT).HandlerFunc
	updateGatewayDesc := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/description", obsidian.PUT).HandlerFunc
	updateGatewayDevice := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/device", obsidian.PUT).HandlerFunc
	updateGatewayConfig := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/gateways/:gateway_id/magmad", obsidian.PUT).HandlerFunc

	// validation error name
	tc := tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/name",
		Handler:        updateGatewayName,
		Payload:        tests.JSONMarshaler(""),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 400,
		ExpectedError:  " in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	// happy path name
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/name",
		Handler:        updateGatewayName,
		Payload:        tests.JSONMarshaler("newname"),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entity, err := configurator.LoadEntity(
		"n1", orc8r.MagmadGatewayType, "g1",
		configurator.EntityLoadCriteria{LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, "newname", entity.Name)

	// happy path desc
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g1/description",
		Handler:        updateGatewayDesc,
		Payload:        tests.JSONMarshaler("newdesc"),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g1"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entity, err = configurator.LoadEntity(
		"n1", orc8r.MagmadGatewayType, "g1",
		configurator.EntityLoadCriteria{LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, "newdesc", entity.Description)

	// happy path device
	tc = tests.Test{
		Method:  "PUT",
		URL:     testURLRoot + "/g2/device",
		Handler: updateGatewayDevice,
		Payload: tests.JSONMarshaler(&models.GatewayDevice{HardwareID: "hw2",
			Key: &models.ChallengeKey{KeyType: "ECHO"}}),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	device, err := device.GetDevice("n1", orc8r.AccessGatewayRecordType, "hw2", serdes.Device)
	assert.NoError(t, err)
	assert.Equal(t, &models.GatewayDevice{HardwareID: "hw2", Key: &models.ChallengeKey{KeyType: "ECHO"}}, device)

	// happy case magmad config
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g2/magmad",
		Handler:        updateGatewayConfig,
		Payload:        gwConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g2"},
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	entity, err = configurator.LoadEntity(
		"n1", orc8r.MagmadGatewayType, "g1",
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, gwConfig, entity.Config)

	// 404 magmad config (no gateway registered)
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot + "/g3/magmad",
		Handler:        updateGatewayConfig,
		Payload:        gwConfig,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{"n1", "g3"},
		ExpectedStatus: 400,
		ExpectedError:  "Gateway g3 does not exist",
	}
	tests.RunUnitTest(t, e, tc)
}
