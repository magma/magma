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
from typing import Optional

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.redis.client import get_default_client
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.mobilityd.ip_address_man import IPAddressManager
from magma.mobilityd.ip_allocator_base import OverlappedIPBlocksError
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP
from magma.mobilityd.ip_allocator_multi_apn import IPAllocatorMultiAPNWrapper
from magma.mobilityd.ip_allocator_pool import IpAllocatorPool
from magma.mobilityd.ip_allocator_static import IPAllocatorStaticWrapper
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mobility_store import MobilityStore
from magma.mobilityd.rpc_servicer import MobilityServiceRpcServicer

DEFAULT_IPV6_PREFIX_ALLOC_MODE = 'RANDOM'


def _get_ipv4_allocator(store: MobilityStore, allocator_type: int,
                        static_ip_enabled: bool, multi_apn: bool,
                        dhcp_iface: str, dhcp_retry_limit: int,
                        subscriberdb_rpc_stub: SubscriberDBStub = None):
    # Read default GW, this is required for static IP allocation.
    store.dhcp_gw_info.read_default_gw()

    if allocator_type == mconfigs_pb2.MobilityD.IP_POOL:
        ip_allocator = IpAllocatorPool(store)
    elif allocator_type == mconfigs_pb2.MobilityD.DHCP:
        ip_allocator = IPAllocatorDHCP(store=store,
                                       iface=dhcp_iface,
                                       retry_limit=dhcp_retry_limit)
    else:
        raise ValueError(
            "Unknown IP allocator type: %s" % allocator_type)

    if static_ip_enabled:
        ip_allocator = IPAllocatorStaticWrapper(
            store=store, subscriberdb_rpc_stub=subscriberdb_rpc_stub,
            ip_allocator=ip_allocator)

    if multi_apn:
        ip_allocator = IPAllocatorMultiAPNWrapper(store=store,
                                                  subscriberdb_rpc_stub=subscriberdb_rpc_stub,
                                                  ip_allocator=ip_allocator)

    logging.info("Using allocator: %s static ip: %s multi_apn %s",
                 allocator_type,
                 static_ip_enabled,
                 multi_apn)
    return ip_allocator


def _get_ip_block(ip_block_str: str) -> Optional[ipaddress.ip_network]:
    """ Convert string into ipaddress.ip_network. Support both IPv4 or IPv6
    addresses.

        Args:
            ip_block_str(string): network address, e.g. "192.168.0.0/24".

        Returns:
            ip_block(ipaddress.ip_network)
    """
    if not ip_block_str:
        logging.error("Empty IP block")
        return None
    try:
        ip_block = ipaddress.ip_network(ip_block_str)
    except ValueError:
        logging.error("Invalid IP block format: %s", ip_block_str)
        return None
    return ip_block


def main():
    """ main() for MobilityD """
    service = MagmaService('mobilityd', mconfigs_pb2.MobilityD())

    # Optionally pipe errors to Sentry
    sentry_init()

    # Load service configs and mconfig
    config = service.config
    mconfig = service.mconfig

    multi_apn = config.get('multi_apn', mconfig.multi_apn_ip_alloc)
    static_ip_enabled = config.get('static_ip', mconfig.static_ip_enabled)
    allocator_type = mconfig.ip_allocator_type

    dhcp_iface = config.get('dhcp_iface', 'dhcp0')
    dhcp_retry_limit = config.get('retry_limit', 300)

    # TODO: consider adding gateway mconfig to decide whether to
    # persist to Redis
    client = get_default_client()
    store = MobilityStore(client, config.get('persist_to_redis', False),
                          config.get('redis_port', 6380))

    chan = ServiceRegistry.get_rpc_channel('subscriberdb',
                                           ServiceRegistry.LOCAL)
    ipv4_allocator = _get_ipv4_allocator(store, allocator_type,
                                         static_ip_enabled, multi_apn,
                                         dhcp_iface, dhcp_retry_limit,
                                         SubscriberDBStub(chan))

    # Init IPv6 allocator, for now only IP_POOL mode is supported for IPv6
    ipv6_prefix_allocation_type = mconfig.ipv6_prefix_allocation_type or \
                                  DEFAULT_IPV6_PREFIX_ALLOC_MODE
    ipv6_allocator = IPv6AllocatorPool(
        store=store, session_prefix_alloc_mode=ipv6_prefix_allocation_type)

    # Load IPAddressManager
    ip_address_man = IPAddressManager(ipv4_allocator, ipv6_allocator, store)

    # Add all servicers to the server
    mobility_service_servicer = MobilityServiceRpcServicer(
        ip_address_man, config.get('print_grpc_payload', False))
    mobility_service_servicer.add_to_server(service.rpc_server)

    # Load IPv4 and IPv6 blocks from the configurable mconfig file
    # No dynamic reloading support for now, assume restart after updates
    logging.info('Adding IPv4 block')
    ipv4_block = _get_ip_block(mconfig.ip_block)
    if ipv4_block is not None:
        try:
            mobility_service_servicer.add_ip_block(ipv4_block)
        except OverlappedIPBlocksError:
            logging.warning("Overlapped IPv4 block: %s", ipv4_block)

    logging.info('Adding IPv6 block')
    ipv6_block = _get_ip_block(mconfig.ipv6_block)
    if ipv6_block is not None:
        try:
            mobility_service_servicer.add_ip_block(ipv6_block)
        except OverlappedIPBlocksError:
            logging.warning("Overlapped IPv6 block: %s", ipv6_block)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
