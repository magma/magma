"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import ipaddress
import logging

import grpc
from lte.protos.mobilityd_pb2 import AllocateIPRequest, IPAddress, IPBlock, \
    ListAddedIPBlocksResponse, ListAllocatedIPsResponse, RemoveIPBlockResponse, \
    SubscriberIPTable
from lte.protos.mobilityd_pb2_grpc import MobilityServiceServicer, \
    add_MobilityServiceServicer_to_server
from lte.protos.subscriberdb_pb2 import SubscriberID
from magma.common.rpc_utils import return_void
from magma.subscriberdb.sid import SIDUtils

from .ip_address_man import IPAddressManager, IPNotInUseError, MappingNotFoundError
from .ip_allocator_base import IPAllocatorType

from .ip_allocator_static import IPBlockNotFoundError, NoAvailableIPError, \
    OverlappedIPBlocksError

from .ip_allocator_base import DuplicatedIPAllocationError

def _get_ip_block(ip_block_str):
    """ Convert string into ipaddress.ip_network. Support both IPv4 or IPv6
    addresses.

        Args:
            ip_block_str(string): network address, e.g. "192.168.0.0/24".

        Returns:
            ip_block(ipaddress.ip_network)
    """
    try:
        ip_block = ipaddress.ip_network(ip_block_str)
    except ValueError:
        logging.error("Invalid IP block format: %s", ip_block_str)
        return None
    return ip_block


class MobilityServiceRpcServicer(MobilityServiceServicer):
    """ gRPC based server for the IPAllocator. """

    def __init__(self, mconfig, config):
        # TODO: consider adding gateway mconfig to decide whether to
        # persist to Redis
        if config['allocator_type'] == 'ip_pool':
            config_allocator_type = IPAllocatorType.IP_POOL

        self._ipv4_allocator = IPAddressManager(
            persist_to_redis=config['persist_to_redis'],
            redis_port=config['redis_port'],
            allocator_type=config_allocator_type)

        # Load IP block from the configurable mconfig file
        # No dynamic reloading support for now, assume restart after updates
        ip_block = _get_ip_block(mconfig.ip_block)
        if ip_block is not None:
            try:
                self.add_ip_block(ip_block)
            except OverlappedIPBlocksError:
                logging.error("Overlapped IP block: %s", ip_block)
            except IPVersionNotSupportedError:
                logging.error(
                    "IP version not supported for IP block: %s", ip_block)

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_MobilityServiceServicer_to_server(self, server)

    def add_ip_block(self, ip_block):
        """ Add IP block to the IP allocator. Currently, only IPv4 is supported.

            Args:
                ipblock (ipaddress.ip_network): ip network to add
                e.g. ipaddress.ip_network("10.0.0.0/24")

            Raise:
                OverlappedIPBlocksError: if the given IP block overlaps with
                existing ones
                IPVersionNotSupportedError: if given IP version of the IP block
                is not supported
        """
        if ip_block.version == 4:
            self._ipv4_allocator.add_ip_block(ip_block)
            logging.info("Added block %s to the IPv4 address pool", ip_block)
        else:
            raise IPVersionNotSupportedError

    @return_void
    def AddIPBlock(self, ipblock_msg, context):
        """ Add a range of IP addresses into the free IP pool

        Args:
            ipblock_msg (IPBlock): ip block to add. ipblock_msg has the
            type IPBlock, a protobuf message type for the gRPC interface.
            Internal representation of ip blocks uses the ipaddress.ip_network
            type and is named as ipblock.
        """
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
        resp = ListAddedIPBlocksResponse()

        ip_blocks = self._ipv4_allocator.list_added_ip_blocks()
        ip_block_msg_list = [IPBlock(version=IPAddress.IPV4,
                                     net_address=block.network_address.packed,
                                     prefix_len=block.prefixlen)
                             for block in ip_blocks]
        resp.ip_block_list.extend(ip_block_msg_list)

        return resp

    def ListAllocatedIPs(self, ipblock_msg, context):
        """ Return a list of IPs allocated from a IP block

        Args:
            ipblock_msg (IPBlock): ip block to add. ipblock_msg has the
            type IPBlock, a protobuf message type for the gRPC interface.
            Internal representation of ip blocks uses the ipaddress.ip_network
            type and is named as ipblock.
        """
        resp = ListAllocatedIPsResponse()

        ipblock = self._ipblock_msg_to_ipblock(ipblock_msg, context)
        if ipblock is None:
            return resp

        if ipblock_msg.version == IPBlock.IPV4:
            try:
                ips = self._ipv4_allocator.list_allocated_ips(ipblock)
                ip_msg_list = [IPAddress(version=IPAddress.IPV4,
                                         address=ip.packed) for ip in ips]

                resp.ip_list.extend(ip_msg_list)
            except IPBlockNotFoundError:
                context.set_details('IP block not found: %s' % ipblock)
                context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
        else:
            self._unimplemented_ip_version_error(context)

        return resp

    def AllocateIPAddress(self, request, context):
        """ Allocate an IP address from the free IP pool """
        resp = IPAddress()
        if request.version == AllocateIPRequest.IPV4:
            try:
                composite_sid = SIDUtils.to_str(request.sid)
                if request.apn:
                    composite_sid = composite_sid + "." + request.apn

                ip = self._ipv4_allocator.alloc_ip_address(composite_sid)
                logging.info("Allocated IPv4 %s for sid %s for apn %s"
                             % (ip, SIDUtils.to_str(request.sid), request.apn))
                resp.version = IPAddress.IPV4
                resp.address = ip.packed
            except NoAvailableIPError:
                context.set_details('No free IPv4 IP available')
                context.set_code(grpc.StatusCode.RESOURCE_EXHAUSTED)
            except DuplicatedIPAllocationError:
                context.set_details('IP has been allocated for this subscriber')
                context.set_code(grpc.StatusCode.ALREADY_EXISTS)
        else:
            self._unimplemented_ip_version_error(context)
        return resp

    @return_void
    def ReleaseIPAddress(self, request, context):
        """ Release an allocated IP address """
        if request.ip.version == IPAddress.IPV4:
            try:
                ip = ipaddress.ip_address(request.ip.address)
                composite_sid = SIDUtils.to_str(request.sid)
                if request.apn:
                    composite_sid = composite_sid + "." + request.apn
                self._ipv4_allocator.release_ip_address(
                    composite_sid, ip)
                logging.info("Released IPv4 %s for sid %s"
                             % (ip, SIDUtils.to_str(request.sid)))
            except IPNotInUseError:
                context.set_details('IP %s not in use' % ip)
                context.set_code(grpc.StatusCode.NOT_FOUND)
            except MappingNotFoundError:
                context.set_details('(SID, IP) map not found: (%s, %s)'
                                    % (SIDUtils.to_str(request.sid), ip))
                context.set_code(grpc.StatusCode.NOT_FOUND)
        else:
            self._unimplemented_ip_version_error(context)

    def RemoveIPBlock(self, request, context):
        """ Attempt to remove IP blocks and return the removed blocks """
        removed_blocks = self._ipv4_allocator.remove_ip_blocks(
            *[self._ipblock_msg_to_ipblock(ipblock_msg, context)
                for ipblock_msg in request.ip_blocks],
            force=request.force)
        removed_block_msgs = [IPBlock(version=IPAddress.IPV4,
                                      net_address=block.network_address.packed,
                                      prefix_len=block.prefixlen)
                                      for block in removed_blocks]

        resp = RemoveIPBlockResponse()
        resp.ip_blocks.extend(removed_block_msgs)
        return resp

    def GetIPForSubscriber(self, request, context):
        composite_sid = SIDUtils.to_str(request.sid)
        if request.apn:
            composite_sid = composite_sid + "." + request.apn

        ip = self._ipv4_allocator.get_ip_for_sid(composite_sid)
        if ip is None:
            context.set_details('SID %s not found'
                                % SIDUtils.to_str(request.sid))
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return IPAddress()

        version = IPAddress.IPV4 if ip.version == 4 else IPAddress.IPV6
        return IPAddress(version=version, address=ip.packed)

    def GetSubscriberIDFromIP(self, ip_addr, context):
        sent_ip = ipaddress.ip_address(ip_addr.address)
        sid = self._ipv4_allocator.get_sid_for_ip(sent_ip)

        if sid is None:
            context.set_details('IP address %s not found' % str(sent_ip))
            context.set_code(grpc.StatusCode.NOT_FOUND)
            return SubscriberID()
        else:
            #handle composite key case
            sid, *rest = sid.partition('.')
            return SIDUtils.to_pb(sid)

    def GetSubscriberIPTable(self, void, context):
        """ Get the full subscriber table """
        logging.debug("Listing subscriber IP table")
        resp = SubscriberIPTable()

        csid_ip_pairs = self._ipv4_allocator.get_sid_ip_table()
        for composite_sid, ip in csid_ip_pairs:
            #handle composite sid to sid and apn mapping
            sid, _, apn = composite_sid.partition('.')
            sid_pb = SIDUtils.to_pb(sid)
            version = IPAddress.IPV4 if ip.version == 4 else IPAddress.IPV6
            ip_msg = IPAddress(version=version, address=ip.packed)
            resp.entries.add(sid=sid_pb, ip=ip_msg, apn=apn)
        return resp

    def _ipblock_msg_to_ipblock(self, ipblock_msg, context):
        """ convert IPBlock to ipaddress.ip_network """
        try:
            ip = ipaddress.ip_address(ipblock_msg.net_address)
            subnet = "%s/%s" % (str(ip), ipblock_msg.prefix_len)
            ipblock = ipaddress.ip_network(subnet)
            return ipblock
        except ValueError:
            context.set_details('Invalid IPBlock format: version=%s,'
                                'net_address=%s, prefix_len=%s' %
                                (ipblock_msg.version, ipblock_msg.net_address,
                                 ipblock_msg.prefix_len))
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return None

    def _unimplemented_ip_version_error(self, context):
        context.set_details("IPv6 is not yet supported")
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)


class IPVersionNotSupportedError(Exception):
    """ Exception thrown when an IP version is not supported """
    pass
