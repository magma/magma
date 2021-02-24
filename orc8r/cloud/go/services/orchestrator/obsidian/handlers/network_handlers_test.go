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

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func Test_GetNetworkHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	obsidianHandlers := handlers.GetObsidianHandlers()
	listNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.GET).HandlerFunc
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.GET).HandlerFunc

	// Test empty case
	tc := tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, tc)

	// Test 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "no_such_network"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"no_such_network"},
		Handler:        getNetwork,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// register a network
	networkName1 := "network1"
	network1 := configurator.Network{
		ID:   "n1",
		Name: networkName1,
	}
	err := configurator.CreateNetwork(network1, serdes.Network)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:  "GET",
		URL:     testURLRoot,
		Payload: nil,

		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1"}),
	}
	tests.RunUnitTest(t, e, tc)

	expectedNetwork1 := models.Network{
		ID:   models1.NetworkID("n1"),
		Name: models1.NetworkName(networkName1),
	}

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// add network features
	networkFeatures1 := models.NewDefaultFeaturesConfig()
	update1 := configurator.NetworkUpdateCriteria{
		ID:                   "n1",
		ConfigsToAddOrUpdate: map[string]interface{}{orc8r.NetworkFeaturesConfig: networkFeatures1},
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update1}, serdes.Network)
	assert.NoError(t, err)

	expectedNetwork1 = models.Network{
		ID:       models1.NetworkID("n1"),
		Name:     models1.NetworkName(networkName1),
		Features: networkFeatures1,
	}

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// add dnsd configs and a description
	dnsdConfig := models.NewDefaultDNSConfig()
	description1 := "A Network"
	update1 = configurator.NetworkUpdateCriteria{
		ID:                   "n1",
		NewDescription:       &description1,
		ConfigsToAddOrUpdate: map[string]interface{}{orc8r.DnsdNetworkType: dnsdConfig},
	}
	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{update1}, serdes.Network)
	assert.NoError(t, err)

	expectedNetwork1 = models.Network{
		ID:          models1.NetworkID("n1"),
		Name:        models1.NetworkName(networkName1),
		Description: models1.NetworkDescription("A Network"),
		Features:    networkFeatures1,
		DNS:         dnsdConfig,
	}

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(expectedNetwork1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// register a second network
	networkID2 := "test_network2"
	networkName2 := "network2"
	network2 := configurator.Network{
		ID:   networkID2,
		Name: networkName2,
	}
	err = configurator.CreateNetwork(network2, serdes.Network)
	assert.NoError(t, err)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", networkID2}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_PostNetworkHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	obsidianHandlers := handlers.GetObsidianHandlers()
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.POST).HandlerFunc
	listNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.GET).HandlerFunc
	getNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.GET).HandlerFunc

	// test empty name, description
	network1 := models.NewDefaultNetwork("n1", "", "")
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"description in body should be at least 1 chars long\n" +
			"name in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	// test bad networkID format
	network1 = models.NewDefaultNetwork("Network*1", "name", "desc")
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"id in body should match '^[\\da-z_-]+$'",
	}
	tests.RunUnitTest(t, e, tc)

	// test no DNSConfig
	network1 = models.NewDefaultNetwork("n1", "name", "desc")
	network1.DNS = nil
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndns in body is required",
	}
	tests.RunUnitTest(t, e, tc)

	// test bad DNSConfig - domain
	network1 = models.NewDefaultNetwork("n1", "name", "desc")
	network1.DNS.Records[0].Domain = ""
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"validation failure list:\n" +
			"domain in body is required",
	}
	tests.RunUnitTest(t, e, tc)

	// test bad DNSConfig - ARecord
	network1 = models.NewDefaultNetwork("n1", "name", "desc")
	network1.DNS.Records[0].ARecord[0] = "not ipv4"
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"validation failure list:\n" +
			"a_record.0 in body must be of type ipv4: \"not ipv4\"",
	}
	tests.RunUnitTest(t, e, tc)

	// test bad DNSConfig - AaaaRecord
	network1 = models.NewDefaultNetwork("n1", "name", "desc")
	network1.DNS.Records[0].AaaaRecord[0] = "not ipv6"
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"validation failure list:\n" +
			"aaaa_record.0 in body must be of type ipv6: \"not ipv6\"",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	network1 = models.NewDefaultNetwork("n1", "name", "desc")
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("n1"),
	}
	tests.RunUnitTest(t, e, tc)

	actualNetwork1, err := configurator.LoadNetwork("n1", true, true, serdes.Network)
	assert.NoError(t, err)
	expectedNetwork1 := configurator.Network{
		ID:          string(network1.ID),
		Type:        string(network1.Type),
		Name:        string(network1.Name),
		Description: string(network1.Description),
		Configs: map[string]interface{}{
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
		},
	}
	assert.Equal(t, expectedNetwork1, actualNetwork1)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1"}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_DeleteNetworkHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	obsidianHandlers := handlers.GetObsidianHandlers()
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.POST).HandlerFunc
	listNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.GET).HandlerFunc
	deleteNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.DELETE).HandlerFunc

	// add a couple of networks
	network1 := models.NewDefaultNetwork("n1", "name", "desc")
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("n1"),
	}
	tests.RunUnitTest(t, e, tc)

	networkID2 := "test_network2"
	network2 := models.NewDefaultNetwork(networkID2, "name", "desc")
	tc = tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network2),
		Handler:        createNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler(networkID2),
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{"n1", networkID2}),
	}
	tests.RunUnitTest(t, e, tc)

	// delete and get
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteNetwork,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{networkID2}),
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, networkID2),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{networkID2},
		Handler:        deleteNetwork,
		ExpectedStatus: 204,
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            testURLRoot,
		Payload:        nil,
		Handler:        listNetwork,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler([]string{}),
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_PutNetworkHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	obsidianHandlers := handlers.GetObsidianHandlers()
	createNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks", obsidian.POST).HandlerFunc
	updateNetwork := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.PUT).HandlerFunc
	getNetworkHandler := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id", obsidian.GET).HandlerFunc

	// happy path
	// add a network
	network1 := models.NewDefaultNetwork("n1", "name", "desc")
	tc := tests.Test{
		Method:         "POST",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        createNetwork,
		ExpectedStatus: 201,
		ExpectedResult: tests.JSONMarshaler("n1"),
	}
	tests.RunUnitTest(t, e, tc)

	// change meta data
	network1.Name = models1.NetworkName("name2")
	network1.Type = "wifi"
	network1.Description = models1.NetworkDescription("desc2")
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkHandler,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
	}
	tests.RunUnitTest(t, e, tc)

	// change configs
	network1.DNS.EnableCaching = swag.Bool(false)
	network1.Features.Features["new-feature"] = "foobar"
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        updateNetwork,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getNetworkHandler,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(network1),
	}
	tests.RunUnitTest(t, e, tc)

	// try do delete DNS config
	network1.DNS = nil
	tc = tests.Test{
		Method:         "PUT",
		URL:            testURLRoot,
		Payload:        tests.JSONMarshaler(network1),
		Handler:        updateNetwork,
		ExpectedStatus: 400,
		ExpectedError:  "validation failure list:\ndns in body is required",
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_GetNetworkMetadataHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	getName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/name", obsidian.GET).HandlerFunc
	getType := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/type", obsidian.GET).HandlerFunc
	getDesc := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/description", obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getName,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("network1"),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/type/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getType,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("type1"),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	getNetworkDesc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/description/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getDesc,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler("network 1"),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkDesc)
}

func Test_PutNetworkMetadataHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	// register a network
	seedNetworks(t)
	expectedNetwork1 := configurator.Network{
		ID:          "n1",
		Type:        "type1",
		Name:        "network1",
		Description: "network 1",
		Configs:     map[string]interface{}{},
	}

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateName := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/name", obsidian.PUT).HandlerFunc
	updateType := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/type", obsidian.PUT).HandlerFunc
	updateDesc := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/description", obsidian.PUT).HandlerFunc

	// check for validity
	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(""),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateName,
		ExpectedStatus: 400,
		ExpectedError:  " in body should be at least 1 chars long",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/name/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler("new_name"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateName,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	actualNetwork, err := configurator.LoadNetwork("n1", true, false, serdes.Network)
	assert.NoError(t, err)
	expectedNetwork1.Version = 1
	expectedNetwork1.Name = "new_name"
	assert.Equal(t, expectedNetwork1, actualNetwork)

	// happy case
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/type/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler("new_type"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateType,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	actualNetwork, err = configurator.LoadNetwork("n1", true, false, serdes.Network)
	assert.NoError(t, err)
	expectedNetwork1.Type = "new_type"
	expectedNetwork1.Version = 2
	assert.Equal(t, expectedNetwork1, actualNetwork)

	// happy case
	putNetworkDesc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/description/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler("new_name"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateDesc,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, putNetworkDesc)

	actualNetwork, err = configurator.LoadNetwork("n1", true, false, serdes.Network)
	assert.NoError(t, err)
	expectedNetwork1.Description = "new_name"
	expectedNetwork1.Version = 3
	assert.Equal(t, expectedNetwork1, actualNetwork)
}

func Test_GetNetworkFeaturesHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	getFeatures := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/features", obsidian.GET).HandlerFunc

	getNetworkFeatures := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/features/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getFeatures,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultFeaturesConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, getNetworkFeatures)
}

func Test_PutNetworkFeaturesHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateFeatures := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/features", obsidian.PUT).HandlerFunc

	// update full feature happy case
	newFeatures := &models.NetworkFeatures{
		Features: map[string]string{
			"hello": "world!!",
		},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/features/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newFeatures),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateFeatures,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	config, err := configurator.LoadNetworkConfig("n1", orc8r.NetworkFeaturesConfig, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, newFeatures, config)
}

func Test_GetNetworkDNSHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	getDNS := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns", obsidian.GET).HandlerFunc
	getDNSRecords := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records", obsidian.GET).HandlerFunc
	getDNSRecordByDomain := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records/:domain", obsidian.GET).HandlerFunc

	tc := tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getDNS,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultDNSConfig()),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/dns/records", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        getDNSRecords,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultDNSConfig().Records),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/dns/records/example.com", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "example.com"},
		Handler:        getDNSRecordByDomain,
		ExpectedStatus: 200,
		ExpectedResult: tests.JSONMarshaler(models.NewDefaultDNSConfig().Records[0]),
		ExpectedError:  "",
	}
	tests.RunUnitTest(t, e, tc)

	// 404
	tc = tests.Test{
		Method:         "GET",
		URL:            fmt.Sprintf("%s/%s/dns/records/google.com", testURLRoot, "n1"),
		Payload:        nil,
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "google.com"},
		Handler:        getDNSRecordByDomain,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_PutNetworkDNSHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	updateDNS := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns", obsidian.PUT).HandlerFunc
	updateDNSRecords := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records", obsidian.PUT).HandlerFunc
	updateDNSRecordByDomain := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records/:domain", obsidian.PUT).HandlerFunc

	// update full dns validation failure
	newDNS := models.NewDefaultDNSConfig()
	newDNS.Records = []*models.DNSConfigRecord{
		{
			ARecord:     []strfmt.IPv4{"192-88-99-142"},
			AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			CnameRecord: []string{"facebook.com"},
			Domain:      "facebook.com",
		},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newDNS),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateDNS,
		ExpectedStatus: 400,
		ExpectedError: "validation failure list:\n" +
			"validation failure list:\n" +
			"a_record.0 in body must be of type ipv4: \"192-88-99-142\"",
	}
	tests.RunUnitTest(t, e, tc)

	// update full DNS happy case
	newDNS = models.NewDefaultDNSConfig()
	newDNS.Records = []*models.DNSConfigRecord{
		{
			ARecord:     []strfmt.IPv4{"192.88.99.142"},
			AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			CnameRecord: []string{"facebook.com"},
			Domain:      "facebook.com",
		},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(newDNS),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateDNS,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	config, err := configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, newDNS, config)

	// update the records only
	records := []*models.DNSConfigRecord{
		{
			ARecord:     []strfmt.IPv4{"192.88.99.142"},
			AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			CnameRecord: []string{"yahoo.com"},
			Domain:      "yahoo.com",
		},
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/records/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(records),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateDNSRecords,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	config, err = configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, models.NetworkDNSRecords(records), config.(*models.NetworkDNSConfig).Records)

	// updating a nonexistent record should fail
	record := &models.DNSConfigRecord{
		ARecord:     []strfmt.IPv4{"192.88.99.142"},
		AaaaRecord:  []strfmt.IPv6{"1234:0db8:85a3:0000:0000:8a2e:0370:1234"},
		CnameRecord: []string{"google.com"},
		Domain:      "google.com",
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/records/google.com/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(record),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "google.com"},
		Handler:        updateDNSRecordByDomain,
		ExpectedStatus: 404,
		ExpectedError:  "Not Found",
	}
	tests.RunUnitTest(t, e, tc)

	// happy case update yahoo.com record
	record = &models.DNSConfigRecord{
		ARecord:     []strfmt.IPv4{"192.88.99.142"},
		AaaaRecord:  []strfmt.IPv6{"1234:0db8:85a3:0000:0000:8a2e:0370:1234"},
		CnameRecord: []string{"yahoo.com"},
		Domain:      "yahoo.com",
	}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/records/yahoo.com/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(record),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "yahoo.com"},
		Handler:        updateDNSRecordByDomain,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	config, err = configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.NoError(t, err)
	assert.Equal(t, record, config.(*models.NetworkDNSConfig).Records[0])

	// delete all records
	records = []*models.DNSConfigRecord{}
	tc = tests.Test{
		Method:         "PUT",
		URL:            fmt.Sprintf("%s/%s/dns/records/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(records),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        updateDNSRecords,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)
	config, err = configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.NoError(t, err)
	assert.Empty(t, config.(*models.NetworkDNSConfig).Records)
}

func Test_CreateNetworkDNSRecord(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	seedNetworks(t)

	obsidianHandlers := handlers.GetObsidianHandlers()
	postDNSRecord := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records/:domain", obsidian.POST).HandlerFunc

	// validation failure
	record := &models.DNSConfigRecord{
		ARecord:     []strfmt.IPv4{"192.88.99.142"},
		AaaaRecord:  []strfmt.IPv6{"a2e:0370:1234"},
		CnameRecord: []string{"yahoo.com"},
		Domain:      "yahoo.com",
	}
	tc := tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/dns/records/yahoo.com/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(record),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "yahoo.com"},
		Handler:        postDNSRecord,
		ExpectedError:  "validation failure list:\naaaa_record.0 in body must be of type ipv6: \"a2e:0370:1234\"",
		ExpectedStatus: 400,
	}
	tests.RunUnitTest(t, e, tc)

	// cannot register a record with an existing domain
	record = models.NewDefaultDNSConfig().Records[0]
	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/dns/records/example.com/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(record),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "example.com"},
		Handler:        postDNSRecord,
		ExpectedError:  "A record with domain:example.com already exists",
		ExpectedStatus: 400,
	}
	tests.RunUnitTest(t, e, tc)

	// happy case
	record = &models.DNSConfigRecord{
		ARecord:     []strfmt.IPv4{"192.88.99.142"},
		AaaaRecord:  []strfmt.IPv6{"1234:0db8:85a3:0000:0000:8a2e:0370:1234"},
		CnameRecord: []string{"google.com"},
		Domain:      "google.com",
	}
	tc = tests.Test{
		Method:         "POST",
		URL:            fmt.Sprintf("%s/%s/dns/records/google.com/", testURLRoot, "n1"),
		Payload:        tests.JSONMarshaler(record),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "google.com"},
		Handler:        postDNSRecord,
		ExpectedStatus: 201,
	}
	tests.RunUnitTest(t, e, tc)
}

func Test_DeleteNetworkDNSHandlers(t *testing.T) {
	test_init.StartTestService(t)

	e := echo.New()
	testURLRoot := "/magma/v1/networks"

	obsidianHandlers := handlers.GetObsidianHandlers()
	deleteDNS := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns", obsidian.DELETE).HandlerFunc
	deleteDNSByDomain := tests.GetHandlerByPathAndMethod(t, obsidianHandlers, "/magma/v1/networks/:network_id/dns/records/:domain", obsidian.DELETE).HandlerFunc

	seedNetworks(t)

	tc := tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/%s/dns/records/%s", testURLRoot, "n1", "example.com"),
		ParamNames:     []string{"network_id", "domain"},
		ParamValues:    []string{"n1", "example.com"},
		Handler:        deleteDNSByDomain,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	config, err := configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.NoError(t, err)
	dnsConfig := config.(*models.NetworkDNSConfig)
	assert.Empty(t, dnsConfig.Records)

	tc = tests.Test{
		Method:         "DELETE",
		URL:            fmt.Sprintf("%s/%s/dns/", testURLRoot, "n1"),
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        deleteDNS,
		ExpectedStatus: 204,
	}
	tests.RunUnitTest(t, e, tc)

	_, err = configurator.LoadNetworkConfig("n1", orc8r.DnsdNetworkType, serdes.Network)
	assert.EqualError(t, err, "Not found")

}

func seedNetworks(t *testing.T) {
	_, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				ID:          "n1",
				Type:        "type1",
				Name:        "network1",
				Description: "network 1",
				Configs: map[string]interface{}{
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
		},
		serdes.Network,
	)
	assert.NoError(t, err)
}
