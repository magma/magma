

#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#include "tc_core.h"


struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 16384);
    __type(key, struct teid_key);
    __type(value, struct session_value);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} teid_map SEC(".maps");




struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 16384);
    __type(key, struct ue_ip_key);
    __type(value, struct session_value);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} ue_ip_map SEC(".maps");




struct {
    __uint(type, BPF_MAP_TYPE_DEVMAP);
    __uint(max_entries, 256);
    __uint(key_size, sizeof(uint32_t));
    __uint(value_size, sizeof(uint32_t));
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} dev_map SEC(".maps");




struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 256);
    __type(key, struct counter_key);
    __type(value, struct counter_value);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} counters_map SEC(".maps");




struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(uint32_t));
    __uint(value_size, sizeof(uint32_t));
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} perf_events SEC(".maps");

