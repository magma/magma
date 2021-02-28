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
import socket
import ipaddress
from lte.protos.pipelined_pb2 import (
    UESessionSet,
    UESessionState,
    IPFlowDL
)
from lte.protos.mobilityd_pb2 import IPAddress
from magma.subscriberdb.sid import SIDUtils

class CreatePGSessionUtils:

    def __init__(self, imsi:str, priority:int, ue_v4_addr:str, ue_v6_addr:str, enb_ipv4_addr:str,
                 apn:str, vlan:int, in_teid:int, out_teid:int, ue_state:str, flow_dl:str):
        """
        Create session place holder with subs_id, f_teid and session_version.
        This will anchor point for create sessions with PDR, FAR & QER.
        """
        ue_ipv4_addr=None
        if (len(ue_v4_addr)):
            ue_ipv4_addr = IPAddress(version=IPAddress.IPV4,
                                     address=socket.inet_pton(socket.AF_INET, ue_v4_addr))
                                 #address=args.ue_ipv4_addr.encode('utf-8'))

        ue_ipv6_addr=None
        if (len(ue_v6_addr)):
            ue_ipv6_addr = IPAddress(version=IPAddress.IPV6,
                                     address=socket.inet_pton(socket.AF_INET6, ue_v6_addr))
                                 #address=args.ue_ipv6_addr.encode('utf-8'))

        enb_ip_addr=None
        if (len(enb_ipv4_addr)):
            enb_ip_addr = IPAddress(version=IPAddress.IPV4,
                                    address=socket.inet_pton(socket.AF_INET, enb_ipv4_addr))
                                    #address=args.enb_ip_addr.encode('utf-8'))

        if (ue_state == 'ADD'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.ACTIVE)
        elif (ue_state == 'DEL'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.UNREGISTERED)
        elif (ue_state == 'ADD_IDLE'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.INSTALL_IDLE)
        elif (ue_state == 'DEL_IDLE'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.UNINSTALL_IDLE)
        elif (ue_state == 'SUSPENDED'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.SUSPENDED_DATA)
        elif (ue_state == 'RESUME'):
            config_ue_state=UESessionState(ue_config_state=UESessionState.RESUME_DATA)

        ip_flow_dl = None
        if (flow_dl == "ENABLE"):
            dest_ip=IPAddress(version=IPAddress.IPV4,
                             address=socket.inet_pton(socket.AF_INET, "192.168.128.12"))
            src_ip=IPAddress(version=IPAddress.IPV4,
                             address=socket.inet_pton(socket.AF_INET, "192.168.129.64"))

            ip_flow_dl = IPFlowDL(set_params=71, tcp_dst_port=0,
                                  tcp_src_port=5002, udp_dst_port=0, udp_src_port=0, ip_proto=6,
                                  dest_ip=dest_ip,src_ip=src_ip)

        self._set_pg_session = UESessionSet(subscriber_id=SIDUtils.to_pb(imsi),
                                            precedence=priority,
                                            ue_ipv4_address=ue_ipv4_addr,
                                            ue_ipv6_address=ue_ipv6_addr,
                                            enb_ip_address=enb_ip_addr,
                                            apn=apn,
                                            vlan=vlan,
                                            in_teid=in_teid,
                                            out_teid=out_teid,
                                            ue_session_state=config_ue_state,
                                            ip_flow_dl=ip_flow_dl)
