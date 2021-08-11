#!/usr/bin/env python3

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

import os
from unittest import TestCase

from magma.common.health.service_state_wrapper import ServiceStateWrapper
from orc8r.protos.service_status_pb2 import ServiceExitStatus
from scripts.service_util import update_stats


# Integration test for service_util script. Requires Redis to be running.
class CertUtilsTest(TestCase):
    def setUp(self):
        self.serviceName = "TestService"
        self.serviceName2 = "TestService2"

    def tearDown(self):
        wrapper = ServiceStateWrapper()
        wrapper.cleanup_service_status()

    @staticmethod
    def _success_fixture():
        os.environ["SERVICE_RESULT"] = "success"
        os.environ["EXIT_CODE"] = "exited"
        os.environ["EXIT_STATUS"] = "0"

    @staticmethod
    def _core_dump_fixture():
        os.environ["SERVICE_RESULT"] = "core-dump"
        os.environ["EXIT_CODE"] = "dumped"
        os.environ["EXIT_STATUS"] = "ABRT"

    @staticmethod
    def _non_zero_exit_fixture():
        os.environ["SERVICE_RESULT"] = "exit-code"
        os.environ["EXIT_CODE"] = "exited"
        os.environ["EXIT_STATUS"] = "2"

    def test_success_exit(self) -> None:
        """
        Test successful exits for a service multiple times and for multiple
        services
        """
        self._success_fixture()
        update_stats(self.serviceName)
        wrapper = ServiceStateWrapper()
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 0, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("SUCCESS"),
            "Service result check",
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
            "Service exit check",
        )
        self.assertEqual(
            status.num_clean_exits, 1,
            "Clean exit check",
        )

        # Multiple restarts
        update_stats(self.serviceName)
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 0, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("SUCCESS"),
            "Service result check",
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
            "Service exit check",
        )
        self.assertEqual(
            status.num_clean_exits, 2,
            "Clean exit check",
        )

        # Multiple service restarts
        update_stats(self.serviceName2)
        status = wrapper.get_service_status(self.serviceName2)
        self.assertEqual(status.latest_rc, 0, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("SUCCESS"),
            "Service result check",
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
            "Service exit code check",
        )
        self.assertEqual(
            status.num_clean_exits, 1,
            "Clean exit check",
        )
        self.assertEqual(status.num_fail_exits, 0, "Failure exit status")

    def test_coredump_exit(self) -> None:
        """
        Test core dump exit and also recovery after a core dump
        """
        self._core_dump_fixture()
        update_stats(self.serviceName)
        wrapper = ServiceStateWrapper()
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 0, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("CORE_DUMP"),
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("DUMPED"),
        )
        self.assertEqual(status.num_fail_exits, 1)
        self.assertEqual(status.num_clean_exits, 0)

        update_stats(self.serviceName)
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.num_fail_exits, 2)
        self.assertEqual(status.num_clean_exits, 0)

        # Test that we can do a clean update after exits
        self._success_fixture()
        update_stats(self.serviceName)
        wrapper = ServiceStateWrapper()
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 0, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("SUCCESS"),
            "Service result check",
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
            "Service exit check",
        )
        self.assertEqual(status.num_clean_exits, 1, "Clean exit check")
        self.assertEqual(status.num_fail_exits, 2)

    def test_non_zero_exit(self) -> None:
        """
        Test non zero service exit multiple times
        """
        self._non_zero_exit_fixture()
        update_stats(self.serviceName)
        wrapper = ServiceStateWrapper()
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 2, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("EXIT_CODE"),
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
        )
        self.assertEqual(status.num_fail_exits, 1)
        self.assertEqual(status.num_clean_exits, 0)

        os.environ["EXIT_STATUS"] = "3"

        update_stats(self.serviceName)
        status = wrapper.get_service_status(self.serviceName)
        self.assertEqual(status.latest_rc, 3, "service exit status check")
        self.assertEqual(
            status.latest_service_result,
            ServiceExitStatus.ServiceResult.Value("EXIT_CODE"),
        )
        self.assertEqual(
            status.latest_exit_code,
            ServiceExitStatus.ExitCode.Value("EXITED"),
        )
        self.assertEqual(status.num_fail_exits, 2)
        self.assertEqual(status.num_clean_exits, 0)
