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
from enum import Enum
from typing import Any, Dict, List

import requests
import swagger_client
from integ_tests.cloud.fixtures import (
    DEFAULT_GATEWAY_CELLULAR_CONFIG,
    DEFAULT_GATEWAY_CONFIG,
    DEFAULT_NETWORK_CELLULAR_CONFIG,
    DEFAULT_NETWORK_DNSD_CONFIG,
)
from swagger_client.rest import ApiException

TEST_CLOUD_HOSTNAME = 'controller.magma.test'
TEST_CLOUD_PORT = '9443'
CERT_DIR = '/home/vagrant/magma/.cache/test_certs'


class CloudManager(object):
    """
    Manager for magma cloud
    """

    def __init__(
        self, hostname=TEST_CLOUD_HOSTNAME, port=TEST_CLOUD_PORT,
        ssl_ca_cert=CERT_DIR + '/rootCA.pem',
        cert_file=CERT_DIR + '/admin_operator.pem',
        key_file=CERT_DIR + '/admin_operator.key.pem',
    ):
        self._hostname = hostname
        self._host = 'https://{hostname}:{port}/magma'.format(
            hostname=hostname, port=port,
        )
        self._managed_networks = []     # type: List[str]

        self._ssl_ca = ssl_ca_cert
        self._ssl_cert = cert_file
        self._ssl_key = key_file

        # Set SSL config. NOTE: must be done before creating ApiClient
        configuration = swagger_client.Configuration()
        configuration.ssl_ca_cert = ssl_ca_cert
        configuration.cert_file = cert_file
        configuration.key_file = key_file
        configuration.verify_ssl = True
        configuration.assert_hostname = True

        api_client = swagger_client.ApiClient(host=self._host)

        self.subscribers_api = swagger_client.SubscribersApi(
            api_client=api_client,
        )
        self.networks_api = swagger_client.NetworksApi(api_client=api_client)
        self.gateways_api = swagger_client.GatewaysApi(api_client=api_client)
        self.tiers_api = swagger_client.TiersApi(api_client=api_client)
        self.channels_api = swagger_client.ChannelsApi(api_client=api_client)

    def create_network(self, network_id: str):
        """
        Synchronously create a network, a default tier for that network, and
        a default cellular configuration for the network.

        If the network ID already exists, no action will be performed

        Args:
            network_id (str): requested network ID
        """
        # Always mark the network ID as managed so it gets cleaned up
        self._managed_networks.append(network_id)
        if network_id in self.networks_api.networks_get():
            return

        network_record = swagger_client.NetworkRecord(name='integ_net_name')
        self.networks_api.networks_post(
            network_record,
            requested_id=network_id,
        )

        self.networks_api.networks_network_id_configs_cellular_post(
            network_id,
            DEFAULT_NETWORK_CELLULAR_CONFIG,
        )
        self.networks_api.networks_network_id_configs_dns_post(
            network_id,
            DEFAULT_NETWORK_DNSD_CONFIG,
        )

        tier_list = self.tiers_api.networks_network_id_tiers_get(network_id)
        tier_exists = (len(tier_list) > 0)
        if not tier_exists:
            self.tiers_api.networks_network_id_tiers_post(
                network_id,
                swagger_client.Tier(
                    id='default',
                    name='default', version='0.0.0-0',
                ),
            )

    def register_gateway(self, network_id, gateway_id, gateway_hw_id):
        """
        Synchronously register a gateway and create default magmad and cellular
        configs for that gateway.

        If the gateway already exists, no action will be performed.

        Args:
            network_id (str): network ID to register on
            gateway_id (str): requested gateway ID
            gateway_hw_id (str): gateway hardware ID
        """
        if gateway_id in self.gateways_api. \
                networks_network_id_gateways_get(network_id):
            return

        hw_id = swagger_client.HwGatewayId(id=gateway_hw_id)
        challenge_key = swagger_client.ChallengeKey(key_type="ECHO")
        gateway_record = swagger_client.AccessGatewayRecord(
            hw_id=hw_id,
            name='integ_test_gw',
            key=challenge_key,
        )
        self.gateways_api.networks_network_id_gateways_post(
            network_id,
            gateway_record,
            requested_id=gateway_id,
        )

        self.gateways_api.networks_network_id_gateways_gateway_id_configs_post(
            network_id, gateway_id,
            DEFAULT_GATEWAY_CONFIG,
        )
        self.gateways_api.networks_network_id_gateways_gateway_id_configs_cellular_post(
            network_id, gateway_id,
            DEFAULT_GATEWAY_CELLULAR_CONFIG,
        )

    def clean_up(self):
        """
        Force-delete all networks that were created via this cloud manager.

        NOTE: this uses an undocumented feature of the cloud API to
        force-destroy a network, regardless of whether it is empty.

        NOTE: for safety, this is currently only allowed on test cloud
        """
        if self._hostname != TEST_CLOUD_HOSTNAME:
            raise ValueError('clean_up only allowed on test cloud')

        self.delete_networks(self._managed_networks)

        # Assert no managed networks exist
        all_networks = set(self.networks_api.networks_get())
        for network_id in self._managed_networks:
            assert network_id not in all_networks
        self._managed_networks = []

    def delete_networks(self, network_ids: List[str]):
        """
        Delete networks and all associated gateways and configs

        Args:
            network_ids: IDs of the networks to delete
        """
        self._clean_up_configs(network_ids)
        for network_id in network_ids:
            self._delete_network(network_id)

    class NetworkConfigType(Enum):
        DNS = 1
        CELLULAR = 2

    def update_network_configs(
        self, network_id: str,
        configs_by_type: Dict[NetworkConfigType, Any],
    ):
        """
        Update multiple configs for a network.

        Args:
            network_id: Network to update
            configs_by_type: New configs to PUT keyed by type.
        """
        for t, value in configs_by_type.items():
            self.update_network_config(network_id, t, value)

    class GatewayConfigType(Enum):
        MAGMAD = 1
        CELLULAR = 2

    def update_gateway_configs(
        self, network_id: str, gateway_id: str,
        configs_by_type: Dict[GatewayConfigType, Any],
    ):
        """
        Update multiple configs for a gateway.

        Args:
            network_id: Network of the gateway
            gateway_id: Gateway to update
            configs_by_type: New configs to PUT keyed by type
        """
        for t, value in configs_by_type.items():
            self.update_gateway_config(network_id, gateway_id, t, value)

    def update_network_config(
        self, network_id: str,
        config_type: NetworkConfigType, value: Any,
    ):
        """
        Update a config for a network.

        Args:
            network_id: Network to update
            config_type: Config type to update (specifies the endpoint)
            value: Config value
        """
        if config_type == CloudManager.NetworkConfigType.DNS:
            self.networks_api.networks_network_id_configs_dns_put(
                network_id, value,
            )
        elif config_type == CloudManager.NetworkConfigType.CELLULAR:
            self.networks_api.networks_network_id_configs_cellular_put(
                network_id, value,
            )
        else:
            raise ValueError(
                'Network config type {} not recognized/supported'.format(
                    config_type,
                ),
            )

    def update_gateway_config(
        self, network_id: str, gateway_id: str,
        config_type: GatewayConfigType, value: Any,
    ):
        """
        Update a config for a gateway.

        Args:
            network_id: Network of gateway
            gateway_id: ID of gateway
            config_type: Config type to update (specifies the endpoint)
            value: Config value
        """
        if config_type == CloudManager.GatewayConfigType.MAGMAD:
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_put(
                network_id, gateway_id, value,
            )
        elif config_type == CloudManager.GatewayConfigType.CELLULAR:
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_cellular_put(
                network_id, gateway_id, value,
            )
        else:
            raise ValueError(
                'Gateway config type {} not recognized/supported'.format(
                    config_type,
                ),
            )

    def _delete_network(self, network_id: str):
        networks = set(self.networks_api.networks_get())
        if network_id not in networks:
            return

        request_url = '{host}/networks/{network}'.format(
            host=self._host, network=network_id,
        )
        resp = requests.delete(
            request_url, cert=(self._ssl_cert, self._ssl_key),
            verify=self._ssl_ca, params={'mode': 'force'},
        )
        resp.raise_for_status()

    def _clean_up_configs(self, network_ids: List[str]):
        # Have to manually delete all configs because we don't hook config
        # service into network force-delete
        gateways_by_network = self._get_managed_gateways_by_network(
            network_ids,
        )
        for network_id, gateway_ids in gateways_by_network.items():
            for gateway_id in gateway_ids:
                self._delete_gateway_configs(network_id, gateway_id)

        for network_id in network_ids:
            self._delete_network_configs(network_id)

    def _delete_gateway_configs(self, network_id: str, gateway_id: str):
        self._delete_config(
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_get,
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_delete,
            network_id,
            gateway_id,
        )
        self._delete_config(
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_cellular_get,
            self.gateways_api.networks_network_id_gateways_gateway_id_configs_cellular_delete,
            network_id,
            gateway_id,
        )

    def _delete_network_configs(self, network_id: str):
        self._delete_config(
            self.networks_api.networks_network_id_configs_dns_get,
            self.networks_api.networks_network_id_configs_dns_delete,
            network_id,
        )
        self._delete_config(
            self.networks_api.networks_network_id_configs_cellular_get,
            self.networks_api.networks_network_id_configs_cellular_delete,
            network_id,
        )

    @staticmethod
    def _delete_config(get, delete, *args):
        try:
            get(*args)
            delete(*args)
        except ApiException as e:
            if e.status != 404:
                raise Exception from e

    def _get_managed_gateways_by_network(
            self,
            network_ids: List[str],
    ) -> Dict[str, List[str]]:
        ret = {}
        for network_id in network_ids:
            ret[network_id] = self.gateways_api.networks_network_id_gateways_get(
                network_id,
            )
        return ret

    def get_gateway_status(self, network_id, gateway_id):
        """ Get gateway status """
        resp = self.gateways_api.networks_network_id_gateways_gateway_id_status_get(
            network_id, gateway_id,
        )
        return resp

    def get_gateway_sw_version(self, network_id, gateway_id):
        """ Get gateway software version """
        resp = self.get_gateway_status(network_id, gateway_id)
        return resp.version

    def get_channel_sw_versions(self, channel):
        """ Get list of available gateway versions """
        resp = self.channels_api.channels_channel_id_get(channel)
        return resp.supported_versions

    def upgrade_tier(self, network, tier, version):
        """ Upgrade SW version for a tier """
        resp = self.tiers_api.networks_network_id_tiers_tier_id_get(
            network, tier,
        )
        resp.version = version
        self.tiers_api.networks_network_id_tiers_tier_id_put(
            network, tier, resp,
        )
