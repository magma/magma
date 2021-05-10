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
import redis
from magma.configuration.service_configs import get_service_config_value


def get_default_client():
    """
    Return a default redis client using the configured port in redis.yml
    """
    redis_port = get_service_config_value('redis', 'port', 6379)
    redis_addr = get_service_config_value('redis', 'bind', 'localhost')
    return redis.Redis(host=redis_addr, port=redis_port)
