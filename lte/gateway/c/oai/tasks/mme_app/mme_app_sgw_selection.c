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

/*! \file mme_app_sgw_selection.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdint.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "mme_config.h"
#include "bstrlib.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "TrackingAreaIdentity.h"
#include "mme_app_sgw_selection.h"
#include "mme_app_edns_emulation.h"

//------------------------------------------------------------------------------
void mme_app_select_sgw(
    const tai_t* const tai, struct sockaddr* const sgw_in_addr) {
  extern mme_config_t mme_config;

  ((struct sockaddr_in*) sgw_in_addr)->sin_addr.s_addr =
      mme_config.e_dns_emulation.sgw_ip_addr[0].s_addr;
  ((struct sockaddr_in*) sgw_in_addr)->sin_family = AF_INET;

  printf("Received SGW IP Address %p\n", (void*) sgw_in_addr);
  OAILOG_DEBUG(
      LOG_MME_APP, "SGW  returned %s\n",
      inet_ntoa(((struct sockaddr_in*) sgw_in_addr)->sin_addr));
  return;

  OAILOG_WARNING(
      LOG_MME_APP, "Failed SGW lookup for TAI " TAI_FMT "\n", TAI_ARG(tai));
  ((struct sockaddr_in*) sgw_in_addr)->sin_addr.s_addr = 0;
  return;
}
