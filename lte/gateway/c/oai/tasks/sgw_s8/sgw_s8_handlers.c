/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include "log.h"
#include "common_defs.h"
#include "sgw_context_manager.h"
#include "spgw_types.h"
#include "sgw_s8_state.h"
#include "sgw_s8_s11_handlers.h"
#include "s8_client_api.h"

uint32_t sgw_get_new_s1u_teid(sgw_state_t* state) {
  if (state->s1u_teid == 0) {
    state->s1u_teid = INITIAL_SGW_S8_S1U_TEID;
  }
  __sync_fetch_and_add(&state->s1u_teid, 1);
  return state->s1u_teid;
}

uint32_t sgw_get_new_s5s8u_teid(sgw_state_t* state) {
  __sync_fetch_and_add(&state->s5s8u_teid, 1);
  return (state->s5s8u_teid);
}

// Re-using the spgw_ue context structure, that contains the list sgw_s11_teids
// and is common across both sgw_s8 and spgw tasks.
spgw_ue_context_t* sgw_create_or_get_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p = NULL;
  hashtable_ts_get(
      sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
      (void**) &ue_context_p);
  if (!ue_context_p) {
    ue_context_p = (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
    if (ue_context_p) {
      LIST_INIT(&ue_context_p->sgw_s11_teid_list);
      hashtable_ts_insert(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
          (void*) ue_context_p);
    } else {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64, "Failed to allocate memory for UE context \n");
    }
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, ue_context_p);
}

int sgw_update_teid_in_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64, teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p =
      sgw_create_or_get_ue_context(sgw_state, imsi64);
  if (!ue_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to get UE context for sgw_s11_teid " TEID_FMT "\n", teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_teid_t* sgw_s11_teid_p =
      (sgw_s11_teid_t*) calloc(1, sizeof(sgw_s11_teid_t));
  if (!sgw_s11_teid_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to allocate memory for sgw_s11_teid:" TEID_FMT "\n", teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_teid_p->sgw_s11_teid = teid;
  LIST_INSERT_HEAD(&ue_context_p->sgw_s11_teid_list, sgw_s11_teid_p, entries);
  OAILOG_DEBUG(
      LOG_SGW_S8,
      "Inserted sgw_s11_teid to list of teids of UE context" TEID_FMT "\n",
      teid);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

sgw_eps_bearer_context_information_t*
sgw_create_bearer_context_information_in_collection(teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* new_sgw_bearer_context_information =
      calloc(1, sizeof(sgw_eps_bearer_context_information_t));

  if (new_sgw_bearer_context_information == NULL) {
    OAILOG_ERROR(
        LOG_SGW_S8,
        "Failed to create new sgw bearer context information object for "
        "sgw_s11_teid " TEID_FMT "\n",
        teid);
    return NULL;
  }
  // Insert the new tunnel with sgw_s11_teid into the hash list.
  hash_table_ts_t* state_imsi_ht = get_sgw_ue_state();
  hashtable_ts_insert(
      state_imsi_ht, (const hash_key_t) teid,
      new_sgw_bearer_context_information);

  OAILOG_DEBUG(
      LOG_SGW_S8,
      "Inserted new sgw eps bearer context into hash list,state_imsi_ht with "
      "key as sgw_s11_teid " TEID_FMT "\n ",
      teid);
  return new_sgw_bearer_context_information;
}

void sgw_s8_handle_s11_create_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_create_session_request_t* const session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64, "Received S11 CREATE SESSION REQUEST from MME_APP\n");
  sgw_eps_bearer_context_information_t* new_sgw_eps_context = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel                           = {0};
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p                  = NULL;

  increment_counter("sgw_s8_create_session", 1, NO_LABELS);
  if (session_req_pP->rat_type != RAT_EUTRAN) {
    OAILOG_WARNING_UE(
        LOG_SGW_S8, imsi64,
        "Received session request with RAT != RAT_TYPE_EUTRAN: type %d\n",
        session_req_pP->rat_type);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
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
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "sender_fteid_incorrect_parameters");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  sgw_get_new_S11_tunnel_id(&sgw_state->tunnel_id);
  sgw_s11_tunnel.local_teid  = sgw_state->tunnel_id;
  sgw_s11_tunnel.remote_teid = session_req_pP->sender_fteid_for_cp.teid;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "Rx CREATE-SESSION-REQUEST MME S11 teid " TEID_FMT
      "SGW S11 teid " TEID_FMT " APN %s EPS bearer Id %d\n",
      sgw_s11_tunnel.remote_teid, sgw_s11_tunnel.local_teid,
      session_req_pP->apn,
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
          .eps_bearer_id);

  if (sgw_update_teid_in_ue_context(
          sgw_state, imsi64, sgw_s11_tunnel.local_teid) == RETURNerror) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to update sgw_s11_teid" TEID_FMT " in UE context \n",
        sgw_s11_tunnel.local_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  new_sgw_eps_context = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  if (!new_sgw_eps_context) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Could not create new sgw context for create session req message for "
        "mme_s11_teid " TEID_FMT "\n",
        session_req_pP->sender_fteid_for_cp.teid);
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "internal_software_error");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  memcpy(
      new_sgw_eps_context->imsi.digit, session_req_pP->imsi.digit,
      IMSI_BCD_DIGITS_MAX);
  new_sgw_eps_context->imsi64 = imsi64;

  new_sgw_eps_context->imsi_unauthenticated_indicator = 1;
  new_sgw_eps_context->mme_teid_S11 = session_req_pP->sender_fteid_for_cp.teid;
  new_sgw_eps_context->s_gw_teid_S11_S4 = sgw_s11_tunnel.local_teid;
  new_sgw_eps_context->trxn             = session_req_pP->trxn;

  // Update PDN details
  if (session_req_pP->apn) {
    new_sgw_eps_context->pdn_connection.apn_in_use =
        strdup(session_req_pP->apn);
  } else {
    new_sgw_eps_context->pdn_connection.apn_in_use = strdup("NO APN");
  }
  new_sgw_eps_context->pdn_connection.s_gw_teid_S5_S8_cp =
      sgw_s11_tunnel.local_teid;
  bearer_context_to_be_created_t csr_bearer_context =
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0];
  new_sgw_eps_context->pdn_connection.default_bearer =
      csr_bearer_context.eps_bearer_id;

  /* creating an eps bearer entry
   * copy informations from create session request to bearer context information
   */

  eps_bearer_ctxt_p = sgw_cm_create_eps_bearer_ctxt_in_collection(
      &new_sgw_eps_context->pdn_connection, csr_bearer_context.eps_bearer_id);
  if (eps_bearer_ctxt_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64, "Failed to create new EPS bearer entry\n");
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "internal_software_error");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  eps_bearer_ctxt_p->eps_bearer_qos = csr_bearer_context.bearer_level_qos;
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = sgw_get_new_s1u_teid(sgw_state);
  eps_bearer_ctxt_p->s_gw_teid_S5_S8_up = sgw_get_new_s5s8u_teid(sgw_state);
  csr_bearer_context.s5_s8_u_sgw_fteid.teid =
      eps_bearer_ctxt_p->s_gw_teid_S5_S8_up;
  csr_bearer_context.s5_s8_u_sgw_fteid.ipv4 = 1;
  csr_bearer_context.s5_s8_u_sgw_fteid.ipv4_address =
      sgw_state->sgw_ip_address_S5S8_up;

  send_s8_create_session_request(
      sgw_s11_tunnel.local_teid, session_req_pP, imsi64);
  sgw_display_s11_bearer_context_information(new_sgw_eps_context);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}
