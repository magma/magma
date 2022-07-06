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
import os
import subprocess
import sys
from typing import Any, List, Optional

import urllib3
from fabric.api import cd

sys.path.append('../../orc8r')
import tools.fab.dev_utils as dev_utils
import tools.fab.types as types
from tools.fab.hosts import vagrant_setup

LTE_NETWORK_TYPE = 'lte'
FEG_LTE_NETWORK_TYPE = 'feg_lte'
NIDS_BY_TYPE = {
    LTE_NETWORK_TYPE: 'test',
    FEG_LTE_NETWORK_TYPE: 'feg_lte_test',
}

FEG_FAB_PATH = '../../feg/gateway/'

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
                gx_gy_relay_enabled=False,
                hss_relay_enabled=False,
                network_services=[],
                mobility=MobilityConfig(
                    ip_allocation_mode='NAT',
                    nat=NatConfig(['192.168.128.0/24']),
                    reserved_addresses=[],
                ),
            ),
            ran=NetworkRANConfig(
                bandwidth_mhz=20,
                tdd_config=NetworkTDDConfig(
                    earfcndl=44590,
                    subframe_assignment=2, special_subframe_pattern=7,
                ),
            ),
            feg_network_id="",
        ),
        dns=types.NetworkDNSConfig(enable_caching=False, local_ttl=60),
    )
    _register_network(LTE_NETWORK_TYPE, network_payload)
    _register_agw(LTE_NETWORK_TYPE)


def register_vm_remote(certs_dir: str, network_id: str, url: str):
    """
    Register local VM gateway with remote controller.

    Example usage:
    fab -f dev_tools.py register_vm_remote:certs_dir=~/certs,network_id=test,url=https://api.stable.magmaeng.org
    """
    if None in {certs_dir, url, network_id}:
        print()
        print('==============================================================')
        print('Must provide the following arguments: certs_dir,url,network_id')
        print('==============================================================')
        return

    admin_cert = types.ClientCert(
        cert=os.path.expanduser(f'{certs_dir}/admin_operator.pem'),
        key=os.path.expanduser(f'{certs_dir}/admin_operator.key.pem'),
    )
    full_url = url + '/magma/v1/'

    _register_agw(
        LTE_NETWORK_TYPE,
        url=full_url,
        admin_cert=admin_cert,
        network_id=network_id,
    )


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
                gx_gy_relay_enabled=False,
                hss_relay_enabled=True,
                network_services=['dpi', 'policy_enforcement'],
                mobility=MobilityConfig(
                    ip_allocation_mode='NAT',
                    nat=NatConfig(['192.168.128.0/24']),
                    reserved_addresses=[],
                ),
            ),
            ran=NetworkRANConfig(
                bandwidth_mhz=20,
                tdd_config=NetworkTDDConfig(
                    earfcndl=44590,
                    subframe_assignment=2, special_subframe_pattern=7,
                ),
            ),
            feg_network_id='feg_test',
        ),
        dns=types.NetworkDNSConfig(enable_caching=False, local_ttl=60),
        federation=FederationNetworkConfig(feg_network_id='feg_test'),
    )
    _register_network(FEG_LTE_NETWORK_TYPE, network_payload)
    # registering gateway with LTE type. FEG_LTE doesn't have gateway endpoint
    _register_agw(FEG_LTE_NETWORK_TYPE)


def deregister_agw():
    """
    Remove AGW gateway from orc8r and remove certs from FEG gateway
    """
    dev_utils.delete_gateway_certs_from_vagrant('magma')
    _deregister_agw(LTE_NETWORK_TYPE)


def deregister_federated_agw():
    """
    Remove AGW gateway from orc8r and remove certs from FEG gateway
    """
    dev_utils.delete_gateway_certs_from_vagrant('magma')
    _deregister_agw(FEG_LTE_NETWORK_TYPE)


def register_feg_gw():
    """
    Registers FEG AGW gateway on orc8r
    """
    subprocess.check_call(
        'fab register_feg_gw', shell=True, cwd=FEG_FAB_PATH,
    )


def deregister_feg_gw():
    """
    Remove FEG gateway from orc8r and remove certs from FEG gateway
    """
    subprocess.check_call(
        'fab deregister_feg_gw', shell=True, cwd=FEG_FAB_PATH,
    )


def check_agw_cloud_connectivity(timeout=10):
    """
    Check connectivity of AGW with the cloud using checkin_cli.py
    Args:
        timeout: amount of time the command will retry
    """
    vagrant_setup("magma", destroy_vm=False, force_provision=False)
    with cd("/home/vagrant/build/python/bin/"):
        dev_utils.run_remote_command_with_repetition("./checkin_cli.py", timeout)


def check_agw_feg_connectivity(timeout=10):
    """
    Check connectivity of AGW with FEG feg_hello_cli.py
    Args:
        timeout: amount of time the command will retry
    """
    vagrant_setup("magma", destroy_vm=False, force_provision=False)
    with cd("/home/vagrant/build/python/bin/"):
        dev_utils.run_remote_command_with_repetition("./feg_hello_cli.py m 0", timeout)


def _register_network(network_type: str, payload: Any):
    network_id = NIDS_BY_TYPE[network_type]
    if not dev_utils.does_network_exist(network_id):
        dev_utils.cloud_post(network_type, payload)
    dev_utils.create_tier_if_not_exists(network_id, 'default')


def _register_agw(
        network_type: str,
        url: Optional[str] = None,
        admin_cert: Optional[types.ClientCert] = None,
        network_id: Optional[str] = None,
):
    network_id = network_id or NIDS_BY_TYPE[network_type]

    dev_utils.create_tier_if_not_exists(
        network_id,
        'default',
        url=url,
        admin_cert=admin_cert,
    )

    hw_id = dev_utils.get_gateway_hardware_id_from_vagrant(vm_name='magma')
    already_registered, registered_as = dev_utils.is_hw_id_registered(
        network_id,
        hw_id,
        url=url,
        admin_cert=admin_cert,
    )
    if already_registered:
        print()
        print(f'===========================================')
        print(f'VM is already registered as {registered_as}')
        print(f'===========================================')
        return

    gw_id = dev_utils.get_next_available_gateway_id(
        network_id,
        url=url,
        admin_cert=admin_cert,
    )
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
    dev_utils.cloud_post(
        f'lte/{network_id}/gateways',
        gw_payload,
        url=url,
        admin_cert=admin_cert,
    )
    print()
    print(f'=========================================')
    print(f'Gateway {gw_id} successfully provisioned!')
    print(f'=========================================')


def _deregister_agw(network_type: str):
    network_id = NIDS_BY_TYPE[network_type]
    hw_id = dev_utils.get_gateway_hardware_id_from_vagrant(vm_name='magma')
    already_registered, registered_as = dev_utils.is_hw_id_registered(
        network_id, hw_id,
    )
    if not already_registered:
        print()
        print('===========================================')
        print(f'VM is not registered (hwid: {hw_id} )')
        print('===========================================')
        return

    dev_utils.cloud_delete(f'lte/{network_id}/gateways/{registered_as}')
    print()
    print('=========================================')
    print(f'AGW Gateway {registered_as} successfully removed!')
    print(f'(restart AGW services on magma vm)')
    print('=========================================')


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


class NatConfig:
    def __init__(self, ip_blocks: List[str]):
        self.ip_blocks = ip_blocks


class MobilityConfig:
    def __init__(
        self, ip_allocation_mode: str, nat: NatConfig,
        reserved_addresses: List[str],
    ):
        self.ip_allocation_mode = ip_allocation_mode
        self.nat = nat
        self.reserved_addresses = reserved_addresses


class NetworkEPCConfig:
    def __init__(
        self, lte_auth_amf: str, lte_auth_op: str,
        mcc: str, mnc: str, tac: int,
        gx_gy_relay_enabled: bool, hss_relay_enabled: bool,
        network_services: List[str],
        mobility: MobilityConfig,
    ):
        self.lte_auth_amf = lte_auth_amf
        self.lte_auth_op = lte_auth_op
        self.mcc = mcc
        self.mnc = mnc
        self.tac = tac
        self.gx_gy_relay_enabled = gx_gy_relay_enabled
        self.hss_relay_enabled = hss_relay_enabled
        self.default_rule_id = "default_rule_1"
        self.network_services = network_services
        self.mobility = mobility


class NetworkCellularConfig:
    def __init__(
        self, epc: NetworkEPCConfig, ran: NetworkRANConfig,
        feg_network_id: str,
    ):
        self.epc = epc
        self.ran = ran
        self.feg_network_id = feg_network_id


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
        self.dns_primary = "8.8.8.8"
        self.dns_secondary = "8.8.4.4"


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
