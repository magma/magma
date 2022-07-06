"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import importlib
import os


class Config(object):
    """
    Configuration class for Metricsd client
    """
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'DEBUG')

    # gRPC
    ORC8R_METRICSD_GRPC_SERVICE = os.environ.get('ORC8R_METRICSD_GRPC_SERVICE', 'orc8r-metricsd')
    ORC8R_METRICSD_GRPC_PORT = os.environ.get('ORC8R_METRICSD_GRPC_PORT', 9190)
    MAGMA_CLIENT_CERT_SERIAL_KEY = 'x-magma-client-cert-serial'
    MAGMA_CLIENT_CERT_SERIAL_VALUE = '7ZZXAF7CAETF241KL22B8YRR7B5UF401'


class TestConfig(Config):
    """
    Test configuration class for Metricsd client
    """
    pass  # noqa: WPS604


class ProductionConfig(Config):
    """
    Production configuration class for Metricsd client
    """
    pass  # noqa: WPS604


def get_config() -> Config:
    """
    Get configuration controller configuration
    """
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module(
        '.'.join(
            f"magma.metricsd_client.config.{app_config}".split('.')[
                :-1
            ],
        ),
    )
    config_class = getattr(config_module, app_config.split('.')[-1])
    return config_class()
