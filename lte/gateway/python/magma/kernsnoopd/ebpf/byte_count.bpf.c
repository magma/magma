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
int kprobe__tcp_sendmsg(struct pt_regs *ctx, struct sock *sk, struct msghdr *msg, size_t size) {
  u16 dport = 0, family = sk->__sk_common.skc_family;

  // only IPv4
  if (family == AF_INET) {
    struct key_t key;
    key.daddr = sk->__sk_common.skc_daddr;

    // ignore packets destined for localhost
    // IPv4 127.0.0.1 == 16777343, 0.0.0.0 == 0
    if ((key.daddr == 16777343 || key.daddr == 0)) {
      return 0;
    }

    // read binary name, pid, source and destination IP address and port
    bpf_get_current_comm(key.comm, TASK_COMM_LEN);
    key.pid = bpf_get_current_pid_tgid() >> 32;;
    key.saddr = sk->__sk_common.skc_rcv_saddr;
    key.lport = sk->__sk_common.skc_num;
    dport = sk->__sk_common.skc_dport;
    key.dport = ntohs(dport);

    // increment the counters
    dest_counters.increment(key, size);
  }
  return 0;
}
