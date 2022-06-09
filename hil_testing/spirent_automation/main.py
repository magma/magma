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
import sys
import time
from datetime import datetime
from typing import Dict, List

import config
import pytz
from pytz import timezone
from TS.hil_lib import hil_lib

LOGS_LOC = os.path.join(os.path.abspath(sys.path[0]), "logs")
sys.path.append(
    os.path.join(
        os.path.dirname(os.path.abspath(sys.path[0])), "Magma_Automations/scripts",
    ),
)
import base
import get_ports

LOG_FORMAT = "%(asctime)-15s HIL_AUTOMATION %(levelname)s %(message)s"
TIME_FORMAT = "%Y-%m-%d %H:%M:%S"

# TODO we should create separate library for CI automation, this would not change
LIBRARY_NAME = "sms/AGW Scale"

STATUS_UPDATE_INTERVAL = 10  # in second


def get_arg_parse():
    """
    get and parse all CLI option"""
    common_parser = argparse.ArgumentParser(add_help=False)
    common_parser.add_argument(
        "--log-level",
        default="WARNING",
        choices=["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
        type=str.upper,
        help="Log at specified level and higher",
    )
    parser = argparse.ArgumentParser(description="Run Hardware in Loop testing")
    subparsers = parser.add_subparsers(help="commands", dest="command")
    list_parser = subparsers.add_parser(
        "list", help="List supported test_suites", parents=[common_parser],
    )
    list_parser.add_argument(
        "only_list", help="Only list this test suites", type=str.upper,
    )
    run_parser = subparsers.add_parser(
        "run",
        parents=[common_parser],
        help="Run test suite ",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    run_parser.add_argument(
        "--credentials-file", "-f", help="Full path to credentials file. JSON format",
    )
    run_parser.add_argument(
        "test_suite",
        help="Run this group of tests",
        choices=config.MAGMA["Test_suite"],
        type=str.upper,
    )
    run_parser.add_argument("only_run", nargs="*", help="Only run this test ")

    run_parser.add_argument(
        "--gateway",
        dest="gateway",
        type=str.lower,
        default="mj_vgw",
        choices=config.MAGMA["AGW"],
        help="select gateway",
    )

    run_parser.add_argument(
        "--release",
        dest="rel",
        type=str.lower,
        default="ci",
        choices=config.MAGMA["REL"],
        help="select magma release",
    )

    run_parser.add_argument(
        "--build",
        dest="build",
        type=str.lower,
        default="LATEST",
        help="specify Magma build to use for Testing",
    )

    run_parser.add_argument(
        "--epc",
        dest="epc",
        type=str.lower,
        default="magma",
        choices=config.MAGMA["EPC"],
        help="select EPC provider",
    )

    run_parser.add_argument(
        "--no-output-text",
        dest="output_text",
        action="store_false",
        help="Whether or not to output ascii text tables",
    )

    run_parser.add_argument(
        "--output-s3",
        dest="output_s3",
        action="store_true",
        help="Whether or not to send results file to s3",
    )
    run_parser.add_argument(
        "--upgrade",
        dest="upgrade",
        action="store_true",
        help="Whether or not to upgrade SUT",
    )
    run_parser.add_argument(
        "--skip-build-check",
        dest="build_check",
        action="store_true",
        help="Whether or not to run test on same old SUT build",
    )
    run_parser.add_argument(
        "--reboot",
        dest="reboot",
        action="store_true",
        help="Whether or not to reboot SUT before running test",
    )
    return parser


def allow_positive(value: int) -> int:
    val = int(value)
    if val < 0:
        raise argparse.ArgumentTypeError("%s please enter a positive value" % val)
    return val


if __name__ == "__main__":
    """starts here get argument, find which tests and suite needs to be run
    run test, trigger all supported functions
    # TODO need to make this smaller push some to class functions
    """
    args = get_arg_parse().parse_args()
    logging.basicConfig(level=args.log_level, format=LOG_FORMAT)
    logging.log(
        getattr(logging, args.log_level, None), "Logging set to %s", args.log_level,
    )
    # record time before test start to retrive logs from SUT later
    _now = datetime.now(tz=pytz.utc)
    now = _now.strftime(TIME_FORMAT)
    now_pt = _now.astimezone(timezone(config.MAGMA["timezone"])).strftime(TIME_FORMAT)
    hil_runner = hil_lib(args, now_pt)
    if args.command == "list":
        hil_runner.hil_list()
    if args.command == "run":
        hil_runner.hil_run()
