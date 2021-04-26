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
import subprocess  # noqa: S404
from typing import List

from lte.protos.subscriberdb_pb2 import SubscriberID

IMPORT_PATH = '/home/vagrant/magma'
RESULTS_PATH = '/var/tmp'


def generate_subs(num_subs: int) -> List[SubscriberID]:
    """Return a list of num_subs many SubscriberIDs

    Args:
        num_subs (int): number of SubscriberIDs to generate

    Returns:
        List[SubscriberID]: Created list of SubscriberIDs
    """
    subs = []
    digit_num = 15
    for index in range(1, num_subs):
        sid = SubscriberID(id=str(index).zfill(digit_num))
        subs.append(sid)
    return subs


def make_output_file_path(
    request_type: str,
) -> str:
    """Return the output file path for the given request type

    Args:
        request_type (str): GRPC request type

    Returns:
        str: full output file path
    """
    return '%s/result_%s.json' % (RESULTS_PATH, request_type)


def make_full_request_type(
    service_name: str,
    request_type: str,
) -> str:
    """Return the full GRPC request type by combining service name and request type

    Args:
        service_name (str): ex: magma.lte.LocalSessionManager
        request_type (str): ex: CreateSession

    Returns:
        str: full request type
    """
    return '%s/%s' % (service_name, request_type)


def benchmark_grpc_request(
    proto_path: str,
    full_request_type: str,
    input_file: str,
    output_file: str,
    num_reqs: int,
    address: str,
):
    """Run GHZ based GRPC benchmarking

    Args:
        proto_path (str): full path to the proto file with definitions
        full_request_type (str): grpc service name + request type
        input_file (str): a path to where data is placed
        output_file (str): a path where result is written to
        num_reqs (int): number of requests to send
        address (str): address to the service being benchmarked
    """
    cmd_list = [
        'ghz',
        '--insecure', '--proto', proto_path, '-i', IMPORT_PATH, '--total',
        str(num_reqs), '--call', full_request_type, '-D', input_file,
        '-O', 'json', '-o', output_file, address,
    ]

    try:
        # call grpc GHZ load test tool
        subprocess.call(cmd_list)  # noqa: S603
    except subprocess.CalledProcessError as e:
        print(e.output)
        print('Check if gRPC GHZ tool is installed')
