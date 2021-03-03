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

/*! \file sgw_s8_handlers.c
  \brief
  \author
  \company
  \email:
*/
#include <stdio.h>
#include <stdint.h>
#include "log.h"
#include "common_defs.h"
#include "intertask_interface.h"
#include "common_types.h"
#include "sgw_context_manager.h"

int sgw_s8_handle_s11_create_session_request(
    const itti_s11_create_session_request_t* const session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64, "Received S11 CREATE SESSION REQUEST from MME_APP\n");
#if 0
  sgw_eps_bearer_context_information_t* sgw_eps_bearer_ctxt_info_p = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel                                  = {0};

  increment_counter("spgw_create_session", 1, NO_LABELS);
  if (session_req_pP->rat_type != RAT_EUTRAN) {
    OAILOG_WARNING_UE(
        LOG_SGW_S8, imsi64,
        "Received session request with RAT != RAT_TYPE_EUTRAN: type %d\n",
        session_req_pP->rat_type);
  }
  /*
   * As we are abstracting GTP-C transport, FTeid ip address is useless.
   * We just use the teid to identify MME tunnel. Normally we received either:
   * - ipv4 address if ipv4 flag is set
   * - ipv6 address if ipv6 flag is set
   * - ipv4 and ipv6 if both flags are set
   * Communication between MME and S-GW involves S11 interface so we are
   * expecting S11_MME_GTP_C (11) as interface_type.
   */
  if ((session_req_pP->sender_fteid_for_cp.teid == 0) &&
      (session_req_pP->sender_fteid_for_cp.interface_type != S11_MME_GTP_C)) {
    // MME sent request with teid = 0. This is not valid...
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Received invalid teid \n");
    increment_counter(
        "spgw_create_session", 1, 2, "result", "failure", "cause",
        "sender_fteid_incorrect_parameters");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  sgw_s11_tunnel.local_teid  = 1;
  sgw_s11_tunnel.remote_teid = session_req_pP->sender_fteid_for_cp.teid;
  mme_sgw_tunnel_t* new_endpoint_p = NULL;
  new_endpoint_p =
      sgw_cm_create_s11_tunnel(session_req_pP->sender_fteid_for_cp.teid, 1);
#endif

  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}
