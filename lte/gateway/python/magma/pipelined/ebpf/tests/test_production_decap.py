#!/usr/bin/env python3
"""
Production eBPF GTP Decapsulation Tests

Tests the ACTUAL ebpf_gtp_decap.c program — compiles it with BCC,
loads it, attaches to a veth pair, sends real GTP-U packets, and
verifies decapsulated output.

Topology:
    [send GTP-U] --> s1u_veth --[gtp_decap_handler ingress]--> ovs_veth --> [capture inner IP]

The decap handler strips outer IP/UDP/GTP headers, reconstructs the
Ethernet header using MACs from the session map, and redirects the
inner IP packet to the OVS-side veth via bpf_redirect().

Requirements:
- Root access (sudo)
- BCC library (python3-bpfcc)
- Scapy with GTP support
- Linux kernel >= 4.18
- Production ebpf_gtp_decap.c in parent directory

Run: sudo python3 -m pytest test_production_decap.py -v -s
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
    from scapy.all import Ether, IP, UDP, TCP, ICMP, Raw, sendp, conf
    from scapy.contrib.gtp import GTP_U_Header
    conf.verb = 0
    SCAPY_AVAILABLE = True
except ImportError:
    SCAPY_AVAILABLE = False

try:
    from lib.config import (
        GTPUTestConfig, DECAP_SOURCE, PRODUCTION_CFLAGS,
        load_production_source, populate_session, set_config, read_stat,
        ip_to_host_order, compute_ue_mark, golden_decap,
        attach_tc, detach_tc, has_map,
        STATS_TOTAL_PROCESSED, STATS_GTP_DECAP_SUCCESS, STATS_SESSION_MISS,
        STATS_TEID_MISMATCH, STATS_PKT_FORWARDED, STATS_PKT_DROPPED,
        STATS_PKT_TOO_SHORT, STATS_INACTIVE_SESSION, STATS_UL_ERRORS,
        CONFIG_OVS_IFINDEX,
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
class TestDecapCompilation(unittest.TestCase):
    """Test that the production decap program compiles and has expected maps"""

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(DECAP_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {DECAP_SOURCE}")
        cls.source = load_production_source(DECAP_SOURCE)

    def test_01_compiles(self):
        """ebpf_gtp_decap.c compiles with production cflags"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        self.assertIsNotNone(bpf)

    def test_02_has_decap_handler(self):
        """gtp_decap_handler function exists and is loadable"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        fn = bpf.load_func("gtp_decap_handler", BPF.SCHED_CLS)
        self.assertIsNotNone(fn)

    def test_03_has_session_map(self):
        """ue_session_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_decap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "ue_session_map"))

    def test_04_has_config_map(self):
        """config_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_decap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "config_map"))

    def test_05_has_stats_map(self):
        """stats_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_decap_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "stats_map"))


@unittest.skipUnless(LIB_AVAILABLE, "Test library not available")
@unittest.skipUnless(BCC_AVAILABLE, "BCC not available")
@unittest.skipUnless(SCAPY_AVAILABLE, "Scapy not available")
class TestDecapEndToEnd(unittest.TestCase):
    """
    End-to-end tests of the production GTP decapsulation program.

    Each test:
    1. Creates veth pair (s1u_veth <-> ovs_veth)
    2. Compiles and attaches gtp_decap_handler on s1u_veth ingress
    3. Configures session and config maps
    4. Sends packets on s1u_veth
    5. Captures on ovs_veth
    6. Validates results
    """

    _test_counter = 0

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(DECAP_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {DECAP_SOURCE}")
        cls.source = load_production_source(DECAP_SOURCE)

    def setUp(self):
        """Use unique veth names per test to avoid kernel state conflicts."""
        TestDecapEndToEnd._test_counter += 1
        n = TestDecapEndToEnd._test_counter
        self.config = GTPUTestConfig(
            tx_iface=f"dec_s1u{n}",
            rx_iface=f"dec_ovs{n}",
        )

    def _setup_decap(self):
        """Compile, load, attach decap program. Returns BPF object."""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        attach_tc(bpf, self.config.tx_iface, "gtp_decap_handler", "ingress")
        return bpf

    def _add_session(self, bpf, ovs_ifindex, mac_src=None, mac_dst=None,
                     active=True, teid_ul=None):
        """Add a UE session pointing redirect to ovs_veth."""
        if mac_src is None:
            mac_src = get_mac_address(self.config.rx_iface)
        if mac_dst is None:
            mac_dst = get_mac_address(self.config.rx_iface)

        populate_session(
            bpf,
            ue_ip_str=self.config.ue_ip,
            enb_ip_str=self.config.enb_ip,
            teid_ul_in=teid_ul or self.config.teid_ul,
            teid_dl_out=self.config.teid_dl,
            ovs_ifindex=ovs_ifindex,
            mac_src=mac_src,
            mac_dst=mac_dst,
            active=active,
        )

    def _build_gtp_packet(self, payload, teid=None, inner_src=None, inner_dst=None):
        """Build a GTP-U encapsulated packet using Scapy."""
        pkt = (
            Ether(dst="ff:ff:ff:ff:ff:ff") /
            IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
            UDP(sport=2152, dport=2152) /
            GTP_U_Header(teid=teid or self.config.teid_ul) /
            IP(src=inner_src or self.config.ue_ip,
               dst=inner_dst or self.config.external_ip) /
            UDP(sport=12345, dport=80) /
            Raw(load=payload)
        )
        return pkt

    def _send_to_ingress(self, pkt):
        """Send packet so it arrives on tx_iface INGRESS.

        TC ingress fires on packets arriving at the interface.
        On a veth pair a<->b, sending from b makes the packet arrive
        on a's ingress. So we send from rx_iface (ovs side) to hit
        the ingress handler on tx_iface (s1u side).
        """
        sendp(pkt, iface=self.config.rx_iface, verbose=False)

    def test_01_attach_to_tc(self):
        """Production decap attaches to TC ingress"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()
            try:
                import subprocess
                result = subprocess.run(
                    ['tc', 'filter', 'show', 'dev', self.config.tx_iface, 'ingress'],
                    capture_output=True, text=True
                )
                self.assertIn('bpf', result.stdout.lower())
                print(f"\n  Attached gtp_decap_handler to {self.config.tx_iface} ingress")
            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_02_decap_golden_model(self):
        """Golden model: validate every output field against Scapy reference"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                # Use known MACs so we can verify the reconstructed Ethernet header
                session_mac_src = "02:aa:bb:cc:dd:01"
                session_mac_dst = "02:aa:bb:cc:dd:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                payload = _generate_payload(64)

                cap = PacketCapture(self.config.rx_iface, filter="ip")
                cap.start()
                time.sleep(0.1)

                pkt = self._build_gtp_packet(payload)
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                # --- Stats validation ---
                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                forwarded = read_stat(bpf, STATS_PKT_FORWARDED)
                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                ul_errors = read_stat(bpf, STATS_UL_ERRORS)
                print(f"\n  Stats: total={total}, decap_ok={decap_ok}, forwarded={forwarded}, ul_errors={ul_errors}")

                # If ul_errors > 0 and decap_ok == 0, bpf_skb_set_tunnel_key()
                # failed at line 392 of ebpf_gtp_decap.c. This kernel helper
                # requires CONFIG_LWTUNNEL / CONFIG_NET_IP_TUNNEL. Try:
                #   sudo modprobe ip_tunnel && sudo modprobe vxlan
                if decap_ok == 0 and ul_errors > 0:
                    self.fail(
                        f"bpf_skb_set_tunnel_key() failed (ul_errors={ul_errors}). "
                        f"Ensure kernel has lightweight tunnel support: "
                        f"modprobe ip_tunnel vxlan")

                self.assertGreater(decap_ok, 0, "Decap success counter not incremented")
                self.assertGreater(forwarded, 0, "Forwarded counter not incremented")
                self.assertGreater(len(packets), 0, "No packets captured")

                # --- Classify all captured packets ---
                # pcap on rx_iface sees TWO things:
                #   1. Our TX: the original GTP packet we sent from rx_iface
                #   2. Redirect RX: the decapsulated inner IP from bpf_redirect
                # In production, only (2) would appear on gtp_veth0 because
                # GTP arrives from eth1, not gtp_veth0. We verify both are
                # present and correctly formed.
                gtp_packets = []    # Our sent GTP (TX on capture iface)
                decap_packets = []  # Decapsulated output (RX via redirect)

                for p in packets:
                    p_bytes = bytes(p)
                    if payload not in p_bytes:
                        continue
                    parsed = Ether(p_bytes)
                    if UDP in parsed and parsed[UDP].dport == 2152:
                        gtp_packets.append(parsed)
                    else:
                        decap_packets.append(parsed)

                print(f"  Captured: {len(gtp_packets)} GTP (our TX) + "
                      f"{len(decap_packets)} decapsulated (redirect RX)")

                # Must have exactly 1 GTP (our send) and 1 decapsulated (redirect)
                self.assertEqual(len(gtp_packets), 1,
                                 f"Expected 1 GTP TX packet, got {len(gtp_packets)}")
                self.assertEqual(len(decap_packets), 1,
                                 f"Expected 1 decapsulated packet, got {len(decap_packets)} "
                                 f"— eBPF may be leaking or duplicating packets")

                decap_pkt = decap_packets[0]

                # --- Ethernet header validation ---
                # eBPF reconstructs Ethernet with MACs from session map
                self.assertEqual(
                    decap_pkt[Ether].dst, session_mac_dst,
                    f"Dst MAC mismatch: {decap_pkt[Ether].dst} != {session_mac_dst}")
                self.assertEqual(
                    decap_pkt[Ether].src, session_mac_src,
                    f"Src MAC mismatch: {decap_pkt[Ether].src} != {session_mac_src}")
                self.assertEqual(
                    decap_pkt[Ether].type, 0x0800,
                    f"EtherType not IPv4: 0x{decap_pkt[Ether].type:04x}")

                # --- GTP outer headers must be GONE ---
                self.assertNotIn(
                    GTP_U_Header, decap_pkt,
                    "GTP-U header still present after decap")
                if UDP in decap_pkt:
                    self.assertNotEqual(
                        decap_pkt[UDP].dport, 2152,
                        "UDP port 2152 still present — outer headers not stripped")

                # --- Inner IP header validation ---
                self.assertIn(IP, decap_pkt, "No IP layer in decapsulated packet")
                inner_ip = decap_pkt[IP]

                self.assertEqual(
                    inner_ip.src, self.config.ue_ip,
                    f"Inner src IP: {inner_ip.src} != {self.config.ue_ip}")
                self.assertEqual(
                    inner_ip.dst, self.config.external_ip,
                    f"Inner dst IP: {inner_ip.dst} != {self.config.external_ip}")
                self.assertEqual(
                    inner_ip.version, 4,
                    f"Inner IP version: {inner_ip.version} != 4")
                self.assertEqual(
                    inner_ip.proto, 17,
                    f"Inner IP protocol: {inner_ip.proto} != 17 (UDP)")

                # --- Inner transport validation ---
                self.assertIn(UDP, decap_pkt, "No UDP layer in decapsulated packet")
                inner_udp = decap_pkt[UDP]
                self.assertEqual(inner_udp.sport, 12345, "Inner UDP sport mismatch")
                self.assertEqual(inner_udp.dport, 80, "Inner UDP dport mismatch")

                # --- Payload validation ---
                self.assertIn(Raw, decap_pkt, "No Raw payload in decapsulated packet")
                self.assertEqual(
                    decap_pkt[Raw].load, payload,
                    "Payload bytes don't match original")

                # --- Session counter validation ---
                session_map = bpf.get_table("ue_session_map")
                key = session_map.Key()
                key.ue_ip = ip_to_host_order(self.config.ue_ip)
                session = session_map[key]

                self.assertGreater(session.ul_packets, 0,
                                   "Session ul_packets not incremented")
                self.assertGreater(session.ul_bytes, 0,
                                   "Session ul_bytes not incremented")
                self.assertGreater(session.last_seen, 0,
                                   "Session last_seen not updated")

                # --- Metadata mark validation ---
                ue_ip_int = ip_to_host_order(self.config.ue_ip)
                expected_mark = compute_ue_mark(ue_ip_int)
                self.assertEqual(
                    session.metadata_mark, expected_mark,
                    f"metadata_mark 0x{session.metadata_mark:08x} != "
                    f"expected 0x{expected_mark:08x}")

                # --- Golden model byte-exact comparison ---
                expected_bytes = golden_decap(
                    bytes(pkt), session_mac_src, session_mac_dst)
                actual_bytes = bytes(decap_pkt)
                self.assertEqual(
                    actual_bytes, expected_bytes,
                    f"Decap output does not match golden model. "
                    f"Actual {len(actual_bytes)} bytes vs "
                    f"expected {len(expected_bytes)} bytes")

                print(f"  Ethernet: dst={decap_pkt[Ether].dst} src={decap_pkt[Ether].src}")
                print(f"  Inner IP: {inner_ip.src} -> {inner_ip.dst} proto={inner_ip.proto}")
                print(f"  Payload: {len(payload)} bytes verified")
                print(f"  Golden model: byte-exact match")
                print(f"  Session: ul_pkts={session.ul_packets} ul_bytes={session.ul_bytes}")
                print(f"  Mark: 0x{session.metadata_mark:08x} (expected 0x{expected_mark:08x})")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_03_decap_no_session(self):
        """Packet with unknown UE IP is dropped — verify no leak"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                marker = b"NO_SESSION_MARKER"
                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                pkt = self._build_gtp_packet(marker, inner_src="172.16.99.99")
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                miss = read_stat(bpf, STATS_SESSION_MISS)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: session_miss={miss}, dropped={dropped}")

                self.assertGreater(miss, 0, "Session miss counter not incremented")
                self.assertGreater(dropped, 0, "Drop counter not incremented")

                # Verify no decapsulated packet leaked to OVS side
                # (only our TX should appear, classified as GTP)
                leaked = []
                for p in packets:
                    parsed = Ether(bytes(p))
                    if marker in bytes(p) and not (UDP in parsed and parsed[UDP].dport == 2152):
                        leaked.append(p)
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} decapsulated packets leaked — "
                                 f"handler should drop unknown sessions")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_04_decap_teid_mismatch(self):
        """Packet with wrong TEID is dropped — verify no leak"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)
                self._add_session(bpf, ovs_ifindex, teid_ul=0xAAAA0001)

                marker = b"TEID_MISMATCH_MARKER"
                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                pkt = self._build_gtp_packet(marker, teid=0xBBBB0002)
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                mismatch = read_stat(bpf, STATS_TEID_MISMATCH)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: teid_mismatch={mismatch}, dropped={dropped}")

                self.assertGreater(mismatch, 0, "TEID mismatch counter not incremented")
                self.assertGreater(dropped, 0, "Drop counter not incremented for TEID mismatch")

                leaked = []
                for p in packets:
                    parsed = Ether(bytes(p))
                    if marker in bytes(p) and not (UDP in parsed and parsed[UDP].dport == 2152):
                        leaked.append(p)
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} decapsulated packets leaked — "
                                 f"handler should drop TEID mismatches")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_05_decap_inactive_session(self):
        """Packet for inactive session is dropped — verify no leak"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)
                self._add_session(bpf, ovs_ifindex, active=False)

                marker = b"INACTIVE_MARKER"
                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                pkt = self._build_gtp_packet(marker)
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                inactive = read_stat(bpf, STATS_INACTIVE_SESSION)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: inactive_session={inactive}, dropped={dropped}")

                self.assertGreater(inactive, 0, "Inactive session counter not incremented")
                self.assertGreater(dropped, 0, "Drop counter not incremented for inactive session")

                leaked = []
                for p in packets:
                    parsed = Ether(bytes(p))
                    if marker in bytes(p) and not (UDP in parsed and parsed[UDP].dport == 2152):
                        leaked.append(p)
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} decapsulated packets leaked — "
                                 f"handler should drop inactive sessions")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_06_decap_multiple_packets(self):
        """Multiple GTP-U packets are all decapsulated — classify output"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)
                self._add_session(bpf, ovs_ifindex)

                num_packets = 5
                payloads = [_generate_payload(64) for _ in range(num_packets)]

                cap = PacketCapture(self.config.rx_iface, filter="ip")
                cap.start()
                time.sleep(0.1)

                for payload in payloads:
                    pkt = self._build_gtp_packet(payload)
                    self._send_to_ingress(pkt)
                    time.sleep(0.05)

                time.sleep(0.1)
                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)

                # Classify captured packets for each payload
                decap_count = 0
                for payload in payloads:
                    for p in packets:
                        p_bytes = bytes(p)
                        if payload not in p_bytes:
                            continue
                        parsed = Ether(p_bytes)
                        # Only count non-GTP packets (actual decap output)
                        if not (UDP in parsed and parsed[UDP].dport == 2152):
                            decap_count += 1
                            break

                print(f"\n  Sent {num_packets}, decap_ok={decap_ok}, "
                      f"decap_captured={decap_count}")

                self.assertEqual(decap_ok, num_packets,
                                 f"decap_ok={decap_ok} != {num_packets} sent")
                self.assertEqual(decap_count, num_packets,
                                 f"Only {decap_count}/{num_packets} decapsulated "
                                 f"packets captured — output missing")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_07_decap_non_gtp_passthrough(self):
        """Non-GTP UDP packets pass through — not dropped, not decapped.

        TC_ACT_OK on ingress means the packet continues into the local
        kernel stack on tx_iface. It is NOT redirected to rx_iface, so
        we cannot capture it. We verify passthrough by checking that:
        - No GTP counters fired (decap_ok, session_miss)
        - No drop counter incremented
        - The handler processed it (total incremented)
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                # Let initial ARP/NDP from veth creation settle before
                # capturing baselines (Finding 3: non-IPv4 increments
                # STATS_PKT_DROPPED even though TC_ACT_OK passes them).
                time.sleep(0.1)

                initial_total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                initial_dropped = read_stat(bpf, STATS_PKT_DROPPED)

                # Non-GTP packet (UDP port 80, not 2152)
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src="1.1.1.1", dst="2.2.2.2") /
                    UDP(sport=1234, dport=80) /
                    Raw(load=b"NOT_GTP_TRAFFIC_MARKER_PAD" * 3)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                miss = read_stat(bpf, STATS_SESSION_MISS)
                dropped = read_stat(bpf, STATS_PKT_DROPPED)
                print(f"\n  Stats: total={total}, decap_ok={decap_ok}, "
                      f"session_miss={miss}, dropped={dropped}")

                # Handler must have processed the packet
                self.assertGreater(total, initial_total,
                                   "Handler did not process the non-GTP packet")

                # Non-GTP IPv4/UDP should not trigger any GTP processing
                self.assertEqual(decap_ok, 0,
                                 "Non-GTP packet should not increment decap_ok")
                self.assertEqual(miss, 0,
                                 "Non-GTP packet should not trigger session_miss")
                self.assertEqual(dropped - initial_dropped, 0,
                                 "IPv4/UDP packet should not increment "
                                 "STATS_PKT_DROPPED")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_08_decap_small_packet(self):
        """Small packets arriving on GTP port must not be silently dropped.

        Finding 2: The handler calls bpf_skb_load_bytes(skb, 0, pkt_data, 64)
        which drops packets under 64 bytes as PKT_TOO_SHORT before any GTP
        parsing begins. A truncated or minimal GTP packet on port 2152 should
        be parsed or explicitly rejected as invalid GTP — not silently dropped
        at a generic size check.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                # Small UDP packet to GTP port (2152).
                # Ether(14) + IP(20) + UDP(8) + Raw(5) = 47 bytes < 64.
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src="10.0.0.1", dst="10.0.0.2") /
                    UDP(sport=2152, dport=2152) /
                    Raw(load=b"SMALL")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                too_short = read_stat(bpf, STATS_PKT_TOO_SHORT)
                print(f"\n  Stats: total={total}, too_short={too_short}")

                self.assertEqual(too_short, 0,
                                 f"FINDING 2: Packet silently dropped as too short. "
                                 f"bpf_skb_load_bytes(64) rejects packets < 64 bytes "
                                 f"before GTP parsing. Should handle gracefully.")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_09_decap_malformed_gtp(self):
        """Malformed GTP packets are rejected — not decapsulated, not leaked"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)
                self._add_session(bpf, ovs_ifindex)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                # Craft malformed GTP packets using raw bytes after UDP.
                # Valid GTP flags = 0x30 (V1=1, PT=1), type = 0xFF (G-PDU)
                inner_ip = bytes(
                    IP(src=self.config.ue_ip, dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"MALFORMED_GTP_TEST_PAYLOAD_PAD" * 3)
                )

                cases = [
                    # (description, gtp_flags, gtp_type)
                    ("GTP version 2 (not v1)", 0x40, 0xFF),   # V=2, PT=0
                    ("PT flag not set",        0x20, 0xFF),   # V=1, PT=0
                    ("Echo Request type",      0x30, 0x01),   # V=1, PT=1, type=Echo
                    ("End Marker type",        0x30, 0xFE),   # V=1, PT=1, type=EndMarker
                ]

                for desc, flags, msg_type in cases:
                    # Build: Ether/IP/UDP(2152)/raw_GTP/inner_IP
                    teid = self.config.teid_ul
                    gtp_len = len(inner_ip)
                    gtp_hdr = struct.pack('!BBHI', flags, msg_type, gtp_len, teid)

                    pkt = (
                        Ether(dst="ff:ff:ff:ff:ff:ff") /
                        IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                        UDP(sport=2152, dport=2152) /
                        Raw(load=gtp_hdr + inner_ip)
                    )
                    self._send_to_ingress(pkt)
                    time.sleep(0.05)

                time.sleep(0.1)
                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                total = read_stat(bpf, STATS_TOTAL_PROCESSED)
                print(f"\n  Stats: total={total}, decap_ok={decap_ok}")
                print(f"  Sent {len(cases)} malformed GTP packets")

                # Verify handler actually processed the packets
                self.assertGreaterEqual(total, len(cases),
                                        f"Handler only processed {total}/{len(cases)} "
                                        f"packets — malformed packets may not have "
                                        f"reached the TC pipeline")

                # None should be decapsulated
                self.assertEqual(decap_ok, 0,
                                 f"Malformed GTP packets should not be decapsulated "
                                 f"(decap_ok={decap_ok})")

                # No decapsulated packets should appear on output
                marker = b"MALFORMED_GTP_TEST_PAYLOAD_PAD"
                leaked = []
                for p in packets:
                    parsed = Ether(bytes(p))
                    if marker in bytes(p) and not (UDP in parsed and parsed[UDP].dport == 2152):
                        leaked.append(p)
                self.assertEqual(len(leaked), 0,
                                 f"{len(leaked)} malformed packets leaked as decapsulated")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_10_decap_gtp_with_seq_number(self):
        """GTP with S flag (sequence number) decapsulates correctly.

        S flag adds 4 optional bytes (seq, NPDU, next-ext-type) at
        ebpf_gtp_decap.c:314-319. Inner packet offset shifts by 4.
        Real eNBs commonly set this flag.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:10:01"
                session_mac_dst = "02:de:ca:00:10:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                # GTP with S flag: adds seq_num(2) + npdu(1) + next_ext(1)
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                    UDP(sport=2152, dport=2152) /
                    GTP_U_Header(teid=self.config.teid_ul, S=1,
                                 seq=0x1234) /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"SEQ_NUMBER_TEST_PAYLOAD")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "GTP with S flag should be decapsulated")

                # Find decapped packet (non-GTP with our marker)
                marker = b"SEQ_NUMBER_TEST_PAYLOAD"
                decapped = [p for p in packets
                            if marker in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected exactly 1 decapsulated packet")

                # Verify inner packet preserved via golden model
                expected = golden_decap(bytes(pkt),
                                        session_mac_src, session_mac_dst)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "Decap output does not match golden model "
                                 "for GTP with S flag")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_11_decap_gtp_with_pdu_session_container(self):
        """GTP with E flag and PDU Session Container (0x85) extension.

        Extension parsing loop at ebpf_gtp_decap.c:336-371.
        PDU Session Container is the standard 5G extension.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:11:01"
                session_mac_dst = "02:de:ca:00:11:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                # Build GTP with E flag + PDU Session Container manually
                inner_ip = bytes(
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"PDU_SESSION_CONTAINER_TEST")
                )
                teid = self.config.teid_ul
                # GTP flags: V1=1, PT=1, E=1 → 0x34
                gtp_flags = 0x34
                gtp_len = len(inner_ip) + 4 + 4  # +4 opt fields, +4 ext
                gtp_hdr = struct.pack('!BBHI', gtp_flags, 0xFF,
                                      gtp_len, teid)
                # Optional fields: seq=0, npdu=0, next_ext_type=0x85
                opt_fields = struct.pack('!HBB', 0, 0, 0x85)
                # PDU Session Container: len=1 (4 bytes), qfi=5, next=0
                pdu_ext = struct.pack('!BBH', 1, 0x05, 0)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                    UDP(sport=2152, dport=2152) /
                    Raw(load=gtp_hdr + opt_fields + pdu_ext + inner_ip)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "GTP with PDU Session Container "
                                   "should be decapsulated")

                marker = b"PDU_SESSION_CONTAINER_TEST"
                decapped = [p for p in packets
                            if marker in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected 1 decapsulated packet")

                # Strict golden comparison (raw GTP, so build expected
                # directly instead of using golden_decap which needs
                # Scapy GTP_U_Header layer)
                expected = (bytes(Ether(src=session_mac_src,
                                        dst=session_mac_dst,
                                        type=0x0800)) +
                            inner_ip)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "PDU Session Container decap output "
                                 "does not match golden model")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_12_decap_gtp_with_multiple_extensions(self):
        """GTP with chained extension headers decapsulates correctly.

        Extension loop at ebpf_gtp_decap.c:336 processes up to 5 extensions.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:12:01"
                session_mac_dst = "02:de:ca:00:12:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                inner_ip = bytes(
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"MULTI_EXTENSION_TEST_PAD_PAD")
                )
                teid = self.config.teid_ul
                # GTP flags: V1=1, PT=1, E=1 → 0x34
                gtp_flags = 0x34
                # Extension 1: PDU Session Container (0x85), 4 bytes
                ext1 = struct.pack('!BBH', 1, 0x05, 0)  # next=0 placeholder
                # Extension 2: Generic extension (type 0x40), 4 bytes
                ext2 = struct.pack('!BBH', 1, 0x00, 0)  # next=0, end chain
                # Fix chain: ext1 next_type = 0x40
                ext1 = struct.pack('!BBBB', 1, 0x05, 0, 0x40)
                # ext2 next_type = 0 (end)
                ext2 = struct.pack('!BBBB', 1, 0x00, 0, 0x00)
                opt_fields = struct.pack('!HBB', 0, 0, 0x85)  # next=0x85
                gtp_len = len(inner_ip) + 4 + len(ext1) + len(ext2)
                gtp_hdr = struct.pack('!BBHI', gtp_flags, 0xFF,
                                      gtp_len, teid)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                    UDP(sport=2152, dport=2152) /
                    Raw(load=gtp_hdr + opt_fields + ext1 + ext2 +
                        inner_ip)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "GTP with chained extensions "
                                   "should be decapsulated")

                marker = b"MULTI_EXTENSION_TEST_PAD_PAD"
                decapped = [p for p in packets
                            if marker in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected 1 decapsulated packet")

                # Strict golden comparison
                expected = (bytes(Ether(src=session_mac_src,
                                        dst=session_mac_dst,
                                        type=0x0800)) +
                            inner_ip)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "Multi-extension decap output "
                                 "does not match golden model")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_13_decap_outer_ipv4_options(self):
        """Outer IPv4 with options — decap uses IHL correctly.

        Unlike encap, the decap handler correctly computes ip_hlen from
        the IHL field (ebpf_gtp_decap.c:266-268). This test verifies
        that outer IPv4 options do not break decapsulation.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:13:01"
                session_mac_dst = "02:de:ca:00:13:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                from scapy.all import IPOption
                inner_payload = b"OUTER_OPTIONS_TEST_PAYLOAD_XX"
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip,
                       options=[IPOption(b'\x01')] * 4) /
                    UDP(sport=2152, dport=2152) /
                    GTP_U_Header(teid=self.config.teid_ul) /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=inner_payload)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "Decap with outer IPv4 options should "
                                   "succeed (handler uses IHL correctly)")

                decapped = [p for p in packets
                            if inner_payload in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected 1 decapsulated packet")

                # Strict golden comparison
                expected = golden_decap(bytes(pkt),
                                        session_mac_src, session_mac_dst)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "Outer IPv4 options decap output "
                                 "does not match golden model")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_14_decap_inner_tcp(self):
        """Inner TCP packet preserved byte-exact after decapsulation"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:14:01"
                session_mac_dst = "02:de:ca:00:14:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                inner_payload = b"TCP_INNER_TEST_PAYLOAD_DATA_XX"
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                    UDP(sport=2152, dport=2152) /
                    GTP_U_Header(teid=self.config.teid_ul) /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    TCP(sport=80, dport=12345, flags='A') /
                    Raw(load=inner_payload)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "Inner TCP packet should be decapsulated")

                decapped = [p for p in packets
                            if inner_payload in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected 1 decapsulated packet")

                # Verify inner IP bytes preserved via golden model
                expected = golden_decap(bytes(pkt),
                                        session_mac_src, session_mac_dst)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "Inner TCP packet not preserved "
                                 "byte-exact after decap")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_15_decap_inner_icmp(self):
        """Inner ICMP packet preserved after decapsulation"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_decap()

            try:
                ovs_ifindex = get_ifindex(self.config.rx_iface)
                set_config(bpf, CONFIG_OVS_IFINDEX, ovs_ifindex)

                session_mac_src = "02:de:ca:00:15:01"
                session_mac_dst = "02:de:ca:00:15:02"
                self._add_session(bpf, ovs_ifindex,
                                  mac_src=session_mac_src,
                                  mac_dst=session_mac_dst)

                cap = PacketCapture(self.config.rx_iface)
                cap.start()
                time.sleep(0.1)

                inner_payload = b"ICMP_INNER_TEST_PAYLOAD_DATA_X"
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.enb_ip, dst=self.config.sgw_ip) /
                    UDP(sport=2152, dport=2152) /
                    GTP_U_Header(teid=self.config.teid_ul) /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    ICMP(type=8, code=0, id=0x1234, seq=1) /
                    Raw(load=inner_payload)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                packets = cap.stop()

                decap_ok = read_stat(bpf, STATS_GTP_DECAP_SUCCESS)
                print(f"\n  Stats: decap_ok={decap_ok}")

                self.assertGreater(decap_ok, 0,
                                   "Inner ICMP packet should be decapsulated")

                decapped = [p for p in packets
                            if inner_payload in bytes(p) and
                            not (UDP in Ether(bytes(p)) and
                                 Ether(bytes(p))[UDP].dport == 2152)]
                self.assertEqual(len(decapped), 1,
                                 "Expected 1 decapsulated packet")

                # Verify inner IP bytes preserved via golden model
                expected = golden_decap(bytes(pkt),
                                        session_mac_src, session_mac_dst)
                actual = bytes(decapped[0])
                self.assertEqual(
                    len(actual), len(expected),
                    f"Decap output length {len(actual)} != "
                    f"golden model {len(expected)}")
                self.assertEqual(actual, expected,
                                 "Inner ICMP packet not preserved "
                                 "byte-exact after decap")

            finally:
                detach_tc(self.config.tx_iface, "ingress")


if __name__ == '__main__':
    if os.geteuid() != 0:
        print("ERROR: Must run as root (sudo)")
        sys.exit(1)
    unittest.main(verbosity=2)
