"""
Test Configuration Module

Centralized configuration for eBPF GTP tests.
Matches production struct layouts from ebpf_gtp_decap.c / ebpf_gtp_encap.c.
"""

import os
import socket
import struct
import subprocess
from dataclasses import dataclass


# Path to production eBPF source files
# __file__ = tests/lib/config.py -> dirname x3 = ebpf/
EBPF_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
DECAP_SOURCE = os.path.join(EBPF_DIR, 'ebpf_gtp_decap.c')
ENCAP_SOURCE = os.path.join(EBPF_DIR, 'ebpf_gtp_encap.c')
MARK_SOURCE = os.path.join(EBPF_DIR, 'ebpf_gtp_veth0_mark.c')

# Compilation flags matching production (ebpf_gtp_manager.py line 1239)
PRODUCTION_CFLAGS = ['-I', EBPF_DIR, '-DDISABLE_DEBUG', '-O2']

# Stats counter indices (from ebpf_gtp_decap.c / ebpf_gtp_encap.c)
STATS_UL_PACKETS = 0
STATS_UL_BYTES = 1
STATS_DL_PACKETS = 2
STATS_DL_BYTES = 3
STATS_UL_ERRORS = 4
STATS_DL_ERRORS = 5
STATS_SESSION_MISS = 6
STATS_TEID_MISMATCH = 7
STATS_GTP_DECAP_SUCCESS = 8
STATS_GTP_ENCAP_SUCCESS = 9
STATS_PKT_TOO_SHORT = 10
STATS_INVALID_GTP = 11
STATS_ADJUST_HEAD_FAIL = 12
STATS_TOTAL_PROCESSED = 13
STATS_UE_ATTACH = 14
STATS_UE_DETACH = 15
STATS_PKT_FORWARDED = 16
STATS_PKT_DROPPED = 17
STATS_SESSION_ACTIVE = 18
STATS_QOS_APPLIED = 19
STATS_INACTIVE_SESSION = 20
STATS_DOUBLE_ENCAP_AVOIDED = 21

# Mark handler stats (from ebpf_gtp_veth0_mark.c)
STATS_VETH0_PACKETS_PROCESSED = 50
STATS_VETH0_MARK_RESTORED = 51
STATS_VETH0_MARK_FALLBACK = 52
STATS_VETH0_SESSION_MISS = 53

# Config map keys (from production code)
CONFIG_S1U_IFINDEX = 0
CONFIG_SGI_IFINDEX = 1
CONFIG_OVS_IFINDEX = 2
CONFIG_DEBUG_LEVEL = 3
CONFIG_SGI_IP = 4
CONFIG_EBPF_VETH_IFINDEX = 5


def ip_to_host_order(ip_str: str) -> int:
    """Convert IP string to host-order integer (as used in production BPF maps)"""
    return struct.unpack("!I", socket.inet_aton(ip_str))[0]


def mac_str_to_bytes(mac_str: str) -> bytes:
    """Convert MAC string (aa:bb:cc:dd:ee:ff) to bytes"""
    return bytes(int(b, 16) for b in mac_str.split(':'))


@dataclass
class GTPUTestConfig:
    """Configuration for GTP-U eBPF tests"""

    # Network interfaces (veth pair)
    tx_iface: str = "ebpf_test_tx"   # Simulates gtp_veth1 (S1-U side)
    rx_iface: str = "ebpf_test_rx"   # Simulates gtp_veth0 (OVS side)

    # IP addresses
    enb_ip: str = "10.0.2.1"         # eNodeB IP (GTP tunnel endpoint)
    sgw_ip: str = "10.0.2.2"         # SGW IP (our side, GTP source for encap)
    ue_ip: str = "192.168.128.100"   # UE IP address
    external_ip: str = "8.8.8.8"     # External destination (internet)

    # GTP tunnel identifiers
    teid_ul: int = 0x12345678        # Uplink TEID (eNB -> SGW)
    teid_dl: int = 0x87654321        # Downlink TEID (SGW -> eNB)

    # Test parameters
    payload_size: int = 64
    capture_timeout: float = 2.0
    num_test_packets: int = 5

    def __post_init__(self):
        if self.tx_iface == self.rx_iface:
            raise ValueError("TX and RX interfaces must be different")
        if self.teid_ul == 0 or self.teid_dl == 0:
            raise ValueError("TEIDs cannot be zero (reserved)")


def compute_ue_mark(ue_ip_int: int) -> int:
    """
    Python port of compute_ue_mark() from ebpf_gtp_decap.c.

    Used to validate the metadata_mark the eBPF program writes to the session.
    """
    safe_mark = ue_ip_int & 0x7FFFFFFE
    if safe_mark == 0x7FFFFFFF or safe_mark == 0:
        safe_mark = (ue_ip_int >> 8) | 0x12345600
    if safe_mark < 0x10000000:
        safe_mark |= 0x12000000
    return safe_mark


def load_production_source(path: str) -> str:
    """Read production eBPF source file"""
    with open(path, 'r') as f:
        return f.read()


def populate_session(bpf, ue_ip_str, enb_ip_str, teid_ul_in, teid_dl_out,
                     s1u_ifindex=0, sgi_ifindex=0, ovs_ifindex=0,
                     mac_src=None, mac_dst=None, bearer_id=5, qfi=9,
                     active=True):
    """
    Add a UE session to the production ue_session_map.

    Uses BCC auto-generated ctypes to match the exact struct layout.
    All IPs stored in host byte order (matching production code).
    """
    session_map = bpf.get_table("ue_session_map")

    key = session_map.Key()
    key.ue_ip = ip_to_host_order(ue_ip_str)

    val = session_map.Leaf()
    val.enb_ip = ip_to_host_order(enb_ip_str)
    val.teid_ul_in = teid_ul_in
    val.teid_ul_out = teid_ul_in  # Same for testing
    val.teid_dl_in = teid_dl_out
    val.teid_dl_out = teid_dl_out
    val.s1u_ifindex = s1u_ifindex
    val.sgi_ifindex = sgi_ifindex
    val.ovs_ifindex = ovs_ifindex
    val.bearer_id = bearer_id
    val.session_flags = 1 if active else 0
    val.qfi = qfi

    if mac_src:
        mac_bytes = mac_str_to_bytes(mac_src)
        for i in range(6):
            val.ul_mac_src[i] = mac_bytes[i]

    if mac_dst:
        mac_bytes = mac_str_to_bytes(mac_dst)
        for i in range(6):
            val.ul_mac_dst[i] = mac_bytes[i]

    session_map[key] = val
    return key, val


def set_config(bpf, key_id, value):
    """Set a value in the production config_map"""
    config_map = bpf.get_table("config_map")
    k = config_map.Key()
    k.key = key_id
    v = config_map.Leaf()
    v.value = value
    config_map[k] = v


def read_stat(bpf, counter_id):
    """Read a counter from stats_map"""
    stats = bpf.get_table("stats_map")
    import ctypes
    key = ctypes.c_uint32(counter_id)
    try:
        val = stats[key]
        return val.value if hasattr(val, 'value') else int(val)
    except KeyError:
        return 0


def attach_tc(bpf, iface, func_name, direction):
    """
    Attach eBPF function to TC hook via pyroute2 netlink.

    Matches production code (ebpf_gtp_manager.py _attach_tc_program).
    The tc CLI can't pass raw fds — only netlink can.
    """
    from bcc import BPF as _BPF
    from pyroute2 import IPRoute

    fn = bpf.load_func(func_name, _BPF.SCHED_CLS)
    ipr = IPRoute()
    try:
        ifindex = ipr.link_lookup(ifname=iface)[0]
        parent = "ffff:fff2" if direction == "ingress" else "ffff:fff3"
        ipr.tc("add-filter", "bpf", ifindex, ":1",
               fd=fn.fd, name=fn.name, parent=parent,
               classid=1, direct_action=True)
    finally:
        ipr.close()


def detach_tc(iface, direction):
    """Remove all TC filters from an interface direction."""
    subprocess.run(
        ['tc', 'filter', 'del', 'dev', iface, direction],
        capture_output=True, check=False,
    )


def has_map(bpf, map_name):
    """Check if a BPF map exists (compatible with all BCC versions)."""
    try:
        bpf.get_table(map_name)
        return True
    except Exception:
        return False


def golden_decap(gtp_packet_bytes, session_mac_src, session_mac_dst):
    """Build expected decap output using Scapy as independent oracle.

    The eBPF decap handler should produce:
      Ether(session MACs, type=0x0800) / inner_ip_packet

    Args:
        gtp_packet_bytes: Raw bytes of the input GTP-U packet (with Ether)
        session_mac_src: MAC address the handler writes as Ether src
        session_mac_dst: MAC address the handler writes as Ether dst

    Returns:
        Expected output packet as raw bytes
    """
    from scapy.all import Ether, IP, Raw
    from scapy.contrib.gtp import GTP_U_Header

    pkt = Ether(gtp_packet_bytes)
    # Extract inner IP: everything after GTP_U_Header
    gtp = pkt[GTP_U_Header]
    inner_ip_bytes = bytes(gtp.payload)

    # Build expected output
    expected = (
        Ether(src=session_mac_src, dst=session_mac_dst, type=0x0800) /
        Raw(load=inner_ip_bytes)
    )
    return bytes(expected)


def golden_encap(ip_packet_bytes, session_mac_src, session_mac_dst,
                 sgw_ip, enb_ip, teid_dl, qfi=9):
    """Build expected encap output using Scapy as independent oracle.

    The eBPF encap handler should produce:
      Ether(session MACs) / IP(sgw->enb) / UDP(2152) /
      GTP_U_Header(teid, E=1) / optional(seq,npdu,next_ext=0x85) /
      PDU_Session_Container(qfi) / original_inner_ip

    Args:
        ip_packet_bytes: Raw bytes of inner IP packet (WITHOUT Ether header)
        session_mac_src: MAC from session ul_mac_src
        session_mac_dst: MAC from session ul_mac_dst
        sgw_ip: SGW IP (outer source)
        enb_ip: eNB IP (outer destination)
        teid_dl: Downlink TEID
        qfi: QoS Flow Identifier

    Returns:
        Expected output packet as raw bytes
    """
    from scapy.all import Ether, IP, UDP, Raw

    inner_len = len(ip_packet_bytes)

    # GTP extension: 4 bytes optional fields + 4 bytes PDU Session Container
    gtp_ext = bytes([
        0x00, 0x00,  # Seq number
        0x00,        # N-PDU
        0x85,        # Next ext: PDU Session Container
        0x01,        # Ext length (1 = 4 bytes)
        0x10,        # PDU Type: DL PDU SESSION INFO
        qfi & 0x3F,  # QFI
        0x00,        # Next ext: none
    ])

    gtp_payload_len = len(gtp_ext) + inner_len  # 8 + inner

    # GTP header: flags=0x34 (V1, PT, E), type=0xFF
    gtp_hdr = struct.pack('!BBHI',
                          0x34,           # flags
                          0xFF,           # type (T-PDU)
                          gtp_payload_len,
                          teid_dl)

    udp_payload = gtp_hdr + gtp_ext + ip_packet_bytes
    udp_len = 8 + len(udp_payload)
    total_ip_len = 20 + udp_len

    expected = (
        Ether(src=session_mac_src, dst=session_mac_dst, type=0x0800) /
        IP(src=sgw_ip, dst=enb_ip, ttl=64, id=0, flags='DF',
           proto=17, len=total_ip_len) /
        UDP(sport=2152, dport=2152, len=udp_len, chksum=0) /
        Raw(load=gtp_hdr + gtp_ext + ip_packet_bytes)
    )

    # Let Scapy calculate the IP checksum
    expected_bytes = bytes(expected)
    return expected_bytes
