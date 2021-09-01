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
from abc import ABC

from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import IPFlowDL, UESessionSet, UESessionState
from magma.subscriberdb.sid import SIDUtils


class CreateMMESessionUtils(ABC):
    """Class for creating the MME Session"""

    def __init__(
        self,
        imsi: str,
        priority: int,
        ue_v4_addr: str,
        ue_v6_addr: str,
        enb_ipv4_addr: str,
        apn: str,
        vlan: int,
        in_teid: int,
        out_teid: int,
        ue_state: int,
        flow_dl: str,
    ):
        """Do create the MME Sessions

        Args:
            imsi: Subscriber for the UE
            priority: Priority for the rule
            ue_v4_addr: UE IPv4 address
            ue_v6_addr : UE IPv6 address
            enb_ipv4_addr: GNB Ipv4 address
            apn: APN Name
            vlan: Vlan Info
            in_teid: Incoming teid
            out_teid: Outgoing teid
            ue_state: UE State
            flow_dl: Flow DL to be included
        """
        ue_ipv4_addr = None
        if (ue_v4_addr):
            ue_ipv4_addr = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, ue_v4_addr),
            )
            # address=args.ue_ipv4_addr.encode('utf-8'))

        ue_ipv6_addr = None
        if (ue_v6_addr):
            ue_ipv6_addr = IPAddress(
                version=IPAddress.IPV6,
                address=socket.inet_pton(socket.AF_INET6, ue_v6_addr),
            )
            # address=args.ue_ipv6_addr.encode('utf-8'))

        enb_ip_addr = None
        if (enb_ipv4_addr):
            enb_ip_addr = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, enb_ipv4_addr),
            )
            # address=args.enb_ip_addr.encode('utf-8'))

        config_ue_state = {
            'ADD':
            UESessionState(ue_config_state=UESessionState.ACTIVE),
            'DEL':
            UESessionState(ue_config_state=UESessionState.UNREGISTERED),
            'ADD_IDLE':
            UESessionState(ue_config_state=UESessionState.INSTALL_IDLE),
            'DEL_IDLE':
            UESessionState(ue_config_state=UESessionState.UNINSTALL_IDLE),
            'SUSPENDED':
            UESessionState(ue_config_state=UESessionState.SUSPENDED_DATA),
            'RESUME':
            UESessionState(ue_config_state=UESessionState.RESUME_DATA),
        }

        ip_flow_dl = None
        if (flow_dl == "ENABLE"):
            dest_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.128.12"),
            )
            src_ip = IPAddress(
                version=IPAddress.IPV4,
                address=socket.inet_pton(socket.AF_INET, "192.168.129.64"),
            )

            params = 71
            tcp_sport = 5002
            proto_type = 6
            ip_flow_dl = IPFlowDL(
                set_params=params, tcp_src_port=tcp_sport,
                ip_proto=proto_type,
                dest_ip=dest_ip, src_ip=src_ip,
            )

        self._set_pg_session = UESessionSet(
            subscriber_id=SIDUtils.to_pb(imsi),
            precedence=priority,
            ue_ipv4_address=ue_ipv4_addr,
            ue_ipv6_address=ue_ipv6_addr,
            enb_ip_address=enb_ip_addr,
            apn=apn,
            vlan=vlan,
            in_teid=in_teid,
            out_teid=out_teid,
            ue_session_state=config_ue_state.get(ue_state),
            ip_flow_dl=ip_flow_dl,
        )
