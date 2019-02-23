#!/usr/bin/env python3
"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import argparse
import logging
import sys


def create_parser():
    parser = argparse.ArgumentParser(
        "Get magma managed configs for the specified service. (mconfig)")
    parser.add_argument(
        "-s", "--service",
        required=True,
        help="Magma service name")
    parser.add_argument(
        "-v", "--variable",
        help="Config variable name. "
             "If not specified, then JSON dump all configs for this service.")
    parser.add_argument(
        "-t", "--test", action="store_true",
        help="Do a truthy test on v. "
             "If True then return code is 0, otherwise return code is 2")
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()

    # import after parsing command line because import is sluggish
    from magma.configuration.mconfig_managers import load_service_mconfig

    # set up logging
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s')

    mconfig = load_service_mconfig(args.service)

    # if a variable was not specified, pretty print config and exit
    if args.variable is None:
        print(mconfig, end="")
        sys.exit(0)

    var = getattr(mconfig, args.variable)

    if args.test:
        if var:
            # if true, then return 0 (zero means success)
            sys.exit(0)
        # exit code 2 to distinguish from exit code 1,
        #    which is returned after python exceptions.
        sys.exit(2)

    # not a boolean, print the config
    print(var)
    sys.exit(0)


if __name__ == "__main__":
    main()
