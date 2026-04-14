"""
GTP-U Protocol Module

Spec-compliant GTP-U packet building and parsing using Scapy.
Based on 3GPP TS 29.281 v16.2.0

Reference: .tmp/ebpf_review/reference/3gpp_ts_29.281_v16.2.0.md
"""

import struct
import socket
from dataclasses import dataclass
from typing import Optional, Tuple

# Scapy imports
from scapy.all import Ether, IP, UDP, Raw, Packet
from scapy.contrib.gtp import GTP_U_Header


class GTPUConstants:
    """
    GTP-U constants per 3GPP TS 29.281

    Section 5.1: Header format
    Section 6.1: Message types (Table 6.1-1)
    Section 4.4.2.1: UDP port
    """
    # UDP Port (Section 4.4.2.1)
    PORT = 2152

    # Version (Section 5.1) - must be 1 for GTPv1-U
    VERSION = 1

    # Protocol Type (Section 5.1) - 1 = GTP, 0 = GTP'
    PT_GTP = 1
    PT_GTP_PRIME = 0

    # Message Types (Table 6.1-1)
    MSG_ECHO_REQUEST = 1
    MSG_ECHO_RESPONSE = 2
    MSG_ERROR_INDICATION = 26
    MSG_SUPPORTED_EXT_HDR = 31
    MSG_END_MARKER = 254
    MSG_GPDU = 255  # G-PDU: User data

    # Header lengths
    HEADER_MIN_LEN = 8   # Minimum header (no optional fields)
    HEADER_OPT_LEN = 12  # With optional fields (seq, npdu, ext)


@dataclass
class GTPUHeaderInfo:
    """Parsed GTP-U header information"""
    version: int
    pt: int              # Protocol Type
    e_flag: bool         # Extension header flag
    s_flag: bool         # Sequence number flag
    pn_flag: bool        # N-PDU number flag
    msg_type: int
    length: int          # Payload length (excludes first 8 bytes)
    teid: int
    seq_num: Optional[int] = None
    npdu_num: Optional[int] = None
    next_ext: Optional[int] = None
    header_len: int = 8  # Actual header length


@dataclass
class GTPUPacketInfo:
    """Parsed GTP-U packet information"""
    header: GTPUHeaderInfo
    outer_src_ip: str
    outer_dst_ip: str
    inner_payload: bytes
    inner_src_ip: Optional[str] = None
    inner_dst_ip: Optional[str] = None


def build_gtpu_header(
    teid: int,
    payload_len: int,
    msg_type: int = GTPUConstants.MSG_GPDU,
    seq_num: Optional[int] = None,
    npdu_num: Optional[int] = None,
    next_ext: Optional[int] = None,
) -> bytes:
    """
    Build GTP-U header per 3GPP TS 29.281 Section 5.1

    Octet 1: |Version(3)|PT(1)|*(1)|E(1)|S(1)|PN(1)|
    Octet 2: Message Type
    Octet 3-4: Length
    Octet 5-8: TEID
    [Octet 9-12: Seq Num, N-PDU, Next Ext - if any flag set]

    Args:
        teid: Tunnel Endpoint Identifier
        payload_len: Length of payload after GTP header
        msg_type: Message type (default: G-PDU = 255)
        seq_num: Optional sequence number (sets S flag)
        npdu_num: Optional N-PDU number (sets PN flag)
        next_ext: Optional next extension header type (sets E flag)

    Returns:
        GTP-U header bytes
    """
    # Determine flags
    e_flag = 1 if next_ext is not None else 0
    s_flag = 1 if seq_num is not None else 0
    pn_flag = 1 if npdu_num is not None else 0

    # Build flags byte: Version(3) | PT(1) | *(1) | E(1) | S(1) | PN(1)
    flags = (GTPUConstants.VERSION << 5) | \
            (GTPUConstants.PT_GTP << 4) | \
            (0 << 3) | \
            (e_flag << 2) | \
            (s_flag << 1) | \
            pn_flag

    # If any optional flag is set, include optional fields
    if e_flag or s_flag or pn_flag:
        length = payload_len + 4  # Optional fields are part of "payload" length
        header = struct.pack('!BBHI',
            flags,
            msg_type,
            length,
            teid
        )
        # Add optional fields
        seq = seq_num if seq_num is not None else 0
        npdu = npdu_num if npdu_num is not None else 0
        ext = next_ext if next_ext is not None else 0
        header += struct.pack('!HBB', seq, npdu, ext)
    else:
        # Minimum 8-byte header
        header = struct.pack('!BBHI',
            flags,
            msg_type,
            payload_len,
            teid
        )

    return header


def parse_gtpu_header(data: bytes) -> Optional[GTPUHeaderInfo]:
    """
    Parse GTP-U header per 3GPP TS 29.281 Section 5.1

    Args:
        data: Raw bytes starting at GTP-U header

    Returns:
        GTPUHeaderInfo or None if invalid
    """
    if len(data) < GTPUConstants.HEADER_MIN_LEN:
        return None

    flags, msg_type, length, teid = struct.unpack('!BBHI', data[:8])

    version = (flags >> 5) & 0x07
    pt = (flags >> 4) & 0x01
    e_flag = bool((flags >> 2) & 0x01)
    s_flag = bool((flags >> 1) & 0x01)
    pn_flag = bool(flags & 0x01)

    result = GTPUHeaderInfo(
        version=version,
        pt=pt,
        e_flag=e_flag,
        s_flag=s_flag,
        pn_flag=pn_flag,
        msg_type=msg_type,
        length=length,
        teid=teid,
        header_len=8
    )

    # Parse optional fields if any flag is set
    if e_flag or s_flag or pn_flag:
        if len(data) < GTPUConstants.HEADER_OPT_LEN:
            return None
        seq_num, npdu_num, next_ext = struct.unpack('!HBB', data[8:12])
        result.seq_num = seq_num if s_flag else None
        result.npdu_num = npdu_num if pn_flag else None
        result.next_ext = next_ext if e_flag else None
        result.header_len = 12

    return result


def build_gtpu_packet(
    outer_src_ip: str,
    outer_dst_ip: str,
    inner_src_ip: str,
    inner_dst_ip: str,
    teid: int,
    payload: bytes,
    inner_sport: int = 12345,
    inner_dport: int = 80,
    src_mac: str = "02:00:00:00:00:01",
    dst_mac: str = "02:00:00:00:00:02",
) -> Tuple[bytes, bytes]:
    """
    Build complete GTP-U encapsulated packet using Scapy.

    Packet structure:
        Ethernet | Outer IP | UDP(2152) | GTP-U | Inner IP | Inner UDP | Payload

    Args:
        outer_src_ip: Outer IP source (e.g., eNB IP)
        outer_dst_ip: Outer IP destination (e.g., SGW IP)
        inner_src_ip: Inner IP source (e.g., UE IP)
        inner_dst_ip: Inner IP destination (e.g., external IP)
        teid: GTP Tunnel Endpoint Identifier
        payload: User payload data
        inner_sport: Inner UDP source port
        inner_dport: Inner UDP destination port
        src_mac: Source MAC address
        dst_mac: Destination MAC address

    Returns:
        Tuple of (complete_packet_bytes, original_payload)
    """
    # Build inner packet (what the UE "sent")
    inner_pkt = IP(src=inner_src_ip, dst=inner_dst_ip) / \
                UDP(sport=inner_sport, dport=inner_dport) / \
                Raw(load=payload)

    # Build outer packet with GTP-U header
    # Using Scapy's GTP_U_Header
    outer_pkt = Ether(src=src_mac, dst=dst_mac) / \
                IP(src=outer_src_ip, dst=outer_dst_ip) / \
                UDP(sport=GTPUConstants.PORT, dport=GTPUConstants.PORT) / \
                GTP_U_Header(teid=teid, gtp_type=GTPUConstants.MSG_GPDU) / \
                inner_pkt

    return bytes(outer_pkt), payload


def build_gtpu_packet_raw(
    outer_src_ip: str,
    outer_dst_ip: str,
    inner_src_ip: str,
    inner_dst_ip: str,
    teid: int,
    payload: bytes,
    inner_sport: int = 12345,
    inner_dport: int = 80,
    src_mac: str = "02:00:00:00:00:01",
    dst_mac: str = "02:00:00:00:00:02",
) -> Tuple[bytes, bytes]:
    """
    Build GTP-U packet using raw header construction.

    Use this when you need precise control over the GTP header
    or when testing header parsing.

    Returns:
        Tuple of (complete_packet_bytes, original_payload)
    """
    # Build inner packet
    inner_pkt = IP(src=inner_src_ip, dst=inner_dst_ip) / \
                UDP(sport=inner_sport, dport=inner_dport) / \
                Raw(load=payload)
    inner_bytes = bytes(inner_pkt)

    # Build GTP-U header manually
    gtpu_hdr = build_gtpu_header(teid, len(inner_bytes))

    # Build outer packet
    outer_pkt = Ether(src=src_mac, dst=dst_mac) / \
                IP(src=outer_src_ip, dst=outer_dst_ip) / \
                UDP(sport=GTPUConstants.PORT, dport=GTPUConstants.PORT) / \
                Raw(load=gtpu_hdr + inner_bytes)

    return bytes(outer_pkt), payload


def parse_gtpu_packet(pkt: Packet) -> Optional[GTPUPacketInfo]:
    """
    Parse a GTP-U packet (Scapy Packet object).

    Args:
        pkt: Scapy packet object

    Returns:
        GTPUPacketInfo or None if not a valid GTP-U packet
    """
    if IP not in pkt:
        return None

    outer_ip = pkt[IP]

    if UDP not in pkt:
        return None

    udp = pkt[UDP]
    if udp.dport != GTPUConstants.PORT and udp.sport != GTPUConstants.PORT:
        return None

    # Try to parse GTP-U header
    if GTP_U_Header in pkt:
        gtp = pkt[GTP_U_Header]
        header = GTPUHeaderInfo(
            version=gtp.version,
            pt=gtp.PT,
            e_flag=bool(gtp.E),
            s_flag=bool(gtp.S),
            pn_flag=bool(gtp.PN),
            msg_type=gtp.gtp_type,
            length=gtp.length,
            teid=gtp.teid,
            header_len=8 if not (gtp.E or gtp.S or gtp.PN) else 12
        )

        # Extract inner payload
        inner_payload = bytes(gtp.payload) if gtp.payload else b''

        # Try to parse inner IP
        inner_src = None
        inner_dst = None
        if len(inner_payload) >= 20:
            # Check if it looks like an IP packet
            version = (inner_payload[0] >> 4) & 0x0F
            if version == 4:
                inner_src = socket.inet_ntoa(inner_payload[12:16])
                inner_dst = socket.inet_ntoa(inner_payload[16:20])

        return GTPUPacketInfo(
            header=header,
            outer_src_ip=outer_ip.src,
            outer_dst_ip=outer_ip.dst,
            inner_payload=inner_payload,
            inner_src_ip=inner_src,
            inner_dst_ip=inner_dst
        )

    # Fallback: parse from raw UDP payload
    if Raw in pkt:
        raw_data = bytes(pkt[Raw])
        header = parse_gtpu_header(raw_data)
        if header is None:
            return None

        inner_payload = raw_data[header.header_len:]

        inner_src = None
        inner_dst = None
        if len(inner_payload) >= 20:
            version = (inner_payload[0] >> 4) & 0x0F
            if version == 4:
                inner_src = socket.inet_ntoa(inner_payload[12:16])
                inner_dst = socket.inet_ntoa(inner_payload[16:20])

        return GTPUPacketInfo(
            header=header,
            outer_src_ip=outer_ip.src,
            outer_dst_ip=outer_ip.dst,
            inner_payload=inner_payload,
            inner_src_ip=inner_src,
            inner_dst_ip=inner_dst
        )

    return None


def extract_inner_payload(pkt: Packet) -> Optional[bytes]:
    """
    Extract the innermost payload from a packet.

    Works for both GTP-U encapsulated and plain packets.

    Args:
        pkt: Scapy packet

    Returns:
        Inner payload bytes or None
    """
    # Check for GTP-U first
    gtpu_info = parse_gtpu_packet(pkt)
    if gtpu_info:
        # Find the actual user payload inside inner packet
        inner = gtpu_info.inner_payload
        if len(inner) > 28:  # IP(20) + UDP(8) minimum
            # Skip inner IP and UDP headers to get raw payload
            ip_hlen = (inner[0] & 0x0F) * 4
            udp_offset = ip_hlen + 8
            return inner[udp_offset:]
        return inner

    # Not GTP-U, try to find Raw payload
    if Raw in pkt:
        return bytes(pkt[Raw])

    return None
