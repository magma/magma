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
from typing import List

from google.protobuf import json_format
from load_tests.common import (
    benchmark_grpc_request,
    generate_subs,
    make_full_request_type,
    make_output_file_path,
)
from lte.protos.session_manager_pb2 import (
    TGPP_LTE,
    CommonSessionContext,
    LocalCreateSessionRequest,
    LocalEndSessionRequest,
    LTESessionContext,
    RatSpecificContext,
)
from lte.protos.subscriberdb_pb2 import SubscriberID

PROTO_DIR = 'lte/protos'
TEST_APN = 'magma.ipv4'
CREATE_SESSION_FILENAME = '/tmp/create_session_data.json'
END_SESSION_FILENAME = '/tmp/end_session_data.json'
SESSIOND_PORT = '0.0.0.0:50065'
IMPORT_PATH = '/home/vagrant/magma'
RESULTS_PATH = '/var/tmp'
PROTO_PATH = PROTO_DIR + '/session_manager.proto'
SERVICE_NAME = 'magma.lte.LocalSessionManager'


def _handle_create_session_benchmarking(subs: List[SubscriberID]):
    _build_create_session_data(subs)
    request_type = 'CreateSession'
    benchmark_grpc_request(
        PROTO_PATH,
        make_full_request_type(SERVICE_NAME, request_type),
        CREATE_SESSION_FILENAME,
        make_output_file_path(request_type),
        len(subs),
        SESSIOND_PORT,
    )


def _build_create_session_data(subs: List[SubscriberID]):
    reqs = []
    for sid in subs:
        # build request
        req = LocalCreateSessionRequest(
            common_context=CommonSessionContext(
                sid=sid,
                apn=TEST_APN,
                rat_type=TGPP_LTE,
            ),
            rat_specific_context=RatSpecificContext(
                lte_context=LTESessionContext(
                    bearer_id=1,
                ),
            ),
        )
        req_dict = json_format.MessageToDict(req)
        # Dumping AllocateIP request into json
        reqs.append(req_dict)
    with open(CREATE_SESSION_FILENAME, 'w') as file:
        json.dump(reqs, file, separators=(',', ':'))


def _handle_end_session_benchmarking(subs: List[SubscriberID]):
    _build_end_session_data(subs)
    request_type = 'EndSession'
    benchmark_grpc_request(
        PROTO_PATH,
        make_full_request_type(SERVICE_NAME, request_type),
        END_SESSION_FILENAME,
        make_output_file_path(request_type),
        len(subs),
        SESSIOND_PORT,
    )


def _build_end_session_data(subs: List[SubscriberID]):
    reqs = []
    for sid in subs:
        # build request
        req = LocalEndSessionRequest(
            sid=sid,
            apn=TEST_APN,
        )
        req_dict = json_format.MessageToDict(req)
        # Dumping AllocateIP request into json
        reqs.append(req_dict)
    with open(END_SESSION_FILENAME, 'w') as file:
        json.dump(reqs, file, separators=(',', ':'))


def _create_parser() -> argparse.ArgumentParser:
    """Create the argparse parser with all the arguments

    Returns:
        the created parser
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title="subcommands", dest="cmd")
    parser_create = subparsers.add_parser(
        "create",
        help="Send a CreateSessionRequest to SessionD",
    )
    parser_end = subparsers.add_parser(
        "end",
        help="Send a EndSession to SessionD",
    )

    # Add arguments
    for cmd in (parser_create, parser_end):
        cmd.add_argument("--num", default=5, help="--num")
        cmd.add_argument(
            '--service_name',
            default='magma.lte.LocalSessionManager',
            help='proto service name',
        )

    # Add function callbacks
    parser_create.set_defaults(func=parser_create)
    parser_end.set_defaults(func=parser_end)
    return parser


def main():
    """Create a parser for running SessionD loadtests"""
    parser = _create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    subs = generate_subs(int(args.num))
    if args.cmd == 'create':
        _handle_create_session_benchmarking(subs)

    elif args.cmd == 'end':
        _handle_end_session_benchmarking(subs)

    print('Done')


if __name__ == "__main__":
    main()
