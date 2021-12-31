// SPDX-License-Identifier: GPL-2.0
/*
 * Copyright 2021 The Magma Authors.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of version 2 of the GNU General Public
 * License as published by the Free Software Foundation.
 */
#include <bcc/proto.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <linux/if_ether.h>
#include <uapi/linux/ipv6.h>

// UE sessions map definitions.
struct dl_map_key {
  u32 ue_ip;
};

struct dl_map_info {
  u32 remote_ipv4;
  u32 tunnel_id;
};

// The map is pinned so that it can be accessed by pipelined or debugging tool
// to examine datapath state.
BPF_TABLE_PINNED(
    "hash", struct dl_map_key, struct dl_map_info, dl_map, 1024 * 512,
    "/sys/fs/bpf/dl_map");

// Ingress handler for Uplink traffic.
int gtpu_egress_handler(struct __sk_buff* skb) {
  int ret;
  struct iphdr* iph;
  void* data;
  void* data_end;

  int ip_hdr = sizeof(struct ethhdr) + sizeof(struct iphdr);

  data     = (void*) (long) skb->data;
  data_end = (void*) (long) skb->data_end;
  int len  = data_end - data;

  // 1. a. check packet len
  if ((data + ip_hdr) > data_end) {
    bpf_trace_printk(
        "ERR: truncated packet: len: %d, data sz %d\n", skb->len, len);
    return TC_ACT_OK;
  }

  // 1. b. check UE map
  iph                     = data + sizeof(struct ethhdr);
  struct dl_map_key key   = {iph->daddr};
  struct dl_map_info* fwd = dl_map.lookup(&key);
  if (!fwd) {
    bpf_trace_printk("ERR: UE for IP %x not found\n", iph->daddr);
    return TC_ACT_OK;
  }

  // 2. set tunnel info
  struct bpf_tunnel_key tun_key;
  __builtin_memset(&key, 0x0, sizeof(tun_key));
  tun_key.remote_ipv4 = fwd->remote_ipv4;
  tun_key.tunnel_id = fwd->tunnel_id;
  tun_key.tunnel_tos = 0;
  tun_key.tunnel_ttl = 64;

  ret = bpf_skb_set_tunnel_key(skb, &key, sizeof(key),
                               BPF_F_ZERO_CSUM_TX | BPF_F_SEQ_NUMBER);
  if (ret < 0) {
      bpf_trace_printk("ERR: bpf_skb_set_tunnel_key failed with %d", ret);
      return TC_ACT_SHOT;
  }
  bpf_trace_printk("INFO: set: key %d remote ip 0x%x ret = %d\n", key.tunnel_id, key.remote_ipv4, ret)

  //TODO save in map
  int if_idx = 51;

  return bpf_redirect(if_idx, 0);
}
