#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

/* TC action return codes */
#define TC_ACT_OK 0
#define TC_ACT_SHOT 2

/* GTP-U constants */
#define GTPU_UDP_PORT 2152
#define GTPU_FLAGS_MASK 0x30  /* Check for extension/header flags if needed */

/* Perf event: event types */
enum event_type {
    EVT_SESSION_LOOKUP_HIT = 1,
    EVT_SESSION_LOOKUP_MISS = 2,
    EVT_DROPPED_BY_POLICY = 3,
};

/* Session action types (simple) */
enum session_action {
    ACTION_COUNT_ONLY = 0,
    ACTION_REDIRECT_IFINDEX = 1,
    ACTION_DROP = 2,
};

/* Key used for TEID map (4-byte TEID) */
struct teid_key {
    __u32 teid;
} __attribute__((packed));

/* Key used for UE IP map (IPv4 address) */
struct ue_ip_key {
    __u32 ip4; /* network (big-endian) order */
} __attribute__((packed));

/* Session value stored in maps */
struct session_value {
    __u32 ifindex;
    __u32 rx_packets;
    __u64 rx_bytes;
    __u32 action;
    __u32 flags;
} __attribute__((packed));

/* Counters map indexed by a small id (0..N-1). Could be per-UE or global */
struct counter_key {
    __u32 id;
};

struct counter_value {
    __u64 packets;
    __u64 bytes;
};

/* Maps */

/* TEID -> session (hash) */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 16384);
    __type(key, struct teid_key);
    __type(value, struct session_value);
} teid_map SEC(".maps");

/* UE IP -> session (hash) */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 16384);
    __type(key, struct ue_ip_key);
    __type(value, struct session_value);
} ue_ip_map SEC(".maps");

/* dev map (optional) for redirect via devmap/ifindex */
struct {
    __uint(type, BPF_MAP_TYPE_DEVMAP);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
    __uint(max_entries, 256);
} dev_map SEC(".maps");

/* small counters map */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 256);
    __type(key, struct counter_key);
    __type(value, struct counter_value);
} counters_map SEC(".maps");

/* perf ring buffer for events */
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
} perf_events SEC(".maps");

/* Event pushed to perf */
struct my_perf_event {   // <-- Renamed to avoid conflict
    __u32 type;
    __u32 reason;
    __u32 ifindex;
    __u32 pad;
    __u64 pkt_len;
    __u32 teid;
    __u32 src_ip;
    __u32 dst_ip;
};

/* helpers for parsing safely: returns pointer or NULL */
static __always_inline void *parse_eth(struct __sk_buff *skb, void *data, void *data_end, struct ethhdr **ethh)
{
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return NULL;
    *ethh = eth;
    return (void *)(eth + 1);
}

static __always_inline struct iphdr *parse_ipv4(void *cursor, void *data_end)
{
    struct iphdr *iph = cursor;
    if ((void *)(iph + 1) > data_end)
        return NULL;
    if (iph->version != 4)
        return NULL;
    return iph;
}

static __always_inline struct udphdr *parse_udp(void *cursor, void *data_end)
{
    struct udphdr *udph = cursor;
    if ((void *)(udph + 1) > data_end)
        return NULL;
    return udph;
}

/* GTP-U header (basic 8-byte header; optional extension headers ignored here) */
struct gtpu_hdr {
    __u8 flags;
    __u8 msg_type;
    __u16 msg_len;
    __u32 teid;
} __attribute__((packed));

static __always_inline struct gtpu_hdr *parse_gtpu(void *cursor, void *data_end)
{
    struct gtpu_hdr *g = cursor;
    if ((void *)(g + 1) > data_end)
        return NULL;
    return g;
}

/* Update counters in counters_map for key 0 (global) */
static __always_inline void update_global_counters(__u64 bytes)
{
    struct counter_key key = { .id = 0 };
    struct counter_value *val;
    val = bpf_map_lookup_elem(&counters_map, &key);
    if (!val) {
        struct counter_value init = { .packets = 1, .bytes = bytes };
        bpf_map_update_elem(&counters_map, &key, &init, BPF_ANY);
    } else {
        __sync_fetch_and_add(&val->packets, 1);
        __sync_fetch_and_add(&val->bytes, bytes);
    }
}

/* Emit perf event */
static __always_inline void emit_perf_event(struct __sk_buff *skb, enum event_type type, __u32 ifindex, __u32 teid, __u32 src_ip, __u32 dst_ip)
{
    struct my_perf_event evt = {};   // <-- updated here
    evt.type = (uint32_t)type;
    evt.ifindex = ifindex;
    evt.pkt_len = skb->len;
    evt.teid = teid;
    evt.src_ip = src_ip;
    evt.dst_ip = dst_ip;

    bpf_perf_event_output(skb, &perf_events, BPF_F_CURRENT_CPU, &evt, sizeof(evt));
}

/* The main TC program */
SEC("tc")
int tc_core_prog(struct __sk_buff *skb)
{
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    struct ethhdr *eth;
    void *cursor;
    __u16 eth_proto;

    cursor = parse_eth(skb, data, data_end, &eth);
    if (!cursor)
        return TC_ACT_OK;

    eth_proto = bpf_ntohs(eth->h_proto);
    if (eth_proto == ETH_P_8021Q || eth_proto == ETH_P_8021AD) {
        if ((void *)cursor + 4 > data_end)
            return TC_ACT_OK;
        eth_proto = bpf_ntohs(*(__u16 *)((__u8 *)cursor + 2));
        cursor = (void *)cursor + 4;
    }

    if (eth_proto != ETH_P_IP)
        return TC_ACT_OK;

    struct iphdr *iph = parse_ipv4(cursor, data_end);
    if (!iph)
        return TC_ACT_OK;

    __u32 ihl = iph->ihl * 4;
    if (ihl < sizeof(struct iphdr))
        return TC_ACT_OK;
    if ((void *)iph + ihl > data_end)
        return TC_ACT_OK;

    if (iph->protocol != IPPROTO_UDP) {
        struct ue_ip_key ukey = { .ip4 = iph->saddr };
        struct session_value *sval = bpf_map_lookup_elem(&ue_ip_map, &ukey);
        if (sval) {
            update_global_counters(skb->len);
            emit_perf_event(skb, EVT_SESSION_LOOKUP_HIT, skb->ifindex, 0, iph->saddr, iph->daddr);

            if (sval->action == ACTION_DROP) {
                emit_perf_event(skb, EVT_DROPPED_BY_POLICY, skb->ifindex, 0, iph->saddr, iph->daddr);
                return TC_ACT_SHOT;
            } else if (sval->action == ACTION_REDIRECT_IFINDEX) {
                __u32 key = sval->ifindex;
                int ret = bpf_redirect_map(&dev_map, key, 0);
                return TC_ACT_OK;
            }
            return TC_ACT_OK;
        } else {
            emit_perf_event(skb, EVT_SESSION_LOOKUP_MISS, skb->ifindex, 0, iph->saddr, iph->daddr);
            return TC_ACT_OK;
        }
    }

    struct udphdr *udph = (void *)iph + ihl;
    if ((void *)(udph + 1) > data_end)
        return TC_ACT_OK;

    __u16 udp_dport = bpf_ntohs(udph->dest);
    __u16 udp_sport = bpf_ntohs(udph->source);
    if (udp_dport != GTPU_UDP_PORT && udp_sport != GTPU_UDP_PORT)
        return TC_ACT_OK;

    struct gtpu_hdr *g = (void *)(udph + 1);
    if (!g || (void *)(g + 1) > data_end)
        return TC_ACT_OK;

    __u32 teid = bpf_ntohl(g->teid);
    struct teid_key tkey = { .teid = teid };
    struct session_value *s = bpf_map_lookup_elem(&teid_map, &tkey);
    if (!s) {
        emit_perf_event(skb, EVT_SESSION_LOOKUP_MISS, skb->ifindex, teid, iph->saddr, iph->daddr);
        return TC_ACT_OK;
    }

    update_global_counters(skb->len);
    emit_perf_event(skb, EVT_SESSION_LOOKUP_HIT, skb->ifindex, teid, iph->saddr, iph->daddr);

    if (s->action == ACTION_DROP) {
        emit_perf_event(skb, EVT_DROPPED_BY_POLICY, skb->ifindex, teid, iph->saddr, iph->daddr);
        return TC_ACT_SHOT;
    } else if (s->action == ACTION_REDIRECT_IFINDEX) {
        __u32 dst_if = s->ifindex;
        int r = bpf_redirect_map(&dev_map, dst_if, 0);
        if (r == TC_ACT_OK || r == TC_ACT_SHOT)
            return r;
        bpf_clone_redirect(skb, dst_if, 0);
        return TC_ACT_OK;
    }

    return TC_ACT_OK;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
