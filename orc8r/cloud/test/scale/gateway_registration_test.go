/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package scale_test

import (
	"flag"
	oclient "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/assert"
	"magma/orc8r/cloud/go/obsidian/swagger/v1/client"
	"magma/orc8r/cloud/go/obsidian/swagger/v1/client/lte_gateways"
	"magma/orc8r/cloud/go/obsidian/swagger/v1/client/lte_networks"
	"magma/orc8r/cloud/go/obsidian/swagger/v1/client/upgrades"
	"magma/orc8r/cloud/test/testlib"
	"net/http"
	"testing"
	//"magma/orc8r/cloud/go/obsidian/swagger/v1/models"
	//"fmt"
)

func TestSanity(t *testing.T) {
	flag.Parse()

	cluster, err := testlib.GetClusterInfo()
	assert.NoError(t, err)

	err = testlib.SetClusterInfo(cluster)
	assert.NoError(t, err)
	t.Logf("Cluster Information %+v", cluster)

	tlsConfig, err := testlib.GetTLSConfig()
	assert.NoError(t, err)

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	openAPIClient := oclient.NewWithClient(
		"api.staging.testminster.com",
		client.DefaultBasePath,
		client.DefaultSchemes,
		httpClient)
	c := client.New(openAPIClient, nil)

	// Create test network
	networkID := "scale_test"
	_, err = c.LTENetworks.PostLTE(
		lte_networks.NewPostLTEParams().WithLTENetwork(
			testlib.GetDefaultLteNetwork(networkID)))
	//assert.NoError(t, err)

	// check if network was successfully created
	res, err := c.LTENetworks.GetLTE(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, []string{networkID}, res.Payload)
	assert.NoError(t, err)

	// Register tier
	_, err = c.Upgrades.PostNetworksNetworkIDTiers(upgrades.NewPostNetworksNetworkIDTiersParams().WithTier(testlib.GetDefaultTier()).WithNetworkID(networkID))
	if err != nil {
		payload := err.(*upgrades.PostNetworksNetworkIDTiersDefault).Payload
		t.Logf("Error: %s\n", *payload.Message)
	}
	// assert.NoError(t, err)

	// Register gateways and check sanity
	for _, gateway := range cluster.Gateways {
		_, err = c.LTEGateways.DeleteLTENetworkIDGatewaysGatewayID(
			lte_gateways.NewDeleteLTENetworkIDGatewaysGatewayIDParams().
				WithGatewayID(gateway.ID).
				WithNetworkID(networkID))
		if err != nil {
			// payload := err.(*lte_gateways.DeleteLTENetworkIDGatewaysGatewayIDNoContent).Payload
			// t.Logf("Error: %s\n", *payload.Message)
		}
		//assert.NoError(t, err)

		_, err = c.LTEGateways.PostLTENetworkIDGateways(
			lte_gateways.NewPostLTENetworkIDGatewaysParams().
				WithGateway(testlib.GetDefaultLteGateway(gateway.ID, gateway.HardwareID)).
				WithNetworkID(networkID))
		if err != nil {
			payload := err.(*lte_gateways.PostLTENetworkIDGatewaysDefault).Payload
			t.Logf("Error: %s\n", *payload.Message)
		}
		assert.NoError(t, err)
	}

	// wait for gateways to be up and healthy

	// tear down the gateways

	// teardown the network
	//_, err = c.LTENetworks.DeleteLTENetworkID(
	//	lte_networks.NewDeleteLTENetworkIDParams().WithNetworkID(networkID))
	//assert.NoError(t, err)
}
