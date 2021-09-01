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

import sys
from typing import Any, List

import urllib3

sys.path.append('../../orc8r')
import tools.fab.dev_utils as dev_utils
import tools.fab.types as types

LTE_NETWORK_TYPE = 'lte'
FEG_LTE_NETWORK_TYPE = 'feg_lte'
NIDS_BY_TYPE = {
    LTE_NETWORK_TYPE: 'test',
    FEG_LTE_NETWORK_TYPE: 'feg_lte_test',
}

# Disable warnings about SSL verification since its a local VM
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def register_vm():
    network_payload = LTENetwork(
        id=NIDS_BY_TYPE[LTE_NETWORK_TYPE],
        name='Test Network', description='Test Network',
        cellular=NetworkCellularConfig(
            epc=NetworkEPCConfig(
                lte_auth_amf='gAA=',
                lte_auth_op='EREREREREREREREREREREQ==',
                mcc='001', mnc='01', tac=1,
                relay_enabled=False,
            ),
            ran=NetworkRANConfig(
                bandwidth_mhz=20,
                tdd_config=NetworkTDDConfig(
                    earfcndl=44590,
                    subframe_assignment=2, special_subframe_pattern=7,
                ),
            ),
        ),
        dns=types.NetworkDNSConfig(enable_caching=False, local_ttl=60),
    )
    _register_network(LTE_NETWORK_TYPE, network_payload)
    _register_agw(LTE_NETWORK_TYPE)


def register_federated_vm():
    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    network_payload = FederatedLTENetwork(
        id=NIDS_BY_TYPE[FEG_LTE_NETWORK_TYPE],
        name='Test Network', description='Test Network',
        cellular=NetworkCellularConfig(
            epc=NetworkEPCConfig(
                lte_auth_amf='gAA=',
                lte_auth_op='EREREREREREREREREREREQ==',
                mcc='001', mnc='01', tac=1,
                relay_enabled=False,
            ),
            ran=NetworkRANConfig(
                bandwidth_mhz=20,
                tdd_config=NetworkTDDConfig(
                    earfcndl=44590,
                    subframe_assignment=2, special_subframe_pattern=7,
                ),
            ),
        ),
        dns=types.NetworkDNSConfig(enable_caching=False, local_ttl=60),
        federation=FederationNetworkConfig(feg_network_id='feg_test'),
    )
    _register_network(FEG_LTE_NETWORK_TYPE, network_payload)
    _register_agw(FEG_LTE_NETWORK_TYPE)


def _register_network(network_type: str, payload: Any):
    network_id = NIDS_BY_TYPE[network_type]
    if not dev_utils.does_network_exist(network_id):
        dev_utils.cloud_post(network_type, payload)

    dev_utils.create_tier_if_not_exists(network_id, 'default')


def _register_agw(network_type: str):
    network_id = NIDS_BY_TYPE[network_type]
    hw_id = dev_utils.get_hardware_id_from_vagrant(vm_name='magma')
    already_registered, registered_as = dev_utils.is_hw_id_registered(
        network_id, hw_id,
    )
    if already_registered:
        print()
        print(f'===========================================')
        print(f'VM is already registered as {registered_as}')
        print(f'===========================================')
        return

    gw_id = dev_utils.get_next_available_gateway_id(network_id)
    md_gw = dev_utils.construct_magmad_gateway_payload(gw_id, hw_id)
    gw_payload = LTEGateway(
        device=md_gw.device,
        id=gw_id, name=md_gw.name, description=md_gw.description,
        magmad=md_gw.magmad, tier=md_gw.tier,
        cellular=GatewayCellularConfig(
            epc=GatewayEPCConfig(
                ip_block='192.168.128.0/24',
                nat_enabled=True,
            ),
            ran=GatewayRANConfig(pci=260, transmit_enabled=True),
        ),
        connected_enodeb_serials=[],
    )
    dev_utils.cloud_post(f'{network_type}/{network_id}/gateways', gw_payload)
    print()
    print(f'=========================================')
    print(f'Gateway {gw_id} successfully provisioned!')
    print(f'=========================================')


class NetworkTDDConfig:
    def __init__(
        self, earfcndl: int,
        subframe_assignment: int, special_subframe_pattern: int,
    ):
        self.earfcndl = earfcndl
        self.subframe_assignment = subframe_assignment
        self.special_subframe_pattern = special_subframe_pattern


class NetworkRANConfig:
    def __init__(self, bandwidth_mhz: int, tdd_config: NetworkTDDConfig):
        self.bandwidth_mhz = bandwidth_mhz
        self.tdd_config = tdd_config


class NetworkEPCConfig:
    def __init__(
        self, lte_auth_amf: str, lte_auth_op: str,
        mcc: str, mnc: str, tac: int,
        relay_enabled: bool,
    ):
        self.lte_auth_amf = lte_auth_amf
        self.lte_auth_op = lte_auth_op
        self.mcc = mcc
        self.mnc = mnc
        self.tac = tac
        self.gx_gy_relay_enabled = relay_enabled
        self.hss_relay_enabled = relay_enabled


class NetworkCellularConfig:
    def __init__(self, epc: NetworkEPCConfig, ran: NetworkRANConfig):
        self.epc = epc
        self.ran = ran


class LTENetwork:
    def __init__(
        self, id: str, name: str, description: str,
        cellular: NetworkCellularConfig,
        dns: types.NetworkDNSConfig,
    ):
        self.id = id
        self.name = name
        self.description = description
        self.cellular = cellular
        self.dns = dns


class FederationNetworkConfig:
    def __init__(self, feg_network_id: str):
        self.feg_network_id = feg_network_id


class FederatedLTENetwork:
    def __init__(
        self, id: str, name: str, description: str,
        cellular: NetworkCellularConfig,
        dns: types.NetworkDNSConfig,
        federation: FederationNetworkConfig,
    ):
        self.id = id
        self.name = name
        self.description = description
        self.cellular = cellular
        self.dns = dns
        self.federation = federation


class GatewayRANConfig:
    def __init__(self, pci: int, transmit_enabled: bool):
        self.pci = pci
        self.transmit_enabled = transmit_enabled


class GatewayEPCConfig:
    def __init__(self, ip_block: str, nat_enabled: bool):
        self.ip_block = ip_block
        self.nat_enabled = nat_enabled


class GatewayCellularConfig:
    def __init__(self, epc: GatewayEPCConfig, ran: GatewayRANConfig):
        self.epc = epc
        self.ran = ran


class LTEGateway:
    def __init__(
        self, device: types.GatewayDevice,
        id: str, name: str, description: str,
        magmad: types.MagmadGatewayConfigs,
        tier: str,
        cellular: GatewayCellularConfig,
        connected_enodeb_serials: List[str],
    ):
        self.device = device
        self.id = id
        self.name = name
        self.description = description
        self.magmad = magmad
        self.tier = tier
        self.cellular = cellular
        self.connected_enodeb_serials = connected_enodeb_serials
