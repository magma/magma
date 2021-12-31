"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import ctypes
import logging
import socket
import struct
import subprocess
from builtins import input
from socket import AF_INET, htons
from subprocess import call
from sys import argv
from threading import Thread

import netifaces
from bcc import BPF
from magma.pipelined.mobilityd_client import get_mobilityd_gw_info
from pyroute2 import IPDB, IPRoute, NetlinkError, NetNS, NSPopen, protocols
from scapy.layers.l2 import getmacbyip

LOG = logging.getLogger("pipelined.ebpf")
# LOG.setLevel(logging.DEBUG)

BASE_MAP_FS = "/sys/fs/bpf/"
BPF_UL_FILE = "/var/opt/magma/ebpf/ebpf_ul_handler.c"
BPF_DL_FILE = "/var/opt/magma/ebpf/ebpf_dl_handler.c"
UL_MAP_NAME = "ul_map"
DL_MAP_NAME = "dl_map"
DL_CFG_ARRAY_NAME = "cfg_array"

"""
    Use pipelineD configuration to initialize eBPF manager for AGW.
"""


def get_ebpf_manager(config):

    if 'ebpf' in config:
        enabled = config['ebpf']['enabled']
    else:
        enabled = False
    gw_info = get_mobilityd_gw_info()
    for gw in gw_info:
        if gw.vlan == "":
            bpf_man = ebpf_manager(config['nat_iface'], config['enodeb_iface'], gw.ip, enabled)
            if enabled:
                # TODO: For Development purpose dettch and attach latest eBPF code.
                # Remove this for production deployment
                bpf_man.detach_ul_ebpf()
                bpf_man.attach_dl_ebpf()
                LOG.info("eBPF manager: initilized: enabled: %s", enabled)
            return bpf_man

    LOG.info("eBPF manager: Not initilized")
    return None


"""
    eBPF manager for AGW.
    Initialize eBPF based datapath for AGW as per the pipelineD config.

    Returns:
        eBPF manager object.
"""


class ebpf_manager:
    def __init__(self, sgi_if_name: str, s1_if_name: str, gw_ip: str, enabled=True, bpf_ul_file: str = BPF_UL_FILE, bpf_dl_file: str = BPF_DL_FILE):

        self.b_ul = BPF(src_file=bpf_ul_file, cflags=[''])
        self.b_dl = BPF(src_file=bpf_dl_file, cflags=[''])
        self.s1_fn = self.b_ul.load_func("gtpu_ingress_handler", BPF.SCHED_CLS)
        self.sgi_fn = self.b_dl.load_func("gtpu_egress_handler", BPF.SCHED_ACT)
        self.ul_map = self.b_ul.get_table(UL_MAP_NAME)
        self.dl_map = self.b_dl.get_table(DL_MAP_NAME)
        self.cfg_array = self.b_dl.get_table(DL_CFG_ARRAY_NAME)
        self.sgi_if_name = sgi_if_name
        self.s1_if_name = s1_if_name
        self.ul_src_mac = self._get_mac_address(sgi_if_name)
        self.ul_gw_mac = self._get_mac_address_of_ip(gw_ip)
        self.sgi_if_index = self._get_ifindex(self.sgi_if_name)
        self.enabled = enabled

    """Attach eBPF Uplink traffic handler
    """

    def attach_ul_ebpf(self):
        s1_if_index = self._get_ifindex(self.s1_if_name)

        ipr = IPRoute()
        try:
            ipr.tc("add", "clsact", s1_if_index)
        except NetlinkError as ex:
            LOG.error("error adding ingress ")

        try:
            ipr.tc(
                "add-filter", "bpf", s1_if_index, ":1", fd=self.s1_fn.fd, name=self.s1_fn.name,
                parent="ffff:fff2", classid=1, direct_action=True,
            )
        except NetlinkError as ex:
            LOG.error("error adding ingress ")

        LOG.debug("Attach done")


    def attach_dl_ebpf(self):
        """
        Attach eBPF downlink traffic handler
        """

        ipr = IPRoute()
        try:
            ipr.tc("add", "clsact", self.sgi_if_index)
        except NetlinkError as ex:
            LOG.error("error adding ingress clasct for dl %s", ex)

        try:

            action = {"kind": "bpf", "fd": self.sgi_fn.fd, "name": self.sgi_fn.name, "action": "ok"}
            ipr.tc(
                "add-filter", "u32", self.sgi_if_index, ":1", parent="ffff:fff2", action=[action],
                protocol=protocols.ETH_P_ALL, classid=1, target=0x10002, keys=['0x0/0x0+0'],
            )
            '''
            This failes for some reason, using above as a workaround
            TODO figure out why
            ipr.tc(
                "add-filter", "bpf", self.sgi_if_index, ":1", fd=self.sgi_fn.fd, name=self.sgi_fn.name,
                parent="ffff:fff2", classid=1, direct_action=True,
            )
            '''
        except NetlinkError as ex:
            LOG.error("error adding ingress filter for dl %s", ex)

        key = self.cfg_array.Key(0)
        # TODO add as pipelined.yml
        ifindex = self._get_ifindex('gtpu_sys_2152')
        val = self.cfg_array.Leaf(ifindex)
        self.cfg_array[key] = val

        LOG.debug("Attach done")

    """Remove the Uplink eBPF handler and associated maps.
    """

    def detach_ul_ebpf(self):
        s1_if_index = self._get_ifindex(self.s1_if_name)

        ipr = IPRoute()
        try:
            ipr.tc("del", "ingress", s1_if_index, "ffff:")
        except NetlinkError as ex:
            pass
        sys_file = BASE_MAP_FS + UL_MAP_NAME
        out1 = subprocess.run(["unlink", sys_file], capture_output=True)
        LOG.debug(out1)



    def detach_dl_ebpf(self):
        """
        Remove the Downlink eBPF handler and associated maps.
        """

        ipr = IPRoute()
        try:
            ipr.tc("del", "clsact", self.sgi_if_index)
        except NetlinkError as ex:
            pass
        sys_file = BASE_MAP_FS + DL_MAP_NAME
        out1 = subprocess.run(["unlink", sys_file], capture_output=True)
        LOG.debug(out1)
        sys_file = BASE_MAP_FS + DL_CFG_ARRAY_NAME
        out1 = subprocess.run(["unlink", sys_file], capture_output=True)
        LOG.debug(out1)

    """Add uplink session entry
    """

    def add_ul_entry(self, mark: int, ue_ip: str):
        if not self.enabled:
            return
        sz = len(self.ul_map)
        ip_addr = self._pack_ip(ue_ip)
        LOG.debug(
            "Add entry: ip: %x mac src %s mac dst: %s" %
            (ip_addr, self._unpack_mac_addr(self.ul_src_mac), self._unpack_mac_addr(self.ul_gw_mac)),
        )

        key = self.ul_map.Key(ip_addr)
        val = self.ul_map.Leaf(mark, self.sgi_if_index, self.ul_src_mac, self.ul_gw_mac)
        self.ul_map[key] = val

    def add_dl_entry(self, ue_ip: str, remote_ipv4: str, tunnel_id: int):
        """
        Add downlink session entry
        """
        if not self.enabled:
            return
        ip_addr = self._pack_ip(ue_ip)
        LOG.debug(
            "Add entry: ip: %x remote ipv4 %s tunnel id: %d" %
            (ip_addr, remote_ipv4, tunnel_id),
        )

        key = self.dl_map.Key(ip_addr)
        val = self.dl_map.Leaf(self._pack_ip(remote_ipv4),
                               socket.htonl(tunnel_id))
        self.dl_map[key] = val

    """Delete uplink session entry
    """

    def del_ul_entry(self, ue_ip: str):
        ip_addr = self._pack_ip(ue_ip)
        key = self.ul_map.Key(ip_addr)

        self.ul_map.pop(key, None)


    def del_dl_entry(self, ue_ip: str):
        """
        Delete downlink session entry
        """
        ip_addr = self._pack_ip(ue_ip)
        key = self.dl_map.Key(ip_addr)

        self.dl_map.pop(key, None)

    """Dump entire ulink session eBPF map
    """

    def print_ul_map(self):

        for k, v in self.ul_map.items():
            ue_ip = self._unpack_ip(k.ue_ip)
            mark = v.mark
            egress_dev_index = v.e_if_index
            egress_dev_name = self._get_if_name(egress_dev_index)
            dst_mac = self._unpack_mac_addr(v.mac_dst)
            src_mac = self._unpack_mac_addr(v.mac_src)

            print(
                "UE: %s -> {mark: %d, dev: %s (%d), src_mac %s dst_mac %s" %
                (ue_ip, mark, egress_dev_name, egress_dev_index, src_mac, dst_mac),
            )

    def print_dl_map(self):
        """
        Dump entire downlink session eBPF map
        """
        for k, v in self.dl_map.items():
            ue_ip = self._unpack_ip(k.ue_ip)
            remote_ipv4 = self._unpack_ip(v.remote_ipv4)
            tunnel_id = socket.ntohl(v.tunnel_id)

            print(
                "UE: %s -> {remote_ipv4: %s, tunnel_id: %d" %
                (ue_ip, remote_ipv4, tunnel_id),
            )

    def print_dl_cfg(self):
        """
        Dump entire cfg array session eBPF map
        """
        for k, v in self.cfg_array.items():
            ifindex = v.if_idx

            print(
                "0: %d " %
                (ifindex),
            )

    def _get_ifindex(self, if_name: str):
        sys_file = "/sys/class/net/" + if_name + "/ifindex"
        ifindex = subprocess.run(["cat", sys_file], capture_output=True)
        return int(ifindex.stdout.decode('utf-8'))

    def _get_if_name(self, if_index: int):
        for if_name in netifaces.interfaces():
            idx = self._get_ifindex(if_name)
            if idx == if_index:
                return if_name
        return None

    def _get_mac_address(self, if_name: str):
        addr_str = netifaces.ifaddresses(self.sgi_if_name)[netifaces.AF_LINK][0]['addr']
        LOG.debug("if-name: %s, mac: %s" % (if_name, addr_str))
        return self._pack_mac_addr(addr_str)

    def _get_mac_address_of_ip(self, ip_addr: str):
        addr_str = getmacbyip(ip_addr)
        LOG.debug("IP: %s, mac: %s" % (ip_addr, addr_str))
        return self._pack_mac_addr(addr_str)

    def _pack_ip(self, ip_str: str):
        packedIP = socket.inet_aton(ip_str)
        return socket.htonl(struct.unpack("!L", packedIP)[0])

    def _unpack_ip(self, ip: int):
        ip_ = socket.ntohl(ip).to_bytes(4, 'big')
        return socket.inet_ntoa(ip_)

    def _pack_mac_addr(self, mac_addr: str):
        mac_bytes = bytes.fromhex(mac_addr.replace(':', ''))
        return (ctypes.c_ubyte * 6).from_buffer(bytearray(mac_bytes))

    def _unpack_mac_addr(self, mac_addr: ctypes.c_ubyte):
        mac_bytes = bytearray(mac_addr)
        return mac_bytes.hex(":")

import time
# for debugging
if __name__ == "__main__":
    bm = ebpf_manager("sgi0", "enb0", "10.0.2.2", bpf_ul_file=BPF_UL_FILE, bpf_dl_file=BPF_DL_FILE, enabled=True)

    bm.detach_ul_ebpf()
    bm.attach_ul_ebpf()
    bm.add_ul_entry(204, '192.168.128.11')
    bm.add_ul_entry(204, '192.168.128.12')
    bm.print_ul_map()
    bm.del_ul_entry('192.168.128.12')
    bm.print_ul_map()
    bm.del_ul_entry('192.168.128.12')
    bm.print_ul_map()

    bm.detach_dl_ebpf()
    bm.attach_dl_ebpf()
    bm.add_dl_entry('192.168.128.11', '10.1.1.1', 123)
    bm.add_dl_entry('192.168.128.12', '10.2.2.2', 555)
    bm.print_dl_map()
    bm.del_dl_entry('192.168.128.12')
    bm.print_dl_map()
    bm.del_dl_entry('192.168.128.12')
    bm.print_dl_map()
    bm.print_dl_cfg()
