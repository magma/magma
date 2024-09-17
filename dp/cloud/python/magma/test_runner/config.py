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
    ORC8R_DP_GRPC_SERVICE = os.environ.get(
        'ORC8R_DP_GRPC_SERVICE', 'orc8r-dp',
    )
    ORC8R_DP_GRPC_PORT = int(os.environ.get('ORC8R_DP_GRPC_PORT', 9180))
    HTTP_SERVER = os.environ.get(
        'HTTP_SERVER', 'https://orc8r-nginx-proxy',
    )
    ORC8R_METRICSD_GRPC_SERVICE = os.environ.get('ORC8R_METRICSD_GRPC_SERVICE', 'orc8r-metricsd')
    ORC8R_METRICSD_GRPC_PORT = os.environ.get('ORC8R_METRICSD_GRPC_PORT', 9190)

    # Security
    DP_CERT_PATH = os.environ.get(
        'DP_CERT_PATH', '/backend/test_runner/certs/admin_operator.pem',
    )
    DP_SSL_KEY_PATH = os.environ.get(
        'DP_SSL_KEY_PATH', '/backend/test_runner/certs/admin_operator.key.pem',
    )

    # Test Elasticsearch
    ELASTICSEARCH_SERVICE_HOST = os.environ.get('ELASTICSEARCH_SERVICE_HOST', '')
    ELASTICSEARCH_SERVICE_PORT = int(os.environ.get('ELASTICSEARCH_SERVICE_PORT', 9200))
    ELASTICSEARCH_INDEX = os.environ.get('ELASTICSEARCH_INDEX', 'dp')
    ELASTICSEARCH_URL = f"http://{ELASTICSEARCH_SERVICE_HOST}:{ELASTICSEARCH_SERVICE_PORT}"

    # Test Fluentd
    FLUENTD_SERVICE_HOST = os.environ.get('DOMAIN_PROXY_FLUENTD_SERVICE_HOST', '')
    FLUENTD_SERVICE_PORT = int(os.environ.get('DOMAIN_PROXY_FLUENTD_SERVICE_PORT', 9888))

    # Test Prometheus
    PROMETHEUS_SERVICE_HOST = os.environ.get('PROMETHEUS_SERVICE_HOST', 'orc8r-prometheus')
    PROMETHEUS_SERVICE_PORT = int(os.environ.get('PROMETHEUS_SERVICE_PORT', 9090))
