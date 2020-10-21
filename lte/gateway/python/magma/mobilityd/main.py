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
import logging

from magma.common.redis.client import get_default_client
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.mobilityd.ip_address_man import IPAddressManager
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP
from magma.mobilityd.ip_allocator_pool import IpAllocatorPool
from magma.mobilityd.ip_allocator_static import IPAllocatorStaticWrapper
from magma.mobilityd.ip_allocator_multi_apn import IPAllocatorMultiAPNWrapper
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.rpc_servicer import MobilityServiceRpcServicer
from magma.mobilityd.mobility_store import MobilityStore
from lte.protos.mconfig import mconfigs_pb2
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub


def _get_ipv4_allocator(store: MobilityStore, allocator_type: int,
                        static_ip_enabled: bool, multi_apn: bool,
                        dhcp_iface: str, dhcp_retry_limit: int,
                        subscriberdb_rpc_stub: SubscriberDBStub = None):
    if allocator_type == mconfigs_pb2.MobilityD.IP_POOL:
        store.dhcp_gw_info.read_default_gw()
        ip_allocator = IpAllocatorPool(store.assigned_ip_blocks,
                                       store.ip_state_map,
                                       store.sid_ips_map)
    elif allocator_type == mconfigs_pb2.MobilityD.DHCP:
        ip_allocator = IPAllocatorDHCP(
            assigned_ip_blocks=store.assigned_ip_blocks,
            ip_state_map=store.ip_state_map,
            iface=dhcp_iface,
            retry_limit=dhcp_retry_limit,
            dhcp_store=store.dhcp_store,
            gw_info=store.dhcp_gw_info)
    else:
        raise ValueError(
            "Unknown IP allocator type: %s" % allocator_type)

    if static_ip_enabled:
        ip_allocator = IPAllocatorStaticWrapper(
            subscriberdb_rpc_stub=subscriberdb_rpc_stub,
            ip_allocator=ip_allocator,
            gw_info=store.dhcp_gw_info,
            assigned_ip_blocks=store.assigned_ip_blocks,
            ip_state_map=store.ip_state_map)

    if multi_apn:
        ip_allocator = IPAllocatorMultiAPNWrapper(
            subscriberdb_rpc_stub=subscriberdb_rpc_stub,
            ip_allocator=ip_allocator)

    logging.info("Using allocator: %s static ip: %s multi_apn %s",
                 allocator_type,
                 static_ip_enabled,
                 multi_apn)
    return ip_allocator


def main():
    """ main() for MobilityD """
    service = MagmaService('mobilityd', mconfigs_pb2.MobilityD())

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
    ipv6_allocation_type = config['ipv6_ip_allocator_type']
    ipv6_allocator = IPv6AllocatorPool(
        session_prefix_alloc_mode=ipv6_allocation_type,
        sid_ips_map=store.sid_ips_map,
        ip_states_map=store.ip_state_map,
        allocated_iid=store.allocated_iid,
        sid_session_prefix_map=store.sid_session_prefix_allocated)

    # Load IPAddressManager
    ip_address_man = IPAddressManager(ipv4_allocator, ipv6_allocator, store)

    # Add all servicers to the server
    mobility_service_servicer = MobilityServiceRpcServicer(ip_address_man,
                                                           mconfig.ip_block,
                                                           config.get(
                                                               'ipv6_prefix_block'))
    mobility_service_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
