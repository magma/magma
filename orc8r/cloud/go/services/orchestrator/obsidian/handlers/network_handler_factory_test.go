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
	"encoding/json"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/lib/go/errors"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type (
	ID struct {
		Name string
	}
	TestFeature1 struct {
		ID   *ID
		Desc string
	}
)

func Test_GetPartialReadNetworkHandler(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	networkSerdes := serde.NewRegistry(configurator.NewNetworkConfigSerde("test", &TestFeature1{}))
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network without any configs
	networkID := "test-network"
	network := configurator.Network{
		ID:          networkID,
		Name:        "Test Network 1",
		Description: "Test Network 1",
	}
	assert.NoError(t, configurator.CreateNetwork(network, networkSerdes))

	networkURL := fmt.Sprintf("%s/%s", testURLRoot, networkID)

	// Test 404
	getFullConfig := handlers.GetPartialReadNetworkHandler(networkURL, &TestFeature1{}, networkSerdes)
	getFeatures := tests.Test{
		Method:         "GET",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        nil,
		Handler:        getFullConfig.HandlerFunc,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, getFeatures)

	getPartialConfig := handlers.GetPartialReadNetworkHandler(networkURL, &ID{}, networkSerdes)
	getName := tests.Test{
		Method:         "GET",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        nil,
		Handler:        getPartialConfig.HandlerFunc,
		ExpectedStatus: 404,
		ExpectedError:  "Not found",
	}
	tests.RunUnitTest(t, e, getName)

	// add config to network
	update := configurator.NetworkUpdateCriteria{
		ID: networkID,
		ConfigsToAddOrUpdate: map[string]interface{}{
			"test": &TestFeature1{ID: &ID{Name: "hello!"}, Desc: "goodbye!"},
		},
	}
	assert.NoError(t, configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update}, networkSerdes))

	// happy full case
	getFullConfig = handlers.GetPartialReadNetworkHandler(networkURL, &TestFeature1{}, networkSerdes)
	getFeatures = tests.Test{
		Method:         "GET",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        nil,
		Handler:        getFullConfig.HandlerFunc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(&TestFeature1{ID: &ID{Name: "hello!"}, Desc: "goodbye!"}),
	}
	tests.RunUnitTest(t, e, getFeatures)

	// happy partial case
	getPartialConfig = handlers.GetPartialReadNetworkHandler(networkURL, &ID{}, networkSerdes)
	getName = tests.Test{
		Method:         "GET",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        nil,
		Handler:        getPartialConfig.HandlerFunc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(&ID{Name: "hello!"}),
	}
	tests.RunUnitTest(t, e, getName)
}

func TestGetUpdateNetworkConfigHandler(t *testing.T) {
	networkSerdes := serde.NewRegistry(configurator.NewNetworkConfigSerde("test", &TestFeature1{}))
	configuratorTestInit.StartTestService(t)
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	networkID := "test-network"
	network := configurator.Network{
		ID:          networkID,
		Type:        "lte",
		Name:        "Test Network 1",
		Description: "Test Network 1",
		Configs:     map[string]interface{}{"test": &TestFeature1{ID: &ID{Name: "hello!"}, Desc: "goodbye!"}},
	}
	assert.NoError(t, configurator.CreateNetwork(network, networkSerdes))

	networkURL := fmt.Sprintf("%s/%s", testURLRoot, networkID)

	updateConfigsFull := handlers.GetPartialUpdateNetworkHandler(networkURL, &TestFeature1{}, networkSerdes)
	updateConfigsPartial := handlers.GetPartialUpdateNetworkHandler(networkURL, &ID{}, networkSerdes)

	// name is empty
	badConfig := &TestFeature1{ID: &ID{Name: ""}, Desc: "goodbye!"}
	updateFullConfig := tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        tests.JSONMarshaler(badConfig),
		Handler:        updateConfigsFull.HandlerFunc,
		ExpectedStatus: 400,
		ExpectedError:  "Name cannot be nil",
	}
	tests.RunUnitTest(t, e, updateFullConfig)

	badPartialConfig := &ID{Name: ""}
	updatePartialConfig := tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        tests.JSONMarshaler(badPartialConfig),
		Handler:        updateConfigsFull.HandlerFunc,
		ExpectedStatus: 400,
		ExpectedError:  "Cannot be nil",
	}
	tests.RunUnitTest(t, e, updatePartialConfig)

	expectedConfig := &TestFeature1{ID: &ID{Name: "hello world!"}, Desc: "goodbye world!"}
	// happy full case
	updateFullConfig = tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        tests.JSONMarshaler(expectedConfig),
		Handler:        updateConfigsFull.HandlerFunc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, updateFullConfig)

	config, err := configurator.LoadNetworkConfig(networkID, "test", networkSerdes)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)

	expectedConfig.ID.Name = "foo! bar!"
	// happy partial case
	updateFullConfig = tests.Test{
		Method:         "PUT",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Payload:        tests.JSONMarshaler(expectedConfig.ID),
		Handler:        updateConfigsPartial.HandlerFunc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, updateFullConfig)

	config, err = configurator.LoadNetworkConfig(networkID, "test", networkSerdes)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}

func TestGetDeleteNetworkConfigHandler(t *testing.T) {
	networkSerdes := serde.NewRegistry(configurator.NewNetworkConfigSerde("test", &TestFeature1{}))
	configuratorTestInit.StartTestService(t)
	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	networkID := "test-network"
	network := configurator.Network{
		ID:          networkID,
		Type:        "lte",
		Name:        "Test Network 1",
		Description: "Test Network 1",
		Configs:     map[string]interface{}{"test": &TestFeature1{ID: &ID{Name: "hello!"}, Desc: "goodbye!"}},
	}
	assert.NoError(t, configurator.CreateNetwork(network, networkSerdes))

	networkURL := fmt.Sprintf("%s/%s", testURLRoot, networkID)

	deleteHandler := handlers.GetPartialDeleteNetworkHandler(networkURL, "test", networkSerdes)
	deleteTestConfig := tests.Test{
		Method:         "DELETE",
		URL:            networkURL,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID},
		Handler:        deleteHandler.HandlerFunc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, deleteTestConfig)

	_, err := configurator.LoadNetworkConfig(networkID, "test", networkSerdes)
	assert.EqualError(t, err, errors.ErrNotFound.Error())
}

func (m *ID) Validate(_ strfmt.Registry) error {
	if m == nil {
		return fmt.Errorf("Cannot be nil")
	}
	if m.Name != "" {
		return nil
	}
	return fmt.Errorf("Name cannot be nil")
}

func (m *ID) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

func (m *ID) GetFromNetwork(network configurator.Network) interface{} {
	feature := (&TestFeature1{}).GetFromNetwork(network)
	if feature == nil {
		return nil
	}
	return feature.(*TestFeature1).ID
}

func (m *ID) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	feature := (&TestFeature1{}).GetFromNetwork(network)
	feature.(*TestFeature1).ID = m
	return configurator.NetworkUpdateCriteria{
		ID:                   network.ID,
		ConfigsToAddOrUpdate: map[string]interface{}{"test": feature},
	}, nil
}

func (m *TestFeature1) Validate(str strfmt.Registry) error {
	if err := m.ID.Validate(str); err != nil {
		return err
	}
	return nil
}

func (m *TestFeature1) ValidateModel() error {
	return m.Validate(strfmt.Default)
}

// MarshalBinary interface implementation
func (m *TestFeature1) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary interface implementation
func (m *TestFeature1) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, m)
}

func (m *TestFeature1) GetFromNetwork(network configurator.Network) interface{} {
	if network.Configs == nil {
		return nil
	}
	config, exists := network.Configs["test"]
	if !exists {
		return nil
	}
	return config
}

func (m *TestFeature1) ToUpdateCriteria(network configurator.Network) (configurator.NetworkUpdateCriteria, error) {
	return configurator.NetworkUpdateCriteria{
		ID:                   network.ID,
		ConfigsToAddOrUpdate: map[string]interface{}{"test": m},
	}, nil
}
