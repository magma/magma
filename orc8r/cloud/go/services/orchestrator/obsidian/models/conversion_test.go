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

package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
)

func Test_Conversions(t *testing.T) {
	cNetwork := configurator.Network{
		ID:          "test",
		Name:        "name",
		Type:        "type",
		Description: "desc",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
			orc8r.NetworkSentryConfig:   models.NewDefaultSentryConfig(),
			orc8r.StateConfig:           models.NewDefaultStateConfig(),
		},
	}
	generatedSNetwork := (&models.Network{}).FromConfiguratorNetwork(cNetwork)
	sNetwork := models.Network{
		ID:           models1.NetworkID("test"),
		Name:         models1.NetworkName("name"),
		Type:         "type",
		Description:  models1.NetworkDescription("desc"),
		Features:     models.NewDefaultFeaturesConfig(),
		DNS:          models.NewDefaultDNSConfig(),
		SentryConfig: models.NewDefaultSentryConfig(),
		StateConfig:  models.NewDefaultStateConfig(),
	}
	generatedCNetwork := sNetwork.ToConfiguratorNetwork()

	assert.Equal(t, sNetwork, *generatedSNetwork)
	assert.Equal(t, cNetwork, generatedCNetwork)
}

func TestLastGatewayCheckInWasRecent(t *testing.T) {
	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 20,
	}

	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-60 * time.Second).UnixMilli())
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-101 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-10 * time.Hour).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, &magmadConfig))
}

func TestLastGatewayCheckInWasRecentHandlingNil(t *testing.T) {
	assert.False(t, models.LastGatewayCheckInWasRecent(nil, nil))

	magmadConfig := models.MagmadGatewayConfigs{
		CheckinInterval: 20,
	}
	assert.False(t, models.LastGatewayCheckInWasRecent(nil, &magmadConfig))

	gatewayStatus := models.GatewayStatus{
		CheckinTime: uint64(time.Now().UnixMilli()),
	}
	assert.True(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil))

	gatewayStatus.CheckinTime = uint64(time.Now().Add(-301 * time.Second).UnixMilli())
	assert.False(t, models.LastGatewayCheckInWasRecent(&gatewayStatus, nil))
}
