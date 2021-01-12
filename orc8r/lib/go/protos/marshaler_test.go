/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protos_test

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/protos/mconfig"
)

var marshaledSrc = `{
 "configsByKey": {
  "aaa_server": {
   "@type": "type.googleapis.com/magma.mconfig.AAAConfig",
   "IdleSessionTimeoutMs": 21600000,
   "logLevel": "INFO"
  },
  "control_proxy": {
   "@type": "type.googleapis.com/magma.mconfig.ControlProxy",
   "logLevel": "INFO"
  },
  "directoryd": {
   "@type": "type.googleapis.com/magma.mconfig.DirectoryD",
   "logLevel": "INFO"
  },
  "dnsd": {
   "@type": "type.googleapis.com/magma.mconfig.DnsD",
   "dhcpServerEnabled": true,
   "enableCaching": false,
   "localTTL": 0,
   "logLevel": "INFO",
   "records": []
  },
  "eap_aka": {
   "@type": "type.googleapis.com/magma.mconfig.EapAkaConfig",
   "logLevel": "INFO",
   "timeout": {
    "ChallengeMs": 20000,
    "ErrorNotificationMs": 10000,
    "SessionMs": 43200000,
    "SessionAuthenticatedMs": 5000
   }
  },
  "magmad": {
   "@type": "type.googleapis.com/magma.mconfig.MagmaD",
   "logLevel": "INFO",
   "checkinInterval": 60,
   "checkinTimeout": 10,
   "autoupgradeEnabled": false,
   "autoupgradePollInterval": 3000,
   "packageVersion": "1.0.2-1580416046-abae140c",
   "images": [
    {
     "name": "string",
     "order": "0"
    }
   ],
   "tierId": "",
   "featureFlags": {
   },
   "dynamicServices": [
   ]
  },
  "metricsd": {
   "@type": "type.googleapis.com/magma.mconfig.MetricsD",
   "logLevel": "INFO"
  },
  "pipelined": {
   "@type": "type.googleapis.com/magma.mconfig.PipelineD",
   "defaultRuleId": "default_rule_1",
   "logLevel": "INFO",
   "relayEnabled": true,
   "services": [
    "ENFORCEMENT"
   ],
   "ueIpBlock": "192.168.128.0/24"
  },
  "redirectd": {
   "@type": "type.googleapis.com/magma.mconfig.RedirectD",
   "logLevel": "INFO"
  },
  "sessiond": {
   "@type": "type.googleapis.com/magma.mconfig.SessionD",
   "logLevel": "INFO",
   "relayEnabled": true
  },
  "td-agent-bit": {
   "@type": "type.googleapis.com/magma.mconfig.FluentBit",
   "extraTags": {
    "gateway_id": "mwc_cwf_test",
    "network_id": "mwc_cwf_net"
   },
   "throttleRate": 1000,
   "throttleWindow": 5,
   "throttleInterval": "1m",
   "filesByTag": {
   }
  }
 },
 "metadata": {
  "createdAt": "1587114828",
  "digest": {
   "md5HexDigest": "894065bf04e6e3423976eb32db15706a"
  }
 }
}`

var pipelinedExpectedVal = `{
 "@type": "type.googleapis.com/magma.mconfig.PipelineD",
 "defaultRuleId": "default_rule_1",
 "logLevel": "INFO",
 "relayEnabled": true,
 "services": [
  "ENFORCEMENT"
 ],
 "ueIpBlock": "192.168.128.0/24"
}`

func TestMconfigMarshal(t *testing.T) {
	cfg := &protos.GatewayConfigs{}
	err := protos.UnmarshalMconfig([]byte(marshaledSrc), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cfg.ConfigsByKey)

	mdval, ok := cfg.GetConfigsByKey()["magmad"]
	assert.True(t, ok)
	assert.NotNil(t, mdval)

	mdcfg := new(mconfig.MagmaD)
	err = ptypes.UnmarshalAny(mdval, mdcfg)
	assert.NoError(t, err)

	assert.Equal(t, int32(3000), mdcfg.AutoupgradePollInterval)
	assert.Equal(t, "1.0.2-1580416046-abae140c", mdcfg.PackageVersion)

	// There is no Pipelined message registered in orc9r/lib/go/proto module, check if we can marshal it as JSON
	assert.NotNil(t, cfg.ConfigsByKey["pipelined"])
	pdAny := cfg.GetConfigsByKey()["pipelined"]
	assert.NotNil(t, pdAny)
	pdJson, err := protos.MarshalMconfigToString(pdAny)
	assert.NoError(t, err)
	assert.Equal(t, pipelinedExpectedVal, pdJson)

	marshaledRes, err := protos.MarshalMconfigToString(cfg)
	assert.NoError(t, err)
	assert.Equal(t, marshaledSrc, marshaledRes)
}
