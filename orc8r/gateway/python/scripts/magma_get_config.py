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
import sys


def create_parser():
    parser = argparse.ArgumentParser(
        "Get magma managed configs for the specified service. (mconfig)",
    )
    parser.add_argument(
        "-s", "--service",
        required=True,
        help="Magma service name",
    )
    parser.add_argument(
        "-v", "--variable",
        help="Config variable name. "
             "If not specified, then JSON dump all configs for this service.",
    )
    parser.add_argument(
        "-t", "--test", action="store_true",
        help="Do a truthy test on v. "
             "If True then return code is 0, otherwise return code is 2",
    )
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()

    # import after parsing command line because import is sluggish
    from magma.configuration.mconfig_managers import (
        load_service_mconfig_as_json,
    )

    # set up logging
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )

    mconfig_json = load_service_mconfig_as_json(args.service)

    # if a variable was not specified, pretty print config and exit
    if args.variable is None:
        for k, v in mconfig_json.items():
            # Keys shouldn't have spaces in them, but just in case
            # Values also shouldn't have newlines, but if they do, this will
            # print differently than if called with --variable
            print(k.replace(" ", "_"), str(v).replace("\n", r"\n"))
        sys.exit(0)

    var = mconfig_json[args.variable]

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
