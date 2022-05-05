// SPDX-License-Identifier: GPL-2.0
/*
 * Copyright 2021 The Magma Authors.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of version 2 of the GNU General Public
 * License as published by the Free Software Foundation.
 */

#include "orc8r/gateway/c/common/ebpf/EbpfMap.h"

#include <bcc/proto.h>
#include <linux/if_packet.h>
#include <linux/ip.h>
#include <linux/socket.h>
#include <linux/pkt_cls.h>
#include <linux/erspan.h>
#include <linux/udp.h>

// The map is pinned so that it can be accessed by pipelined or debugging tool
// to examine datapath state.
BPF_TABLE_PINNED("hash", struct ul_map_key, struct ul_map_info, ul_map,
                 1024 * 512, "/sys/fs/bpf/ul_map");

// Ingress handler for Uplink traffic.
int gtpu_ingress_handler(struct __sk_buff* skb) {
  int ret;
  struct iphdr* iph;
  void* data;
  void* data_end;

  int gtp_offset =
      sizeof(struct ethhdr) + sizeof(struct iphdr) + sizeof(struct udphdr);
  int inner_ip_hdr = sizeof(struct ethhdr) + sizeof(struct iphdr) +
                     sizeof(struct udphdr) + gtp_hdr_size +
                     sizeof(struct iphdr);
  // 1. a. check GTP HDR
  data = (void*)(long)skb->data;
  data_end = (void*)(long)skb->data_end;
  int len = data_end - data;

  if ((data + inner_ip_hdr) > data_end) {
    bpf_trace_printk("ERR: truncated packet: len: %d, data sz %d\n", skb->len,
                     len);
    return TC_ACT_OK;
  }
  struct gtp1_header* gtp1 = (struct gtp1_header*)(data + gtp_offset);
  int tid = ntohl(gtp1->tid);

  struct udphdr* uh = data + sizeof(struct ethhdr) + sizeof(struct iphdr);
  if (uh->dest != htons(GTP_PORT_NO)) {
    // not GTP packet, let it continue.
    bpf_trace_printk("ERR: Not GTP packet: %d\n", ntohs(uh->dest));
    return TC_ACT_OK;
  }
  // 1. b. check UE map
  iph = data + gtp_offset + gtp_hdr_size;
  struct ul_map_key key = {iph->saddr};
  struct ul_map_info* fwd = ul_map.lookup(&key);
  if (!fwd) {
    bpf_trace_printk("ERR: UE for IP %x not found\n", iph->saddr);
    // No UE entry.
    return TC_ACT_OK;
  }

  // 2. process inner packet.
  iph = NULL;
  int rem_hdrs = sizeof(struct iphdr) + sizeof(struct udphdr) + gtp_hdr_size;
  ret = bpf_skb_adjust_room(skb, -rem_hdrs, BPF_ADJ_ROOM_MAC, 0);
  if (ret) {
    bpf_trace_printk("ERR: get: error adjust %d proto: %x, offset %d\n", ret,
                     skb->protocol, skb->len);
    return TC_ACT_OK;
  }
  data = (void*)(long)skb->data;
  data_end = (void*)(long)skb->data_end;
  len = data_end - data;

  if ((data + sizeof(struct ethhdr)) > data_end) {
    bpf_trace_printk("ERR: truncated inner packet: len: %d, data sz %d\n",
                     skb->len, len);
    return TC_ACT_OK;
  }

  struct ethhdr* eth = data;

  // 2.1. Update dest mac
  __builtin_memcpy(eth->h_dest, fwd->mac_dst, ETH_ALEN);
  // 2.2. Update src  mac
  __builtin_memcpy(eth->h_source, fwd->mac_src, ETH_ALEN);

  // 2.3 skb mark for qos
  skb->mark = fwd->mark;

  // 2.4 increment bytes
  // TODO: Add lock for accessing bytes
  fwd->bytes += len;

  bpf_trace_printk("INFO: UL-fwd: tid %d egress: %d mark: %d\n", tid,
                   fwd->e_if_index, fwd->mark);
  return bpf_redirect(fwd->e_if_index, 0);
}
