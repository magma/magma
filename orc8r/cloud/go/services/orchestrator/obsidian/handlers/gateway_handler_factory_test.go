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
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/configurator/test_utils"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_GetPartialReadGatewayHandler(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network without any configs
	networkID := "test-network"
	network := configurator.Network{
		ID:          networkID,
		Name:        "Test Network 1",
		Description: "Test Network 1",
	}

	gatewayRoot := fmt.Sprintf("%s/:network_id/gateways/:gateway_id", testURLRoot)

	// Test 404
	getGatewayName := handlers.GetPartialReadGatewayHandler(fmt.Sprintf("%s/Name", gatewayRoot), &testName{}, nil)
	tc := tests.Test{
		Method:         "GET",
		URL:            gatewayRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "gw1"},
		Payload:        nil,
		Handler:        getGatewayName.HandlerFunc,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	assert.NoError(t, configurator.CreateNetwork(network, serdes.Network))
	gateway := configurator.NetworkEntity{
		Key:  "gw1",
		Type: orc8r.MagmadGatewayType,
		Name: "gateway 1",
	}
	_, err := configurator.CreateEntity(networkID, gateway, serdes.Entity)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            gatewayRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "gw1"},
		Payload:        nil,
		Handler:        getGatewayName.HandlerFunc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(testName{Name: "gateway 1"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GetPartialUpdateGatewayHandler(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network without any configs
	networkID := "test-network"
	network := configurator.Network{
		ID:          networkID,
		Name:        "Test Network 1",
		Description: "Test Network 1",
	}

	gatewayRoot := fmt.Sprintf("%s/:network_id/gateways/:gateway_id", testURLRoot)

	// Test 404
	updateGatewayName := handlers.GetPartialUpdateGatewayHandler(fmt.Sprintf("%s/Name", gatewayRoot), &testName{}, nil)
	tc := tests.Test{
		Method:         "PUT",
		URL:            gatewayRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "test_gateway_1"},
		Payload:        tests.JSONMarshaler(&testName{Name: "updated Name!"}),
		Handler:        updateGatewayName.HandlerFunc,
		ExpectedStatus: 400,
		ExpectedError:  "Gateway test_gateway_1 does not exist",
	}
	tests.RunUnitTest(t, e, tc)

	assert.NoError(t, configurator.CreateNetwork(network, serdes.Network))
	Gateway := configurator.NetworkEntity{
		Key:  "test_gateway_1",
		Type: orc8r.MagmadGatewayType,
		Name: "Gateway 1",
	}
	_, err := configurator.CreateEntity(networkID, Gateway, serdes.Entity)
	assert.NoError(t, err)

	// validation failure
	tc = tests.Test{
		Method:         "PUT",
		URL:            gatewayRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "test_gateway_1"},
		Payload:        tests.JSONMarshaler(&testName{Name: ""}),
		Handler:        updateGatewayName.HandlerFunc,
		ExpectedStatus: 400,
		ExpectedError:  "Name cannot be empty",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "PUT",
		URL:            gatewayRoot,
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "test_gateway_1"},
		Payload:        tests.JSONMarshaler(&testName{Name: "updated Name!"}),
		Handler:        updateGatewayName.HandlerFunc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	Gateway, err = configurator.LoadEntity(
		networkID, orc8r.MagmadGatewayType, "test_gateway_1",
		configurator.EntityLoadCriteria{LoadMetadata: true},
		serdes.Entity,
	)
	assert.NoError(t, err)
	assert.Equal(t, "updated Name!", Gateway.Name)
}

func Test_GetGatewayDeviceHandler(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	deviceTestInit.StartTestService(t)
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network without any configs
	networkID := "test-network"
	test_utils.RegisterNetwork(t, networkID, "Name")

	gatewayRoot := fmt.Sprintf("%s/:network_id/gateways/:gateway_id", testURLRoot)
	getDevice := handlers.GetReadGatewayDeviceHandler(fmt.Sprintf("%s/device", gatewayRoot), serdes.Device)

	// 404
	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/device", gatewayRoot),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "test_gateway_1"},
		Handler:        getDevice.HandlerFunc,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, tc)

	test_utils.RegisterGateway(t, networkID, "test_gateway_1", &models.GatewayDevice{HardwareID: "test_hardware_id"})

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/device", gatewayRoot),
		ParamNames:     []string{"network_id", "gateway_id"},
		ParamValues:    []string{networkID, "test_gateway_1"},
		Handler:        getDevice.HandlerFunc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(&models.GatewayDevice{HardwareID: "test_hardware_id"}),
	}
	tests.RunUnitTest(t, e, tc)
}

type testName struct {
	Name string
}

func (m *testName) ValidateModel() error {
	if m == nil {
		return fmt.Errorf("Cannot be nil")
	}
	if len(m.Name) == 0 {
		return fmt.Errorf("Name cannot be empty")
	}
	return nil
}

func (m *testName) FromBackendModels(networkID string, gatewayID string) error {
	entity, err := configurator.LoadEntity(
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{LoadMetadata: true},
		serdes.Entity,
	)
	if err != nil {
		return err
	}
	m.Name = entity.Name
	return nil
}

func (m *testName) ToUpdateCriteria(networkID string, gatewayID string) ([]configurator.EntityUpdateCriteria, error) {
	exists, err := configurator.DoesEntityExist(networkID, orc8r.MagmadGatewayType, gatewayID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Gateway %s does not exist", gatewayID)
	}
	return []configurator.EntityUpdateCriteria{
		{
			Type:    orc8r.MagmadGatewayType,
			Key:     gatewayID,
			NewName: swag.String(m.Name),
		},
	}, nil
}
