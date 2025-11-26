

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
    __u32 ifindex;        /* interface index to redirect to (if action == REDIRECT_IFINDEX) */
    __u32 rx_packets;     /* simple counters kept in user-space but mirrored here */
    __u64 rx_bytes;
    __u32 action;         /* enum session_action */
    __u32 flags;          /* misc session flags */
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

/* dev map (optional) for redirect via devmap/ifindex) - user-space can populate */
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
struct perf_event {
    __u32 type;
    __u32 reason;     /* more details */
    __u32 ifindex;    /* where we would redirect or the ingress ifindex */
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

/* Update counters in counters_map for key 0 (global) - keep it tiny for verifier */
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
    struct perf_event evt = {};
    evt.type = (uint32_t)type;
    evt.ifindex = ifindex;
    evt.pkt_len = skb->len;
    evt.teid = teid;
    evt.src_ip = src_ip;
    evt.dst_ip = dst_ip;

    /* index 0 for cpu mapping; pfunc handles per-cpu delivery */
    bpf_perf_event_output(skb, &perf_events, BPF_F_CURRENT_CPU, &evt, sizeof(evt));
}

/* The main TC program. We'll attach this to both ingress/egress as needed. */
SEC("tc")
int tc_core_prog(struct __sk_buff *skb)
{
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    struct ethhdr *eth;
    void *cursor;
    __u16 eth_proto;
    __u64 offset = 0;

    /* parse eth */
    cursor = parse_eth(skb, data, data_end, &eth);
    if (!cursor)
        return TC_ACT_OK;

    eth_proto = bpf_ntohs(eth->h_proto);
    /* skip VLANs (simple) */
    if (eth_proto == ETH_P_8021Q || eth_proto == ETH_P_8021AD) {
        /* 802.1Q: 4 bytes VLAN tag */
        if ((void *)cursor + 4 > data_end)
            return TC_ACT_OK;
        eth_proto = bpf_ntohs(*(__u16 *)((__u8 *)cursor + 2));
        cursor = (void *)cursor + 4;
    }

    /* handle IPv4 only for now */
    if (eth_proto != ETH_P_IP)
        return TC_ACT_OK;

    struct iphdr *iph = parse_ipv4(cursor, data_end);
    if (!iph)
        return TC_ACT_OK;

    /* compute ip header length */
    __u32 ihl = iph->ihl * 4;
    if (ihl < sizeof(struct iphdr))
        return TC_ACT_OK;

    if ((void *)iph + ihl > data_end)
        return TC_ACT_OK;

    /* only handle UDP for GTP-U */
    if (iph->protocol != IPPROTO_UDP) {
        /* non-UDP: could come from UE towards SGi; do UE-IP map lookup */
        struct ue_ip_key ukey = {};
        ukey.ip4 = iph->saddr; /* network order */
        struct session_value *sval = bpf_map_lookup_elem(&ue_ip_map, &ukey);
        if (sval) {
            /* session found for UE IP: take action */
            update_global_counters(skb->len);
            emit_perf_event(skb, EVT_SESSION_LOOKUP_HIT, skb->ifindex, 0, iph->saddr, iph->daddr);

            if (sval->action == ACTION_DROP) {
                emit_perf_event(skb, EVT_DROPPED_BY_POLICY, skb->ifindex, 0, iph->saddr, iph->daddr);
                return TC_ACT_SHOT;
            } else if (sval->action == ACTION_REDIRECT_IFINDEX) {
                /* attempt a redirect via devmap (preferred) or clone redirect */
                /* try devmap redirect first: key is ifindex */
                __u32 key = sval->ifindex;
                /* bpf_redirect_map is preferable, but here we call bpf_redirect_map via dev_map helper by writing to dev_map */
                /* dev_map will be configured in user-space to map indices to devices; use bpf_redirect_map helper (libbpf provides wrapper) */
                int ret = bpf_redirect_map(&dev_map, key, 0);
                /* if redirect helper returned non-zero, verify: TC expects us to return action codes */
                if (ret == 0)
                    return TC_ACT_OK; /* redirected */
                else
                    return TC_ACT_OK;
            } else {
                /* count only */
                return TC_ACT_OK;
            }
        } else {
            /* not found: let packet pass */
            emit_perf_event(skb, EVT_SESSION_LOOKUP_MISS, skb->ifindex, 0, iph->saddr, iph->daddr);
            return TC_ACT_OK;
        }
    }

    /* For UDP, parse udp header */
    struct udphdr *udph = (void *)iph + ihl;
    if ((void *)(udph + 1) > data_end)
        return TC_ACT_OK;

    __u16 udp_dport = bpf_ntohs(udph->dest);
    __u16 udp_sport = bpf_ntohs(udph->source);

    /* GTP-U typically uses dest port 2152 */
    if (udp_dport != GTPU_UDP_PORT && udp_sport != GTPU_UDP_PORT)
        return TC_ACT_OK;

    /* locate GTP-U header */
    struct gtpu_hdr *g = (void *)(udph + 1);
    if (!g || (void *)(g + 1) > data_end)
        return TC_ACT_OK;

    /* read TEID (network order in header) */
    __u32 teid = bpf_ntohl(g->teid);
    struct teid_key tkey = { .teid = teid };
    struct session_value *s = bpf_map_lookup_elem(&teid_map, &tkey);
    if (!s) {
        /* missing session: count miss and pass */
        emit_perf_event(skb, EVT_SESSION_LOOKUP_MISS, skb->ifindex, teid, iph->saddr, iph->daddr);
        return TC_ACT_OK;
    }

    /* session hit: update counters & take action */
    update_global_counters(skb->len);
    emit_perf_event(skb, EVT_SESSION_LOOKUP_HIT, skb->ifindex, teid, iph->saddr, iph->daddr);

    if (s->action == ACTION_DROP) {
        emit_perf_event(skb, EVT_DROPPED_BY_POLICY, skb->ifindex, teid, iph->saddr, iph->daddr);
        return TC_ACT_SHOT;
    } else if (s->action == ACTION_REDIRECT_IFINDEX) {
        __u32 dst_if = s->ifindex;
        /* Try redirect using devmap if populated */
        int r = bpf_redirect_map(&dev_map, dst_if, 0);
        /* bpf_redirect_map returns TC_ACT_* semantics; if that isn't available, fallback to clone redirect */
        if (r == TC_ACT_OK || r == TC_ACT_SHOT)
            return r;
        /* fallback: clone/redirect to ifindex */
        /* NOTE: clone redirect returns 0 on success and we must return TC_ACT_SHOT or OK depending on policy.
         * Use bpf_clone_redirect when preserving original expected behaviour is needed.
         */
        bpf_clone_redirect(skb, dst_if, 0);
        return TC_ACT_OK;
    } else {
        /* ACTION_COUNT_ONLY or unknown: just pass */
        return TC_ACT_OK;
    }

    return TC_ACT_OK;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
