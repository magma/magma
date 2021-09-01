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
from typing import Any, Optional

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
RETRY_LIMIT = 300
DEFAULT_REDIS_PORT = 6380


def _get_ipv4_allocator(
    store: MobilityStore, allocator_type: int,
    static_ip_enabled: bool, multi_apn: bool,
    dhcp_iface: str, dhcp_retry_limit: int,
    subscriberdb_rpc_stub: SubscriberDBStub = None,
):
    # Read default GW, this is required for static IP allocation.
    store.dhcp_gw_info.read_default_gw()

    if allocator_type == mconfigs_pb2.MobilityD.IP_POOL:
        ip_allocator = IpAllocatorPool(store)
    elif allocator_type == mconfigs_pb2.MobilityD.DHCP:
        ip_allocator = IPAllocatorDHCP(
            store=store,
            iface=dhcp_iface,
            retry_limit=dhcp_retry_limit,
        )
    else:
        raise ValueError(
            "Unknown IP allocator type: %s" % allocator_type,
        )

    if static_ip_enabled:
        ip_allocator = IPAllocatorStaticWrapper(
            store=store, subscriberdb_rpc_stub=subscriberdb_rpc_stub,
            ip_allocator=ip_allocator,
        )

    if multi_apn:
        ip_allocator = IPAllocatorMultiAPNWrapper(
            store=store,
            subscriberdb_rpc_stub=subscriberdb_rpc_stub,
            ip_allocator=ip_allocator,
        )

    logging.info(
        "Using allocator: %s static ip: %s multi_apn %s",
        allocator_type,
        static_ip_enabled,
        multi_apn,
    )
    return ip_allocator


def _get_ip_block(
    ip_block_str: str,
    ip_type: str,
) -> Optional[ipaddress.ip_network]:
    """Convert string into ipaddress.ip_network

    Support both IPv4 or IPv6 addresses

    Args:
        ip_block_str (str): network address, e.g. "192.168.0.0/24"
        ip_type (str): ipv4 or ipv6 used for logging

    Returns:
        Optional[ipaddress.ip_network]:
    """
    if not ip_block_str:
        logging.warning(
            "%s ip block is not specified in mconfig, skipping", ip_type,
        )
        return None
    try:
        ip_block = ipaddress.ip_network(ip_block_str)
    except ValueError:
        logging.error("Invalid IP block format: %s", ip_block_str)
        return None
    return ip_block


def _get_value_or_default(val: Any, default: Any) -> Any:
    return val or default


def main():
    """Start mobilityd"""
    service = MagmaService('mobilityd', mconfigs_pb2.MobilityD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name)

    # Load service configs and mconfig
    config = service.config
    mconfig = service.mconfig

    multi_apn = config.get('multi_apn', mconfig.multi_apn_ip_alloc)
    static_ip_enabled = config.get('static_ip', mconfig.static_ip_enabled)
    allocator_type = mconfig.ip_allocator_type

    dhcp_iface = config.get('dhcp_iface', 'dhcp0')
    dhcp_retry_limit = config.get('retry_limit', RETRY_LIMIT)

    # TODO: consider adding gateway mconfig to decide whether to
    # persist to Redis
    client = get_default_client()
    store = MobilityStore(
        client, config.get('persist_to_redis', False),
        config.get('redis_port', DEFAULT_REDIS_PORT),
    )

    chan = ServiceRegistry.get_rpc_channel(
        'subscriberdb',
        ServiceRegistry.LOCAL,
    )
    ipv4_allocator = _get_ipv4_allocator(
        store, allocator_type,
        static_ip_enabled, multi_apn,
        dhcp_iface, dhcp_retry_limit,
        SubscriberDBStub(chan),
    )

    # Init IPv6 allocator, for now only IP_POOL mode is supported for IPv6
    ipv6_allocator = IPv6AllocatorPool(
        store=store,
        session_prefix_alloc_mode=_get_value_or_default(
            mconfig.ipv6_prefix_allocation_type,
            DEFAULT_IPV6_PREFIX_ALLOC_MODE,
        ),
    )

    # Load IPAddressManager
    ip_address_man = IPAddressManager(ipv4_allocator, ipv6_allocator, store)

    # Load IPv4 and IPv6 blocks from the configurable mconfig file
    # No dynamic reloading support for now, assume restart after updates
    ipv4_block = _get_ip_block(mconfig.ip_block, "IPv4")
    if ipv4_block is not None:
        logging.info('Adding IPv4 block')
        try:
            allocated_ip_blocks = ip_address_man.list_added_ip_blocks()
            if ipv4_block not in allocated_ip_blocks:
                # Cleanup previously allocated IP blocks
                ip_address_man.remove_ip_blocks(*allocated_ip_blocks, force=True)
                ip_address_man.add_ip_block(ipv4_block)
        except OverlappedIPBlocksError:
            logging.warning("Overlapped IPv4 block: %s", ipv4_block)

    ipv6_block = _get_ip_block(mconfig.ipv6_block, "IPv6")
    if ipv6_block is not None:
        logging.info('Adding IPv6 block')
        try:
            allocated_ipv6_block = ip_address_man.get_assigned_ipv6_block()
            if ipv6_block != allocated_ipv6_block:
                ip_address_man.add_ip_block(ipv6_block)
        except OverlappedIPBlocksError:
            logging.warning("Overlapped IPv6 block: %s", ipv6_block)

    # Add all servicers to the server
    mobility_service_servicer = MobilityServiceRpcServicer(
        ip_address_man, config.get('print_grpc_payload', False),
    )
    mobility_service_servicer.add_to_server(service.rpc_server)
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
