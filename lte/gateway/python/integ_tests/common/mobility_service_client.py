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

import grpc
from integ_tests.gateway.rpc import get_rpc_channel
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
                err_code ==
                grpc.StatusCode.FAILED_PRECONDITION
            ):
                logging.info("Ignoring FAILED_PRECONDITION exception")
            else:
                raise

    def remove_ip_blocks(self, blocks, force=False):
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
                RemoveIPBlockRequest(ip_blocks=ip_blocks, force=force),
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
                err_code ==
                grpc.StatusCode.FAILED_PRECONDITION
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
