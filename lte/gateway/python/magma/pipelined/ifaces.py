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
import asyncio
from typing import Optional

import netifaces
from magma.pipelined.metrics import NETWORK_IFACE_STATUS

POLL_INTERVAL_SECONDS = 3


@asyncio.coroutine
def monitor_ifaces(iface_names):
    """
    Call to poll the network interfaces and set the corresponding metric
    """
    while True:
        active = set(netifaces.interfaces())
        for iface in iface_names:
            status = 1 if iface in active else 0
            NETWORK_IFACE_STATUS.labels(iface_name=iface).set(status)
        yield from asyncio.sleep(POLL_INTERVAL_SECONDS)


def get_mac_address_from_iface(iface: str) -> str:
    ifaddress = netifaces.ifaddresses(iface)
    if netifaces.AF_LINK in ifaddress:
        return ifaddress[netifaces.AF_LINK][0]['addr']
    raise ValueError(f"No mac address found for interface {iface}")


def get_mac_address_from_ip4(ip4: str) -> str:
    return _get_mac_from_ip(ip4, "ip4")


def get_mac_address_from_ip6(ip6: str) -> str:
    return _get_mac_from_ip(ip6, "ip6")


def _get_mac_from_ip(ip: str, ipv: str) -> str:
    af_type = {
        "ip4": netifaces.AF_INET,
        "ip6": netifaces.AF_INET6,
    }[ipv]
    for iface in netifaces.interfaces():
        try:
            ifaddress = netifaces.ifaddresses(iface)
            if ifaddress[af_type][0]['addr'].split("%")[0] == ip:
                return ifaddress[netifaces.AF_LINK][0]['addr']
        except (KeyError, IndexError, ValueError):
            continue
    raise ValueError(f"No mac address found for {ipv} {ip}")
