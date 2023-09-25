/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software *
 *distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 *WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *License for the specific language governing permissions and limitations under
 *the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file pgw_handlers.cpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define S5_HANDLERS_C

#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/include/ip_forward_messages_types.h"
#include "lte/gateway/c/core/oai/include/pgw_config.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/service303.hpp"
#include "lte/gateway/c/core/oai/include/sgw_ie_defs.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.hpp"
#include "lte/gateway/c/core/oai/lib/pcef/pcef_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_pco.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

extern void print_bearer_ids_helper(const ebi_t*, uint32_t);
extern task_zmq_ctx_t sgw_s8_task_zmq_ctx;
extern spgw_config_t spgw_config;

static void delete_temporary_dedicated_bearer_context(
    teid_t s1_u_sgw_fteid, ebi_t lbi,
    magma::lte::oai::S11BearerContext* spgw_context_p);

static void spgw_handle_s5_response_with_error(
    spgw_state_t* spgw_state,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    teid_t context_teid, ebi_t eps_bearer_id,
    s5_create_session_response_t* s5_response);

//--------------------------------------------------------------------------------

void handle_s5_create_session_request(
    spgw_state_t* spgw_state,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    teid_t context_teid, ebi_t eps_bearer_id) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  s5_create_session_response_t s5_response = {0};

  std::string imsi;
  std::string apn;

  if (!new_bearer_ctxt_info_p) {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "Failed to fetch sgw bearer context from the received context "
                 "teid" TEID_FMT "\n",
                 context_teid);
    s5_response.status = SGI_STATUS_ERROR_CONTEXT_NOT_FOUND;
    spgw_handle_s5_response_with_error(spgw_state, new_bearer_ctxt_info_p,
                                       context_teid, eps_bearer_id,
                                       &s5_response);
  }
  magma::lte::oai::SgwEpsBearerContextInfo sgw_context =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context();
  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, sgw_context.imsi64(),
      "Handle s5_create_session_request, for context sgw s11 teid, " TEID_FMT
      "EPS bearer id %u\n",
      context_teid, eps_bearer_id);

  imsi = sgw_context.imsi();

  apn = sgw_context.pdn_connection().apn_in_use();

  switch (sgw_context.saved_message().pdn_type()) {
    case IPv4:
      pgw_handle_allocate_ipv4_address(imsi, apn, "ipv4", context_teid,
                                       eps_bearer_id);
      break;

    case IPv6:
      pgw_handle_allocate_ipv6_address(imsi, apn, "ipv6", context_teid,
                                       eps_bearer_id);
      break;

    case IPv4_AND_v6:
      pgw_handle_allocate_ipv4v6_address(imsi, apn, "ipv4v6", context_teid,
                                         eps_bearer_id);
      break;

    default:
      OAILOG_ERROR(LOG_SPGW_APP, "BAD paa.pdn_type %d",
                   sgw_context.saved_message().pdn_type());
      break;
  }
}

void spgw_handle_s5_response_with_error(
    spgw_state_t* spgw_state,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    teid_t context_teid, ebi_t eps_bearer_id,
    s5_create_session_response_t* s5_response) {
  s5_response->context_teid = context_teid;
  s5_response->eps_bearer_id = eps_bearer_id;
  s5_response->failure_cause = S5_OK;

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
      "Sending S5 Create Session Response to SGW: with context teid, " TEID_FMT
      "EPS Bearer Id = %u\n",
      s5_response->context_teid, s5_response->eps_bearer_id);
  handle_s5_create_session_response(spgw_state, new_bearer_ctxt_info_p,
                                    (*s5_response));
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

void spgw_handle_pcef_create_session_response(
    spgw_state_t* spgw_state,
    const itti_pcef_create_session_response_t* const pcef_csr_resp_p,
    imsi64_t imsi64) {
  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Received PCEF-CREATE-SESSION-RESPONSE");

  s5_create_session_response_t s5_response = {0};
  s5_response.context_teid = pcef_csr_resp_p->teid;
  s5_response.eps_bearer_id = pcef_csr_resp_p->eps_bearer_id;
  s5_response.status = pcef_csr_resp_p->sgi_status;
  s5_response.failure_cause = S5_OK;

  magma::lte::oai::S11BearerContext* bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(pcef_csr_resp_p->teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  if (sgw_cm_get_eps_bearer_entry(
          bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
              ->mutable_pdn_connection(),
          s5_response.eps_bearer_id, &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to find eps bearer entry for bearer id %u",
                    s5_response.eps_bearer_id);
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  std::string apn = bearer_ctxt_info_p->sgw_eps_bearer_context()
                        .pdn_connection()
                        .apn_in_use();
  eps_bearer_ctxt.set_update_teids(true);
  // Updating eps_bearer_ctxt
  if (sgw_update_eps_bearer_entry(
          bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
              ->mutable_pdn_connection(),
          s5_response.eps_bearer_id, &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SPGW_APP,
                    bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
                    "Failed to update bearer context for bearer id :%u\n",
                    s5_response.eps_bearer_id);
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  std::string imsi_str = bearer_ctxt_info_p->sgw_eps_bearer_context().imsi();
  if (pcef_csr_resp_p->rpc_status != PCEF_STATUS_OK) {
    struct in_addr ue_ipv4 = {.s_addr = 0};
    struct in6_addr ue_ipv6;
    convert_proto_ip_to_standard_ip_fmt(eps_bearer_ctxt.mutable_ue_ip_paa(),
                                        &ue_ipv4, &ue_ipv6, true);
    if ((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4) ||
        (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) {
      release_ipv4_address(imsi_str, apn, &ue_ipv4);
    }
    if ((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6) ||
        (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) {
      release_ipv6_address(imsi_str, apn, &ue_ipv6);
    }
    s5_response.failure_cause = PCEF_FAILURE;
  }
  handle_s5_create_session_response(spgw_state, bearer_ctxt_info_p,
                                    s5_response);
}

/*
 * Handle NW initiated Dedicated Bearer Activation from SPGW service
 */
status_code_e spgw_handle_nw_initiated_bearer_actv_req(
    spgw_state_t* spgw_state,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    imsi64_t imsi64, gtpv2c_cause_value_t* failed_cause) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  status_code_e rc = RETURNok;
  state_teid_map_t* state_teid_map = nullptr;
  magma::lte::oai::S11BearerContext* spgw_ctxt_p = nullptr;
  bool is_imsi_found = false;
  bool is_lbi_found = false;

  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Received Create Bearer Req from PCRF with lbi:%u ",
                 bearer_req_p->lbi);

  state_teid_map = get_spgw_teid_state();
  if (state_teid_map == nullptr) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "Failed to get state_teid_map");
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  /* On reception of Dedicated Bearer Activation Request from PCRF,
   * SPGW shall identify whether valid PDN session exists for the UE
   * using IMSI and LBI, for which Dedicated Bearer Activation is requested.
   */
  for (auto itr = state_teid_map->map->begin();
       itr != state_teid_map->map->end(); itr++) {
    if (!is_lbi_found) {
      state_teid_map->get(itr->first, &spgw_ctxt_p);
      if (spgw_ctxt_p != nullptr) {
        if (!strncmp((const char*)spgw_ctxt_p->sgw_eps_bearer_context()
                         .imsi()
                         .c_str(),
                     (const char*)bearer_req_p->imsi,
                     strlen((const char*)bearer_req_p->imsi))) {
          is_imsi_found = true;
          if (spgw_ctxt_p->sgw_eps_bearer_context()
                  .pdn_connection()
                  .default_bearer() == bearer_req_p->lbi) {
            is_lbi_found = true;
            break;
          }
        }
      }
    }
  }

  if ((!is_imsi_found) || (!is_lbi_found)) {
    OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                   "is_imsi_found (%d), is_lbi_found (%d)\n", is_imsi_found,
                   is_lbi_found);
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Sending dedicated_bearer_actv_rsp with REQUEST_REJECTED cause to "
        "NW\n");
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  teid_t s1_u_sgw_fteid = spgw_get_new_s1u_teid(spgw_state);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  if (sgw_cm_get_eps_bearer_entry(spgw_ctxt_p->mutable_sgw_eps_bearer_context()
                                      ->mutable_pdn_connection(),
                                  bearer_req_p->lbi,
                                  &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "Failed to retrieve bearer ctxt:%u\n",
                    bearer_req_p->lbi);
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  // Create temporary dedicated bearer context
  rc = create_temporary_dedicated_bearer_context(
      spgw_ctxt_p->mutable_sgw_eps_bearer_context(), bearer_req_p,
      eps_bearer_ctxt.sgw_s1u_s12_s4_up_ip_addr().pdn_type(),
      spgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
      &spgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid, 0,
      LOG_SPGW_APP);
  if (rc != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Failed to create temporary dedicated bearer context for lbi: %u \n ",
        bearer_req_p->lbi);
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  // Build and send ITTI message, s11_create_bearer_request to MME APP
  rc = sgw_build_and_send_s11_create_bearer_request(
      spgw_ctxt_p->mutable_sgw_eps_bearer_context(), bearer_req_p,
      eps_bearer_ctxt.sgw_s1u_s12_s4_up_ip_addr().pdn_type(),
      spgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr,
      &spgw_state->sgw_ipv6_address_S1u_S12_S4_up, s1_u_sgw_fteid,
      LOG_SPGW_APP);
  if (rc != RETURNok) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Failed to build and send S11 Create Bearer Request for lbi :%u \n",
        bearer_req_p->lbi);

    *failed_cause = REQUEST_REJECTED;
    delete_temporary_dedicated_bearer_context(s1_u_sgw_fteid, bearer_req_p->lbi,
                                              spgw_ctxt_p);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e spgw_handle_nw_initiated_bearer_deactv_req(
    const itti_gx_nw_init_deactv_bearer_request_t* const bearer_req_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  status_code_e rc = RETURNok;
  state_teid_map_t* state_teid_map = nullptr;
  magma::lte::oai::S11BearerContext* spgw_ctxt_p = nullptr;
  bool is_lbi_found = false;
  bool is_imsi_found = false;
  bool is_ebi_found = false;
  ebi_t ebi_to_be_deactivated[BEARERS_PER_UE] = {0};
  uint32_t no_of_bearers_to_be_deact = 0;
  uint32_t no_of_bearers_rej = 0;
  ebi_t invalid_bearer_id[BEARERS_PER_UE] = {0};

  OAILOG_INFO_UE(
      LOG_SPGW_APP, imsi64,
      "Received nw_initiated_deactv_bearer_req from SPGW service \n");
  print_bearer_ids_helper(bearer_req_p->ebi, bearer_req_p->no_of_bearers);

  state_teid_map = get_spgw_teid_state();
  if (state_teid_map == nullptr) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "No s11_bearer_context_information is found in state_teid_map\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  // Check if valid LBI and EBI recvd
  /* For multi PDN, same IMSI can have multiple sessions, which means there
   * will be multiple entries for different sessions with the same IMSI. Hence
   * even though IMSI is found search the entire list for the LBI
   */
  for (auto itr = state_teid_map->map->begin();
       itr != state_teid_map->map->end(); itr++) {
    if (!is_lbi_found) {
      state_teid_map->get(itr->first, &spgw_ctxt_p);
      if (spgw_ctxt_p != nullptr) {
        if (!strncmp((const char*)spgw_ctxt_p->sgw_eps_bearer_context()
                         .imsi()
                         .c_str(),
                     (const char*)bearer_req_p->imsi,
                     strlen((const char*)bearer_req_p->imsi))) {
          is_imsi_found = true;
          if ((bearer_req_p->lbi != 0) &&
              (bearer_req_p->lbi == spgw_ctxt_p->sgw_eps_bearer_context()
                                        .pdn_connection()
                                        .default_bearer())) {
            is_lbi_found = true;
            // Check if the received EBI is valid
            for (uint32_t itrn = 0; itrn < bearer_req_p->no_of_bearers;
                 itrn++) {
              magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
              if (sgw_cm_get_eps_bearer_entry(
                      spgw_ctxt_p->mutable_sgw_eps_bearer_context()
                          ->mutable_pdn_connection(),
                      bearer_req_p->ebi[itrn],
                      &eps_bearer_ctxt) == magma::PROTO_MAP_OK) {
                is_ebi_found = true;
                ebi_to_be_deactivated[no_of_bearers_to_be_deact] =
                    bearer_req_p->ebi[itrn];
                no_of_bearers_to_be_deact++;
              } else {
                invalid_bearer_id[no_of_bearers_rej] = bearer_req_p->ebi[itrn];
                no_of_bearers_rej++;
              }
            }
            break;
          }
        }
      }
    }
  }

  /* Send reject to NW if we did not find ebi/lbi/imsi.
   * Also in case of multiple bearers, if some EBIs are valid and some are
   * not, send reject to those for which we did not find EBI. Proceed with
   * deactivation by sending s5_nw_init_deactv_bearer_request to SGW for valid
   * EBIs
   */
  if ((!is_ebi_found) || (!is_lbi_found) || (!is_imsi_found) ||
      (no_of_bearers_rej > 0)) {
    OAILOG_INFO_UE(
        LOG_SPGW_APP, imsi64,
        "is_imsi_found (%d), is_lbi_found (%d), is_ebi_found (%d) \n",
        is_imsi_found, is_lbi_found, is_ebi_found);
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Sending dedicated bearer deactivation reject to NW\n");
    print_bearer_ids_helper(invalid_bearer_id, no_of_bearers_rej);
    // TODO-Uncomment once implemented at PCRF
    /* rc = send_dedicated_bearer_deactv_rsp(invalid_bearer_id,
         REQUEST_REJECTED);*/
  }

  if (no_of_bearers_to_be_deact > 0) {
    bool delete_default_bearer =
        (bearer_req_p->lbi == bearer_req_p->ebi[0]) ? true : false;
    rc = spgw_build_and_send_s11_deactivate_bearer_req(
        imsi64, no_of_bearers_to_be_deact, ebi_to_be_deactivated,
        delete_default_bearer,
        spgw_ctxt_p->sgw_eps_bearer_context().mme_teid_s11(), LOG_SPGW_APP);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

// Send ITTI message,S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST to mme_app
status_code_e spgw_build_and_send_s11_deactivate_bearer_req(
    imsi64_t imsi64, uint8_t no_of_bearers_to_be_deact,
    ebi_t* ebi_to_be_deactivated, bool delete_default_bearer,
    teid_t mme_teid_S11, log_proto_t module) {
  OAILOG_FUNC_IN(module);
  MessageDef* message_p = itti_alloc_new_message(
      (module == LOG_SPGW_APP ? TASK_SPGW_APP : TASK_SGW_S8),
      S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        module, imsi64,
        "itti_alloc_new_message failed for nw_initiated_deactv_bearer_req\n");
    OAILOG_FUNC_RETURN(module, RETURNerror);
  }
  itti_s11_nw_init_deactv_bearer_request_t* s11_bearer_deactv_request =
      &message_p->ittiMsg.s11_nw_init_deactv_bearer_request;
  memset(s11_bearer_deactv_request, 0,
         sizeof(itti_s11_nw_init_deactv_bearer_request_t));

  s11_bearer_deactv_request->s11_mme_teid = mme_teid_S11;
  /* If default bearer has to be deleted then the EBI list in the received
   * pgw_nw_init_deactv_bearer_request message contains a single entry at 0th
   * index and LBI == bearer_req_p->ebi[0]
   */
  s11_bearer_deactv_request->delete_default_bearer = delete_default_bearer;
  s11_bearer_deactv_request->no_of_bearers = no_of_bearers_to_be_deact;

  memcpy(s11_bearer_deactv_request->ebi, ebi_to_be_deactivated,
         (sizeof(ebi_t) * no_of_bearers_to_be_deact));
  print_bearer_ids_helper(s11_bearer_deactv_request->ebi,
                          s11_bearer_deactv_request->no_of_bearers);

  message_p->ittiMsgHeader.imsi = imsi64;
  OAILOG_INFO_UE(module, imsi64,
                 "Sending nw_initiated_deactv_bearer_req to mme_app "
                 "with delete_default_bearer flag set to %d\n",
                 s11_bearer_deactv_request->delete_default_bearer);
  status_code_e rc = RETURNerror;
  if (module == LOG_SPGW_APP) {
    rc = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  } else if (module == LOG_SGW_S8) {
    rc = send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);
  } else {
    OAILOG_ERROR_UE(module, imsi64, "Invalid module \n");
  }
  OAILOG_FUNC_RETURN(module, rc);
}

//------------------------------------------------------------------------------
status_code_e spgw_send_nw_init_activate_bearer_rsp(
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_ctx,
    uint8_t default_bearer_id, char* policy_rule_name) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  status_code_e rc = RETURNok;

  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Sending Create Bearer Rsp to PCRF with EBI %d with "
                 "cause: %d linked bearer id: %d policy rule name: %s\n",
                 bearer_ctx->eps_bearer_id, cause, default_bearer_id,
                 policy_rule_name);
  // Send Dedicated Bearer ID and Policy Rule ID binding to PCRF
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*)imsi_str, IMSI_BCD_DIGITS_MAX);
  if (cause == REQUEST_ACCEPTED) {
    pcef_send_policy2bearer_binding(imsi_str, default_bearer_id,
                                    policy_rule_name, bearer_ctx->eps_bearer_id,
                                    bearer_ctx->s1u_sgw_fteid.teid,
                                    bearer_ctx->s1u_enb_fteid.teid);
  } else {
    // Send 0 as dedicated bearer id if the create bearer request
    // was not accepted. Session manager should delete the policy rule
    // for this bearer. Set the tunnel IDs to zero as well.
    pcef_send_policy2bearer_binding(imsi_str, default_bearer_id,
                                    policy_rule_name, 0, 0, 0);
  }

  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

//------------------------------------------------------------------------------
status_code_e spgw_handle_nw_init_deactivate_bearer_rsp(gtpv2c_cause_t cause,
                                                        ebi_t lbi) {
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  OAILOG_INFO(
      LOG_SPGW_APP,
      "To be implemented: Sending Delete Bearer Rsp to PCRF with LBI %u with "
      "cause :%d\n",
      lbi, cause.cause_value);
  // Send Delete Bearer Rsp to PCRF
  // TODO-Uncomment once implemented at PCRF
  // rc = send_dedicated_bearer_deactv_rsp(lbi, cause);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

// Build and send ITTI message, s11_create_bearer_request to MME APP
status_code_e sgw_build_and_send_s11_create_bearer_request(
    magma::lte::oai::SgwEpsBearerContextInfo*
        sgw_eps_bearer_context_information,
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
    OAILOG_ERROR_UE(module, sgw_eps_bearer_context_information->imsi64(),
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
      sgw_eps_bearer_context_information->mme_teid_s11();
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
  message_p->ittiMsgHeader.imsi = sgw_eps_bearer_context_information->imsi64();
  OAILOG_INFO_UE(module, sgw_eps_bearer_context_information->imsi64(),
                 "Sending S11 Create Bearer Request to MME_APP for LBI %d \n",
                 bearer_req_p->lbi);
  if (module == LOG_SPGW_APP) {
    rc = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  } else if (module == LOG_SGW_S8) {
    rc = send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME_APP, message_p);
  } else {
    OAILOG_ERROR_UE(module, sgw_eps_bearer_context_information->imsi64(),
                    "Invalid module \n");
  }
  OAILOG_FUNC_RETURN(module, rc);
}

// Create temporary dedicated bearer context
status_code_e create_temporary_dedicated_bearer_context(
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_ctxt_p,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    uint32_t sequence_number, log_proto_t module) {
  OAILOG_FUNC_IN(module);
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;

  // Copy PAA from default bearer cntxt
  magma::lte::oai::SgwEpsBearerContext default_eps_bearer_entry;
  if (sgw_cm_get_eps_bearer_entry(sgw_ctxt_p->mutable_pdn_connection(),
                                  sgw_ctxt_p->pdn_connection().default_bearer(),
                                  &default_eps_bearer_entry) !=
      magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(module, sgw_ctxt_p->imsi64(),
                    "Failed to get default bearer context\n");
    OAILOG_FUNC_RETURN(module, RETURNerror);
  }

  eps_bearer_ctxt.set_eps_bearer_id(0);
  eps_bearer_ctxt.mutable_ue_ip_paa()->MergeFrom(
      default_eps_bearer_entry.ue_ip_paa());
  // SGW FTEID
  eps_bearer_ctxt.set_sgw_teid_s1u_s12_s4_up(s1_u_sgw_fteid);

  if (pdn_type == IPv4 || pdn_type == IPv4_AND_v6) {
    eps_bearer_ctxt.mutable_sgw_s1u_s12_s4_up_ip_addr()->set_pdn_type(IPv4);
    char ip4_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &sgw_ip_address_S1u_S12_S4_up, ip4_str, INET_ADDRSTRLEN);
    eps_bearer_ctxt.mutable_sgw_s1u_s12_s4_up_ip_addr()->set_ipv4_addr(ip4_str);
  } else {
    char ip6_str[INET6_ADDRSTRLEN];
    eps_bearer_ctxt.mutable_sgw_s1u_s12_s4_up_ip_addr()->set_pdn_type(IPv6);
    inet_ntop(AF_INET6, sgw_ipv6_address_S1u_S12_S4_up, ip6_str,
              INET6_ADDRSTRLEN);
    eps_bearer_ctxt.mutable_sgw_s1u_s12_s4_up_ip_addr()->set_ipv6_addr(ip6_str);
  }
  // DL TFT
  traffic_flow_template_to_proto(&bearer_req_p->dl_tft,
                                 eps_bearer_ctxt.mutable_tft());
  // QoS
  eps_bearer_qos_to_proto(&bearer_req_p->eps_bearer_qos,
                          eps_bearer_ctxt.mutable_eps_bearer_qos());
  // Save Policy Rule Name
  eps_bearer_ctxt.set_policy_rule_name(
      bearer_req_p->policy_rule_name,
      strlen(bearer_req_p->policy_rule_name) + 1);
  eps_bearer_ctxt.set_sgw_sequence_number(sequence_number);
  OAILOG_INFO_UE(module, sgw_ctxt_p->imsi64(),
                 "Number of DL packet filter rules: %d\n",
                 eps_bearer_ctxt.tft().number_of_packet_filters());

  // Create temporary bearer context for NW initiated dedicated bearer request
  magma::lte::oai::PgwCbrProcedure* pgw_ni_cbr_proc =
      sgw_ctxt_p->add_pending_procedures();
  pgw_ni_cbr_proc->set_type(
      PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST);
  magma::lte::oai::SgwEpsBearerContext* eps_bearer_proto =
      pgw_ni_cbr_proc->add_pending_eps_bearers();
  eps_bearer_proto->MergeFrom(eps_bearer_ctxt);

  OAILOG_FUNC_RETURN(module, RETURNok);
}

// Deletes temporary dedicated bearer context
static void delete_temporary_dedicated_bearer_context(
    teid_t s1_u_sgw_fteid, ebi_t lbi,
    magma::lte::oai::S11BearerContext* spgw_context_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  magma::lte::oai::PgwCbrProcedure* pgw_ni_cbr_proc = nullptr;
  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      spgw_context_p->mutable_sgw_eps_bearer_context();

  if (!(sgw_context_p->pending_procedures_size())) {
    OAILOG_ERROR_UE(LOG_SPGW_APP,
                    spgw_context_p->sgw_eps_bearer_context().imsi64(),
                    "No temporary bearer contexts stored for lbi :%u \n", lbi);
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }

  OAILOG_INFO_UE(LOG_SPGW_APP,
                 spgw_context_p->sgw_eps_bearer_context().imsi64(),
                 "Delete temporary bearer context for lbi :%u \n", lbi);

  uint8_t num_of_bearers_deleted = 0;
  uint8_t num_of_pending_procedures = sgw_context_p->pending_procedures_size();
  for (uint8_t proc_index = 0;
       proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
    pgw_ni_cbr_proc = sgw_context_p->mutable_pending_procedures(proc_index);
    if (!pgw_ni_cbr_proc) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, spgw_context_p->sgw_eps_bearer_context().imsi64(),
          "Pending procedure within sgw_context is null for "
          "proc_index:%u and s1u_teid " TEID_FMT,
          proc_index,
          pgw_ni_cbr_proc->pending_eps_bearers(0).sgw_teid_s1u_s12_s4_up());
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }
    if (pgw_ni_cbr_proc->type() ==
        PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST) {
      num_of_bearers_deleted = pgw_ni_cbr_proc->pending_eps_bearers_size();
      for (uint8_t bearer_index = 0;
           bearer_index < pgw_ni_cbr_proc->pending_eps_bearers_size();
           bearer_index++) {
        magma::lte::oai::SgwEpsBearerContext bearer_context =
            pgw_ni_cbr_proc->pending_eps_bearers(bearer_index);
        if (bearer_context.sgw_teid_s1u_s12_s4_up() == s1_u_sgw_fteid) {
          --num_of_bearers_deleted;
        }
      }  // end of bearer index loop
      if (num_of_bearers_deleted == 0) {
        pgw_ni_cbr_proc->clear_pending_eps_bearers();
        --num_of_pending_procedures;
      }
    }
  }  // end of procedure index loop
  if (num_of_pending_procedures == 0) {
    sgw_context_p->clear_pending_procedures();
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}
