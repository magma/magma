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
from typing import Dict, List

import urllib3
from fabric.api import hide, lcd, local

sys.path.append('../../orc8r')
import tools.fab.dev_utils as dev_utils  # NOQA
import tools.fab.types as types

NETWORK_ID = 'feg_test'

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def register_feg():
    _register_federation_network()
    _register_feg()


class RadiusConfig:
    def __init__(
        self,
        DAE_addr: str = '127.0.0.1:3799',
        acct_addr: str = '127.0.0.1:1813',
        auth_addr: str = '127.0.0.1:1812',
        network: str = 'udp',
        secret: str = 'MTIzNDU2',
    ):
        self.DAE_addr = DAE_addr
        self.acct_addr = acct_addr
        self.auth_addr = auth_addr
        self.network = network
        self.secret = secret


class AAAServerConfig:
    def __init__(
        self,
        accounting_enabled: bool = True,
        create_session_on_auth: bool = True,
        event_logging_enabled: bool = False,
        idle_session_timeout_ms: int = 21600000,
        radius_config: RadiusConfig = RadiusConfig(),
    ):
        self.accounting_enabled = accounting_enabled
        self.create_session_on_auth = create_session_on_auth
        self.event_logging_enabled = event_logging_enabled
        self.idle_session_timeout_ms = idle_session_timeout_ms
        self.radius_config = radius_config


class EapAkaTimeout:
    def __init__(
        self,
        challenge_ms: int = 20000, error_notification_ms: int = 10000,
        session_authenticated_ms: int = 5000,
        session_ms: int = 43200000,
    ):
        self.challenge_ms = challenge_ms
        self.error_notification_ms = error_notification_ms
        self.session_authenticated_ms = session_authenticated_ms
        self.session_ms = session_ms


class EapAkaConfig:
    def __init__(
        self, plmn_ids: List[str] = None,
        timeout: EapAkaTimeout = EapAkaTimeout(),
    ):
        if plmn_ids is None:
            plmn_ids = ['123456']
        self.plmn_ids = plmn_ids
        self.timeout = timeout


class DiamServerConfig:
    def __init__(
        self,
        address: str = 'localhost:1234',
        dest_host: str = 'magma-fedgw.magma.com',
        dest_realm: str = 'magma.com',
        disable_dest_host: bool = False,
        host: str = 'string',
        local_address: str = ':56789',
        overwrite_dest_host: bool = False,
        product_name: str = 'string',
        protocol: str = 'tcp',
        realm: str = 'string',
        retransmits: int = 0,
        retry_count: int = 0,
        watchdog_interval: int = 0,
    ):
        self.address = address
        self.dest_host = dest_host
        self.dest_realm = dest_realm
        self.disable_dest_host = disable_dest_host
        self.host = host
        self.local_address = local_address
        self.overwrite_dest_host = overwrite_dest_host
        self.product_name = product_name
        self.protocol = protocol
        self.realm = realm
        self.retransmits = retransmits
        self.retry_count = retry_count
        self.watchdog_interval = watchdog_interval


class GxConfig:
    def __init__(
        self,
        disableGx: bool = False,
        servers: List[DiamServerConfig] = None,
    ):
        self.disableGx = disableGx
        if servers is None:
            servers = [DiamServerConfig()]
        self.servers = servers


class GyConfig:
    def __init__(
        self,
        disableGy: bool = False,
        init_method: int = 2,
        servers: List[DiamServerConfig] = None,
    ):
        self.disableGy = disableGy
        self.init_method = init_method
        if servers is None:
            servers = [DiamServerConfig()]
        self.servers = servers


class HealthConfigs:
    def __init__(
        self,
        cloud_disable_period_secs: int = 10,
        cpu_utilization_threshold: float = 0.9,
        health_services: List[str] = None,
        local_disable_period_secs: int = 1,
        memory_available_threshold: float = 0.75,
        minimum_request_threshold: int = 1,
        request_failure_threshold: float = 0.5,
        update_failure_threshold: int = 3,
        update_interval_secs: int = 10,
    ):
        self.cloud_disable_period_secs = cloud_disable_period_secs
        self.cpu_utilization_threshold = cpu_utilization_threshold
        if health_services is None:
            health_services = ['SESSION_PROXY', 'SWX_PROXY']
        self.health_services = health_services
        self.local_disable_period_secs = local_disable_period_secs
        self.memory_available_threshold = memory_available_threshold
        self.minimum_request_threshold = minimum_request_threshold
        self.request_failure_threshold = request_failure_threshold
        self.update_failure_threshold = update_failure_threshold
        self.update_interval_secs = update_interval_secs


class SubProfile:
    def __init__(
        self,
        max_dl_bit_rate: int = 20000000,
        max_ul_bit_rate: int = 10000000,
    ):
        self.max_dl_bit_rate = max_dl_bit_rate
        self.max_ul_bit_rate = max_ul_bit_rate


class HssServer:
    def __init__(
        self,
        address: str = 'localhost:1234',
        dest_host: str = 'magma-fedgw.magma.com',
        dest_realm: str = 'magma.com',
        local_address: str = ':56789',
        protocol: str = 'tcp',
    ):
        self.address = address
        self.dest_host = dest_host
        self.dest_realm = dest_realm
        self.local_address = local_address
        self.protocol = protocol


class HssConfigs:
    def __init__(
        self,
        default_sub_profile: SubProfile = SubProfile(),
        lte_auth_amf: str = 'gAA=',
        lte_auth_op: str = 'EREREREREREREREREREREQ==',
        server: HssServer = HssServer(),
        stream_subscribers: bool = False,
        sub_profiles: Dict[str, SubProfile] = None,
    ):
        self.default_sub_profile = default_sub_profile
        self.lte_auth_amf = lte_auth_amf
        self.lte_auth_op = lte_auth_op
        self.server = server
        self.stream_subscribers = stream_subscribers
        if sub_profiles is None:
            sub_profiles = {'additionalProp1': SubProfile()}
        self.sub_profiles = sub_profiles


class S6aConfigs:
    def __init__(
        self,
        plmn_ids: List[str] = None,
        server: DiamServerConfig = DiamServerConfig(),
    ):
        if plmn_ids is None:
            plmn_ids = ["123456"]
        self.plmn_ids = plmn_ids
        self.server = server


class SwxConfigs:
    def __init__(
        self,
        cache_TTL_seconds: int = 10800,
        derive_unregister_realm: bool = False,
        hlr_plmn_ids: List[str] = None,
        register_on_auth: bool = False,
        servers: List[DiamServerConfig] = None,
        verify_authorization: bool = False,
    ):
        self.cache_TTL_seconds = cache_TTL_seconds
        self.derive_unregister_realm = derive_unregister_realm
        if hlr_plmn_ids is None:
            hlr_plmn_ids = ['00101']
        self.hlr_plmn_ids = hlr_plmn_ids
        self.register_on_auth = register_on_auth
        if servers is None:
            servers = [DiamServerConfig()]
        self.servers = servers
        self.verify_authorization = verify_authorization


class SubConfig:
    def __init__(
        self,
        network_wide_base_names: List[str] = None,
        network_wide_rule_names: List[str] = None,
    ):
        if network_wide_base_names is None:
            network_wide_base_names = ['base_1']
        self.network_wide_base_names = network_wide_base_names
        if network_wide_rule_names is None:
            network_wide_rule_names = ['rule_1']
        self.network_wide_rule_names = network_wide_rule_names


class FederationNetworkConfigs:
    def __init__(
        self,
        served_network_ids: List[str] = None,
        aaa_server: AAAServerConfig = AAAServerConfig(),
        eap_aka: EapAkaConfig = EapAkaConfig(),
        gx: GxConfig = GxConfig(),
        gy: GyConfig = GyConfig(),
        health: HealthConfigs = HealthConfigs(),
        hss: HssConfigs = HssConfigs(),
        s6a: S6aConfigs = S6aConfigs(),
        swx: SwxConfigs = SwxConfigs(),
    ):
        if served_network_ids is None:
            served_network_ids = ['feg_lte_test']
        self.served_network_ids = served_network_ids
        self.aaa_server = aaa_server
        self.eap_aka = eap_aka
        self.gx = gx
        self.gy = gy
        self.health = health
        self.hss = hss
        self.s6a = s6a
        self.swx = swx


class FederationNetwork:
    def __init__(
        self,
        id: str = NETWORK_ID,
        name: str = 'Testing',
        description: str = 'Test federation network',
        federation: FederationNetworkConfigs = FederationNetworkConfigs(),
        dns: types.NetworkDNSConfig = types.NetworkDNSConfig(),
        subscriber_config: SubConfig = SubConfig(),
    ):
        self.id = id
        self.name = name
        self.description = description
        self.dns = dns
        self.federation = federation
        self.subscriber_config = subscriber_config


class FederationGateway:
    def __init__(
        self,
        id: str, name: str, description: str,
        device: types.GatewayDevice,
        magmad: types.MagmadGatewayConfigs,
        tier: str = 'default',
        federation: FederationNetworkConfigs = FederationNetworkConfigs(),
    ):
        self.id = id
        self.name = name
        self.description = description
        self.device = device
        self.magmad = magmad
        self.tier = tier
        self.federation = federation


def _register_federation_network(payload: FederationNetwork = FederationNetwork()):
    nid = payload.id
    if not dev_utils.does_network_exist(nid):
        dev_utils.cloud_post('feg', payload)

    dev_utils.create_tier_if_not_exists(nid, 'default')


def _register_feg():
    with lcd('docker'), hide('output', 'running', 'warnings'):
        hw_id = local(
            'docker-compose exec magmad bash -c "cat /etc/snowflake"',
            capture=True,
        )
    already_registered, registered_as = dev_utils.is_hw_id_registered(
        NETWORK_ID, hw_id,
    )
    if already_registered:
        print()
        print(f'============================================')
        print(f'Feg is already registered as {registered_as}')
        print(f'============================================')
        return

    gw_id = dev_utils.get_next_available_gateway_id(NETWORK_ID)
    md_gw = dev_utils.construct_magmad_gateway_payload(gw_id, hw_id)
    gw_payload = FederationGateway(
        id=md_gw.id, name=md_gw.name, description=md_gw.description,
        device=md_gw.device,
        magmad=md_gw.magmad, tier=md_gw.tier,
    )
    dev_utils.cloud_post(f'feg/{NETWORK_ID}/gateways', gw_payload)
    print()
    print(f'=====================================')
    print(f'Feg {gw_id} successfully provisioned!')
    print(f'=====================================')
