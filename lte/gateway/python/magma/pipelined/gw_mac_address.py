"""
Copyright 2022 The Magma Authors.

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

from lte.protos.mobilityd_pb2 import IPAddress
from scapy.arch import get_if_addr, get_if_hwaddr
from scapy.data import ETH_P_ALL, ETHER_BROADCAST
from scapy.error import Scapy_Exception
from scapy.layers.inet6 import getmacbyip6
from scapy.layers.l2 import ARP, Dot1Q, Ether
from scapy.sendrecv import srp1


def get_gw_mac_address(ip: IPAddress, vlan: str, non_nat_arp_egress_port: str) -> str:
    if ip.version == IPAddress.IPV4:
        return _get_gw_mac_address_v4(ip, vlan, non_nat_arp_egress_port)
    if ip.version == IPAddress.IPV6:
        if vlan == "NO_VLAN":
            return _get_gw_mac_address_v6(ip)
        else:
            gw_ip = ipaddress.ip_address(ip.address)
            logging.error("Not supported: GW IPv6: %s over vlan %d", str(gw_ip), vlan)
            return ""
    return ""


def _get_gw_mac_address_v4(ip: IPAddress, vlan: str, non_nat_arp_egress_port: str) -> str:
    try:
        gw_ip = ipaddress.ip_address(ip.address)
        logging.debug(
            "sending arp via egress: %s",
            non_nat_arp_egress_port,
        )
        eth_mac_src = get_if_hwaddr(non_nat_arp_egress_port)
        psrc = "0.0.0.0"
        egress_port_ip = get_if_addr(non_nat_arp_egress_port)
        if egress_port_ip:
            psrc = egress_port_ip

        pkt = Ether(dst=ETHER_BROADCAST, src=eth_mac_src)
        if vlan.isdigit():
            pkt /= Dot1Q(vlan=int(vlan))
        pkt /= ARP(op="who-has", pdst=gw_ip, hwsrc=eth_mac_src, psrc=psrc)
        logging.debug("ARP Req pkt %s", pkt.show(dump=True))

        res = srp1(
            pkt,
            type=ETH_P_ALL,
            iface=non_nat_arp_egress_port,
            timeout=1,
            verbose=0,
            nofilter=1,
            promisc=0,
        )

        if res is None:
            logging.debug("Got Null response")
            return ""

        logging.debug("ARP Res pkt %s", res.show(dump=True))
        if str(res[ARP].psrc) != str(gw_ip):
            logging.warning(
                "Unexpected IP in ARP response. expected: %s pkt: %s",
                str(gw_ip),
                res.show(dump=True),
            )
            return ""
        if vlan.isdigit():
            if Dot1Q in res and str(res[Dot1Q].vlan) == vlan:
                mac = res[ARP].hwsrc
            else:
                logging.warning(
                    "Unexpected vlan in ARP response. expected: %s pkt: %s",
                    vlan,
                    res.show(dump=True),
                )
                return ""
        else:
            mac = res[ARP].hwsrc
        return mac

    except Scapy_Exception as ex:
        logging.warning("Error in probing Mac address: err %s", ex)
        return ""
    except ValueError:
        logging.warning(
            "Invalid GW Ip address: [%s] or vlan %s",
            str(ip), vlan,
        )
        return ""


def _get_gw_mac_address_v6(ip: IPAddress) -> str:
    try:
        gw_ip = ipaddress.ip_address(ip.address)
        mac = getmacbyip6(str(gw_ip))
        logging.debug("Got mac %s for IP: %s", mac, gw_ip)
        return mac

    except Scapy_Exception as ex:
        logging.warning("Error in probing Mac address: err %s", ex)
        return ""
    except ValueError:
        logging.warning(
            "Invalid GW Ip address: [%s]",
            str(ip),
        )
        return ""
