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
    HTTP_SERVER = os.environ.get(
        'HTTP_SERVER', 'https://orc8r-nginx-proxy',
    )

    # Security
    DP_CERT_PATH = os.environ.get(
        'DP_CERT_PATH', '/backend/test_runner/certs/admin_operator.pem',
    )
    DP_SSL_KEY_PATH = os.environ.get(
        'DP_SSL_KEY_PATH', '/backend/test_runner/certs/admin_operator.key.pem',
    )

    # Test Elasticsearch
    ELASTICSEARCH_HOST = os.environ.get('ELASTICSEARCH_HOST', 'elasticsearch-service')
    ELASTICSEARCH_PORT = int(os.environ.get('ELASTICSEARCH_PORT', 9200))
    ELASTICSEARCH_INDEX = os.environ.get('ELASTICSEARCH_INDEX', 'dp')
    ELASTICSEARCH_URL = f"http://{ELASTICSEARCH_HOST}:{ELASTICSEARCH_PORT}"

    # Test Fluentd
    FLUENTD_HOST = os.environ.get('FLUENTD_HOST', 'fluentd-service')
    FLUENTD_PORT = int(os.environ.get('FLUENTD_PORT', 24224))
