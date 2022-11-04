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

import getmac
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


def get_mac_address(network_request=True, **kwargs) -> str:
    if len(kwargs) != 1:
        raise ValueError('Must specify exactly one keyword argument for get_mac_address')
    mac_address = getmac.get_mac_address(
        interface=kwargs.get('interface'), ip=kwargs.get('ip'), ip6=kwargs.get('ip6'),
        hostname=kwargs.get('hostname'), network_request=network_request,
    )
    if mac_address is None:
        raise ValueError(f"No mac address found for {list(kwargs.keys())[0]} {list(kwargs.values())[0]}")
    return mac_address
