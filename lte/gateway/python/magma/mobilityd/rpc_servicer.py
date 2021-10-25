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

import ipaddress
import logging

import grpc
from google.protobuf.json_format import MessageToJson
from lte.protos.mobilityd_pb2 import (
    AllocateIPAddressResponse,
    AllocateIPRequest,
    IPAddress,
    IPBlock,
    ListAddedIPBlocksResponse,
    ListAllocatedIPsResponse,
    ListGWInfoResponse,
    RemoveIPBlockResponse,
    SubscriberIPTable,
)
from lte.protos.mobilityd_pb2_grpc import (
    MobilityServiceServicer,
    add_MobilityServiceServicer_to_server,
)
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.common.rpc_utils import return_void
from magma.subscriberdb.sid import SIDUtils

from .ip_address_man import (
    IPAddressManager,
    IPNotInUseError,
    MappingNotFoundError,
)
from .ip_allocator_base import (
    DuplicatedIPAllocationError,
    DuplicateIPAssignmentError,
    IPBlockNotFoundError,
    NoAvailableIPError,
    OverlappedIPBlocksError,
)
from .ipv6_allocator_pool import MaxCalculationError
from .subscriberdb_client import (
    SubscriberDBConnectionError,
    SubscriberDBMultiAPNValueError,
    SubscriberDBStaticIPValueError,
)


class MobilityServiceRpcServicer(MobilityServiceServicer):
    """ gRPC based server for the IPAllocator. """

    def __init__(
        self, ip_address_manager: IPAddressManager,
        print_grpc_payload: bool = False,
    ):
        """Initialize mobilityd GRPC endpoints."""
        self._ip_address_man = ip_address_manager
        self._print_grpc_payload = print_grpc_payload

        if self._print_grpc_payload:
            logging.info("Printing GRPC messages")

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_MobilityServiceServicer_to_server(self, server)

    def add_ip_block(self, ip_block):
        """ Add IP block to the IP allocator.

            Args:
                ip_block (ipaddress.ip_network): ip network to add
                e.g. ipaddress.ip_network("10.0.0.0/24")

            Raise:
                OverlappedIPBlocksError: if the given IP block overlaps with
                existing ones
                IPVersionNotSupportedError: if given IP version of the IP block
                is not supported
        """
        self._ip_address_man.add_ip_block(ip_block)

    @return_void
    def AddIPBlock(self, ipblock_msg, context):
        """ Add a range of IP addresses into the free IP pool

        Args:
            ipblock_msg (IPBlock): ip block to add. ipblock_msg has the
            type IPBlock, a protobuf message type for the gRPC interface.
            Internal representation of ip blocks uses the ipaddress.ip_network
            type and is named as ipblock.
        """
        logging.debug("Received AddIPBlock")
        self._print_grpc(ipblock_msg)
        ipblock = self._ipblock_msg_to_ipblock(ipblock_msg, context)
        if ipblock is None:
            return

        try:
            self.add_ip_block(ipblock)
        except OverlappedIPBlocksError:
            context.set_details('Overlapped ip block: %s' % ipblock)
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        except IPVersionNotSupportedError:
            self._unimplemented_ip_version_error(context)

    def ListAddedIPv4Blocks(self, void, context):
        """ Return a list of IPv4 blocks assigned """
        logging.debug("Received ListAddedIPv4Blocks")
        self._print_grpc(void)
        resp = ListAddedIPBlocksResponse()

        ip_blocks = self._ip_address_man.list_added_ip_blocks()
        ip_block_msg_list = [
            IPBlock(
                version=IPAddress.IPV4,
                net_address=block.network_address.packed,
                prefix_len=block.prefixlen,
            )
            for block in ip_blocks
        ]
        resp.ip_block_list.extend(ip_block_msg_list)

        self._print_grpc(resp)
        return resp

    def ListAllocatedIPs(self, ipblock_msg, context):
        """ Return a list of IPs allocated from a IP block

        Args:
            ipblock_msg (IPBlock): ip block to add. ipblock_msg has the
            type IPBlock, a protobuf message type for the gRPC interface.
            Internal representation of ip blocks uses the ipaddress.ip_network
            type and is named as ipblock.
        """
        logging.debug("Received ListAllocatedIPs")
        self._print_grpc(ipblock_msg)
        resp = ListAllocatedIPsResponse()

        ipblock = self._ipblock_msg_to_ipblock(ipblock_msg, context)
        if ipblock is None:
            self._print_grpc(resp)
            return resp

        if ipblock_msg.version == IPBlock.IPV4:
            try:
                ips = self._ip_address_man.list_allocated_ips(ipblock)
                ip_msg_list = [
                    IPAddress(
                        version=IPAddress.IPV4,
                        address=ip.packed,
                    ) for ip in ips
                ]

                resp.ip_list.extend(ip_msg_list)
            except IPBlockNotFoundError:
                context.set_details('IP block not found: %s' % ipblock)
                context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        else:
            self._unimplemented_ip_version_error(context)

        self._print_grpc(resp)
        return resp

    def AllocateIPAddress(self, request, context):
        """ Allocate an IP address from the free IP pool """
        logging.debug("Received AllocateIPAddress")
        self._print_grpc(request)
        composite_sid = SIDUtils.to_str(request.sid)
        if request.apn:
            composite_sid = composite_sid + "." + request.apn

        if request.version == AllocateIPRequest.IPV4:
            resp = self._get_allocate_ip_response(
                composite_sid + ",ipv4",
                IPAddress.IPV4, context,
                request,
            )
        elif request.version == AllocateIPRequest.IPV6:
            resp = self._get_allocate_ip_response(
                composite_sid + ",ipv6",
                IPAddress.IPV6, context,
                request,
            )
        elif request.version == AllocateIPRequest.IPV4V6:
            ipv4_response = self._get_allocate_ip_response(
                composite_sid + ",ipv4", IPAddress.IPV4,
                context, request,
            )
            ipv6_response = self._get_allocate_ip_response(
                composite_sid + ",ipv6", IPAddress.IPV6,
                context, request,
            )
            ipv4_addr = ipv4_response.ip_list[0]
            ipv6_addr = ipv6_response.ip_list[0]
            # Get vlan from IPv4 Allocate response
            resp = AllocateIPAddressResponse(
                ip_list=[ipv4_addr, ipv6_addr],
                vlan=ipv4_response.vlan,
            )
        else:
            resp = AllocateIPAddressResponse()

        self._print_grpc(resp)
        return resp

    @return_void
    def ReleaseIPAddress(self, request, context):
        """ Release an allocated IP address """
        logging.debug("Received ReleaseIPAddress")
        self._print_grpc(request)

        ip = ipaddress.ip_address(request.ip.address)
        composite_sid = SIDUtils.to_str(request.sid)
        if request.apn:
            composite_sid = composite_sid + "." + request.apn

        if request.ip.version == IPAddress.IPV4:
            composite_sid = composite_sid + ",ipv4"
        elif request.ip.version == IPAddress.IPV6:
            composite_sid = composite_sid + ",ipv6"

        try:
            self._ip_address_man.release_ip_address(
                composite_sid, ip,
                request.ip.version,
            )
            logging.info(
                "Released IP %s for sid %s"
                % (ip, SIDUtils.to_str(request.sid)),
            )
        except IPNotInUseError:
            context.set_details('IP %s not in use' % ip)
            context.set_code(grpc.StatusCode.NOT_FOUND)
        except MappingNotFoundError:
            context.set_details(
                '(SID, IP) map not found: (%s, %s)'
                % (SIDUtils.to_str(request.sid), ip),
            )
            context.set_code(grpc.StatusCode.NOT_FOUND)

    def RemoveIPBlock(self, request, context):
        """ Attempt to remove IP blocks and return the removed blocks """
        logging.debug("Received RemoveIPBlock")
        self._print_grpc(request)

        removed_blocks = self._ip_address_man.remove_ip_blocks(
            *[
                self._ipblock_msg_to_ipblock(ipblock_msg, context)
                for ipblock_msg in request.ip_blocks
            ],
            force=request.force,
        )

        removed_block_msgs = []
        for block in removed_blocks:
            if block.version == 4:
                removed_block_msgs.append(
                    IPBlock(
                        version=IPAddress.IPV4,
                        net_address=block.network_address.packed,
                        prefix_len=block.prefixlen,
                    ),
                )
            elif block.version == 6:
                removed_block_msgs.append(
                    IPBlock(
                        version=IPAddress.IPV6,
                        net_address=block.network_address.packed,
                        prefix_len=block.prefixlen,
                    ),
                )

        resp = RemoveIPBlockResponse()
        resp.ip_blocks.extend(removed_block_msgs)
        self._print_grpc(resp)
        return resp

    def GetIPForSubscriber(self, request, context):
        logging.debug("Received GetIPForSubscriber")
        self._print_grpc(request)

        composite_sid = SIDUtils.to_str(request.sid)
        if request.apn:
            composite_sid = composite_sid + "." + request.apn

        if request.version == IPAddress.IPV4:
            composite_sid += ",ipv4"
        elif request.version == IPAddress.IPV6:
            composite_sid += ",ipv6"

        ip = self._ip_address_man.get_ip_for_sid(composite_sid)
        if ip is None:
            context.set_details(
                'SID %s not found'
                % SIDUtils.to_str(request.sid),
            )
            context.set_code(grpc.StatusCode.NOT_FOUND)
            resp = IPAddress()
        else:
            version = IPAddress.IPV4 if ip.version == 4 else IPAddress.IPV6
            resp = IPAddress(version=version, address=ip.packed)
        self._print_grpc(resp)
        return resp

    def GetSubscriberIDFromIP(self, ip_addr, context):
        logging.debug("Received GetSubscriberIDFromIP")
        self._print_grpc(ip_addr)

        sent_ip = ipaddress.ip_address(ip_addr.address)
        sid = self._ip_address_man.get_sid_for_ip(sent_ip)

        if sid is None:
            context.set_details('IP address %s not found' % str(sent_ip))
            context.set_code(grpc.StatusCode.NOT_FOUND)
            resp = SubscriberID()
        else:
            # handle composite key case
            sid, *rest = sid.partition('.')
            resp = SIDUtils.to_pb(sid)
        self._print_grpc(resp)
        return resp

    def GetSubscriberIPTable(self, void, context):
        """ Get the full subscriber table """
        logging.debug("Received GetSubscriberIPTable")
        self._print_grpc(void)

        resp = SubscriberIPTable()

        csid_ip_pairs = self._ip_address_man.get_sid_ip_table()
        for composite_sid, ip in csid_ip_pairs:
            # handle composite sid to sid and apn mapping
            sid, _, apn_part = composite_sid.partition('.')
            apn, _ = apn_part.split(',')
            sid_pb = SIDUtils.to_pb(sid)
            version = IPAddress.IPV4 if ip.version == 4 else IPAddress.IPV6
            ip_msg = IPAddress(version=version, address=ip.packed)
            resp.entries.add(sid=sid_pb, ip=ip_msg, apn=apn)
        self._print_grpc(resp)
        return resp

    def ListGatewayInfo(self, void, context):
        logging.debug("Received ListGatewayInfo")
        self._print_grpc(void)

        resp = ListGWInfoResponse()
        gw_info_list = self._ip_address_man.list_gateway_info()
        if gw_info_list:
            resp.gw_list.extend(gw_info_list)
        self._print_grpc(resp)
        return resp

    @return_void
    def SetGatewayInfo(self, info, context):
        logging.debug("Received SetGatewayInfo")
        self._print_grpc(info)
        self._ip_address_man.set_gateway_info(info)

    def _get_allocate_ip_response(
        self, composite_sid, version, context,
        request,
    ):
        try:
            ip, vlan = self._ip_address_man.alloc_ip_address(
                composite_sid,
                version,
            )
            logging.info(
                "Allocated IP %s for sid %s for apn %s"
                % (ip, SIDUtils.to_str(request.sid), request.apn),
            )
            ip_addr = IPAddress(address=ip.packed, version=version)
            return AllocateIPAddressResponse(
                ip_list=[ip_addr],
                vlan=str(vlan),
            )
        except NoAvailableIPError:
            context.set_details('No free IP available')
            context.set_code(grpc.StatusCode.RESOURCE_EXHAUSTED)
        except DuplicatedIPAllocationError:
            context.set_details(
                'IP has been allocated for this subscriber',
            )
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
        except DuplicateIPAssignmentError:
            context.set_details(
                'IP has been allocated for other subscriber',
            )
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
        except MaxCalculationError:
            context.set_details(
                'Reached maximum IPv6 calculation tries',
            )
            context.set_code(grpc.StatusCode.RESOURCE_EXHAUSTED)
        except SubscriberDBConnectionError:
            context.set_details(
                'Could not connect to SubscriberDB',
            )
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        except SubscriberDBStaticIPValueError:
            context.set_details(
                'Could not parse static IP response from SubscriberDB',
            )
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        except SubscriberDBMultiAPNValueError:
            context.set_details(
                'Could not parse MultiAPN IP response from SubscriberDB',
            )
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        return AllocateIPAddressResponse()

    def _ipblock_msg_to_ipblock(self, ipblock_msg, context):
        """ convert IPBlock to ipaddress.ip_network """
        try:
            ip = ipaddress.ip_address(ipblock_msg.net_address)
            subnet = "%s/%s" % (str(ip), ipblock_msg.prefix_len)
            ipblock = ipaddress.ip_network(subnet)
            return ipblock
        except ValueError:
            context.set_details(
                'Invalid IPBlock format: version=%s,'
                'net_address=%s, prefix_len=%s' %
                (
                    ipblock_msg.version, ipblock_msg.net_address,
                    ipblock_msg.prefix_len,
                ),
            )
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return None

    def _unimplemented_ip_version_error(self, context):
        context.set_details("IPv6 is not yet supported")
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)

    def _print_grpc(self, message):
        if self._print_grpc_payload:
            log_msg = "{} {}".format(
                message.DESCRIPTOR.full_name,
                MessageToJson(message),
            )
            # add indentation
            padding = 2 * ' '
            log_msg = ''.join(
                "{}{}".format(padding, line)
                for line in log_msg.splitlines(True)
            )

            log_msg = "GRPC message:\n{}".format(log_msg)
            logging.info(log_msg)


class IPVersionNotSupportedError(Exception):
    """ Exception thrown when an IP version is not supported """
    pass


class UnknownIPAllocatorError(Exception):
    """ Exception thrown when an IP version is not supported """
    pass
