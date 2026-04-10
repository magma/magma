import socket
import struct
import logging
from typing import Optional

from ebpf_utils import write_bpf_map, delete_bpf_map_entry

LOG = logging.getLogger(__name__)


class EbpfGtpManager:
    """
    Userspace manager for GTP eBPF datapath.

    Owns lifecycle of UE sessions stored in ue_session_map.
    """

    def __init__(self, ue_session_map):
        """
        ue_session_map: handle or fd passed from loader (bcc/libbpf)
        """
        self.ue_session_map = ue_session_map

    # ------------------------------------------------------------------
    # ADD
    # ------------------------------------------------------------------
    def add_ue_session(
        self,
        ue_ip: str,
        enb_ip: str,
        teid_ul_in: int,
        teid_ul_out: int,
        teid_dl_in: int,
        teid_dl_out: int,
        s1u_ifindex: int,
        sgi_ifindex: int,
        ovs_ifindex: int,
        ul_mac_src: bytes,
        ul_mac_dst: bytes,
        qfi: int = 0,
        active: bool = True,
    ) -> None:
        """
        Create or update a UE session in ue_session_map.
        """

        # ---- key (struct ue_session_key { __be32 ue_ip; }) ----
        ue_ip_be = struct.unpack("!I", socket.inet_aton(ue_ip))[0]
        key = struct.pack("!I", ue_ip_be)

        session_flags = 1 if active else 0

        # ---- value (struct ue_session_info) ----
        value = struct.pack(
            # enb_ip
            "!I"
            # teid_ul_in, teid_ul_out, teid_dl_in, teid_dl_out
            "IIII"
            # s1u_ifindex, sgi_ifindex, ovs_ifindex
            "III"
            # ul_mac_src, ul_mac_dst
            "6s6s"
            # qos_mark, bearer_id
            "II"
            # ul_bytes, dl_bytes, ul_packets, dl_packets
            "QQQQ"
            # last_seen, session_flags
            "QI"
            # imsi[16], imsi_len
            "16sI"
            # encoded_imsi
            "Q"
            # qfi
            "B"
            # tunnel_id
            "I"
            # tun_ipv4_dst
            "I"
            # tun_flags
            "B"
            # direction
            "B"
            # original_port
            "I"
            # reserved[3]
            "3s"
            # metadata_mark (network order, set later by decap)
            "I",
            # ---- values ----
            struct.unpack("!I", socket.inet_aton(enb_ip))[0],
            teid_ul_in,
            teid_ul_out,
            teid_dl_in,
            teid_dl_out,
            s1u_ifindex,
            sgi_ifindex,
            ovs_ifindex,
            ul_mac_src,
            ul_mac_dst,
            0,                  # qos_mark
            0,                  # bearer_id
            0, 0, 0, 0,          # counters
            0,                  # last_seen
            session_flags,       # session_flags
            b"",                 # imsi
            0,                  # imsi_len
            0,                  # encoded_imsi
            qfi,                # qfi
            0,                  # tunnel_id
            0,                  # tun_ipv4_dst
            0,                  # tun_flags
            0,                  # direction
            0,                  # original_port
            b"\x00\x00\x00",
            0,                  # metadata_mark (filled by gtp_decap)
        )

        LOG.info("Adding UE session for %s", ue_ip)
        write_bpf_map(self.ue_session_map, key, value)

    # ------------------------------------------------------------------
    # REMOVE
    # ------------------------------------------------------------------
    def remove_ue_session(self, ue_ip: str) -> None:
        """
        Remove a UE session from ue_session_map.
        """

        ue_ip_be = struct.unpack("!I", socket.inet_aton(ue_ip))[0]
        key = struct.pack("!I", ue_ip_be)

        LOG.info("Removing UE session for %s", ue_ip)
        delete_bpf_map_entry(self.ue_session_map, key)