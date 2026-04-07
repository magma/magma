#!/usr/bin/env python3
"""
GTP-U Packet Building and Parsing Tests

Tests the lib/gtpu.py module for spec compliance with 3GPP TS 29.281.

These tests do NOT require root - they only test packet construction.

Run: python3 -m pytest test_gtpu_packets.py -v
"""

import struct
import unittest
import hashlib
import random
import string

# Import may fail if scapy not installed - tests will be skipped
try:
    from scapy.all import Ether, IP, UDP, Raw
    from scapy.contrib.gtp import GTP_U_Header
    SCAPY_AVAILABLE = True
except ImportError:
    SCAPY_AVAILABLE = False

# Import our library
try:
    from lib.gtpu import (
        GTPUConstants,
        build_gtpu_header,
        parse_gtpu_header,
        build_gtpu_packet,
        build_gtpu_packet_raw,
        parse_gtpu_packet,
        extract_inner_payload,
    )
    LIB_AVAILABLE = True
except ImportError:
    LIB_AVAILABLE = False


def generate_random_payload(size: int = 64) -> bytes:
    """Generate random payload for testing"""
    data = ''.join(random.choices(string.ascii_letters + string.digits, k=size))
    return data.encode('utf-8')


@unittest.skipUnless(LIB_AVAILABLE, "lib.gtpu not available")
class TestGTPUConstants(unittest.TestCase):
    """Test GTP-U constants match the spec"""

    def test_port_number(self):
        """TS 29.281 Section 4.4.2.1: UDP port 2152"""
        self.assertEqual(GTPUConstants.PORT, 2152)

    def test_version(self):
        """TS 29.281 Section 5.1: Version must be 1"""
        self.assertEqual(GTPUConstants.VERSION, 1)

    def test_message_types(self):
        """TS 29.281 Table 6.1-1: Message type values"""
        self.assertEqual(GTPUConstants.MSG_ECHO_REQUEST, 1)
        self.assertEqual(GTPUConstants.MSG_ECHO_RESPONSE, 2)
        self.assertEqual(GTPUConstants.MSG_ERROR_INDICATION, 26)
        self.assertEqual(GTPUConstants.MSG_END_MARKER, 254)
        self.assertEqual(GTPUConstants.MSG_GPDU, 255)

    def test_header_lengths(self):
        """TS 29.281 Section 5.1: Minimum header is 8 bytes"""
        self.assertEqual(GTPUConstants.HEADER_MIN_LEN, 8)
        self.assertEqual(GTPUConstants.HEADER_OPT_LEN, 12)


@unittest.skipUnless(LIB_AVAILABLE, "lib.gtpu not available")
class TestGTPUHeaderBuilding(unittest.TestCase):
    """Test GTP-U header construction"""

    def test_minimal_header(self):
        """Build minimal 8-byte header"""
        teid = 0x12345678
        payload_len = 100

        header = build_gtpu_header(teid, payload_len)

        self.assertEqual(len(header), 8)

        # Parse and verify
        flags, msg_type, length, parsed_teid = struct.unpack('!BBHI', header)

        # Check version = 1 (bits 5-7)
        version = (flags >> 5) & 0x07
        self.assertEqual(version, 1)

        # Check PT = 1 (bit 4)
        pt = (flags >> 4) & 0x01
        self.assertEqual(pt, 1)

        # Check E/S/PN flags are 0
        self.assertEqual(flags & 0x07, 0)

        # Check message type (G-PDU)
        self.assertEqual(msg_type, 255)

        # Check length
        self.assertEqual(length, payload_len)

        # Check TEID
        self.assertEqual(parsed_teid, teid)

    def test_header_with_sequence_number(self):
        """Build header with sequence number (S flag)"""
        teid = 0xDEADBEEF
        payload_len = 50
        seq_num = 12345

        header = build_gtpu_header(teid, payload_len, seq_num=seq_num)

        # Should be 12 bytes with optional fields
        self.assertEqual(len(header), 12)

        flags = header[0]

        # S flag should be set
        s_flag = (flags >> 1) & 0x01
        self.assertEqual(s_flag, 1)

        # Parse sequence number
        seq_parsed = struct.unpack('!H', header[8:10])[0]
        self.assertEqual(seq_parsed, seq_num)

    def test_header_with_extension(self):
        """Build header with extension header flag"""
        teid = 0x11223344
        payload_len = 64
        next_ext = 0x85  # PDU Session Container

        header = build_gtpu_header(teid, payload_len, next_ext=next_ext)

        self.assertEqual(len(header), 12)

        flags = header[0]

        # E flag should be set
        e_flag = (flags >> 2) & 0x01
        self.assertEqual(e_flag, 1)

        # Parse next extension header type
        ext_parsed = header[11]
        self.assertEqual(ext_parsed, next_ext)


@unittest.skipUnless(LIB_AVAILABLE, "lib.gtpu not available")
class TestGTPUHeaderParsing(unittest.TestCase):
    """Test GTP-U header parsing"""

    def test_parse_minimal_header(self):
        """Parse minimal 8-byte header"""
        # Build a header
        header = build_gtpu_header(0x12345678, 100)

        # Parse it
        info = parse_gtpu_header(header)

        self.assertIsNotNone(info)
        self.assertEqual(info.version, 1)
        self.assertEqual(info.pt, 1)
        self.assertFalse(info.e_flag)
        self.assertFalse(info.s_flag)
        self.assertFalse(info.pn_flag)
        self.assertEqual(info.msg_type, 255)
        self.assertEqual(info.teid, 0x12345678)
        self.assertEqual(info.header_len, 8)

    def test_parse_header_with_options(self):
        """Parse header with optional fields"""
        header = build_gtpu_header(0xAABBCCDD, 200, seq_num=9999)

        info = parse_gtpu_header(header)

        self.assertIsNotNone(info)
        self.assertTrue(info.s_flag)
        self.assertEqual(info.seq_num, 9999)
        self.assertEqual(info.header_len, 12)

    def test_parse_invalid_header(self):
        """Parsing too-short data returns None"""
        info = parse_gtpu_header(b'\x00\x00\x00')  # Too short

        self.assertIsNone(info)

    def test_roundtrip(self):
        """Build and parse roundtrip"""
        teid = 0x99887766
        payload_len = 512
        seq_num = 42

        header = build_gtpu_header(teid, payload_len, seq_num=seq_num)
        info = parse_gtpu_header(header)

        self.assertEqual(info.teid, teid)
        self.assertEqual(info.length, payload_len + 4)  # +4 for optional fields
        self.assertEqual(info.seq_num, seq_num)


@unittest.skipUnless(SCAPY_AVAILABLE and LIB_AVAILABLE, "scapy or lib not available")
class TestGTPUPacketBuilding(unittest.TestCase):
    """Test complete GTP-U packet construction"""

    def test_build_packet_scapy(self):
        """Build packet using Scapy GTP"""
        payload = b"Hello GTP-U!"

        pkt_bytes, orig_payload = build_gtpu_packet(
            outer_src_ip="10.0.2.1",
            outer_dst_ip="10.0.2.2",
            inner_src_ip="192.168.128.100",
            inner_dst_ip="8.8.8.8",
            teid=0x12345678,
            payload=payload
        )

        self.assertIsInstance(pkt_bytes, bytes)
        self.assertEqual(orig_payload, payload)

        # Parse to verify structure
        from scapy.all import Ether
        pkt = Ether(pkt_bytes)

        self.assertIn(IP, pkt)
        self.assertIn(UDP, pkt)
        self.assertIn(GTP_U_Header, pkt)

        # Check outer IP
        outer_ip = pkt[IP]
        self.assertEqual(outer_ip.src, "10.0.2.1")
        self.assertEqual(outer_ip.dst, "10.0.2.2")

        # Check UDP port
        udp = pkt[UDP]
        self.assertEqual(udp.dport, 2152)

        # Check GTP header
        gtp = pkt[GTP_U_Header]
        self.assertEqual(gtp.teid, 0x12345678)
        self.assertEqual(gtp.gtp_type, 255)

    def test_build_packet_raw(self):
        """Build packet using raw header construction"""
        payload = b"Raw packet test"

        pkt_bytes, orig_payload = build_gtpu_packet_raw(
            outer_src_ip="10.0.2.1",
            outer_dst_ip="10.0.2.2",
            inner_src_ip="192.168.128.100",
            inner_dst_ip="8.8.8.8",
            teid=0xDEADBEEF,
            payload=payload
        )

        # Verify it's a valid packet
        from scapy.all import Ether
        pkt = Ether(pkt_bytes)

        self.assertIn(IP, pkt)
        self.assertIn(UDP, pkt)

        # Payload should be somewhere in the packet
        self.assertIn(payload, pkt_bytes)

    def test_payload_preserved(self):
        """Original payload preserved through build"""
        payload = generate_random_payload(128)
        payload_hash = hashlib.md5(payload).hexdigest()

        pkt_bytes, returned_payload = build_gtpu_packet(
            outer_src_ip="10.0.0.1",
            outer_dst_ip="10.0.0.2",
            inner_src_ip="192.168.1.1",
            inner_dst_ip="1.1.1.1",
            teid=0x11111111,
            payload=payload
        )

        # Returned payload should match
        self.assertEqual(returned_payload, payload)

        # Hash should match
        self.assertEqual(hashlib.md5(returned_payload).hexdigest(), payload_hash)

        # Payload should be in packet bytes
        self.assertIn(payload, pkt_bytes)


@unittest.skipUnless(SCAPY_AVAILABLE and LIB_AVAILABLE, "scapy or lib not available")
class TestGTPUPacketParsing(unittest.TestCase):
    """Test GTP-U packet parsing"""

    def test_parse_gtpu_packet(self):
        """Parse a GTP-U packet"""
        payload = b"Test payload data"

        pkt_bytes, _ = build_gtpu_packet(
            outer_src_ip="10.0.2.1",
            outer_dst_ip="10.0.2.2",
            inner_src_ip="192.168.128.100",
            inner_dst_ip="8.8.8.8",
            teid=0x12345678,
            payload=payload
        )

        pkt = Ether(pkt_bytes)
        info = parse_gtpu_packet(pkt)

        self.assertIsNotNone(info)
        self.assertEqual(info.outer_src_ip, "10.0.2.1")
        self.assertEqual(info.outer_dst_ip, "10.0.2.2")
        self.assertEqual(info.header.teid, 0x12345678)
        self.assertEqual(info.inner_src_ip, "192.168.128.100")
        self.assertEqual(info.inner_dst_ip, "8.8.8.8")

    def test_extract_inner_payload(self):
        """Extract inner payload from GTP-U packet"""
        payload = b"Extract me!"

        pkt_bytes, _ = build_gtpu_packet(
            outer_src_ip="10.0.2.1",
            outer_dst_ip="10.0.2.2",
            inner_src_ip="192.168.128.100",
            inner_dst_ip="8.8.8.8",
            teid=0x12345678,
            payload=payload
        )

        pkt = Ether(pkt_bytes)
        extracted = extract_inner_payload(pkt)

        self.assertIsNotNone(extracted)
        self.assertEqual(extracted, payload)

    def test_non_gtpu_packet(self):
        """Parsing non-GTP-U packet returns None"""
        # Build a plain IP packet
        pkt = Ether() / IP(src="1.1.1.1", dst="2.2.2.2") / UDP(sport=1234, dport=80) / Raw(b"plain")

        info = parse_gtpu_packet(pkt)
        self.assertIsNone(info)


@unittest.skipUnless(LIB_AVAILABLE, "lib.gtpu not available")
class TestEdgeCases(unittest.TestCase):
    """Test edge cases and boundary conditions"""

    def test_zero_teid_rejected_in_config(self):
        """TEID 0 is reserved per spec (except for Echo)"""
        # Note: We don't reject TEID 0 in header building,
        # but config validation should catch it
        from lib.config import GTPUTestConfig

        with self.assertRaises(ValueError):
            GTPUTestConfig(teid_ul=0)

    def test_large_payload(self):
        """Handle large payloads"""
        payload = generate_random_payload(1400)  # Near MTU

        header = build_gtpu_header(0x12345678, len(payload))
        info = parse_gtpu_header(header)

        self.assertEqual(info.length, len(payload))

    def test_min_payload(self):
        """Handle minimum payload"""
        payload = b"X"

        header = build_gtpu_header(0x12345678, len(payload))
        info = parse_gtpu_header(header)

        self.assertEqual(info.length, 1)


if __name__ == '__main__':
    unittest.main(verbosity=2)
