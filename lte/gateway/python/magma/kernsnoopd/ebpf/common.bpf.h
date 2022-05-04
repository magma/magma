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
 * This is the common header for all eBPF programs. eBPF programs to be loaded
 * into the kernel will be joined together and prefixed with this header.
 * All include statements and common definitions should go here.
 */

#include <bcc/proto.h>
#include <linux/sched.h>
#include <uapi/linux/ptrace.h>
#include <net/sock.h>

struct key_t {
  // binary name (task->comm in the kernel)
  char comm[TASK_COMM_LEN];
  u32 pid;
  // destination IP address and port
  unsigned __int128 daddr;
  u16 dport;
  // IP version
  u16 family;
};
