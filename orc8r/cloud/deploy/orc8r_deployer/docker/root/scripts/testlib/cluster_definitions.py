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


from dataclasses import dataclass
from enum import Enum
from typing import Dict, List

from dataclasses_json import dataclass_json


class ClusterType(Enum):
    LOCAL = 0
    AWS = 1


@dataclass_json
@dataclass
class Orc8rTemplate:
    infra: Dict[str, any]
    platform: Dict[str, any]
    service: Dict[str, any]


@dataclass_json
@dataclass
class GatewayTemplate:
    prefix: str
    count: int
    ami: str
    cloudstrapper_ami: str
    region: str
    az: str
    service_config: dict[str, any]


@dataclass_json
@dataclass
class ClusterTemplate:
    orc8r: Orc8rTemplate
    gateway: GatewayTemplate


@dataclass_json
@dataclass
class GatewayConfig:
    gateway_id: str
    hardware_id: str
    hostname: str


@dataclass_json
@dataclass
class ClusterInternalConfig:
    bastion_ip: str


@dataclass_json
@dataclass
class ClusterConfig:
    uuid: str
    cluster_type: ClusterType
    internal_config: ClusterInternalConfig
    template: ClusterTemplate
    gateways: List[GatewayConfig]


class ClusterCreateError(Exception):
    def __init__(self, msg):
        super().__init__(msg)


class ClusterDestroyError(Exception):
    def __init__(self, msg):
        super().__init__(msg)


class ClusterVerifyError(Exception):
    def __init__(self):
        super().__init__()
