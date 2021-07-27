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

// Create a hash map with `key_t` as key type and `counters_t` as value type
BPF_HASH(dest_counters, struct key_t, struct counters_t, 1000);

// Attach hook for the `net_dev_start_xmit` kernel trace event
TRACEPOINT_PROBE(net, net_dev_start_xmit) {
  struct sk_buff* skb = (struct sk_buff*) args->skbaddr;
  struct sock* sk     = skb->sk;

  struct key_t key = {};
  bpf_probe_read(&key.daddr, sizeof(sk->sk_daddr), &sk->sk_daddr);
  // ignore packets destined for localhost
  // IPv4 127.0.0.1 == 16777343, 0.0.0.0 == 0
  if ((key.daddr == 16777343 || key.daddr == 0)) {
    return 0;
  }

  // read binary name, pid, destination IP address, port
  bpf_get_current_comm(key.comm, TASK_COMM_LEN);
  key.pid = bpf_get_current_pid_tgid() >> 32;
  bpf_probe_read(&key.dport, sizeof(sk->sk_dport), &sk->sk_dport);

  // lookup or initialize item in dest_counters
  struct counters_t empty;
  __builtin_memset(&empty, 0, sizeof(empty));
  struct counters_t* data = dest_counters.lookup_or_try_init(&key, &empty);

  // increment the counters
  if (data) {
    data->bytes += skb->len;
    data->packets++;
  }
  return 0;
}
