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
import warnings

from integ_tests.cloud.cloud_manager import CloudManager
from integ_tests.cloud.fixtures import GATEWAY_ID, NETWORK_ID
from integ_tests.gateway.rpc import get_gateway_hw_id
from swagger_client import rest

MAX_GATEWAY_CHECKS = 12
GATEWAY_POLL_INTERVAL_SEC = 10


class TestCheckin(unittest.TestCase):
    STATUS_RESPONSE_ATTRIBUTES = (
        'checkin_time', 'hardware_id', 'version',
        'system_status',
    )

    def setUp(self):
        self._cloud_manager = CloudManager()
        # We want to start with a fresh network every run
        self._cloud_manager.delete_networks([NETWORK_ID])

        self._cloud_manager.create_network(NETWORK_ID)
        self._cloud_manager.register_gateway(
            NETWORK_ID, GATEWAY_ID,
            get_gateway_hw_id(),
        )

    def tearDown(self):
        self._cloud_manager.clean_up()

    def test_checkin(self):
        """
        Basic test case to ensure that the gateway VM checks in with cloud.
        """
        # Unfortunately there's an issue with swagger's api exceptions, where
        # it never closes the socket, so we should ignore this warning
        warnings.simplefilter("ignore")

        for _ in range(MAX_GATEWAY_CHECKS):
            try:
                response = self._cloud_manager.get_gateway_status(
                    NETWORK_ID, GATEWAY_ID,
                )
                # Verify response body
                for attr in self.STATUS_RESPONSE_ATTRIBUTES:
                    self.assertTrue(hasattr(response, attr))
                return
            except rest.ApiException as e:
                self.assertEqual(e.status, 404)
                print(
                    'Gateway did not check in, waiting for {} seconds'.format(
                        GATEWAY_POLL_INTERVAL_SEC,
                    ),
                )
                time.sleep(GATEWAY_POLL_INTERVAL_SEC)

        self.assertTrue(False)


if __name__ == "__main__":
    unittest.main()
