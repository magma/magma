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

/* File : Ebpf.cpp
 */

#include <sys/socket.h>
#include <glog/logging.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <errno.h>

#include "Ebpf.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/session_manager.pb.h"
#include "magma_logging.h"


#define DL_MAP_PATH "/sys/fs/bpf/dl_map"
#define UL_MAP_PATH "/sys/fs/bpf/ul_map"

namespace magma {

///////

static uint64_t ptr_to_u64(const void* ptr) {
  return (uint64_t) ptr;
}

static inline int sys_bpf(enum bpf_cmd cmd, union bpf_attr* attr, uint size) {
  return syscall(__NR_bpf, cmd, attr, size);
}

int bpf_obj_get(const char* pathname) {
  union bpf_attr attr;

  bzero(&attr, sizeof(attr));
  attr.pathname = ptr_to_u64((const void*) pathname);

  return sys_bpf(BPF_OBJ_GET, &attr, sizeof(attr));
}

int bpf_map_update_elem(int fd, void* key, void* value, uint64_t flags) {
  union bpf_attr attr;

  bzero(&attr, sizeof(attr));
  attr.map_fd = fd;
  attr.key    = ptr_to_u64(key);
  attr.value  = ptr_to_u64(value);
  attr.flags  = flags;

  return sys_bpf(BPF_MAP_UPDATE_ELEM, &attr, sizeof(attr));
}

int bpf_map_delete_elem(int fd, void* key) {
  union bpf_attr attr;

  memset(&attr, 0, sizeof(attr));
  attr.map_fd = fd;
  attr.key    = ptr_to_u64(key);

  return sys_bpf(BPF_MAP_DELETE_ELEM, &attr, sizeof(attr));
}

int bpf_map_get_next_key(int fd, void *key, void *next_key) {
	union bpf_attr attr;

	bzero(&attr, sizeof(attr));
	attr.map_fd = fd;
	attr.key = ptr_to_u64(key);
	attr.next_key = ptr_to_u64(next_key);

	return sys_bpf(BPF_MAP_GET_NEXT_KEY, &attr, sizeof(attr));
}

int bpf_map_lookup_elem(int fd, void *key, void *value)
{
	union bpf_attr attr;

	bzero(&attr, sizeof(attr));
	attr.map_fd = fd;
	attr.key = ptr_to_u64(key);
	attr.value = ptr_to_u64(value);

	return sys_bpf(BPF_MAP_LOOKUP_ELEM, &attr, sizeof(attr));
}

///////


struct bpf_map_val {
  uint32_t ip;
  uint32_t tei;
  char imsi[16];
  uint64_t bytes;
};

int get_dl_map_fd() {
  return bpf_obj_get(DL_MAP_PATH);
}

int get_ul_map_fd() {
  return bpf_obj_get(UL_MAP_PATH);
}
//
//void delete_ebpf_dl_map_entry(int hash_fd, struct in_addr ue) {
//  uint32_t nkey = htonl(ue.s_addr);
//  bpf_map_delete_elem(hash_fd, &nkey);
//}

magma::RuleRecordTable GetEbpfTable() {
    magma::RuleRecordTable records;
    int fd = get_dl_map_fd();
    MLOG(MERROR) << "Got fd" << fd;
    //int ul_map_fd = get_ul_map_fd();
    uint32_t key, prev_key;
    struct bpf_map_val val;
    int res;
    prev_key=-1;
    while(bpf_map_get_next_key(fd, &prev_key, &key) == 0) {
        MLOG(MERROR) << "prev key  " << prev_key;
        MLOG(MERROR) << "Got key " << key;
        RuleRecord record;
        res = bpf_map_lookup_elem(fd, &key, &val);
        if (res < 0) {
            MLOG(MERROR) << "No value??\n";
        } else {
            //MLOG(MERROR) <<"%lld\n"<< val;
            record.set_sid(val.imsi);
            // convert uint to string
            struct in_addr ip;
            ip.s_addr = key;
            record.set_ue_ipv4(inet_ntoa(ip));
            record.set_bytes_tx(val.bytes);
            record.set_teid(val.tei);
            //record.set_bytes_rx(ap_name);
        }
        prev_key=key;

        records.mutable_records()->Add()->CopyFrom(record);
    }
    return records;
}
}
