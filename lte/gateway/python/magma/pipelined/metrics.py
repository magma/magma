"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from prometheus_client import Counter, Gauge


DP_SEND_MSG_ERROR = Counter('dp_send_msg_error',
                            'Total datapath message send errors', ['cause'])
ARP_DEFAULT_GW_MAC_ERROR = Counter('arp_default_gw_mac_error',
                                   'Error with default gateway MAC resolution',
                                   [],
                                   )
OPENFLOW_ERROR_MSG = Counter(
    'openflow_error_msg',
    'Total openflow error messages received by code and type',
    ['error_type', 'error_code'])

UNKNOWN_PACKET_DIRECTION = Counter(
    'unknown_pkt_direction',
    'Counts number of times a packet is missing its flow direction',
    [],
)

NETWORK_IFACE_STATUS = Gauge(
    'network_iface_status',
    'Status of a network interface required for data pipeline',
    ['iface_name'],
)
