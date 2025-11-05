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
 * eBPF GTP-U Decapsulation Program - TC eBPF Implementation
 * Handles GTP decapsulation in TC ingress context
 */

#ifdef COMPILE_TIME
static const int compile_timestamp = COMPILE_TIME;
#endif

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <linux/in.h>
#include <linux/pkt_cls.h>

#ifndef BPF_F_TUNINFO_IPV4
#define BPF_F_TUNINFO_IPV4 0
#endif

#ifndef ETH_P_IP
#define ETH_P_IP 0x0800
#endif

#ifndef ETH_ALEN
#define ETH_ALEN 6
#endif

#ifndef ETH_HLEN
#define ETH_HLEN 14
#endif

#ifndef TC_ACT_OK
#define TC_ACT_OK 0
#endif
#ifndef TC_ACT_SHOT
#define TC_ACT_SHOT 2
#endif
#ifndef TC_ACT_REDIRECT
#define TC_ACT_REDIRECT 7
#endif

#define BPF_ADJ_ROOM_NET 0
#define BPF_ADJ_ROOM_MAC 1

#define BPF_F_ADJ_ROOM_FIXED_GSO (1ULL << 0)

struct ue_session_key {
    __be32 ue_ip;
};

struct ue_session_info {
    __be32 enb_ip;
    __u32 teid_ul_in;    // TEID for uplink (eNB->UE)
    __u32 teid_ul_out;   // TEID for uplink response
    __u32 teid_dl_in;    // TEID for downlink (UE->eNB)
    __u32 teid_dl_out;   // TEID for downlink response
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
    __u32 metadata_mark;    // CRITICAL: Store metadata mark for gtp_veth0 program
};

struct config_key {
    __u32 key;
};

struct config_value {
    __u32 value;
};

BPF_HASH(ue_session_map, struct ue_session_key, struct ue_session_info, 1024);
BPF_HASH(config_map, struct config_key, struct config_value, 16);

BPF_ARRAY(stats_map, __u64, 64);

#define STATS_UL_PACKETS 0
#define STATS_UL_BYTES 1
#define STATS_DL_PACKETS 2
#define STATS_DL_BYTES 3
#define STATS_UL_ERRORS 4
#define STATS_DL_ERRORS 5
#define STATS_SESSION_MISS 6
#define STATS_TEID_MISMATCH 7
#define STATS_GTP_DECAP_SUCCESS 8
#define STATS_GTP_ENCAP_SUCCESS 9
#define STATS_PKT_TOO_SHORT 10
#define STATS_INVALID_GTP 11
#define STATS_ADJUST_HEAD_FAIL 12
#define STATS_TOTAL_PROCESSED 13
#define STATS_UE_ATTACH 14
#define STATS_UE_DETACH 15
#define STATS_PKT_FORWARDED 16
#define STATS_PKT_DROPPED 17
#define STATS_SESSION_ACTIVE 18
#define STATS_QOS_APPLIED 19
#define STATS_INACTIVE_SESSION 20

#define CONFIG_S1U_IFINDEX 0
#define CONFIG_SGI_IFINDEX 1
#define CONFIG_OVS_IFINDEX 2
#define CONFIG_DEBUG_LEVEL 3
#define CONFIG_SGI_IP 4
#define CONFIG_EBPF_VETH_IFINDEX 5

static inline void update_stats(__u32 counter_id, __u64 value) {
    __u64* count = stats_map.lookup(&counter_id);
    if (count) {
        *count += value;
    }
}

static inline __u32 get_outer_src_ip(struct __sk_buff *skb) {
    __u8 ip_data[4];
    if (bpf_skb_load_bytes(skb, ETH_HLEN + 12, ip_data, 4) < 0) {
        return 0;
    }
    return ((__u32)ip_data[0] << 24) | ((__u32)ip_data[1] << 16) |
           ((__u32)ip_data[2] << 8) | (__u32)ip_data[3];
}

static inline __u32 get_outer_dst_ip(struct __sk_buff *skb) {
    __u8 ip_data[4];
    if (bpf_skb_load_bytes(skb, ETH_HLEN + 16, ip_data, 4) < 0) {
        return 0;
    }
    return ((__u32)ip_data[0] << 24) | ((__u32)ip_data[1] << 16) |
           ((__u32)ip_data[2] << 8) | (__u32)ip_data[3];
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

static inline __u16 ip_checksum(__u8 *data, int len) {
    __u32 sum = 0;

    #pragma unroll
    for (int i = 0; i < 10; i++) {
        if (i * 2 < len) {
            sum += ((__u16)data[i * 2] << 8) | (__u16)data[i * 2 + 1];
        }
    }

    #pragma unroll
    for (int i = 0; i < 2; i++) {  
        if (sum >> 16) {
            sum = (sum & 0xFFFF) + (sum >> 16);
        }
    }

    return ~sum;
}

// GTP-U Protocol Constants (3GPP TS 29.281)
#define GTP_HDR_SIZE_MIN 8
#define GTP_VERSION_1 0x01
#define GTP_PT_FLAG 0x01
#define GTP_MSG_TPDU 0xFF

#define GTP_FLAG_VERSION_MASK 0xE0
#define GTP_FLAG_PT 0x10
#define GTP_FLAG_E 0x04
#define GTP_FLAG_S 0x02
#define GTP_FLAG_PN 0x01

struct gtp1_header {
    __u8 flags;
    __u8 type;
    __be16 length;
    __be32 teid;
} __attribute__((packed));


int gtp_decap_handler(struct __sk_buff* skb) {
    bpf_trace_printk("IN: skb len = %u\n", skb->len);
    update_stats(STATS_TOTAL_PROCESSED, 1);

    __u8 pkt_data[64];
    if (bpf_skb_load_bytes(skb, 0, pkt_data, sizeof(pkt_data)) < 0) {
        update_stats(STATS_PKT_TOO_SHORT, 1);
        bpf_trace_printk("[GTP] Failed to load first 64 bytes\n");
        return TC_ACT_SHOT;  // Drop malformed packets
    }

    bpf_trace_printk("[ETH] Src bytes: %u %u %u\n", pkt_data[6], pkt_data[7], pkt_data[8]);
    bpf_trace_printk("[ETH] Src bytes cont: %u %u %u\n", pkt_data[9], pkt_data[10], pkt_data[11]);
    bpf_trace_printk("[ETH] Dst bytes: %u %u %u\n", pkt_data[0], pkt_data[1], pkt_data[2]);
    bpf_trace_printk("[ETH] Dst bytes cont: %u %u %u\n", pkt_data[3], pkt_data[4], pkt_data[5]);

    __u16 eth_type = (__u16)pkt_data[12] << 8 | (__u16)pkt_data[13];
    bpf_trace_printk("[ETH] Ethertype: 0x%x\n", eth_type);

    if (eth_type != ETH_P_IP) {
        update_stats(STATS_PKT_DROPPED, 1);
        bpf_trace_printk("[GTP] Non-IPv4 outer pkt, skipping\n");
        return TC_ACT_OK;  
    }
    update_stats(22, 1); 
    bpf_trace_printk("[Parsing IP Header...]\n");

    if (ETH_HLEN + sizeof(struct iphdr) > sizeof(pkt_data)) {
        bpf_trace_printk("[GTP] Not enough data for IPv4 header\n");
        return TC_ACT_SHOT;
    }

    __u8 ip_version_ihl = pkt_data[ETH_HLEN + 0];
    __u8 ip_protocol     = pkt_data[ETH_HLEN + 9];
    __u8 ip_version      = (ip_version_ihl >> 4) & 0xF;
    __u8 ihl             = ip_version_ihl & 0xF;

    bpf_trace_printk("[IP] Ver:%u IHL:%u Proto:%u\n", ip_version, ihl, ip_protocol);
    bpf_trace_printk("[IP] Src part1: %u.%u", pkt_data[14], pkt_data[15]);
    bpf_trace_printk(".%u.%u\n", pkt_data[16], pkt_data[17]);
    bpf_trace_printk("[IP] Dst part1: %u.%u", pkt_data[18], pkt_data[19]);
    bpf_trace_printk(".%u.%u\n", pkt_data[20], pkt_data[21]);

    if (ip_version != 4 || ip_protocol != IPPROTO_UDP) {
        return TC_ACT_OK; 
    }
    update_stats(23, 1); 

    __u32 ip_hlen;
    if (ihl >= 5 && ihl <= 15) ip_hlen = ihl * 4;
    else return TC_ACT_SHOT;  

    __u8 udp_data[8];
    if (bpf_skb_load_bytes(skb, ETH_HLEN + ip_hlen, udp_data, sizeof(udp_data)) < 0) {
        bpf_trace_printk("[GTP] Unable to load UDP header\n");
        return TC_ACT_SHOT;  
    }

    __u16 udp_src_port = (__u16)udp_data[0] << 8 | (__u16)udp_data[1];
    __u16 udp_dest_port= (__u16)udp_data[2] << 8 | (__u16)udp_data[3];
    bpf_trace_printk("[UDP] Src:%u Dst:%u\n", udp_src_port, udp_dest_port);

    if (udp_dest_port != 2152) {
        bpf_trace_printk("[GTP] Not GTP-U (udp dst != 2152)\n");
        return TC_ACT_OK;  
    }
    update_stats(24, 1); 
    bpf_trace_printk("[GTP] GTP-U detected, proceeding to decap\n");

    __u32 gtp_offset = ETH_HLEN + ip_hlen + sizeof(struct udphdr);
    if (gtp_offset > 78) {
        bpf_trace_printk("[GTP] GTP offset too large\n");
        return TC_ACT_SHOT;
    }

    __u8 gtp_data[8];
    if (bpf_skb_load_bytes(skb, gtp_offset, gtp_data, sizeof(gtp_data)) < 0) {
        bpf_trace_printk("[GTP] Unable to load GTP header\n");
        return TC_ACT_SHOT;
    }

    __u8 gtp_flags = gtp_data[0];
    __u8 gtp_type  = gtp_data[1];
    __u32 gtp_teid = (__u32)gtp_data[4] << 24 | (__u32)gtp_data[5] << 16 |
                     (__u32)gtp_data[6] << 8  | (__u32)gtp_data[7];
    __u8 gtp_version = (gtp_flags & GTP_FLAG_VERSION_MASK) >> 5;

    bpf_trace_printk("[GTP] Flags: 0x%x Ver:%u Type:%u\n", gtp_flags, gtp_version, gtp_type);
    bpf_trace_printk("[GTP] TEID: %u\n", gtp_teid);

    if (gtp_version != GTP_VERSION_1 || gtp_type != GTP_MSG_TPDU || !(gtp_flags & GTP_FLAG_PT)) {
        bpf_trace_printk("[GTP] Invalid GTP header\n");
        return TC_ACT_SHOT;
    }

    __u32 gtp_hdr_len = GTP_HDR_SIZE_MIN;
    if (gtp_flags & (GTP_FLAG_E | GTP_FLAG_S | GTP_FLAG_PN)) {
        gtp_hdr_len += 4;  
        if (gtp_offset + gtp_hdr_len > skb->len) {
            bpf_trace_printk("[GTP] GTP optional fields exceed skb len\n");
            return TC_ACT_SHOT;
        }
    }

    if (gtp_flags & GTP_FLAG_E) {
        __u8 next_ext_type = 0;
        __u32 ext_offset = gtp_offset + gtp_hdr_len - 1;  // Point to next extension header type field
        
        bpf_trace_printk("[GTP] Extension flag set, parsing headers at offset %u\n", ext_offset);
        
        if (bpf_skb_load_bytes(skb, ext_offset, &next_ext_type, 1) < 0) {
            bpf_trace_printk("[GTP] Failed to read next ext header type\n");
            return TC_ACT_SHOT;
        }
        
        bpf_trace_printk("[GTP] First ext header type: 0x%x at offset %u\n", next_ext_type, ext_offset);
        
        #pragma unroll
        for (int i = 0; i < 5 && next_ext_type != 0; i++) {  // Max 5 extensions to avoid verifier issues
            __u8 ext_hdr[4];
            ext_offset++;  // Move to start of extension header
            
            if (bpf_skb_load_bytes(skb, ext_offset, ext_hdr, sizeof(ext_hdr)) < 0) {
                bpf_trace_printk("[GTP] Failed to load ext header %d\n", i);
                return TC_ACT_SHOT;
            }
            
            __u8 ext_len = ext_hdr[0];  // Length in 4-byte units
            
            // QoS Flow Identifier (QFI) for 5G
            if (next_ext_type == 0x85) {  // PDU Session Container
                __u8 qfi = ext_hdr[1] & 0x3F;
                bpf_trace_printk("[5G-QoS] QFI: %u\n", qfi);
            }
            
            __u32 ext_hdr_size = ext_len * 4;
            if (ext_hdr_size < 4) {
                bpf_trace_printk("[GTP] Invalid ext header length: %u\n", ext_len);
                return TC_ACT_SHOT;
            }
            
            bpf_trace_printk("[GTP] Ext %d: size=%u, hdr_len=%u\n", 
                             i, ext_hdr_size, gtp_hdr_len);
            
            gtp_hdr_len += ext_hdr_size;
            ext_offset += ext_hdr_size - 1;  // Move to next ext header type field
            
            if (bpf_skb_load_bytes(skb, ext_offset, &next_ext_type, 1) < 0) {
                bpf_trace_printk("[GTP] Failed to read next ext type\n");
                break;
            }
            
            bpf_trace_printk("[GTP] Ext header %d: len=%u, next=0x%x\n", i, ext_len, next_ext_type);
        }
    }
    
    bpf_trace_printk("[GTP] Total GTP header len: %u bytes\n", gtp_hdr_len);

    __u32 outer_hdr_len = ip_hlen + sizeof(struct udphdr) + gtp_hdr_len;
    __u32 inner_ip_offset = ETH_HLEN + outer_hdr_len;
    
    bpf_trace_printk("[GTP] Offsets: outer_hdr=%u, inner_ip=%u\n", outer_hdr_len, inner_ip_offset);
    
    if (inner_ip_offset > skb->len) {
        bpf_trace_printk("[GTP] inner offset %u beyond skb len %u\n", inner_ip_offset, skb->len);
        return TC_ACT_SHOT;
    }

    struct bpf_tunnel_key tun_key = {};
    tun_key.tunnel_id = (__u64)gtp_teid;     // TEID as tunnel ID
    tun_key.remote_ipv4 = get_outer_src_ip(skb);  // eNB IP
    tun_key.tunnel_tos = 0;
    tun_key.tunnel_ttl = 64;
    
    if (bpf_skb_set_tunnel_key(skb, &tun_key, sizeof(tun_key), 0) < 0) {
        bpf_trace_printk("[GTP] Failed to set tunnel metadata\n");
        update_stats(STATS_UL_ERRORS, 1);
        return TC_ACT_SHOT;
    }
    
    struct gtpu_metadata {
        __u8 ver;
        __u8 flags;
        __u8 type;
    } gtpu_opts = {
        .ver = GTP_VERSION_1,
        .flags = gtp_flags,
        .type = gtp_type
    };
    
    if (bpf_skb_set_tunnel_opt(skb, &gtpu_opts, sizeof(gtpu_opts)) < 0) {
        bpf_trace_printk("[GTP] Failed to set GTP tunnel options\n");
    }
    
    __u8 inner_ip_data[20];
    struct ue_session_info* session_info = NULL;
    
    __u8 debug_bytes[8];
    if (bpf_skb_load_bytes(skb, inner_ip_offset - 4, debug_bytes, sizeof(debug_bytes)) == 0) {
        bpf_trace_printk("[GTP] Before inner IP: %x %x %x\n",
                         debug_bytes[0], debug_bytes[1], debug_bytes[2]);
        bpf_trace_printk("[GTP] IP start: %x %x %x\n",
                         debug_bytes[4], debug_bytes[5], debug_bytes[6]);
    }
    
    if (bpf_skb_load_bytes(skb, inner_ip_offset, inner_ip_data, sizeof(inner_ip_data)) < 0) {
        bpf_trace_printk("[GTP] Failed to load inner IP header at offset %u\n", inner_ip_offset);
        update_stats(STATS_PKT_DROPPED, 1);
        return TC_ACT_SHOT;
    }
    
    __u8 inner_version = (inner_ip_data[0] >> 4) & 0xF;
    __u8 inner_ihl = inner_ip_data[0] & 0xF;
    
    bpf_trace_printk("[GTP] Inner IP: ver=%u, IHL=%u, byte=%x\n", 
                     inner_version, inner_ihl, inner_ip_data[0]);
    
    if (inner_version != 4) {
        bpf_trace_printk("[GTP] ERROR: Inner packet not IPv4 (ver=%u)\n", inner_version);
        bpf_trace_printk("[GTP] Invalid inner packet at offset %u\n", inner_ip_offset);
        return TC_ACT_SHOT;
    }
    
    __be32 inner_src_ip_be = *((__be32 *)&inner_ip_data[12]);
    __be32 inner_dst_ip_be = *((__be32 *)&inner_ip_data[16]);

    __u32 inner_src_ip = bpf_ntohl(inner_src_ip_be);
    __u32 inner_dst_ip = bpf_ntohl(inner_dst_ip_be);

    bpf_trace_printk("[GTP] Inner src IP (host order): 0x%x\n", inner_src_ip);

    struct ue_session_key session_key = {.ue_ip = inner_src_ip};
    bpf_trace_printk("[GTP] Looking up session for IP: 0x%x\n", session_key.ue_ip);
    
    session_info = ue_session_map.lookup(&session_key);
    
    if (session_info == NULL) {
        update_stats(STATS_SESSION_MISS, 1);
        update_stats(STATS_PKT_DROPPED, 1);
        bpf_trace_printk("[GTP] Session NOT FOUND for UE IP 0x%x\n", inner_src_ip);
        return TC_ACT_SHOT;
    }
    
    bpf_trace_printk("[GTP] Session FOUND! TEID exp: %u\n", session_info->teid_ul_in);
    
    if (!(session_info->session_flags & 1)) {
        update_stats(STATS_INACTIVE_SESSION, 1);
        update_stats(STATS_PKT_DROPPED, 1);
        bpf_trace_printk("[GTP] Session inactive\n");
        return TC_ACT_SHOT;
    }
    
    if (session_info->teid_ul_in != 0 &&
        gtp_teid != 0x7FFFFFFF &&
        gtp_teid != session_info->teid_ul_in) {
        update_stats(STATS_TEID_MISMATCH, 1);
        update_stats(STATS_PKT_DROPPED, 1);
        bpf_trace_printk("[GTP] TEID mismatch: got %u exp %u\n",
                         gtp_teid, session_info->teid_ul_in);
        return TC_ACT_SHOT;
    }
    
    session_info->ul_packets++;
    session_info->ul_bytes += skb->len;
    session_info->last_seen = bpf_ktime_get_ns();
    
    __u32 ue_ip_int = ((__u32)inner_ip_data[12] << 24) |
                      ((__u32)inner_ip_data[13] << 16) |
                      ((__u32)inner_ip_data[14] << 8) |
                      (__u32)inner_ip_data[15];

    __u32 calculated_mark = compute_ue_mark(ue_ip_int);

    session_info->metadata_mark = calculated_mark;
    skb->mark = calculated_mark;

    bpf_trace_printk("[GTP] Calculated UE mark: 0x%x for UE IP 0x%x\n",
                     calculated_mark, ue_ip_int);

    __s32 adjust_len = -(__s32)outer_hdr_len;
    int ret = bpf_skb_adjust_room(skb, adjust_len, BPF_ADJ_ROOM_MAC, 0);
    if (ret < 0) {
        bpf_trace_printk("[GTP] bpf_skb_adjust_room failed: %d\n", ret);
        update_stats(STATS_ADJUST_HEAD_FAIL, 1);
        update_stats(STATS_UL_ERRORS, 1);
        return TC_ACT_SHOT;
    }

    __u8 new_eth_hdr[ETH_HLEN];
    
    new_eth_hdr[0] = session_info->ul_mac_dst[0];
    new_eth_hdr[1] = session_info->ul_mac_dst[1];
    new_eth_hdr[2] = session_info->ul_mac_dst[2];
    new_eth_hdr[3] = session_info->ul_mac_dst[3];
    new_eth_hdr[4] = session_info->ul_mac_dst[4];
    new_eth_hdr[5] = session_info->ul_mac_dst[5];
    
    new_eth_hdr[6] = session_info->ul_mac_src[0];
    new_eth_hdr[7] = session_info->ul_mac_src[1];
    new_eth_hdr[8] = session_info->ul_mac_src[2];
    new_eth_hdr[9] = session_info->ul_mac_src[3];
    new_eth_hdr[10] = session_info->ul_mac_src[4];
    new_eth_hdr[11] = session_info->ul_mac_src[5];
    
    new_eth_hdr[12] = 0x08;
    new_eth_hdr[13] = 0x00;
    
    if (bpf_skb_store_bytes(skb, 0, new_eth_hdr, ETH_HLEN, 0) < 0) {
        bpf_trace_printk("[GTP] Failed to store new ethernet header\n");
        update_stats(STATS_UL_ERRORS, 1);
        return TC_ACT_SHOT;
    }

    update_stats(STATS_GTP_DECAP_SUCCESS, 1);
    update_stats(STATS_PKT_FORWARDED, 1);

    struct config_key ovs_key = {.key = CONFIG_OVS_IFINDEX};
    struct config_value* ovs_val = config_map.lookup(&ovs_key);
    if (ovs_val && ovs_val->value > 0) {
        bpf_trace_printk("[GTP] Redirecting to gtp_veth0 ifindex: %u (ingress)\n", ovs_val->value);
        return bpf_redirect(ovs_val->value, BPF_F_INGRESS);
    }

    if (session_info->ovs_ifindex > 0) {
        bpf_trace_printk("[GTP] Redirecting to session ovs ifindex: %u (ingress)\n", session_info->ovs_ifindex);
        return bpf_redirect(session_info->ovs_ifindex, BPF_F_INGRESS);
    }

    bpf_trace_printk("[GTP] ERROR: No valid OVS interface found, dropping packet\n");
    update_stats(STATS_UL_ERRORS, 1);
    return TC_ACT_SHOT;
}

int gtp_echo_handler(struct __sk_buff* skb) {
    __u8 pkt_data[64];
    if (bpf_skb_load_bytes(skb, 0, pkt_data, sizeof(pkt_data)) < 0) {
        return TC_ACT_SHOT;
    }

    __u16 eth_type = (__u16)pkt_data[12] << 8 | (__u16)pkt_data[13];
    if (eth_type != ETH_P_IP) return TC_ACT_OK;

    __u8 ip_protocol = pkt_data[ETH_HLEN + 9];
    if (ip_protocol != IPPROTO_UDP) return TC_ACT_OK;

    __u8 ip_version_ihl = pkt_data[ETH_HLEN + 0];
    __u8 ihl = ip_version_ihl & 0xF;
    __u32 ip_hlen = ihl * 4;

    __u8 udp_data[8];
    if (bpf_skb_load_bytes(skb, ETH_HLEN + ip_hlen, udp_data, sizeof(udp_data)) < 0) {
        return TC_ACT_OK;
    }

    __u16 udp_dest_port = (__u16)udp_data[2] << 8 | (__u16)udp_data[3];
    if (udp_dest_port != 2152) return TC_ACT_OK;

    __u32 gtp_offset = ETH_HLEN + ip_hlen + 8;
    __u8 gtp_data[8];
    if (bpf_skb_load_bytes(skb, gtp_offset, gtp_data, sizeof(gtp_data)) < 0) {
        return TC_ACT_SHOT;
    }

    __u8 gtp_type = gtp_data[1];
    
    if (gtp_type == 1 || gtp_type == 2) {
        bpf_trace_printk("[GTP] Echo %s handled\n", gtp_type == 1 ? "Request" : "Response");
        
        if (gtp_type == 1) {
            return TC_ACT_OK;
        } else {
            update_stats(STATS_GTP_DECAP_SUCCESS, 1);
            return TC_ACT_OK;
        }
    }

    return TC_ACT_OK;
}

int gtp_veth0_mark_handler(struct __sk_buff* skb) {
    __u8 pkt_data[40];
    if (bpf_skb_load_bytes(skb, 0, pkt_data, sizeof(pkt_data)) < 0) {
        return TC_ACT_OK;  // Pass through on error
    }

    __u16 eth_type = (__u16)pkt_data[12] << 8 | (__u16)pkt_data[13];
    if (eth_type != ETH_P_IP) {
        return TC_ACT_OK;  // Pass through non-IPv4
    }

    if (ETH_HLEN + 20 > sizeof(pkt_data)) {
        return TC_ACT_OK;  // Pass through if insufficient data
    }

    __u8 ip_version = (pkt_data[ETH_HLEN] >> 4) & 0xF;
    if (ip_version != 4) {
        return TC_ACT_OK;  // Pass through non-IPv4
    }

    __be32 src_ip_be = *((__be32 *)&pkt_data[ETH_HLEN + 12]);

    __u32 src_ip = bpf_ntohl(src_ip_be);

    bpf_trace_printk("[VETH0] Processing packet from UE IP: 0x%x (be: 0x%x)\n", src_ip, src_ip_be);

    struct ue_session_key session_key = {.ue_ip = src_ip};
    struct ue_session_info* session_info = ue_session_map.lookup(&session_key);

    if (session_info == NULL) {
        bpf_trace_printk("[VETH0] No session found for IP 0x%x\n", src_ip);
        return TC_ACT_OK;  // Pass through if no session
    }

    if (!(session_info->session_flags & 1)) {
        bpf_trace_printk("[VETH0] Session inactive for IP 0x%x\n", src_ip);
        return TC_ACT_OK;  // Pass through inactive sessions
    }

    __u32 ue_ip_int = ((__u32)pkt_data[ETH_HLEN + 12] << 24) |
                      ((__u32)pkt_data[ETH_HLEN + 13] << 16) |
                      ((__u32)pkt_data[ETH_HLEN + 14] << 8) |
                      (__u32)pkt_data[ETH_HLEN + 15];

    __u32 restored_mark = compute_ue_mark(ue_ip_int);

    session_info->metadata_mark = restored_mark;
    skb->mark = restored_mark;

    bpf_trace_printk("[VETH0] Restored metadata mark: 0x%x for UE 0x%x\n",
                     restored_mark, src_ip);

    update_stats(STATS_PKT_FORWARDED, 1);

    return TC_ACT_OK;  // Continue to OVS with mark set
}

int gtp_passthrough_handler(struct __sk_buff* skb) {
    return TC_ACT_OK;
}
