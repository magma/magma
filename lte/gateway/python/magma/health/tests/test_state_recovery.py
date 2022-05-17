"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest
from unittest.mock import MagicMock

from lte.gateway.python.magma.health.state_recovery import StateRecoveryJob


class MockServiceResult():
    def __init__(self, num_fail_exits):
        self.num_fail_exits = num_fail_exits


class TestStateRecovery(unittest.TestCase):
    """
    Tests for the StateRecoveryJob python class
    """

    def test_get_last_service_restarts_without_services(self):
        services = []
        job = StateRecoveryJob(
            MagicMock(), MagicMock(),
            services,
            MagicMock(), MagicMock(),
            MagicMock(), MagicMock(),
        )
        last_services_restarts = job._get_last_service_restarts()
        self.assertDictEqual({}, last_services_restarts)

    def test_get_last_service_restarts_with_restarts(self):
        services = ["service1", "service2"]
        job = StateRecoveryJob(
            MagicMock(), MagicMock(),
            services,
            MagicMock(), MagicMock(),
            MagicMock(), MagicMock(),
        )
        job._get_service_status = MagicMock(return_value=MockServiceResult(5))
        last_services_restarts = job._get_last_service_restarts()
        self.assertDictEqual({services[0]: 5, services[1]: 5}, last_services_restarts)

    def test_get_last_service_restarts_without_restarts(self):
        services = ["service1"]
        job = StateRecoveryJob(
            MagicMock(), MagicMock(),
            services,
            MagicMock(), MagicMock(),
            MagicMock(), MagicMock(),
        )
        job._get_service_status = MagicMock(return_value=None)
        last_services_restarts = job._get_last_service_restarts()
        self.assertDictEqual({services[0]: 0}, last_services_restarts)
