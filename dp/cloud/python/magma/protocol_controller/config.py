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


class Config(object):
    """
    Configuration class for protocol controller
    """
    # General
    TESTING = False
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'DEBUG')
    API_PREFIX = os.environ.get('API_PREFIX', '/sas/v1')
    RC_RESPONSE_WAIT_TIMEOUT_SEC = int(
        os.environ.get('RC_RESPONSE_WAIT_TIMEOUT_SEC', 60),
    )
    RC_RESPONSE_WAIT_INTERVAL_SEC = int(
        os.environ.get('RC_RESPONSE_WAIT_INTERVAL_SEC', 1),
    )

    PROTOCOL_PLUGIN = os.environ.get(
        'PROTOCOL_PLUGIN', 'magma.protocol_controller.plugins.cbsd_sas.plugin.CBSDSASProtocolPlugin',
    )

    # gRPC
    GRPC_SERVICE = os.environ.get(
        'GRPC_SERVICE', 'domain-proxy-radio-controller',
    )
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))

    JSON_ADD_STATUS = True
    JSON_STATUS_FIELD_NAME = '__status'
    JSON_JSONIFY_HTTP_ERRORS = True
    JSON_USE_ENCODE_METHODS = True


class DevelopmentConfig(Config):
    """
    Development configuration class for protocol controller
    """

    pass  # noqa: WPS604


class TestConfig(Config):
    """
    Test class for protocol controller
    """

    pass  # noqa: WPS604


class ProductionConfig(Config):
    """
    Prodiction class for protocol controller
    """

    TESTING = False
