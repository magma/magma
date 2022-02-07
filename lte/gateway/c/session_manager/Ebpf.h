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

#include <unistd.h>
#include <stdint.h>
#include <linux/bpf.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include "lte/protos/session_manager.pb.h"

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

namespace magma {
using namespace lte;
    int bpf_obj_get(const char* pathname);
    int bpf_map_update_elem(int fd, void*   key, void* value, uint64_t flags);
    int bpf_map_lookup_elem(int fd, void* key, void* value);
    int bpf_map_get_next_key(int fd, void *key, void *next_key);
    int bpf_map_lookup_elem(int fd, void *key, void *value);
    RuleRecordTable GetEbpfTable();
}
