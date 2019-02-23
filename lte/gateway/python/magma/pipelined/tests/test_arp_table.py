"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import time
import warnings

from magma.pipelined.openflow.registers import DIRECTION_REG
from magma.pipelined.tests.app.flow_query import RyuRestFlowQuery
from magma.pipelined.tests.app.table_isolation import RyuRestTableIsolator,\
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder,\
    ARPPacketBuilder


def _pkt_total(stats):
    return sum(n.packets for n in stats)


@unittest.skip("Rest tests currently disabled and are left as an api example")
class ARPTableTest(unittest.TestCase):
    TID = 2
    IFACE = "gtp_br0"
    MAC_DEST = "0e:9f:0f:0d:98:4e"
    IP_DEST = "192.168.128.0"

    def setUp(self):
        warnings.simplefilter("ignore")

    def test_rest_arp_flow(self):
        """
        Sends an arp request to the ARP table

        Assert:
            The arp rule is matched 2 times for each arp packet
            No other rule is matched
        """
        isolator = RyuRestTableIsolator(
            RyuForwardFlowArgsBuilder(self.TID).set_reg_value(DIRECTION_REG,
                                                              0x10)
                                               .build_requests()
        )
        flow_query = RyuRestFlowQuery(
            self.TID,
            match={
                'eth_type': 2054,
                DIRECTION_REG: 16,
                'arp_tpa': self.IP_DEST + '/255.255.255.0'
            }
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packets = ARPPacketBuilder().set_arp_layer(self.IP_DEST + "/28")\
                                    .build()

        # 16 as the bitmask was /28
        num_pkts = 16
        arp_start = flow_query.lookup()[0].packets
        total_start = _pkt_total(RyuRestFlowQuery(self.TID).lookup())

        with isolator:
            pkt_sender.get_response(packets)
            time.sleep(2.5)

        arp_final = flow_query.lookup()[0].packets
        total_final = _pkt_total(RyuRestFlowQuery(self.TID).lookup())

        self.assertEqual(arp_final - arp_start, num_pkts * 2)
        self.assertEqual(total_final - total_start, num_pkts * 2)

    def test_rest_ip_flow(self):
        """
        Sends an ip packet

        Assert:
            The correct ip rule is matched
            No other rule is matched
        """
        isolator = RyuRestTableIsolator(
            RyuForwardFlowArgsBuilder(self.TID).set_reg_value(DIRECTION_REG,
                                                              0x1)
                                               .build_requests()
        )
        flow_query = RyuRestFlowQuery(
            self.TID, match={
                'eth_type': 2048,
                DIRECTION_REG: 1
            }
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder()\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:04")\
            .set_ip_layer(self.IP_DEST, "10.0.0.0")\
            .build()

        num_pkts = 42
        ip_start = flow_query.lookup()[0].packets
        total_start = _pkt_total(RyuRestFlowQuery(self.TID).lookup())

        with isolator:
            pkt_sender.send(packet, num_pkts)
            time.sleep(2.5)

        total_final = _pkt_total(RyuRestFlowQuery(self.TID).lookup())
        ip_final = flow_query.lookup()[0].packets

        self.assertEqual(ip_final - ip_start, num_pkts)
        self.assertEqual(total_final - total_start, num_pkts)


if __name__ == "__main__":
    unittest.main()
