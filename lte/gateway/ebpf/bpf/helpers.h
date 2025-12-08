// helpers.h
#ifndef HELPERS_H
#define HELPERS_H

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <linux/tcp.h>

// Logging macro using bpf_trace_printk
#define bpf_log(fmt, ...) \
    ({ char ____fmt[] = fmt; bpf_printk(____fmt, ##__VA_ARGS__); })

// Compute simple checksum over data
static __always_inline __maybe_unused __u16 csum16(const void *data, __u32 len) {
    const __u16 *buf = data;
    __u32 sum = 0;

    for (__u32 i = 0; i < (len / 2); i++) {
        sum += buf[i];
    }

    if (len & 1) {
        sum += *(__u8 *)((__u8 *)data + len - 1);
    }

    // Fold 32-bit sum to 16 bits
    while (sum >> 16)
        sum = (sum & 0xFFFF) + (sum >> 16);

    return ~sum;
}

// Helper to parse Ethernet header
static __always_inline __maybe_unused int parse_ethhdr(void *data, void *data_end, struct ethhdr **ethh) {
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return -1;
    *ethh = eth;
    return 0;
}

// Helper to parse IPv4 header
static __always_inline __maybe_unused int parse_iphdr(void *data, void *data_end, struct iphdr **iph) {
    struct iphdr *ip = data;
    if ((void *)(ip + 1) > data_end)
        return -1;
    *iph = ip;
    return 0;
}

// Helper to parse UDP header
static __always_inline __maybe_unused int parse_udphdr(void *data, void *data_end, struct udphdr **udph) {
    struct udphdr *udp = data;
    if ((void *)(udp + 1) > data_end)
        return -1;
    *udph = udp;
    return 0;
}

// Helper to parse TCP header
static __always_inline __maybe_unused int parse_tcphdr(void *data, void *data_end, struct tcphdr **tcph) {
    struct tcphdr *tcp = data;
    if ((void *)(tcp + 1) > data_end)
        return -1;
    *tcph = tcp;
    return 0;
}

#endif // HELPERS_H
