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
from magma.common.redis.client import get_default_client
from magma.common.service import MagmaService
from magma.mobilityd.ip_address_man import IPAddressManager
from magma.mobilityd.rpc_servicer import MobilityServiceRpcServicer
from magma.mobilityd.mobility_store import MobilityStore
from lte.protos.mconfig import mconfigs_pb2
from magma.common.service_registry import ServiceRegistry
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub



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

    ipv6_allocator_type = config['ipv6_ip_allocator_type']

    # TODO: consider adding gateway mconfig to decide whether to
    # persist to Redis
    client = get_default_client()
    mobility_store = MobilityStore(client,
                                   config.get('persist_to_redis', False),
                                   config.get('redis_port', 6380))

    # Load IPAddressManager
    chan = ServiceRegistry.get_rpc_channel('subscriberdb', ServiceRegistry.LOCAL)
    ip_address_man = IPAddressManager(allocator_type=allocator_type,
                                      store=mobility_store,
                                      multi_apn=multi_apn,
                                      static_ip_enabled=static_ip_enabled,
                                      dhcp_iface=dhcp_iface,
                                      dhcp_retry_limit=dhcp_retry_limit,
                                      ipv6_allocation_type=ipv6_allocator_type,
                                      subscriberdb_rpc_stub=SubscriberDBStub(chan))

    # Add all servicers to the server
    mobility_service_servicer = MobilityServiceRpcServicer(ip_address_man,
                                                           mconfig.ip_block,
                                                           config.get('ipv6_prefix_block'))
    mobility_service_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
