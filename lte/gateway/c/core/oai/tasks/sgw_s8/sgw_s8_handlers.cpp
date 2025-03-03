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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/sgw_s8_state.hpp"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/lib/directoryd/directoryd.hpp"
#include "lte/gateway/c/core/oai/lib/s8_proxy/s8_client_api.hpp"
#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_s11_handlers.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

extern task_zmq_ctx_t sgw_s8_task_zmq_ctx;
extern struct gtp_tunnel_ops* gtp_tunnel_ops;

#if !MME_UNIT_TEST
static int sgw_s8_add_gtp_up_tunnel(
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    sgw_eps_bearer_context_information_t* sgw_context_p);
#endif

static void sgw_send_modify_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64);

static void sgw_s8_send_failed_delete_session_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64);

static void insert_sgw_cp_and_up_teid_to_directoryd(sgw_state_t* state,
                                                    imsi64_t imsi64,
                                                    uint8_t teid_type);

uint32_t sgw_get_new_s1u_teid(sgw_state_t* state) {
  if (state->s1u_teid == 0) {
    state->s1u_teid = INITIAL_SGW_S8_S1U_TEID;
  }
  __sync_fetch_and_add(&state->s1u_teid, 1);
  return state->s1u_teid;
}

static void sgw_s8_populate_mbr_bearer_contexts_modified(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p);

static sgw_eps_bearer_context_information_t* update_sgw_context_to_s11_teid_map(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64);

// TODO(rsarwad): shall be removed while porting sgw_s8 context to protobuf
// github issue: 11191

static void populate_s8_sgi_end_point_update(
    uint8_t sgi_rsp_idx, uint8_t idx,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    itti_sgi_update_end_point_response_t* sgi_update_end_point_resp);

static void sgw_s8_process_release_access_bearer_request(
    log_proto_t module, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context);

static void sgw_s8_release_all_enb_related_information(
    sgw_eps_bearer_ctxt_t* const eps_bearer_ctxt);

status_code_e sgw_build_and_send_s11_create_bearer_request(
    sgw_eps_bearer_context_information_t* sgw_eps_bearer_context_information,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    log_proto_t module);

void sgw_s8_populate_mbr_bearer_contexts_removed(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p);

sgw_eps_bearer_ctxt_t* sgw_cm_insert_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection,
    sgw_eps_bearer_ctxt_t* const sgw_eps_bearer_ctxt);

bool is_sgw_enb_ip_address_same(const fteid_t* fte_p, ip_address_t* ip_p);

void sgw_remove_sgw_bearer_context_information(sgw_state_t* sgw_state,
                                               teid_t teid, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  map_uint32_sgw_eps_bearer_context_t* state_teid_map = get_s8_state_teid_map();
  if (!state_teid_map) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Failed to get state_teid_map");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  if (state_teid_map->remove(teid) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to free teid :" TEID_FMT "from state_teid_map",
                    teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  spgw_ue_context_t* ue_context_p = nullptr;
  sgw_state->imsi_ue_context_map.get(imsi64, &ue_context_p);
  if (ue_context_p) {
    sgw_s11_teid_t* p1 = LIST_FIRST(&(ue_context_p->sgw_s11_teid_list));
    while (p1) {
      if (p1->sgw_s11_teid == teid) {
        LIST_REMOVE(p1, entries);
        free_cpp_wrapper(reinterpret_cast<void**>(&p1));
        break;
      }
      p1 = LIST_NEXT(p1, entries);
    }
    if (LIST_EMPTY(&ue_context_p->sgw_s11_teid_list)) {
      if (sgw_state->imsi_ue_context_map.remove(imsi64) !=
          magma::PROTO_MAP_OK) {
        OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                        "Failed to remove imsi64 from imsi_ue_context_map\n");
        OAILOG_FUNC_OUT(LOG_SGW_S8);
      }
      delete_sgw_ue_state(imsi64);
    }
    insert_sgw_cp_and_up_teid_to_directoryd(sgw_state, imsi64, UP_TEID);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

sgw_eps_bearer_context_information_t* sgw_get_sgw_eps_bearer_context(
    teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_context_information_t* sgw_bearer_context_info = nullptr;
  map_uint32_sgw_eps_bearer_context_t* state_teid_map = get_s8_state_teid_map();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to get state_teid_map");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, nullptr);
  }

  state_teid_map->get(teid, &sgw_bearer_context_info);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_bearer_context_info);
}

// Re-using the spgw_ue context structure, that contains the list sgw_s11_teids
// and is common across both sgw_s8 and spgw tasks.
spgw_ue_context_t* sgw_create_or_get_ue_context(sgw_state_t* sgw_state,
                                                imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p = nullptr;
  sgw_state->imsi_ue_context_map.get(imsi64, &ue_context_p);
  if (!ue_context_p) {
    ue_context_p = new spgw_ue_context_t();
    if (ue_context_p) {
      LIST_INIT(&ue_context_p->sgw_s11_teid_list);
      sgw_state->imsi_ue_context_map.insert(imsi64, ue_context_p);
    }
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, ue_context_p);
}

status_code_e sgw_update_teid_in_ue_context(sgw_state_t* sgw_state,
                                            imsi64_t imsi64, teid_t teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  spgw_ue_context_t* ue_context_p =
      sgw_create_or_get_ue_context(sgw_state, imsi64);
  if (!ue_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to get UE context for sgw_s11_teid " TEID_FMT "\n",
                    teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_teid_t* sgw_s11_teid_p = new sgw_s11_teid_t();
  sgw_s11_teid_p->sgw_s11_teid = teid;
  LIST_INSERT_HEAD(&ue_context_p->sgw_s11_teid_list, sgw_s11_teid_p, entries);
  OAILOG_DEBUG(LOG_SGW_S8,
               "Inserted sgw_s11_teid to list of teids of UE context" TEID_FMT
               "\n",
               teid);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

sgw_eps_bearer_context_information_t*
sgw_create_bearer_context_information_in_collection(
    sgw_state_t* sgw_state, uint32_t* temporary_create_session_procedure_id_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  *temporary_create_session_procedure_id_p = (uint32_t)rand();
  sgw_eps_bearer_context_information_t* new_sgw_bearer_context_information =
      new sgw_eps_bearer_context_information_t();

  sgw_state->temporary_create_session_procedure_id_map.insert(
      *temporary_create_session_procedure_id_p,
      new_sgw_bearer_context_information);

  OAILOG_DEBUG(LOG_SGW_S8,
               "Inserted new sgw eps bearer context into map, "
               "temporary_create_session_procedure_id_map with "
               "key as temporary_create_session_procedure_id_p :%u \n",
               *temporary_create_session_procedure_id_p);
  return new_sgw_bearer_context_information;
}

bool check_empty_apn(char* apn) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
#define MAX_APN_LEN (ACCESS_POINT_NAME_MAX_LENGTH + 1)

  char zerobuf[MAX_APN_LEN] = {0};
  if (memcmp(apn, zerobuf, MAX_APN_LEN) == 0) {
    OAILOG_FUNC_RETURN(LOG_SGW_S8, true);
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, false);
}

int sgw_update_bearer_context_information_on_csreq(
    sgw_state_t* sgw_state,
    sgw_eps_bearer_context_information_t* new_sgw_eps_context,
    itti_s11_create_session_request_t* session_req_pP, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = nullptr;
  memcpy(new_sgw_eps_context->imsi.digit, session_req_pP->imsi.digit,
         IMSI_BCD_DIGITS_MAX);
  new_sgw_eps_context->imsi.length = session_req_pP->imsi.length;
  new_sgw_eps_context->imsi64 = imsi64;
  new_sgw_eps_context->imsi_unauthenticated_indicator = 1;
  new_sgw_eps_context->mme_teid_S11 = session_req_pP->sender_fteid_for_cp.teid;
  new_sgw_eps_context->trxn = session_req_pP->trxn;
  // Update PDN details
  if (check_empty_apn(session_req_pP->apn)) {
    new_sgw_eps_context->pdn_connection.apn_in_use = strdup("NO APN");
  } else {
    new_sgw_eps_context->pdn_connection.apn_in_use =
        strdup(session_req_pP->apn);
  }
  bearer_context_to_be_created_t* csr_bearer_context =
      &session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0];
  new_sgw_eps_context->pdn_connection.default_bearer =
      csr_bearer_context->eps_bearer_id;
  /* creating an eps bearer entry
   * copy informations from create session request to bearer context information
   */

  eps_bearer_ctxt_p = sgw_cm_create_eps_bearer_ctxt_in_collection(
      &new_sgw_eps_context->pdn_connection, csr_bearer_context->eps_bearer_id);
  if (eps_bearer_ctxt_p == nullptr) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to create new EPS bearer entry\n");
    increment_counter("sgw_s8_create_session", 1, 2, "result", "failure",
                      "cause", "internal_software_error");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  eps_bearer_ctxt_p->eps_bearer_qos = csr_bearer_context->bearer_level_qos;
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = sgw_get_new_s1u_teid(sgw_state);
  csr_bearer_context->s5_s8_u_sgw_fteid.ipv4 = 1;
  csr_bearer_context->s5_s8_u_sgw_fteid.ipv4_address =
      sgw_state->sgw_ip_address_S5S8_up;
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

status_code_e sgw_s8_handle_s11_create_session_request(
    sgw_state_t* sgw_state, itti_s11_create_session_request_t* session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_INFO_UE(LOG_SGW_S8, imsi64,
                 "Received S11 CREATE SESSION REQUEST from MME_APP\n");
  sgw_eps_bearer_context_information_t* new_sgw_eps_context = NULL;
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};

  increment_counter("sgw_s8_create_session", 1, NO_LABELS);
  if (session_req_pP->rat_type != RAT_EUTRAN) {
    OAILOG_WARNING_UE(
        LOG_SGW_S8, imsi64,
        "Received session request with RAT != RAT_TYPE_EUTRAN: type %d\n",
        session_req_pP->rat_type);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
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
    increment_counter("sgw_s8_create_session", 1, 2, "result", "failure",
                      "cause", "sender_fteid_incorrect_parameters");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  sgw_s11_tunnel.remote_teid = session_req_pP->sender_fteid_for_cp.teid;
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "Rx CREATE-SESSION-REQUEST MME S11 teid " TEID_FMT
      " APN %s EPS bearer Id %u\n",
      sgw_s11_tunnel.remote_teid, session_req_pP->apn,
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
          .eps_bearer_id);

  uint32_t temporary_create_session_procedure_id = 0;
  new_sgw_eps_context = sgw_create_bearer_context_information_in_collection(
      sgw_state, &temporary_create_session_procedure_id);
  if (!new_sgw_eps_context) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Could not create new sgw context for create session req message for "
        "mme_s11_teid " TEID_FMT "\n",
        session_req_pP->sender_fteid_for_cp.teid);
    increment_counter("sgw_s8_create_session", 1, 2, "result", "failure",
                      "cause", "internal_software_error");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (sgw_update_bearer_context_information_on_csreq(
          sgw_state, new_sgw_eps_context, session_req_pP, imsi64) != RETURNok) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to update sgw_eps_bearer_context_information for "
                    "mme_s11_teid " TEID_FMT "\n",
                    session_req_pP->sender_fteid_for_cp.teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  send_s8_create_session_request(temporary_create_session_procedure_id,
                                 session_req_pP, imsi64);
  sgw_s8_display_s11_bearer_context_information(LOG_SGW_S8,
                                                new_sgw_eps_context);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

int sgw_update_bearer_context_information_on_csrsp(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const s8_create_session_response_t* const session_rsp_p,
    sgw_state_t* sgw_state) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  sgw_eps_bearer_ctxt_t* default_bearer_ctx_p = sgw_s8_cm_get_eps_bearer_entry(
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
  FTEID_T_2_IP_ADDRESS_T((&s5s8_bearer_context.pgw_s8_up),
                         (&default_bearer_ctx_p->p_gw_address_in_use_up));
  default_bearer_ctx_p->p_gw_teid_S5_S8_up = s5s8_bearer_context.pgw_s8_up.teid;
  default_bearer_ctx_p->s_gw_teid_S5_S8_up = session_rsp_p->sgw_s8_up_teid;
  insert_sgw_cp_and_up_teid_to_directoryd(sgw_state, sgw_context_p->imsi64,
                                          UP_TEID);

  memcpy(&default_bearer_ctx_p->eps_bearer_qos, &s5s8_bearer_context.qos,
         sizeof(bearer_qos_t));
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

static int sgw_s8_send_create_session_response(
    sgw_state_t* sgw_state, sgw_eps_bearer_context_information_t* sgw_context_p,
    s8_create_session_response_t* session_rsp_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p = NULL;
  itti_s11_create_session_response_t* create_session_response_p = NULL;

  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_CREATE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_create_session_response \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "sgw_context_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  create_session_response_p = &message_p->ittiMsg.s11_create_session_response;
  if (!create_session_response_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "create_session_response_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  create_session_response_p->teid = sgw_context_p->mme_teid_S11;
  create_session_response_p->cause.cause_value =
      (gtpv2c_cause_value_t)session_rsp_p->cause;
  create_session_response_p->s11_sgw_fteid.teid =
      sgw_context_p->s_gw_teid_S11_S4;
  create_session_response_p->s11_sgw_fteid.interface_type = S11_SGW_GTP_C;
  create_session_response_p->s11_sgw_fteid.ipv4 = 1;
  create_session_response_p->s11_sgw_fteid.ipv4_address.s_addr =
      spgw_config.sgw_config.ipv4.S11.s_addr;

  if (session_rsp_p->cause == REQUEST_ACCEPTED) {
    sgw_eps_bearer_ctxt_t* default_bearer_ctxt_p =
        sgw_s8_cm_get_eps_bearer_entry(&sgw_context_p->pdn_connection,
                                       session_rsp_p->eps_bearer_id);
    if (!default_bearer_ctxt_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, sgw_context_p->imsi64,
          "Failed to get default eps bearer context for sgw_s11_teid " TEID_FMT
          "and bearer_id :%u \n",
          session_rsp_p->context_teid, session_rsp_p->eps_bearer_id);
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    memcpy(&create_session_response_p->paa, &default_bearer_ctxt_p->paa,
           sizeof(paa_t));
    create_session_response_p->bearer_contexts_created.num_bearer_context = 1;
    bearer_context_created_t* bearer_context =
        &create_session_response_p->bearer_contexts_created.bearer_contexts[0];

    bearer_context->s1u_sgw_fteid.teid =
        default_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up;
    bearer_context->s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;

    if (session_rsp_p->pdn_type == IPv4 ||
        session_rsp_p->pdn_type == IPv4_AND_v6) {
      bearer_context->s1u_sgw_fteid.ipv4 = 1;
      bearer_context->s1u_sgw_fteid.ipv4_address.s_addr =
          sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr;
    } else {
      bearer_context->s1u_sgw_fteid.ipv6 = 1;
      memcpy(&bearer_context->s1u_sgw_fteid.ipv6_address,
             &sgw_state->sgw_ipv6_address_S1u_S12_S4_up,
             sizeof(bearer_context->s1u_sgw_fteid.ipv6_address));
    }

    bearer_context->eps_bearer_id = session_rsp_p->eps_bearer_id;
    /*
     * Set the Cause information from bearer context created.
     * "Request accepted" is returned when the GTPv2 entity has accepted a
     * control plane request.
     */
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .cause.cause_value = (gtpv2c_cause_value_t)session_rsp_p->cause;
    if (session_rsp_p->pco.num_protocol_or_container_id) {
      copy_protocol_configuration_options(&create_session_response_p->pco,
                                          &session_rsp_p->pco);
      clear_protocol_configuration_options(&session_rsp_p->pco);
    }
  } else {
    create_session_response_p->bearer_contexts_marked_for_removal
        .num_bearer_context = 1;
    bearer_context_marked_for_removal_t* bearer_context =
        &create_session_response_p->bearer_contexts_marked_for_removal
             .bearer_contexts[0];
    bearer_context->cause.cause_value =
        (gtpv2c_cause_value_t)session_rsp_p->cause;
    bearer_context->eps_bearer_id = session_rsp_p->eps_bearer_id;
    create_session_response_p->trxn = sgw_context_p->trxn;
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

status_code_e sgw_s8_handle_create_session_response(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  if (!session_rsp_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Received null create session response from s8_proxy\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  OAILOG_INFO_UE(LOG_SGW_S8, imsi64,
                 " Rx S5S8_CREATE_SESSION_RSP for context_teid " TEID_FMT "\n",
                 session_rsp_p->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      update_sgw_context_to_s11_teid_map(sgw_state, session_rsp_p, imsi64);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to fetch sgw_eps_bearer_context_info from map of "
                    "temporary_create_session_procedure_id_map for "
                    "temporary_create_session_procedure_id:%u \n",
                    session_rsp_p->temporary_create_session_procedure_id);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (sgw_update_teid_in_ue_context(
          sgw_state, imsi64, session_rsp_p->context_teid) == RETURNerror) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to update sgw_s11_teid" TEID_FMT
                    " in UE context \n",
                    session_rsp_p->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (session_rsp_p->cause == REQUEST_ACCEPTED) {
    // update pdn details received from PGW
    sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp =
        session_rsp_p->pgw_s8_cp_teid.teid;
    // update bearer context details received from PGW
    if ((sgw_update_bearer_context_information_on_csrsp(
            sgw_context_p, session_rsp_p, sgw_state)) != RETURNok) {
      send_s8_delete_session_request(
          sgw_context_p->imsi64, sgw_context_p->imsi,
          sgw_context_p->s_gw_teid_S11_S4,
          sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp,
          sgw_context_p->pdn_connection.default_bearer, NULL);
      session_rsp_p->cause = CONTEXT_NOT_FOUND;
    }
  }
  // send Create session response to mme
  if ((sgw_s8_send_create_session_response(sgw_state, sgw_context_p,
                                           session_rsp_p)) != RETURNok) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
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
    insert_sgw_cp_and_up_teid_to_directoryd(sgw_state, imsi64, CP_TEID);
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

// The function generates comma separated list of sgw's control plane teid of
// s8 interface for each pdn session and updates the list to directoryd.
static void insert_sgw_cp_and_up_teid_to_directoryd(sgw_state_t* sgw_state,
                                                    imsi64_t imsi64,
                                                    uint8_t teid_type) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  char teidString[512] = {0};
  uint32_t teidlist[16] = {0};
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  uint8_t teid_list_idx = 0;
  IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);
  spgw_ue_context_t* ue_context_p = nullptr;

  if (!((teid_type == CP_TEID) || (teid_type == UP_TEID))) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Invalid teid type :%u \n", teid_type);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  sgw_state->imsi_ue_context_map.get(imsi64, &ue_context_p);
  if (!ue_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Failed to find ue context");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }  // end of ue context
  sgw_s11_teid_t* s11_teid_p = nullptr;
  LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
    if (!s11_teid_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64,
          "Failed to find s11 teid entry within ue context's s11 teid list");
      OAILOG_FUNC_OUT(LOG_SGW_S8);
    }
    if (teid_type == CP_TEID) {
      teidlist[teid_list_idx++] = s11_teid_p->sgw_s11_teid;
    } else if (teid_type == UP_TEID) {
      uint8_t idx = 0;
      sgw_eps_bearer_context_information_t* sgw_session_ctxt_p =
          sgw_get_sgw_eps_bearer_context(s11_teid_p->sgw_s11_teid);
      if (sgw_session_ctxt_p) {
        for (; idx < BEARERS_PER_UE; idx++) {
          if (sgw_session_ctxt_p->pdn_connection.sgw_eps_bearers_array[idx]) {
            teidlist[teid_list_idx++] =
                sgw_session_ctxt_p->pdn_connection.sgw_eps_bearers_array[idx]
                    ->s_gw_teid_S5_S8_up;
          }
        }  // end of bearer list
      }  // check for sgw_session_ctxt_p
    }  // switch based on teid type
  }  // end of s11_teid_list

  char separator = ',';
  for (uint8_t idx = 0; idx < teid_list_idx; idx++) {
    if (idx == teid_list_idx - 1) {
      separator = '\0';
    }
    snprintf(teidString + strlen(teidString),
             sizeof(teidString) - strlen(teidString), "%u%c", teidlist[idx],
             separator);
  }
  std::string ptr = "sgw_c_teid";
  if (teid_type == UP_TEID) {
    ptr.assign("sgw_u_teid");
  }
  directoryd_update_record_field(imsi_str, ptr, teidString);
  OAILOG_DEBUG_UE(
      LOG_SGW_S8, imsi64,
      "teid type is %d and teid string to be written to directoryd is %s",
      teid_type, teidString);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void sgw_s8_handle_modify_bearer_request(
    sgw_state_t* state,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  uint8_t idx = 0;
  uint8_t sgi_rsp_idx = 0;
  itti_sgi_update_end_point_response_t sgi_update_end_point_resp = {0};
  struct in_addr enb = {.s_addr = 0};
  struct in6_addr* enb_ipv6 = NULL;
  struct in_addr pgw = {.s_addr = 0};
  struct in6_addr* pgw_ipv6 = NULL;
  sgw_eps_bearer_ctxt_t* bearer_ctx_p = NULL;

  OAILOG_INFO_UE(LOG_SGW_S8, imsi64,
                 "Rx MODIFY_BEARER_REQUEST, teid " TEID_FMT "\n",
                 modify_bearer_pP->teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(modify_bearer_pP->teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
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
  sgw_context_p->trxn = modify_bearer_pP->trxn;
  sgi_update_end_point_resp.context_teid = modify_bearer_pP->teid;
  // Traversing through the list of bearers to be modified
  for (; idx <
         modify_bearer_pP->bearer_contexts_to_be_modified.num_bearer_context;
       idx++) {
    bearer_context_to_be_modified_t mbr_bearer_ctxt_p =
        modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx];
    bearer_ctx_p = sgw_s8_cm_get_eps_bearer_entry(
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

      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
          bearer_ctx_p->enb_ip_address_S1u.pdn_type == IPv6) {
        enb_ipv6 = &bearer_ctx_p->enb_ip_address_S1u.address.ipv6_address;
      }
      pgw.s_addr =
          bearer_ctx_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;
      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
          bearer_ctx_p->p_gw_address_in_use_up.pdn_type == IPv6) {
        pgw_ipv6 = &bearer_ctx_p->p_gw_address_in_use_up.address.ipv6_address;
      }

#if !MME_UNIT_TEST
      // Send end marker to eNB and then delete the tunnel if enb_ip is
      // different
      if (does_sgw_bearer_context_hold_valid_enb_ip(
              bearer_ctx_p->enb_ip_address_S1u) &&
          is_sgw_enb_ip_address_same(&mbr_bearer_ctxt_p.s1_eNB_fteid,
                                     &bearer_ctx_p->enb_ip_address_S1u) ==
              false) {
        struct in_addr ue_ipv4 = bearer_ctx_p->paa.ipv4_address;
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
        gtpv1u_del_s8_tunnel(enb, enb_ipv6, pgw, pgw_ipv6, ue_ipv4, ue_ipv6,
                             bearer_ctx_p->s_gw_teid_S1u_S12_S4_up,
                             bearer_ctx_p->s_gw_teid_S5_S8_up);
      }
#endif
      populate_s8_sgi_end_point_update(sgi_rsp_idx, idx, modify_bearer_pP,
                                       bearer_ctx_p,
                                       &sgi_update_end_point_resp);
      sgi_rsp_idx++;
    }
  }  // for loop

  sgi_rsp_idx = 0;
  for (idx = 0;
       idx < modify_bearer_pP->bearer_contexts_to_be_removed.num_bearer_context;
       idx++) {
    bearer_ctx_p = sgw_s8_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection,
        modify_bearer_pP->bearer_contexts_to_be_removed.bearer_contexts[idx]
            .eps_bearer_id);
    if (bearer_ctx_p) {
      sgi_update_end_point_resp.bearer_contexts_to_be_removed[sgi_rsp_idx++] =
          bearer_ctx_p->eps_bearer_id;
      sgi_update_end_point_resp.num_bearers_removed++;
    }
  }
  sgw_send_modify_bearer_response(sgw_context_p, &sgi_update_end_point_resp,
                                  imsi64);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void sgw_send_modify_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const itti_sgi_update_end_point_response_t* const resp_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  itti_s11_modify_bearer_response_t* modify_response_p = NULL;
  MessageDef* message_p = NULL;

  OAILOG_DEBUG_UE(LOG_SGW_S8, imsi64,
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
    modify_response_p->teid = sgw_context_p->mme_teid_S11;
    modify_response_p->cause.cause_value = REQUEST_ACCEPTED;
    modify_response_p->trxn = sgw_context_p->trxn;
    message_p->ittiMsgHeader.imsi = imsi64;

    sgw_s8_populate_mbr_bearer_contexts_modified(resp_pP, imsi64, sgw_context_p,
                                                 modify_response_p);
    sgw_s8_populate_mbr_bearer_contexts_removed(resp_pP, imsi64, sgw_context_p,
                                                modify_response_p);
    sgw_populate_mbr_bearer_contexts_not_found(LOG_SGW_S8, resp_pP,
                                               modify_response_p);
    send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

static void sgw_s8_populate_mbr_bearer_contexts_modified(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t rsp_idx = 0;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = NULL;

  for (uint8_t idx = 0; idx < resp_pP->num_bearers_modified; idx++) {
    eps_bearer_ctxt_p = sgw_s8_cm_get_eps_bearer_entry(
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

#if !MME_UNIT_TEST
      // setup GTPv1-U tunnels, both s1-u and s8-u tunnels
      sgw_s8_add_gtp_up_tunnel(eps_bearer_ctxt_p, sgw_context_p);
#endif
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

#if !MME_UNIT_TEST
// Helper function to add gtp tunnels for default and dedicated bearers
static int sgw_s8_add_gtp_up_tunnel(
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    sgw_eps_bearer_context_information_t* sgw_context_p) {
  int rv = RETURNok;
  struct in_addr enb = {.s_addr = 0};
  struct in6_addr* enb_ipv6 = NULL;
  struct in_addr pgw = {.s_addr = 0};
  struct in6_addr* pgw_ipv6 = NULL;
  pgw.s_addr =
      eps_bearer_ctxt_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;
  if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
      eps_bearer_ctxt_p->p_gw_address_in_use_up.pdn_type == IPv6) {
    pgw_ipv6 = &eps_bearer_ctxt_p->p_gw_address_in_use_up.address.ipv6_address;
  }
  if ((pgw.s_addr == 0) && (eps_bearer_ctxt_p->p_gw_teid_S5_S8_up == 0)) {
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "bearer context has invalid pgw_s8_teid " TEID_FMT
                    "pgw_ip address :%x \n",
                    eps_bearer_ctxt_p->p_gw_teid_S5_S8_up, pgw.s_addr);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  enb.s_addr =
      eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr;
  if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
      eps_bearer_ctxt_p->enb_ip_address_S1u.pdn_type == IPv6) {
    enb_ipv6 = &eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv6_address;
  }

  struct in_addr ue_ipv4 = {.s_addr = 0};
  struct in6_addr* ue_ipv6 = NULL;
  ue_ipv4.s_addr = eps_bearer_ctxt_p->paa.ipv4_address.s_addr;
  if ((eps_bearer_ctxt_p->paa.pdn_type == IPv6) ||
      (eps_bearer_ctxt_p->paa.pdn_type == IPv4_AND_v6)) {
    ue_ipv6 = &eps_bearer_ctxt_p->paa.ipv6_address;
  }

  int vlan = eps_bearer_ctxt_p->paa.vlan;
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
        ue_ipv4, ue_ipv6, vlan, enb, enb_ipv6, pgw, pgw_ipv6,
        eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
        eps_bearer_ctxt_p->enb_teid_S1u, eps_bearer_ctxt_p->s_gw_teid_S5_S8_up,
        eps_bearer_ctxt_p->p_gw_teid_S5_S8_up, imsi);
    if (rv < 0) {
      OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                      "ERROR in setting up TUNNEL err=%d\n", rv);
    }
  } else {
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
    for (int i = 0; i < eps_bearer_ctxt_p->tft.numberofpacketfilters; ++i) {
      rv = gtpv1u_add_s8_tunnel(ue_ipv4, ue_ipv6, vlan, enb, enb_ipv6, pgw,
                                pgw_ipv6,
                                eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
                                eps_bearer_ctxt_p->enb_teid_S1u,
                                eps_bearer_ctxt_p->s_gw_teid_S5_S8_up,
                                eps_bearer_ctxt_p->p_gw_teid_S5_S8_up, imsi);
      if (rv < 0) {
        OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                        "ERROR in setting up TUNNEL err=%d\n", rv);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, rv);
}
#endif

status_code_e sgw_s8_handle_s11_delete_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  gtpv2c_cause_value_t gtpv2c_cause = (gtpv2c_cause_value_t)0;
  if (!delete_session_req_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Received NULL delete_session_req_p from mme app \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Received S11 DELETE SESSION REQUEST for sgw_s11_teid " TEID_FMT "\n",
      delete_session_req_p->teid);
  increment_counter("sgw_delete_session", 1, NO_LABELS);
  if (delete_session_req_p->indication_flags.oi) {
    OAILOG_DEBUG_UE(LOG_SGW_S8, imsi64,
                    "OI flag is set for this message indicating the request"
                    "should be forwarded to P-GW entity\n");
  }

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(delete_session_req_p->teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to fetch sgw_eps_bearer_context_info from "
                    "sgw_s11_teid " TEID_FMT " \n",
                    delete_session_req_p->teid);
    gtpv2c_cause = CONTEXT_NOT_FOUND;
    sgw_s8_send_failed_delete_session_response(
        sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if ((delete_session_req_p->sender_fteid_for_cp.ipv4) &&
      (delete_session_req_p->sender_fteid_for_cp.ipv6)) {
    // Sender F-TEID IE present
    if (delete_session_req_p->teid != sgw_context_p->mme_teid_S11) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64,
          "Mismatch in MME Teid for CP teid recevied in delete session "
          "req: " TEID_FMT " teid present in sgw_context :" TEID_FMT "\n",
          delete_session_req_p->teid, sgw_context_p->mme_teid_S11);
      gtpv2c_cause = INVALID_PEER;
      sgw_s8_send_failed_delete_session_response(
          sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
  }
  if (delete_session_req_p->lbi !=
      sgw_context_p->pdn_connection.default_bearer) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Mismatch in default eps bearer_id, bearer_id recevied in delete "
        "session req :%d and bearer_id present in sgw_context :%d \n",
        delete_session_req_p->lbi,
        sgw_context_p->pdn_connection.default_bearer);
    sgw_s8_send_failed_delete_session_response(
        sgw_context_p, gtpv2c_cause, sgw_state, delete_session_req_p, imsi64);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  send_s8_delete_session_request(
      sgw_context_p->imsi64, sgw_context_p->imsi,
      sgw_context_p->s_gw_teid_S11_S4,
      sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp,
      sgw_context_p->pdn_connection.default_bearer, delete_session_req_p);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

#if !MME_UNIT_TEST
static void delete_userplane_tunnels(
    sgw_eps_bearer_context_information_t* sgw_context_p) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  struct in_addr enb = {.s_addr = 0};
  struct in6_addr* enb_ipv6 = NULL;
  struct in_addr pgw = {.s_addr = 0};
  struct in6_addr* pgw_ipv6 = NULL;
  sgw_eps_bearer_ctxt_t* bearer_ctxt_p = NULL;
  int rv = RETURNerror;
  struct in_addr ue_ipv4 = {.s_addr = 0};

  for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
    ebi_t ebi = INDEX_TO_EBI(ebix);
    bearer_ctxt_p =
        sgw_s8_cm_get_eps_bearer_entry(&sgw_context_p->pdn_connection, ebi);

    if (bearer_ctxt_p) {
      enb.s_addr =
          bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr;

      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
          bearer_ctxt_p->enb_ip_address_S1u.pdn_type == IPv6) {
        enb_ipv6 = &bearer_ctxt_p->enb_ip_address_S1u.address.ipv6_address;
      }
      pgw.s_addr =
          bearer_ctxt_p->p_gw_address_in_use_up.address.ipv4_address.s_addr;
      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
          bearer_ctxt_p->p_gw_address_in_use_up.pdn_type == IPv6) {
        pgw_ipv6 = &bearer_ctxt_p->p_gw_address_in_use_up.address.ipv6_address;
      }
      struct in6_addr* ue_ipv6 = NULL;
      if ((bearer_ctxt_p->paa.pdn_type == IPv6) ||
          (bearer_ctxt_p->paa.pdn_type == IPv4_AND_v6)) {
        ue_ipv6 = &bearer_ctxt_p->paa.ipv6_address;
      }
      ue_ipv4 = bearer_ctxt_p->paa.ipv4_address;
      // Delete S1-U tunnel and S8-U tunnel
      rv = gtpv1u_del_s8_tunnel(enb, enb_ipv6, pgw, pgw_ipv6, ue_ipv4, ue_ipv6,
                                bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
                                bearer_ctxt_p->s_gw_teid_S5_S8_up);
      if (rv < 0) {
        OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                        "ERROR in deleting S1-U TUNNEL " TEID_FMT
                        " (eNB) <-> (SGW) " TEID_FMT "\n",
                        bearer_ctxt_p->enb_teid_S1u,
                        bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up);
      }
      // delete paging rule
      char* ip_str = inet_ntoa(ue_ipv4);
      rv = gtp_tunnel_ops->delete_paging_rule(ue_ipv4, ue_ipv6);
      if (rv < 0) {
        OAILOG_ERROR(LOG_SGW_S8,
                     "ERROR in deleting paging rule for IP Addr: %s\n", ip_str);
      } else {
        OAILOG_DEBUG(LOG_SGW_S8, "Stopped paging for IP Addr: %s\n", ip_str);
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}
#endif

status_code_e sgw_s8_handle_delete_session_response(
    sgw_state_t* sgw_state, s8_delete_session_response_t* session_rsp_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p = NULL;
  itti_s11_delete_session_response_t* delete_session_response_p = NULL;

  if (!session_rsp_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Received null delete session response from s8_proxy\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      " Rx S5S8_DELETE_SESSION_RSP from s8_proxy for context_teid " TEID_FMT
      "\n",
      session_rsp_p->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(session_rsp_p->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to fetch sgw_eps_bearer_context_info from "
                    "context_teid " TEID_FMT " \n",
                    session_rsp_p->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  message_p = itti_alloc_new_message(TASK_SGW_S8, S11_DELETE_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_CRITICAL_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to allocate memory for S11_delete_session_response for "
        "context_teid " TEID_FMT "\n",
        session_rsp_p->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  delete_session_response_p = &message_p->ittiMsg.s11_delete_session_response;
  message_p->ittiMsgHeader.imsi = imsi64;
  if (!delete_session_response_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "delete_session_response_p is NULL \n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  delete_session_response_p->teid = sgw_context_p->mme_teid_S11;
  delete_session_response_p->cause.cause_value =
      (gtpv2c_cause_value_t)session_rsp_p->cause;
  delete_session_response_p->trxn = sgw_context_p->trxn;
  delete_session_response_p->lbi = sgw_context_p->pdn_connection.default_bearer;

#if !MME_UNIT_TEST
  // Delete ovs rules
  delete_userplane_tunnels(sgw_context_p);
#endif
  sgw_remove_sgw_bearer_context_information(
      sgw_state, session_rsp_p->context_teid, imsi64);
  // send delete session response to mme
  if (send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to send delete session response to mme for "
                    "sgw_s11_teid " TEID_FMT "\n",
                    session_rsp_p->context_teid);
    increment_counter("sgw_s8_delete_session", 1, 1, "result", "failed");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  increment_counter("sgw_s8_delete_session", 1, 1, "result", "success");
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

static void sgw_s8_send_failed_delete_session_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  MessageDef* message_p = NULL;
  itti_s11_delete_session_response_t* delete_session_response_p = NULL;
  teid_t teid = 0;

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
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "delete_session_response_p is NULL \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  if (sgw_context_p) {
    delete_session_response_p->teid = sgw_context_p->mme_teid_S11;
    delete_session_response_p->cause.cause_value = cause;
    delete_session_response_p->trxn = sgw_context_p->trxn;
    delete_session_response_p->lbi =
        sgw_context_p->pdn_connection.default_bearer;

    // Delete ovs rules
#if !MME_UNIT_TEST
    delete_userplane_tunnels(sgw_context_p);
#endif
    sgw_remove_sgw_bearer_context_information(
        sgw_state, sgw_context_p->s_gw_teid_S11_S4, imsi64);
    teid = sgw_context_p->s_gw_teid_S11_S4;
  } else {
    if (delete_session_req_p) {
      delete_session_response_p->teid = delete_session_req_p->local_teid;
      delete_session_response_p->cause.cause_value = cause;
      delete_session_response_p->trxn = delete_session_req_p->trxn;
      delete_session_response_p->lbi = delete_session_req_p->lbi;
      teid = delete_session_req_p->teid;
    }
  }
  increment_counter("sgw_delete_session", 1, 1, "result", "failed");
  // send delete session response to mme
  if (send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
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
    sgw_state_t* sgw_state,
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  OAILOG_DEBUG_UE(LOG_SGW_S8, imsi64,
                  "Release Access Bearer Request Received in SGW_S8 task");

  spgw_ue_context_t* ue_context_p = nullptr;
  gtpv2c_cause_value_t cause = CONTEXT_NOT_FOUND;

  sgw_state->imsi_ue_context_map.get(imsi64, &ue_context_p);

  if (!ue_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "Failed to find ue context");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }  // end of ue context
  sgw_s11_teid_t* s11_teid_p = nullptr;
  LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
    if (!s11_teid_p) {
      OAILOG_ERROR_UE(
          LOG_SGW_S8, imsi64,
          "Failed to find s11 teid entry within ue context's s11 teid list");
      OAILOG_FUNC_OUT(LOG_SGW_S8);
    }
    sgw_eps_bearer_context_information_t* sgw_session_ctxt_p =
        sgw_get_sgw_eps_bearer_context(s11_teid_p->sgw_s11_teid);
    if (sgw_session_ctxt_p) {
      sgw_s8_process_release_access_bearer_request(LOG_SGW_S8, imsi64,
                                                   sgw_session_ctxt_p);
      cause = REQUEST_ACCEPTED;
    }
  }
  sgw_send_release_access_bearer_response(
      LOG_SGW_S8, TASK_SGW_S8, imsi64, cause, release_access_bearers_req_pP,
      release_access_bearers_req_pP->local_teid);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

int update_pgw_info_to_temp_dedicated_bearer_context(
    sgw_eps_bearer_context_information_t* sgw_context_p, teid_t s1_u_sgw_fteid,
    s8_bearer_context_t* bc_cbreq, sgw_state_t* sgw_state, char* pgw_cp_ip_port,
    teid_t sgw_s8_u_teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_context_p);
  if (!pgw_ni_cbr_proc) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, sgw_context_p->imsi64,
        "Failed to get Create bearer procedure from temporary stored contexts "
        "for lbi :%u \n",
        bc_cbreq->eps_bearer_id);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  struct sgw_eps_bearer_entry_wrapper_s* spgw_eps_bearer_entry_p =
      LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  while (spgw_eps_bearer_entry_p &&
         spgw_eps_bearer_entry_p->sgw_eps_bearer_entry) {
    if (s1_u_sgw_fteid == spgw_eps_bearer_entry_p->sgw_eps_bearer_entry
                              ->s_gw_teid_S1u_S12_S4_up) {
      // update PGW teid and ip address
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->p_gw_teid_S5_S8_up =
          bc_cbreq->pgw_s8_up.teid;
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->p_gw_address_in_use_up
          .address.ipv4_address.s_addr =
          bc_cbreq->pgw_s8_up.ipv4_address.s_addr;
      if (pgw_cp_ip_port) {
        uint8_t pgw_ip_port_len = strlen(pgw_cp_ip_port);
        spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->pgw_cp_ip_port =
            reinterpret_cast<char*>(calloc(1, pgw_ip_port_len + 1));
        memcpy(spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->pgw_cp_ip_port,
               pgw_cp_ip_port, pgw_ip_port_len);
        spgw_eps_bearer_entry_p->sgw_eps_bearer_entry
            ->pgw_cp_ip_port[pgw_ip_port_len] = '\0';
      }
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S5_S8_up =
          sgw_s8_u_teid;
      break;
    }
    spgw_eps_bearer_entry_p = LIST_NEXT(spgw_eps_bearer_entry_p, entries);
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

imsi64_t sgw_s8_handle_create_bearer_request(
    sgw_state_t* sgw_state, const s8_create_bearer_request_t* const cb_req,
    gtpv2c_cause_value_t* cause_value) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint8_t bearer_idx = 0;

  if (!cb_req) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Received null create bearer request from s8_proxy\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }
  OAILOG_INFO(
      LOG_SGW_S8,
      "Received S8_CREATE_BEARER_REQ from s8_proxy for context_teid " TEID_FMT
      "\n",
      cb_req->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(cb_req->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to fetch sgw_eps_bearer_context_info from "
                 "context_teid " TEID_FMT " \n",
                 cb_req->context_teid);
    *cause_value = CONTEXT_NOT_FOUND;
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }

  if (sgw_context_p->pdn_connection.default_bearer !=
      cb_req->linked_eps_bearer_id) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "No matching lbi found for context_teid: " TEID_FMT
        "lbi within create bearer request: %u, lbi with sgw_context: %u "
        "Sending dedicated_bearer_actv_rsp with REQUEST_REJECTED cause to NW\n",
        cb_req->context_teid, cb_req->linked_eps_bearer_id,
        sgw_context_p->pdn_connection.default_bearer);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req->bearer_context[bearer_idx];

  sgw_eps_bearer_ctxt_t* bearer_ctxt_p = NULL;
  bearer_ctxt_p = sgw_s8_cm_get_eps_bearer_entry(&sgw_context_p->pdn_connection,
                                                 cb_req->linked_eps_bearer_id);
  if (bearer_ctxt_p == NULL) {
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "Failed to retrieve bearer ctxt\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  itti_bearer_req.lbi = cb_req->linked_eps_bearer_id;

  memcpy(&itti_bearer_req.ul_tft, &bc_cbreq.tft,
         sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.dl_tft, &bc_cbreq.tft,
         sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  int rc = create_temporary_s8_dedicated_bearer_context(
      sgw_context_p, &itti_bearer_req,
      bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.pdn_type,
      sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
      &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
      cb_req->sequence_number, LOG_SGW_S8);
  if (rc != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to create temporary dedicated bearer context for lbi: %u"
        " and context_teid " TEID_FMT "\n ",
        itti_bearer_req.lbi, cb_req->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }

  rc = update_pgw_info_to_temp_dedicated_bearer_context(
      sgw_context_p, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
      cb_req->pgw_cp_address, cb_req->sgw_s8_up_teid);
  free_wrapper((void**)&cb_req->pgw_cp_address);
  if (rc != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to update PGW info to temporary dedicated bearer context for "
        "lbi %u and context_teid " TEID_FMT " \n ",
        itti_bearer_req.lbi, cb_req->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }

  if (sgw_build_and_send_s11_create_bearer_request(
          sgw_context_p, &itti_bearer_req,
          bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.pdn_type,
          sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
          &sgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
          LOG_SGW_S8) != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, sgw_context_p->imsi64,
        "Failed to send create bearer request from s8_proxy for lbi:%u "
        "context_teid:" TEID_FMT " \n",
        itti_bearer_req.lbi, cb_req->context_teid);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, INVALID_IMSI64);
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_context_p->imsi64);
}

void sgw_s8_proc_s11_create_bearer_rsp(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    bearer_context_within_create_bearer_response_t* bc_cbrsp,
    itti_s11_nw_init_actv_bearer_rsp_t* s11_actv_bearer_rsp, imsi64_t imsi64,
    sgw_state_t* sgw_state) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  struct sgw_eps_bearer_entry_wrapper_s* sgw_eps_bearer_entry_p = nullptr;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = nullptr;
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = nullptr;
  pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(sgw_context_p);

  if (!pgw_ni_cbr_proc) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to get create bearer procedure from temporary stored context, "
        "so did not create new EPS bearer entry for EBI %u for "
        "sgw_s11_teid " TEID_FMT "\n",
        bc_cbrsp->eps_bearer_id, s11_actv_bearer_rsp->sgw_s11_teid);
    handle_failed_s8_create_bearer_response(
        sgw_context_p, s11_actv_bearer_rsp->cause.cause_value, imsi64, bc_cbrsp,
        nullptr, LOG_SGW_S8);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }

  sgw_eps_bearer_entry_p = LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  while (sgw_eps_bearer_entry_p) {
    if (bc_cbrsp->s1u_sgw_fteid.teid ==
        sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up) {
      eps_bearer_ctxt_p = sgw_eps_bearer_entry_p->sgw_eps_bearer_entry;
      if (eps_bearer_ctxt_p) {
        eps_bearer_ctxt_p->eps_bearer_id = bc_cbrsp->eps_bearer_id;

        // Store enb-s1u teid and ip address
        get_fteid_ip_address(&bc_cbrsp->s1u_enb_fteid,
                             &eps_bearer_ctxt_p->enb_ip_address_S1u);
        eps_bearer_ctxt_p->enb_teid_S1u = bc_cbrsp->s1u_enb_fteid.teid;
        sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_entry_p =
            sgw_cm_insert_eps_bearer_ctxt_in_collection(
                &sgw_context_p->pdn_connection, eps_bearer_ctxt_p);
        if (eps_bearer_ctxt_entry_p == nullptr) {
          OAILOG_ERROR_UE(
              LOG_SGW_S8, imsi64,
              "Failed to create new EPS bearer entry for bearer_id :%u \n",
              eps_bearer_ctxt_p->eps_bearer_id);
          increment_counter("s11_actv_bearer_rsp", 1, 2, "result", "failure",
                            "cause", "internal_software_error");
        } else {
          OAILOG_INFO_UE(
              LOG_SGW_S8, imsi64,
              "Successfully created new EPS bearer entry with EBI %d ip:%x "
              "enb_s1u_teid :" TEID_FMT "\n",
              eps_bearer_ctxt_p->eps_bearer_id,
              eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr,
              eps_bearer_ctxt_p->enb_teid_S1u);
#if !MME_UNIT_TEST
          if (sgw_s8_add_gtp_up_tunnel(eps_bearer_ctxt_p, sgw_context_p) ==
              RETURNerror) {
            OAILOG_ERROR_UE(
                LOG_SGW_S8, imsi64,
                "Failed to create OVS rules for dedicated bearer with "
                "bearer_id:%u \n",
                eps_bearer_ctxt_p->eps_bearer_id);
            bc_cbrsp->cause.cause_value = REQUEST_REJECTED;
          }
#endif
          bc_cbrsp->cause.cause_value = REQUEST_ACCEPTED;
        }
      }
      // Remove the temporary spgw entry
      LIST_REMOVE(sgw_eps_bearer_entry_p, entries);
      free_cpp_wrapper(reinterpret_cast<void**>(&sgw_eps_bearer_entry_p));
      break;
    }
    sgw_eps_bearer_entry_p = LIST_NEXT(sgw_eps_bearer_entry_p, entries);
  }
  if (pgw_ni_cbr_proc && (LIST_EMPTY(pgw_ni_cbr_proc->pending_eps_bearers))) {
    pgw_base_proc_t* base_proc1 = LIST_FIRST(sgw_context_p->pending_procedures);
    LIST_REMOVE(base_proc1, entries);
    free_cpp_wrapper(
        reinterpret_cast<void**>(&sgw_context_p->pending_procedures));
    free_cpp_wrapper(
        reinterpret_cast<void**>(&pgw_ni_cbr_proc->pending_eps_bearers));
    pgw_free_procedure_create_bearer((pgw_ni_cbr_proc_t**)&pgw_ni_cbr_proc);
  }
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void sgw_s8_handle_s11_create_bearer_response(
    sgw_state_t* sgw_state,
    itti_s11_nw_init_actv_bearer_rsp_t* s11_actv_bearer_rsp, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  uint32_t msg_bearer_index = 0;
  sgw_eps_bearer_ctxt_t dedicated_bearer_ctxt = {};
  bearer_context_within_create_bearer_response_t bc_cbrsp = {0};

  if (!s11_actv_bearer_rsp) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Received null itti:s11_actv_bearer_rsp message from MME \n");
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  bc_cbrsp =
      s11_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index];
  OAILOG_INFO_UE(LOG_SGW_S8, imsi64,
                 "Received S11_create_bearer_response from MME with EBI %u for "
                 "sgw_s11_teid " TEID_FMT "\n",
                 bc_cbrsp.eps_bearer_id, s11_actv_bearer_rsp->sgw_s11_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(s11_actv_bearer_rsp->sgw_s11_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to retrieve sgw_context from sgw_s11_teid " TEID_FMT
                    "\n",
                    s11_actv_bearer_rsp->sgw_s11_teid);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  // If UE did not accept the request send reject to NW
  if (s11_actv_bearer_rsp->cause.cause_value != REQUEST_ACCEPTED) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Did not create new EPS bearer entry as "
        "UE rejected the request for EBI %u for sgw_s11_teid " TEID_FMT "\n",
        bc_cbrsp.eps_bearer_id, s11_actv_bearer_rsp->sgw_s11_teid);
    handle_failed_s8_create_bearer_response(
        sgw_context_p, s11_actv_bearer_rsp->cause.cause_value, imsi64,
        &bc_cbrsp, &dedicated_bearer_ctxt, LOG_SGW_S8);
    sgw_s8_send_failed_create_bearer_response(
        sgw_state, dedicated_bearer_ctxt.sgw_sequence_number,
        dedicated_bearer_ctxt.pgw_cp_ip_port,
        s11_actv_bearer_rsp->cause.cause_value, sgw_context_p->imsi,
        sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp);
    OAILOG_FUNC_OUT(LOG_SGW_S8);
  }
  sgw_s8_proc_s11_create_bearer_rsp(sgw_context_p, &bc_cbrsp,
                                    s11_actv_bearer_rsp, imsi64, sgw_state);

  insert_sgw_cp_and_up_teid_to_directoryd(sgw_state, imsi64, UP_TEID);
  char* pgw_cp_ip_port = NULL;
  sgw_eps_bearer_ctxt_t* dedicated_bearer_ctx_p = NULL;
  for (uint8_t idx = 0;
       idx < s11_actv_bearer_rsp->bearer_contexts.num_bearer_context; idx++) {
    bearer_context_within_create_bearer_response_t* bc_cbresp_msg =
        &s11_actv_bearer_rsp->bearer_contexts.bearer_contexts[idx];
    dedicated_bearer_ctx_p = sgw_s8_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, bc_cbresp_msg->eps_bearer_id);
    if (!dedicated_bearer_ctx_p) {
      OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                      "Failed to get dedicated eps bearer context for context "
                      "teid " TEID_FMT "and bearer_id :%u \n",
                      s11_actv_bearer_rsp->sgw_s11_teid,
                      bc_cbrsp.eps_bearer_id);
      OAILOG_FUNC_OUT(LOG_SGW_S8);
    }
    bc_cbresp_msg->s5_s8_u_sgw_fteid.teid =
        dedicated_bearer_ctx_p->s_gw_teid_S5_S8_up;
    bc_cbresp_msg->s5_s8_u_sgw_fteid.ipv4 = 1;
    bc_cbresp_msg->s5_s8_u_sgw_fteid.ipv4_address =
        sgw_state->sgw_ip_address_S5S8_up;
    memcpy(&bc_cbresp_msg->bearer_level_qos,
           &dedicated_bearer_ctx_p->eps_bearer_qos, sizeof(bearer_qos_t));
    pgw_cp_ip_port = dedicated_bearer_ctx_p->pgw_cp_ip_port;

    bc_cbresp_msg->s5_s8_u_pgw_fteid.teid =
        dedicated_bearer_ctx_p->p_gw_teid_S5_S8_up;
    bc_cbresp_msg->s5_s8_u_pgw_fteid.ipv4 = 1;
    bc_cbresp_msg->s5_s8_u_pgw_fteid.ipv4_address =
        dedicated_bearer_ctx_p->p_gw_address_in_use_up.address.ipv4_address;
  }

  send_s8_create_bearer_response(
      s11_actv_bearer_rsp, sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp,
      dedicated_bearer_ctx_p->sgw_sequence_number, pgw_cp_ip_port,
      sgw_context_p->imsi);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

void sgw_s8_send_failed_create_bearer_response(
    sgw_state_t* sgw_state, uint32_t sequence_number, char* pgw_cp_address,
    gtpv2c_cause_value_t cause_value, Imsi_t imsi, teid_t pgw_s8_teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp = {};
  s11_actv_bearer_rsp.cause.cause_value = cause_value;
  s11_actv_bearer_rsp.bearer_contexts.num_bearer_context = 1;
  s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0].cause.cause_value =
      cause_value;
  send_s8_create_bearer_response(&s11_actv_bearer_rsp, pgw_s8_teid,
                                 sequence_number, pgw_cp_address, imsi);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

int sgw_s8_handle_delete_bearer_request(
    sgw_state_t* sgw_state, const s8_delete_bearer_request_t* const db_req) {
  OAILOG_FUNC_IN(LOG_SGW_S8);

  ebi_t ebi_to_be_deactivated[BEARERS_PER_UE] = {0};
  bool is_ebi_found = false;
  uint32_t no_of_bearers_to_be_deact = 0;
  uint32_t no_of_bearers_rej = 0;
  ebi_t invalid_bearer_id[BEARERS_PER_UE] = {0};

  if (!db_req) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Received null delete bearer request from s8_proxy\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  OAILOG_INFO(
      LOG_SGW_S8,
      "Received S8_DELETE_BEARER_REQ from s8_proxy for context_teid " TEID_FMT
      "\n",
      db_req->context_teid);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(db_req->context_teid);
  if (!sgw_context_p) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to fetch sgw_eps_bearer_context_info from "
                 "context_teid " TEID_FMT " \n",
                 db_req->context_teid);
    Imsi_t imsi = {0};
    sgw_s8_send_failed_delete_bearer_response(db_req, CONTEXT_NOT_FOUND, imsi,
                                              0);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  // Check if the received EBI is valid
  for (uint8_t idx = 0; idx < db_req->num_eps_bearer_id; idx++) {
    sgw_eps_bearer_ctxt_t* bearer_ctxt_p = sgw_s8_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, db_req->eps_bearer_id[idx]);
    if (bearer_ctxt_p) {
      bearer_ctxt_p->sgw_sequence_number = db_req->sequence_number;
      is_ebi_found = true;
      ebi_to_be_deactivated[no_of_bearers_to_be_deact] =
          db_req->eps_bearer_id[idx];
      no_of_bearers_to_be_deact++;
    } else {
      invalid_bearer_id[no_of_bearers_rej] = db_req->eps_bearer_id[idx];
      no_of_bearers_rej++;
    }
  }
  /* Send reject to NW if we did not find ebi/lbi
   * Also in case of multiple bearers, if some EBIs are valid and some are not,
   * send reject to those for which we did not find EBI.
   * Proceed with deactivation by sending s5_nw_init_deactv_bearer_request to
   * SGW for valid EBIs
   */
  if ((!is_ebi_found) || (no_of_bearers_rej > 0)) {
    OAILOG_INFO_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                   "is_ebi_found: %d no_of_bearers_rej: %d\n", is_ebi_found,
                   no_of_bearers_rej);
    OAILOG_ERROR_UE(LOG_SGW_S8, sgw_context_p->imsi64,
                    "Sending dedicated bearer deactivation reject to NW\n");
    print_bearer_ids_helper(invalid_bearer_id, no_of_bearers_rej);
    sgw_s8_send_failed_delete_bearer_response(
        db_req, REQUEST_REJECTED, sgw_context_p->imsi,
        sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  if (no_of_bearers_to_be_deact > 0) {
    bool delete_default_bearer =
        (sgw_context_p->pdn_connection.default_bearer ==
         db_req->eps_bearer_id[0])
            ? true
            : false;
    spgw_build_and_send_s11_deactivate_bearer_req(
        sgw_context_p->imsi64, no_of_bearers_to_be_deact, ebi_to_be_deactivated,
        delete_default_bearer, sgw_context_p->mme_teid_S11, LOG_SGW_S8);
  }
  OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNok);
}

// Handle NW-initiated dedicated bearer dectivation response from MME
status_code_e sgw_s8_handle_s11_delete_bearer_response(
    sgw_state_t* sgw_state,
    const itti_s11_nw_init_deactv_bearer_rsp_t* const
        s11_delete_bearer_response_p,
    imsi64_t imsi64) {
  status_code_e rc = RETURNok;
  ebi_t ebi = {0};
  uint32_t sequence_number = 0;
  char* pgw_cp_ip_port = NULL;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = NULL;

  if (!s11_delete_bearer_response_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Received null delete bearer response from MME\n");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }
  OAILOG_INFO_UE(
      LOG_SGW_S8, imsi64,
      "Received S11 delete bearer response from MME for teid " TEID_FMT "\n",
      s11_delete_bearer_response_p->s_gw_teid_s11_s4);

  sgw_eps_bearer_context_information_t* sgw_context_p =
      sgw_get_sgw_eps_bearer_context(
          s11_delete_bearer_response_p->s_gw_teid_s11_s4);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(
        LOG_SGW_S8, imsi64,
        "Failed to fetch sgw_eps_bearer_context_info from teid " TEID_FMT "\n",
        s11_delete_bearer_response_p->s_gw_teid_s11_s4);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
  }

  /* If delete bearer request is initiated for default bearer
   * then delete all the dedicated bearers linked to this default bearer
   */
  if (s11_delete_bearer_response_p->delete_default_bearer) {
    if (!s11_delete_bearer_response_p->lbi) {
      OAILOG_ERROR_UE(LOG_SGW_S8, imsi64, "LBI received from MME is NULL\n");
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    eps_bearer_ctxt_p = sgw_s8_cm_get_eps_bearer_entry(
        &sgw_context_p->pdn_connection, *(s11_delete_bearer_response_p->lbi));
    if (!eps_bearer_ctxt_p) {
      OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                      "Failed to get bearer context for bearer_id :%u\n",
                      *(s11_delete_bearer_response_p->lbi));
      OAILOG_FUNC_RETURN(LOG_SGW_S8, RETURNerror);
    }
    sequence_number = eps_bearer_ctxt_p->sgw_sequence_number;
    uint8_t pgw_ip_port_len = strlen(eps_bearer_ctxt_p->pgw_cp_ip_port) + 1;
    pgw_cp_ip_port = reinterpret_cast<char*>(calloc(1, pgw_ip_port_len));
    memcpy(pgw_cp_ip_port, eps_bearer_ctxt_p->pgw_cp_ip_port, pgw_ip_port_len);
#if !MME_UNIT_TEST
    // Delete ovs rules
    delete_userplane_tunnels(sgw_context_p);
#endif
    sgw_remove_sgw_bearer_context_information(
        sgw_state, s11_delete_bearer_response_p->s_gw_teid_s11_s4, imsi64);

    OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                   "Remove default bearer context for (ebi = %u)\n",
                   *(s11_delete_bearer_response_p->lbi));
  } else {
    // Remove the dedicated bearer/s context
    uint32_t no_of_bearers =
        s11_delete_bearer_response_p->bearer_contexts.num_bearer_context;
    for (uint8_t i = 0; i < no_of_bearers; i++) {
      ebi = s11_delete_bearer_response_p->bearer_contexts.bearer_contexts[i]
                .eps_bearer_id;
      eps_bearer_ctxt_p =
          sgw_s8_cm_get_eps_bearer_entry(&sgw_context_p->pdn_connection, ebi);
      if (eps_bearer_ctxt_p) {
        OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                       "Removed bearer context for (ebi = %u)\n", ebi);
        struct in6_addr* ue_ipv6 = NULL;
        if ((eps_bearer_ctxt_p->paa.pdn_type == IPv6) ||
            (eps_bearer_ctxt_p->paa.pdn_type == IPv4_AND_v6)) {
          ue_ipv6 = &eps_bearer_ctxt_p->paa.ipv6_address;
        }
        struct in_addr ue_ipv4 = eps_bearer_ctxt_p->paa.ipv4_address;
        struct in_addr enb = {.s_addr = 0};
        struct in6_addr* enb_ipv6 = NULL;
        struct in_addr pgw = {.s_addr = 0};
        struct in6_addr* pgw_ipv6 = NULL;
        enb.s_addr =
            eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv4_address.s_addr;

        if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
            eps_bearer_ctxt_p->enb_ip_address_S1u.pdn_type == IPv6) {
          enb_ipv6 =
              &eps_bearer_ctxt_p->enb_ip_address_S1u.address.ipv6_address;
        }
        pgw.s_addr = eps_bearer_ctxt_p->p_gw_address_in_use_up.address
                         .ipv4_address.s_addr;
        if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
            eps_bearer_ctxt_p->p_gw_address_in_use_up.pdn_type == IPv6) {
          pgw_ipv6 =
              &eps_bearer_ctxt_p->p_gw_address_in_use_up.address.ipv6_address;
        }
        OAILOG_INFO_UE(
            LOG_SGW_S8, imsi64,
            "Successfully created new EPS bearer entry with enb_ip:%x "
            "pgw_ip :%x"
            "sgw_enb_s1u_teid :" TEID_FMT "sgw_s5s8u_teid " TEID_FMT "\n",
            enb.s_addr, pgw.s_addr, eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
            eps_bearer_ctxt_p->s_gw_teid_S5_S8_up);

#if !MME_UNIT_TEST
        if (RETURNok != (gtpv1u_del_s8_tunnel(
                            enb, enb_ipv6, pgw, pgw_ipv6, ue_ipv4, ue_ipv6,
                            eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up,
                            eps_bearer_ctxt_p->s_gw_teid_S5_S8_up))) {
          OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                          "ERROR in deleting TUNNEL " TEID_FMT
                          " (eNB) <-> (SGW) " TEID_FMT "\n",
                          eps_bearer_ctxt_p->s_gw_teid_S5_S8_up,
                          eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up);
        }
#endif
        uint8_t pgw_ip_port_len = strlen(eps_bearer_ctxt_p->pgw_cp_ip_port) + 1;
        pgw_cp_ip_port = reinterpret_cast<char*>(calloc(1, pgw_ip_port_len));
        memcpy(pgw_cp_ip_port, eps_bearer_ctxt_p->pgw_cp_ip_port,
               pgw_ip_port_len);
        sequence_number = eps_bearer_ctxt_p->sgw_sequence_number;
        sgw_s8_free_eps_bearer_context(&eps_bearer_ctxt_p);
        sgw_context_p->pdn_connection.sgw_eps_bearers_array[EBI_TO_INDEX(ebi)] =
            NULL;
        insert_sgw_cp_and_up_teid_to_directoryd(sgw_state, imsi64, UP_TEID);
        break;
      }
    }
  }
  send_s8_delete_bearer_response(
      s11_delete_bearer_response_p,
      sgw_context_p->pdn_connection.p_gw_teid_S5_S8_cp, sequence_number,
      pgw_cp_ip_port, sgw_context_p->imsi);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

void sgw_s8_send_failed_delete_bearer_response(
    const s8_delete_bearer_request_t* const db_req,
    gtpv2c_cause_value_t cause_value, Imsi_t imsi, teid_t pgw_s8_teid) {
  OAILOG_FUNC_IN(LOG_SGW_S8);
  itti_s11_nw_init_deactv_bearer_rsp_t s11_delete_bearer_rsp = {0};
  s11_delete_bearer_rsp.cause.cause_value = cause_value;
  s11_delete_bearer_rsp.bearer_contexts.num_bearer_context =
      db_req->num_eps_bearer_id;
  for (uint8_t idx = 0; idx < db_req->num_eps_bearer_id; idx++) {
    s11_delete_bearer_rsp.bearer_contexts.bearer_contexts[idx]
        .cause.cause_value = cause_value;
    s11_delete_bearer_rsp.bearer_contexts.bearer_contexts[idx].eps_bearer_id =
        db_req->eps_bearer_id[idx];
  }
  send_s8_delete_bearer_response(&s11_delete_bearer_rsp, pgw_s8_teid,
                                 db_req->sequence_number,
                                 db_req->pgw_cp_address, imsi);
  OAILOG_FUNC_OUT(LOG_SGW_S8);
}

sgw_eps_bearer_context_information_t* update_sgw_context_to_s11_teid_map(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64) {
  /* Once sgw_s8_teid is obtained from orc8r, move sgw_eps_bearer_context
   * from temporary_create_session_procedure_id_map to state_teid_map
   */
  sgw_eps_bearer_context_information_t* sgw_context_p = nullptr;
  sgw_state->temporary_create_session_procedure_id_map.get(
      session_rsp_p->temporary_create_session_procedure_id, &sgw_context_p);
  if (!sgw_context_p) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to fetch sgw_eps_bearer_context_info from "
                    "temporary_create_session_procedure_id:%u \n",
                    session_rsp_p->temporary_create_session_procedure_id);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_context_p);
  }

  if (sgw_state->temporary_create_session_procedure_id_map.remove(
          session_rsp_p->temporary_create_session_procedure_id, false) !=
      magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SGW_S8, imsi64,
                    "Failed to remove sgw_eps_bearer_context_info from map "
                    "with temporary procedure id"
                    "temporary_create_session_procedure_id:%u \n",
                    session_rsp_p->temporary_create_session_procedure_id);
    OAILOG_FUNC_RETURN(LOG_SGW_S8, nullptr);
  }

  /* Teid shall remain same for both sgw's s11 interface and s8 interface as
   * teid is allocated per PDN
   */
  sgw_context_p->s_gw_teid_S11_S4 = session_rsp_p->context_teid;
  sgw_context_p->pdn_connection.s_gw_teid_S5_S8_cp =
      session_rsp_p->context_teid;
  // Insert the new tunnel with sgw_s11_teid into the state_teid_map
  map_uint32_sgw_eps_bearer_context_t* state_teid_map = get_s8_state_teid_map();
  if (!state_teid_map) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to get state_teid_map");
    OAILOG_FUNC_RETURN(LOG_SGW_S8, nullptr);
  }
  state_teid_map->insert(session_rsp_p->context_teid, sgw_context_p);

  OAILOG_DEBUG(LOG_SGW_S8,
               "Inserted new sgw eps bearer context into state_teid_map with "
               "key as sgw_s11_teid " TEID_FMT "\n ",
               session_rsp_p->context_teid);
  OAILOG_FUNC_RETURN(LOG_SGW_S8, sgw_context_p);
}

//------------------------------------------------------------------------------
static void sgw_s8_release_all_enb_related_information(
    sgw_eps_bearer_ctxt_t* const eps_bearer_ctxt) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  if (eps_bearer_ctxt) {
    memset(&eps_bearer_ctxt->enb_ip_address_S1u, 0,
           sizeof(eps_bearer_ctxt->enb_ip_address_S1u));
    eps_bearer_ctxt->enb_teid_S1u = INVALID_TEID;
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

void sgw_s8_process_release_access_bearer_request(
    log_proto_t module, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context) {
  OAILOG_FUNC_IN(module);
  int rv = RETURNok;
  /*
   * Release the tunnels so that in idle state, DL packets are not sent
   * towards eNB.
   * These tunnels will be added again when UE moves back to connected mode.
   */
  for (uint8_t ebx = 0; ebx < BEARERS_PER_UE; ebx++) {
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt =
        sgw_context->pdn_connection.sgw_eps_bearers_array[ebx];
    if (eps_bearer_ctxt) {
      struct in_addr enb = {.s_addr = 0};
      struct in6_addr* enb_ipv6 = nullptr;
      enb.s_addr =
          eps_bearer_ctxt->enb_ip_address_S1u.address.ipv4_address.s_addr;
      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
          eps_bearer_ctxt->enb_ip_address_S1u.pdn_type == IPv6) {
        enb_ipv6 = &eps_bearer_ctxt->enb_ip_address_S1u.address.ipv6_address;
      }

      struct in6_addr* ue_ipv6 = nullptr;
      if ((eps_bearer_ctxt->paa.pdn_type == IPv6) ||
          (eps_bearer_ctxt->paa.pdn_type == IPv4_AND_v6)) {
        ue_ipv6 = &eps_bearer_ctxt->paa.ipv6_address;
      }
      struct in_addr ue_ipv4 = {.s_addr = 0};
      ue_ipv4.s_addr = eps_bearer_ctxt->paa.ipv4_address.s_addr;
      struct in_addr pgw = {.s_addr = 0};
      pgw.s_addr =
          eps_bearer_ctxt->p_gw_address_in_use_up.address.ipv4_address.s_addr;
      struct in6_addr* pgw_ipv6 = nullptr;
      if ((eps_bearer_ctxt->p_gw_address_in_use_up.pdn_type == IPv6) ||
          (eps_bearer_ctxt->p_gw_address_in_use_up.pdn_type == IPv4_AND_v6)) {
        pgw_ipv6 =
            &eps_bearer_ctxt->p_gw_address_in_use_up.address.ipv6_address;
      }
      OAILOG_DEBUG_UE(
          module, imsi64,
          "Deleting tunnel for bearer_id %u ue addr %x enb_ip %x "
          "sgw_teid_s1u_s12_s4_up " TEID_FMT "enb_teid_S1u " TEID_FMT
          " pgw_up_ip %x "
          "pgw_up_teid " TEID_FMT
          "s_gw_ip_address_S5_S8_up %x"
          "s_gw_teid_S5_S8_up " TEID_FMT,
          eps_bearer_ctxt->eps_bearer_id, ue_ipv4.s_addr, enb.s_addr,
          eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up,
          eps_bearer_ctxt->enb_teid_S1u, pgw.s_addr,
          eps_bearer_ctxt->p_gw_teid_S5_S8_up,
          eps_bearer_ctxt->s_gw_ip_address_S5_S8_up.address.ipv4_address.s_addr,
          eps_bearer_ctxt->s_gw_teid_S5_S8_up);
#if !MME_UNIT_TEST  // skip tunnel deletion for unit tests
      if (module == LOG_SPGW_APP) {
        rv = gtp_tunnel_ops->del_tunnel(
            enb, enb_ipv6, eps_bearer_ctxt->paa.ipv4_address, ue_ipv6,
            eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up,
            eps_bearer_ctxt->enb_teid_S1u, nullptr);
      } else if (module == LOG_SGW_S8) {
        rv = gtpv1u_del_s8_tunnel(enb, enb_ipv6, pgw, pgw_ipv6,
                                  eps_bearer_ctxt->paa.ipv4_address, ue_ipv6,
                                  eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up,
                                  eps_bearer_ctxt->s_gw_teid_S5_S8_up);
      }

      // TODO Need to add handling on failing to delete s1-u tunnel rules from
      // ovs flow table
      if (rv < 0) {
        OAILOG_ERROR_UE(module, imsi64,
                        "ERROR in deleting TUNNEL " TEID_FMT
                        " (eNB) <-> (SGW) " TEID_FMT "\n",
                        eps_bearer_ctxt->enb_teid_S1u,
                        eps_bearer_ctxt->s_gw_teid_S1u_S12_S4_up);
      }
      // Paging is performed without packet buffering
      rv = gtp_tunnel_ops->add_paging_rule(
          sgw_context->imsi, eps_bearer_ctxt->paa.ipv4_address, ue_ipv6);
      // Convert to string for logging
      char* ip_str = inet_ntoa(eps_bearer_ctxt->paa.ipv4_address);
      if (rv < 0) {
        OAILOG_ERROR_UE(module, imsi64,
                        "ERROR in setting paging rule for IP Addr: %s\n",
                        ip_str);
      }
      if ((eps_bearer_ctxt->paa.pdn_type == IPv6) ||
          (eps_bearer_ctxt->paa.pdn_type == IPv4_AND_v6)) {
        char ip6_str[INET6_ADDRSTRLEN];
        inet_ntop(AF_INET6,
                  reinterpret_cast<void*>(&eps_bearer_ctxt->paa.ipv6_address),
                  ip6_str, INET6_ADDRSTRLEN);
        OAILOG_DEBUG(module, "Set the paging rule for IPv6 Addr: %s\n",
                     ip6_str);
      } else {
        OAILOG_DEBUG(module, "Set the paging rule for IP Addr: %s\n", ip_str);
      }
#endif
      sgw_s8_release_all_enb_related_information(eps_bearer_ctxt);
    }
  }
  OAILOG_FUNC_OUT(module);
}

// Create temporary dedicated bearer context
status_code_e create_temporary_s8_dedicated_bearer_context(
    sgw_eps_bearer_context_information_t* sgw_ctxt_p,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    uint32_t sequence_number, log_proto_t module) {
  OAILOG_FUNC_IN(module);
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = new sgw_eps_bearer_ctxt_t();

  // Copy PAA from default bearer cntxt
  sgw_eps_bearer_ctxt_t* default_eps_bearer_entry_p =
      sgw_s8_cm_get_eps_bearer_entry(&sgw_ctxt_p->pdn_connection,
                                     sgw_ctxt_p->pdn_connection.default_bearer);

  if (!default_eps_bearer_entry_p) {
    OAILOG_ERROR_UE(module, sgw_ctxt_p->imsi64,
                    "Failed to get default bearer context for bearer id :%u",
                    sgw_ctxt_p->pdn_connection.default_bearer);
    OAILOG_FUNC_RETURN(module, RETURNerror);
  }

  eps_bearer_ctxt_p->eps_bearer_id = 0;
  eps_bearer_ctxt_p->paa = default_eps_bearer_entry_p->paa;
  // SGW FTEID
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = s1_u_sgw_fteid;

  if (pdn_type == IPv4 || pdn_type == IPv4_AND_v6) {
    eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.pdn_type = IPv4;
    eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.address.ipv4_address
        .s_addr = sgw_ip_address_S1u_S12_S4_up;
  } else {
    eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.pdn_type = IPv6;
    memcpy(
        &eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.address.ipv6_address,
        &sgw_ipv6_address_S1u_S12_S4_up,
        sizeof(eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.address
                   .ipv6_address));
  }
  // DL TFT
  memcpy(&eps_bearer_ctxt_p->tft, &bearer_req_p->dl_tft,
         sizeof(traffic_flow_template_t));
  // QoS
  memcpy(&eps_bearer_ctxt_p->eps_bearer_qos, &bearer_req_p->eps_bearer_qos,
         sizeof(bearer_qos_t));
  // Save Policy Rule Name
  snprintf(eps_bearer_ctxt_p->policy_rule_name, POLICY_RULE_NAME_MAXLEN + 1,
           "%s", bearer_req_p->policy_rule_name);
  eps_bearer_ctxt_p->sgw_sequence_number = sequence_number;
  OAILOG_INFO_UE(module, sgw_ctxt_p->imsi64,
                 "Number of DL packet filter rules: %d\n",
                 eps_bearer_ctxt_p->tft.numberofpacketfilters);

  // Create temporary spgw bearer context entry
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_ctxt_p);
  if (!pgw_ni_cbr_proc) {
    OAILOG_DEBUG_UE(module, sgw_ctxt_p->imsi64,
                    "Creating a new temporary eps bearer context entry\n");
    pgw_ni_cbr_proc = pgw_create_procedure_create_bearer(sgw_ctxt_p);
    if (!pgw_ni_cbr_proc) {
      OAILOG_ERROR_UE(module, sgw_ctxt_p->imsi64,
                      "Failed to create temporary eps bearer context entry\n");
      OAILOG_FUNC_RETURN(module, RETURNerror);
    }
  }
  struct sgw_eps_bearer_entry_wrapper_s* sgw_eps_bearer_entry_p =
      new sgw_eps_bearer_entry_wrapper_s();
  sgw_eps_bearer_entry_p->sgw_eps_bearer_entry = eps_bearer_ctxt_p;
  LIST_INSERT_HEAD((pgw_ni_cbr_proc->pending_eps_bearers),
                   sgw_eps_bearer_entry_p, entries);
  OAILOG_FUNC_RETURN(module, RETURNok);
}

// Build and send ITTI message, s11_create_bearer_request to MME APP
status_code_e sgw_build_and_send_s11_create_bearer_request(
    sgw_eps_bearer_context_information_t* sgw_eps_bearer_context_information,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    log_proto_t module) {
  OAILOG_FUNC_IN(module);
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNerror;

  message_p = itti_alloc_new_message(
      (module == LOG_SPGW_APP ? TASK_SPGW_APP : TASK_SGW_S8),
      S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST);
  if (!message_p) {
    OAILOG_ERROR_UE(module, sgw_eps_bearer_context_information->imsi64,
                    "Failed to allocate message_p for"
                    "S11_NW_INITIATED_BEARER_ACTV_REQUEST\n");
    OAILOG_FUNC_RETURN(module, rc);
  }

  itti_s11_nw_init_actv_bearer_request_t* s11_actv_bearer_request =
      &message_p->ittiMsg.s11_nw_init_actv_bearer_request;
  memset(s11_actv_bearer_request, 0,
         sizeof(itti_s11_nw_init_actv_bearer_request_t));
  // Context TEID
  s11_actv_bearer_request->s11_mme_teid =
      sgw_eps_bearer_context_information->mme_teid_S11;
  // LBI
  s11_actv_bearer_request->lbi = bearer_req_p->lbi;
  // UL TFT to be sent to UE
  memcpy(&s11_actv_bearer_request->tft, &bearer_req_p->ul_tft,
         sizeof(traffic_flow_template_t));
  // QoS
  memcpy(&s11_actv_bearer_request->eps_bearer_qos,
         &bearer_req_p->eps_bearer_qos, sizeof(bearer_qos_t));
  // S1U SGW F-TEID
  s11_actv_bearer_request->s1_u_sgw_fteid.teid = s1_u_sgw_fteid;
  s11_actv_bearer_request->s1_u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
  // Set IPv4 address type bit

  if (pdn_type == IPv4 || pdn_type == IPv4_AND_v6) {
    s11_actv_bearer_request->s1_u_sgw_fteid.ipv4 = true;
    s11_actv_bearer_request->s1_u_sgw_fteid.ipv4_address.s_addr =
        sgw_ip_address_S1u_S12_S4_up;
  } else {
    s11_actv_bearer_request->s1_u_sgw_fteid.ipv6 = true;
    memcpy(&s11_actv_bearer_request->s1_u_sgw_fteid.ipv6_address,
           sgw_ipv6_address_S1u_S12_S4_up,
           sizeof(s11_actv_bearer_request->s1_u_sgw_fteid.ipv6_address));
  }
  message_p->ittiMsgHeader.imsi = sgw_eps_bearer_context_information->imsi64;
  OAILOG_INFO_UE(module, sgw_eps_bearer_context_information->imsi64,
                 "Sending S11 Create Bearer Request to MME_APP for LBI %d \n",
                 bearer_req_p->lbi);
  if (module == LOG_SPGW_APP) {
    rc = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  } else if (module == LOG_SGW_S8) {
    rc = send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);
  } else {
    OAILOG_ERROR_UE(module, sgw_eps_bearer_context_information->imsi64,
                    "Invalid module \n");
  }
  OAILOG_FUNC_RETURN(module, rc);
}

//-----------------------------------------------------------------------------
sgw_eps_bearer_ctxt_t* sgw_cm_insert_eps_bearer_ctxt_in_collection(
    sgw_pdn_connection_t* const sgw_pdn_connection,
    sgw_eps_bearer_ctxt_t* const sgw_eps_bearer_ctxt) {
  if (!sgw_eps_bearer_ctxt) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "Failed to insert EPS bearer context : nullptr context\n");
    return nullptr;
  }

  if (!sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
          sgw_eps_bearer_ctxt->eps_bearer_id)]) {
    sgw_pdn_connection->sgw_eps_bearers_array[EBI_TO_INDEX(
        sgw_eps_bearer_ctxt->eps_bearer_id)] = sgw_eps_bearer_ctxt;
    OAILOG_DEBUG(LOG_SPGW_APP,
                 "Inserted new EPS bearer entry for EPS bearer id %u \n",
                 sgw_eps_bearer_ctxt->eps_bearer_id);
  } else {
    OAILOG_WARNING(
        LOG_SPGW_APP,
        "Could not create new EPS bearer ctxt for EPS bearer id %u : already "
        "exist\n",
        sgw_eps_bearer_ctxt->eps_bearer_id);
  }
  return sgw_eps_bearer_ctxt;
}

bool is_sgw_enb_ip_address_same(const fteid_t* fte_p, ip_address_t* ip_p) {
  bool rc = true;

  switch ((ip_p)->pdn_type) {
    case IPv4:
      if ((ip_p)->address.ipv4_address.s_addr != (fte_p)->ipv4_address.s_addr) {
        rc = false;
      }
      break;
    case IPv4_AND_v6:
    case IPv6:
      if (memcmp(&(ip_p)->address.ipv6_address, &(fte_p)->ipv6_address,
                 sizeof((ip_p)->address.ipv6_address)) != 0) {
        rc = false;
      }
      break;
    default:
      rc = true;
      break;
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

void sgw_s8_populate_mbr_bearer_contexts_removed(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    sgw_eps_bearer_context_information_t* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint8_t rsp_idx = 0;
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p = nullptr;
  for (uint8_t idx = 0; idx < resp_pP->num_bearers_removed; idx++) {
    eps_bearer_ctxt_p = sgw_s8_cm_get_eps_bearer_entry(
        &(sgw_context_p->pdn_connection),
        resp_pP->bearer_contexts_to_be_removed[idx]);
    /* If context is found, delete the context and set cause as
     * REQUEST_ACCEPTED. If context is not found set the cause as
     * CONTEXT_NOT_FOUND. MME App sends bearer deactivation message to UE for
     * the bearers with cause CONTEXT_NOT_FOUND
     */
    if (eps_bearer_ctxt_p != nullptr) {
      sgw_s8_free_eps_bearer_context(
          &(sgw_context_p->pdn_connection.sgw_eps_bearers_array[EBI_TO_INDEX(
              eps_bearer_ctxt_p->eps_bearer_id)]));
      modify_response_p->bearer_contexts_marked_for_removal
          .bearer_contexts[rsp_idx]
          .cause.cause_value = REQUEST_ACCEPTED;
    } else {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, imsi64,
          "Rx SGI_UPDATE_ENDPOINT_RESPONSE: eps_bearer_ctxt_p not found for "
          "bearer to be removed ebi %u\n",
          resp_pP->bearer_contexts_to_be_removed[idx]);
      modify_response_p->bearer_contexts_marked_for_removal
          .bearer_contexts[rsp_idx]
          .cause.cause_value = CONTEXT_NOT_FOUND;
    }
    modify_response_p->bearer_contexts_marked_for_removal
        .bearer_contexts[rsp_idx++]
        .eps_bearer_id = resp_pP->bearer_contexts_to_be_removed[idx];
    modify_response_p->bearer_contexts_marked_for_removal.num_bearer_context++;
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

void populate_s8_sgi_end_point_update(
    uint8_t sgi_rsp_idx, uint8_t idx,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p,
    itti_sgi_update_end_point_response_t* sgi_update_end_point_resp) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  FTEID_T_2_IP_ADDRESS_T(
      (&modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
            .s1_eNB_fteid),
      (&eps_bearer_ctxt_p->enb_ip_address_S1u));
  eps_bearer_ctxt_p->enb_teid_S1u =
      modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
          .s1_eNB_fteid.teid;
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .sgw_S1u_teid = eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up;
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .enb_S1u_teid = eps_bearer_ctxt_p->enb_teid_S1u;
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .eps_bearer_id = eps_bearer_ctxt_p->eps_bearer_id;
  sgi_update_end_point_resp->num_bearers_modified++;

  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}
