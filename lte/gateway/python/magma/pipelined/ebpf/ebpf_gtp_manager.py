"""
Copyright 2025 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Author: Nitin Rajput (coRAN LABS)

eBPF GTP Manager

Manages eBPF-based GTP-U encapsulation/decapsulation as a replacement
for the kernel GTP module while maintaining OVS integration.

This manager:
1. Loads and manages eBPF GTP programs
2. Creates and manages veth pair for OVS integration
3. Handles UE session management via BPF maps
4. Provides statistics collection
5. Maintains compatibility with existing OpenFlow pipeline
"""

import logging
import socket
import struct
import subprocess
import time
import threading
import select
import mmap
import os
import re
import ctypes as ct
import sys
from typing import Dict, List, Optional, Tuple, Any

from bcc import BPF


class BPFMapWrapper:
    """Wrapper for BPF map operations via bpftool when using external TC maps"""
    def __init__(self, map_id):
        self.map_id = map_id
        self._use_bpftool = True

    def __setitem__(self, key, value):
        """Add/update entry in map via bpftool"""
        if hasattr(key, 'ue_ip'):
            ue_ip_int = key.ue_ip
        else:
            ue_ip_int = key

        # Build value bytes
        value_bytes = bytearray(160)  # ue_session_info size
        if hasattr(value, 'enb_ip'):
            import struct
            struct.pack_into('<I', value_bytes, 0, value.enb_ip)
            struct.pack_into('<I', value_bytes, 4, value.teid_ul_in)
            struct.pack_into('<I', value_bytes, 8, value.teid_ul_out)
            struct.pack_into('<I', value_bytes, 12, value.teid_dl_in)
            struct.pack_into('<I', value_bytes, 16, value.teid_dl_out)
            struct.pack_into('<I', value_bytes, 20, value.s1u_ifindex)
            value_bytes[48] = value.bearer_id if hasattr(value, 'bearer_id') else 5
            struct.pack_into('<I', value_bytes, 88, 1)  # session_flags = active

        # Convert to hex
        import struct
        key_hex = ' '.join([f'0x{b:02x}' for b in struct.pack('<I', ue_ip_int)])
        value_hex = ' '.join([f'{b:02x}' for b in value_bytes])

        # Use bpftool
        cmd = f"sudo bpftool map update id {self.map_id} key {key_hex} value hex {value_hex}"
        result = subprocess.run(cmd, shell=True, capture_output=True)
        if result.returncode != 0:
            raise Exception(f"Failed to update map: {result.stderr.decode()}")

    def __getitem__(self, key):
        """Get entry from map via bpftool"""
        if hasattr(key, 'ue_ip'):
            ue_ip_int = key.ue_ip
        else:
            ue_ip_int = key

        import struct
        key_hex = ' '.join([f'0x{b:02x}' for b in struct.pack('<I', ue_ip_int)])
        cmd = f"sudo bpftool map lookup id {self.map_id} key {key_hex}"
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)

        if result.returncode == 0:
            # Parse the output and return a mock object
            class MapValue:
                pass
            return MapValue()
        else:
            raise KeyError(f"Key {ue_ip_int} not found in map")

    def __delitem__(self, key):
        """Delete entry from map via bpftool"""
        if hasattr(key, 'ue_ip'):
            ue_ip_int = key.ue_ip
        else:
            ue_ip_int = key

        import struct
        key_hex = ' '.join([f'0x{b:02x}' for b in struct.pack('<I', ue_ip_int)])
        cmd = f"sudo bpftool map delete id {self.map_id} key {key_hex}"
        subprocess.run(cmd, shell=True, capture_output=True)

    def __contains__(self, key):
        """Check if key exists in map"""
        try:
            self.__getitem__(key)
            return True
        except KeyError:
            return False

    def Key(self, ue_ip):
        """Create a key object"""
        class MapKey:
            def __init__(self, ip):
                self.ue_ip = ip
        return MapKey(ue_ip)

    def Leaf(self, *args, **kwargs):
        """Create a value object"""
        class MapValue:
            def __init__(self):
                for k, v in kwargs.items():
                    setattr(self, k, v)
        val = MapValue()
        return val


class ExternalBPFMap:
    """Wrapper to manipulate external BPF maps via bpftool"""

    def __init__(self, map_id):
        self.map_id = map_id
        self.logger = logging.getLogger('pipelined.ebpf_gtp')

    def update(self, key_int, value_struct):
        """Update map entry using bpftool"""
        try:
            # Convert key to bytes (little-endian for x86)
            key_bytes = key_int.to_bytes(4, 'little')
            key_hex = ' '.join(f'0x{b:02x}' for b in key_bytes)

            # Convert value struct to bytes
            if hasattr(value_struct, '__bytes__'):
                value_bytes = bytes(value_struct)
            elif isinstance(value_struct, bytes):
                value_bytes = value_struct
            else:
                # Pack the ue_session_info structure
                value_bytes = bytes(160)  # Default

            value_hex = ' '.join(f'0x{b:02x}' for b in value_bytes)

            cmd = f'bpftool map update id {self.map_id} key {key_hex} value {value_hex}'
            result = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=5)

            if result.returncode != 0:
                self.logger.error(f"Failed to update map: {result.stderr}")
                return False

            self.logger.debug(f"Updated external map {self.map_id} with key 0x{key_int:08x}")
            return True

        except Exception as e:
            self.logger.error(f"Error updating external map: {e}")
            return False

    def __setitem__(self, key, value):
        """Allow map[key] = value syntax for BCC compatibility"""
        if hasattr(key, 'value'):
            key_int = key.value
        else:
            key_int = key
        return self.update(key_int, value)

    def Key(self, value):
        """Create a key object for BCC compatibility"""
        class KeyObj:
            def __init__(self, v):
                self.value = v
        return KeyObj(value)

    def Leaf(self):
        """Create a value object for BCC compatibility"""
        # Return a ctypes structure matching ue_session_info size
        class Value(ct.Structure):
            _fields_ = [("data", ct.c_ubyte * 160)]
        return Value
from magma.pipelined.gw_mac_address import get_mac_by_ip4
from magma.pipelined.ifaces import get_mac_address_from_iface
from pyroute2 import IPRoute, NetlinkError
from magma.pipelined.ebpf.ebpf_manager import EbpfManager
from magma.pipelined.ebpf.ebpf_utils import INGRESS, EGRESS
from magma.pipelined.bridge_util import BridgeTools

LOG = logging.getLogger("pipelined.ebpf_gtp")

# Constants
BPF_FS_PATH = "/sys/fs/bpf"
EBPF_GTP_DIR = "/var/opt/magma/ebpf"
GTP_MAPS_DIR = "/sys/fs/bpf/gtp_maps"  # Directory for pinned maps
MAX_RETRY_ATTEMPTS = 3  # Maximum retry attempts for eBPF loading
RETRY_DELAY = 2  # Delay between retries in seconds

# Define the C structures using ctypes to match the eBPF program exactly
class UeSessionKey(ct.Structure):
    _fields_ = [
        ('ue_ip', ct.c_uint32),  # __be32 in network byte order
    ]

class UeSessionInfo(ct.Structure):
    _fields_ = [
        ('enb_ip', ct.c_uint32),           # __be32 enb_ip
        ('teid_ul_in', ct.c_uint32),       # __u32 teid_ul_in
        ('teid_ul_out', ct.c_uint32),      # __u32 teid_ul_out
        ('teid_dl_in', ct.c_uint32),       # __u32 teid_dl_in
        ('teid_dl_out', ct.c_uint32),      # __u32 teid_dl_out
        ('s1u_ifindex', ct.c_uint32),      # __u32 s1u_ifindex
        ('bearer_id', ct.c_uint32),        # __u32 bearer_id
        ('ul_bytes', ct.c_uint64),         # __u64 ul_bytes
        ('dl_bytes', ct.c_uint64),         # __u64 dl_bytes
        ('ul_packets', ct.c_uint64),       # __u64 ul_packets
        ('dl_packets', ct.c_uint64),       # __u64 dl_packets
        ('last_seen', ct.c_uint64),        # __u64 last_seen
        ('session_flags', ct.c_uint32),    # __u32 session_flags
        ('imsi', ct.c_ubyte * 16),         # __u8 imsi[16]
        ('imsi_len', ct.c_uint32),         # __u32 imsi_len
        ('_pad1', ct.c_ubyte * 3),         # padding for alignment
        ('direction', ct.c_uint8),         # __u8 direction
        ('_pad2', ct.c_ubyte * 2),         # padding for alignment
        ('original_port', ct.c_uint32),    # __u32 original_port
        ('reserved', ct.c_ubyte * 3),      # __u8 reserved[3]
        ('_pad3', ct.c_ubyte),             # padding for alignment
        ('metadata_mark', ct.c_uint32),    # __u32 metadata_mark
        # No trailing padding needed - struct size is 160 bytes
    ]
GTP_DECAP_PROG = "ebpf_gtp_decap.c"
GTP_ENCAP_PROG = "ebpf_gtp_encap.c"
GTP_VETH0_MARK_PROG = "ebpf_gtp_veth0_mark.c"
GTP_VETH_OVS = "gtp_veth0"
GTP_VETH_EBPF = "gtp_veth1"
GTP_TUNNEL = "gtp0"
GTP_PORT_NO = 2152

# Map names (must match eBPF programs)
UE_SESSION_MAP = "ue_session_map"
CONFIG_MAP = "config_map"
STATS_MAP = "stats_map"

# Configuration keys (must match eBPF programs)
CONFIG_S1U_IFINDEX = 0
CONFIG_SGI_IFINDEX = 1
CONFIG_OVS_IFINDEX = 2
CONFIG_DEBUG_LEVEL = 3
CONFIG_SGI_IP = 4
CONFIG_EBPF_VETH_IFINDEX = 5

# Statistics keys - Synchronized with eBPF programs
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
STATS_MAX_COUNTERS = 32

# Struct formats (must match C layout exactly)
SESSION_STRUCT_FMT = '!8I6s6s2I5QI16sIQB3xIIBB2xI3sI9x'
SESSION_STRUCT_SIZE = struct.calcsize(SESSION_STRUCT_FMT)
SESSION_STRUCT_FMT_LEGACY = '!8I6s6s2I5QI16sIQBIIBBI3s'
SESSION_STRUCT_SIZE_LEGACY = struct.calcsize(SESSION_STRUCT_FMT_LEGACY)


def _compute_ue_mark_from_ip_int(ue_ip_int: int) -> int:
    """Compute UE mark matching eBPF/classifier logic."""

    safe_mark = ue_ip_int & 0x7FFFFFFE

    if safe_mark in (0x7FFFFFFF, 0):
        safe_mark = (ue_ip_int >> 8) | 0x12345600

    if safe_mark < 0x10000000:
        safe_mark |= 0x12000000

    return safe_mark


class EbpfGtpManager:
    """
    eBPF GTP Manager - Replaces kernel GTP module with eBPF implementation
    """

    def __init__(self, config: Dict):
        """
        Initialize eBPF GTP Manager
        
        Args:
            config: Pipelined configuration dictionary
        """
        self.config = config
        self.enabled = False
        # Detect the actual S1U interface where GTP packets arrive
        self.s1u_iface = self._detect_s1u_interface(config)
        self.sgi_iface = config.get('nat_iface', 'eth0')
        self.ovs_bridge = config.get('bridge_name', 'gtp_br0')
        self.debug_level = config.get('debug_level', 0)
        
        # Control plane echo handling (CRITICAL for kernel module replacement)
        self._echo_socket = None
        self._echo_thread = None
        self._stop_echo_handler = False
        
        # eBPF programs and maps
        self.decap_bpf = None
        self.encap_bpf = None
        self.combined_bpf = None
        self.veth0_mark_bpf = None
        self.ue_session_map = None
        self.config_map = None
        self.stats_map = None
        self._maps_pinned = False  # Track if maps are pinned
        self._maps_initialized = False
        self._use_external_maps = False  # Using existing TC maps
        self._external_map_id = None  # ID of external ue_session_map
        self._bpf_maps_dir = None
        
        # Interface management
        self.ipr = IPRoute()
        self.s1u_ifindex = None
        # Phase 1: self.sgi_ifindex = None
        # Phase 1: self.ovs_ifindex = None
        
        # Session tracking
        
        # Interface refresh callbacks (CRITICAL for dynamic port handling)
        self._interface_callbacks = []
        
        LOG.info("eBPF GTP Manager initialized - S1U: %s, SGi: %s", 
                 self.s1u_iface, self.sgi_iface)

    def _detect_s1u_interface(self, config: Dict) -> str:
        """
        For TC+TC approach, always return the actual interface where packets arrive
        
        Args:
            config: Pipelined configuration
            
        Returns:
            Interface name for TC attachment
        """
        # Use configured interface from pipelined.yml
        configured_iface = config.get('enodeb_iface', 'eth1')
        
        LOG.info("Using %s for TC eBPF attachment (from config)", configured_iface)

        return configured_iface

    def _detect_attached_tc_maps(self) -> Dict[str, int]:
        """
        Detect BPF maps from currently attached TC programs

        This ensures pipelined writes to the SAME maps that the attached eBPF programs read from,
        avoiding the mismatch where pipelined writes to new maps but TC programs read from old maps.

        Returns:
            Dictionary mapping map names to map IDs, or None if detection fails
        """
        try:
            import subprocess
            import re

            # Check gtp_veth0 egress for ENCAP program (most critical for downlink)
            result = subprocess.run(
                ['tc', 'filter', 'show', 'dev', GTP_VETH_OVS, 'egress'],
                capture_output=True, text=True, timeout=5
            )

            # Extract program ID from tc filter output
            # Example: "filter protocol all pref 49152 bpf chain 0 handle 0x1 ... id 356 ..."
            prog_id_match = re.search(r'\bid\s+(\d+)\b', result.stdout)

            if not prog_id_match:
                LOG.warning("No TC program found on %s egress", GTP_VETH_OVS)
                return None

            prog_id = int(prog_id_match.group(1))
            LOG.info(f"Found attached TC ENCAP program ID: {prog_id}")

            # Get map IDs from this program using bpftool
            result = subprocess.run(
                ['bpftool', 'prog', 'show', 'id', str(prog_id)],
                capture_output=True, text=True, timeout=5
            )

            # Extract map_ids line
            # Example: "xlated 4768B  jited 2703B  memlock 8192B  map_ids 413,415,414"
            map_ids_match = re.search(r'map_ids\s+([\d,]+)', result.stdout)

            if not map_ids_match:
                LOG.warning(f"No map_ids found for program {prog_id}")
                return None

            map_ids_str = map_ids_match.group(1)
            map_ids = [int(x) for x in map_ids_str.split(',')]

            LOG.info(f"Program {prog_id} uses map IDs: {map_ids}")

            # Identify which map is which by checking map names
            maps = {}
            for map_id in map_ids:
                result = subprocess.run(
                    ['bpftool', 'map', 'show', 'id', str(map_id)],
                    capture_output=True, text=True, timeout=5
                )

                # Extract map name
                # Example: "415: hash  name ue_session_map  flags 0x0"
                name_match = re.search(r'name\s+(\w+)', result.stdout)
                if name_match:
                    map_name = name_match.group(1)

                    # Map to our expected names
                    if 'ue_session' in map_name:
                        maps['ue_session'] = map_id
                    elif 'config' in map_name:
                        maps['config'] = map_id
                    elif 'stats' in map_name:
                        maps['stats'] = map_id

            if 'ue_session' in maps:
                LOG.info(f"Detected attached TC maps: {maps}")
                return maps
            else:
                LOG.warning("ue_session_map not found in attached program's maps")
                return None

        except Exception as e:
            LOG.warning(f"Failed to detect attached TC maps: {e}")
            return None

    def _cleanup_stale_xdp_programs(self):
        """Clean up any stale XDP programs from previous runs"""
        # XDP cleanup disabled - using TC-only approach
        # try:
        #     LOG.info("Checking for stale XDP programs...")
        #     
        #     # List all interfaces that might have XDP programs
        #     interfaces_to_check = [self.s1u_iface, GTP_VETH_EBPF, 'eth1', 'eth0']
        #     
        #     for iface in interfaces_to_check:
        #         try:
        #             # Check if interface exists
        #             if not os.path.exists(f"/sys/class/net/{iface}"):
        #                 continue
        #                 
        #             # Use ip link show to check for XDP programs
        #             result = subprocess.run(['ip', 'link', 'show', iface], 
        #                                   capture_output=True, text=True)
        #             
        #             if 'xdp' in result.stdout.lower():
        #                 LOG.warning(f"Found XDP program on {iface}, removing...")
        #                 # Remove XDP program
        #                 subprocess.run(['ip', 'link', 'set', 'dev', iface, 'xdp', 'off'],
        #                              capture_output=True, check=False)
        #                 LOG.info(f"Removed XDP program from {iface}")
        #                 
        #         except Exception as e:
        #             LOG.debug(f"Error checking {iface} for XDP programs: {e}")
        #     
        #     # Clean up any pinned maps from previous runs
        #     maps_dir = f"{BPF_FS_PATH}/gtp_maps"
        #     if os.path.exists(maps_dir):
        #         try:
        #             for map_file in os.listdir(maps_dir):
        #                 map_path = os.path.join(maps_dir, map_file)
        #                 os.unlink(map_path)
        #                 LOG.debug(f"Removed stale pinned map: {map_file}")
        #         except Exception as e:
        #             LOG.debug(f"Error cleaning up pinned maps: {e}")
        #             
        #     LOG.info("Stale XDP program cleanup completed")
        #     
        # except Exception as e:
        #     LOG.warning(f"Error during XDP cleanup: {e}")
        #     # Non-fatal error - continue with initialization
        # Clean up old pinned maps if they exist
        self._cleanup_pinned_maps()
    
    def _cleanup_pinned_maps(self):
        """Clean up old pinned BPF maps"""
        if os.path.exists(GTP_MAPS_DIR):
            try:
                import shutil
                shutil.rmtree(GTP_MAPS_DIR)
                LOG.info("Cleaned up old pinned maps at %s", GTP_MAPS_DIR)
            except Exception as e:
                LOG.warning(f"Failed to cleanup pinned maps: {e}")

    def _bind_gtp_port(self):
        """DEPRECATED: Removed UDP socket binding - using TC eBPF instead"""
        # LOG.info("UDP socket binding skipped - using AF_XDP for packet processing")
        LOG.info("UDP socket binding skipped - using TC eBPF for packet processing")
        pass

    def _handle_gtp_packets(self):
        """Handle GTP packets received on UDP socket (temporary until TC eBPF is ready)"""
        # XDP/AF_XDP code commented out - using TC eBPF
        # try:
        #     while True:
        #         try:
        #             # Receive packet from UDP socket
        #             data, addr = self.gtp_socket.recvfrom(2048)
        #             
        #             # Log packet info (temporary - will be replaced by AF_XDP processing)
        #             LOG.debug("Received GTP packet from %s, length: %d", addr, len(data))
        #             
        #             # For now, just acknowledge receipt - AF_XDP will handle processing
        #             # TODO: Replace with AF_XDP packet processing
        #             
        #         except socket.error as e:
        #             if e.errno != 9:  # Ignore "Bad file descriptor" on socket close
        #                 LOG.debug("Socket error in GTP handler: %s", e)
        #             break
        #         except Exception as e:
        #             LOG.debug("Error in GTP packet handler: %s", e)
        #             
        # except Exception as e:
        #     LOG.error("GTP packet handler thread failed: %s", e)
        pass

    def _setup_af_xdp_socket(self) -> bool:
        """Setup AF_XDP socket for high-performance packet processing"""
        # AF_XDP code commented out - using TC eBPF
        # try:
        #     if not self.af_xdp_enabled:
        #         LOG.info("AF_XDP disabled, falling back to kernel XDP")
        #         return False
        #         
        #     # Import our AF_XDP helper with proper error handling
        #     try:
        #         from magma.pipelined.ebpf.af_xdp_helper import AF_XDPSocket
        #     except ImportError as ie:
        #         LOG.error("AF_XDP helper module not found: %s", ie)
        #         LOG.error("Check if af_xdp_helper.py is in the correct path")
        #         return False
        #     
        #     # Create AF_XDP socket for S1-U interface (eth1)
        #     self.af_xdp_socket = AF_XDPSocket(
        #         ifname=self.s1u_iface,
        #         queue_id=0,
        #         frame_size=4096,
        #         num_frames=4096,
        #         ring_size=2048
        #     )
        #     
        #     # Create and configure the AF_XDP socket
        #     if not self.af_xdp_socket.create():
        #         LOG.error("Failed to create AF_XDP socket")
        #         return False
        #         
        #     LOG.info("AF_XDP socket successfully created for %s", self.s1u_iface)
        #     
        #     return True
        #     
        # except Exception as e:
        #     LOG.error("Failed to setup AF_XDP socket: %s", e)
        #     LOG.warning("AF_XDP setup failed, will use kernel XDP instead")
        #     if self.af_xdp_socket:
        #         self.af_xdp_socket.close()
        #         self.af_xdp_socket = None
        #     # AF_XDP failure is not fatal - we can still use kernel XDP
        #     self.af_xdp_enabled = False
        #     return False
        return False

    def _handle_af_xdp_packets(self):
        """Handle packets from AF_XDP socket (EUPF-style processing)"""
        # AF_XDP code commented out - using TC eBPF
        # try:
        #     while self.af_xdp_socket:
        #         try:
        #             # Receive packet from AF_XDP socket
        #             data = self.af_xdp_socket.recv_packet()
        #             
        #             if data:
        #                 # Process GTP packet using our enhanced processing
        #                 self._process_gtp_packet(data)
        #             else:
        #                 # No packet available, sleep briefly
        #                 time.sleep(0.001)  # 1ms sleep
        #                 
        #         except Exception as e:
        #             LOG.debug("Error in AF_XDP packet handler: %s", e)
        #             time.sleep(0.001)
        #             
        # except Exception as e:
        #     LOG.error("AF_XDP packet handler thread failed: %s", e)
        pass
            
    def _process_gtp_packet(self, packet_data: bytes):
        """Process GTP packet received from eBPF (for debugging purposes)"""
        try:
            # Parse Ethernet header
            if len(packet_data) < 14:
                return
                
            eth_header = packet_data[:14]
            eth_type = struct.unpack(">H", eth_header[12:14])[0]
            
            # Check if IPv4
            if eth_type != 0x0800:
                return
                
            # Parse IP header
            ip_header = packet_data[14:34]
            if len(ip_header) < 20:
                return
                
            ip_version = (ip_header[0] >> 4) & 0xF
            ip_protocol = ip_header[9]
            
            if ip_version != 4 or ip_protocol != 17:  # UDP
                return
                
            # Parse UDP header
            udp_start = 14 + ((ip_header[0] & 0xF) * 4)
            udp_header = packet_data[udp_start:udp_start + 8]
            
            if len(udp_header) < 8:
                return
                
            _, dst_port, udp_len, _ = struct.unpack(">HHHH", udp_header)
            
            # Check if GTP-U port
            if dst_port != 2152:
                return
                
            # Parse GTP header
            gtp_start = udp_start + 8
            gtp_header = packet_data[gtp_start:gtp_start + 8]
            
            if len(gtp_header) < 8:
                return
                
            gtp_flags, gtp_type, gtp_length, gtp_teid = struct.unpack(">BBHI", gtp_header)
            
            # Validate GTP header
            if (gtp_flags >> 5) != 1 or gtp_type != 0xFF:  # GTP v1 T-PDU
                return
                
            # Extract inner IP packet
            inner_ip_start = gtp_start + 8
            inner_ip = packet_data[inner_ip_start:]
            
            if len(inner_ip) < 20:
                return
                
            # Get UE IP (destination of inner packet)
            ue_ip_bytes = inner_ip[16:20]
            ue_ip = socket.inet_ntoa(ue_ip_bytes)

            # Construct decapsulated packet (Ethernet + Inner IP)
            decap_packet = bytearray()
            
            # New Ethernet header
            decap_packet.extend(b'\x02\x00\x00\x00\x00\x01')  # dst MAC
            decap_packet.extend(b'\x02\x00\x00\x00\x00\x02')  # src MAC
            decap_packet.extend(b'\x08\x00')  # IPv4 EtherType
            
            # Inner IP packet
            decap_packet.extend(inner_ip)
            
            # Send decapsulated packet to OVS via veth interface
            
            # Update statistics
            if self.stats_map:
                self._update_stats_counters(len(packet_data), True)
            
            LOG.debug("Successfully processed GTP packet for UE %s", ue_ip)
            
        except Exception as e:
            LOG.debug("Error processing GTP packet: %s", e)
            
    def _start_echo_handler(self) -> bool:
        """
        Start GTP control plane echo handler (CRITICAL for kernel module replacement)
        This replaces the kernel module's echo request/response handling
        """
        try:
            # Create raw socket for GTP control packets
            self._echo_socket = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_UDP)
            self._echo_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            
            # Bind to GTP-U port for echo handling
            gtp_bind_ip = self._get_interface_ip(self.config.get('enodeb_iface', 'eth1'))
            if not gtp_bind_ip:
                gtp_bind_ip = '0.0.0.0'
            
            self._echo_socket.bind((gtp_bind_ip, GTP_PORT_NO))
            LOG.info("GTP echo handler bound to %s:%d", gtp_bind_ip, GTP_PORT_NO)
            
            # Start echo handler thread
            self._stop_echo_handler = False
            self._echo_thread = threading.Thread(target=self._handle_gtp_echo_packets, daemon=True)
            self._echo_thread.start()
            
            LOG.info("GTP control plane echo handler started")
            return True
            
        except Exception as e:
            LOG.error("Failed to start GTP echo handler: %s", e)
            if self._echo_socket:
                self._echo_socket.close()
                self._echo_socket = None
            return False
    
    def _handle_gtp_echo_packets(self):
        """
        Handle GTP echo requests and responses (replaces kernel module functionality)
        This is CRITICAL for maintaining GTP tunnel state with eNBs
        """
        LOG.info("GTP echo handler thread started")
        
        while not self._stop_echo_handler:
            try:
                # Use select with timeout to allow clean shutdown
                ready = select.select([self._echo_socket], [], [], 1.0)
                if not ready[0]:
                    continue
                    
                # Receive GTP packet
                data, addr = self._echo_socket.recvfrom(2048)
                if len(data) < 8:  # Minimum GTP header
                    continue
                
                # Parse GTP header
                gtp_flags = data[0]
                gtp_type = data[1]
                gtp_length = struct.unpack(">H", data[2:4])[0]
                gtp_teid = struct.unpack(">I", data[4:8])[0]
                
                LOG.debug("Received GTP packet: type=%d, teid=%d from %s", gtp_type, gtp_teid, addr[0])
                
                # Handle GTP Echo Request (type 1)
                if gtp_type == 1:
                    self._send_echo_response(addr[0], data)
                    LOG.debug("Sent GTP echo response to %s", addr[0])
                    
                    # Update statistics
                    if self.stats_map:
                        self._update_stats_counters(len(data), False)
                
                # Handle GTP Echo Response (type 2)
                elif gtp_type == 2:
                    LOG.debug("Received GTP echo response from %s", addr[0])
                    # Update keepalive state for this peer
                    
            except socket.timeout:
                continue
            except Exception as e:
                if not self._stop_echo_handler:
                    LOG.debug("Error in GTP echo handler: %s", e)
                    time.sleep(0.1)  # Brief pause before retry
        
        LOG.info("GTP echo handler thread stopped")
    
    def _send_echo_response(self, peer_ip: str, request_data: bytes):
        """
        Send GTP echo response (replicates kernel module behavior)
        """
        try:
            if len(request_data) < 8:
                return
                
            # Extract sequence number if present
            seq_num = 0
            if len(request_data) >= 12 and (request_data[0] & 0x02):  # S flag set
                seq_num = struct.unpack(">H", request_data[8:10])[0]
            
            # Build echo response
            response = bytearray()
            response.append(0x32)  # GTP v1, PT=1, S=1  
            response.append(0x02)  # Echo response
            response.extend(struct.pack(">H", 6))  # Length (header + seq + spare)
            response.extend(struct.pack(">I", 0))  # TEID = 0 for echo
            
            # Add sequence number and spare bytes
            response.extend(struct.pack(">H", seq_num))
            response.extend(b'\x00\x00')  # Spare bytes
            
            # Send response
            response_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
            response_socket.sendto(bytes(response), (peer_ip, GTP_PORT_NO))
            response_socket.close()
            
            LOG.debug("Sent GTP echo response to %s (seq=%d)", peer_ip, seq_num)
            
        except Exception as e:
            LOG.debug("Failed to send echo response: %s", e)
            
    def _update_stats_counters(self, packet_size: int, is_uplink: bool):
        """Update statistics counters in BPF maps"""
        try:
            import ctypes
            
            if is_uplink:
                # Update UL counters
                ul_pkts_key = ctypes.c_uint32(STATS_UL_PACKETS)
                ul_bytes_key = ctypes.c_uint32(STATS_UL_BYTES)
                
                ul_pkts_val = self.stats_map.get(ul_pkts_key, ctypes.c_uint64(0))
                ul_bytes_val = self.stats_map.get(ul_bytes_key, ctypes.c_uint64(0))
                
                self.stats_map[ul_pkts_key] = ctypes.c_uint64(ul_pkts_val.value + 1)
                self.stats_map[ul_bytes_key] = ctypes.c_uint64(ul_bytes_val.value + packet_size)
                
        except Exception as e:
            LOG.debug("Error updating statistics: %s", e)
            
    def initialize(self) -> bool:
        """
        Initialize eBPF GTP system
        
        Returns:
            True if initialization successful, False otherwise
        """
        try:
            # 0. Clean up any existing TC filters to prevent duplicate maps
            self._cleanup_existing_tc_filters()

            # 1. Setup interfaces
            if not self._setup_interfaces():
                return False

            # 2. Load eBPF programs with retry logic
            if not self._load_ebpf_programs_with_retry():
                LOG.error("Failed to initialize eBPF GTP manager")
                return False

            # 3. Attach TC programs
            if not self._attach_programs():
                return False
            LOG.info("TC eBPF programs attached successfully")
            
            # 4. Restore sessions if available
            if self._maps_initialized:
                self._restore_sessions_from_file()
            
            # 5. Configure BPF maps (enabled for step 3) - CRITICAL for GTP processing
            if not self._maps_initialized:
                LOG.error("CRITICAL: BPF maps were not initialized during program loading")
                LOG.error("This will prevent all GTP tunnel processing - aborting initialization")
                return False

            if self.config_map is None:
                LOG.error("CRITICAL: BPF config_map is None - initialization failed")
                return False

            LOG.info("BPF maps configured successfully")

            # 6. Start control plane echo handler (CRITICAL for GTP compliance)
            if not self._start_echo_handler():
                LOG.warning("Failed to start GTP echo handler - control plane may not work")
                # Non-fatal, continue with data plane
            
            self.enabled = True
            # LOG.info("eBPF GTP Manager successfully initialized (AF_XDP: %s)", af_xdp_success)
            LOG.info("eBPF GTP Manager successfully initialized with TC eBPF")
            return True
            
        except Exception as e:
            LOG.error("Failed to initialize eBPF GTP Manager: %s", e)
            self.cleanup()
            return False

    def _find_existing_tc_maps(self):
        """Find existing TC-attached BPF maps that we can reuse"""
        try:
            # Check if TC has programs attached
            result = subprocess.run(['tc', 'filter', 'show', 'dev', GTP_VETH_EBPF, 'ingress'],
                                  capture_output=True, text=True, timeout=2)
            if 'gtp_decap' not in result.stdout:
                return None

            # Extract program ID from TC output
            import re
            match = re.search(r'id (\d+)', result.stdout)
            if not match:
                return None

            prog_id = match.group(1)

            # Get map IDs for this program
            result = subprocess.run(['sudo', 'bpftool', 'prog', 'show', 'id', prog_id],
                                  capture_output=True, text=True, timeout=2)

            # Extract map IDs
            match = re.search(r'map_ids ([0-9,]+)', result.stdout)
            if not match:
                return None

            map_ids = [int(x) for x in match.group(1).split(',')]

            # Identify which map is which
            tc_maps = {}
            for map_id in map_ids:
                result = subprocess.run(['sudo', 'bpftool', 'map', 'show', 'id', str(map_id)],
                                      capture_output=True, text=True, timeout=2)
                if 'ue_session' in result.stdout:
                    tc_maps['ue_session'] = map_id
                elif 'stats' in result.stdout:
                    tc_maps['stats'] = map_id
                elif 'config' in result.stdout:
                    tc_maps['config'] = map_id

            if tc_maps:
                LOG.info(f"Found existing TC maps: {tc_maps}")
                return tc_maps

        except Exception as e:
            LOG.debug(f"Error finding TC maps: {e}")

        return None

    def _cleanup_existing_tc_filters(self):
        """Clean up any existing TC filters and qdiscs to prevent duplicate programs/maps"""
        try:
            # First check if we should reuse existing maps
            self._tc_map_ids = self._find_existing_tc_maps()
            if self._tc_map_ids:
                LOG.info(f"Found existing TC maps to reuse: {self._tc_map_ids}")
                # Don't clean up if we're going to reuse
                return

            interfaces_to_clean = [GTP_VETH_EBPF, GTP_VETH_OVS, self.s1u_iface]

            for iface in interfaces_to_clean:
                try:
                    # Remove TC filters if interface exists
                    subprocess.run(['tc', 'filter', 'del', 'dev', iface, 'ingress'],
                                 capture_output=True, timeout=2)
                    subprocess.run(['tc', 'filter', 'del', 'dev', iface, 'egress'],
                                 capture_output=True, timeout=2)
                    # Also remove clsact qdisc to fully detach BPF programs
                    # This ensures old programs and maps are released
                    subprocess.run(['tc', 'qdisc', 'del', 'dev', iface, 'clsact'],
                                 capture_output=True, timeout=2)
                except:
                    pass  # Interface might not exist yet

            # Clear any existing BPF program references to allow garbage collection
            # This helps prevent duplicate map creation
            if hasattr(self, 'decap_bpf'):
                self.decap_bpf = None
            if hasattr(self, 'encap_bpf'):
                self.encap_bpf = None
            if hasattr(self, 'mark_bpf'):
                self.mark_bpf = None
            self._maps_initialized = False

            LOG.debug("Cleaned up existing TC filters, qdiscs, and BPF references")

        except Exception as e:
            LOG.debug(f"Error during TC/BPF cleanup: {e}")

    def _setup_interfaces(self) -> bool:
        """
        Setup veth pair and OVS integration
        
        Returns:
            True if successful, False otherwise
        """
        try:
            # Get interface indices
            self.s1u_ifindex = self._get_ifindex(self.s1u_iface)
            self.sgi_ifindex = self._get_ifindex(self.sgi_iface)
            
            if not self.s1u_ifindex :
                LOG.error("Failed to get interface indices")
                return False
            
            # Remove existing veth pair if it exists
            try:
                self.ipr.link("del", ifname=GTP_VETH_OVS)
            except NetlinkError:
                pass  # Interface doesn't exist, which is fine
            
            # Create veth pair
            self.ipr.link("add", 
                         ifname=GTP_VETH_OVS, 
                         kind="veth", 
                         peer=GTP_VETH_EBPF)
            
            # Bring up both ends of veth pair
            self.ipr.link("set", ifname=GTP_VETH_OVS, state="up")
            self.ipr.link("set", ifname=GTP_VETH_EBPF, state="up")
            
            # Enable mark preservation on veth interfaces for TEID flow matching
            try:
                subprocess.run(
                    ["sysctl", "-w", f"net.ipv4.conf.{GTP_VETH_OVS}.src_valid_mark=1"],
                    capture_output=True, check=True
                )
                subprocess.run(
                    ["sysctl", "-w", f"net.ipv4.conf.{GTP_VETH_EBPF}.src_valid_mark=1"],
                    capture_output=True, check=True
                )
                LOG.info("Enabled mark preservation on veth interfaces for TEID-based flow matching")
            except subprocess.CalledProcessError as e:
                LOG.warning(f"Failed to enable mark preservation (non-fatal): {e}")
                # This is non-fatal as we have IP-based flow fallback
            
            # Get OVS side interface index
            
                LOG.error("Failed to get OVS veth interface index")
                return False
            
            # Add OVS side to bridge with correct port number

            # CRITICAL FIX: Notify classifier to refresh flows for new port
            self._notify_interface_created()

            # Setup dynamic ARP handling for UE subnet
            if not self._setup_ue_subnet_proxy_arp():
                LOG.warning("Failed to setup proxy ARP for UE subnet, manual ARP entries will be used")

            return True
            
        except Exception as e:
            LOG.error("Failed to setup interfaces: %s", e)
            return False
    
    def _ensure_bpf_fs_mounted(self):
        """Ensure BPF filesystem is mounted"""
        try:
            if not os.path.exists(BPF_FS_PATH):
                LOG.error("BPF filesystem not found at %s", BPF_FS_PATH)
                return False

            # Create directory for pinned maps
            if not os.path.exists(GTP_MAPS_DIR):
                os.makedirs(GTP_MAPS_DIR, mode=0o755)
                LOG.info("Created BPF maps directory: %s", GTP_MAPS_DIR)

            return True
        except Exception as e:
            LOG.error(f"Failed to ensure BPF filesystem: {e}")
            return False

    def _find_existing_maps(self):
        """Find existing eBPF maps from TC-attached programs"""
        try:
            # Check if there are existing GTP programs attached to TC
            import subprocess
            result = subprocess.run(['tc', 'filter', 'show', 'dev', GTP_VETH_EBPF, 'ingress'],
                                    capture_output=True, text=True, timeout=5)

            if 'gtp_decap' not in result.stdout:
                LOG.debug("No existing GTP decap program found on TC")
                return None

            # Extract program ID from TC output
            import re
            prog_id_match = re.search(r'id (\d+)', result.stdout)
            if not prog_id_match:
                return None

            prog_id = prog_id_match.group(1)
            LOG.info(f"Found existing eBPF program ID {prog_id} attached to TC")

            # Get map IDs for this program
            result = subprocess.run(['bpftool', 'prog', 'show', 'id', prog_id],
                                    capture_output=True, text=True, timeout=5)

            map_ids_match = re.search(r'map_ids ([0-9,]+)', result.stdout)
            if not map_ids_match:
                return None

            map_ids = map_ids_match.group(1).split(',')
            LOG.info(f"Found existing maps: {map_ids}")

            # Get map details to find the correct ones
            maps = {}
            for map_id in map_ids:
                result = subprocess.run(['bpftool', 'map', 'show', 'id', map_id],
                                        capture_output=True, text=True, timeout=5)

                if 'ue_session_map' in result.stdout:
                    maps['ue_session_map'] = int(map_id)
                elif 'config_map' in result.stdout:
                    maps['config_map'] = int(map_id)
                elif 'stats_map' in result.stdout:
                    maps['stats_map'] = int(map_id)

            if maps:
                LOG.info(f"Found existing maps to reuse: {maps}")
                return maps

            return None

        except Exception as e:
            LOG.debug(f"Error finding existing maps: {e}")
            return None

    def _generate_map_reuse_code(self, existing_maps):
        """Generate BPF code to reuse existing maps"""
        # BCC doesn't directly support map reuse via FD, but we can try pinning
        # For now, return empty string and handle map access differently
        # In a full implementation, we'd use libbpf or modify BCC
        return ""
    
    def _load_ebpf_programs_with_retry(self):
        """Load eBPF programs with retry logic"""
        for attempt in range(1, MAX_RETRY_ATTEMPTS + 1):
            try:
                LOG.info(f"Attempting to load eBPF programs (attempt {attempt}/{MAX_RETRY_ATTEMPTS})")
                if self._load_ebpf_programs():
                    return True
                
                if attempt < MAX_RETRY_ATTEMPTS:
                    LOG.warning(f"eBPF loading failed, retrying in {RETRY_DELAY} seconds...")
                    time.sleep(RETRY_DELAY)
            except Exception as e:
                LOG.error(f"eBPF loading attempt {attempt} failed: {e}")
                if attempt < MAX_RETRY_ATTEMPTS:
                    time.sleep(RETRY_DELAY)
        
        LOG.error("Failed to load eBPF programs after all retries")
        return False

    def _check_and_use_existing_tc_maps(self):
        """Check if eBPF programs are already attached via TC and get their map IDs"""
        try:
            import subprocess
            import re

            # Check if TC has GTP program attached
            result = subprocess.run(['tc', 'filter', 'show', 'dev', GTP_VETH_EBPF, 'ingress'],
                                    capture_output=True, text=True, timeout=5)

            if 'gtp_decap' not in result.stdout:
                return False

            # Extract program ID
            prog_id_match = re.search(r'id (\d+)', result.stdout)
            if not prog_id_match:
                return False

            prog_id = prog_id_match.group(1)
            LOG.info(f"Found existing TC eBPF program ID {prog_id}")

            # Get map IDs for this program
            result = subprocess.run(['bpftool', 'prog', 'show', 'id', prog_id],
                                    capture_output=True, text=True, timeout=5)

            map_ids_match = re.search(r'map_ids ([0-9,]+)', result.stdout)
            if not map_ids_match:
                return False

            map_ids = map_ids_match.group(1).split(',')

            # Find the ue_session_map
            for map_id in map_ids:
                result = subprocess.run(['bpftool', 'map', 'show', 'id', map_id],
                                        capture_output=True, text=True, timeout=5)
                if 'ue_session_map' in result.stdout:
                    self._external_map_id = int(map_id)
                    LOG.info(f"Will use external ue_session_map ID {self._external_map_id}")

                    # Create a wrapper for external map operations
                    self.ue_session_map = ExternalBPFMap(self._external_map_id)
                    return True

            return False

        except Exception as e:
            LOG.debug(f"Error checking existing TC maps: {e}")
            return False

    def _load_ebpf_programs(self) -> bool:
        """
        Load eBPF GTP program using Strategy 1: Single program with multiple entry points
        Compatible with Ubuntu 20.04 and kernel 5.4

        Merges decap and encap programs into a single BPF program where:
        - Maps are naturally shared (single program = single map instance)
        - Both gtp_decap_handler and gtp_encap_handler exist in same program
        - Each handler is attached to different TC hooks

        Returns:
            True if successful, False otherwise
        """
        try:
            # Check if we should skip loading and reuse existing TC programs
            if hasattr(self, '_tc_map_ids') and self._tc_map_ids:
                LOG.info("Skipping BPF program loading - will reuse existing TC programs")
                # Set dummy BPF objects to satisfy later checks
                self.combined_bpf = None
                self.decap_bpf = None
                self.encap_bpf = None
                self.veth0_mark_bpf = None

                # Use wrapper for external maps
                self.ue_session_map = BPFMapWrapper(self._tc_map_ids.get('ue_session'))
                self.config_map = BPFMapWrapper(self._tc_map_ids.get('config'))
                self.stats_map = BPFMapWrapper(self._tc_map_ids.get('stats'))
                self._maps_initialized = True

                LOG.info("Configured to use existing TC maps via bpftool wrapper")
                return True
            # Include paths for BPF compilation (BCC provides its own BPF helpers)
            # Disable debug statements to prevent BPF verifier issues
            cflags = ['-I', EBPF_GTP_DIR, '-DDISABLE_DEBUG', '-O2']
            
            LOG.info("Loading eBPF GTP programs...")
            # Load decapsulation program
            decap_path = f"{EBPF_GTP_DIR}/{GTP_DECAP_PROG}"
            LOG.info(f"Decap program path: {decap_path}")
            
            # Force BCC to recompile by reading the source and passing as text
            # This prevents BCC from using cached bytecode
            with open(decap_path, 'r') as f:
                decap_src = f.read()
            
            LOG.info("Loading eBPF GTP encapsulation program...")
            # Load encapsulation program  
            encap_path = f"{EBPF_GTP_DIR}/{GTP_ENCAP_PROG}"
            LOG.info(f"Encap program path: {encap_path}")
            
            # For proper map sharing, we need to modify the encap program source
            # to reuse the pinned maps instead of creating new ones
            # Read the encap program source
            with open(encap_path, 'r') as f:
                encap_src = f.read()
            
            # Strategy 1: Single BPF program with multiple entry points
            # This approach works reliably with Ubuntu 20.04 and kernel 5.4
            
            LOG.info("Using Strategy 1: Single BPF program with multiple handlers")
            
            # Extract the unique encap handler function from encap source
            # First, we need to extract unique defines from encap that are needed
            unique_defines = []
            
            # Extract GTP_PORT_NO which is only in encap
            if "#define GTP_PORT_NO" in encap_src and "#define GTP_PORT_NO" not in decap_src:
                import re
                gtp_port_match = re.search(r'#define\s+GTP_PORT_NO\s+\d+', encap_src)
                if gtp_port_match:
                    unique_defines.append(gtp_port_match.group(0))
            
            # Add STATS_DOUBLE_ENCAP_AVOIDED
            unique_defines.append("#define STATS_DOUBLE_ENCAP_AVOIDED 21")
            
            # Find the start of gtp_encap_handler function
            handler_start = encap_src.find("int gtp_encap_handler(struct __sk_buff *skb)")
            if handler_start == -1:
                # Try alternative format
                handler_start = encap_src.find("TC_PROG(gtp_encap_handler)")
            
            if handler_start == -1:
                raise Exception("Could not find gtp_encap_handler in encap source")
            
            # Extract from handler start to end of file
            encap_handler_raw = encap_src[handler_start:]
            
            # Remove the duplicate gtp_passthrough_handler if it exists
            # Find where gtp_encap_handler ends and gtp_passthrough_handler might start
            passthrough_start = encap_handler_raw.find("int gtp_passthrough_handler")
            if passthrough_start > 0:
                # Only keep up to the passthrough handler
                encap_handler = encap_handler_raw[:passthrough_start].rstrip()
            else:
                encap_handler = encap_handler_raw.rstrip()
            
            # Create merged source: decap (with all common parts) + unique defines + encap handler
            merged_src = decap_src.rstrip() + "\n\n// Added from encap program\n" + "\n".join(unique_defines) + "\n\n" + encap_handler
            
            LOG.info("Merged programs into single source with both handlers")
            
            # Load the merged program with timestamp to force recompilation
            timestamp = int(time.time())
            cflags_with_timestamp = cflags + [f'-DCOMPILE_TIME={timestamp}']

            # We always load fresh programs since we cleaned up TC filters
            LOG.info("Loading merged BPF program...")
            self.combined_bpf = BPF(text=merged_src, cflags=cflags_with_timestamp)
            LOG.info("Merged BPF program loaded successfully")

            # Set both references to the same BPF object
            # This maintains compatibility with the attach code
            self.decap_bpf = self.combined_bpf
            self.encap_bpf = self.combined_bpf

            LOG.info("BPF program loaded with both gtp_decap_handler and gtp_encap_handler")

            # We're not reusing maps in this path
            self._bpf_maps_dir = GTP_MAPS_DIR
            
            LOG.info("Encapsulation program loaded successfully")
            
            # For step 3, enable BPF maps for session lookup
            try:
                # CRITICAL FIX: Always detect currently attached TC programs and use THEIR maps
                # This ensures pipelined writes to the same maps that the attached programs read from
                attached_maps = self._detect_attached_tc_maps()

                if attached_maps:
                    LOG.info(f"Detected attached TC programs with maps: {attached_maps}")
                    # Use wrapper for external maps from currently attached TC programs
                    self.ue_session_map = BPFMapWrapper(attached_maps.get('ue_session'))
                    self.config_map = BPFMapWrapper(attached_maps.get('config'))
                    self.stats_map = BPFMapWrapper(attached_maps.get('stats'))
                    self._maps_initialized = True
                    LOG.info("✅ Using maps from ATTACHED TC programs - writes will be visible to eBPF!")
                elif hasattr(self, '_tc_map_ids') and self._tc_map_ids:
                    LOG.info(f"Falling back to cached TC maps: {self._tc_map_ids}")
                    # Use wrapper for external maps
                    self.ue_session_map = BPFMapWrapper(self._tc_map_ids.get('ue_session'))
                    self.config_map = BPFMapWrapper(self._tc_map_ids.get('config'))
                    self.stats_map = BPFMapWrapper(self._tc_map_ids.get('stats'))
                    self._maps_initialized = True
                    LOG.info("Using cached TC maps via bpftool wrapper")
                else:
                    # Normal BCC map initialization (only if no TC programs attached yet)
                    LOG.info("No attached TC programs detected, using newly loaded program's maps")
                    self.ue_session_map = self.decap_bpf.get_table(UE_SESSION_MAP)
                    self.config_map = self.decap_bpf.get_table(CONFIG_MAP)
                    self.stats_map = self.decap_bpf.get_table(STATS_MAP)

                    # Validate maps are actually accessible
                    if self.ue_session_map is None or self.config_map is None or self.stats_map is None:
                        raise Exception("One or more BPF maps are None after initialization")

                    # Test map access to ensure they're working
                    test_key = self.config_map.Key(CONFIG_DEBUG_LEVEL)
                    test_val = self.config_map.Leaf(1)
                    self.config_map[test_key] = test_val

                    # Maps are shared via file descriptor passing
                    LOG.info("Maps are shared between encap and decap programs (FD sharing)")

                    LOG.info("BPF maps initialized and validated successfully")
                    self._maps_initialized = True
                
                # Pin maps if BPF filesystem is available
                if self._ensure_bpf_fs_mounted():
                    self._pin_maps()
            except Exception as e:
                LOG.error("CRITICAL: BPF maps initialization failed: %s", e)
                LOG.error("This will prevent all GTP processing - aborting initialization")
                self.ue_session_map = None
                self.config_map = None  
                self.stats_map = None
                self._maps_initialized = False
                return False
            
            # Use the same combined program for gtp_veth0 mark restoration
            # The gtp_veth0_mark_handler function is already in the main program
            LOG.info("Using combined eBPF program for gtp_veth0 mark restoration")
            self.veth0_mark_bpf = self.combined_bpf
            LOG.info("gtp_veth0 mark restoration program ready (shared with main program)")

            LOG.info("eBPF GTP programs loaded successfully")
            return True
            
        except Exception as e:
            LOG.error("Failed to load eBPF programs: %s", e)
            import traceback
            LOG.error("Full traceback: %s", traceback.format_exc())
            return False
    
    def _pin_maps(self):
        """Pin BPF maps to filesystem for persistence"""
        try:
            # Pin each map
            maps_to_pin = [
                ("ue_sessions", self.ue_session_map),
                ("config", self.config_map),
                ("stats", self.stats_map)
            ]
            
            for map_name, map_obj in maps_to_pin:
                if map_obj:
                    pin_path = os.path.join(GTP_MAPS_DIR, map_name)
                    # Note: BCC doesn't have direct pinning API, would need to use libbpf
                    # For now, we'll track that maps should be persisted
                    # In production, you'd use libbpf or bpftool for proper pinning
                    LOG.info(f"Map {map_name} marked for persistence at {pin_path}")
            
            self._maps_pinned = True
            LOG.info("BPF maps marked for persistence")
            
            # Save current sessions to file for recovery
            self._save_sessions_to_file()
            
        except Exception as e:
            LOG.error(f"Failed to pin maps: {e}")
    
    def _save_sessions_to_file(self):
        """Save current sessions to file for recovery after restart"""
        try:
            session_file = os.path.join(EBPF_GTP_DIR, "saved_sessions.json")
            sessions = []
            
            for key in self.ue_session_map.keys():
                value = self.ue_session_map[key]
                if value:
                    session_info = self._parse_session_value_ctypes(value)
                    sessions.append({
                        "ue_ip_str": session_info['ue_ip_str'],
                        "session_info": session_info
                    })
            
            import json
            with open(session_file, 'w') as f:
                json.dump(sessions, f, indent=2)
            
            LOG.info(f"Saved {len(sessions)} sessions to {session_file}")
            
        except Exception as e:
            LOG.warning(f"Failed to save sessions to file: {e}")
    
    def _restore_sessions_from_file(self):
        """Restore sessions from file after restart"""
        try:
            session_file = os.path.join(EBPF_GTP_DIR, "saved_sessions.json")
            if not os.path.exists(session_file):
                return
            
            import json
            with open(session_file, 'r') as f:
                sessions = json.load(f)
            
            restored = 0
            for session in sessions:
                try:
                    info = session['session_info']
                    self.add_ue_session(
                        ue_ip_str=info['ue_ip_str'],
                        enb_ip_int=info['enb_ip_int'],
                        ue_teid=info['ue_teid'],
                        enb_teid=info['enb_teid'],
                        bearer_id=info.get('bearer_id', 0)
                    )
                    restored += 1
                except Exception as e:
                    LOG.warning(f"Failed to restore session: {e}")
            
            LOG.info(f"Restored {restored} sessions from file")
            
            # Remove file after successful restoration
            os.remove(session_file)
            
        except Exception as e:
            LOG.warning(f"Failed to restore sessions from file: {e}")

    def _attach_programs(self) -> bool:
        """
        Attach eBPF programs using TC+TC approach with flower filter
        
        Returns:
            True if successful, False otherwise
        """
        try:
            LOG.info("Attaching eBPF programs using TC+TC approach...")
            
            # Get interface indices
            s1u_ifindex = self._get_ifindex(self.s1u_iface)
            if not s1u_ifindex:
                LOG.error("Failed to get interface index for %s", self.s1u_iface)
                return False
                
            ebpf_ifindex = self._get_ifindex(GTP_VETH_EBPF)
            if not ebpf_ifindex:
                LOG.error("Failed to get interface index for %s", GTP_VETH_EBPF)
                return False
            
            # Only clean up filters if we're NOT reusing existing TC programs
            if not (hasattr(self, '_tc_map_ids') and self._tc_map_ids):
                # Clean up ALL existing filters on interfaces to ensure fresh start
                try:
                    # Clean all filters on s1u interface
                    subprocess.run(['tc', 'filter', 'del', 'dev', self.s1u_iface, 'ingress'],
                                  capture_output=True, check=False, timeout=2)
                    # Clean all filters on veth interfaces
                    subprocess.run(['tc', 'filter', 'del', 'dev', GTP_VETH_EBPF, 'ingress'],
                                  capture_output=True, check=False, timeout=2)
                    subprocess.run(['tc', 'filter', 'del', 'dev', GTP_VETH_EBPF, 'egress'],
                                  capture_output=True, check=False, timeout=2)
                    subprocess.run(['tc', 'filter', 'del', 'dev', GTP_VETH_OVS, 'ingress'],
                                  capture_output=True, check=False, timeout=2)
                    LOG.debug("Cleaned up any existing filters on all interfaces")
                except:
                    pass
            else:
                LOG.debug("Keeping existing TC filters since we're reusing programs")
            
            # Step 1: Setup flower filter on s1u_iface (eth0) to redirect GTP traffic
            LOG.info("Setting up flower filter on %s to redirect GTP traffic to %s", 
                     self.s1u_iface, GTP_VETH_EBPF)
            self._setup_flower_redirect(self.s1u_iface, s1u_ifindex, ebpf_ifindex)
            
            # Step 2: Setup clsact qdisc on gtp_veth1
            try:
                self.ipr.tc("add", "clsact", ebpf_ifindex)
                LOG.debug("Added clsact qdisc to %s", GTP_VETH_EBPF)
            except NetlinkError as e:
                if "File exists" in str(e):
                    LOG.debug("clsact qdisc already exists on %s", GTP_VETH_EBPF)
                else:
                    raise
            

            LOG.info("switched the interaces: ha 1")
            # Step 3: Attach TC ingress for decap (uplink GTP traffic)
            # Skip if we're reusing existing TC programs
            if hasattr(self, '_tc_map_ids') and self._tc_map_ids:
                LOG.info("Skipping TC attachment - reusing existing programs with maps %s", self._tc_map_ids)
                return True

            LOG.info("Attaching GTP decapsulation to %s ingress", GTP_VETH_EBPF)
            self._attach_tc_program(
                interface=GTP_VETH_EBPF,
                ifindex=ebpf_ifindex,
                bpf_prog=self.decap_bpf,
                func_name="gtp_decap_handler",
                direction="ingress"   # testing purpose ingress-> egress
            )

            # Step 4: Attach gtp_veth0 mark restoration program
            ovs_ifindex = self._get_ifindex(GTP_VETH_OVS)
            if not ovs_ifindex:
                LOG.error("Failed to get interface index for %s", GTP_VETH_OVS)
                return False

            LOG.info("Attaching mark restoration to %s ingress", GTP_VETH_OVS)
            self._attach_tc_program(
                interface=GTP_VETH_OVS,
                ifindex=ovs_ifindex,
                bpf_prog=self.veth0_mark_bpf,
                func_name="gtp_veth0_mark_handler",
                direction="ingress"
            )

            # Step 5: Attach GTP encapsulation to gtp_veth0 egress for downlink
            LOG.info("Attaching GTP encapsulation to %s egress", GTP_VETH_OVS)
            self._attach_tc_program(
                interface=GTP_VETH_OVS,
                ifindex=ovs_ifindex,
                bpf_prog=self.encap_bpf,
                func_name="gtp_encap_handler",
                direction="egress"
            )

            LOG.info("Successfully attached TC+TC programs with mark restoration and downlink encap")
            LOG.info("Uplink: %s --flower--> %s(ingress/decap) --> %s(ingress/mark) --> OVS",
                     self.s1u_iface, GTP_VETH_EBPF, GTP_VETH_OVS)
            LOG.info("Downlink: OVS --> %s(egress/encap) --> Physical NIC --> RAN",
                     GTP_VETH_OVS)

            # NOTE: Map sharing is automatic because all programs (decap, encap, veth0_mark)
            # are loaded from the same combined_bpf object via load_func().
            # They share the same map instances by design in BCC.
            LOG.info("BPF maps are shared across all programs via combined_bpf object")

            return True
            
        except Exception as e:
            LOG.error("Failed to attach eBPF programs: %s", e)
            import traceback
            LOG.error("Traceback: %s", traceback.format_exc())
            return False

    def _attach_tc_program(self, interface: str, ifindex: int, bpf_prog: BPF,
                          func_name: str, direction: str):
        """
        Attach eBPF program to interface via TC

        Args:
            interface: Interface name
            ifindex: Interface index
            bpf_prog: BPF program object
            func_name: Function name to attach
            direction: "ingress" or "egress"
        """
        try:
            # Add clsact qdisc
            try:
                self.ipr.tc("add", "clsact", ifindex)
                LOG.debug("Added clsact qdisc to %s", interface)
            except NetlinkError as e:
                # Already exists, which is fine
                LOG.debug("clsact qdisc already exists on %s: %s", interface, e)
                pass

            # Attach BPF program
            parent = "ffff:fff2" if direction == "ingress" else "ffff:fff3"

            LOG.debug("Loading function %s from BPF program", func_name)
            fn = bpf_prog.load_func(func_name, BPF.SCHED_CLS)

            # CRITICAL FIX: Mark handler on gtp_veth0 must NOT use direct_action
            # With direct_action=True, TC_ACT_OK stops packet processing
            # We need packets to continue to OVS after mark restoration
            use_direct_action = True
            if interface == GTP_VETH_OVS and func_name == "gtp_veth0_mark_handler":
                LOG.info("Mark handler on %s: Disabling direct_action to allow packets to reach OVS", interface)
                use_direct_action = False

            LOG.debug("Attaching BPF filter to %s %s (fd=%d, direct_action=%s)",
                     interface, direction, fn.fd, use_direct_action)
            self.ipr.tc("add-filter", "bpf", ifindex, ":1",
                       fd=fn.fd, name=fn.name, parent=parent,
                       classid=1, direct_action=use_direct_action)
            
            LOG.info("Successfully attached %s to %s %s", func_name, interface, direction)
            
        except Exception as e:
            LOG.error("Failed to attach %s to %s %s: %s", func_name, interface, direction, e)
            raise

    def _setup_flower_redirect(self, interface: str, ifindex: int, target_ifindex: int):
        """
        Setup flower filter to redirect GTP traffic to gtp_veth1
        Using tc command directly for better compatibility
        Now with separate filters for data and control plane (CRITICAL for kernel module replacement)
        """
        try:
            # Ensure clsact qdisc exists on source interface
            try:
                subprocess.run([
                    'tc', 'qdisc', 'add', 'dev', interface, 'clsact'
                ], capture_output=True, check=False)
            except:
                pass  # OK if already exists
            
            # Remove any existing flower filters
            for prio in ['1', '2']:
                try:
                    subprocess.run([
                        'tc', 'filter', 'del', 'dev', interface, 
                        'ingress', 'pref', prio
                    ], capture_output=True, check=False)
                except:
                    pass  # OK if doesn't exist
            
            # Filter 1: GTP-U data traffic (T-PDU, type 255)
            cmd_data = [
                'tc', 'filter', 'add', 'dev', interface,
                'ingress', 'prio', '1', 'protocol', 'ip',
                'flower',
                'ip_proto', 'udp',
                'dst_port', '2152',
                'action', 'mirred', 'ingress', 'redirect', 'dev', GTP_VETH_EBPF
            ]
            
            result = subprocess.run(cmd_data, capture_output=True, text=True)
            if result.returncode != 0:
                LOG.error("Failed to add data plane filter: %s", result.stderr)
                raise Exception(f"Data plane tc command failed: {result.stderr}")
            
            # Filter 2: GTP-C control traffic (echo, handover, etc.)
            # This ensures control plane packets also reach our eBPF handler
            cmd_control = [
                'tc', 'filter', 'add', 'dev', interface,
                'ingress', 'prio', '2', 'protocol', 'ip', 
                'flower',
                'ip_proto', 'udp',
                'dst_port', '2123',  # GTP-C port
                'action', 'mirred', 'ingress', 'redirect', 'dev', GTP_VETH_EBPF
            ]
            
            # Note: GTP-C on 2123 is optional, most control is on 2152
            subprocess.run(cmd_control, capture_output=True, text=True)
            
            LOG.info("Flower filters successfully added on %s to redirect GTP to %s", 
                     interface, GTP_VETH_EBPF)
            LOG.info("Data plane (2152) and Control plane (2123) filters active")
            
        except Exception as e:
            LOG.error("Failed to setup flower filter: %s", e)
            raise

    def _attach_xdp_program(self, interface: str, ifindex: int, bpf_prog: BPF, func_name: str):
        """
        Attach eBPF program to interface via XDP
        
        Args:
            interface: Interface name
            ifindex: Interface index
            bpf_prog: BPF program object
            func_name: Function name to attach
        """
        # XDP attachment code commented out - using TC eBPF
        # try:
        #     LOG.debug("Loading XDP function %s from BPF program", func_name)
        #     fn = bpf_prog.load_func(func_name, BPF.XDP)
        #     
        #     LOG.debug("Attaching XDP program to %s (fd=%d)", interface, fn.fd)
        #     # Attach XDP program to interface using Native mode for better packet interception
        #     try:
        #         # Try XDP Native mode first (more reliable for packet capture)
        #         bpf_prog.attach_xdp(interface, fn, BPF.XDP_FLAGS_DRV_MODE)
        #         LOG.info("Successfully attached %s to %s using XDP Native mode", func_name, interface)
        #     except Exception as native_err:
        #         LOG.warning("XDP Native mode failed (%s), falling back to Generic mode", native_err)
        #         # Fallback to Generic mode if Native mode is not supported
        #         bpf_prog.attach_xdp(interface, fn, BPF.XDP_FLAGS_SKB_MODE)
        #     
        #     LOG.info("Successfully attached %s to %s (XDP)", func_name, interface)
        #     
        # except Exception as e:
        #     LOG.error("Failed to attach %s to %s (XDP): %s", func_name, interface, e)
        #     raise
        pass

    def _attach_af_xdp_programs(self) -> bool:
        """Attach XDP programs configured for AF_XDP redirect"""
        # AF_XDP code commented out - using TC eBPF
        # try:
        #     LOG.info("Attaching XDP programs for AF_XDP redirect...")
        #     
        #     # Create AF_XDP redirect program dynamically
        #     af_xdp_redirect_prog = f"""
        # #include <linux/bpf.h>
        # #include <linux/if_ether.h>
        # #include <linux/ip.h>
        # #include <linux/udp.h>
        # #include <linux/in.h>
        # 
        # // AF_XDP socket map
        # BPF_XSKMAP(xsks_map, 64);
        # 
        # // Statistics map
        # BPF_ARRAY(stats_map, __u64, 16);
        # 
        # // Statistics counters
        # #define STATS_AF_XDP_REDIRECT 10
        # #define STATS_AF_XDP_PACKETS 11
        # 
        # static inline void update_stats(__u32 counter_id, __u64 value) {{
        #     __u64* count = stats_map.lookup(&counter_id);
        #     if (count) {{
        #         __sync_fetch_and_add(count, value);
        #     }}
        # }}
        # 
        # int af_xdp_redirect_handler(struct xdp_md* ctx) {{
        #     void* data = (void*)(long)ctx->data;
        #     void* data_end = (void*)(long)ctx->data_end;
        #     
        #     // Parse Ethernet header
        #     struct ethhdr* eth = data;
        #     if ((void*)(eth + 1) > data_end) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Check if IPv4
        #     if (eth->h_proto != __constant_htons(ETH_P_IP)) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Parse IP header
        #     struct iphdr* ip = (struct iphdr*)(eth + 1);
        #     if ((void*)(ip + 1) > data_end) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Check if UDP
        #     if (ip->protocol != IPPROTO_UDP) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Parse UDP header
        #     struct udphdr* udp = (struct udphdr*)((void*)ip + (ip->ihl * 4));
        #     if ((void*)(udp + 1) > data_end) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Check if GTP-U port (2152)
        #     if (udp->dest != __constant_htons(2152)) {{
        #         return XDP_PASS;
        #     }}
        #     
        #     // Update statistics
        #     update_stats(STATS_AF_XDP_PACKETS, 1);
        #     
        #     // Redirect to AF_XDP socket
        #     int queue_id = 0;  // Use queue 0 for simplicity
        #     int ret = bpf_redirect_map(&xsks_map, queue_id, 0);
        #     if (ret == XDP_REDIRECT) {{
        #         update_stats(STATS_AF_XDP_REDIRECT, 1);
        #     }}
        #     
        #     return ret;
        # }}
        # """
        #     
        #     # Load AF_XDP redirect program
        #     cflags = ['-I', EBPF_GTP_DIR]
        #     self.af_xdp_bpf = BPF(text=af_xdp_redirect_prog, cflags=cflags)
        #     
        #     # Get the XSK map and register our AF_XDP socket
        #     xsks_map = self.af_xdp_bpf.get_table("xsks_map")
        #     if self.af_xdp_socket and self.af_xdp_socket.socket:
        #         xsks_map[0] = self.af_xdp_socket.socket.fileno()
        #     
        #     # Attach AF_XDP redirect program to S1-U interface
        #     LOG.info("Attaching AF_XDP redirect program to %s", self.s1u_iface)
        #     fn = self.af_xdp_bpf.load_func("af_xdp_redirect_handler", BPF.XDP)
        #     self.af_xdp_bpf.attach_xdp(self.s1u_iface, fn, BPF.XDP_FLAGS_SKB_MODE)
        #     
        #     LOG.info("AF_XDP redirect programs attached successfully")
        #     return True
        #     
        # except Exception as e:
        #     LOG.error("Failed to attach AF_XDP redirect programs: %s", e)
        #     import traceback
        #     LOG.error("Full traceback: %s", traceback.format_exc())
        #     return False
        return False

    def get_statistics(self) -> Dict:
        """
        Get eBPF GTP statistics - Step 6: Full statistics (per-UE metrics)
        
        Returns:
            Dictionary containing comprehensive statistics
        """
        try:
            stats = {}
            
            # Global statistics - all counters
            stats_names = {
                STATS_UL_PACKETS: "ul_packets",
                STATS_DL_PACKETS: "dl_packets", 
                STATS_UL_BYTES: "ul_bytes",
                STATS_DL_BYTES: "dl_bytes",
                STATS_UL_ERRORS: "ul_errors",
                STATS_DL_ERRORS: "dl_errors",
                STATS_GTP_INVALID: "gtp_invalid",
                STATS_SESSION_MISS: "session_miss",
                STATS_GTP_DECAP_SUCCESS: "gtp_decap_success",
                STATS_GTP_ENCAP_SUCCESS: "gtp_encap_success",
                STATS_TEID_MISMATCH: "teid_mismatch",
                STATS_UE_ATTACH: "ue_attach",
                STATS_UE_DETACH: "ue_detach",
                STATS_PKT_DROPPED: "packets_dropped",
                STATS_PKT_FORWARDED: "packets_forwarded"
            }
            
            global_stats = {}
            for key, name in stats_names.items():
                try:
                    value = self.stats_map[key].value
                    global_stats[name] = value
                except (KeyError, AttributeError):
                    global_stats[name] = 0
            
            stats['global'] = global_stats
            
            # Per-UE statistics with enhanced metrics
            ue_stats = {}
            for session_key, session_info in self.ue_session_map.items():
                ue_ip = socket.inet_ntoa(struct.pack("!I", session_key.ue_ip))
                ue_stats[ue_ip] = {
                    # Traffic counters
                    'ul_packets': session_info.ul_packets,
                    'dl_packets': session_info.dl_packets,
                    'ul_bytes': session_info.ul_bytes,
                    'dl_bytes': session_info.dl_bytes,
                    
                    # Session metadata
                    'last_seen': session_info.last_seen,
                    'session_flags': session_info.session_flags,
                    
                    # Tunnel information
                    'teid_in': session_info.teid_in,
                    'teid_out': session_info.teid_out,
                    'enb_ip': socket.inet_ntoa(struct.pack("!I", session_info.enb_ip)),
                    
                    # Derived metrics
                    'total_packets': session_info.ul_packets + session_info.dl_packets,
                    'total_bytes': session_info.ul_bytes + session_info.dl_bytes,
                    'is_active': bool(session_info.session_flags & 1)
                }
            
            stats['ue_sessions'] = ue_stats
            
            # Additional derived statistics  
            stats['summary'] = {
                'total_ue_sessions': len(ue_stats),
                'active_ue_sessions': len([s for s in ue_stats.values() if s['is_active']]),
                'total_ul_packets': global_stats['ul_packets'],
                'total_dl_packets': global_stats['dl_packets'],
                'total_packets': global_stats['ul_packets'] + global_stats['dl_packets'],
                'total_ul_bytes': global_stats['ul_bytes'],
                'total_dl_bytes': global_stats['dl_bytes'],
                'total_bytes': global_stats['ul_bytes'] + global_stats['dl_bytes'],
                'error_rate': ((global_stats['ul_errors'] + global_stats['dl_errors']) / 
                              max(global_stats['ul_packets'] + global_stats['dl_packets'], 1)) * 100,
                'session_miss_rate': (global_stats['session_miss'] / 
                                     max(global_stats['ul_packets'] + global_stats['dl_packets'], 1)) * 100
            }
            
            return stats
            
        except Exception as e:
            LOG.error("Failed to get statistics: %s", e)
            import traceback
            LOG.error("Full traceback: %s", traceback.format_exc())
            return {}

    def _get_interface_ip(self, interface: str) -> Optional[str]:
        """
        Get IP address of an interface
        
        Args:
            interface: Interface name
            
        Returns:
            IP address as string or None
        """
        try:
            pass
        except:
            pass
        return None
    
    def _get_ifindex(self, interface: str) -> Optional[int]:
        """Get interface index by name"""
        try:
            for link in self.ipr.get_links():
                if link.get_attr('IFLA_IFNAME') == interface:
                    return link['index']
            return None
        except Exception as e:
            LOG.error("Failed to get interface index for %s: %s", interface, e)
            return None

    def _get_mac_bytes(self, interface: str) -> bytes:
        """Get MAC address as bytes"""
        try:
            mac_str = get_mac_address_from_iface(interface)
            return bytes.fromhex(mac_str.replace(':', ''))
        except Exception:
            return b'\x02\x00\x00\x00\x00\x01'  # Default MAC

    def register_interface_callback(self, callback):
        """
        Register a callback to be called when gtp_veth0 interface is created/recreated
        CRITICAL for notifying classifier to refresh flows when port numbers change
        """
        self._interface_callbacks.append(callback)
        LOG.debug("Registered interface callback: %s", callback)
    
    def _notify_interface_created(self):
        """
        Notify all registered callbacks that gtp_veth0 interface was created/recreated
        CRITICAL for ensuring OVS flows are updated when port numbers change
        """
        try:
            for callback in self._interface_callbacks:
                try:
                    callback()
                    LOG.debug("Interface callback executed successfully")
                except Exception as e:
                    LOG.error("Interface callback failed: %s", e)
        except Exception as e:
            LOG.error("Failed to notify interface callbacks: %s", e)

    def _setup_ue_subnet_proxy_arp(self):
        """
        Setup proxy ARP for entire UE subnet (192.168.128.0/24) to handle
        dynamic UE IPs without requiring individual ARP entries.

        Phase 1: Basic proxy ARP setup for datapath functionality
        """
        try:
            import subprocess

            LOG.info("Setting up proxy ARP for UE subnet 192.168.128.0/24")

            # Enable proxy ARP on gtp_br0
            try:
                subprocess.run([
                    'sysctl', '-w', 'net.ipv4.conf.gtp_br0.proxy_arp=1'
                ], check=True, capture_output=True)
                LOG.info("Enabled proxy ARP on gtp_br0")
            except subprocess.CalledProcessError as e:
                LOG.warning("Failed to enable proxy ARP on gtp_br0: %s", e)

            # Enable proxy ARP PVLAN for better subnet handling
            try:
                subprocess.run([
                    'sysctl', '-w', 'net.ipv4.conf.gtp_br0.proxy_arp_pvlan=1'
                ], check=True, capture_output=True)
                LOG.info("Enabled proxy ARP PVLAN on gtp_br0")
            except subprocess.CalledProcessError as e:
                LOG.warning("Failed to enable proxy ARP PVLAN on gtp_br0: %s", e)

            # Add the UE subnet route to gtp_veth0
            try:
                subprocess.run([
                    'ip', 'route', 'add', '192.168.128.0/24', 'dev', GTP_VETH_OVS, 'scope', 'link'
                ], capture_output=True)
                LOG.info("Added UE subnet route via %s", GTP_VETH_OVS)
            except subprocess.CalledProcessError as e:
                LOG.debug("Route setup had issues (may already exist): %s", e)

            return True

        except Exception as e:
            LOG.error("Failed to setup proxy ARP for UE subnet: %s", e)
            return False

    def debug_session_byte_order(self) -> None:
        """
        Debug method to check byte order issues in stored sessions
        """
        try:
            if not self.enabled or not self.ue_session_map:
                LOG.warning("eBPF not enabled or map not available")
                return
                
            LOG.info("=== Debugging eBPF Session Byte Order ===")
            
            # Iterate through all sessions
            for key, value in self.ue_session_map.items():
                # Key is 4 bytes (IP address)
                if len(key) >= 4:
                    # Extract IP as different byte orders
                    ip_be = struct.unpack("!I", key[:4])[0]  # Big-endian (network order)
                    ip_le = struct.unpack("<I", key[:4])[0]  # Little-endian
                    
                    # Convert to dotted notation
                    ip_be_str = socket.inet_ntoa(struct.pack("!I", ip_be))
                    ip_le_str = socket.inet_ntoa(struct.pack("!I", ip_le))
                    
                    LOG.info("Session key: 0x%08x (BE: %s, LE: %s)", 
                             ip_be, ip_be_str, ip_le_str)
                    
                    # Check which one looks correct (should be in 192.168.128.0/24)
                    if ip_be_str.startswith("192.168.128."):
                        LOG.info("  -> Correct byte order (big-endian)")
                    elif ip_le_str.startswith("192.168.128."):
                        LOG.warning("  -> WRONG byte order! Should be: %s", ip_le_str)
                    else:
                        LOG.warning("  -> IP not in expected range")
                        
        except Exception as e:
            LOG.error("Failed to debug byte order: %s", e)

    def cleanup(self):
        """
        Cleanup all eBPF programs and filters
        """
        LOG.info("Cleaning up eBPF GTP manager...")
        
        try:
            # Remove flower filter from s1u interface
            s1u_ifindex = self._get_ifindex(self.s1u_iface)
            if s1u_ifindex:
                try:
                    subprocess.run([
                        'tc', 'filter', 'del', 'dev', self.s1u_iface,
                        'ingress', 'pref', '1'
                    ], capture_output=True, check=False)
                    LOG.info("Removed flower filter from %s", self.s1u_iface)
                except:
                    pass
        except:
            pass
        
        try:
            # Remove TC filters from gtp_veth1
            ebpf_ifindex = self._get_ifindex(GTP_VETH_EBPF)
            if ebpf_ifindex:
                try:
                    subprocess.run([
                        'tc', 'qdisc', 'del', 'dev', GTP_VETH_EBPF, 'clsact'
                    ], capture_output=True, check=False)
                    LOG.info("Removed clsact qdisc from %s", GTP_VETH_EBPF)
                except:
                    pass
        except:
            pass
        
        # Close BPF programs
        if hasattr(self, 'combined_bpf'):
            del self.combined_bpf
        if hasattr(self, 'decap_bpf'):
            del self.decap_bpf
        if hasattr(self, 'encap_bpf'):
            del self.encap_bpf
        if hasattr(self, 'veth0_mark_bpf'):
            del self.veth0_mark_bpf

        # No map pinning cleanup needed with Strategy 1
        # Maps are automatically cleaned up when BPF program is deleted
        
        # Remove veth pair
        self._cleanup_veth_pair()
        
        # Close control plane echo handler
        if hasattr(self, '_echo_thread') and self._echo_thread:
            self._stop_echo_handler = True
            self._echo_thread.join(timeout=2.0)
        
        # Close IPRoute
        if hasattr(self, 'ipr'):
            self.ipr.close()
        
        self.enabled = False
        LOG.info("eBPF GTP manager cleanup completed")
    
    def _cleanup_veth_pair(self):
        """Remove veth pair if it exists"""
        try:
            self.ipr.link("del", ifname=GTP_VETH_OVS)
            LOG.info("Removed veth pair %s", GTP_VETH_OVS)
        except:
            pass  # Already removed
    
    def is_ebpf_gtp_enabled(self) -> bool:
        """
        Check if eBPF GTP is enabled
        
        Returns:
            bool: True if eBPF GTP is enabled and initialized
        """
        return self.enabled

    def _parse_session_value_ctypes(self, bpf_value) -> Optional[Dict[str, Any]]:
        """Parse BPF session value using ctypes."""
        try:
            # If bpf_value is already a UeSessionInfo ctypes structure, use it directly
            if isinstance(bpf_value, UeSessionInfo):
                value = bpf_value
            else:
                # This shouldn't happen with proper ctypes usage
                LOG.warning("Unexpected value type in _parse_session_value_ctypes: %s", type(bpf_value))
                return None
            
            # Convert network byte order IPs to host strings
            if sys.byteorder == 'little':
                enb_ip = socket.inet_ntoa(socket.htonl(value.enb_ip).to_bytes(4, 'big'))
                tun_ipv4_dst = socket.inet_ntoa(socket.htonl(value.tun_ipv4_dst).to_bytes(4, 'big'))
            else:
                enb_ip = socket.inet_ntoa(value.enb_ip.to_bytes(4, 'big'))
                tun_ipv4_dst = socket.inet_ntoa(value.tun_ipv4_dst.to_bytes(4, 'big'))
            
            # Extract IMSI string
            imsi_str = ''
            if value.imsi_len > 0:
                imsi_bytes = bytes(value.imsi[:value.imsi_len])
                try:
                    imsi_str = imsi_bytes.decode('utf-8')
                except:
                    imsi_str = imsi_bytes.hex()
            
            # Convert MAC addresses to hex strings
            ul_mac_src = ':'.join(f'{b:02x}' for b in value.ul_mac_src)
            ul_mac_dst = ':'.join(f'{b:02x}' for b in value.ul_mac_dst)
            
            # Determine session state
            state = 'active' if value.session_flags & 0x1 else 'inactive'
            
            return {
                'enb_ip': enb_ip,
                'teid_ul_in': value.teid_ul_in,
                'teid_ul_out': value.teid_ul_out,
                'teid_dl_in': value.teid_dl_in,
                'teid_dl_out': value.teid_dl_out,
                's1u_ifindex': value.s1u_ifindex,
                'ul_mac_src': ul_mac_src,
                'ul_mac_dst': ul_mac_dst,
                'qos_mark': value.qos_mark,
                'bearer_id': value.bearer_id,
                'ul_bytes': value.ul_bytes,
                'dl_bytes': value.dl_bytes,
                'ul_packets': value.ul_packets,
                'dl_packets': value.dl_packets,
                'last_seen_ns': value.last_seen,
                'session_flags': value.session_flags,
                'state': state,
                'imsi': imsi_str,
                'imsi_len': value.imsi_len,
                'encoded_imsi': value.encoded_imsi,
                'qfi': value.qfi,
                'tunnel_id': value.tunnel_id,
                'tun_ipv4_dst': tun_ipv4_dst,
                'tun_flags': value.tun_flags,
                'direction': value.direction,
                'original_port': value.original_port,
                'metadata_mark': value.metadata_mark
            }
            
        except Exception as e:
            LOG.error("Failed to parse session value with ctypes: %s", e)
            return None

    def _parse_session_value(self, bpf_value) -> Optional[Dict[str, Any]]:
        """Parse raw BPF session value into a dict."""
        try:
            if hasattr(bpf_value, 'value'):
                raw_data = bpf_value.value
            elif isinstance(bpf_value, bytes):
                raw_data = bpf_value
            else:
                raw_data = bytes(bpf_value)

            try:
                if len(raw_data) >= SESSION_STRUCT_SIZE:
                    unpacked = struct.unpack(SESSION_STRUCT_FMT, raw_data[:SESSION_STRUCT_SIZE])
                    metadata_mark = unpacked[-1]
                    reserved_bytes = unpacked[-2]
                elif len(raw_data) >= SESSION_STRUCT_SIZE_LEGACY:
                    unpacked = struct.unpack(SESSION_STRUCT_FMT_LEGACY, raw_data[:SESSION_STRUCT_SIZE_LEGACY])
                    metadata_mark = None
                    reserved_bytes = b'\x00\x00\x00'
                else:
                    LOG.warning("Session data too short: %d bytes", len(raw_data))
                    return None
            except struct.error as e:
                LOG.warning("Failed to parse session struct (size=%d): %s", len(raw_data), e)
                return None

            idx = 0
            enb_ip_int = unpacked[idx]; idx += 1
            teid_ul_in = unpacked[idx]; idx += 1
            teid_ul_out = unpacked[idx]; idx += 1
            teid_dl_in = unpacked[idx]; idx += 1
            teid_dl_out = unpacked[idx]; idx += 1
            s1u_ifindex = unpacked[idx]; idx += 1
            sgi_ifindex = unpacked[idx]; idx += 1
            ovs_ifindex = unpacked[idx]; idx += 1
            ul_mac_src = unpacked[idx]; idx += 1
            ul_mac_dst = unpacked[idx]; idx += 1
            qos_mark = unpacked[idx]; idx += 1
            bearer_id = unpacked[idx]; idx += 1
            ul_bytes = unpacked[idx]; idx += 1
            dl_bytes = unpacked[idx]; idx += 1
            ul_packets = unpacked[idx]; idx += 1
            dl_packets = unpacked[idx]; idx += 1
            last_seen = unpacked[idx]; idx += 1
            session_flags = unpacked[idx]; idx += 1
            imsi_bytes = unpacked[idx]; idx += 1
            imsi_len = unpacked[idx]; idx += 1

            encoded_imsi = qfi = tunnel_id = tun_ipv4_dst_int = tun_flags = direction = original_port = 0

            if len(unpacked) > idx:
                encoded_imsi = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                qfi = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                tunnel_id = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                tun_ipv4_dst_int = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                tun_flags = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                direction = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                original_port = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                reserved_bytes = unpacked[idx]; idx += 1
            if len(unpacked) > idx:
                metadata_mark = unpacked[idx]

            if metadata_mark is None:
                metadata_mark = _compute_ue_mark_from_ip_int(teid_ul_in)

            state = 'active' if (session_flags & 0x1) else 'inactive'

            return {
                'enb_ip': socket.inet_ntop(socket.AF_INET, struct.pack('!I', enb_ip_int)),
                'teid_ul_in': teid_ul_in,
                'teid_ul_out': teid_ul_out,
                'teid_dl_in': teid_dl_in,
                'teid_dl_out': teid_dl_out,
                's1u_ifindex': s1u_ifindex,
                'sgi_ifindex': sgi_ifindex,
                'ovs_ifindex': ovs_ifindex,
                'ul_mac_src': ul_mac_src,
                'ul_mac_dst': ul_mac_dst,
                'qos_mark': qos_mark,
                'bearer_id': bearer_id,
                'ul_bytes': ul_bytes,
                'dl_bytes': dl_bytes,
                'ul_packets': ul_packets,
                'dl_packets': dl_packets,
                'last_seen': last_seen,
                'session_flags': session_flags,
                'imsi': imsi_bytes[:imsi_len],
                'imsi_len': imsi_len,
                'encoded_imsi': encoded_imsi,
                'qfi': qfi,
                'tunnel_id': tunnel_id,
                'tun_ipv4_dst': socket.inet_ntop(socket.AF_INET, struct.pack('!I', tun_ipv4_dst_int)) if tun_ipv4_dst_int else None,
                'tun_flags': tun_flags,
                'direction': direction,
                'original_port': original_port,
                'metadata_mark': metadata_mark,
                'state': state,
                'created_time': int(time.time()),
            }

        except Exception as e:
            LOG.error("Failed to parse session value: %s", e)
            return None

    def get_gtp_veth_ovs_port(self) -> Optional[int]:
        """
        Get the dynamically detected OVS port number for gtp_veth0
        
        Returns:
            OVS port number or None if not yet detected
        """
        if self._ovs_port_number:
            return self._ovs_port_number
        
        # Try to detect it now if not already known
        try:
            ofport = BridgeTools.get_ofport(GTP_VETH_OVS)
            return ofport
        except Exception as e:
            LOG.warning(f"Failed to get OVS port number: {e}")
        
        return None


# Global singleton instance to prevent garbage collection
_ebpf_gtp_manager_instance = None
_ebpf_gtp_manager_lock = threading.Lock()

def get_ebpf_gtp_manager(config: Dict) -> Optional[EbpfGtpManager]:
    """
    Get eBPF GTP manager instance (thread-safe singleton pattern)

    Args:
        config: Pipelined configuration

    Returns:
        EbpfGtpManager instance if enabled, None otherwise
    """
    global _ebpf_gtp_manager_instance

    # Fast path - return existing instance if available
    if _ebpf_gtp_manager_instance is not None:
        # Return existing instance even if maps aren't initialized yet
        # (the controller may have created it)
        LOG.debug("Returning existing EbpfGtpManager instance")
        return _ebpf_gtp_manager_instance
    
    # Slow path - need to create or recreate instance
    with _ebpf_gtp_manager_lock:
        # Double-check after acquiring lock
        if _ebpf_gtp_manager_instance is not None:
            if hasattr(_ebpf_gtp_manager_instance, '_maps_initialized') and _ebpf_gtp_manager_instance._maps_initialized:
                LOG.debug("Another thread created instance while waiting for lock")
                return _ebpf_gtp_manager_instance
            else:
                LOG.warning("Existing instance found but BPF maps not properly initialized")
                LOG.warning("Recreating EbpfGtpManager instance")
                _ebpf_gtp_manager_instance = None  # Reset for recreation
        
        # Check for conflicts with existing eBPF manager
        if config.get('ebpf', {}).get('enabled', False):
            LOG.error("eBPF GTP cannot be enabled when original eBPF manager is active")
            LOG.error("Please disable 'ebpf.enabled' in pipelined.yml before enabling eBPF GTP")
            return None
        
        # Check for kernel GTP conflicts
        try:
            import subprocess
            result = subprocess.run(['lsmod'], capture_output=True, text=True)
            if 'gtp' in result.stdout:
                LOG.warning("Kernel GTP module detected, attempting to remove...")
                subprocess.run(['rmmod', 'gtp'], capture_output=True)
                LOG.info("Kernel GTP module removed")
        except Exception as e:
            LOG.warning(f"Could not check/remove kernel GTP module: {e}")
        
        if not config.get('ebpf_gtp', {}).get('enabled', False):
            LOG.info("eBPF GTP not enabled in configuration")
            return None
        
        # Create new instance (still under lock)
        manager = EbpfGtpManager(config)
        if manager.initialize():
            _ebpf_gtp_manager_instance = manager  # Store global reference
            LOG.info("Created and stored EbpfGtpManager singleton instance")
            return manager
        else:
            LOG.error("Failed to initialize eBPF GTP manager")
            return None


def reset_ebpf_gtp_manager():
    """
    Reset the global eBPF GTP manager instance (for testing/cleanup)
    """
    global _ebpf_gtp_manager_instance
    if _ebpf_gtp_manager_instance is not None:
        LOG.info("Resetting EbpfGtpManager singleton instance")
        _ebpf_gtp_manager_instance = None


def get_current_ebpf_gtp_manager() -> Optional[EbpfGtpManager]:
    """
    Get the current eBPF GTP manager instance without creating a new one
    
    Returns:
        EbpfGtpManager instance if available, None otherwise
    """
    global _ebpf_gtp_manager_instance
    return _ebpf_gtp_manager_instance
