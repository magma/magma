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
