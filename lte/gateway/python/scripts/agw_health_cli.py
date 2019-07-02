#!/usr/bin/env python3

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import ipaddress
import sys

import fire
from lte.protos.enodebd_pb2_grpc import EnodebdStub
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from lte.protos.mobilityd_pb2 import IPAddress
from orc8r.protos.common_pb2 import Void
from magma.common.service_registry import ServiceRegistry
from magma.configuration.mconfig_managers import load_service_mconfig_as_json


class AGWHealthSummary:
    def __init__(self, relay_enabled, nb_enbs_connected, allocated_ips):
        self.relay_enabled = relay_enabled
        self.nb_enbs_connected = nb_enbs_connected
        self.allocated_ips = allocated_ips

    def __str__(self):
        return """
{}
#eNBs connected: {}
#IPs allocated: {}
""".format(
            'Using Feg' if self.relay_enabled else 'Using subscriberdb',
            self.nb_enbs_connected,
            len(self.allocated_ips),
        )


def get_allocated_ips():
    chan = ServiceRegistry.get_rpc_channel('mobilityd', ServiceRegistry.LOCAL)
    client = MobilityServiceStub(chan)
    res = []

    list_blocks_resp = client.ListAddedIPv4Blocks(Void())
    for block_msg in list_blocks_resp.ip_block_list:

        list_ips_resp = client.ListAllocatedIPs(block_msg)
        for ip_msg in list_ips_resp.ip_list:
            if ip_msg.version == IPAddress.IPV4:
                ip = ipaddress.IPv4Address(ip_msg.address)
            elif ip_msg.address == IPAddress.IPV6:
                ip = ipaddress.IPv6Address(ip_msg.address)
            else:
                continue
            res.append(ip)
    return res


def gateway_health_status():
    config = load_service_mconfig_as_json('mme')

    ''' eNB status for #eNBs connected '''
    chan = ServiceRegistry.get_rpc_channel('enodebd', ServiceRegistry.LOCAL)
    client = EnodebdStub(chan)
    status = client.GetAllEnodebStatus(Void())
    status_list = status.enb_status_list

    ''' List IPs allocated '''
    allocated_ips = get_allocated_ips()

    health_summary = AGWHealthSummary(
        relay_enabled=config['relayEnabled'],
        nb_enbs_connected={enb_status.ip_address: enb_status.connected
                           for enb_status in status_list},
        allocated_ips=allocated_ips,
    )
    return str(health_summary)


if __name__ == '__main__':
    print('Access Gateway health summary')
    if len(sys.argv) == 1:
        fire.Fire(gateway_health_status)
    else:
        fire.Fire({
            'status': gateway_health_status,
        })
