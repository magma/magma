/**
 * Copyright 2020 The Magma Authors.
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

/**
 * This program attaches to the packet transmit event net_dev_start_xmit. It
 * counts the number of bytes and packets sent by each linux binary.
 * Expect considerable performance impact.
 * https://patchwork.ozlabs.org/patch/309437/
 */

// Create a hash map with `key_t` as key type and u64 as value type
BPF_HASH(dest_counters, struct key_t);

// Attach kprobe for the `tcp_sendmsg` syscall
int kprobe__tcp_sendmsg(struct pt_regs *ctx, struct sock *sk,
                        struct msghdr *msg, size_t size) {
  u16 family = sk->sk_family;

  // both IPv4 and IPv6
  if (family == AF_INET || family == AF_INET6) {
    struct key_t key = {};

    // read destination IP address and port
    if (family == AF_INET) {
      key.daddr = sk->sk_daddr;
    } else if (family == AF_INET6) {
      bpf_probe_read(&key.daddr, sizeof(key.daddr), &sk->sk_v6_daddr.s6_addr32);
    }
    u16 dport = sk->sk_dport;
    key.dport = ntohs(dport);

    // ignore packets destined for localhost unless dport is PROXY_PORT

    // IPv4 127.0.0.1 == 16777343, 0.0.0.0 == 0
    if (family == AF_INET && key.dport != {{PROXY_PORT}}) {
      if (key.daddr == 0x100007F || key.daddr == 0) {
        return 0;
      }
    }
    // IPv6 localhost has embedded IPv4 in case of control_proxy
    // uint128 is split into two u64 and high is FFFF:"127.0.0.1"
    else if (family == AF_INET6 && key.dport != {{PROXY_PORT}}) {
      uint64_t low = (uint64_t)key.daddr;
      uint64_t high = (key.daddr >> 64);
      if (low == 0 && high == 0x100007FFFFF0000) {
        return 0;
      }
    }

    // read binary name, pid, and IP version
    bpf_get_current_comm(key.comm, TASK_COMM_LEN);
    key.pid = bpf_get_current_pid_tgid() >> 32;
    key.family = family;

    // increment the counters
    dest_counters.increment(key, size);
  }
  return 0;
}
