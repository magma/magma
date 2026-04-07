#!/usr/bin/env python3
"""
Production eBPF GTP Encapsulation Tests

Tests the ACTUAL ebpf_gtp_encap.c program — compiles it with BCC,
loads it, attaches to a veth pair, sends plain IP packets, and
verifies GTP-U encapsulated output.

Topology:
    [send plain IP] --> ovs_veth --[gtp_encap_handler egress]--> s1u_veth --> [capture GTP-U]

The encap handler looks up the UE session by destination IP, adds
outer IP/UDP/GTP/extension headers (58 bytes total), and redirects
the encapsulated packet to s1u_ifindex via bpf_redirect().

Requirements:
- Root access (sudo)
- BCC library (python3-bpfcc)
- Scapy with GTP support
- Linux kernel >= 4.18
- Production ebpf_gtp_encap.c in parent directory

Run: sudo python3 -m pytest test_production_encap.py -v -s
"""

import os
import sys
import time
import unittest
import random
import string
import struct

if os.geteuid() != 0:
    raise unittest.SkipTest("Root access required")

try:
    from bcc import BPF
    BCC_AVAILABLE = True
except ImportError:
    BCC_AVAILABLE = False

try:
    from scapy.all import Ether, IP, UDP, ARP, Raw, conf
    from scapy.contrib.gtp import GTP_U_Header
    conf.verb = 0
    SCAPY_AVAILABLE = True
except ImportError:
    SCAPY_AVAILABLE = False

try:
    from lib.config import (
        GTPUTestConfig, ENCAP_SOURCE, PRODUCTION_CFLAGS,
        load_production_source, populate_session, set_config, read_stat,
        ip_to_host_order, golden_encap, attach_tc, detach_tc, has_map,
        STATS_TOTAL_PROCESSED, STATS_GTP_ENCAP_SUCCESS, STATS_SESSION_MISS,
        STATS_PKT_FORWARDED, STATS_PKT_DROPPED, STATS_PKT_TOO_SHORT,
        STATS_INACTIVE_SESSION,
        STATS_DOUBLE_ENCAP_AVOIDED, STATS_DL_ERRORS,
        CONFIG_SGI_IP,
    )
    from lib.network import veth_pair, setup_tc_qdisc, get_ifindex, get_mac_address
    from lib.capture import PacketCapture
    LIB_AVAILABLE = True
except ImportError as e:
    LIB_AVAILABLE = False
    IMPORT_ERROR = str(e)


def _generate_payload(size=64):
    data = ''.join(random.choices(string.ascii_letters + string.digits, k=size))
    return data.encode('utf-8')


@unittest.skipUnless(BCC_AVAILABLE, "BCC not available")
class TestEncapCompilation(unittest.TestCase):
    """Test that the production encap program compiles and has expected maps"""

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(ENCAP_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {ENCAP_SOURCE}")
        cls.source = load_production_source(ENCAP_SOURCE)

    def test_01_compiles(self):
        """ebpf_gtp_encap.c compiles with production cflags"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        self.assertIsNotNone(bpf)

    def test_02_has_encap_handler(self):
        """gtp_encap_handler function exists and is loadable"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        fn = bpf.load_func("gtp_encap_handler", BPF.SCHED_CLS)
        self.assertIsNotNone(fn)

    def test_03_has_session_map(self):
        """ue_session_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_encap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "ue_session_map"))

    def test_04_has_config_map(self):
        """config_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_encap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "config_map"))

    def test_05_has_stats_map(self):
        """stats_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_encap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "stats_map"))


@unittest.skipUnless(LIB_AVAILABLE, "Test library not available")
@unittest.skipUnless(BCC_AVAILABLE, "BCC not available")
@unittest.skipUnless(SCAPY_AVAILABLE, "Scapy not available")
class TestEncapEndToEnd(unittest.TestCase):
    """
    End-to-end tests of the production GTP encapsulation program.

    Each test:
    1. Creates veth pair (ovs_veth <-> s1u_veth)
    2. Compiles and attaches gtp_encap_handler on ovs_veth egress
    3. Configures session and config maps
    4. Sends plain IP packets on ovs_veth
    5. Captures GTP-U packets on s1u_veth
    6. Validates GTP headers, TEID, outer IPs, extension headers
    """

    _test_counter = 0

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(ENCAP_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {ENCAP_SOURCE}")
        cls.source = load_production_source(ENCAP_SOURCE)

    def setUp(self):
        """Use unique veth names per test to avoid kernel state conflicts."""
        TestEncapEndToEnd._test_counter += 1
        n = TestEncapEndToEnd._test_counter
        self.config = GTPUTestConfig(
            tx_iface=f"enc_ovs{n}",
            rx_iface=f"enc_s1u{n}",
        )

    def _setup_encap(self):
        """Compile, load, attach encap program on tx_iface egress."""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        attach_tc(bpf, self.config.tx_iface, "gtp_encap_handler", "egress")
        return bpf

    def _add_session(self, bpf, s1u_ifindex, mac_src=None, mac_dst=None,
                     active=True, teid_dl=None, ue_ip=None, enb_ip=None):
        """Add a UE session for encapsulation."""
        if mac_src is None:
            mac_src = get_mac_address(self.config.rx_iface)
        if mac_dst is None:
            mac_dst = get_mac_address(self.config.rx_iface)

        populate_session(
            bpf,
            ue_ip_str=ue_ip or self.config.ue_ip,
            enb_ip_str=enb_ip or self.config.enb_ip,
            teid_ul_in=self.config.teid_ul,
            teid_dl_out=teid_dl if teid_dl is not None else self.config.teid_dl,
            s1u_ifindex=s1u_ifindex,
            mac_src=mac_src,
            mac_dst=mac_dst,
            active=active,
        )

    def test_01_attach_to_tc(self):
        """Production encap attaches to TC egress"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()
            try:
                import subprocess
                result = subprocess.run(
                    ['tc', 'filter', 'show', 'dev', self.config.tx_iface, 'egress'],
                    capture_output=True, text=True
                )
                self.assertIn('bpf', result.stdout.lower())
                print(f"\n  Attached gtp_encap_handler to {self.config.tx_iface} egress")
            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_02_encap_golden_model(self):
        """Golden model: validate every output field against Scapy reference"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))

                session_mac_src = "02:ee:ff:00:11:01"
                session_mac_dst = "02:ee:ff:00:11:02"
                test_qfi = 7
                self._add_session(bpf, s1u_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)
                # Set QFI on the session
                session_map = bpf.get_table("ue_session_map")
                key = session_map.Key()
                key.ue_ip = ip_to_host_order(self.config.ue_ip)
                session = session_map[key]
                session.qfi = test_qfi
                session_map[key] = session

                payload = _generate_payload(64)

                # Build the input packet and compute its inner IP total length
                input_pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=80, dport=12345) /
                    Raw(load=payload)
                )
                # Force Scapy to calculate lengths
                input_pkt = Ether(bytes(input_pkt))
                inner_ip_total_len = input_pkt[IP].len

                cap = PacketCapture(self.config.rx_iface, filter="udp port 2152")
                cap.start()
                time.sleep(0.1)

                self._send_on_egress(input_pkt)
                time.sleep(0.1)

                packets = cap.stop()

                # --- Stats validation ---
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                forwarded = read_stat(bpf, STATS_PKT_FORWARDED)
                self.assertGreater(encap_ok, 0, "Encap success counter not incremented")
                self.assertGreater(len(packets), 0, "No packets captured")

                # Find GTP-U packet
                encap_pkt = None
                for p in packets:
                    if GTP_U_Header in p:
                        encap_pkt = p
                        break
                self.assertIsNotNone(encap_pkt, "No GTP-U packet found")

                # --- Ethernet header ---
                self.assertEqual(
                    encap_pkt[Ether].dst, session_mac_dst,
                    f"Dst MAC: {encap_pkt[Ether].dst} != {session_mac_dst}")
                self.assertEqual(
                    encap_pkt[Ether].src, session_mac_src,
                    f"Src MAC: {encap_pkt[Ether].src} != {session_mac_src}")
                self.assertEqual(encap_pkt[Ether].type, 0x0800, "EtherType not IPv4")

                # --- Outer IP header ---
                outer_ip = encap_pkt[IP]
                self.assertEqual(outer_ip.version, 4, f"IP version: {outer_ip.version}")
                self.assertEqual(outer_ip.ihl, 5, f"IP IHL: {outer_ip.ihl}")
                self.assertEqual(outer_ip.tos, 0, f"IP TOS: {outer_ip.tos}")
                self.assertEqual(outer_ip.id, 0, f"IP ID: {outer_ip.id}")
                # DF flag: Scapy represents flags as int, DF=0x2
                self.assertTrue(outer_ip.flags.DF, "DF flag not set")
                self.assertEqual(outer_ip.frag, 0, f"Frag offset: {outer_ip.frag}")
                self.assertEqual(outer_ip.ttl, 64, f"TTL: {outer_ip.ttl}")
                self.assertEqual(outer_ip.proto, 17, f"Protocol: {outer_ip.proto} != UDP(17)")
                self.assertEqual(
                    outer_ip.src, self.config.sgw_ip,
                    f"Outer src IP: {outer_ip.src} != {self.config.sgw_ip}")
                self.assertEqual(
                    outer_ip.dst, self.config.enb_ip,
                    f"Outer dst IP: {outer_ip.dst} != {self.config.enb_ip}")

                # IP total length: 20(IP) + 8(UDP) + 8(GTP) + 8(ext) + inner_ip_total_len
                expected_ip_len = 20 + 8 + 8 + 8 + inner_ip_total_len
                self.assertEqual(
                    outer_ip.len, expected_ip_len,
                    f"IP total len: {outer_ip.len} != {expected_ip_len}")

                # IP checksum validation: let Scapy recalculate and compare
                ebpf_checksum = outer_ip.chksum
                # Rebuild to get Scapy's calculated checksum
                outer_ip_copy = outer_ip.copy()
                del outer_ip_copy.chksum
                outer_ip_recalc = IP(bytes(outer_ip_copy))
                self.assertEqual(
                    ebpf_checksum, outer_ip_recalc.chksum,
                    f"IP checksum: eBPF=0x{ebpf_checksum:04x} != "
                    f"Scapy=0x{outer_ip_recalc.chksum:04x}")

                # --- UDP header ---
                udp = encap_pkt[UDP]
                self.assertEqual(udp.sport, 2152, f"UDP sport: {udp.sport}")
                self.assertEqual(udp.dport, 2152, f"UDP dport: {udp.dport}")
                self.assertEqual(udp.chksum, 0, f"UDP checksum: {udp.chksum} (should be 0)")

                # UDP length: 8(UDP) + 8(GTP) + 8(ext) + inner_ip_total_len
                expected_udp_len = 8 + 8 + 8 + inner_ip_total_len
                self.assertEqual(
                    udp.len, expected_udp_len,
                    f"UDP len: {udp.len} != {expected_udp_len}")

                # --- GTP-U header ---
                gtp = encap_pkt[GTP_U_Header]
                self.assertEqual(gtp.gtp_type, 255, f"GTP type: {gtp.gtp_type} != 255 (T-PDU)")
                self.assertEqual(
                    gtp.teid, self.config.teid_dl,
                    f"TEID: 0x{gtp.teid:08X} != 0x{self.config.teid_dl:08X}")

                # GTP length: 8(ext) + inner_ip_total_len
                expected_gtp_len = 8 + inner_ip_total_len
                self.assertEqual(
                    gtp.length, expected_gtp_len,
                    f"GTP length: {gtp.length} != {expected_gtp_len}")

                # --- GTP extension header (PDU Session Container) ---
                # Parse raw bytes to validate extension fields
                pkt_bytes = bytes(encap_pkt)
                # Find GTP header start: after Ether(14) + IP(20) + UDP(8) = offset 42
                gtp_offset = 14 + 20 + 8
                # GTP flags byte
                gtp_flags = pkt_bytes[gtp_offset]
                # Flags should be 0x34: version=1(0x20), PT=1(0x10), E=1(0x04)
                self.assertEqual(
                    gtp_flags, 0x34,
                    f"GTP flags: 0x{gtp_flags:02x} != 0x34 "
                    f"(V={gtp_flags>>5}, PT={(gtp_flags>>4)&1}, E={(gtp_flags>>2)&1})")

                # Optional fields start at gtp_offset + 8
                opt_offset = gtp_offset + 8
                seq_num = struct.unpack("!H", pkt_bytes[opt_offset:opt_offset+2])[0]
                npdu = pkt_bytes[opt_offset + 2]
                next_ext = pkt_bytes[opt_offset + 3]
                self.assertEqual(seq_num, 0, f"Seq num: {seq_num}")
                self.assertEqual(npdu, 0, f"N-PDU: {npdu}")
                self.assertEqual(
                    next_ext, 0x85,
                    f"Next ext type: 0x{next_ext:02x} != 0x85 (PDU Session Container)")

                # PDU Session Container extension (4 bytes after optional fields)
                ext_offset = opt_offset + 4
                ext_len = pkt_bytes[ext_offset]
                pdu_type = pkt_bytes[ext_offset + 1]
                qfi_byte = pkt_bytes[ext_offset + 2]
                next_ext_end = pkt_bytes[ext_offset + 3]

                self.assertEqual(ext_len, 1, f"Ext length: {ext_len} (should be 1 = 4 bytes)")
                self.assertEqual(
                    pdu_type, 0x10,
                    f"PDU type: 0x{pdu_type:02x} != 0x10 (DL PDU SESSION INFO)")
                self.assertEqual(
                    qfi_byte & 0x3F, test_qfi,
                    f"QFI: {qfi_byte & 0x3F} != {test_qfi}")
                self.assertEqual(
                    next_ext_end, 0x00,
                    f"Final next_ext: 0x{next_ext_end:02x} != 0x00")

                # --- Inner packet byte-exact preservation ---
                # Extract inner IP from captured output (after GTP extension)
                inner_offset = ext_offset + 4  # After the 4-byte PDU Session Container
                captured_inner = pkt_bytes[inner_offset:]
                # Original inner IP bytes (input packet minus Ethernet header)
                original_inner = bytes(input_pkt)[14:]  # Skip Ether(14)

                self.assertEqual(
                    captured_inner, original_inner,
                    f"Inner packet not preserved byte-for-byte. "
                    f"Captured {len(captured_inner)} bytes vs "
                    f"original {len(original_inner)} bytes")

                # --- Session counter validation ---
                session = session_map[key]
                self.assertGreater(session.dl_packets, 0,
                                   "Session dl_packets not incremented")
                self.assertGreater(session.dl_bytes, 0,
                                   "Session dl_bytes not incremented")
                self.assertGreater(session.last_seen, 0,
                                   "Session last_seen not updated")

                # --- Golden model byte-exact comparison ---
                expected_bytes = golden_encap(
                    original_inner,
                    session_mac_src, session_mac_dst,
                    self.config.sgw_ip, self.config.enb_ip,
                    self.config.teid_dl, qfi=test_qfi)
                # Strict length and content equality
                self.assertEqual(
                    len(pkt_bytes), len(expected_bytes),
                    f"Encap output length {len(pkt_bytes)} != "
                    f"golden model {len(expected_bytes)}")
                self.assertEqual(
                    pkt_bytes, expected_bytes,
                    "Encap output does not match golden model byte-for-byte")

                print(f"\n  Ethernet: dst={encap_pkt[Ether].dst} src={encap_pkt[Ether].src}")
                print(f"  Outer IP: {outer_ip.src} -> {outer_ip.dst} len={outer_ip.len} ttl={outer_ip.ttl}")
                print(f"  IP checksum: eBPF=0x{ebpf_checksum:04x} Scapy=0x{outer_ip_recalc.chksum:04x}")
                print(f"  UDP: {udp.sport}->{udp.dport} len={udp.len} chksum={udp.chksum}")
                print(f"  GTP: flags=0x{gtp_flags:02x} type={gtp.gtp_type} len={gtp.length} TEID=0x{gtp.teid:08X}")
                print(f"  Extension: type=0x85 len={ext_len} pdu_type=0x{pdu_type:02x} QFI={qfi_byte & 0x3F}")
                print(f"  Golden model: byte-exact match ({len(pkt_bytes)} bytes)")
                print(f"  Session: dl_pkts={session.dl_packets} dl_bytes={session.dl_bytes}")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def _send_on_egress(self, pkt):
        """Send packet on tx_iface egress using raw AF_PACKET socket.

        Scapy's L2Socket can fail with ENOBUFS when the TC eBPF handler
        drops (TC_ACT_SHOT) or redirects (bpf_redirect) on egress, because
        virtual devices process TC synchronously during send().
        A raw AF_PACKET socket tolerates this — the packet enters the TC
        pipeline and the handler runs even if send() returns an error.

        In production, OVS sends via dev_queue_xmit (kernel-internal) which
        does not surface these errors to userspace.
        """
        import socket as _socket
        raw_bytes = bytes(pkt)
        s = _socket.socket(_socket.AF_PACKET, _socket.SOCK_RAW,
                           _socket.htons(0x0003))
        try:
            s.bind((self.config.tx_iface, 0))
            s.send(raw_bytes)
        except OSError:
            pass  # TC processes packet before returning error
        finally:
            s.close()

    # Padding to ensure packets exceed 64 bytes at L2 for tests that need
    # to reach the handler's session logic (past the bpf_skb_load_bytes check).
    _PAD = b'\x00' * 50

    def test_03_encap_small_packet(self):
        """Small valid IP packets to a known UE must be encapsulated.

        Finding 2: The handler calls bpf_skb_load_bytes(skb, 0, pkt_data, 64)
        which drops packets under 64 bytes as PKT_TOO_SHORT. This test sends
        a small but valid IP packet (TCP ACK size) to a UE with an active
        session. It should be encapsulated — if it's dropped, the
        implementation has a bug that affects real traffic (TCP ACKs, ICMP,
        small DNS responses).
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))
                self._add_session(bpf, s1u_ifindex)

                cap = PacketCapture(self.config.rx_iface, filter="udp port 2152")
                cap.start()
                time.sleep(0.1)

                # Small packet — valid IP/UDP to known UE, ~52 bytes at L2
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=80, dport=12345) /
                    Raw(load=b"SMALL")
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                too_short = read_stat(bpf, STATS_PKT_TOO_SHORT)
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                print(f"\n  Stats: total={total}, too_short={too_short}, encap_ok={encap_ok}")

                self.assertEqual(too_short, 0,
                                 f"FINDING 2: Valid small packet dropped as too short. "
                                 f"bpf_skb_load_bytes(64) rejects packets < 64 bytes. "
                                 f"TCP ACKs, ICMP, small DNS would be silently dropped.")
                self.assertGreater(encap_ok, 0,
                                   "Small packet with active session should be encapsulated")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_04_encap_no_session(self):
        """Packet for unknown UE is dropped — verify no packet leaks through"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                marker = b"NO_SESSION_MARKER"
                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src="1.1.1.1", dst="172.16.99.99") /
                    UDP(sport=80, dport=12345) /
                    Raw(load=marker + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                miss = read_stat(bpf, STATS_SESSION_MISS)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: total={total}, session_miss={miss}, dropped={dropped}")

                self.assertGreater(total, 0,
                                   "Handler never ran — packet didn't reach TC pipeline")
                self.assertGreater(miss, 0, "Session miss counter not incremented")

                leaked = [p for p in packets if marker in bytes(p)]
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} packets leaked through — "
                                 f"handler should drop unknown sessions")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_05_encap_inactive_session(self):
        """Packet for inactive session is dropped — verify no packet leaks"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                self._add_session(bpf, s1u_ifindex, active=False)

                marker = b"INACTIVE_MARKER"
                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=80, dport=12345) /
                    Raw(load=marker + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                inactive = read_stat(bpf, STATS_INACTIVE_SESSION)
                print(f"\n  Stats: total={total}, inactive_session={inactive}")

                self.assertGreater(total, 0,
                                   "Handler never ran — packet didn't reach TC pipeline")
                self.assertGreater(inactive, 0, "Inactive session counter not incremented")

                leaked = [p for p in packets if marker in bytes(p)]
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} packets leaked — "
                                 f"handler should drop inactive sessions")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_06_double_encap_avoided(self):
        """Already-GTP packets are not double-encapsulated"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                self._add_session(bpf, s1u_ifindex)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=2152, dport=2152) /
                    Raw(load=b"ALREADY_GTP" + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                avoided = read_stat(bpf, STATS_DOUBLE_ENCAP_AVOIDED)
                print(f"\n  Stats: total={total}, double_encap_avoided={avoided}")

                self.assertGreater(total, 0,
                                   "Handler never ran — packet didn't reach TC pipeline")
                self.assertGreater(avoided, 0, "Double encap avoidance not counted")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_07_multiple_ue_sessions(self):
        """Packets for different UEs get correct TEIDs"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))

                sessions = [
                    ("192.168.128.101", 0x00000101),
                    ("192.168.128.102", 0x00000102),
                    ("192.168.128.103", 0x00000103),
                ]

                for ue_ip, teid in sessions:
                    self._add_session(bpf, s1u_ifindex, teid_dl=teid, ue_ip=ue_ip)

                cap = PacketCapture(self.config.rx_iface, filter="udp port 2152")
                cap.start()
                time.sleep(0.1)

                for ue_ip, _ in sessions:
                    pkt = (
                        Ether(dst="ff:ff:ff:ff:ff:ff") /
                        IP(src=self.config.external_ip, dst=ue_ip) /
                        UDP(sport=80, dport=12345) /
                        Raw(load=f"to_{ue_ip}".encode() + self._PAD)
                    )
                    self._send_on_egress(pkt)
                    time.sleep(0.05)

                time.sleep(0.1)
                packets = cap.stop()

                found_teids = set()
                for pkt in packets:
                    if GTP_U_Header in pkt:
                        found_teids.add(pkt[GTP_U_Header].teid)

                expected_teids = {teid for _, teid in sessions}
                print(f"\n  Expected TEIDs: {[hex(t) for t in expected_teids]}")
                print(f"  Found TEIDs: {[hex(t) for t in found_teids]}")

                self.assertEqual(found_teids, expected_teids,
                                 f"TEID mismatch: expected {expected_teids}, got {found_teids}")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_08_encap_double_encap_with_ipv4_options(self):
        """Double-encap detection fails when outer IPv4 has options.

        Finding 4: The encap handler checks UDP ports at hardcoded offsets
        34-37 (ebpf_gtp_encap.c:197) assuming a 20-byte IPv4 header.
        Protocol and dst IP are at fixed positions in IPv4 and read correctly,
        but UDP ports are AFTER the IP header, so options shift them.
        With IHL=6 (24 bytes), the handler reads option bytes instead of
        UDP ports, misses port 2152, and re-encapsulates an already-GTP packet.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))
                self._add_session(bpf, s1u_ifindex)

                time.sleep(0.1)
                initial_avoided = read_stat(bpf, STATS_DOUBLE_ENCAP_AVOIDED)

                # Already-GTP packet with IPv4 options
                # Without options, this triggers STATS_DOUBLE_ENCAP_AVOIDED.
                # With options, the handler reads wrong port bytes and
                # misses the double-encap check.
                from scapy.all import IPOption
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip,
                       options=[IPOption(b'\x01')] * 4) /
                    UDP(sport=2152, dport=2152) /
                    Raw(load=b"ALREADY_GTP_WITH_OPTIONS" + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                avoided = read_stat(bpf, STATS_DOUBLE_ENCAP_AVOIDED)
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                print(f"\n  Stats: total={total}, "
                      f"double_encap_avoided={avoided}, "
                      f"encap_ok={encap_ok}")

                self.assertGreater(total, 0,
                                   "Handler did not process packet")
                self.assertGreater(
                    avoided, initial_avoided,
                    "FINDING 4: Double-encap detection broken with IPv4 "
                    "options. Handler reads UDP ports at hardcoded offsets "
                    "(bytes 34-37) which point into IP options instead of "
                    "actual UDP header. Already-GTP packet was NOT detected "
                    f"as double-encap (encap_ok={encap_ok}).")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_09_encap_missing_sgi_ip(self):
        """Missing CONFIG_SGI_IP must fail closed — no packet emitted.

        Finding 5: The handler should reject the packet when CONFIG_SGI_IP
        is not configured, since the outer source IP is unknown. Expected:
        encap_ok == 0, an error/drop counter incremented, no GTP packet
        emitted. The current implementation silently uses 0.0.0.0.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                # Intentionally do NOT set CONFIG_SGI_IP
                self._add_session(bpf, s1u_ifindex)

                cap = PacketCapture(self.config.rx_iface, filter="udp port 2152")
                cap.start()
                time.sleep(0.1)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=80, dport=12345) /
                    Raw(load=b"MISSING_SGI_IP_TEST" + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                dl_errors = read_stat(bpf, STATS_DL_ERRORS)
                print(f"\n  Stats: total={total}, encap_ok={encap_ok}, "
                      f"dropped={dropped}, dl_errors={dl_errors}")

                self.assertGreater(total, 0,
                                   "Handler did not process the packet")

                # Missing CONFIG_SGI_IP should prevent encapsulation
                self.assertEqual(
                    encap_ok, 0,
                    "FINDING 5: Missing CONFIG_SGI_IP should fail closed. "
                    "Handler must not encapsulate without a valid source IP.")

                # No GTP packet should be emitted
                marker = b"MISSING_SGI_IP_TEST"
                leaked = [p for p in packets if marker in bytes(p)]
                self.assertEqual(
                    len(leaked), 0,
                    "FINDING 5: Packet emitted despite missing CONFIG_SGI_IP. "
                    "Handler silently used 0.0.0.0 as outer source IP.")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_10_encap_qfi_values(self):
        """QFI boundary values: default (0→9), normal (1), maximum (63)"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            s1u_ifindex = get_ifindex(self.config.rx_iface)

            cases = [
                # (qfi_in, expected_qfi_out, ue_ip_suffix, description)
                (0, 9, 200, "default: qfi=0 should output 9"),
                (1, 1, 201, "normal: qfi=1"),
                (63, 63, 202, "maximum: qfi=63 (6-bit max)"),
            ]

            for qfi_in, expected_qfi, ip_suffix, desc in cases:
                with self.subTest(desc):
                    bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
                    attach_tc(bpf, self.config.tx_iface,
                              "gtp_encap_handler", "egress")
                    try:
                        set_config(bpf, CONFIG_SGI_IP,
                                   ip_to_host_order(self.config.sgw_ip))
                        ue_ip = f"192.168.128.{ip_suffix}"
                        populate_session(
                            bpf,
                            ue_ip_str=ue_ip,
                            enb_ip_str=self.config.enb_ip,
                            teid_ul_in=self.config.teid_ul,
                            teid_dl_out=self.config.teid_dl,
                            s1u_ifindex=s1u_ifindex,
                            mac_src=get_mac_address(self.config.rx_iface),
                            mac_dst=get_mac_address(self.config.rx_iface),
                            qfi=qfi_in,
                            active=True,
                        )

                        cap = PacketCapture(self.config.rx_iface,
                                            filter="udp port 2152")
                        cap.start()
                        time.sleep(0.1)

                        pkt = (
                            Ether(dst="ff:ff:ff:ff:ff:ff") /
                            IP(src=self.config.external_ip, dst=ue_ip) /
                            UDP(sport=80, dport=12345) /
                            Raw(load=f"QFI_{qfi_in}_TEST".encode() +
                                self._PAD)
                        )
                        self._send_on_egress(pkt)
                        time.sleep(0.1)

                        packets = cap.stop()
                        self.assertGreater(len(packets), 0,
                                           f"No encap output for {desc}")

                        pkt_bytes = bytes(packets[0])
                        # PDU Session Container QFI is at byte 56:
                        # headers[54]=ext_len, [55]=PDU_type(0x10),
                        # [56]=QFI, [57]=next_ext(0)
                        qfi_offset = 56
                        qfi_byte = pkt_bytes[qfi_offset] & 0x3F
                        print(f"\n  {desc}: input={qfi_in}, "
                              f"output={qfi_byte}")
                        self.assertEqual(
                            qfi_byte, expected_qfi,
                            f"QFI mismatch for {desc}: "
                            f"got {qfi_byte}, expected {expected_qfi}")
                    finally:
                        detach_tc(self.config.tx_iface, "egress")

    def test_11_encap_non_ipv4_passthrough(self):
        """Non-IPv4 packets (ARP) pass through without encapsulation.

        The handler returns TC_ACT_OK for non-IPv4 (ebpf_gtp_encap.c:185-187).
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))
                self._add_session(bpf, s1u_ifindex)

                time.sleep(0.1)
                initial_total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                initial_encap = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                initial_dropped = read_stat(bpf, STATS_PKT_DROPPED)

                # ARP packet — non-IPv4, should pass through
                arp_pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    ARP(pdst=self.config.ue_ip) /
                    Raw(load=b'\x00' * 50)
                )
                self._send_on_egress(arp_pkt)
                time.sleep(0.1)

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: total={total}, encap_ok={encap_ok}, "
                      f"dropped={dropped}")

                self.assertGreater(total, initial_total,
                                   "Handler did not process ARP packet")
                self.assertEqual(encap_ok, initial_encap,
                                 "ARP should not be encapsulated")
                self.assertEqual(dropped, initial_dropped,
                                 "ARP should not increment drop counter")

            finally:
                detach_tc(self.config.tx_iface, "egress")

    def test_12_encap_teid_zero(self):
        """Session with teid_dl_out=0 is treated as inactive — packet dropped"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_encap()

            try:
                s1u_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_SGI_IP, ip_to_host_order(self.config.sgw_ip))
                self._add_session(bpf, s1u_ifindex, teid_dl=0)

                cap = PacketCapture(self.config.rx_iface, filter="udp port 2152")
                cap.start()
                time.sleep(0.1)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.external_ip, dst=self.config.ue_ip) /
                    UDP(sport=80, dport=12345) /
                    Raw(load=b"TEID_ZERO_TEST" + self._PAD)
                )
                self._send_on_egress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                inactive = read_stat(bpf, STATS_INACTIVE_SESSION)
                encap_ok = read_stat(bpf, STATS_GTP_ENCAP_SUCCESS)
                print(f"\n  Stats: total={total}, inactive={inactive}, "
                      f"encap_ok={encap_ok}")

                self.assertGreater(total, 0,
                                   "Handler did not process packet")
                self.assertGreater(inactive, 0,
                                   "teid_dl_out=0 should be treated as inactive")
                self.assertEqual(encap_ok, 0,
                                 "Packet with teid_dl_out=0 should not be encapsulated")

                leaked = [p for p in packets
                          if b"TEID_ZERO_TEST" in bytes(p)]
                self.assertEqual(len(leaked), 0,
                                 "Packet leaked through despite teid_dl_out=0")

            finally:
                detach_tc(self.config.tx_iface, "egress")


if __name__ == '__main__':
    if os.geteuid() != 0:
        print("ERROR: Must run as root (sudo)")
        sys.exit(1)
    unittest.main(verbosity=2)
