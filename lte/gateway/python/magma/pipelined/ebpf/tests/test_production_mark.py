#!/usr/bin/env python3
"""
Production eBPF GTP Mark Restoration Tests

Tests the ACTUAL ebpf_gtp_veth0_mark.c program -- compiles it with BCC,
loads it, attaches to a veth pair, sends decapped inner IP packets, and
verifies mark restoration via session map and stats counters.

Topology:
    [send inner IP] --> peer_veth --arrival--> mark_veth --[gtp_veth0_mark_handler ingress]--> kernel

The mark handler looks up the UE session by source IP, computes
compute_ue_mark(ue_ip), and writes it to session_info->metadata_mark
and skb->mark. Returns TC_ACT_OK (passthrough) for all packets.

Requirements:
- Root access (sudo)
- BCC library (python3-bpfcc)
- Scapy
- Linux kernel >= 4.18
- Production ebpf_gtp_veth0_mark.c in parent directory

Run: sudo python3 -m pytest test_production_mark.py -v -s
"""

import os
import sys
import time
import unittest

if os.geteuid() != 0:
    raise unittest.SkipTest("Root access required")

try:
    from bcc import BPF
    BCC_AVAILABLE = True
except ImportError:
    BCC_AVAILABLE = False

try:
    from scapy.all import Ether, IP, UDP, ARP, Raw, sendp, conf
    conf.verb = 0
    SCAPY_AVAILABLE = True
except ImportError:
    SCAPY_AVAILABLE = False

try:
    from lib.config import (
        GTPUTestConfig, MARK_SOURCE, PRODUCTION_CFLAGS,
        load_production_source, populate_session, read_stat,
        ip_to_host_order, compute_ue_mark, attach_tc, detach_tc, has_map,
        STATS_VETH0_PACKETS_PROCESSED, STATS_VETH0_MARK_RESTORED,
        STATS_VETH0_SESSION_MISS,
    )
    from lib.network import veth_pair, setup_tc_qdisc
    LIB_AVAILABLE = True
except ImportError as e:
    LIB_AVAILABLE = False
    IMPORT_ERROR = str(e)


@unittest.skipUnless(LIB_AVAILABLE, "Test library not available")
@unittest.skipUnless(BCC_AVAILABLE, "BCC not available")
class TestMarkCompilation(unittest.TestCase):
    """Test that the production mark program compiles and has expected maps"""

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(MARK_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {MARK_SOURCE}")
        cls.source = load_production_source(MARK_SOURCE)

    def test_01_compiles(self):
        """ebpf_gtp_veth0_mark.c compiles with production cflags"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        self.assertIsNotNone(bpf)

    def test_02_has_mark_handler(self):
        """gtp_veth0_mark_handler function is loadable"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        fn = bpf.load_func("gtp_veth0_mark_handler", BPF.SCHED_CLS)
        self.assertIsNotNone(fn)

    def test_03_has_session_map(self):
        """ue_session_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_veth0_mark_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "ue_session_map"))

    def test_04_has_stats_map(self):
        """stats_map exists"""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        bpf.load_func("gtp_veth0_mark_handler", BPF.SCHED_CLS)
        self.assertTrue(has_map(bpf, "stats_map"))


@unittest.skipUnless(LIB_AVAILABLE, "Test library not available")
@unittest.skipUnless(BCC_AVAILABLE, "BCC not available")
@unittest.skipUnless(SCAPY_AVAILABLE, "Scapy not available")
class TestMarkEndToEnd(unittest.TestCase):
    """
    End-to-end tests of the production GTP mark restoration program.

    Each test:
    1. Creates veth pair (mark_veth <-> peer_veth)
    2. Compiles and attaches gtp_veth0_mark_handler on mark_veth ingress
    3. Populates session map
    4. Sends inner IP packets from peer_veth (arrives on mark_veth ingress)
    5. Verifies mark restoration via session map and stats counters
    """

    _test_counter = 0

    @classmethod
    def setUpClass(cls):
        if not os.path.exists(MARK_SOURCE):
            raise unittest.SkipTest(f"Production source not found: {MARK_SOURCE}")
        cls.source = load_production_source(MARK_SOURCE)

    def setUp(self):
        TestMarkEndToEnd._test_counter += 1
        n = TestMarkEndToEnd._test_counter
        self.config = GTPUTestConfig(
            tx_iface=f"mrk_v0_{n}",
            rx_iface=f"mrk_pr_{n}",
        )

    # Verifier BPF: reads skb->mark after the mark handler and stores
    # it in a map keyed by source IP, proving the actual packet mark.
    VERIFIER_SOURCE = """
    #include <linux/bpf.h>
    #include <linux/pkt_cls.h>
    #include <linux/if_ether.h>
    #include <linux/ip.h>

    BPF_HASH(observed_marks, __u32, __u32, 64);

    int mark_verifier(struct __sk_buff *skb) {
        __u8 hdr[34];
        if (bpf_skb_load_bytes(skb, 0, hdr, sizeof(hdr)) < 0)
            return TC_ACT_OK;
        __u16 eth_type = (__u16)hdr[12] << 8 | (__u16)hdr[13];
        if (eth_type != 0x0800)
            return TC_ACT_OK;
        __u32 src_ip = (__u32)hdr[26] << 24 | (__u32)hdr[27] << 16 |
                       (__u32)hdr[28] << 8 | (__u32)hdr[29];
        __u32 mark = skb->mark;
        observed_marks.update(&src_ip, &mark);
        return TC_ACT_OK;
    }
    """

    def _setup_mark(self):
        """Compile, load, attach mark program on tx_iface ingress."""
        bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        attach_tc(bpf, self.config.tx_iface,
                  "gtp_veth0_mark_handler", "ingress")
        return bpf

    def _setup_mark_with_verifier(self):
        """Attach mark handler (direct_action=False) + verifier probe.

        Production uses direct_action=False so packets continue after
        TC_ACT_OK. The verifier runs as a second filter on the same
        ingress hook and records skb->mark in observed_marks map.
        """
        from pyroute2 import IPRoute

        mark_bpf = BPF(text=self.source, cflags=PRODUCTION_CFLAGS)
        verifier_bpf = BPF(text=self.VERIFIER_SOURCE)

        ipr = IPRoute()
        try:
            ifindex = ipr.link_lookup(ifname=self.config.tx_iface)[0]
            parent = "ffff:fff2"  # ingress

            # Filter 1 (prio 1): mark handler, direct_action=False
            fn_mark = mark_bpf.load_func(
                "gtp_veth0_mark_handler", BPF.SCHED_CLS)
            ipr.tc("add-filter", "bpf", ifindex, ":1",
                   fd=fn_mark.fd, name=fn_mark.name, parent=parent,
                   classid=1, direct_action=False, prio=1)

            # Filter 2 (prio 2): verifier, direct_action=True
            fn_verify = verifier_bpf.load_func(
                "mark_verifier", BPF.SCHED_CLS)
            ipr.tc("add-filter", "bpf", ifindex, ":1",
                   fd=fn_verify.fd, name=fn_verify.name, parent=parent,
                   classid=1, direct_action=True, prio=2)
        finally:
            ipr.close()

        return mark_bpf, verifier_bpf

    def _read_observed_mark(self, verifier_bpf, ue_ip):
        """Read observed skb->mark from verifier map."""
        import ctypes
        marks = verifier_bpf.get_table("observed_marks")
        key = ctypes.c_uint32(ip_to_host_order(ue_ip))
        try:
            val = marks[key]
            return val.value
        except KeyError:
            return None

    def _add_session(self, bpf, active=True, ue_ip=None):
        """Add a UE session for mark restoration."""
        populate_session(
            bpf,
            ue_ip_str=ue_ip or self.config.ue_ip,
            enb_ip_str=self.config.enb_ip,
            teid_ul_in=self.config.teid_ul,
            teid_dl_out=self.config.teid_dl,
            active=active,
        )

    def _send_to_ingress(self, pkt):
        """Send packet so it arrives on tx_iface ingress.

        Send from rx_iface (peer) so packet arrives at tx_iface ingress.
        """
        sendp(pkt, iface=self.config.rx_iface, verbose=False)

    def _read_session_mark(self, bpf, ue_ip=None):
        """Read metadata_mark from session map for a UE."""
        session_map = bpf.get_table("ue_session_map")
        key = session_map.Key()
        key.ue_ip = ip_to_host_order(ue_ip or self.config.ue_ip)
        try:
            val = session_map[key]
            return val.metadata_mark
        except KeyError:
            return None

    def test_01_attach_to_tc(self):
        """Production mark handler attaches to TC ingress"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()
            try:
                import subprocess
                result = subprocess.run(
                    ['tc', 'filter', 'show', 'dev',
                     self.config.tx_iface, 'ingress'],
                    capture_output=True, text=True
                )
                self.assertIn('bpf', result.stdout.lower())
                print(f"\n  Attached gtp_veth0_mark_handler to "
                      f"{self.config.tx_iface} ingress")
            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_02_mark_restoration(self):
        """Active session gets correct metadata_mark from compute_ue_mark"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                self._add_session(bpf)

                # Verify mark starts at 0
                initial_mark = self._read_session_mark(bpf)
                self.assertEqual(initial_mark, 0,
                                 "metadata_mark should start at 0")

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"MARK_RESTORE_TEST")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                processed = read_stat(bpf, STATS_VETH0_PACKETS_PROCESSED)
                restored = read_stat(bpf, STATS_VETH0_MARK_RESTORED)
                print(f"\n  Stats: processed={processed}, "
                      f"restored={restored}")

                self.assertGreater(processed, 0,
                                   "Handler did not process packet")
                self.assertGreater(restored, 0,
                                   "Mark not restored")

                # Verify mark matches compute_ue_mark
                ue_ip_int = ip_to_host_order(self.config.ue_ip)
                expected_mark = compute_ue_mark(ue_ip_int)
                actual_mark = self._read_session_mark(bpf)
                print(f"  Expected mark: 0x{expected_mark:08x}")
                print(f"  Actual mark:   0x{actual_mark:08x}")

                self.assertEqual(
                    actual_mark, expected_mark,
                    f"metadata_mark 0x{actual_mark:08x} != "
                    f"compute_ue_mark 0x{expected_mark:08x}")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_03_mark_session_miss(self):
        """Unknown UE triggers session miss, packet passes through"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                # No session added — source IP won't match
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src="10.99.99.99", dst="8.8.8.8") /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"UNKNOWN_UE_MARK_TEST")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                processed = read_stat(bpf, STATS_VETH0_PACKETS_PROCESSED)
                miss = read_stat(bpf, STATS_VETH0_SESSION_MISS)
                restored = read_stat(bpf, STATS_VETH0_MARK_RESTORED)
                print(f"\n  Stats: processed={processed}, "
                      f"miss={miss}, restored={restored}")

                self.assertGreater(processed, 0,
                                   "Handler did not process packet")
                self.assertGreater(miss, 0,
                                   "Session miss not counted")
                self.assertEqual(restored, 0,
                                 "Mark should not be restored for "
                                 "unknown UE")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_04_mark_inactive_session(self):
        """Inactive session does not get mark restored"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                self._add_session(bpf, active=False)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"INACTIVE_MARK_TEST_PAD")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                processed = read_stat(bpf, STATS_VETH0_PACKETS_PROCESSED)
                restored = read_stat(bpf, STATS_VETH0_MARK_RESTORED)
                mark = self._read_session_mark(bpf)
                print(f"\n  Stats: processed={processed}, "
                      f"restored={restored}, mark=0x{mark:08x}")

                self.assertGreater(processed, 0,
                                   "Handler did not process packet")
                self.assertEqual(restored, 0,
                                 "Mark should not be restored for "
                                 "inactive session")
                self.assertEqual(mark, 0,
                                 "metadata_mark should remain 0 for "
                                 "inactive session")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_05_mark_non_ipv4_passthrough(self):
        """Non-IPv4 packets pass through without mark processing"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                self._add_session(bpf)

                time.sleep(0.1)
                initial_processed = read_stat(
                    bpf, STATS_VETH0_PACKETS_PROCESSED)
                initial_restored = read_stat(
                    bpf, STATS_VETH0_MARK_RESTORED)
                initial_miss = read_stat(
                    bpf, STATS_VETH0_SESSION_MISS)

                # ARP packet — non-IPv4
                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    ARP(pdst=self.config.ue_ip)
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                processed = read_stat(bpf, STATS_VETH0_PACKETS_PROCESSED)
                restored = read_stat(bpf, STATS_VETH0_MARK_RESTORED)
                miss = read_stat(bpf, STATS_VETH0_SESSION_MISS)
                print(f"\n  Stats: processed={processed}, "
                      f"restored={restored}, miss={miss}")

                # Handler processes the packet (increments counter)
                self.assertGreater(processed, initial_processed,
                                   "Handler did not process ARP packet")
                # But does not attempt mark restore or session lookup
                self.assertEqual(restored, initial_restored,
                                 "ARP should not trigger mark restore")
                self.assertEqual(miss, initial_miss,
                                 "ARP should not trigger session miss")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_06_mark_multiple_ues(self):
        """Different UEs get correct distinct marks"""
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                ues = [
                    "192.168.128.101",
                    "192.168.128.102",
                    "192.168.128.103",
                ]
                for ue_ip in ues:
                    self._add_session(bpf, ue_ip=ue_ip)

                for ue_ip in ues:
                    pkt = (
                        Ether(dst="ff:ff:ff:ff:ff:ff") /
                        IP(src=ue_ip, dst="8.8.8.8") /
                        UDP(sport=12345, dport=80) /
                        Raw(load=f"MULTI_UE_{ue_ip}".encode())
                    )
                    self._send_to_ingress(pkt)
                    time.sleep(0.05)

                time.sleep(0.1)

                restored = read_stat(bpf, STATS_VETH0_MARK_RESTORED)
                self.assertGreaterEqual(restored, len(ues),
                                        f"Expected >= {len(ues)} marks "
                                        f"restored, got {restored}")

                for ue_ip in ues:
                    ue_ip_int = ip_to_host_order(ue_ip)
                    expected = compute_ue_mark(ue_ip_int)
                    actual = self._read_session_mark(bpf, ue_ip=ue_ip)
                    print(f"  UE {ue_ip}: expected=0x{expected:08x} "
                          f"actual=0x{actual:08x}")
                    self.assertEqual(
                        actual, expected,
                        f"UE {ue_ip}: mark 0x{actual:08x} != "
                        f"expected 0x{expected:08x}")

                # Verify marks are distinct
                marks = set()
                for ue_ip in ues:
                    marks.add(self._read_session_mark(bpf, ue_ip=ue_ip))
                self.assertEqual(len(marks), len(ues),
                                 "Marks should be distinct per UE")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_07_mark_computation_correctness(self):
        """Mark values for boundary IPs match compute_ue_mark() exactly.

        Note: Python compute_ue_mark() is a port of the C function, not
        an independent specification. This test catches byte-order, wiring,
        and transcription bugs, but not a shared algorithmic error.
        The mark formula is defined solely in the eBPF C code.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            bpf = self._setup_mark()

            try:
                # Test IPs that exercise compute_ue_mark edge cases
                test_ips = [
                    "192.168.128.1",   # Normal range
                    "10.0.0.1",        # Low range
                    "172.16.0.1",      # Mid range
                    "1.0.0.1",         # mark < 0x10000000 (OR 0x12000000)
                    "0.0.0.1",         # Near-zero (mask → 0, fallback)
                ]

                for ue_ip in test_ips:
                    self._add_session(bpf, ue_ip=ue_ip)

                for ue_ip in test_ips:
                    pkt = (
                        Ether(dst="ff:ff:ff:ff:ff:ff") /
                        IP(src=ue_ip, dst="8.8.8.8") /
                        UDP(sport=12345, dport=80) /
                        Raw(load=f"BOUNDARY_{ue_ip}".encode())
                    )
                    self._send_to_ingress(pkt)
                    time.sleep(0.05)

                time.sleep(0.1)

                for ue_ip in test_ips:
                    ue_ip_int = ip_to_host_order(ue_ip)
                    expected = compute_ue_mark(ue_ip_int)
                    actual = self._read_session_mark(bpf, ue_ip=ue_ip)
                    print(f"  {ue_ip} (0x{ue_ip_int:08x}): "
                          f"eBPF=0x{actual:08x} "
                          f"Python=0x{expected:08x}")
                    self.assertEqual(
                        actual, expected,
                        f"{ue_ip}: eBPF mark 0x{actual:08x} != "
                        f"Python compute_ue_mark 0x{expected:08x}")

            finally:
                detach_tc(self.config.tx_iface, "ingress")

    def test_08_skb_mark_verified(self):
        """Verify skb->mark is actually set on the packet, not just the map.

        Chains two TC ingress filters:
        1. gtp_veth0_mark_handler (direct_action=False, prio 1)
        2. mark_verifier probe (direct_action=True, prio 2)

        The verifier reads skb->mark and writes it to observed_marks map
        keyed by source IP. This proves the packet itself carries the
        correct mark, not just session_info->metadata_mark.
        """
        with veth_pair(self.config.tx_iface, self.config.rx_iface):
            setup_tc_qdisc(self.config.tx_iface)
            mark_bpf, verifier_bpf = self._setup_mark_with_verifier()

            try:
                self._add_session(mark_bpf)

                pkt = (
                    Ether(dst="ff:ff:ff:ff:ff:ff") /
                    IP(src=self.config.ue_ip,
                       dst=self.config.external_ip) /
                    UDP(sport=12345, dport=80) /
                    Raw(load=b"SKB_MARK_VERIFY_TEST")
                )
                self._send_to_ingress(pkt)
                time.sleep(0.1)

                restored = read_stat(mark_bpf,
                                     STATS_VETH0_MARK_RESTORED)
                expected_mark = compute_ue_mark(
                    ip_to_host_order(self.config.ue_ip))
                observed = self._read_observed_mark(
                    verifier_bpf, self.config.ue_ip)

                print(f"\n  Stats: restored={restored}")
                print(f"  Expected mark: 0x{expected_mark:08x}")
                print(f"  Observed skb->mark: "
                      f"{'0x{:08x}'.format(observed) if observed is not None else 'None'}")

                self.assertGreater(restored, 0,
                                   "Mark handler did not restore mark")
                self.assertIsNotNone(
                    observed,
                    "Verifier did not capture skb->mark — packet "
                    "may not have reached second filter")
                self.assertEqual(
                    observed, expected_mark,
                    f"skb->mark 0x{observed:08x} != expected "
                    f"0x{expected_mark:08x}. Handler may update "
                    f"metadata_mark but fail to set skb->mark.")

            finally:
                detach_tc(self.config.tx_iface, "ingress")


if __name__ == '__main__':
    if os.geteuid() != 0:
        print("ERROR: Must run as root (sudo)")
        sys.exit(1)
    unittest.main(verbosity=2)
