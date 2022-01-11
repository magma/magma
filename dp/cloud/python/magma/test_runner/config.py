"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os


class TestConfig(object):
    """
    Configuration class for test runner
    """
    # General
    CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX = os.environ.get(
        'CBSD_SAS_PROTOCOL_CONTROLLER_API_PREFIX',
        "http://domain-proxy-protocol-controller:8080/sas/v1",
    )
    GRPC_SERVICE = os.environ.get(
        'GRPC_SERVICE', 'domain-proxy-radio-controller',
    )
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))
