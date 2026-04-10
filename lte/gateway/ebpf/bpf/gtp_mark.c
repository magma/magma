
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/pkt_cls.h>

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
``
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/pkt_cls.h>

// Define ETH_P_IP if not defined
#ifndef ETH_P_IP
#define ETH_P_IP 0x0800
#endif

// Define ETH_HLEN if not defined
#ifndef ETH_HLEN
#define ETH_HLEN 14
#endif

// Use standard TC action codes
#ifndef TC_ACT_OK
#define TC_ACT_OK 0
#endif

// UE session structures (must match main GTP program)
struct ue_session_key {
    __be32 ue_ip;
};

struct ue_session_info {
    __be32 enb_ip;
    __u32 teid_ul_in;
    __u32 teid_ul_out;
    __u32 teid_dl_in;
    __u32 teid_dl_out;
    __u32 s1u_ifindex;
    __u32 sgi_ifindex;
    __u32 ovs_ifindex;
    __u8 ul_mac_src[6];
    __u8 ul_mac_dst[6];
    __u32 qos_mark;
    __u32 bearer_id;
    __u64 ul_bytes;
    __u64 dl_bytes;
    __u64 ul_packets;
    __u64 dl_packets;
    __u64 last_seen;
    __u32 session_flags;
    __u8 imsi[16];
    __u32 imsi_len;
    __u64 encoded_imsi;
    __u8 qfi;
    __u32 tunnel_id;
    __be32 tun_ipv4_dst;
    __u8 tun_flags;
    __u8 direction;
    __u32 original_port;
    __u8 reserved[3];
    __u32 metadata_mark;    // CRITICAL: Stored metadata mark from gtp_veth1 program
};

// BPF map definitions - these MUST reference the same maps as main GTP program
// Note: Map will be pinned manually from Python for sharing between programs (BCC 0.12.0 compatible)
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, struct ue_session_key);
    __type(value, struct ue_session_info);
} ue_session_map SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 64);
    __type(key, __u32);
    __type(value, __u64);
} stats_map SEC(".maps");


// Statistics map for monitoring
BPF_ARRAY(stats_map, __u64, 64);

// Statistics counters
#define STATS_VETH0_PACKETS_PROCESSED 50
#define STATS_VETH0_MARK_RESTORED 51
#define STATS_VETH0_MARK_FALLBACK 52
#define STATS_VETH0_SESSION_MISS 53

// Helper function to update statistics
sstatic __always_inline void update_stats(__u32 counter_id, __u64 value)
{
    __u64 *count = bpf_map_lookup_elem(&stats_map, &counter_id);
    if (count)
        *count += value;
}

static inline __u32 compute_ue_mark(__u32 ue_ip_int) {
    __u32 safe_mark = ue_ip_int & 0x7FFFFFFEU;

    if (safe_mark == 0x7FFFFFFF || safe_mark == 0) {
        safe_mark = (ue_ip_int >> 8) | 0x12345600U;
    }

    if (safe_mark < 0x10000000U) {
        safe_mark |= 0x12000000U;
    }

    return safe_mark;
}

/**
 * CRITICAL: gtp_veth0 Metadata Mark Handler
 *
 * This function restores metadata marks that were stored by the
 * gtp_veth1 program but lost during bpf_redirect.
 *
 * Flow:
 * 1. Extract source IP from packet (UE IP)
 * 2. Look up session info in shared ue_session_map
 * 3. Restore metadata_mark from session info to skb->mark
 * 4. Continue to OVS with proper mark set
 */
 SEC("tc")
int gtp_veth0_mark_handler(struct __sk_buff* skb) {
    // Update packet processing stats
    update_stats(STATS_VETH0_PACKETS_PROCESSED, 1);

    // Load first part of packet to extract source IP
    __u8 pkt_data[40];
    if (bpf_skb_load_bytes(skb, 0, pkt_data, sizeof(pkt_data)) < 0) {
        return TC_ACT_OK;  // Pass through on error
    }

    // Check if it's IPv4 packet
    __u16 eth_type = (__u16)pkt_data[12] << 8 | (__u16)pkt_data[13];
    if (eth_type != ETH_P_IP) {
        return TC_ACT_OK;  // Pass through non-IPv4
    }

    // Validate IP header
    if (ETH_HLEN + 20 > sizeof(pkt_data)) {
        return TC_ACT_OK;  // Pass through if insufficient data
    }

    __u8 ip_version = (pkt_data[ETH_HLEN] >> 4) & 0xF;
    if (ip_version != 4) {
        return TC_ACT_OK;  // Pass through non-IPv4
    }

    __be32 src_ip = *((__be32 *)&pkt_data[ETH_HLEN + 12]);
    

    bpf_printk("[VETH0] Processing packet from UE IP: 0x%x", src_ip);

    struct ue_session_key session_key = {.ue_ip = src_ip};
    struct ue_session_info *session_info =
    bpf_map_lookup_elem(&ue_session_map, &session_key);

    if (session_info == NULL) {
        bpf_trace_printk("[VETH0] No session found for IP 0x%x\n", src_ip);
        update_stats(STATS_VETH0_SESSION_MISS, 1);
        return TC_ACT_OK;  // Pass through if no session
    }

    // Check if session is active
    if (!(session_info->session_flags & 1)) {
        bpf_trace_printk("[VETH0] Session inactive for IP 0x%x\n", src_ip);
        return TC_ACT_OK;  
    }

    __u32 ue_ip_int = ((__u32)pkt_data[ETH_HLEN + 12] << 24) |
                      ((__u32)pkt_data[ETH_HLEN + 13] << 16) |
                      ((__u32)pkt_data[ETH_HLEN + 14] << 8) |
                      (__u32)pkt_data[ETH_HLEN + 15];

    // Phase‑1 Fix (Finding 6): Restore metadata mark from network byte order
__u32 restored_mark = bpf_ntohl(session_info->metadata_mark);
skb->mark = restored_mark;

    bpf_trace_printk("[VETH0] Restored metadata mark: 0x%x for UE 0x%x\n",
                     restored_mark, src_ip);
    update_stats(STATS_VETH0_MARK_RESTORED, 1);

    return TC_ACT_OK;  
}
char _license[] SEC("license") = "GPL";