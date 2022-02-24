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

#include <linux/if_ether.h>
#include <linux/types.h>

// UE sessions map definitions.
struct ul_map_key {
  __u32 ue_ip;
};

struct ul_map_info {
  __u32 mark;
  __u32 e_if_index;
  uint64_t bytes;
  __u8 mac_src[ETH_ALEN];
  __u8 mac_dst[ETH_ALEN];
};

struct dl_map_key {
  __u32 ue_ip;
};

struct dl_map_info {
  __u32 remote_ipv4;
  __u32 tunnel_id;
  uint64_t bytes;
  __u8 user_data[64];
};

// GTP protocol definitions
struct gtp1_header { /* According to 3GPP TS 29.060. */
  __u8 flags;
  __u8 type;
  __be16 length;
  __be32 tid;
} __attribute__((packed));

static const int GTP_PORT_NO = 2152;
static const int gtp_hdr_size = 8;
