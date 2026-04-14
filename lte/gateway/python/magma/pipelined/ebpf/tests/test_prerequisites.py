#!/usr/bin/env python3
"""
System Prerequisites Validation Tests

Validates that the system has all required components for eBPF GTP testing:
- Root access
- Kernel version >= 4.18
- BCC library
- Scapy with GTP support
- Network capabilities

These tests can be run without root to check non-privileged requirements,
but full validation requires root access.

Run: python3 -m pytest test_prerequisites.py -v
"""

import os
import re
import unittest
import subprocess
import platform


class TestKernelSupport(unittest.TestCase):
    """Test kernel requirements for eBPF TC programs"""

    def test_kernel_version_minimum(self):
        """Kernel must be >= 4.18 for TC eBPF support"""
        release = platform.release()
        parts = release.split('.')

        try:
            major = int(parts[0])
            minor = int(parts[1].split('-')[0])
        except (IndexError, ValueError) as e:
            self.fail(f"Cannot parse kernel version '{release}': {e}")

        # TC eBPF requires kernel >= 4.18
        is_supported = major > 4 or (major == 4 and minor >= 18)

        self.assertTrue(
            is_supported,
            f"Kernel {release} not supported. eBPF TC requires >= 4.18"
        )

    def test_kernel_config_bpf(self):
        """Check kernel has BPF support compiled in"""
        # Check /proc/config.gz or /boot/config-*
        config_paths = [
            '/proc/config.gz',
            f'/boot/config-{platform.release()}',
        ]

        config_content = None
        for path in config_paths:
            if os.path.exists(path):
                try:
                    if path.endswith('.gz'):
                        import gzip
                        with gzip.open(path, 'rt') as f:
                            config_content = f.read()
                    else:
                        with open(path, 'r') as f:
                            config_content = f.read()
                    break
                except Exception:
                    continue

        if config_content is None:
            self.skipTest("Cannot read kernel config - assuming BPF is enabled")

        # Check for BPF configs
        required_configs = [
            'CONFIG_BPF=y',
            'CONFIG_BPF_SYSCALL=y',
        ]

        for cfg in required_configs:
            cfg_name = cfg.split('=')[0]
            # Match exactly "CONFIG_X=y" or "CONFIG_X=m" at start of line.
            # Rejects "# CONFIG_X is not set" and substring matches.
            if not re.search(rf'^{cfg_name}=[ym]', config_content,
                             re.MULTILINE):
                self.skipTest(f"Cannot verify {cfg_name} is enabled")

    def test_bpf_filesystem_mounted(self):
        """BPF filesystem should be mounted"""
        # Check if bpf fs is mounted
        try:
            result = subprocess.run(
                ['mount'],
                capture_output=True,
                text=True
            )
            if 'bpf' in result.stdout:
                return  # Found BPF mount

            # Check standard location
            if os.path.exists('/sys/fs/bpf'):
                return

            self.skipTest("BPF filesystem not detected - may need to mount")
        except Exception as e:
            self.skipTest(f"Cannot check BPF filesystem: {e}")


class TestTunnelSupport(unittest.TestCase):
    """Test kernel tunnel module requirements for bpf_skb_set_tunnel_key.

    Finding 1: The decap handler requires ip_tunnel/vxlan kernel modules
    for bpf_skb_set_tunnel_key(). Without them, all uplink GTP-U traffic
    is silently dropped.
    """

    def test_ip_tunnel_module(self):
        """ip_tunnel kernel module must be loaded"""
        if os.path.exists('/sys/module/ip_tunnel'):
            return
        try:
            with open('/proc/modules', 'r') as f:
                modules = f.read()
            if 'ip_tunnel' in modules:
                return
        except Exception:
            pass
        self.skipTest(
            "ip_tunnel module not loaded — run: sudo modprobe ip_tunnel")

    def test_vxlan_module(self):
        """vxlan kernel module must be loaded"""
        if os.path.exists('/sys/module/vxlan'):
            return
        try:
            with open('/proc/modules', 'r') as f:
                modules = f.read()
            if 'vxlan' in modules:
                return
        except Exception:
            pass
        self.skipTest(
            "vxlan module not loaded — run: sudo modprobe vxlan")

    def test_lwtunnel_config(self):
        """CONFIG_LWTUNNEL must be enabled for bpf_skb_set_tunnel_key"""
        config_content = None
        for path in [f'/boot/config-{platform.release()}',
                     '/proc/config.gz']:
            if os.path.exists(path):
                try:
                    if path.endswith('.gz'):
                        import gzip
                        with gzip.open(path, 'rt') as f:
                            config_content = f.read()
                    else:
                        with open(path, 'r') as f:
                            config_content = f.read()
                    break
                except Exception:
                    continue
        if config_content is None:
            self.skipTest("Cannot read kernel config")
        if re.search(r'^CONFIG_LWTUNNEL=[ym]', config_content,
                     re.MULTILINE):
            return
        if re.search(r'^# CONFIG_LWTUNNEL is not set',
                     config_content, re.MULTILINE):
            self.fail("CONFIG_LWTUNNEL is disabled — "
                      "bpf_skb_set_tunnel_key will fail")
        self.skipTest("Cannot verify CONFIG_LWTUNNEL")

    def test_net_ip_tunnel_config(self):
        """CONFIG_NET_IP_TUNNEL must be enabled for ip_tunnel support"""
        config_content = None
        for path in [f'/boot/config-{platform.release()}',
                     '/proc/config.gz']:
            if os.path.exists(path):
                try:
                    if path.endswith('.gz'):
                        import gzip
                        with gzip.open(path, 'rt') as f:
                            config_content = f.read()
                    else:
                        with open(path, 'r') as f:
                            config_content = f.read()
                    break
                except Exception:
                    continue
        if config_content is None:
            self.skipTest("Cannot read kernel config")
        if re.search(r'^CONFIG_NET_IP_TUNNEL=[ym]', config_content,
                     re.MULTILINE):
            return
        if re.search(r'^# CONFIG_NET_IP_TUNNEL is not set',
                     config_content, re.MULTILINE):
            self.fail("CONFIG_NET_IP_TUNNEL is disabled — "
                      "ip_tunnel module unavailable")
        self.skipTest("Cannot verify CONFIG_NET_IP_TUNNEL")


class TestBCCLibrary(unittest.TestCase):
    """Test BCC (BPF Compiler Collection) availability"""

    def test_bcc_import(self):
        """BCC library must be importable"""
        try:
            from bcc import BPF
            self.assertIsNotNone(BPF)
        except ImportError as e:
            self.fail(
                f"BCC library not available: {e}\n"
                "Install with: apt-get install python3-bpfcc"
            )

    def test_bcc_version(self):
        """Check BCC version"""
        try:
            from bcc import __version__ as bcc_version
            print(f"\n  BCC version: {bcc_version}")
            # Any version should work, but log it
        except ImportError:
            try:
                from bcc import BPF
                # BCC available but no version info
                print("\n  BCC available (version unknown)")
            except ImportError:
                self.skipTest("BCC not available")

    def test_bcc_compile_simple(self):
        """Test BCC can compile a simple program"""
        try:
            from bcc import BPF

            # Minimal eBPF program
            prog = """
            int test(void *ctx) {
                return 0;
            }
            """

            # This will compile the program
            b = BPF(text=prog)
            self.assertIsNotNone(b)
        except ImportError:
            self.skipTest("BCC not available")
        except Exception as e:
            self.fail(f"BCC compilation failed: {e}")


class TestScapySupport(unittest.TestCase):
    """Test Scapy and GTP support"""

    def test_scapy_import(self):
        """Scapy must be importable"""
        try:
            from scapy.all import Ether, IP, UDP, Raw
            self.assertIsNotNone(Ether)
            self.assertIsNotNone(IP)
        except ImportError as e:
            self.fail(
                f"Scapy not available: {e}\n"
                "Install with: pip install scapy"
            )

    def test_scapy_gtp_support(self):
        """Scapy GTP contrib module must be available"""
        try:
            from scapy.contrib.gtp import GTP_U_Header
            self.assertIsNotNone(GTP_U_Header)
        except ImportError as e:
            self.fail(
                f"Scapy GTP support not available: {e}\n"
                "GTP support should be built into Scapy"
            )

    def test_scapy_build_gtp_packet(self):
        """Test building a GTP-U packet with Scapy"""
        try:
            from scapy.all import Ether, IP, UDP, Raw
            from scapy.contrib.gtp import GTP_U_Header

            # Build a complete GTP-U packet
            pkt = (
                Ether() /
                IP(src="10.0.2.1", dst="10.0.2.2") /
                UDP(sport=2152, dport=2152) /
                GTP_U_Header(teid=0x12345678) /
                IP(src="192.168.1.1", dst="8.8.8.8") /
                UDP(sport=1234, dport=80) /
                Raw(b"test payload")
            )

            # Convert to bytes
            pkt_bytes = bytes(pkt)

            # Should be a valid packet
            self.assertGreater(len(pkt_bytes), 0)

            # Parse it back
            parsed = Ether(pkt_bytes)
            self.assertIn(GTP_U_Header, parsed)
            self.assertEqual(parsed[GTP_U_Header].teid, 0x12345678)

        except ImportError:
            self.skipTest("Scapy or GTP not available")


class TestNetworkCapabilities(unittest.TestCase):
    """Test network-related capabilities"""

    def test_ip_command_available(self):
        """ip command must be available"""
        try:
            result = subprocess.run(
                ['ip', '-V'],
                capture_output=True,
                text=True
            )
            self.assertEqual(result.returncode, 0)
            print(f"\n  {result.stdout.strip()}")
        except FileNotFoundError:
            self.fail("'ip' command not found - install iproute2")

    def test_tc_command_available(self):
        """tc command must be available"""
        try:
            result = subprocess.run(
                ['tc', '-V'],
                capture_output=True,
                text=True
            )
            # tc -V may return non-zero but still output version
            self.assertIn('tc', result.stdout.lower() + result.stderr.lower())
        except FileNotFoundError:
            self.fail("'tc' command not found - install iproute2")

    def test_ethtool_available(self):
        """ethtool should be available (optional but recommended)"""
        try:
            result = subprocess.run(
                ['ethtool', '--version'],
                capture_output=True,
                text=True
            )
            if result.returncode != 0:
                self.skipTest("ethtool not working properly")
        except FileNotFoundError:
            self.skipTest("ethtool not found - checksum offload tests may fail")


@unittest.skipUnless(os.geteuid() == 0, "Root access required")
class TestRootCapabilities(unittest.TestCase):
    """Tests that require root access"""

    def test_can_create_veth(self):
        """Test we can create veth pairs"""
        import random
        suffix = random.randint(1000, 9999)
        name_a = f"prereq_a_{suffix}"
        name_b = f"prereq_b_{suffix}"

        try:
            # Create veth pair
            result = subprocess.run(
                ['ip', 'link', 'add', name_a, 'type', 'veth', 'peer', 'name', name_b],
                capture_output=True,
                text=True
            )
            self.assertEqual(result.returncode, 0, f"Failed to create veth: {result.stderr}")

            # Verify interfaces exist
            self.assertTrue(os.path.exists(f'/sys/class/net/{name_a}'))
            self.assertTrue(os.path.exists(f'/sys/class/net/{name_b}'))

        finally:
            # Cleanup
            subprocess.run(['ip', 'link', 'del', name_a], capture_output=True)

    def test_can_add_tc_qdisc(self):
        """Test we can add TC qdisc for eBPF attachment"""
        import random
        suffix = random.randint(1000, 9999)
        name_a = f"tc_test_a_{suffix}"
        name_b = f"tc_test_b_{suffix}"

        try:
            # Create veth pair
            subprocess.run(
                ['ip', 'link', 'add', name_a, 'type', 'veth', 'peer', 'name', name_b],
                capture_output=True
            )

            # Bring up interface
            subprocess.run(['ip', 'link', 'set', name_a, 'up'], capture_output=True)

            # Add clsact qdisc
            result = subprocess.run(
                ['tc', 'qdisc', 'add', 'dev', name_a, 'clsact'],
                capture_output=True,
                text=True
            )
            self.assertEqual(result.returncode, 0, f"Failed to add clsact: {result.stderr}")

            # Verify qdisc exists
            result = subprocess.run(
                ['tc', 'qdisc', 'show', 'dev', name_a],
                capture_output=True,
                text=True
            )
            self.assertIn('clsact', result.stdout)

        finally:
            # Cleanup
            subprocess.run(['ip', 'link', 'del', name_a], capture_output=True)

    def test_can_load_bpf_program(self):
        """Test we can load an eBPF program"""
        try:
            from bcc import BPF

            # Minimal TC program
            prog = """
            #include <linux/bpf.h>
            #include <linux/pkt_cls.h>

            int test_prog(struct __sk_buff *skb) {
                return TC_ACT_OK;
            }
            """

            b = BPF(text=prog)
            fn = b.load_func("test_prog", BPF.SCHED_CLS)
            self.assertIsNotNone(fn)

        except ImportError:
            self.skipTest("BCC not available")
        except Exception as e:
            self.fail(f"Failed to load eBPF program: {e}")


class TestLibraryAvailability(unittest.TestCase):
    """Test that our test library modules are available"""

    def test_lib_config(self):
        """lib.config module must be importable"""
        try:
            from lib.config import GTPUTestConfig, populate_session
            self.assertIsNotNone(GTPUTestConfig)
            self.assertIsNotNone(populate_session)

            # Test default config
            cfg = GTPUTestConfig()
            self.assertIsNotNone(cfg.tx_iface)
            self.assertIsNotNone(cfg.ue_ip)
        except ImportError as e:
            self.fail(f"lib.config not available: {e}")

    def test_lib_gtpu(self):
        """lib.gtpu module must be importable"""
        try:
            from lib.gtpu import (
                GTPUConstants,
                build_gtpu_header,
                parse_gtpu_header,
                build_gtpu_packet,
            )
            self.assertIsNotNone(GTPUConstants)
            self.assertEqual(GTPUConstants.PORT, 2152)
        except ImportError as e:
            self.fail(f"lib.gtpu not available: {e}")

    def test_lib_network(self):
        """lib.network module must be importable"""
        try:
            from lib.network import (
                veth_pair,
                setup_tc_qdisc,
                get_ifindex,
                check_root,
            )
            self.assertIsNotNone(veth_pair)
            self.assertIsNotNone(setup_tc_qdisc)
        except ImportError as e:
            self.fail(f"lib.network not available: {e}")

    def test_lib_capture(self):
        """lib.capture module must be importable"""
        try:
            from lib.capture import (
                PacketCapture,
                capture_packets,
                packet_contains_payload,
            )
            self.assertIsNotNone(PacketCapture)
        except ImportError as e:
            self.fail(f"lib.capture not available: {e}")

    def test_lib_config_production_paths(self):
        """lib.config must know production eBPF source paths"""
        try:
            from lib.config import (
                DECAP_SOURCE,
                ENCAP_SOURCE,
                PRODUCTION_CFLAGS,
                populate_session,
                set_config,
                read_stat,
            )
            self.assertIsNotNone(DECAP_SOURCE)
            self.assertIsNotNone(ENCAP_SOURCE)
            self.assertIsInstance(PRODUCTION_CFLAGS, list)
        except ImportError as e:
            self.fail(f"lib.config production helpers not available: {e}")


def print_system_info():
    """Print system information for debugging"""
    print("\n" + "=" * 60)
    print("SYSTEM INFORMATION")
    print("=" * 60)

    # OS Info
    print(f"\nOS: {platform.system()} {platform.release()}")
    print(f"Architecture: {platform.machine()}")
    print(f"Python: {platform.python_version()}")

    # Root status
    print(f"Running as root: {os.geteuid() == 0}")

    # BCC version
    try:
        from bcc import __version__ as bcc_version
        print(f"BCC version: {bcc_version}")
    except ImportError:
        print("BCC: NOT INSTALLED")

    # Scapy version
    try:
        import scapy
        print(f"Scapy version: {scapy.__version__}")
    except ImportError:
        print("Scapy: NOT INSTALLED")

    print("=" * 60 + "\n")


if __name__ == '__main__':
    print_system_info()
    unittest.main(verbosity=2)
