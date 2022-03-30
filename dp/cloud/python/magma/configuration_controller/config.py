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

from magma.db_service import config as conf
from magma.mappings.request_mapping import request_mapping


class Config(object):
    """
    Configuration class for Configuration Controller
    """
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'DEBUG')
    REQUEST_PROCESSING_INTERVAL_SEC = int(
        os.environ.get('REQUEST_PROCESSING_INTERVAL_SEC', 10),
    )
    REQUEST_PROCESSING_LIMIT = int(
        os.environ.get('REQUEST_PROCESSING_LIMIT', 100),
    )

    # Services
    SAS_URL = os.environ.get('SAS_URL', 'https://fake-sas-service/v1.2')
    RC_INGEST_URL = os.environ.get('RC_INGEST_URL', '')

    # SQLAlchemy
    SQLALCHEMY_DB_URI = conf.Config().SQLALCHEMY_DB_URI
    SQLALCHEMY_DB_ENCODING = conf.Config().SQLALCHEMY_DB_ENCODING
    SQLALCHEMY_ECHO = conf.Config().SQLALCHEMY_ECHO
    SQLALCHEMY_FUTURE = conf.Config().SQLALCHEMY_FUTURE
    # DB engine connection pool size will default to the amount of request types
    # as each request type has its own query thread
    SQLALCHEMY_ENGINE_POOL_SIZE = int(
        os.environ.get(
            'SQLALCHEMY_ENGINE_POOL_SIZE',
            len(request_mapping),
        ),
    )
    SQLALCHEMY_ENGINE_MAX_OVERFLOW = int(
        os.environ.get(
            'SQLALCHEMY_ENGINE_MAX_OVERFLOW',
            conf.Config().SQLALCHEMY_ENGINE_MAX_OVERFLOW,
        ),
    )

    # Security
    CC_CERT_PATH = os.environ.get(
        'CC_CERT_PATH', '/backend/configuration_controller/certs/tls.crt',
    )
    CC_SSL_KEY_PATH = os.environ.get(
        'CC_SSL_KEY_PATH', '/backend/configuration_controller/certs/tls.key',
    )
    SAS_CERT_PATH = os.environ.get(
        'SAS_CERT_PATH', '/backend/configuration_controller/certs/ca.crt',
    )


class DevelopmentConfig(Config):
    """
    Development configuration class for Configuration Controller
    """

    pass  # noqa: WPS604


class TestConfig(Config):
    """
    Test configuration class for Configuration Controller
    """

    SQLALCHEMY_DB_URI = conf.TestConfig().SQLALCHEMY_DB_URI


class ProductionConfig(Config):
    """
    Production configuration class for Configuration Controller
    """

    SQLALCHEMY_ECHO = False
