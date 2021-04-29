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

/*! \file mme_app_edns_emulation.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdlib.h>
#include <netinet/in.h>
#include <string.h>

#include "bstrlib.h"
#include "obj_hashtable.h"
#include "mme_config.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "mme_app_edns_emulation.h"
#include "hashtable.h"

static obj_hash_table_t* g_e_dns_entries = NULL;

//------------------------------------------------------------------------------
struct in_addr* mme_app_edns_get_sgw_entry(bstring id) {
  struct in_addr* in_addr = NULL;
  obj_hashtable_get(g_e_dns_entries, bdata(id), blength(id), (void**) &in_addr);

  return in_addr;
}

//------------------------------------------------------------------------------
int mme_app_edns_add_sgw_entry(bstring id, struct in_addr in_addr) {
  char* cid = calloc(1, blength(id) + 1);
  if (cid) {
    strncpy(cid, (const char*) id->data, blength(id));

    struct in_addr* data = malloc(sizeof(struct in_addr));
    if (data) {
      data->s_addr = in_addr.s_addr;

      hashtable_rc_t rc =
          obj_hashtable_insert(g_e_dns_entries, cid, strlen(cid), data);
      if (HASH_TABLE_OK == rc) {
        free(cid);
        return RETURNok;
      }
    }
    free(cid);
  }
  return RETURNerror;
}

//------------------------------------------------------------------------------
int mme_app_edns_init(const mme_config_t* mme_config_p) {
  int rc          = RETURNok;
  g_e_dns_entries = obj_hashtable_create(
      OAI_MIN(32, MME_CONFIG_MAX_SGW), NULL, free_wrapper, free_wrapper, NULL);
  if (g_e_dns_entries) {
    for (int i = 0; i < mme_config_p->e_dns_emulation.nb_sgw_entries; i++) {
      rc |= mme_app_edns_add_sgw_entry(
          mme_config_p->e_dns_emulation.sgw_id[i],
          mme_config_p->e_dns_emulation.sgw_ip_addr[i]);
    }
    return rc;
  }
  return RETURNerror;
}
//------------------------------------------------------------------------------
void mme_app_edns_exit(void) {
  obj_hashtable_destroy(g_e_dns_entries);
}
