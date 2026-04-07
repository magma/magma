# eBPF GTP Test Suite

Tests for the production eBPF GTP-U encapsulation/decapsulation programs.

## What Gets Tested

The tests compile, load, and exercise the **actual production code**:

- `ebpf_gtp_decap.c` -- `gtp_decap_handler()` on TC ingress
- `ebpf_gtp_encap.c` -- `gtp_encap_handler()` on TC egress
- `EbpfGtpMap.h` -- shared structs and map definitions

Each test creates a veth pair, attaches the production eBPF program, populates
the session and config maps with the real 160-byte `ue_session_info` struct,
sends real packets, and validates the output.

## Directory Structure

```text
tests/
+-- lib/                           # Reusable test utilities
|   +-- config.py                  # Production struct helpers, golden functions
|   +-- gtpu.py                    # GTP-U packet building (Scapy)
|   +-- network.py                 # veth pairs, TC setup
|   +-- capture.py                 # Packet capture
|
+-- test_prerequisites.py          # System requirements check
+-- test_gtpu_packets.py           # GTP-U packet unit tests (no root)
+-- test_production_decap.py       # E2E decap tests against ebpf_gtp_decap.c
+-- test_production_encap.py       # E2E encap tests against ebpf_gtp_encap.c
+-- README.md
```

## Requirements

```bash
sudo apt-get install python3-bpfcc python3-venv iproute2 ethtool linux-headers-$(uname -r)
sudo modprobe ip_tunnel vxlan  # Required for bpf_skb_set_tunnel_key (Finding 1)
```

## Running

```bash
cd lte/gateway/python/magma/pipelined/ebpf/tests
python3 -m venv venv --system-site-packages
source venv/bin/activate
pip install -r requirements.txt
sudo venv/bin/python -m pytest . -v -s
```

## Test Coverage

### test_production_decap.py

| Test | What it validates |
| ---- | ----------------- |
| test_01_compiles | `ebpf_gtp_decap.c` compiles with production cflags |
| test_02_has_decap_handler | `gtp_decap_handler` function loadable |
| test_03-05_has_maps | `ue_session_map`, `config_map`, `stats_map` exist |
| test_01_attach_to_tc | Attaches to TC ingress |
| test_02_decap_golden_model | Full field-by-field validation: MACs, inner IP, payload, session counters, metadata_mark, no GTP leakage |
| test_03_decap_no_session | Dropped + no leak on output |
| test_04_decap_teid_mismatch | Dropped + no leak on output |
| test_05_decap_inactive_session | Dropped + no leak on output |
| test_06_decap_multiple_packets | N packets in, N decapsulated out, classified by GTP/non-GTP |
| test_07_decap_non_gtp_passthrough | Non-GTP passes through, not dropped, not decapped |
| test_08_decap_small_packet | **EXPECTED FAIL** -- Finding 2: packets < 64 bytes dropped |
| test_09_decap_malformed_gtp | Invalid version/PT/type rejected, no leak |
| test_10_decap_gtp_with_seq_number | GTP S flag (seq number) decapsulated correctly |
| test_11_decap_gtp_with_pdu_session_container | E flag + PDU Session Container (0x85) extension |
| test_12_decap_gtp_with_multiple_extensions | Chained extension headers parsed correctly |
| test_13_decap_outer_ipv4_options | Outer IPv4 options handled via IHL |
| test_14_decap_inner_tcp | Inner TCP preserved byte-exact |
| test_15_decap_inner_icmp | Inner ICMP preserved byte-exact |

### test_production_encap.py

| Test | What it validates |
| ---- | ----------------- |
| test_01_compiles | `ebpf_gtp_encap.c` compiles with production cflags |
| test_02_has_encap_handler | `gtp_encap_handler` function loadable |
| test_03-05_has_maps | Maps exist |
| test_01_attach_to_tc | Attaches to TC egress |
| test_02_encap_golden_model | Full field-by-field: Ether MACs, IP (all fields + checksum), UDP, GTP flags/TEID/length, PDU Session Container QFI, inner packet byte-exact preservation, session counters |
| test_03_encap_small_packet | **EXPECTED FAIL** -- Finding 2: packets < 64 bytes dropped |
| test_04_encap_no_session | Dropped + no leak on output |
| test_05_encap_inactive_session | Dropped + no leak on output |
| test_06_double_encap_avoided | Already-GTP packets not re-encapsulated |
| test_07_multiple_ue_sessions | Different UEs get correct TEIDs |
| test_08_encap_double_encap_with_ipv4_options | **EXPECTED FAIL** -- Finding 4: double-encap detection broken with IPv4 options |
| test_09_encap_missing_sgi_ip | **EXPECTED FAIL** -- Finding 5: silent 0.0.0.0 source IP |
| test_10_encap_qfi_values | QFI boundary: default (0→9), normal (1), max (63) |
| test_11_encap_non_ipv4_passthrough | ARP passes through without encapsulation |
| test_12_encap_teid_zero | teid_dl_out=0 treated as inactive, packet dropped |

## Known Findings

See `.tmp/ebpf_review/FINDINGS.md` for details.

- **Finding 1**: `bpf_skb_set_tunnel_key()` requires `ip_tunnel`/`vxlan` kernel modules (undocumented)
- **Finding 2**: Both handlers drop packets < 64 bytes via `bpf_skb_load_bytes(64)` (affects TCP ACKs, ICMP, small DNS)
- **Finding 3**: `STATS_PKT_DROPPED` incremented for non-IPv4 `TC_ACT_OK` traffic (misleading counter)
- **Finding 4**: Encap double-encap detection broken with IPv4 options (hardcoded UDP port offsets)
- **Finding 5**: Missing `CONFIG_SGI_IP` causes silent 0.0.0.0 outer source IP (no error counter)
