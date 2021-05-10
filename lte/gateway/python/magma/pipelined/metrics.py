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

ENFORCEMENT_RULE_INSTALL_FAIL = Counter(
    'enforcement_rule_install_fail',
    'Counts number of times rule install failed in enforcement app',
    ['rule_id', 'imsi'],
)

ENFORCEMENT_STATS_RULE_INSTALL_FAIL = Counter(
    'enforcement_stats_rule_install_fail',
    'Counts number of times rule install failed in enforcement stats app',
    ['rule_id', 'imsi'],
)

NETWORK_IFACE_STATUS = Gauge(
    'network_iface_status',
    'Status of a network interface required for data pipeline',
    ['iface_name'],
)

GTP_PORT_USER_PLANE_UL_BYTES = Gauge('gtp_port_user_plane_ul_bytes',
                                       'GTP port user plane uplink bytes',
                                       ['ip_addr'],
                                       )

GTP_PORT_USER_PLANE_DL_BYTES = Gauge('gtp_port_user_plane_dl_bytes',
                                       'GTP port user plane downlink bytes',
                                       ['ip_addr'],
                                       )
