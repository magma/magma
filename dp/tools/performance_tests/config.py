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
    Configuration for performance tests
    """
    DB_HOST = os.environ.get('DB_HOST', 'localhost')
    DB_PORT = os.environ.get('DB_PORT', '5532')
    DB_USER = os.environ.get('DB_USER', 'postgres')
    DB_NAME = os.environ.get('DB_NAME', 'dp')
    DB_PASSWORD = os.environ.get('DB_PASSWORD', 'postgres')
    LOG_FILE = os.environ.get('LOG_FILE', 'performance_tests.log')
    MAX_WORKERS = os.environ.get('MAX_WORKERS', 8)
    MIN_CONNECTIONS = os.environ.get('MIN_CONNECTIONS', 1)
    MAX_CONNECTIONS = os.environ.get('MAX_CONNECTIONS', 800)
