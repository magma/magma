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

import abc
import logging

logging.getLogger("scapy.runtime").setLevel(logging.ERROR)
from scapy.all import rdpcap, sendp, srp, wrpcap


class PacketInjector(metaclass=abc.ABCMeta):
    """Packet injection interface"""
    @abc.abstractmethod
    def send(self, pkt, count):
        """
        Send packet
        Args:
            pkt (bytes): packet or array of packets to send
            count (int): number of packets to send
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def get_response(self, pkt, timeout):
        """
        Send pkt and return response
        Args:
            pkt (bytes): packet or array of packets to send
            timeout (int): response waiting time, default 1.5s
        Return:
            A list of packets received
        """
        raise NotImplementedError()


class ScapyPacketInjector(PacketInjector):
    """
    Scapy Packet Injector, the pkt arg for send, get_response
    can be either bytes or Scapy pkts
    """
    def __init__(self, iface, pcap_filename=None):
        self._iface = iface
        self._pcap_filename = pcap_filename

    def send(self, pkt, count=1):
        if self._pcap_filename:
            wrpcap(self._pcap_filename, pkt, append=True)
        sendp(pkt, iface=self._iface, count=count, verbose=False)

    def send_pcap(self, name):
        pkts = rdpcap(name)
        if self._pcap_filename:
            wrpcap(self._pcap_filename, pkts, append=True)
        sendp(pkts, iface=self._iface, verbose=False)

    def get_response(self, pkt, timeout=1.5):
        packets = srp(pkt, iface=self._iface, timeout=timeout, verbose=False)
        if self._pcap_filename:
            wrpcap(self._pcap_filename, pkt, append=True)
            for pkt in packets[0]:
                wrpcap(self._pcap_filename, pkt, append=True)
        return packets[0]
