#!/usr/bin/env python3
"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import argparse
import logging
import shlex
import time

import os
from magma.configuration.mconfig_managers import load_service_mconfig


def create_parser():
    parser = argparse.ArgumentParser(
        "Start a service (or not) depending on magma mconfig")
    parser.add_argument(
        "--service", required=True,
        help="Magma service name")
    parser.add_argument(
        "--variable", required=True,
        help="Magma variable in service. "
             "Perform a truthy test on this variable.")
    parser.add_argument(
        "--not", action="store_true",
        help="Flip boolean test. E.g. for use with 'disabled' variables. "
             "'Start this service if not disabled.'")

    group = parser.add_mutually_exclusive_group()
    group.add_argument(
        "--oneshot", action="store_true",
        help="'command' is a oneshot service.")
    group.add_argument(
        "--forking", action="store_true",
        help="'command' is a forking service.")
    parser.add_argument(
        "-P", "--forking_pid_file",
        help="If forking and this parameter is specified, "
             "then write the PID in the specified file.")
    parser.add_argument("command", type=str)
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
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s')

    mconfig = load_service_mconfig(args.service)

    var = getattr(mconfig, args.variable)

    serviceEnabled = bool(var)
    if getattr(args, "not"):  # 'not' is a reserved keyword
        serviceEnabled = not serviceEnabled

    execArgs = [args.command] + args.args

    if serviceEnabled:
        logging.info(
            "service enabled, starting: %s" %
            " ".join([shlex.quote(a) for a in execArgs]))
        os.execv(execArgs[0], execArgs)
    else:
        info = "service disabled since config %s.%s==%s %%s" % (
            args.service,
            args.variable,
            bool(var),
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
                    args.service, args.variable, bool(var)),
            ]
            os.execv(forkArgs[0], forkArgs)
        else:
            logging.info(info, "(simple)")
            while True:
                time.sleep(600)
    # must never reach here


if __name__ == "__main__":
    main()
