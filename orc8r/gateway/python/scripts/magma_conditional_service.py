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
import shlex
import time

from magma.configuration.mconfig_managers import load_service_mconfig_as_json


def create_parser():
    parser = argparse.ArgumentParser(
        "Start a service (or not) depending on magma mconfig",
    )
    parser.add_argument(
        "--service", required=True,
        help="Magma service name",
    )
    parser.add_argument(
        "--variable", required=True,
        help="Magma variable in service. "
             "Perform a truthy test on this variable.",
    )
    parser.add_argument(
        "--not",
        dest="invert_enable",
        action="store_true",
        help="Flip boolean test. E.g. for use with 'disabled' variables. "
             "'Start this service if not disabled.'",
    )
    parser.add_argument(
        "--enable-by-default",
        action="store_true",
        help="If the value is unset, treat as a truthy value",
    )

    group = parser.add_mutually_exclusive_group()
    group.add_argument(
        "--oneshot", action="store_true",
        help="'command' is a oneshot service.",
    )
    group.add_argument(
        "--forking", action="store_true",
        help="'command' is a forking service.",
    )
    parser.add_argument(
        "-P", "--forking_pid_file",
        help="If forking and this parameter is specified, "
             "then write the PID in the specified file.",
    )
    parser.add_argument("command")
    parser.add_argument("args", nargs=argparse.REMAINDER)
    return parser


def main():
    """
    conditionally start a systemd service based on an mconfig
    Example usage:
        ExecStart=/usr/sbin/magma_conditional_service.py \
                  --service service \
                  --variable enable \
                  service.sh start
    """

    parser = create_parser()
    args = parser.parse_args()

    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )

    mconfig = load_service_mconfig_as_json(args.service)

    service_enabled = bool(mconfig.get(args.variable, args.enable_by_default))

    if args.invert_enable:
        service_enabled = not service_enabled

    execArgs = [args.command] + args.args

    if service_enabled:
        logging.info(
            "service enabled, starting: %s" %
            " ".join([shlex.quote(a) for a in execArgs]),
        )
        os.execv(execArgs[0], execArgs)
    else:
        info = "service disabled since config %s.%s==%s %%s" % (
            args.service,
            args.variable,
            service_enabled,
        )
        if args.oneshot:
            logging.info(info, "(oneshot, exiting...)")
            return 0
        elif args.forking:
            writePIDCmd = ""
            if args.forking_pid_file:
                writePIDCmd = "( echo $! > %s )" % args.forking_pid_file
            logging.info(info, "(forking, pid_file=%s)" % args.forking_pid_file)
            # TODO: use os.fork(), when it works on all devices.
            forkArgs = [
                "/bin/sh", "-c",
                "while true; do sleep 600; done & %s "
                "# conditional_service disabled since config %s.%s==%s" % (
                    writePIDCmd,
                    args.service, args.variable, service_enabled,
                ),
            ]
            os.execv(forkArgs[0], forkArgs)
        else:
            logging.info(info, "(simple)")
            while True:
                time.sleep(600)
    # must never reach here


if __name__ == "__main__":
    main()
