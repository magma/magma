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
import json

from google.protobuf import json_format
from load_tests.common import (
    PROTO_DIR,
    benchmark_grpc_request,
    make_full_request_type,
    make_output_file_path,
)
from lte.protos.policydb_pb2 import (
    DisableStaticRuleRequest,
    EnableStaticRuleRequest,
)

POLICYDB_SERVICE_NAME = 'policydb'
POLICYDB_SERVICE_RPC_PATH = 'magma.lte.PolicyDB'
POLICYDB_PORT = '0.0.0.0:50068'
PROTO_PATH = PROTO_DIR + '/policydb.proto'


def _gen_imsi(num_of_ues):
    imsi = 123000000
    imsi_list = []
    for _ in range(0, num_of_ues):
        imsi_list.append(str(imsi))
        imsi = imsi + 1
    return imsi_list


def _build_enable_static_rules_data(imsi_list: list, input_file: str):
    enable_static_rule_reqs = []

    rule_list = ["p1", "p2", "p3"]
    for index, imsi in enumerate(imsi_list):
        request = EnableStaticRuleRequest(
            imsi=imsi,
            rule_ids=rule_list,
            base_names=["bn1"],
        )
        request_dict = json_format.MessageToDict(request)
        # Dumping EnableStaticRule request into json
        enable_static_rule_reqs.append(request_dict)
    with open(input_file, 'w') as file:
        json.dump(enable_static_rule_reqs, file, separators=(',', ':'))


def _build_disable_static_rules_data(imsi_list: list, input_file: str):
    disable_static_rule_reqs = []

    rule_list = ["p1", "p2", "p3"]
    for index, imsi in enumerate(imsi_list):
        request = DisableStaticRuleRequest(
            imsi=imsi,
            rule_ids=rule_list,
            base_names=["bn1"],
        )
        request_dict = json_format.MessageToDict(request)
        # Dumping DisableStaticRule request into json
        disable_static_rule_reqs.append(request_dict)
    with open(input_file, 'w') as file:
        json.dump(disable_static_rule_reqs, file, separators=(',', ':'))


def enable_static_rules_test(args):
    input_file = 'enable_static_rules.json'
    imsi_list = _gen_imsi(args.num_of_ues)
    _build_enable_static_rules_data(imsi_list, input_file)
    request_type = 'EnableStaticRules'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            POLICYDB_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num_of_ues, address=POLICYDB_PORT,
        import_path=args.import_path,
    )


def disable_static_rules_test(args):
    input_file = 'disable_static_rules.json'
    imsi_list = _gen_imsi(args.num_of_ues)
    _build_disable_static_rules_data(imsi_list, input_file)
    request_type = 'DisableStaticRules'
    benchmark_grpc_request(
        proto_path=PROTO_PATH,
        full_request_type=make_full_request_type(
            POLICYDB_SERVICE_RPC_PATH, request_type,
        ),
        input_file=input_file,
        output_file=make_output_file_path(request_type),
        num_reqs=args.num_of_ues, address=POLICYDB_PORT,
        import_path=args.import_path,
    )


def create_parser():
    """
    Creates the argparse subparser for all args
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    parser_enable_static_rules = subparsers.add_parser(
        'enable_static_rules', help='Enable static rules for a subscriber',
    )

    parser_disable_static_rules = subparsers.add_parser(
        'disable_static_rules', help='Disable static rules for a subscriber',
    )

    for subcmd in [
        parser_enable_static_rules,
        parser_disable_static_rules,
    ]:
        subcmd.add_argument(
            '--num_of_ues', help='Number of total UEs to atach',
            type=int, default=2000,
        )
        subcmd.add_argument(
            '--import_path', default=None, help='Protobuf import path directory',
        )

    parser_enable_static_rules.set_defaults(func=enable_static_rules_test)
    parser_disable_static_rules.set_defaults(func=disable_static_rules_test)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
