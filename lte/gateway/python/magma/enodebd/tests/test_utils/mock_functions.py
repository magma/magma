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

from typing import Any

GET_IP_FROM_IF_PATH = \
    'magma.enodebd.device_config.configuration_init.get_ip_from_if'

LOAD_SERVICE_MCONFIG_PATH = \
    'magma.enodebd.device_config.configuration_init.load_service_mconfig_as_json'


def mock_get_ip_from_if(
    _iface_name: str,
    _preference: Any = None,
) -> str:
    return '192.168.60.142'


def mock_load_service_mconfig_as_json(_service_name: str) -> Any:
    return {}
