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

import argparse
import logging
import os

from magma.common.health.service_state_wrapper import ServiceStateWrapper
from orc8r.protos.service_status_pb2 import ServiceExitStatus


def get_status() -> ServiceExitStatus:
    """
    Populates the Status protobuf with the status environment as populated by
    systemd
    @returns a populated service exit status object
    """
    service_result = os.environ.get("SERVICE_RESULT")
    exit_code = os.environ.get("EXIT_CODE")
    exit_status = os.environ.get("EXIT_STATUS")

    # Populate the service exit status string and exit code.
    status_obj = ServiceExitStatus()
    status_obj.latest_rc = 0
    status_obj.latest_service_result = \
        ServiceExitStatus.ServiceResult.Value(
            service_result.upper().replace('-', '_'),
        )
    status_obj.latest_exit_code = \
        ServiceExitStatus.ExitCode.Value(exit_code.upper())

    if (
        status_obj.latest_service_result
        == ServiceExitStatus.ServiceResult.Value("EXIT_CODE")
    ):
        try:
            status_obj.latest_rc = int(exit_status)
        except ValueError:
            logging.error("Error parsing service exit status", exit_status)
            pass
    return status_obj


def get_service_name() -> str:
    parser = argparse.ArgumentParser(
        description='Systemd service exit utility script',
    )
    parser.add_argument(
        'service_name',
        help='name of the service that is exiting',
    )
    args = parser.parse_args()
    return args.service_name


def update_stats(service: str) -> None:
    status = get_status()
    wrapper_obj = ServiceStateWrapper()
    try:
        wrapper_obj.update_service_status(service, status)
    except Exception as e:
        logging.error('Failed to write to redis, status %s', e)
        logging.error('Logging exit info instead \n%s', status)


if __name__ == "__main__":
    service_name = get_service_name()
    update_stats(service_name)
