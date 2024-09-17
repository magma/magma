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
    Configuration class for db service
    """
    # General
    LOG_LEVEL = os.environ.get('LOG_LEVEL', 'INFO')

    # SQLAlchemy DB URI (scheme + url)
    SQLALCHEMY_DB_URI = os.environ.get(
        'SQLALCHEMY_DB_URI', 'postgresql+psycopg2://postgres:postgres@db:5432/dp',
    )
    SQLALCHEMY_DB_ENCODING = os.environ.get('SQLALCHEMY_DB_ENCODING', 'utf8')
    SQLALCHEMY_ECHO = False
    SQLALCHEMY_FUTURE = False
    SQLALCHEMY_ENGINE_POOL_SIZE = os.environ.get(
        'SQLALCHEMY_ENGINE_POOL_SIZE', 6,
    )
    SQLALCHEMY_ENGINE_MAX_OVERFLOW = os.environ.get(
        'SQLALCHEMY_ENGINE_MAX_OVERFLOW', 10,
    )


class DevelopmentConfig(Config):
    """
    Configuration class for db service
    """
    pass  # noqa: WPS604


class TestConfig(Config):
    """
    Configuration class for db service
    """
    SQLALCHEMY_DB_URI = os.environ.get(
        'SQLALCHEMY_DB_URI', 'postgresql+psycopg2://postgres:postgres@db:5432/dp_test',
    )


class ProductionConfig(Config):
    """
    Configuration class for db service
    """
    pass  # noqa: WPS604
