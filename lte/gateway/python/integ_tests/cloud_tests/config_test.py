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
import time
import unittest

import swagger_client
from integ_tests.cloud import cloud_manager, fixtures
from integ_tests.cloud.cloud_manager import CloudManager
from integ_tests.gateway import rpc


class TestConfigUpdates(unittest.TestCase):
    """
    Test that a newly-registered gateway receives updated configurations from
    the cloud.

    This test should run last in the suite as it modifies mconfig values.
    """

    MAX_CHECKS = 12
    POLL_SEC = 10

    def setUp(self):
        self._cloud_manager = cloud_manager.CloudManager()
        # We want to start with a fresh network every time because we're
        # testing gateway registration -> config update flow
        self._cloud_manager.delete_networks([fixtures.NETWORK_ID])
        # We also want to start off with default mconfigs
        rpc.reset_gateway_mconfigs()

        self._cloud_manager.create_network(fixtures.NETWORK_ID)
        self._cloud_manager.register_gateway(
            fixtures.NETWORK_ID, fixtures.GATEWAY_ID,
            rpc.get_gateway_hw_id(),
        )

    def tearDown(self):
        self._cloud_manager.clean_up()
        rpc.reset_gateway_mconfigs()

    def test_config_update(self):
        # Update configs on cloud
        updated_gw_config = swagger_client.MagmadGatewayConfig(
            **fixtures.DEFAULT_GATEWAY_CONFIG.to_dict(),
        )
        updated_gw_config.checkin_interval = 12
        updated_gw_config.checkin_timeout = 20

        updated_gw_cellular = swagger_client.GatewayCellularConfigs(
            ran=swagger_client.GatewayRanConfigs(
                **fixtures.DEFAULT_GATEWAY_CELLULAR_CONFIG.ran.to_dict(),
            ),
            epc=swagger_client.GatewayEpcConfigs(
                **fixtures.DEFAULT_GATEWAY_CELLULAR_CONFIG.epc.to_dict(),
            ),
        )
        updated_gw_cellular.ran.pci = 261

        updated_network_dnsd = swagger_client.NetworkDnsConfig(
            enable_caching=True,
        )

        updated_network_cellular = swagger_client.NetworkCellularConfigs(
            ran=swagger_client.NetworkRanConfigs(
                **fixtures.DEFAULT_NETWORK_CELLULAR_CONFIG.ran.to_dict(),
            ),
            epc=swagger_client.NetworkEpcConfigs(
                **fixtures.DEFAULT_NETWORK_CELLULAR_CONFIG.epc.to_dict(),
            ),
        )
        updated_network_cellular.epc.mcc = '002'
        updated_network_cellular.epc.mnc = '02'
        updated_network_cellular.epc.tac = 2

        self._cloud_manager.update_network_configs(
            fixtures.NETWORK_ID,
            {
                CloudManager.NetworkConfigType.DNS: updated_network_dnsd,
                CloudManager.NetworkConfigType.CELLULAR: updated_network_cellular,
            },
        )
        self._cloud_manager.update_gateway_configs(
            fixtures.NETWORK_ID, fixtures.GATEWAY_ID,
            {
                CloudManager.GatewayConfigType.MAGMAD: updated_gw_config,
                CloudManager.GatewayConfigType.CELLULAR: updated_gw_cellular,
            },
        )

        # Expected updated mconfig values
        expected = {
            'magmad': {'checkin_interval': 12, 'checkin_timeout': 20},
            'enodebd': {'pci': 261, 'tac': 2},
            'dnsd': {'enable_caching': True},
            'mme': {'mcc': '002', 'mnc': '02'},
        }

        def verify_mconfigs(actual_mconfigs):
            for srv, actual_mconfig in actual_mconfigs.items():
                expected_mconfig = expected[srv]
                for k, expected_v in expected_mconfig.items():
                    actual = getattr(actual_mconfig, k)
                    if actual != expected_v:
                        return False
            return True

        for _ in range(self.MAX_CHECKS):
            mconfigs = rpc.get_gateway_service_mconfigs(
                ['magmad', 'enodebd', 'dnsd', 'mme'],
            )
            if not verify_mconfigs(mconfigs):
                print(
                    'mconfigs do not match expected values, '
                    'will poll again',
                )
                time.sleep(self.POLL_SEC)
            else:
                return

        self.fail('mconfigs did not match expected values within poll limit')
