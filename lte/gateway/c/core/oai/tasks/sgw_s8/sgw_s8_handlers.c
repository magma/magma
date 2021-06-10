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
#include "3gpp_29.274.h"
#include "log.h"
#include "common_defs.h"
#include "intertask_interface.h"
#include "sgw_context_manager.h"
#include "spgw_types.h"
#include "sgw_s8_state.h"
#include "sgw_s8_s11_handlers.h"
#include "s8_client_api.h"
#include "gtpv1u.h"
#include "dynamic_memory_check.h"
#include "sgw_handlers.h"
#include "directoryd.h"
#include "conversions.h"

extern task_zmq_ctx_t sgw_s8_task_zmq_ctx;
extern struct gtp_tunnel_ops* gtp_tunnel_ops;
static int sgw_s8_add_gtp_up_tunnel(
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    sgw_eps_bearer_context_information_t* sgw_context_p);

static void sgw_send_modify_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64);

static void sgw_s8_send_failed_delete_session_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64);

static void insert_sgw_c_teid_to_directoryd(
    sgw_state_t* state, imsi64_t imsi64);

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
static void sgw_s8_populate_mbr_bearer_contexts_modified(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p);

static teid_t sgw_s8_generate_new_cp_teid(void);

void sgw_remove_sgw_bearer_context_information(
    sgw_state_t* sgw_state, teid_t teid, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  int rc = 0;

  hash_table_ts_t* state_imsi_ht = get_sgw_ue_state();
  rc                             = hashtable_ts_free(state_imsi_ht, teid);
  if (rc != HASH_TABLE_OK) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64, "Failed to free teid from state_imsi_ht\n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  spgw_ue_context_t* ue_context_p = NULL;
  hashtable_ts_get(
      sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
      (void**) &ue_context_p);
  if (ue_context_p) {
    sgw_s11_teid_t* p1 = LIST_FIRST(&(ue_context_p->sgw_s11_teid_list));
    while (p1) {
      if (p1->sgw_s11_teid == teid) {
        LIST_REMOVE(p1, entries);
        free_wrapper((void**) &p1);
        break;
      }
      p1 = LIST_NEXT(p1, entries);
    }
    if (LIST_EMPTY(&ue_context_p->sgw_s11_teid_list)) {
      rc = hashtable_ts_free(
          sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64);
      if (rc != HASH_TABLE_OK) {
        OAILOG_ERROR_UE(
            LOG_SGW_S8, imsi64,
            "Failed to free imsi64 from imsi_ue_context_htbl\n");
        OAILOG_FUNC_OUT(LOG_SGW_S8);
      }
      delete_sgw_ue_state(imsi64);
    }
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

sgw_eps_bearer_context_information_t* sgw_get_sgw_eps_bearer_context(
    teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* sgw_bearer_context_info = NULL;
  hash_table_ts_t* state_imsi_ht = get_sgw_ue_state();

  hashtable_ts_get(
      state_imsi_ht, (const hash_key_t) teid,
      (void**) &sgw_bearer_context_info);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_bearer_context_info);
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
          LOG_SGW_S8, imsi64, "Failed to allocate memory for UE context\n");
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
    sgw_state_t* sgw_state, itti_s11_create_session_request_t* session_req_pP,
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

  sgw_s11_tunnel.local_teid  = sgw_s8_generate_new_cp_teid();
  sgw_s11_tunnel.remote_teid = session_req_pP->sender_fteid_for_cp.teid;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "Rx CREATE-SESSION-REQUEST MME S11 teid " TEID_FMT
      " SGW S11 teid " TEID_FMT " APN %s EPS bearer Id %u\n",
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
  new_sgw_eps_context->imsi.length = session_req_pP->imsi.length;
  new_sgw_eps_context->imsi64      = imsi64;

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
  bearer_context_to_be_created_t* csr_bearer_context =
      &session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0];
  new_sgw_eps_context->pdn_connection.default_bearer =
      csr_bearer_context->eps_bearer_id;

  /* creating an eps bearer entry
   * copy informations from create session request to bearer context information
   */

  eps_bearer_ctxt_p = sgw_cm_create_eps_bearer_ctxt_in_collection(
      &new_sgw_eps_context->pdn_connection, csr_bearer_context->eps_bearer_id);
  if (eps_bearer_ctxt_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64, "Failed to create new EPS bearer entry\n");
    increment_counter(
        "sgw_s8_create_session", 1, 2, "result", "failure", "cause",
        "internal_software_error");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  eps_bearer_ctxt_p->eps_bearer_qos = csr_bearer_context->bearer_level_qos;
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = sgw_get_new_s1u_teid(sgw_state);
  eps_bearer_ctxt_p->s_gw_teid_S5_S8_up = sgw_get_new_s5s8u_teid(sgw_state);
  csr_bearer_context->s5_s8_u_sgw_fteid.teid =
      eps_bearer_ctxt_p->s_gw_teid_S5_S8_up;
  csr_bearer_context->s5_s8_u_sgw_fteid.ipv4 = 1;
  csr_bearer_context->s5_s8_u_sgw_fteid.ipv4_address =
      sgw_state->sgw_ip_address_S5S8_up;

  send_s8_create_session_request(
      sgw_s11_tunnel.local_teid, session_req_pP, imsi64);
  sgw_display_s11_bearer_context_information(LOG_SGW_S8, new_sgw_eps_context);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static int update_bearer_context_info(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const s8_create_session_response_t* const session_rsp_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_ctxt_t* default_bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
      &sgw_context_p->pdn_connection, session_rsp_p->eps_bearer_id);
  if (!default_bearer_ctx_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to get default eps bearer context for context teid " TEID_FMT
        "and bearer_id :%u \n",
        session_rsp_p->context_teid, session_rsp_p->eps_bearer_id);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  memcpy(&default_bearer_ctx_p->paa, &session_rsp_p->paa, sizeof(paa_t));
  if (session_rsp_p->eps_bearer_id !=
      session_rsp_p->bearer_context[0].eps_bearer_id) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Mismatch of eps bearer id between bearer context's bearer id:%d and "
        "default eps bearer id:%u for context teid " TEID_FMT "\n",
        session_rsp_p->bearer_context[0].eps_bearer_id,
        session_rsp_p->eps_bearer_id, session_rsp_p->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  s8_bearer_context_t s5s8_bearer_context = session_rsp_p->bearer_context[0];
  FTEID_T_2_IP_ADDRESS_T(
      (&s5s8_bearer_context.pgw_s8_up),
      (&default_bearer_ctx_p->p_gw_address_in_use_up));
  default_bearer_ctx_p->p_gw_teid_S5_S8_up = s5s8_bearer_context.pgw_s8_up.teid;

  memcpy(
      &default_bearer_ctx_p->eps_bearer_qos, &s5s8_bearer_context.qos,
      sizeof(bearer_qos_t));
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

static int sgw_s8_send_create_session_response(
    sgw_state_t* sgw_state, sgw_eps_bearer_context_information_t* sgw_context_p,
    s8_create_session_response_t* session_rsp_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p                                         = NULL;
  itti_s11_create_session_response_t* create_session_response_p = NULL;

  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_CREATE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_create_session_response \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64, "sgw_context_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  create_session_response_p = &message_p->ittiMsg.s11_create_session_response;
  if (!create_session_response_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "create_session_response_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  create_session_response_p->teid              = sgw_context_p->mme_teid_S11;
  create_session_response_p->cause.cause_value = session_rsp_p->cause;
  create_session_response_p->s11_sgw_fteid.teid =
      sgw_context_p->s_gw_teid_S11_S4;
  create_session_response_p->s11_sgw_fteid.interface_type = S11_SGW_GTP_C;
  create_session_response_p->s11_sgw_fteid.ipv4           = 1;
  create_session_response_p->s11_sgw_fteid.ipv4_address.s_addr =
      spgw_config.sgw_config.ipv4.S11.s_addr;

  if (session_rsp_p->cause == REQUEST_ACCEPTED) {
    sgw_eps_bearer_ctxt_t* default_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, session_rsp_p->eps_bearer_id);
    if (!default_bearer_ctxt_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "Failed to get default eps bearer context for sgw_s11_teid " TEID_FMT
          "and bearer_id :%u \n",
          session_rsp_p->context_teid, session_rsp_p->eps_bearer_id);
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    memcpy(
        &create_session_response_p->paa, &default_bearer_ctxt_p->paa,
        sizeof(paa_t));
    create_session_response_p->bearer_contexts_created.num_bearer_context = 1;
    bearer_context_created_t* bearer_context =
        &create_session_response_p->bearer_contexts_created.bearer_contexts[0];

    bearer_context->s1u_sgw_fteid.teid =
        default_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up;
    bearer_context->s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
    bearer_context->s1u_sgw_fteid.ipv4           = 1;
    bearer_context->s1u_sgw_fteid.ipv4_address.s_addr =
        sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr;
    bearer_context->eps_bearer_id = session_rsp_p->eps_bearer_id;
    /*
     * Set the Cause information from bearer context created.
     * "Request accepted" is returned when the GTPv2 entity has accepted a
     * control plane request.
     */
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .cause.cause_value = session_rsp_p->cause;
    if (session_rsp_p->pco.num_protocol_or_container_id) {
      copy_protocol_configuration_options(
          &create_session_response_p->pco, &session_rsp_p->pco);
      clear_protocol_configuration_options(&session_rsp_p->pco);
    }
  } else {
    create_session_response_p->bearer_contexts_marked_for_removal
        .num_bearer_context = 1;
    bearer_context_marked_for_removal_t* bearer_context =
        &create_session_response_p->bearer_contexts_marked_for_removal
             .bearer_contexts[0];
    bearer_context->cause.cause_value = session_rsp_p->cause;
    bearer_context->eps_bearer_id     = session_rsp_p->eps_bearer_id;
    create_session_response_p->trxn   = sgw_context_p->trxn;
  }
  message_p->ittiMsgHeader.imsi = sgw_context_p->imsi64;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, sgw_context_p->imsi64,
      "Sending S11 Create Session Response to mme for mme_s11_teid: " TEID_FMT
      "\n",
      create_session_response_p->teid);

  send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);

  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

void sgw_s8_handle_create_session_response(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  if (!session_rsp_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received null create session response from s8_proxy\n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      " Rx S5S8_CREATE_SESSION_RSP for context_teid " TEID_FMT "\n",
      session_rsp_p->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(session_rsp_p->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from "
        "context_teid " TEID_FMT " \n",
        session_rsp_p->context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  if (session_rsp_p->cause == REQUEST_ACCEPTED) {
    // update pdn details received from PGW
    sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp =
        session_rsp_p->pgw_s8_cp_teid.teid;
    // update bearer context details received from PGW
    if ((update_bearer_context_info(sgw_context_p, session_rsp_p)) !=
        RETURNok) {
      /*TODO need to send delete session request to pgw */
      session_rsp_p->cause = CONTEXT_NOT_FOUND;
    }
  }
  // send Create session response to mme
  if ((sgw_s8_send_create_session_response(
          sgw_state, sgw_context_p, session_rsp_p)) != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to send create session response to mme for "
        "sgw_s11_teid " TEID_FMT "\n",
        session_rsp_p->context_teid);
  }
  if (session_rsp_p->cause != REQUEST_ACCEPTED) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received failed create session response with cause: %d for "
        "context_id: " TEID_FMT "\n",
        session_rsp_p->cause, session_rsp_p->context_teid);
    sgw_remove_sgw_bearer_context_information(
        sgw_state, session_rsp_p->context_teid, imsi64);
  } else {
    insert_sgw_c_teid_to_directoryd(sgw_state, imsi64);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

// The function generates comma separated list of sgw's control plane teid of
// s8 interface for each pdn session and updates the list to directoryd.
static void insert_sgw_c_teid_to_directoryd(
    sgw_state_t* sgw_state, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  char teidString[16]                    = {0};
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);
  spgw_ue_context_t* ue_context_p = NULL;
  // TODO move imsi_ue_context_htbl from sgw_state to sgw_s8's state manager
  hashtable_ts_get(
      sgw_state->imsi_ue_context_htbl, (const hash_key_t) imsi64,
      (void**) &ue_context_p);
  if (ue_context_p) {
    sgw_s11_teid_t* s11_teid_p = NULL;
    LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
      if (s11_teid_p) {
        if (s11_teid_p->entries.le_next) {
          snprintf(
              teidString + strlen(teidString),
              sizeof(teidString) - strlen(teidString), "%u,",
              s11_teid_p->sgw_s11_teid);
        } else {
          snprintf(
              teidString + strlen(teidString),
              sizeof(teidString) - strlen(teidString), "%u",
              s11_teid_p->sgw_s11_teid);
        }
      }
    }
    directoryd_update_record_field(imsi_str, "sgw_c_teid", teidString);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void sgw_s8_handle_modify_bearer_request(
    sgw_state_t* state,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  uint8_t idx                                                    = 0;
  uint8_t sgi_rsp_idx                                            = 0;
  itti_sgi_update_end_point_response_t sgi_update_end_point_resp = {0};
  struct in_addr enb                  = {.s_addr = 0};
  struct in_addr pgw                  = {.s_addr = 0};
  sgw_eps_bearer_ctxt_t* bearer_ctx_p = NULL;

  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64, "Rx MODIFY_BEARER_REQUEST, teid " TEID_FMT "\n",
      modify_bearer_pP->teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(modify_bearer_pP->teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from "
        "context_teid " TEID_FMT " \n",
        modify_bearer_pP->teid);
    if ((send_mbr_failure(LOG_SGW_S8, modify_bearer_pP, imsi64) != RETURNok)) {
      OAILOG_ERROR(
          LOG_SGW_S8,
          "Error in sending modify bearer response to MME App for context "
          "teid " TEID_FMT "\n",
          modify_bearer_pP->teid);
    }
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  sgw_context_p->trxn                    = modify_bearer_pP->trxn;
  sgi_update_end_point_resp.context_teid = modify_bearer_pP->teid;
  // Traversing through the list of bearers to be modified
  for (; idx <
         modify_bearer_pP->bearer_contexts_to_be_modified.num_bearer_context;
       idx++) {
    bearer_context_to_be_modified_t mbr_bearer_ctxt_p =
        modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx];
    bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, mbr_bearer_ctxt_p.eps_bearer_id);
    if (!bearer_ctx_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64,
          "Failed to get eps bearer context for context teid " TEID_FMT
          "and bearer_id :%u \n",
          modify_bearer_pP->teid, mbr_bearer_ctxt_p.eps_bearer_id);
      sgi_update_end_point_resp.bearer_contexts_not_found[sgi_rsp_idx++] =
          mbr_bearer_ctxt_p.eps_bearer_id;
      sgi_update_end_point_resp.num_bearers_not_found++;
    } else {
      enb.s_addr = bearer_ctx_p->enb_ip_address_S1u.address.ipv4_address.s_addr;
      pgw.s_addr =
          bearer_ctx_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;

      // Send end marker to eNB and then delete the tunnel if enb_ip is
      // different
      if (does_bearer_context_hold_valid_enb_ip(
              bearer_ctx_p->enb_ip_address_S1u) &&
          is_enb_ip_address_same(
              &mbr_bearer_ctxt_p.s1_eNB_fteid,
              &bearer_ctx_p->enb_ip_address_S1u) == false) {
        struct in_addr ue_ipv4   = bearer_ctx_p->paa.ipv4_address;
        struct in6_addr* ue_ipv6 = NULL;
        if ((bearer_ctx_p->paa.pdn_type == IPv6) ||
            (bearer_ctx_p->paa.pdn_type == IPv4_AND_v6)) {
          ue_ipv6 = &bearer_ctx_p->paa.ipv6_address;
        }

        OAILOG_DEBUG_UE(
            LOG_SGW_S8, imsi64,
            "Delete GTPv1-U tunnel for sgw_teid:" TEID_FMT "for bearer %u\n",
            bearer_ctx_p->s_gw_teid_S1u_S12_S4_up, bearer_ctx_p->eps_bearer_id);
        // This is best effort, ignore return code.
        gtp_tunnel_ops->send_end_marker(enb, modify_bearer_pP->teid);
        // delete GTPv1-U tunnel
        gtpv1u_del_s8_tunnel(
            enb, pgw, ue_ipv4, ue_ipv6, bearer_ctx_p->s_gw_teid_S1u_S12_S4_up,
            bearer_ctx_p->enb_teid_S1u, NULL);
      }
      populate_sgi_end_point_update(
          sgi_rsp_idx, idx, modify_bearer_pP, bearer_ctx_p,
          &sgi_update_end_point_resp);
      sgi_rsp_idx++;
    }
  }  // for loop

  sgi_rsp_idx = 0;
  for (idx = 0;
       idx < modify_bearer_pP->bearer_contexts_to_be_removed.num_bearer_context;
       idx++) {
    bearer_ctx_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection,
        modify_bearer_pP->bearer_contexts_to_be_removed.bearer_contexts[idx]
            .eps_bearer_id);
    if (bearer_ctx_p) {
      sgi_update_end_point_resp.bearer_contexts_to_be_removed[sgi_rsp_idx++] =
          bearer_ctx_p->eps_bearer_id;
      sgi_update_end_point_resp.num_bearers_removed++;
    }
  }
  sgw_send_modify_bearer_response(
      sgw_context_p, &sgi_update_end_point_resp, imsi64);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void sgw_send_modify_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const itti_sgi_update_end_point_response_t* const resp_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  itti_s11_modify_bearer_response_t* modify_response_p = NULL;
  MessageDef* message_p                                = NULL;

  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "send modify bearer response for Context teid " TEID_FMT "\n",
      resp_pP->context_teid);
  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_MODIFY_BEARER_RESPONSE);

  if (!message_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to allocate memory for S11_MODIFY_BEARER_RESPONSE\n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  modify_response_p = &message_p->ittiMsg.s11_modify_bearer_response;

  if (sgw_context_p) {
    modify_response_p->teid              = sgw_context_p->mme_teid_S11;
    modify_response_p->cause.cause_value = REQUEST_ACCEPTED;
    modify_response_p->trxn              = sgw_context_p->trxn;
    message_p->ittiMsgHeader.imsi        = imsi64;

    sgw_s8_populate_mbr_bearer_contexts_modified(
        resp_pP, imsi64, sgw_context_p, modify_response_p);
    sgw_populate_mbr_bearer_contexts_removed(
        resp_pP, imsi64, sgw_context_p, modify_response_p);
    sgw_populate_mbr_bearer_contexts_not_found(
        LOG_SGW_S8, resp_pP, modify_response_p);
    send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void sgw_s8_populate_mbr_bearer_contexts_modified(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t rsp_idx                          = 0;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = NULL;

  for (uint8_t idx = 0; idx < resp_pP->num_bearers_modified; idx++) {
    eps_bearer_ctxt_p = sgw_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection,
        resp_pP->bearer_contexts_to_be_modified[idx].eps_bearer_id);

    if (NULL != eps_bearer_ctxt_p) {
      OAILOG_DEBUG_UE(
          LOG_SGW_S8, imsi64,
          "Modify bearer request is accepted for bearer_id :%u\n",
          resp_pP->bearer_contexts_to_be_modified[idx].eps_bearer_id);
      modify_response_p->bearer_contexts_modified.bearer_contexts[rsp_idx]
          .eps_bearer_id =
          resp_pP->bearer_contexts_to_be_modified[idx].eps_bearer_id;
      modify_response_p->bearer_contexts_modified.bearer_contexts[rsp_idx++]
          .cause.cause_value = REQUEST_ACCEPTED;
      modify_response_p->bearer_contexts_modified.num_bearer_context++;

      // setup GTPv1-U tunnels, both s1-u and s8-u tunnels
      sgw_s8_add_gtp_up_tunnel(eps_bearer_ctxt_p, sgw_context_p);
      // may be removed TODO rashmi remove after testing
      if (TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX >
          eps_bearer_ctxt_p->num_sdf) {
        int i = 0;
        while ((i < eps_bearer_ctxt_p->num_sdf) &&
               (SDF_ID_NGBR_DEFAULT != eps_bearer_ctxt_p->sdf_id[i]))
          i++;
        if (i >= eps_bearer_ctxt_p->num_sdf) {
          eps_bearer_ctxt_p->sdf_id[eps_bearer_ctxt_p->num_sdf] =
              SDF_ID_NGBR_DEFAULT;
          eps_bearer_ctxt_p->num_sdf += 1;
        }
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

// Helper function to add gtp tunnels for default bearers
static int sgw_s8_add_gtp_up_tunnel(
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    sgw_eps_bearer_context_information_t* sgw_context_p) {
  int rv             = RETURNok;
  struct in_addr enb = {.s_addr = 0};
  struct in_addr pgw = {.s_addr = 0};
  pgw.s_addr =
      eps_bearer_ctxt_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;
  if ((pgw.s_addr == 0) && (eps_bearer_ctxt_p->p_gw_teid_S5_S8_up == 0)) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "bearer context has invalid pgw_s8_teid " TEID_FMT
        "pgw_ip address :%x \n",
        eps_bearer_ctxt_p->p_gw_teid_S5_S8_up, pgw.s_addr);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  enb.s_addr =
      eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr;

  struct in_addr ue_ipv4   = {.s_addr = 0};
  struct in6_addr* ue_ipv6 = NULL;
  ue_ipv4.s_addr           = eps_bearer_ctxt_p->paa.ipv4_address.s_addr;
  if ((eps_bearer_ctxt_p->paa.pdn_type == IPv6) ||
      (eps_bearer_ctxt_p->paa.pdn_type == IPv4_AND_v6)) {
    ue_ipv6 = &eps_bearer_ctxt_p->paa.ipv6_address;
  }

  int vlan    = eps_bearer_ctxt_p->paa.vlan;
  Imsi_t imsi = sgw_context_p->imsi;

  char ip6_str[INET6_ADDRSTRLEN];
  if (ue_ipv6) {
    inet_ntop(AF_INET6, ue_ipv6, ip6_str, INET6_ADDRSTRLEN);
  }
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, sgw_context_p->imsi64,
      "Adding tunnel for bearer_id %u ue addr %x enb %x "
      "s_gw_teid_S1u_S12_S4_up %x, enb_teid_S1u %x pgw_up_ip %x pgw_up_teid %x "
      "s_gw_ip_address_S5_S8_up %x"
      "s_gw_teid_S5_S8_up %x \n ",
      eps_bearer_ctxt_p->eps_bearer_id, ue_ipv4.s_addr, enb.s_addr,
      eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
      eps_bearer_ctxt_p->enb_teid_S1u, pgw.s_addr,
      eps_bearer_ctxt_p->p_gw_teid_S5_S8_up,
      eps_bearer_ctxt_p->s_gw_ip_address_S5_S8_up.address.ipv4_address.s_addr,
      eps_bearer_ctxt_p->s_gw_teid_S5_S8_up);
  if (eps_bearer_ctxt_p->eps_bearer_id ==
      sgw_context_p->pdn_connection.default_bearer) {
    // Set default precedence and tft for default bearer
    if (ue_ipv6) {
      OAILOG_INFO_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "Adding tunnel for ipv6 ue addr %s, enb %x, "
          "s_gw_teid_S5_S8_up %x, s_gw_ip_address_S5_S8_up %x pgw_up_ip %x "
          "pgw_up_teid %x \n",
          ip6_str, enb.s_addr, eps_bearer_ctxt_p->s_gw_teid_S5_S8_up,
          eps_bearer_ctxt_p->s_gw_ip_address_S5_S8_up.address.ipv4_address
              .s_addr,
          pgw.s_addr, eps_bearer_ctxt_p->p_gw_teid_S5_S8_up);
    }
    rv = gtpv1u_add_s8_tunnel(
        ue_ipv4, ue_ipv6, vlan, enb, pgw,
        eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
        eps_bearer_ctxt_p->enb_teid_S1u, eps_bearer_ctxt_p->s_gw_teid_S5_S8_up,
        eps_bearer_ctxt_p->p_gw_teid_S5_S8_up, imsi, NULL, DEFAULT_PRECEDENCE);
    if (rv < 0) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "ERROR in setting up TUNNEL err=%d\n", rv);
    }
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, rv);
}

void sgw_s8_handle_s11_delete_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  gtpv2c_cause_value_t gtpv2c_cause = 0;
  if (!delete_session_req_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received NULL delete_session_req_p from mme app \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Received S11 DELETE SESSION REQUEST for sgw_s11_teid " TEID_FMT "\n",
      delete_session_req_p->teid);
  increment_counter("sgw_delete_session", 1, NO_LABELS);
  if (delete_session_req_p->indication_flags.oi) {
    OAILOG_DEBUG_UE(
        LOG_SPGW_APP, imsi64,
        "OI flag is set for this message indicating the request"
        "should be forwarded to P-GW entity\n");
  }

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(delete_session_req_p->teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from "
        "sgw_s11_teid " TEID_FMT " \n",
        delete_session_req_p->teid);
    gtpv2c_cause = CONTEXT_NOT_FOUND;
    sgw_s8_send_failed_delete_session_response(
        sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  if ((delete_session_req_p->sender_fteid_for_cp.ipv4) &&
      (delete_session_req_p->sender_fteid_for_cp.ipv6)) {
    // Sender F-TEID IE present
    if (delete_session_req_p->teid != sgw_context_p->mme_teid_S11) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, imsi64,
          "Mismatch in MME Teid for CP teid recevied in delete session "
          "req: " TEID_FMT " teid present in sgw_context :" TEID_FMT "\n",
          delete_session_req_p->teid, sgw_context_p->mme_teid_S11);
      gtpv2c_cause = INVALID_PEER;
      sgw_s8_send_failed_delete_session_response(
          sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
      OAILOG_FUNC_OUT(LOG_SGW_S8);
    }
  }
  if (delete_session_req_p->lbi !=
      sgw_context_p->pdn_connection.default_bearer) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Mismatch in default eps bearer_id, bearer_id recevied in delete "
        "session req :%d and bearer_id present in sgw_context :%d \n",
        delete_session_req_p->lbi,
        sgw_context_p->pdn_connection.default_bearer);
    sgw_s8_send_failed_delete_session_response(
        sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  send_s8_delete_session_request(
      sgw_context_p->imsi64, sgw_context_p->imsi,
      sgw_context_p->s_gw_teid_S11_S4,
      sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp,
      sgw_context_p->pdn_connection.default_bearer, delete_session_req_p);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void delete_userplane_tunnels(
    sgw_eps_bearer_context_information_t* sgw_context_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  struct in_addr enb                   = {.s_addr = 0};
  struct in_addr pgw                   = {.s_addr = 0};
  sgw_eps_bearer_ctxt_t* bearer_ctxt_p = NULL;
  int rv                               = RETURNerror;
  struct in_addr ue_ipv4               = {.s_addr = 0};

  for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
    ebi_t ebi = INDEX_TO_EBI(ebix);
    bearer_ctxt_p =
        sgw_cm_get_eps_bearer_entry(&sgw_context_p->pdn_connection, ebi);

    if (bearer_ctxt_p) {
      enb.s_addr =
          bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr;
      pgw.s_addr =
          bearer_ctxt_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;
      struct in6_addr* ue_ipv6 = NULL;
      if ((bearer_ctxt_p->paa.pdn_type == IPv6) ||
          (bearer_ctxt_p->paa.pdn_type == IPv4_AND_v6)) {
        ue_ipv6 = &bearer_ctxt_p->paa.ipv6_address;
      }
      ue_ipv4 = bearer_ctxt_p->paa.ipv4_address;
      // Delete S1-U tunnel and S8-U tunnel
      rv = gtpv1u_del_s8_tunnel(
          enb, pgw, ue_ipv4, ue_ipv6, bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
          bearer_ctxt_p->enb_teid_S1u, NULL);
      if (rv < 0) {
        OAILOG_ERROR_UE(
            LOG_SPGW_APP, sgw_context_p->imsi64,
            "ERROR in deleting S1-U TUNNEL " TEID_FMT
            " (eNB) <-> (SGW) " TEID_FMT "\n",
            bearer_ctxt_p->enb_teid_S1u,
            bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up);
      }
      // delete paging rule
      char* ip_str = inet_ntoa(ue_ipv4);
      rv           = gtp_tunnel_ops->delete_paging_rule(ue_ipv4);
      if (rv < 0) {
        OAILOG_ERROR(
            LOG_SPGW_APP, "ERROR in deleting paging rule for IP Addr: %s\n",
            ip_str);
      } else {
        OAILOG_DEBUG(LOG_SPGW_APP, "Stopped paging for IP Addr: %s\n", ip_str);
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void sgw_s8_handle_delete_session_response(
    sgw_state_t* sgw_state, s8_delete_session_response_t* session_rsp_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p                                         = NULL;
  itti_s11_delete_session_response_t* delete_session_response_p = NULL;

  if (!session_rsp_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received null delete session response from s8_proxy\n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      " Rx S5S8_DELETE_SESSION_RSP from s8_proxy for context_teid " TEID_FMT
      "\n",
      session_rsp_p->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(session_rsp_p->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from "
        "context_teid " TEID_FMT " \n",
        session_rsp_p->context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_DELETE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_delete_session_response for "
        "context_teid " TEID_FMT "\n",
        session_rsp_p->context_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  delete_session_response_p = &message_p->ittiMsg.s11_delete_session_response;
  message_p->ittiMsgHeader.imsi = imsi64;
  if (!delete_session_response_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "delete_session_response_p is NULL \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  delete_session_response_p->teid              = sgw_context_p->mme_teid_S11;
  delete_session_response_p->cause.cause_value = session_rsp_p->cause;
  delete_session_response_p->trxn              = sgw_context_p->trxn;
  delete_session_response_p->lbi = sgw_context_p->pdn_connection.default_bearer;

  // Delete ovs rules
  delete_userplane_tunnels(sgw_context_p);
  sgw_remove_sgw_bearer_context_information(
      sgw_state, session_rsp_p->context_teid, imsi64);
  // send delete session response to mme
  if (send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to send delete session response to mme for "
        "sgw_s11_teid " TEID_FMT "\n",
        session_rsp_p->context_teid);
  }
  increment_counter("sgw_delete_session", 1, 1, "result", "success");
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void sgw_s8_send_failed_delete_session_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p                                         = NULL;
  itti_s11_delete_session_response_t* delete_session_response_p = NULL;
  teid_t teid                                                   = 0;

  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_DELETE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_delete_session_response \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  delete_session_response_p = &message_p->ittiMsg.s11_delete_session_response;
  message_p->ittiMsgHeader.imsi = imsi64;
  if (!delete_session_response_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "delete_session_response_p is NULL \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  if (sgw_context_p) {
    delete_session_response_p->teid              = sgw_context_p->mme_teid_S11;
    delete_session_response_p->cause.cause_value = cause;
    delete_session_response_p->trxn              = sgw_context_p->trxn;
    delete_session_response_p->lbi =
        sgw_context_p->pdn_connection.default_bearer;

    // Delete ovs rules
    delete_userplane_tunnels(sgw_context_p);
    sgw_remove_sgw_bearer_context_information(
        sgw_state, sgw_context_p->s_gw_teid_S11_S4, imsi64);
    teid = sgw_context_p->s_gw_teid_S11_S4;
  } else {
    if (delete_session_req_p) {
      delete_session_response_p->teid = delete_session_req_p->local_teid;
      delete_session_response_p->cause.cause_value = cause;
      delete_session_response_p->trxn              = delete_session_req_p->trxn;
      delete_session_response_p->lbi               = delete_session_req_p->lbi;
      teid                                         = delete_session_req_p->teid;
    }
  }
  increment_counter("sgw_delete_session", 1, 1, "result", "failed");
  // send delete session response to mme
  if (send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to send delete session response to mme for "
        "sgw_s11_teid " TEID_FMT "\n",
        teid);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

/* From GPP TS 23.401 version 11.11.0 Release 11, section 5.3.5 S1 release
   procedure: The S-GW releases all eNodeB related information (address and
   TEIDs) for the UE and responds with a Release Access Bearers Response message
   to the MME. Other elements of the UE's S-GW context are not affected. The
   S-GW retains the S1-U configuration that the S-GW allocated for the UE's
   bearers. The S-GW starts buffering downlink packets received for the UE and
   initiating the "Network Triggered Service Request" procedure, described in
   clause 5.3.4.3, if downlink packets arrive for the UE.
*/
//------------------------------------------------------------------------------
void sgw_s8_handle_release_access_bearers_request(
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "Release Access Bearer Request Received in SGW_S8 task \n");

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(release_access_bearers_req_pP->teid);
  if (sgw_context_p) {
    sgw_send_release_access_bearer_response(
        LOG_SGW_S8, imsi64, REQUEST_ACCEPTED, release_access_bearers_req_pP,
        sgw_context_p->mme_teid_S11);
    sgw_process_release_access_bearer_request(
        LOG_SGW_S8, imsi64, sgw_context_p);
  } else {
    sgw_send_release_access_bearer_response(
        LOG_SGW_S8, imsi64, CONTEXT_NOT_FOUND, release_access_bearers_req_pP,
        0);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

/* Generates random control plane teid and is used for both s11 and
 * s8 interface to uniquely identify sgw_s8_context for pdn session
 */
static teid_t sgw_s8_generate_new_cp_teid(void) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* sgw_context_p = NULL;
  teid_t teid                                         = INVALID_TEID;
  // note srand with seed is initialized at main
  do {
    teid          = (teid_t) rand();
    sgw_context_p = sgw_get_sgw_eps_bearer_context(teid);
  } while (sgw_context_p);

  OAILOG_FUNC_RETURN(LOG_SGW_S8, teid);
}
