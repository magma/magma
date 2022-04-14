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

#pragma once

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>
#include <strings.h>
#include <sys/types.h>
#include <unistd.h>
#include <linux/bpf.h>

#include "orc8r/gateway/c/common/ebpf/EbpfMap.h"

#ifndef __NR_bpf
#if defined(__i386__)
#define __NR_bpf 357
#elif defined(__x86_64__)
#define __NR_bpf 321
#elif defined(__aarch64__)
#define __NR_bpf 280
#elif defined(__sparc__)
#define __NR_bpf 349
#elif defined(__s390__)
#define __NR_bpf 351
#else
#error __NR_bpf not defined. libbpf does not support your arch.
#endif
#endif

#define BPF_ANY 0     /* create new element or update existing */
#define BPF_NOEXIST 1 /* create new element only if it didn't exist */
#define BPF_EXIST 2   /* only update existing element */

#define DL_MAP_PATH "/sys/fs/bpf/dl_map"

static uint64_t ptr_to_u64(const void* ptr) { return (uint64_t)ptr; }

static inline int sys_bpf(enum bpf_cmd cmd, union bpf_attr* attr, uint size) {
  return syscall(__NR_bpf, cmd, attr, size);
}

int bpf_obj_get(const char* pathname);
int bpf_map_update_elem(int fd, void* key, void* value, uint64_t flags);
int bpf_map_lookup_elem(int fd, void* key, void* value);

int bpf_obj_get(const char* pathname) {
  union bpf_attr attr;

  bzero(&attr, sizeof(attr));
  attr.pathname = ptr_to_u64((const void*)pathname);

  return sys_bpf(BPF_OBJ_GET, &attr, sizeof(attr));
}

int bpf_map_update_elem(int fd, void* key, void* value, uint64_t flags) {
  union bpf_attr attr;

  bzero(&attr, sizeof(attr));
  attr.map_fd = fd;
  attr.key = ptr_to_u64(key);
  attr.value = ptr_to_u64(value);
  attr.flags = flags;

  return sys_bpf(BPF_MAP_UPDATE_ELEM, &attr, sizeof(attr));
}

int bpf_map_delete_elem(int fd, void* key) {
  union bpf_attr attr;

  memset(&attr, 0, sizeof(attr));
  attr.map_fd = fd;
  attr.key = ptr_to_u64(key);

  return sys_bpf(BPF_MAP_DELETE_ELEM, &attr, sizeof(attr));
}

int get_map_fd() { return bpf_obj_get(DL_MAP_PATH); }

int add_ebpf_dl_map_entry(int hash_fd, struct in_addr ue, struct in_addr enb,
                          uint32_t o_tei, uint8_t* user_data,
                          size_t user_data_len) {
  struct dl_map_info val = {htonl(enb.s_addr), o_tei, 0, {}};
  if (user_data_len > sizeof(val.user_data)) {
    return -1;
  }
  memcpy(val.user_data, user_data, user_data_len);
  uint32_t nkey = htonl(ue.s_addr);
  return bpf_map_update_elem(hash_fd, &nkey, &val, 0);
}

void delete_ebpf_dl_map_entry(int hash_fd, struct in_addr ue) {
  uint32_t nkey = htonl(ue.s_addr);
  bpf_map_delete_elem(hash_fd, &nkey);
}
