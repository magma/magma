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
import importlib
import os
from distutils.util import strtobool


class Config(object):
    """
    Configuration class for Configuration Controller
    """
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'DEBUG')

    # Elasticsearch
    ELASTICSEARCH_INDEX = os.environ.get('ELASTICSEARCH_INDEX', 'dp')

    # Fluentd
    FLUENTD_SERVICE_HOST = os.environ.get('DOMAIN_PROXY_FLUENTD_SERVICE_HOST', 'domain-proxy-fluentd')
    FLUENTD_SERVICE_PORT = int(os.environ.get('DOMAIN_PROXY_FLUENTD_SERVICE_PORT', 9888))
    FLUENTD_TLS_ENABLED = strtobool(os.environ.get('FLUENTD_TLS_ENABLED', 'False'))
    FLUENTD_CERT_PATH = os.environ.get('FLUENTD_CERT_PATH', '')
    FLUENTD_KEY_PATH = os.environ.get('FLUENTD_KEY_PATH', '')
    FLUENTD_PROTOCOL = 'https' if FLUENTD_TLS_ENABLED else 'http'
    FLUENTD_URL = f'{FLUENTD_PROTOCOL}://{FLUENTD_SERVICE_HOST}:{FLUENTD_SERVICE_PORT}/{ELASTICSEARCH_INDEX}'


class TestConfig(Config):
    """
    Test configuration class for Configuration Controller
    """
    pass  # noqa: WPS604


class ProductionConfig(Config):
    """
    Production configuration class for Configuration Controller
    """
    pass  # noqa: WPS604


def get_config() -> Config:
    """
    Get configuration controller configuration
    """
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module(
        '.'.join(
            f"magma.fluentd_client.config.{app_config}".split('.')[
                :-1
            ],
        ),
    )
    config_class = getattr(config_module, app_config.split('.')[-1])
    return config_class()
