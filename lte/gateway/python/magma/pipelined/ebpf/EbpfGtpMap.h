/**
 * Copyright 2025 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Author: Nitin Rajput (coRAN LABS)
 *
 * eBPF GTP-U Header and Map Definitions
 * Ubuntu 20.04 (Focal) compatible - Kernel 5.4+
 * Replaces GTP kernel module with eBPF implementation
 */

#ifndef _EBPF_GTP_MAP_H
#define _EBPF_GTP_MAP_H

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/udp.h>
#include <linux/in.h>

#define ETH_P_IP_BE 0x0008      // ETH_P_IP (0x0800) in network byte order  
#define GTP_PORT_NO_BE 0x6808   // 2152 (0x0868) in network byte order

static inline __be16 htons_bcc(__u16 val) {
    return (__be16)(((val & 0x00ff) << 8) | ((val & 0xff00) >> 8));
}

static inline __u16 ntohs_bcc(__be16 val) {
    return (__u16)(((val & 0x00ff) << 8) | ((val & 0xff00) >> 8));
}

static inline __be32 htonl_bcc(__u32 val) {
    return (__be32)(((val & 0x000000ff) << 24) |
                    ((val & 0x0000ff00) << 8) |
                    ((val & 0x00ff0000) >> 8) |
                    ((val & 0xff000000) >> 24));
}

static inline __u32 ntohl_bcc(__be32 val) {
    return (__u32)(((val & 0x000000ff) << 24) |
                   ((val & 0x0000ff00) << 8) |
                   ((val & 0x00ff0000) >> 8) |
                   ((val & 0xff000000) >> 24));
}

// GTP-U Protocol Constants (3GPP TS 29.281)
#define GTP_PORT_NO 2152
#define GTP_HDR_SIZE_MIN 8          // Minimum GTP header size
#define GTP_HDR_SIZE_MAX 16         // Maximum GTP header size (with optional fields)
#define GTP_VERSION_1 0x01
#define GTP_PT_FLAG 0x01            // Protocol Type flag
#define GTP_MSG_TPDU 0xFF           // T-PDU message type

// GTP Header Flags
#define GTP_FLAG_VERSION_MASK 0xE0  // Version (3 bits)
#define GTP_FLAG_PT 0x10            // Protocol Type
#define GTP_FLAG_RESERVED 0x08      // Reserved bit
#define GTP_FLAG_E 0x04             // Extension Header flag
#define GTP_FLAG_S 0x02             // Sequence Number flag  
#define GTP_FLAG_PN 0x01            // N-PDU Number flag

// GTP-U header structure (3GPP TS 29.281)
struct gtp1_header {
    __u8 flags;                     // Version, PT, Reserved, E, S, PN flags
    __u8 type;                      // Message type (0xFF for T-PDU)
    __be16 length;                  // Length of payload + optional fields
    __be32 teid;                    // Tunnel Endpoint Identifier
} __attribute__((packed));

struct gtp1_optional {
    __be16 seq;                     // Sequence number (if S flag set)
    __u8 npdu;                      // N-PDU number (if PN flag set)
    __u8 next_ext;                  // Next extension header type (if E flag set)
} __attribute__((packed));

// UE session map key structure
struct ue_session_key {
    __be32 ue_ip;                   // UE IPv4 address
};

// UE session information structure
struct ue_session_info {
    // Tunnel information
    __be32 enb_ip;                  // eNodeB IPv4 address
    __u32 teid_ul_in;               // TEID for uplink (eNB->UE)
    __u32 teid_ul_out;              // TEID for uplink response
    __u32 teid_dl_in;               // TEID for downlink (UE->eNB)
    __u32 teid_dl_out;              // TEID for downlink response
    
    // Interface information
    __u32 s1u_ifindex;              // S1-U interface index (eth1)
    __u32 sgi_ifindex;              // SGi interface index (eth0)
    __u32 ovs_ifindex;              // OVS interface index (gtp_veth0)
    
    // MAC addresses for packet forwarding
    __u8 ul_mac_src[ETH_ALEN];      // UL source MAC (for OVS)
    __u8 ul_mac_dst[ETH_ALEN];      // UL destination MAC (for OVS)
    
    // QoS and classification
    __u32 qos_mark;                 // QoS mark for traffic classification
    __u32 bearer_id;                // EPS Bearer ID
    
    // Statistics counters
    __u64 ul_bytes;                 // Uplink byte counter
    __u64 dl_bytes;                 // Downlink byte counter
    __u64 ul_packets;               // Uplink packet counter
    __u64 dl_packets;               // Downlink packet counter
    
    // Session metadata
    __u64 last_seen;                // Last activity timestamp
    __u32 session_flags;            // Session flags (active, suspended, etc.)
    
    // IMSI storage (for logging/debugging)
    __u8 imsi[16];                  // IMSI bytes
    __u8 imsi_len;                  // IMSI length
    __u64 encoded_imsi;             // Pre-encoded IMSI for OVS metadata
    __u8 qfi;                       // QoS Flow Identifier
    __u32 tunnel_id;                // GTP tunnel ID for metadata
    __be32 tun_ipv4_dst;            // GTP destination IP
    __u8 tun_flags;                 // Tunnel flags
    __u8 direction;                 // Traffic direction marker
    __u32 original_port;            // Original GTP port
    __u8 reserved[3];               // Padding/alignment
    __u32 metadata_mark;            // UE mark for OVS matching
};

// Configuration map key structure
struct config_key {
    __u32 key;                      // Configuration key
};

// Configuration values
#define CONFIG_S1U_IFINDEX 0        // S1-U interface index
#define CONFIG_SGI_IFINDEX 1        // SGi interface index  
#define CONFIG_OVS_IFINDEX 2        // OVS interface index
#define CONFIG_DEBUG_LEVEL 3        // Debug level

// Statistics map key structure
struct stats_key {
    __u32 cpu;                      // CPU number
    __u32 counter_id;               // Counter ID
};

// Global statistics counters - Step 6: Full statistics (per-UE metrics)
#define STATS_UL_PACKETS 0          // Total UL packets
#define STATS_DL_PACKETS 1          // Total DL packets
#define STATS_UL_BYTES 2            // Total UL bytes
#define STATS_DL_BYTES 3            // Total DL bytes
#define STATS_UL_ERRORS 4           // UL processing errors
#define STATS_DL_ERRORS 5           // DL processing errors
#define STATS_GTP_INVALID 6         // Invalid GTP packets
#define STATS_SESSION_MISS 7        // Session lookup misses
#define STATS_GTP_DECAP_SUCCESS 8   // Successful GTP decapsulation
#define STATS_GTP_ENCAP_SUCCESS 9   // Successful GTP encapsulation
#define STATS_TEID_MISMATCH 10      // TEID mismatches
#define STATS_UE_ATTACH 11          // UE attach events
#define STATS_UE_DETACH 12          // UE detach events
#define STATS_PKT_DROPPED 13        // Packets dropped
#define STATS_PKT_FORWARDED 14      // Packets forwarded
#define STATS_MAX_COUNTERS 15       // Maximum number of counters

BPF_HASH(ue_session_map, struct ue_session_key, struct ue_session_info, 1024);
BPF_ARRAY(config_map, __u32, 8);
BPF_PERCPU_ARRAY(stats_map, __u64, STATS_MAX_COUNTERS);

static inline struct gtp1_header* get_gtp_header(void* data, void* data_end) {
    struct ethhdr* eth = data;
    struct iphdr* ip;
    struct udphdr* udp;
    struct gtp1_header* gtp;
    
    if ((void*)(eth + 1) > data_end)
        return NULL;
    
    if (eth->h_proto != ETH_P_IP_BE)
        return NULL;
    
    ip = (struct iphdr*)(eth + 1);
    if ((void*)(ip + 1) > data_end)
        return NULL;
    
    if (ip->protocol != IPPROTO_UDP)
        return NULL;
    
    udp = (struct udphdr*)((void*)ip + (ip->ihl * 4));
    if ((void*)(udp + 1) > data_end)
        return NULL;
    
    if (udp->dest != GTP_PORT_NO_BE)
        return NULL;
    
    gtp = (struct gtp1_header*)(udp + 1);
    if ((void*)(gtp + 1) > data_end)
        return NULL;
    
    return gtp;
}

static inline struct iphdr* get_inner_ip(struct gtp1_header* gtp, void* data_end) {
    struct iphdr* inner_ip;
    __u8* gtp_payload;
    __u32 gtp_hdr_len = GTP_HDR_SIZE_MIN;
    
    if (gtp->flags & (GTP_FLAG_E | GTP_FLAG_S | GTP_FLAG_PN)) {
        gtp_hdr_len += 4;  // Add optional fields size
    }
    
    gtp_payload = (__u8*)gtp + gtp_hdr_len;
    inner_ip = (struct iphdr*)gtp_payload;
    
    if ((void*)(inner_ip + 1) > data_end)
        return NULL;
    
    return inner_ip;
}

static inline void update_stats(__u32 counter_id, __u64 value) {
    __u64* count = bpf_map_lookup_elem(&stats_map, &counter_id);
    if (count) {
        *count += value;  // Simple addition instead of atomic
    }
}

static inline __u64 get_timestamp_ns(void) {
    return bpf_ktime_get_ns();
}

#define GTP_DEBUG(fmt, ...) do {} while (0)

#endif /* _EBPF_GTP_MAP_H */
