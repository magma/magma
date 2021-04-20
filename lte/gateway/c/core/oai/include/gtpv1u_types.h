/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include "common_types.h"
#include "hashtable.h"

#define GTPV1U_UDP_PORT (2152)
#define MAX_BEARERS_PER_UE (11)

#define BUFFER_TO_uint32_t(buf, x)                                             \
  do {                                                                         \
    x = ((uint32_t)((buf)[0])) | ((uint32_t)((buf)[1]) << 8) |                 \
        ((uint32_t)((buf)[2]) << 16) | ((uint32_t)((buf)[3]) << 24);           \
  } while (0)

typedef enum {
  BEARER_DOWN = 0,
  BEARER_IN_CONFIG,
  BEARER_UP,
  BEARER_DL_HANDOVER,
  BEARER_UL_HANDOVER,
  BEARER_MAX,
} s1_bearer_state_t;

typedef struct {
  /* RB tree of UEs */
  // RB_HEAD(gtpv1u_ue_map, gtpv1u_ue_data_s) gtpv1u_ue_map_head;
  /* Local IP address to use */
  struct in_addr sgw_ip_address_for_S1u_S12_S4_up;
  hash_table_t* S1U_mapping;

  // GTP-U kernel interface
  pthread_t reader_thread;
  int fd0;  /* GTP0 file descriptor */
  int fd1u; /* GTP1-U user plane file descriptor */
} gtpv1u_data_t;
