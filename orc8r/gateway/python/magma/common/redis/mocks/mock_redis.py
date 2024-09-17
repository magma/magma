"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from redis.exceptions import RedisError

# For non-failure cases, just use the fakeredis module


class MockUnavailableRedis(object):
    """
    MockUnavailableRedis implements a mock Redis Server that always raises
    a connection exception
    """

    def __init__(self, host, port):
        self.host = host
        self.port = port

    def lock(self, key):
        raise RedisError("mock redis error")

    def keys(self, pattern=".*"):
        """ Mock keys with regex pattern matching."""
        raise RedisError("mock redis error")
