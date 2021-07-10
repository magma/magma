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

from typing import List, NamedTuple

ClientCert = NamedTuple(
    'ClientCert',
    [('cert', str), ('key', str)],
)


class NetworkDNSConfig:
    def __init__(self, enable_caching: bool = False, local_ttl: int = 60):
        self.enable_caching = enable_caching
        self.local_ttl = local_ttl


class GenericNetwork:
    def __init__(
        self, id: str, name: str, description: str,
        dns: NetworkDNSConfig,
    ):
        self.id = id
        self.name = name
        self.description = description
        self.dns = dns


class MagmadGatewayConfigs:
    def __init__(
        self,
        autoupgrade_enabled: bool,
        autoupgrade_poll_interval: int,
        checkin_interval: int,
        checkin_timeout: int,
    ):
        self.autoupgrade_enabled = autoupgrade_enabled
        self.autoupgrade_poll_interval = autoupgrade_poll_interval
        self.checkin_interval = checkin_interval
        self.checkin_timeout = checkin_timeout


class ChallengeKey:
    def __init__(self, key_type: str):
        self.key_type = key_type


class GatewayDevice:
    def __init__(self, hardware_id: str, key: ChallengeKey):
        self.hardware_id = hardware_id
        self.key = key


class Gateway:
    def __init__(
        self,
        id: str,
        name: str, description: str,
        magmad: MagmadGatewayConfigs,
        device: GatewayDevice,
        tier: str,
    ):
        self.id = id
        self.name, self.description = name, description,
        self.magmad = magmad
        self.device = device
        self.tier = tier


class TierImage:
    def __init__(self, name: str, order: int):
        self.name = name
        self.order = order


class Tier:
    def __init__(
        self, id: str, version: str, images: List[TierImage],
        gateways: List[str],
    ):
        self.id = id
        self.version = version
        self.images = images
        self.gateways = gateways
