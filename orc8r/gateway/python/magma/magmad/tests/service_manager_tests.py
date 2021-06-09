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

import asyncio
import os
import subprocess
import time
from unittest import TestCase
from unittest.mock import MagicMock

from magma.magmad.service_manager import ServiceManager, ServiceState

# Allow access to protected variables for unit testing
# pylint: disable=protected-access


class ServiceManagerSystemdTest(TestCase):
    """
    Tests for the service manager class using the systemd init system, and the
    service control sub-class. Can be run as integration tests by setting
    environment variable MAGMA_INTEGRATION_TEST=1.
    Must copy tests/magma@dummy_service.service to /etc/systemd/system/ for
    integration tests to pass.
    """

    def setUp(self):
        """
        Run before each test
        """
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)
        #  Only patch if this is a unit test. Otherwise create dummy mock so
        #  code that affects the mock can run.
        self.subprocess_mock = MagicMock(return_value='')
        if not os.environ.get('MAGMA_INTEGRATION_TEST'):
            subprocess.check_output = self.subprocess_mock

        #  Ensure test process is stopped
        self.dummy_service = ServiceManager.ServiceControl(
            name='dummy_service',
            init_system_spec=ServiceManager.SystemdInitSystem,
        )
        self._loop.run_until_complete(self.dummy_service.stop_process())

    def tearDown(self):
        """
        Run after each test
        """
        #  Ensure test process is stopped
        self._loop.run_until_complete(self.dummy_service.stop_process())
        self._loop.close()

    def test_process_start_stop(self):
        """
        Test that process can be started and stopped
        """
        self.subprocess_mock.return_value = b'inactive\n'
        self.assertEqual(self.dummy_service.status(), ServiceState.Inactive)

        self.subprocess_mock.return_value = b'active\n'
        self._loop.run_until_complete(self.dummy_service.start_process())
        time.sleep(1)  # Make sure that process doesnt immediately die
        self.assertEqual(self.dummy_service.status(), ServiceState.Active)

        self.subprocess_mock.return_value = b'inactive\n'
        self._loop.run_until_complete(self.dummy_service.stop_process())
        self.assertEqual(self.dummy_service.status(), ServiceState.Inactive)

    def test_service_manager_start_stop(self):
        """
        This test exercises the ServiceManager functions, but doesn't really
        verify functionality beyond the tests above
        """
        mgr = ServiceManager(
            ['dummy1', 'dummy2'], init_system='systemd',
            service_poller=MagicMock(),
        )

        self._loop.run_until_complete(mgr.start_services())
        self.subprocess_mock.return_value = b'active\n'
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Active,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Active,
        )

        self._loop.run_until_complete(mgr.stop_services())
        self.subprocess_mock.return_value = b'inactive\n'
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Inactive,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Inactive,
        )

    def test_dynamic_service_manager_start_stop(self):
        """
        This test exercises the ServiceManager functions, but doesn't really
        verify functionality beyond the tests above
        """
        mgr = ServiceManager(
            ['dummy1', 'dummy2'], 'systemd', MagicMock(), ['redirectd'], [],
        )

        self.subprocess_mock.return_value = b'inactive\n'
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Inactive,
        )

        self._loop.run_until_complete(
            mgr.update_dynamic_services(['redirectd']),
        )
        self.subprocess_mock.return_value = b'active\n'
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Active,
        )

        self._loop.run_until_complete(mgr.update_dynamic_services([]))
        self.subprocess_mock.return_value = b'inactive\n'
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Inactive,
        )


class ServiceManagerRunitTest(TestCase):
    """
    Tests for the service manager class using the runit init system, and the
    service control sub-class.
    """

    def setUp(self):
        """
        Run before each test
        """
        self.is_integration_test = bool(
            os.environ.get('MAGMA_INTEGRATION_TEST'),
        )
        #  Only patch if this is a unit test. Otherwise create dummy mock so
        #  code that affects the mock can run.
        self.subprocess_mock = MagicMock(return_value='')
        if not self.is_integration_test:
            subprocess.check_output = self.subprocess_mock

        if not self.is_integration_test:
            self.dummy_service = ServiceManager.ServiceControl(
                name='dummy_service',
                init_system_spec=ServiceManager.RunitInitSystem,
            )
            self.dummy_service.stop_process()

    def tearDown(self):
        """
        Run after each test
        """
        #  Ensure test process is stopped
        if not self.is_integration_test:
            self.dummy_service.stop_process()

    def test_process_start_stop(self):
        """
        Test that runit process can be started and stopped
        """
        # Skip for integration tests, otherwise will fail on systems without
        # runit
        if self.is_integration_test:
            return

        self.subprocess_mock.return_value = (
            'down: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            self.dummy_service.status(),
            ServiceState.Inactive,
        )

        self.subprocess_mock.return_value = (
            'run: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.dummy_service.start_process()
        time.sleep(1)  # Make sure that process doesnt immediately die
        self.assertEqual(
            self.dummy_service.status(),
            ServiceState.Active,
        )

        self.subprocess_mock.return_value = (
            "fail: blabla: can't change to service directory: "
            "No such file or directory\n"
        )
        self.assertEqual(
            self.dummy_service.status(),
            ServiceState.Failed,
        )

    def test_service_manager_start_stop(self):
        """
        This test exercises the ServiceManager functions, but doesn't really
        verify functionality beyond the tests above
        """
        # Skip for integration tests, otherwise will fail on systems without
        # runit
        if self.is_integration_test:
            return

        mgr = ServiceManager(
            ['dummy1', 'dummy2'], init_system='runit',
            service_poller=MagicMock(),
        )

        mgr.start_services()
        self.subprocess_mock.return_value = (
            'run: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Active,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Active,
        )

        mgr.stop_services()
        self.subprocess_mock.return_value = (
            'down: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Inactive,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Inactive,
        )

    def test_dynamic_service_manager_start_stop(self):
        """
        This test exercises the ServiceManager functions, but doesn't really
        verify functionality beyond the tests above
        """
        # Skip for integration tests, otherwise will fail on systems without
        # runit
        if self.is_integration_test:
            return

        mgr = ServiceManager(
            ['dummy1', 'dummy2'], 'runit', MagicMock(), ['redirectd'], [],
        )

        self.subprocess_mock.return_value = (
            'down: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Inactive,
        )

        mgr.update_dynamic_services(['redirectd'])
        self.subprocess_mock.return_value = (
            'run: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Active,
        )

        mgr.update_dynamic_services([])
        self.subprocess_mock.return_value = (
            'down: e2e_controller: 268195s; run: log: (pid 2275) 268195s\n'
        )
        self.assertEqual(
            mgr._service_control['redirectd'].status(),
            ServiceState.Inactive,
        )


class ServiceManagerDockerTest(TestCase):
    """
    Tests for the service manager class using the docker init system, and the
    service control sub-class. Can be run as integration tests by setting
    environment variable MAGMA_INTEGRATION_TEST=1.
    Must run docker-compose up -d in this directory for integration tests
    to pass.
    """

    def setUp(self):
        """
        Run before each test
        """
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)
        #  Only patch if this is a unit test. Otherwise create dummy mock so
        #  code that affects the mock can run.
        self.subprocess_mock = MagicMock(return_value='')
        if not os.environ.get('MAGMA_INTEGRATION_TEST'):
            subprocess.check_output = self.subprocess_mock

        #  Ensure test process is stopped
        self.dummy_service = ServiceManager.ServiceControl(
            name='dummy_service',
            init_system_spec=ServiceManager.DockerInitSystem,
        )
        try:
            self._loop.run_until_complete(self.dummy_service.stop_process())
        except FileNotFoundError:
            pass

    def tearDown(self):
        """
        Run after each test
        """
        #  Ensure test process is stopped
        try:
            self._loop.run_until_complete(self.dummy_service.stop_process())
        except FileNotFoundError:
            pass
        self._loop.close()

    def test_process_start_stop(self):
        """
        Test that process can be started and stopped
        """
        self.subprocess_mock.return_value = b'[{"State": {"Status": "exited"}}]'
        self.assertEqual(self.dummy_service.status(), ServiceState.Inactive)

        self.subprocess_mock.return_value = b'[{"State": {"Status": "running"}}]'
        try:
            self._loop.run_until_complete(self.dummy_service.start_process())
        except FileNotFoundError:
            pass
        time.sleep(1)  # Make sure that process doesnt immediately die
        self.assertEqual(self.dummy_service.status(), ServiceState.Active)

        self.subprocess_mock.return_value = b'[{"State": {"Status": "exited"}}]'
        try:
            self._loop.run_until_complete(self.dummy_service.stop_process())
        except FileNotFoundError:
            pass
        self.assertEqual(self.dummy_service.status(), ServiceState.Inactive)

    def test_service_manager_start_stop(self):
        """
        This test exercises the ServiceManager functions, but doesn't really
        verify functionality beyond the tests above
        """
        mgr = ServiceManager(
            ['dummy1', 'dummy2'], init_system='docker',
            service_poller=MagicMock(),
        )

        try:
            self._loop.run_until_complete(mgr.start_services())
        except FileNotFoundError:
            pass
        self.subprocess_mock.return_value = b'[{"State": {"Status": "running"}}]'
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Active,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Active,
        )

        try:
            self._loop.run_until_complete(mgr.stop_services())
        except FileNotFoundError:
            pass
        self.subprocess_mock.return_value = b'[{"State": {"Status": "exited"}}]'
        self.assertEqual(
            mgr._service_control['dummy1'].status(),
            ServiceState.Inactive,
        )
        self.assertEqual(
            mgr._service_control['dummy2'].status(),
            ServiceState.Inactive,
        )
