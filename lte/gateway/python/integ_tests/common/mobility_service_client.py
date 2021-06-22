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

import abc
import ipaddress
import logging
import time

import grpc
# from integ_tests.cloud.fixtures import GATEWAY_ID, NETWORK_ID
from integ_tests.gateway.rpc import get_gateway_hw_id, get_rpc_channel
from lte.protos.mobilityd_pb2 import IPAddress, IPBlock, RemoveIPBlockRequest
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from orc8r.protos.common_pb2 import Void


class MobilityServiceClient(metaclass=abc.ABCMeta):
    """ Interface for Mobility services. """

    @abc.abstractmethod
    def add_ip_block(self, block):
        """
        Add an ip block.

        Args:
            ip_block (ipaddress.ip_network): IP block to add
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def list_added_blocks(self):
        """
        List all IP blocks in mobilityd. Is blocking.

        Returns:
            blocks (ipaddress.ip_network[]): all added IP blocks
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def remove_ip_blocks(self, blocks):
        """
        Attempt to remove :blocks: from mobilityd.

        Args:
            blocks (ipaddress.ip_network[]): IP blocks to remove
        Returns:
            removed_blocks (ipaddress.ip_network[]): IP blocks
                successfully removed
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def get_subscriber_ip_table(self):
        """
        Retrieve full subscriber table from mobilityd. Is blocking.

        Note: this method will be deprecated once s1aptester exposes
        the IP of the UE.
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def remove_all_ip_blocks(self):
        """ Remove all allocated IP blocks. Is blocking. """
        raise NotImplementedError()

    @abc.abstractmethod
    def wait_for_changes(self):
        raise NotImplementedError()


class MobilityServiceGrpc(MobilityServiceClient):
    """
    Handle mobility actions by making calls over gRPC directly to the
    gateway.
    """

    def __init__(self):
        """ Init the gRPC stub. """
        self._mobility_stub = MobilityServiceStub(get_rpc_channel("mobilityd"))

    @staticmethod
    def _get_ip_data(block):
        """
        Construct an IPBlock message from a given IP block.

        Args:
            block (ipaddress.ip_network): the IP block to embed
        """
        ipblock_msg = IPBlock()
        ipblock_msg.version = IPBlock.IPV4
        ipblock_msg.net_address = block.network_address.packed
        ipblock_msg.prefix_len = block.prefixlen
        return ipblock_msg

    def add_ip_block(self, block):
        mobility_data = self._get_ip_data(block)
        try:
            self._mobility_stub.AddIPBlock(mobility_data)
        except grpc.RpcError as error:
            err_code = error.exception().code()
            if err_code == grpc.StatusCode.FAILED_PRECONDITION:
                logging.info("Ignoring FAILED_PRECONDITION exception")
            else:
                raise

    def list_added_blocks(self):
        try:
            response = self._mobility_stub.ListAddedIPv4Blocks(Void())
            ip_block_list = []
            for block in response.ip_block_list:
                address_bytes = block.net_address
                address_int = int.from_bytes(address_bytes, byteorder='big')
                address = ipaddress.ip_address(address_int)
                ip_block_list.append(
                    ipaddress.ip_network(
                        "%s/%d" % (address, block.prefix_len),
                    ),
                )
            if ip_block_list is not None:
                ip_block_list.sort()
            return ip_block_list
        except grpc.RpcError as error:
            err_code = error.exception().code()
            if (
                err_code
                == grpc.StatusCode.FAILED_PRECONDITION
            ):
                logging.info("Ignoring FAILED_PRECONDITION exception")
            else:
                raise

    def remove_ip_blocks(self, blocks):
        try:
            ip_blocks = [
                IPBlock(
                    version={
                        4: IPAddress.IPV4,
                        6: IPAddress.IPV6,
                    }[block.version],
                    net_address=block.network_address.packed,
                    prefix_len=block.prefixlen,
                )
                for block in blocks
            ]
            response = self._mobility_stub.RemoveIPBlock(
                RemoveIPBlockRequest(ip_blocks=ip_blocks, force=False),
            )
            removed_ip_block_list = ()
            for block in response.ip_blocks:
                address_bytes = block.net_address
                address_int = int.from_bytes(address_bytes, byteorder='big')
                address = ipaddress.ip_address(address_int)
                removed_ip_block_list += (
                    ipaddress.ip_network(
                        "%s/%d" % (address, block.prefix_len),
                    ),
                )
            return removed_ip_block_list
        except grpc.RpcError as error:
            err_code = error.exception().code()
            if (
                err_code
                == grpc.StatusCode.FAILED_PRECONDITION
            ):
                logging.info("Ignoring FAILED_PRECONDITION exception")
            else:
                raise

    def get_subscriber_ip_table(self):
        response = self._mobility_stub.GetSubscriberIPTable(Void())
        table = {}
        for entry in response.entries:
            sid = entry.sid.id
            ip = ipaddress.ip_address(entry.ip.address)
            table[sid] = ip
        return table

    def remove_all_ip_blocks(self):
        blocks = self.list_added_blocks()
        self.remove_ip_blocks(blocks)

    def wait_for_changes(self):
        """ All changes propagate immediately, no need to wait """
        return


# class MobilityServiceRest(MobilityServiceClient):
#    """
#    Handle mobility actions by making REST calls to the cloud.
#
#    Note that, while the gRPC API exposes support for multiple
#    mobilityd IP blocks, the REST API does not. It supports
#    only a single IP block. This class fails an assertion on attempts
#    to add more than 1 block at a time.
#    """
#
#    def __init__(self, cloud_manager, gateway_services):
#        """ Init the REST endpoint. """
#        self._cloud_manager = cloud_manager
#        self._gateway_services = gateway_services
#
#        self._cloud_manager.create_network(NETWORK_ID)
#        self._cloud_manager.register_gateway(
#            NETWORK_ID,
#            GATEWAY_ID,
#            get_gateway_hw_id(),
#        )
#
#        # Create mobility gRPC stub to communicate directly with the gateway
#        self._mobility_grpc = MobilityServiceGrpc()
#
#        # Used for knowing when the gateway should restart after
#        # Initialize to 0 so no restart is required by default
#        self._last_change_time = 0
#
#    def _get_gateway_config_from_cloud(self):
#        """ Get the gateway config record. """
#        gateway_config = self._cloud_manager \
#            .gateways_api \
#            .networks_network_id_gateways_gateway_id_configs_cellular_get(
#                NETWORK_ID, GATEWAY_ID)
#        return gateway_config
#
#    def _update_mobilityd_ip_block(self, block):
#        """
#        Pull the full config record, update its ip_block, then push
#        the updated record.
#
#        Args:
#            block (ipaddress.ip_network): address to update with
#        """
#        gateway_config = self._get_gateway_config_from_cloud()
#
#        old_block = gateway_config.epc.ip_block
#        new_block = block.with_prefixlen
#
#        if old_block != new_block:
#            gateway_config.epc.ip_block = new_block
#            self._cloud_manager.gateways_api. \
#                networks_network_id_gateways_gateway_id_configs_cellular_put(
#                    NETWORK_ID, GATEWAY_ID, gateway_config)
#
#    def _delete_block(self):
#        """
#        Delete an IP block via the cloud REST API.
#
#        TODO: not sure how to delete yet, see T19922441.
#        """
#        deleted_ip_placeholder = ipaddress.ip_network('0.0.0.0/32')
#        self._update_mobilityd_ip_block(deleted_ip_placeholder)
#
#    def _cloud_equal_gateway_ips(self):
#        """
#        Check that cloud IP networks match gateway IP networks.
#
#        HACK: since we don't know how to delete ip blocks from the cloud yet,
#        we can treat the following as equal:
#        cloud: [0.0.0.0/32]
#        gateway: []
#
#        Returns:
#            equal (bool): True if IP networks are equivalent; False otherwise
#        """
#        cloud_ips = self.list_added_blocks()
#        gateway_ips = self._mobility_grpc.list_added_blocks()
#        cloud_ips_set = self._ip_block_list_to_set(cloud_ips)
#        gateway_ips_set = self._ip_block_list_to_set(gateway_ips)
#
#        return cloud_ips_set == gateway_ips_set
#
#    def _ip_block_list_to_set(self, ip_block_list):
#        ip_set = {ip.with_prefixlen for ip in ip_block_list}
#        ip_set.discard('0.0.0.0/32')
#        return ip_set
#
#    def add_ip_block(self, block):
#        # to update gateway, remove old ip blocks and then add new one
#        self._last_change_time = time.time()
#        self._update_mobilityd_ip_block(block)
#
#    def list_added_blocks(self):
#        ip_block = self._get_gateway_config_from_cloud().epc.ip_block
#        network = ipaddress.ip_network(ip_block)
#        return [network]
#
#    def remove_ip_blocks(self, blocks):
#        removed_blocks = []
#        for block in self.list_added_blocks():
#            if block in blocks:
#                self._last_change_time = time.time()
#                self._delete_block()
#                removed_blocks.append(block)
#        return [removed_blocks]
#
#    def get_subscriber_ip_table(self):
#        # HACK: IP tables intentionally not exposed over the REST API
#        ip_table = self._mobility_grpc.get_subscriber_ip_table()
#        return ip_table
#
#    def remove_all_ip_blocks(self):
#        blocks = self.list_added_blocks()
#        self.remove_ip_blocks(blocks)
#
#    def wait_for_changes(self):
#        if self._cloud_equal_gateway_ips():
#            # Shouldn't wait, because gateway won't restart
#            return
#        print("Added IP block to cloud, waiting for gateway restart")
#        self._gateway_services.wait_for_healthy_gateway(
#            after_start_time=self._last_change_time)
#        print("Gateway healthy after restart")
#        assert(self._cloud_equal_gateway_ips())
#
