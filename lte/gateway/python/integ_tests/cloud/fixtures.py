"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
# Network ID for integration tests
import swagger_client

NETWORK_ID = 'integ_net'

# Gateway ID for integration tests
GATEWAY_ID = 'integ_gate'

DEFAULT_NETWORK_DNSD_CONFIG = swagger_client.NetworkDnsConfig(
    enable_caching=False,
    records=[],
)

DEFAULT_NETWORK_CELLULAR_CONFIG = swagger_client.NetworkCellularConfigs(
    ran=swagger_client.NetworkRanConfigs(
        earfcndl=44590,
        bandwidth_mhz=20,
        subframe_assignment=2,
        special_subframe_pattern=7,
    ),
    epc=swagger_client.NetworkEpcConfigs(
        mcc='001',
        mnc='01',
        tac=1,
        lte_auth_op='EREREREREREREREREREREQ==',
        lte_auth_amf='gAA=',
        default_rule_id='default_rule_1',
        relay_enabled=True,
        hss_relay_enabled=True,
        gx_gy_relay_enabled=True,
    ),
)

DEFAULT_GATEWAY_CONFIG = swagger_client.MagmadGatewayConfig(
    checkin_interval=10,
    checkin_timeout=15,
    autoupgrade_enabled=False,
    autoupgrade_poll_interval=300,
    tier='default',
)

DEFAULT_GATEWAY_CELLULAR_CONFIG = swagger_client.GatewayCellularConfigs(
    ran=swagger_client.GatewayRanConfigs(
        pci=260,
        transmit_enabled=True,
    ),
    epc=swagger_client.GatewayEpcConfigs(
        nat_enabled=True,
        ip_block='192.168.128.0/24',
    ),
)
