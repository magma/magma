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

/*! \file sgw_lite.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_SGW_SEEN
#define FILE_SGW_SEEN
#include <stdint.h>
#include <netinet/in.h>

#include "bstrlib.h"
#include "queue.h"
#include "hashtable.h"

#include "nas/commonDef.h"
#include "common_types.h"
#include "sgw_context_manager.h"
#include "gtpv1u_sgw_defs.h"
#include "pgw_pcef_emulation.h"

typedef struct sgw_app_s {
  bstring sgw_if_name_S1u_S12_S4_up;
  struct in_addr sgw_ip_address_S1u_S12_S4_up;

  bstring sgw_if_name_S11_S4;
  struct in_addr sgw_ip_address_S11_S4;

  struct in_addr sgw_ip_address_S5_S8_up;  // unused now

  // key is S11 S-GW local teid
  hash_table_ts_t* s11teid2mme_hashtable;

  // the key of this hashtable is the S11 s-gw local teid.
  hash_table_ts_t* s11_bearer_context_information_hashtable;

  gtpv1u_data_t gtpv1u_data;
} sgw_app_t;

struct ipv4_list_elm_s {
  STAILQ_ENTRY(ipv4_list_elm_s) ipv4_entries;
  struct in_addr addr;
};

typedef struct pgw_app_s {
  STAILQ_HEAD(ipv4_list_free_head_s, ipv4_list_elm_s) ipv4_list_free;
  STAILQ_HEAD(ipv4_list_allocated_head_s, ipv4_list_elm_s) ipv4_list_allocated;
  hash_table_ts_t* deactivated_predefined_pcc_rules;
  hash_table_ts_t* predefined_pcc_rules;
} pgw_app_t;

#endif
