/**
 * Copyright 2021 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "orc8r/gateway/c/common/ebpf/EbpfMap.h"

#include <bcc/proto.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <uapi/linux/ipv6.h>

// The map is pinned so that it can be accessed by pipelined or debugging tool
// to examine datapath state.
BPF_TABLE_PINNED("hash", struct dl_map_key, struct dl_map_info, dl_map,
                 1024 * 512, "/sys/fs/bpf/dl_map");

struct cfg_array_info {
  u32 if_idx;
};

BPF_TABLE_PINNED("array", u32, struct cfg_array_info, cfg_array, 1,
                 "/sys/fs/bpf/cfg_array");

// Ingress handler for Uplink traffic.
int gtpu_egress_handler(struct __sk_buff* skb) {
  int ret;
  struct iphdr* iph;
  void* data;
  void* data_end;

  int ip_hdr = sizeof(struct ethhdr) + sizeof(struct iphdr);

  // 1. a. check pkt length
  data = (void*)(long)skb->data;
  data_end = (void*)(long)skb->data_end;
  int len = data_end - data;

  if ((data + ip_hdr) > data_end) {
    bpf_trace_printk("ERR: truncated packet: len: %d, data sz %d\n", skb->len,
                     len);
    return TC_ACT_OK;
  }

  // 1. b. check UE map
  iph = data + sizeof(struct ethhdr);
  struct dl_map_key lookup_key = {iph->daddr};
  struct dl_map_info* fwd = dl_map.lookup(&lookup_key);
  if (!fwd) {
    bpf_trace_printk("ERR: UE for IP %x not found\n", iph->daddr);
    return TC_ACT_OK;
  }

  // 2. set tunnel info
  struct bpf_tunnel_key tun_info;
  __builtin_memset(&tun_info, 0x0, sizeof(tun_info));
  tun_info.remote_ipv4 = fwd->remote_ipv4;
  tun_info.tunnel_id = fwd->tunnel_id;
  tun_info.tunnel_ttl = 64;

  ret = bpf_skb_set_tunnel_key(skb, &tun_info, sizeof(tun_info),
                               BPF_F_ZERO_CSUM_TX);
  if (ret < 0) {
    bpf_trace_printk("ERR: bpf_skb_set_tunnel_key failed with %d", ret);
    return TC_ACT_SHOT;
  }
  bpf_trace_printk("INFO: set: key %d remote ip 0x%x ret = %d\n",
                   tun_info.tunnel_id, tun_info.remote_ipv4, ret);

  u32 cfg_key = 0;
  struct cfg_array_info* cfg = cfg_array.lookup(&cfg_key);
  if (!cfg) {
    bpf_trace_printk("ERR: Config array lookup failed\n");
    return TC_ACT_OK;
  }

  // TODO: Add lock for accessing bytes
  fwd->bytes += skb->len;

  return bpf_redirect(cfg->if_idx, 0);
}
