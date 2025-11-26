/* SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause */
/*
 * tc_core.h
 *
 * Shared definitions between kernel eBPF program (tc_core.bpf.c) and user-space
 * components (bpf_loader, ebpf_manager). Keep layout stable: structs are packed.
 *
 * Include from kernel code: use <linux/types.h> or <bpf/bpf_helpers.h>.
 * Include from user-space: include <stdint.h>.
 */

#ifndef __TC_CORE_H__
#define __TC_CORE_H__

/* Prefer fixed-width types for consistent layout across user/kern */
#include <stdint.h>

/* Map name macros (useful for pinning and looking up maps from user-space) */
#define TEID_MAP_NAME      "teid_map"
#define UE_IP_MAP_NAME     "ue_ip_map"
#define DEV_MAP_NAME       "dev_map"
#define COUNTERS_MAP_NAME  "counters_map"
#define PERF_EVENTS_NAME   "perf_events"

/* TC action return codes (same semantics as linux/tc_act.h) */
#define TC_ACT_OK    0
#define TC_ACT_SHOT  2

/* GTP-U port */
#define GTPU_UDP_PORT 2152

/* Event types emitted to perf ring (small enum; user-space may interpret) */
enum event_type {
    EVT_SESSION_LOOKUP_HIT = 1,
    EVT_SESSION_LOOKUP_MISS = 2,
    EVT_DROPPED_BY_POLICY  = 3,
};

/* Session action type decisions (simple policy) */
enum session_action {
    ACTION_COUNT_ONLY       = 0,
    ACTION_REDIRECT_IFINDEX = 1,
    ACTION_DROP             = 2,
};

/* TEID map key (4-byte TEID) */
struct teid_key {
    uint32_t teid; /* host or network order? stored/used in network order in BPF code */
} __attribute__((packed));

/* UE IP map key (IPv4) */
struct ue_ip_key {
    uint32_t ip4; /* IPv4 address in network byte order */
} __attribute__((packed));

/* Per-session value stored in maps */
struct session_value {
    uint32_t ifindex;    /* target ifindex for redirect (action==REDIRECT_IFINDEX) */
    uint32_t rx_packets; /* optional shadow counter (user-space may update) */
    uint64_t rx_bytes;   /* optional shadow bytes */
    uint32_t action;     /* enum session_action */
    uint32_t flags;      /* misc flags (reserved) */
} __attribute__((packed));

/* Small global/summmary counter key/value */
struct counter_key {
    uint32_t id;
} __attribute__((packed));

struct counter_value {
    uint64_t packets;
    uint64_t bytes;
} __attribute__((packed));

/* Perf event record (written to perf ring) */
struct perf_event {
    uint32_t type;     /* event_type */
    uint32_t reason;   /* optional reason code or sub-type */
    uint32_t ifindex;  /* ingress or redirect ifindex */
    uint32_t pad;
    uint64_t pkt_len;  /* skb->len */
    uint32_t teid;     /* TEID for GTP events (host-order) */
    uint32_t src_ip;   /* IPv4 src (network order) */
    uint32_t dst_ip;   /* IPv4 dst (network order) */
} __attribute__((packed));


#endif /* __TC_CORE_H__ */
