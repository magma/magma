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

from integ_tests.common.service303_utils import GatewayServicesUtil
from integ_tests.s1aptests import s1ap_wrapper


class TestServiceInfo(unittest.TestCase):

    MAX_ITER = 15

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def _check_timing_metrics(self, service_util):
        start_time = service_util.get_start_time()
        uptime = service_util.get_uptime()
        print("START TIME ", start_time)
        print("UPTIME ", uptime)

        # Sanity checks
        self.assertGreater(start_time, 0)

        self.assertGreaterEqual(uptime, 0)

        # Make sure uptime is updating and start time is not
        time.sleep(1)

        new_uptime = service_util.get_uptime()
        new_start_time = service_util.get_start_time()
        self.assertGreater(new_uptime, uptime)
        self.assertEqual(start_time, new_start_time)
        print("START TIME ", new_start_time)
        print("UPTIME ", new_uptime)

    def test_service_info(self):
        """
        Test oai service303 service info for service state and application
        health.
        """
        mme = GatewayServicesUtil().get_mme_service_util()

        print("************* Checking timing metrics for", mme._service_name)
        self._check_timing_metrics(mme)
        print("************* Service Info OK for", mme._service_name)


if __name__ == "__main__":
    unittest.main()
