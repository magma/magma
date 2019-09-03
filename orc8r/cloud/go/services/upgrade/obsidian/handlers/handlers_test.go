/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/stretchr/testify/assert"
)

func TestReleaseChannels(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	configuratorTestInit.StartTestService(t)
	restPort := tests.StartObsidian(t)
	testUrlRoot := fmt.Sprintf("http://localhost:%d%s/channels", restPort, obsidian.RestRoot)

	// List channels when none exist
	listChannelsTestCase := tests.Testcase{
		Name:     "List Release Channels",
		Method:   "GET",
		Url:      testUrlRoot,
		Payload:  "",
		Expected: "null",
	}
	tests.RunTest(t, listChannelsTestCase)

	// Create 2 release channels
	createChannelTestCase := tests.Testcase{
		Name:                      "Create Release Channel",
		Method:                    "POST",
		Url:                       testUrlRoot,
		Payload:                   `{"name": "stable", "supported_versions": ["1.0.0-0", "1.1.0-0"]}`,
		Skip_payload_verification: true,
	}
	_, channelId, _ := tests.RunTest(t, createChannelTestCase)
	json.Unmarshal([]byte(channelId), &channelId)
	assert.Equal(t, "stable", channelId)

	createChannelTestCase.Payload = `{
		"name": "beta",
		"supported_versions": ["1.2.0-0"]
	}`
	_, channelId, _ = tests.RunTest(t, createChannelTestCase)
	json.Unmarshal([]byte(channelId), &channelId)
	assert.Equal(t, "beta", channelId)

	listChannelsTestCase.Expected = `["beta", "stable"]`
	tests.RunTest(t, listChannelsTestCase)

	// Get release channel
	getChannelTestCase := tests.Testcase{
		Name:     "Get Release Channel",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "stable"),
		Payload:  "",
		Expected: `{"id":"","name":"stable","supported_versions":["1.0.0-0","1.1.0-0"]}`,
	}
	tests.RunTest(t, getChannelTestCase)

	// Update release channel
	updateChannelTestCase := tests.Testcase{
		Name:     "Update Release Channel",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "stable"),
		Payload:  `{"name": "stable", "supported_versions": ["1.3.0-0"]}`,
		Expected: "",
	}
	tests.RunTest(t, updateChannelTestCase)
	getChannelTestCase.Expected = `{"id":"","name":"stable","supported_versions":["1.3.0-0"]}`
	tests.RunTest(t, getChannelTestCase)

	// Delete release channel
	deleteChannelTestCase := tests.Testcase{
		Name:     "Delete Release Channel",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "beta"),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, deleteChannelTestCase)
	listChannelsTestCase.Expected = `["stable"]`
	tests.RunTest(t, listChannelsTestCase)

	// Some error cases

	// Get nonexistent release channel should 404
	status, _, err := tests.SendHttpRequest(
		"GET",
		fmt.Sprintf("%s/%s", testUrlRoot, "beta"),
		"")
	assert.NoError(t, err)
	assert.Equal(t, 404, status)

	// Update name of release channel should 400
	status, _, err = tests.SendHttpRequest(
		"PUT",
		fmt.Sprintf("%s/%s", testUrlRoot, "stable"),
		`{"name": "prod2", "supported_versions": ["1.4.0-0"]}`)
	assert.NoError(t, err)
	assert.Equal(t, 400, status)
	tests.RunTest(t, getChannelTestCase)

	// Delete nonexistent release channel should 500
	status, _, err = tests.SendHttpRequest(
		"DELETE",
		fmt.Sprintf("%s/%s", testUrlRoot, "beta"),
		"")
	assert.NoError(t, err)
	assert.Equal(t, 500, status)
}

// Obsidian integration test for tiers migrated API endpoints backed by configurator
func TestTiers(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	restPort := tests.StartObsidian(t)
	configuratorTestInit.StartTestService(t)
	netUrlRoot := fmt.Sprintf("http://localhost:%d%s/networks", restPort, obsidian.RestRoot)

	registerNetworkTestCase := tests.Testcase{
		Name:                      "Register Network",
		Method:                    "POST",
		Url:                       fmt.Sprintf("%s?requested_id=upgrade_obsidian_test_network", netUrlRoot),
		Payload:                   `{"name":"This Is A Test Network Name"}`,
		Skip_payload_verification: true,
	}
	_, networkId, err := tests.RunTest(t, registerNetworkTestCase)
	assert.NoError(t, err)
	json.Unmarshal([]byte(networkId), &networkId)

	testUrlRoot := fmt.Sprintf("%s/%s/tiers", netUrlRoot, networkId)

	// List tiers when none exist
	listTiersTestCase := tests.Testcase{
		Name:     "List Tiers",
		Method:   "GET",
		Url:      testUrlRoot,
		Payload:  "",
		Expected: "[]",
	}
	tests.RunTest(t, listTiersTestCase)

	// Create 2 tiers
	const tier1contentsA string = `{"gateways":null,"id":"t1","images":null,"name":"t1","version":"1.1.0-0"}`
	const tier2contentsA string = `{"gateways":null,"id":"t2","images":[{"name":"v002","order":10},{"name":"v001","order":15}],"name":"t2","version":"none"}`
	createTierTestCase := tests.Testcase{
		Name:                      "Create Tier",
		Method:                    "POST",
		Url:                       testUrlRoot,
		Payload:                   tier1contentsA,
		Skip_payload_verification: true,
	}
	tests.RunTest(t, createTierTestCase)
	createTierTestCase.Payload = tier2contentsA
	tests.RunTest(t, createTierTestCase)

	listTiersTestCase.Expected = `["t1", "t2"]`
	tests.RunTest(t, listTiersTestCase)

	// Get tier1
	getTierTestCase1 := tests.Testcase{
		Name:     "Get Tier",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "t1"),
		Payload:  "",
		Expected: `{"gateways":null,"id":"t1","images":null,"name":"t1","version":"1.1.0-0"}`,
	}
	tests.RunTest(t, getTierTestCase1)

	// Get tier2
	getTierTestCase2 := tests.Testcase{
		Name:     "Get Tier",
		Method:   "GET",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "t2"),
		Payload:  "",
		Expected: tier2contentsA,
	}
	tests.RunTest(t, getTierTestCase2)

	// Update tier1 to a new version
	updateTierTestCase := tests.Testcase{
		Name:     "Update Tier",
		Method:   "PUT",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "t1"),
		Payload:  `{"gateways":null, "id": "t1", "name": "t1v2", "version": "1.3.0-0", "images": null}`,
		Expected: "",
	}
	tests.RunTest(t, updateTierTestCase)
	getTierTestCase1.Expected = `{"gateways":null,"id":"t1","name":"t1v2","version":"1.3.0-0","images":null}`
	tests.RunTest(t, getTierTestCase1)

	// Update tier1 to have images
	const tier1contentsC string = `{"gateways":null,"id":"t1","images":[{"name":"v003","order":12},{"name":"v002","order":14}],"name":"t1v3","version":"none"}`
	updateTierTestCase.Payload = tier1contentsC
	tests.RunTest(t, updateTierTestCase)
	getTierTestCase1.Expected = tier1contentsC
	tests.RunTest(t, getTierTestCase1)

	// Delete tier
	deleteTierTestCase := tests.Testcase{
		Name:     "Delete Tier",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s", testUrlRoot, "t2"),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, deleteTierTestCase)

	listTiersTestCase.Expected = `["t1"]`
	tests.RunTest(t, listTiersTestCase)

	// Some error cases

	// Get nonexistent tier should 404
	status, _, err := tests.SendHttpRequest(
		"GET",
		fmt.Sprintf("%s/%s", testUrlRoot, "t2"),
		"")
	assert.NoError(t, err)
	assert.Equal(t, 404, status)

	// Delete nonexistent tier should 500
	status, _, err = tests.SendHttpRequest(
		"DELETE",
		fmt.Sprintf("%s/%s", testUrlRoot, "t2"),
		"")
	assert.NoError(t, err)
	assert.Equal(t, 500, status)

	// Remove network
	removeNetworkTestCase := tests.Testcase{
		Name:     "Force Remove Non Empty Network",
		Method:   "DELETE",
		Url:      fmt.Sprintf("%s/%s?mode=force", netUrlRoot, networkId),
		Payload:  "",
		Expected: "",
	}
	tests.RunTest(t, removeNetworkTestCase)
}
