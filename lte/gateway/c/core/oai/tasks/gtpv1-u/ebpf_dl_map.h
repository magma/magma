/**
 * Copyright 2022 The Magma Authors.
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

/* File : ebpf_dl_map.h
 */

#pragma once

#include "lte/gateway/c/core/oai/tasks/gtpv1-u/ebpf.h"

#define DL_MAP_PATH "/sys/fs/bpf/dl_map"

struct bpf_map_val {
  uint32_t ip;
  uint32_t tei;
};

int get_map_fd() {
  return bpf_obj_get(DL_MAP_PATH);
}

void add_ebpf_dl_map_entry(
    int hash_fd, struct in_addr ue, struct in_addr enb, uint32_t o_tei) {
  struct bpf_map_val val = {htonl(enb.s_addr), o_tei};
  uint32_t nkey          = htonl(ue.s_addr);
  bpf_map_update_elem(hash_fd, &nkey, &val, 0);
}

void delete_ebpf_dl_map_entry(int hash_fd, struct in_addr ue) {
  uint32_t nkey = htonl(ue.s_addr);
  bpf_map_delete_elem(hash_fd, &nkey);
}
