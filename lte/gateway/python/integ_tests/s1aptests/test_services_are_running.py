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

import time
import unittest
from typing import Dict, List, Optional

from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import InitMode, MagmadUtil

SERVICE_START_TIME = 'start_time'
SERVICE_ACTIVE = 'active_status'
SERVICE_RESULT = 'result'


class TestServicesAreRunning(unittest.TestCase):
    """
    A simple smoke test that checks if all relevant magma services are running.
    This test was introduced because of GH14007 where a service
    that is not used in the integration tests was damaged in the
    magma debian package artifact.
    """

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            health_service=MagmadUtil.health_service_cmds.ENABLE,
        )

    def test_services_are_running(self):

        if self._s1ap_wrapper.magmad_util.init_system.value \
                != InitMode.SYSTEMD.value:
            self.skipTest("Systemd only test.")

        services = self._s1ap_wrapper.magmad_util.get_magma_services()

        """
        Initialize service status dictionary.
        The dict holds for each service status results that are collected
        each time the status is queried. Also for each service a result
        list holds reporting about issues with the queried status results.
        """
        service_status = {
            service: {
                SERVICE_START_TIME: [],
                SERVICE_ACTIVE: [],
                SERVICE_RESULT: [],
            } for service in services
        }

        self._query_state_of_services(service_status)

        print(
            'Waiting 10 seconds for second check to identify '
            'services in a restart-loop.',
        )
        time.sleep(10)

        self._query_state_of_services(service_status)

        failed_services = self._get_failed_services(service_status)

        for service, state in failed_services.items():
            print(f'{service} failed with reason(s) "{state[SERVICE_RESULT]}"')
            print(self.get_failed_service_info(service))

        assert not failed_services, \
            "Services are not running correctly. See logging above."

    def _query_state_of_services(
        self,
        service_status: Dict[str, Dict[str, List[str]]],
    ):
        for service, status in service_status.items():
            active_state = "not active"
            if self._s1ap_wrapper.magmad_util.is_service_active(service):
                active_state = "active"
            service_name = (
                self._s1ap_wrapper.magmad_util.
                map_service_to_init_system_service_name(service)
            )
            start_time = self._s1ap_wrapper.magmad_util.exec_command_output(
                f'systemctl show {service_name} '
                f'--property=ActiveEnterTimestamp',
            ).strip()

            status[SERVICE_START_TIME].append(start_time)
            status[SERVICE_ACTIVE].append(active_state)

            print(f'checking {service} ...')
            print(f'  {active_state}')
            print(f'  {start_time}')

    def get_failed_service_info(self, failed_service: str):
        return self._s1ap_wrapper.magmad_util.exec_command_capture_output(
            f'sudo systemctl status {failed_service}',
        ).stdout.decode("utf-8", errors='ignore')

    def _get_failed_services(
        self,
        service_status: Dict[str, Dict[str, List[str]]],
    ) -> Optional[Dict[str, Dict[str, List[str]]]]:
        for service, status in service_status.items():
            # was the service not active?
            active_status_first_check = status[SERVICE_ACTIVE][0]
            active_status_second_check = status[SERVICE_ACTIVE][1]
            not_active = (
                active_status_first_check != 'active'
                or active_status_second_check != 'active'
            )
            if not_active:
                status[SERVICE_RESULT].append('not active')

            # was the service restarted?
            start_time_first_check = status[SERVICE_START_TIME][0]
            start_time_second_check = status[SERVICE_START_TIME][1]
            restarted = start_time_first_check != start_time_second_check
            if restarted:
                status[SERVICE_RESULT].append('restarted')

        return {
            service: status for service, status in service_status.items()
            if len(status[SERVICE_RESULT]) > 0
        }
