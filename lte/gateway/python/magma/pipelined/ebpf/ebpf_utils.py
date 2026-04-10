"""
Copyright 2025 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Author: Nitin Rajput (coRAN LABS)

eBPF utilities and constants for Magma pipelined
"""

# TC BPF attachment types
INGRESS = "ingress"
EGRESS = "egress"

# TC actions
TC_ACT_UNSPEC = -1
TC_ACT_OK = 0
TC_ACT_RECLASSIFY = 1
TC_ACT_SHOT = 2
TC_ACT_PIPE = 3
TC_ACT_STOLEN = 4
TC_ACT_QUEUED = 5
TC_ACT_REPEAT = 6
TC_ACT_REDIRECT = 7

# BPF program types
BPF_PROG_TYPE_SCHED_CLS = 3
BPF_PROG_TYPE_XDP = 6

# XDP actions
XDP_ABORTED = 0
XDP_DROP = 1
XDP_PASS = 2
XDP_TX = 3
XDP_REDIRECT = 4

# Interface constants
MAX_IFACE_NUM = 512
IFACE_NAME_MAX_LEN = 16

# Map types
BPF_MAP_TYPE_HASH = 1
BPF_MAP_TYPE_ARRAY = 2
BPF_MAP_TYPE_PROG_ARRAY = 3
BPF_MAP_TYPE_PERF_EVENT_ARRAY = 4
BPF_MAP_TYPE_PERCPU_HASH = 5
BPF_MAP_TYPE_PERCPU_ARRAY = 6
BPF_MAP_TYPE_STACK_TRACE = 7
BPF_MAP_TYPE_CGROUP_ARRAY = 8
BPF_MAP_TYPE_LRU_HASH = 9
BPF_MAP_TYPE_LRU_PERCPU_HASH = 10
BPF_MAP_TYPE_LPM_TRIE = 11
BPF_MAP_TYPE_ARRAY_OF_MAPS = 12
BPF_MAP_TYPE_HASH_OF_MAPS = 13
BPF_MAP_TYPE_DEVMAP = 14
BPF_MAP_TYPE_SOCKMAP = 15
BPF_MAP_TYPE_CPUMAP = 16
BPF_MAP_TYPE_XSKMAP = 17
BPF_MAP_TYPE_SOCKHASH = 18

# Helper functions for working with eBPF programs
def get_tc_direction(direction):
    """Convert direction string to TC direction constant"""
    if direction.lower() == "ingress":
        return INGRESS
    elif direction.lower() == "egress":
        return EGRESS
    else:
        raise ValueError(f"Invalid direction: {direction}")

def is_valid_tc_action(action):
    """Check if TC action code is valid"""
    valid_actions = [
        TC_ACT_OK, TC_ACT_SHOT, TC_ACT_PIPE, 
        TC_ACT_REDIRECT, TC_ACT_STOLEN
    ]
    return action in valid_actions