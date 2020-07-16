import redis
import logging
from collections import defaultdict
from typing import MutableMapping

from magma.mobilityd.ip_descriptor import IPDesc
from lte.protos.mconfig.mconfigs_pb2 import MobilityD

from magma.mobilityd import mobility_store as store

from .ip_allocator_dhcp import IPAllocatorDHCP
from .ip_allocator_base import IPAllocator
from .ip_allocator_static import IpAllocatorStatic
from .ip_descriptor_map import IpDescriptorMap
from .uplink_gw import UplinkGatewayInfo


class IPAllocatorFactory:
    def __init__(self, ip_allocator_type, mobilityd_conf, redis_config):
        """ Creates IP allocator, states and sid maps in redis/memory flavors
        """
        persist_to_redis = mobilityd_conf.get('persist_to_redis', False)
        logging.debug('Persist to Redis: %s', persist_to_redis)

        if not persist_to_redis:
            # Do not use redis backend
            self._assigned_ip_blocks = set()  # {ip_block}
            self._sid_ips_map = defaultdict(IPDesc)  # {SID=>IPDesc}
            self.ip_states = defaultdict(dict)  # {state=>{ip=>ip_desc}}
            self.dhcp_store = {}
            self._dhcp_gw_info = UplinkGatewayInfo(defaultdict(str))
        else:
            # Enable redis backend
            redis_port = redis_config.get('port', 6380)
            redis_host = redis_config.get('bind', 'localhost')
            client = redis.Redis(host=redis_host,
                                 port=redis_port)
            self._assigned_ip_blocks = store.AssignedIpBlocksSet(client)
            self._sid_ips_map = store.IPDescDict(client)
            self.ip_states = store.defaultdict_key(
                lambda key: store.ip_states(client, key))
            self.dhcp_store = store.MacToIP()  # mac => DHCP_State
            self._dhcp_gw_info = UplinkGatewayInfo(store.GatewayInfoMap())

        self._ip_state_map = IpDescriptorMap(self.ip_states)
        self.allocator_type = ip_allocator_type

        logging.info("Using allocator: %s", self.allocator_type)
        if self.allocator_type == MobilityD.IP_POOL:
            self._dhcp_gw_info.read_default_gw()
            self._ip_allocator = IpAllocatorStatic(self._assigned_ip_blocks,
                                                   self._ip_state_map,
                                                   self._sid_ips_map)
        elif self.allocator_type == MobilityD.DHCP:
            iface = mobilityd_conf.get('dhcp_iface', 'dhcp0')
            retry_limit = mobilityd_conf.get('retry_limit', 300)
            self._ip_allocator = IPAllocatorDHCP(self._assigned_ip_blocks,
                                                 self._ip_state_map,
                                                 iface=iface,
                                                 retry_limit=retry_limit,
                                                 dhcp_store=self.dhcp_store,
                                                 gw_info=self._dhcp_gw_info)

    def get_allocator(self) -> IPAllocator:
        return self._ip_allocator

    def get_sid_ips_map(self) -> MutableMapping[str, IPDesc]:
        return self._sid_ips_map

    def get_ip_states_map(self) -> IpDescriptorMap:
        return self._ip_state_map

    def get_dhcp_gw_info(self) -> UplinkGatewayInfo:
        return self._dhcp_gw_info

