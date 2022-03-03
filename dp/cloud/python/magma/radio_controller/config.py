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
import logging
import os
from distutils.util import strtobool

from magma.db_service import config as conf


class Config(object):
    """
    Configuration class for radio controller
    """
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

    # gRPC
    GRPC_PORT = int(os.environ.get('GRPC_PORT', 50053))

    # SQLAlchemy DB URI (scheme + url)
    SQLALCHEMY_DB_URI = conf.Config().SQLALCHEMY_DB_URI
    SQLALCHEMY_DB_ENCODING = conf.Config().SQLALCHEMY_DB_ENCODING
    SQLALCHEMY_ECHO = conf.Config().SQLALCHEMY_ECHO
    SQLALCHEMY_FUTURE = conf.Config().SQLALCHEMY_FUTURE

    # Elasticsearch
    ELASTICSEARCH_INDEX = os.environ.get('ELASTICSEARCH_INDEX', 'dp')

    # Fluentd
    FLUENTD_HOST = os.environ.get('FLUENTD_HOST', 'fluentd-service')
    FLUENTD_PORT = int(os.environ.get('FLUENTD_PORT', 24224))
    FLUENTD_TLS_ENABLED = strtobool(os.environ.get('FLUENTD_TLS_ENABLED', 'False'))
    FLUENTD_CERT_PATH = os.environ.get('FLUENTD_CERT_PATH', '')
    FLUENTD_KEY_PATH = os.environ.get('FLUENTD_KEY_PATH', '')
    FLUENTD_PROTOCOL = 'https' if FLUENTD_TLS_ENABLED else 'http'
    FLUENTD_URL = f'{FLUENTD_PROTOCOL}://{FLUENTD_HOST}:{FLUENTD_PORT}/{ELASTICSEARCH_INDEX}'


class DevelopmentConfig(Config):
    """
    Development configuration class for radio controller
    """

    pass  # noqa: WPS604


class TestConfig(Config):
    """
    Test configuration class for radio controller
    """

    SQLALCHEMY_DB_URI = conf.TestConfig().SQLALCHEMY_DB_URI


class ProductionConfig(Config):
    """
    Production configuration class for radio controller
    """

    pass  # noqa: WPS604


def get_config() -> Config:
    """
    Get Configuration object for radio controller
    """
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module(
        '.'.join(
            f"magma.radio_controller.config.{app_config}".split('.')[:-1],
        ),
    )
    config_class = getattr(config_module, app_config.split('.')[-1])
    logging.info(str(config_class))

    return config_class()
