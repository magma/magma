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
import ipaddress
import threading


class InternalIPAllocator:
    """
    Internal IP Allocator

    This class assigns internal IP from the same subnet as the bridge to the UE,
    this is used primarily for redirection to catch ARPs as they'll be directed
    to the bridge(where we auto respond to them)
    """

    def __init__(self, config):
        # Internal IP allocation is only required in CWF networks
        if 'setup_type' not in config or config['setup_type'] != 'CWF':
            return

        self.internal_ip_network = \
            ipaddress.ip_network(config['internal_ip_subnet'])
        self.internal_ip_iterator = self.internal_ip_network.hosts()
        self._invalid_ips = []
        self._invalid_ips.append(config['bridge_ip_address'])
        if 'quota_check_ip' in config:
            self._invalid_ips.append(config['quota_check_ip'])
        self._lock = threading.Lock()  # write lock

    def next_ip(self):
        with self._lock:
            ip = next(self.internal_ip_iterator, None)
            while ip is None or str(ip) in self._invalid_ips:
                ip = next(self.internal_ip_iterator, None)
                if ip is None:
                    self.internal_ip_iterator = self.internal_ip_network.hosts()
        return ip
