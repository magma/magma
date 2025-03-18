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

/*! \file sgw_handlers.cpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#define SGW
#define S11_HANDLERS_C

#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <netinet/in.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/pgw_config.h"
#include "lte/gateway/c/core/oai/include/sgw_config.h"
#include "lte/gateway/c/core/oai/include/sgw_ie_defs.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/queue.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/service303.hpp"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.hpp"
#include "lte/gateway/c/core/oai/lib/pcef/pcef_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_pco.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_ue_ip_address_alloc.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

extern task_zmq_ctx_t sgw_s8_task_zmq_ctx;
extern spgw_config_t spgw_config;
extern struct gtp_tunnel_ops* gtp_tunnel_ops;
extern void print_bearer_ids_helper(const ebi_t*, uint32_t);
static void add_tunnel_helper(
    magma::lte::oai::S11BearerContext* spgw_context,
    magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt_entry_p,
    imsi64_t imsi64);
static teid_t sgw_generate_new_s11_cp_teid(void);
static void get_session_req_data(
    spgw_state_t* spgw_state,
    const magma::lte::oai::CreateSessionMessage* saved_req,
    struct pcef_create_session_data* data);

constexpr char SPGW_EPS_BEARER_MAP_NAME[] = "spgw_eps_bearer_map";
#define TASK_MME TASK_MME_APP

//------------------------------------------------------------------------------
uint32_t spgw_get_new_s1u_teid(spgw_state_t* state) {
  __sync_fetch_and_add(&state->gtpv1u_teid, 1);
  return (state->gtpv1u_teid) % INITIAL_SGW_S8_S1U_TEID;
}

//------------------------------------------------------------------------------
status_code_e sgw_handle_s11_create_session_request(
    spgw_state_t* state,
    const itti_s11_create_session_request_t* const session_req_pP,
    imsi64_t imsi64) {
  mme_sgw_tunnel_t* new_endpoint_p = nullptr;
  magma::lte::oai::S11BearerContext* s_plus_p_gw_eps_bearer_ctxt_info_p =
      nullptr;
  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p = nullptr;

  OAILOG_FUNC_IN(LOG_SPGW_APP);
  increment_counter("spgw_create_session", 1, NO_LABELS);
  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Received S11 CREATE SESSION REQUEST from MME_APP\n");
  /*
   * Upon reception of create session request from MME,
   * * * * S-GW should create UE, eNB and MME contexts and forward message to
   * P-GW.
   */
  if (session_req_pP->rat_type != RAT_EUTRAN) {
    OAILOG_WARNING_UE(
        LOG_SPGW_APP, imsi64,
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
    /*
     * MME sent request with teid = 0. This is not valid...
     */
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "F-TEID parameter mismatch\n");
    increment_counter("spgw_create_session", 1, 2, "result", "failure", "cause",
                      "sender_fteid_incorrect_parameters");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  new_endpoint_p = sgw_cm_create_s11_tunnel(
      session_req_pP->sender_fteid_for_cp.teid, sgw_generate_new_s11_cp_teid());

  if (new_endpoint_p == nullptr) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Could not create new tunnel endpoint between S-GW and MME "
                    "for S11 abstraction\n");
    increment_counter("spgw_create_session", 1, 2, "result", "failure", "cause",
                      "s11_tunnel_creation_failure");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, imsi64,
      "Rx CREATE-SESSION-REQUEST MME S11 teid " TEID_FMT
      "S-GW S11 teid " TEID_FMT " APN %s EPS bearer Id %d\n",
      new_endpoint_p->remote_teid, new_endpoint_p->local_teid,
      session_req_pP->apn,
      session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
          .eps_bearer_id);

  if (spgw_update_teid_in_ue_context(imsi64, new_endpoint_p->local_teid) ==
      RETURNerror) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to update sgw_s11_teid" TEID_FMT
                    " in UE context \n",
                    new_endpoint_p->local_teid);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  s_plus_p_gw_eps_bearer_ctxt_info_p =
      sgw_cm_create_bearer_context_information_in_collection(
          new_endpoint_p->local_teid);
  if (s_plus_p_gw_eps_bearer_ctxt_info_p) {
    /*
     * We try to create endpoint for S11 interface. A NULL endpoint means that
     * either the teid is already in list of known teid or ENOMEM error has been
     * raised during malloc.
     */
    //--------------------------------------------------
    // copy informations from create session request to bearer context
    // information
    //--------------------------------------------------
    sgw_context_p =
        s_plus_p_gw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context();
    sgw_context_p->set_imsi(session_req_pP->imsi.digit, IMSI_BCD_DIGITS_MAX);
    s_plus_p_gw_eps_bearer_ctxt_info_p->mutable_pgw_eps_bearer_context()
        ->set_imsi(session_req_pP->imsi.digit, IMSI_BCD_DIGITS_MAX);
    sgw_context_p->set_imsi64(imsi64);
    sgw_context_p->set_imsi_unauth_indicator(1);
    s_plus_p_gw_eps_bearer_ctxt_info_p->mutable_pgw_eps_bearer_context()
        ->set_imsi_unauth_indicator(1);
    sgw_context_p->set_mme_teid_s11(session_req_pP->sender_fteid_for_cp.teid);
    sgw_context_p->set_sgw_teid_s11_s4(new_endpoint_p->local_teid);
    if (session_req_pP->trxn) {
      sgw_context_p->set_trxn(reinterpret_cast<char*>(session_req_pP->trxn));
    }

    FTEID_T_2_PROTO_IP((&session_req_pP->sender_fteid_for_cp),
                       sgw_context_p->mutable_mme_cp_ip_address_s11());
    sgw_context_p->clear_pdn_connection();
    magma::lte::oai::SgwPdnConnection* proto_pdn =
        sgw_context_p->mutable_pdn_connection();

    if (session_req_pP->apn) {
      proto_pdn->set_apn_in_use(session_req_pP->apn,
                                strlen(session_req_pP->apn));
    } else {
      proto_pdn->set_apn_in_use("NO APN");
    }
    proto_pdn->set_default_bearer(
        session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
            .eps_bearer_id);

    //--------------------------------------
    // EPS bearer entry
    //--------------------------------------
    // TODO several bearers
    map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
    eps_bearer_map.map = proto_pdn->mutable_eps_bearer_map();
    eps_bearer_map.set_name(SPGW_EPS_BEARER_MAP_NAME);
    magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
    if (eps_bearer_map.get(
            session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
                .eps_bearer_id,
            &eps_bearer_ctxt) == magma::PROTO_MAP_OK) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, imsi64,
          "Bearer contexts already exists for eps bearer id:%u ",
          session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
              .eps_bearer_id);
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
    }
    eps_bearer_ctxt.set_eps_bearer_id(
        session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
            .eps_bearer_id);
    eps_bearer_qos_to_proto(
        &session_req_pP->bearer_contexts_to_be_created.bearer_contexts[0]
             .bearer_level_qos,
        eps_bearer_ctxt.mutable_eps_bearer_qos());
    /*
     * Trying to insert the new tunnel into the tree.
     * If collision_p is not NULL (0), it means tunnel is already present.
     */
    sgw_create_session_message_to_proto(
        session_req_pP,
        s_plus_p_gw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
            ->mutable_saved_message());

    /*
     * Send a create bearer request to PGW and handle respond
     * asynchronously through sgw_handle_s5_create_bearer_response()
     */
    eps_bearer_ctxt.set_sgw_teid_s1u_s12_s4_up(spgw_get_new_s1u_teid(state));
    OAILOG_DEBUG_UE(
        LOG_SPGW_APP, imsi64,
        "Updated eps_bearer_entry_p eps_b_id %u with SGW S1U teid" TEID_FMT
        "\n",
        eps_bearer_ctxt.eps_bearer_id(), new_endpoint_p->local_teid);

    if (eps_bearer_map.insert(eps_bearer_ctxt.eps_bearer_id(),
                              eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Failed to insert eps bearer context for bearer id:%u",
                      eps_bearer_ctxt.eps_bearer_id());
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
    }
    sgw_display_s11_bearer_context_information(
        LOG_SPGW_APP, s_plus_p_gw_eps_bearer_ctxt_info_p);

    handle_s5_create_session_request(state, s_plus_p_gw_eps_bearer_ctxt_info_p,
                                     new_endpoint_p->local_teid,
                                     eps_bearer_ctxt.eps_bearer_id());
  } else {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Could not create new transaction for SESSION_CREATE message\n");
    free_cpp_wrapper(reinterpret_cast<void**>(&new_endpoint_p));
    increment_counter("spgw_create_session", 1, 2, "result", "failure", "cause",
                      "internal_software_error");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  free_cpp_wrapper(reinterpret_cast<void**>(&new_endpoint_p));
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e sgw_handle_sgi_endpoint_created(
    spgw_state_t* state, itti_sgi_create_end_point_response_t* const resp_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  itti_s11_create_session_response_t* create_session_response_p = NULL;
  MessageDef* message_p = NULL;
  status_code_e rv = RETURNok;

  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Rx SGI_CREATE_ENDPOINT_RESPONSE, Context: S11 teid " TEID_FMT
                  "EPS bearer id %u\n",
                  resp_pP->context_teid, resp_pP->eps_bearer_id);

  message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_CREATE_SESSION_RESPONSE);

  if (message_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  create_session_response_p = &message_p->ittiMsg.s11_create_session_response;
  memset(create_session_response_p, 0,
         sizeof(itti_s11_create_session_response_t));

  magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(resp_pP->context_teid);
  if (new_bearer_ctxt_info_p) {
    create_session_response_p->teid =
        new_bearer_ctxt_info_p->sgw_eps_bearer_context().mme_teid_s11();

    /*
     * Preparing to send create session response on S11 abstraction interface.
     * * * *  we set the cause value regarding the S1-U bearer establishment
     * result status.
     */
    if (resp_pP->status == SGI_STATUS_OK) {
      create_session_response_p->ambr.br_dl = 100000000;
      create_session_response_p->ambr.br_ul = 40000000;

      magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
      map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
      eps_bearer_map.map =
          new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
              ->mutable_pdn_connection()
              ->mutable_eps_bearer_map();
      if (eps_bearer_map.get(resp_pP->eps_bearer_id, &eps_bearer_ctxt) !=
          magma::PROTO_MAP_OK) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "Failed to find eps bearer entry for bearer id:%u ",
                        resp_pP->eps_bearer_id);
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
      }
      create_session_response_p->paa.pdn_type =
          (pdn_type_value_t)eps_bearer_ctxt.ue_ip_paa().pdn_type();
      if (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4 ||
          eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6) {
        inet_pton(AF_INET, eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str(),
                  &create_session_response_p->paa.ipv4_address.s_addr);
      }
      if (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6 ||
          eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6) {
        inet_pton(AF_INET6, eps_bearer_ctxt.ue_ip_paa().ipv6_addr().c_str(),
                  &create_session_response_p->paa.ipv6_address);
        create_session_response_p->paa.ipv6_prefix_length =
            eps_bearer_ctxt.ue_ip_paa().ipv6_prefix_length();
        create_session_response_p->paa.vlan =
            eps_bearer_ctxt.ue_ip_paa().vlan();
      }

      copy_protocol_configuration_options(&create_session_response_p->pco,
                                          &resp_pP->pco);
      clear_protocol_configuration_options(&resp_pP->pco);
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.teid = eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up();
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;

      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.ipv4 = 1;
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.ipv4_address.s_addr =
          state->sgw_ip_address_S1u_S12_S4_up.s_addr;
      if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled) {
        create_session_response_p->bearer_contexts_created.bearer_contexts[0]
            .s1u_sgw_fteid.ipv6 = 1;
        memcpy(&create_session_response_p->bearer_contexts_created
                    .bearer_contexts[0]
                    .s1u_sgw_fteid.ipv6_address,
               &state->sgw_ipv6_address_S1u_S12_S4_up,
               sizeof(create_session_response_p->bearer_contexts_created
                          .bearer_contexts[0]
                          .s1u_sgw_fteid.ipv6_address));
      }
      /*
       * Set the Cause information from bearer context created.
       * "Request accepted" is returned when the GTPv2 entity has accepted a
       * control plane request.
       */
      create_session_response_p->cause.cause_value = REQUEST_ACCEPTED;
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .cause.cause_value = REQUEST_ACCEPTED;
    } else {
      create_session_response_p->cause.cause_value = M_PDN_APN_NOT_ALLOWED;
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .cause.cause_value = M_PDN_APN_NOT_ALLOWED;
    }

    create_session_response_p->s11_sgw_fteid.teid = resp_pP->context_teid;
    create_session_response_p->s11_sgw_fteid.interface_type = S11_SGW_GTP_C;
    create_session_response_p->s11_sgw_fteid.ipv4 = 1;
    create_session_response_p->s11_sgw_fteid.ipv4_address.s_addr =
        spgw_config.sgw_config.ipv4.S11.s_addr;

    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .eps_bearer_id = resp_pP->eps_bearer_id;
    create_session_response_p->bearer_contexts_created.num_bearer_context += 1;

    memcpy(create_session_response_p->trxn,
           reinterpret_cast<const void*>(
               new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().c_str()),
           new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().size());

    inet_pton(AF_INET,
              new_bearer_ctxt_info_p->sgw_eps_bearer_context()
                  .mme_cp_ip_address_s11()
                  .ipv4_addr()
                  .c_str(),
              &create_session_response_p->peer_ip.s_addr);
  } else {
    create_session_response_p->cause.cause_value = CONTEXT_NOT_FOUND;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .cause.cause_value = CONTEXT_NOT_FOUND;
    create_session_response_p->bearer_contexts_created.num_bearer_context += 1;
  }

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, imsi64,
      "Tx CREATE-SESSION-RESPONSE SPGW -> TASK_MME, S11 MME teid " TEID_FMT
      " S11 S-GW teid " TEID_FMT " S1U teid " TEID_FMT
      " S1U addr 0x%x EPS bearer id %u status %d\n",
      create_session_response_p->teid,
      create_session_response_p->s11_sgw_fteid.teid,
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.teid,
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .s1u_sgw_fteid.ipv4_address.s_addr,
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .eps_bearer_id,
      create_session_response_p->bearer_contexts_created.bearer_contexts[0]
          .cause.cause_value);

  message_p->ittiMsgHeader.imsi = imsi64;
  rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);
}

//------------------------------------------------------------------------------
/* Populates bearer contexts marked for removal structure in
 * modify bearer rsp message.
 */
void sgw_populate_mbr_bearer_contexts_not_found(
    log_proto_t module,
    const itti_sgi_update_end_point_response_t* const resp_pP,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(module);
  uint8_t rsp_idx = 0;
  for (uint8_t idx = 0; idx < resp_pP->num_bearers_not_found; idx++) {
    modify_response_p->bearer_contexts_marked_for_removal
        .bearer_contexts[rsp_idx]
        .eps_bearer_id = resp_pP->bearer_contexts_not_found[idx];
    modify_response_p->bearer_contexts_marked_for_removal
        .bearer_contexts[rsp_idx++]
        .cause.cause_value = CONTEXT_NOT_FOUND;
    modify_response_p->bearer_contexts_marked_for_removal.num_bearer_context++;
  }
  OAILOG_FUNC_OUT(module);
}
//------------------------------------------------------------------------------
/* Populates bearer contexts marked for removal structure in
 * modify bearer rsp message
 */
void sgw_populate_mbr_bearer_contexts_removed(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint8_t rsp_idx = 0;
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  for (uint8_t idx = 0; idx < resp_pP->num_bearers_removed; idx++) {
    if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                    resp_pP->bearer_contexts_to_be_removed[idx],
                                    &eps_bearer_ctxt) == magma::PROTO_MAP_OK) {
      /* If context is found, delete the context and set cause as
       * REQUEST_ACCEPTED. If context is not found set the cause as
       * CONTEXT_NOT_FOUND. MME App sends bearer deactivation message to UE for
       * the bearers with cause CONTEXT_NOT_FOUND
       */
      sgw_remove_eps_bearer_context(sgw_context_p->mutable_pdn_connection(),
                                    eps_bearer_ctxt.eps_bearer_id());
      modify_response_p->bearer_contexts_marked_for_removal
          .bearer_contexts[rsp_idx]
          .cause.cause_value = REQUEST_ACCEPTED;
    } else {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Rx SGI_UPDATE_ENDPOINT_RESPONSE: eps_bearer_ctxt_p "
                      "not found for "
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

//------------------------------------------------------------------------------
/* Helper function to add gtp tunnels for default and
 * dedicated bearers
 */
static void sgw_add_gtp_tunnel(
    imsi64_t imsi64, magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt_p,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p) {
  int rv = RETURNok;
  struct in_addr enb = {.s_addr = 0};
  struct in6_addr enb_ipv6 = {};
  convert_proto_ip_to_standard_ip_fmt(
      eps_bearer_ctxt_p->mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
      spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

  struct in_addr ue_ipv4 = {.s_addr = 0};
  struct in6_addr ue_ipv6 = {};
  convert_proto_ip_to_standard_ip_fmt(eps_bearer_ctxt_p->mutable_ue_ip_paa(),
                                      &ue_ipv4, &ue_ipv6, true);
  bool is_ue_ipv6_pres = false;
  if (((eps_bearer_ctxt_p->ue_ip_paa().pdn_type() == IPv6 ||
        eps_bearer_ctxt_p->ue_ip_paa().pdn_type() == IPv4_AND_v6)) &&
      !(eps_bearer_ctxt_p->ue_ip_paa().ipv6_addr().empty())) {
    is_ue_ipv6_pres = true;
  }

  int vlan = eps_bearer_ctxt_p->ue_ip_paa().vlan();
  Imsi_t imsi;
  memcpy(imsi.digit,
         new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi().c_str(),
         new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi().size());
  const char* apn = reinterpret_cast<const char*>(
      new_bearer_ctxt_info_p->sgw_eps_bearer_context()
          .pdn_connection()
          .apn_in_use()
          .c_str());
  char ip6_str[INET6_ADDRSTRLEN];

  inet_ntop(AF_INET6, &ue_ipv6, ip6_str, INET6_ADDRSTRLEN);
  /* UE is switching back to EPS services after the CS Fallback
   * If Modify bearer Request is received in UE suspended mode, Resume PS
   * data
   */
  if (new_bearer_ctxt_info_p->sgw_eps_bearer_context()
          .pdn_connection()
          .ue_suspended_for_ps_handover()) {
    rv = gtp_tunnel_ops->forward_data_on_tunnel(
        ue_ipv4, &ue_ipv6, eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(), nullptr,
        DEFAULT_PRECEDENCE);
    if (rv < 0) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "ERROR in forwarding data on TUNNEL err=%d\n", rv);
    }
  } else {
    OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                    "Adding tunnel for bearer %u ue addr %x\n",
                    eps_bearer_ctxt_p->eps_bearer_id(), ue_ipv4.s_addr);
    if (eps_bearer_ctxt_p->eps_bearer_id() ==
        new_bearer_ctxt_info_p->sgw_eps_bearer_context()
            .pdn_connection()
            .default_bearer()) {
      // Set default precedence and tft for default bearer
      OAILOG_INFO_UE(
          LOG_SPGW_APP, imsi64,
          "Adding tunnel for ipv6 ue addr %s, enb %x, "
          "sgw_teid_s1u_s12_s4_up" TEID_FMT " enb_teid_S1u " TEID_FMT,
          ip6_str, enb.s_addr, eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
          eps_bearer_ctxt_p->enb_teid_s1u());

      if (is_ue_ipv6_pres) {
        rv = gtpv1u_add_tunnel(ue_ipv4, &ue_ipv6, vlan, enb, &enb_ipv6,
                               eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
                               eps_bearer_ctxt_p->enb_teid_s1u(), imsi, nullptr,
                               DEFAULT_PRECEDENCE, apn);
      } else {
        rv = gtpv1u_add_tunnel(ue_ipv4, nullptr, vlan, enb, &enb_ipv6,
                               eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
                               eps_bearer_ctxt_p->enb_teid_s1u(), imsi, nullptr,
                               DEFAULT_PRECEDENCE, apn);
      }

      // (@ulaskozat) We only need to update the TEIDs during session creation
      // which triggers rule installments on sessiond. When pipelined needs to
      // use eNB TEID for reporting we need to change this logic into something
      // more meaningful.
      bool update_teids = eps_bearer_ctxt_p->update_teids();
      eps_bearer_ctxt_p->set_update_teids(false);
      if (rv < 0) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "ERROR in setting up TUNNEL err=%d\n", rv);
      } else if (update_teids) {
        pcef_update_teids(reinterpret_cast<char*>(imsi.digit),
                          eps_bearer_ctxt_p->eps_bearer_id(),
                          eps_bearer_ctxt_p->enb_teid_s1u(),
                          eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up());
      }
    } else {
      for (uint8_t itrn = 0;
           itrn < eps_bearer_ctxt_p->tft().number_of_packet_filters(); ++itrn) {
        // Prepare DL flow rule
        magma::lte::oai::PacketFilter create_new_tft =
            eps_bearer_ctxt_p->tft().packet_filter_list().create_new_tft(itrn);
        packet_filter_t packet_filter = {0};
        proto_to_packet_filter(create_new_tft, &packet_filter);
        struct ip_flow_dl dlflow = {0};
        if (is_ue_ipv6_pres) {
          generate_dl_flow(&packet_filter.packetfiltercontents, ue_ipv4.s_addr,
                           &ue_ipv6, &dlflow);
        } else {
          generate_dl_flow(&packet_filter.packetfiltercontents, ue_ipv4.s_addr,
                           nullptr, &dlflow);
        }

        OAILOG_INFO_UE(
            LOG_SPGW_APP, imsi64,
            "Adding tunnel for ded bearer ipv6 ue addr %s, enb %x, "
            "sgw_teid_S1u_S12_S4_up" TEID_FMT " enb_teid_S1u " TEID_FMT,
            ip6_str, enb.s_addr, eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
            eps_bearer_ctxt_p->enb_teid_s1u());

        if (is_ue_ipv6_pres) {
          rv =
              gtpv1u_add_tunnel(ue_ipv4, &ue_ipv6, vlan, enb, &enb_ipv6,
                                eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
                                eps_bearer_ctxt_p->enb_teid_s1u(), imsi,
                                &dlflow, create_new_tft.eval_precedence(), apn);
        } else {
          rv =
              gtpv1u_add_tunnel(ue_ipv4, nullptr, vlan, enb, &enb_ipv6,
                                eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up(),
                                eps_bearer_ctxt_p->enb_teid_s1u(), imsi,
                                &dlflow, create_new_tft.eval_precedence(), apn);
        }

        if (rv < 0) {
          OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                          "ERROR in setting up TUNNEL err=%d\n", rv);
        } else {
          OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                         "Successfully setup flow rule for EPS bearer id %u "
                         "tunnel " TEID_FMT " (eNB) <-> (SGW) " TEID_FMT "\n",
                         eps_bearer_ctxt_p->eps_bearer_id(),
                         eps_bearer_ctxt_p->enb_teid_s1u(),
                         eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up());
        }
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

//------------------------------------------------------------------------------
/* Populates bearer contexts to be modified structure in
 * modify bearer rsp message
 */
static void sgw_populate_mbr_bearer_contexts_modified(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    itti_s11_modify_bearer_response_t* modify_response_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint8_t rsp_idx = 0;
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;

  for (uint8_t idx = 0; idx < resp_pP->num_bearers_modified; idx++) {
    if (sgw_cm_get_eps_bearer_entry(
            new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
                ->mutable_pdn_connection(),
            resp_pP->bearer_contexts_to_be_modified[idx].eps_bearer_id,
            &eps_bearer_ctxt) == magma::PROTO_MAP_OK) {
      OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                      "Rx SGI_UPDATE_ENDPOINT_RESPONSE: REQUEST_ACCEPTED\n");
      modify_response_p->bearer_contexts_modified.bearer_contexts[rsp_idx]
          .eps_bearer_id =
          resp_pP->bearer_contexts_to_be_modified[idx].eps_bearer_id;
      modify_response_p->bearer_contexts_modified.bearer_contexts[rsp_idx++]
          .cause.cause_value = REQUEST_ACCEPTED;
      modify_response_p->bearer_contexts_modified.num_bearer_context++;
      // if default bearer
      // #pragma message  "TODO define constant for default eps_bearer id"

#if !MME_UNIT_TEST  // skip tunnel creation for unit tests
      // setup GTPv1-U tunnel
      sgw_add_gtp_tunnel(imsi64, &eps_bearer_ctxt, new_bearer_ctxt_info_p);
#endif
      // may be removed
      if (TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX >
          eps_bearer_ctxt.num_sdf()) {
        uint8_t i = 0;
        while ((i < eps_bearer_ctxt.num_sdf()) &&
               (SDF_ID_NGBR_DEFAULT != eps_bearer_ctxt.sdf_ids(i)))
          i++;
        if (i >= eps_bearer_ctxt.num_sdf()) {
          eps_bearer_ctxt.add_sdf_ids(SDF_ID_NGBR_DEFAULT);
          eps_bearer_ctxt.set_num_sdf(eps_bearer_ctxt.num_sdf() + 1);
        }
      }
      sgw_update_eps_bearer_entry(
          new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
              ->mutable_pdn_connection(),
          eps_bearer_ctxt.eps_bearer_id(), &eps_bearer_ctxt);
    }
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}
//------------------------------------------------------------------------------
void sgw_handle_sgi_endpoint_updated(
    const itti_sgi_update_end_point_response_t* const resp_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  itti_s11_modify_bearer_response_t* modify_response_p = NULL;
  MessageDef* message_p = NULL;

  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Rx SGI_UPDATE_ENDPOINT_RESPONSE, Context teid " TEID_FMT
                  "\n",
                  resp_pP->context_teid);
  message_p = itti_alloc_new_message(TASK_SPGW_APP, S11_MODIFY_BEARER_RESPONSE);

  if (!message_p) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Failed to allocate memory for S11_MODIFY_BEARER_RESPONSE\n");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }

  modify_response_p = &message_p->ittiMsg.s11_modify_bearer_response;

  magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(resp_pP->context_teid);
  if (new_bearer_ctxt_info_p) {
    modify_response_p->teid =
        new_bearer_ctxt_info_p->sgw_eps_bearer_context().mme_teid_s11();
    modify_response_p->cause.cause_value = REQUEST_ACCEPTED;

    memcpy(modify_response_p->trxn,
           reinterpret_cast<const void*>(
               new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().c_str()),
           new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().size());
    message_p->ittiMsgHeader.imsi = imsi64;

    sgw_populate_mbr_bearer_contexts_modified(
        resp_pP, imsi64, new_bearer_ctxt_info_p, modify_response_p);
    sgw_populate_mbr_bearer_contexts_removed(
        resp_pP, imsi64,
        new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context(),
        modify_response_p);
    sgw_populate_mbr_bearer_contexts_not_found(LOG_SPGW_APP, resp_pP,
                                               modify_response_p);
    send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

//------------------------------------------------------------------------------
int sgw_handle_sgi_endpoint_deleted(
    const itti_sgi_delete_end_point_request_t* const resp_pP, imsi64_t imsi64) {
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  int rv = RETURNok;
  struct in_addr inaddr;
  struct in6_addr in6addr = {};

  OAILOG_FUNC_IN(LOG_SPGW_APP);

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, imsi64,
      "bcom Rx SGI_DELETE_ENDPOINT_REQUEST, Context teid %u, SGW S1U teid %u, "
      "EPS bearer id %u\n",
      resp_pP->context_teid, resp_pP->sgw_S1u_teid, resp_pP->eps_bearer_id);

  magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(resp_pP->context_teid);
  if (new_bearer_ctxt_info_p) {
    if (sgw_cm_get_eps_bearer_entry(
            new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
                ->mutable_pdn_connection(),
            resp_pP->eps_bearer_id, &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Rx SGI_DELETE_ENDPOINT_REQUEST: CONTEXT_NOT_FOUND "
                      "(pdn_connection.sgw_eps_bearers context) for bearer id",
                      resp_pP->eps_bearer_id);
    } else {
      OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                      "Rx SGI_DELETE_ENDPOINT_REQUEST: REQUEST_ACCEPTED\n");

      struct in_addr ue_ipv4 = {.s_addr = 0};
      struct in6_addr ue_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(eps_bearer_ctxt.mutable_ue_ip_paa(),
                                          &ue_ipv4, &ue_ipv6, true);
      // If the forwarding was suspended, first resume it.
      // Note that forward_data_on_tunnel does not install a new forwarding
      // rule, but simply deletes previously installed drop rule by
      // discard_data_on_tunnel.
      if (new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
              ->pdn_connection()
              .ue_suspended_for_ps_handover()) {
        rv = gtp_tunnel_ops->forward_data_on_tunnel(
            ue_ipv4, &ue_ipv6, eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
            nullptr, DEFAULT_PRECEDENCE);
        if (rv < 0) {
          OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                          "ERROR in resume forwarding data on TUNNEL err=%d\n",
                          rv);
        }
      }

#if !MME_UNIT_TEST  // skip tunnel deletion for unit tests
      // delete GTPv1-U tunnel
      struct in_addr enb = {.s_addr = 0};
      struct in6_addr enb_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(
          eps_bearer_ctxt.mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
          spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

      rv = gtp_tunnel_ops->del_tunnel(enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
                                      eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                                      eps_bearer_ctxt.enb_teid_s1u(), nullptr);
      if (rv < 0) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "ERROR in deleting TUNNEL\n");
      }
      // delete paging rule
      const char* ip_str = eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str();
      rv = gtp_tunnel_ops->delete_paging_rule(ue_ipv4, &ue_ipv6);
      if (rv < 0) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "ERROR in deleting paging rule for IP Addr: %s\n",
                        ip_str);
      } else {
        OAILOG_DEBUG(LOG_SPGW_APP, "Stopped paging for IP Addr: %s\n", ip_str);
      }
#endif

      std::string imsi =
          new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi();
      std::string apn = new_bearer_ctxt_info_p->sgw_eps_bearer_context()
                            .pdn_connection()
                            .apn_in_use();
      switch (resp_pP->paa.pdn_type) {
        case IPv4:
          inaddr = resp_pP->paa.ipv4_address;
          release_ue_ipv4_address(imsi, apn, &inaddr);
          OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                          "Released IPv4 PAA for PDN type IPv4\n");
          break;

        case IPv6:
          in6addr = resp_pP->paa.ipv6_address;
          release_ue_ipv6_address(imsi, apn, &in6addr);
          OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                          "Released IPv6 PAA for PDN type IPv6\n");
          break;

        case IPv4_AND_v6:
          inaddr = resp_pP->paa.ipv4_address;
          release_ue_ipv4_address(imsi, apn, &inaddr);
          OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                          "Released IPv4 PAA for PDN type IPv4_AND_v6\n");
          in6addr = resp_pP->paa.ipv6_address;
          release_ue_ipv6_address(imsi, apn, &in6addr);
          OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                          "Released IPv6 PAA for PDN type IPv4v6\n");
          break;

        default:
          Fatal("Bad paa.pdn_type %d", resp_pP->paa.pdn_type);
          break;
      }
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);
    }
  } else {
    OAILOG_DEBUG_UE(
        LOG_SPGW_APP, imsi64,
        "Rx SGI_DELETE_ENDPOINT_RESPONSE: CONTEXT_NOT_FOUND (S11 context)\n");
    /*    modify_response_p->teid = resp_pP->context_teid;    // TO BE CHECKED
    IF IT IS THIS TEID modify_response_p->bearer_present =
    MODIFY_BEARER_RESPONSE_REM;
    modify_response_p->bearer_choice.bearer_for_removal.eps_bearer_id =
    resp_pP->eps_bearer_id;
    modify_response_p->bearer_choice.bearer_for_removal.cause =
    CONTEXT_NOT_FOUND; modify_response_p->cause = CONTEXT_NOT_FOUND;
    modify_response_p->trxn = 0;
    rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);*/
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);
}

//------------------------------------------------------------------------------
// This function populates itti_sgi_update_end_point_response_t message

void populate_sgi_end_point_update(
    uint8_t sgi_rsp_idx, uint8_t idx,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt_p,
    itti_sgi_update_end_point_response_t* sgi_update_end_point_resp) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  FTEID_T_2_PROTO_IP(
      (&modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
            .s1_eNB_fteid),
      (eps_bearer_ctxt_p->mutable_enb_s1u_ip_addr()));
  eps_bearer_ctxt_p->set_enb_teid_s1u(
      modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
          .s1_eNB_fteid.teid);
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .sgw_S1u_teid = eps_bearer_ctxt_p->sgw_teid_s1u_s12_s4_up();
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .enb_S1u_teid = eps_bearer_ctxt_p->enb_teid_s1u();
  sgi_update_end_point_resp->bearer_contexts_to_be_modified[sgi_rsp_idx]
      .eps_bearer_id = eps_bearer_ctxt_p->eps_bearer_id();
  sgi_update_end_point_resp->num_bearers_modified++;

  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

//------------------------------------------------------------------------------
// This function populates and sends MBR failure message to MME APP
status_code_e send_mbr_failure(
    log_proto_t module,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64) {
  status_code_e rv = RETURNok;
  OAILOG_FUNC_IN(module);
  MessageDef* message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_MODIFY_BEARER_RESPONSE);

  if (!message_p) {
    OAILOG_ERROR_UE(module, imsi64,
                    "S11_MODIFY_BEARER_RESPONSE memory allocation failed\n");
    OAILOG_FUNC_RETURN(module, RETURNerror);
  }

  itti_s11_modify_bearer_response_t* modify_response_p =
      &message_p->ittiMsg.s11_modify_bearer_response;

  for (uint8_t idx = 0;
       idx <
       modify_bearer_pP->bearer_contexts_to_be_modified.num_bearer_context;
       idx++) {
    modify_response_p->bearer_contexts_marked_for_removal.bearer_contexts[idx]
        .eps_bearer_id =
        modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
            .eps_bearer_id;
    modify_response_p->bearer_contexts_marked_for_removal.bearer_contexts[idx]
        .cause.cause_value = CONTEXT_NOT_FOUND;
  }
  // Fill mme s11 teid received in modify bearer request
  modify_response_p->teid = modify_bearer_pP->local_teid;
  modify_response_p->bearer_contexts_marked_for_removal.num_bearer_context += 1;
  modify_response_p->cause.cause_value = CONTEXT_NOT_FOUND;
  modify_response_p->trxn = modify_bearer_pP->trxn;
  OAILOG_DEBUG_UE(module, imsi64,
                  "Rx MODIFY_BEARER_REQUEST, teid " TEID_FMT
                  " CONTEXT_NOT_FOUND\n",
                  modify_bearer_pP->teid);
  message_p->ittiMsgHeader.imsi = imsi64;
  rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);

  OAILOG_FUNC_RETURN(module, rv);
}

//------------------------------------------------------------------------------
status_code_e sgw_handle_modify_bearer_request(
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  status_code_e rv = RETURNok;
  uint8_t idx = 0;
  itti_sgi_update_end_point_response_t sgi_update_end_point_resp = {0};
  struct in_addr enb = {.s_addr = 0};

  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Rx MODIFY_BEARER_REQUEST, teid " TEID_FMT "\n",
                  modify_bearer_pP->teid);

  magma::lte::oai::S11BearerContext* bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(modify_bearer_pP->teid);

  if (bearer_ctxt_info_p) {
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
        bearer_ctxt_info_p->mutable_sgw_eps_bearer_context();
    sgw_context_p->mutable_pdn_connection()->set_default_bearer(
        modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[0]
            .eps_bearer_id);
    if (modify_bearer_pP->trxn) {
      sgw_context_p->set_trxn(reinterpret_cast<char*>(modify_bearer_pP->trxn));
    }

    sgi_update_end_point_resp.context_teid = modify_bearer_pP->teid;
    uint8_t sgi_rsp_idx = 0;
    for (idx = 0;
         idx <
         modify_bearer_pP->bearer_contexts_to_be_modified.num_bearer_context;
         idx++) {
      magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
      uint32_t eps_bearer_id =
          modify_bearer_pP->bearer_contexts_to_be_modified.bearer_contexts[idx]
              .eps_bearer_id;
      if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                      eps_bearer_id, &eps_bearer_ctxt) !=
          magma::PROTO_MAP_OK) {
        sgi_update_end_point_resp.bearer_contexts_not_found[sgi_rsp_idx++] =
            eps_bearer_id;
        sgi_update_end_point_resp.num_bearers_not_found++;
      } else {  // eps_bearer_ctxt found
        inet_pton(AF_INET,
                  eps_bearer_ctxt.enb_s1u_ip_addr().ipv4_addr().c_str(),
                  &enb.s_addr);

        struct in6_addr enb_ipv6 = {};
        if (spgw_config.sgw_config.ipv6.s1_ipv6_enabled &&
            eps_bearer_ctxt.enb_s1u_ip_addr().pdn_type() == IPv6) {
          inet_pton(AF_INET6,
                    eps_bearer_ctxt.enb_s1u_ip_addr().ipv6_addr().c_str(),
                    &enb_ipv6);
        }
        // Send end marker to eNB and then delete the tunnel if enb_ip is
        // different
        if (does_bearer_context_hold_valid_enb_ip(
                eps_bearer_ctxt.enb_s1u_ip_addr()) &&
            is_enb_ip_address_same(
                &modify_bearer_pP->bearer_contexts_to_be_modified
                     .bearer_contexts[idx]
                     .s1_eNB_fteid,
                enb, enb_ipv6) == false) {
          struct in_addr ue_ipv4 = {.s_addr = 0};
          inet_pton(AF_INET, eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str(),
                    &ue_ipv4.s_addr);
          struct in6_addr ue_ipv6 = {};
          if ((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6) ||
              (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) {
            inet_pton(AF_INET6, eps_bearer_ctxt.ue_ip_paa().ipv6_addr().c_str(),
                      &ue_ipv6);
          }

          OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                          "Delete GTPv1-U tunnel for sgw_teid " TEID_FMT
                          "for bearer %d\n",
                          eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                          eps_bearer_ctxt.eps_bearer_id());
          if (gtp_tunnel_ops) {
            // This is best effort, ignore return code.
            gtp_tunnel_ops->send_end_marker(enb, modify_bearer_pP->teid);
            // delete GTPv1-U tunnel
            if (gtp_tunnel_ops->del_tunnel(
                    enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
                    eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                    eps_bearer_ctxt.enb_teid_s1u(), nullptr) < 0) {
              OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                              "ERROR in deleting TUNNEL " TEID_FMT
                              " (eNB) <-> (SGW) " TEID_FMT "\n",
                              eps_bearer_ctxt.enb_teid_s1u(),
                              eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up());
            }
          }
        }
        populate_sgi_end_point_update(sgi_rsp_idx, idx, modify_bearer_pP,
                                      &eps_bearer_ctxt,
                                      &sgi_update_end_point_resp);
        sgi_rsp_idx++;
        sgw_update_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                    eps_bearer_ctxt.eps_bearer_id(),
                                    &eps_bearer_ctxt);
      }
    }  // for loop
    sgi_rsp_idx = 0;
    for (idx = 0;
         idx <
         modify_bearer_pP->bearer_contexts_to_be_removed.num_bearer_context;
         idx++) {
      magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
      uint32_t eps_bearer_id =
          modify_bearer_pP->bearer_contexts_to_be_removed.bearer_contexts[idx]
              .eps_bearer_id;
      if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                      eps_bearer_id, &eps_bearer_ctxt) !=
          magma::PROTO_MAP_OK) {
        sgi_update_end_point_resp.bearer_contexts_to_be_removed[sgi_rsp_idx++] =
            eps_bearer_ctxt.eps_bearer_id();
        sgi_update_end_point_resp.num_bearers_removed++;
      }
    }
    sgw_handle_sgi_endpoint_updated(&sgi_update_end_point_resp, imsi64);
  } else {  // bearer_ctxt_info_p not found
    rv = send_mbr_failure(LOG_SPGW_APP, modify_bearer_pP, imsi64);
    if (rv != RETURNok) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, imsi64,
          "Error in sending modify bearer response to MME App for the failed "
          "bearers, teid" TEID_FMT "\n",
          modify_bearer_pP->teid);
    }
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);
}

//------------------------------------------------------------------------------
status_code_e sgw_handle_delete_session_request(
    const itti_s11_delete_session_request_t* const delete_session_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  itti_s11_delete_session_response_t* delete_session_resp_p = NULL;
  MessageDef* message_p = NULL;
  status_code_e rv = RETURNok;

  increment_counter("spgw_delete_session", 1, NO_LABELS);
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_DELETE_SESSION_RESPONSE);

  if (!message_p) {
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  delete_session_resp_p = &message_p->ittiMsg.s11_delete_session_response;
  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Handle delete session request for sgw_s11_teid " TEID_FMT
                 "\n",
                 delete_session_req_pP->teid);

  if (delete_session_req_pP->indication_flags.oi) {
    OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                    "OI flag is set for this message indicating the request"
                    "should be forwarded to P-GW entity\n");
  }

  magma::lte::oai::S11BearerContext* ctx_p =
      sgw_cm_get_spgw_context(delete_session_req_pP->teid);
  if (ctx_p) {
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
        ctx_p->mutable_sgw_eps_bearer_context();
    if ((delete_session_req_pP->sender_fteid_for_cp.ipv4) &&
        (delete_session_req_pP->sender_fteid_for_cp.ipv6)) {
      /*
       * Sender F-TEID IE present
       */
      if (delete_session_req_pP->teid != sgw_context_p->mme_teid_s11()) {
        delete_session_resp_p->teid = sgw_context_p->mme_teid_s11();
        delete_session_resp_p->cause.cause_value = INVALID_PEER;
        OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64, "Mismatch in MME Teid for CP\n");
      } else {
        delete_session_resp_p->teid =
            delete_session_req_pP->sender_fteid_for_cp.teid;
      }
    } else {
      delete_session_resp_p->cause.cause_value = REQUEST_ACCEPTED;
      delete_session_resp_p->teid = sgw_context_p->mme_teid_s11();

      // TODO make async
      std::string imsi = sgw_context_p->imsi();
      std::string apn = sgw_context_p->pdn_connection().apn_in_use();
      pcef_end_session(imsi, apn);
      itti_sgi_delete_end_point_request_t sgi_delete_end_point_request;
      magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;

      for (int ebix = 0; ebix < BEARERS_PER_UE; ebix++) {
        ebi_t ebi = INDEX_TO_EBI(ebix);
        if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                        ebi, &eps_bearer_ctxt) ==
            magma::PROTO_MAP_OK) {
          if (ebi != delete_session_req_pP->lbi) {
            struct in_addr enb = {.s_addr = 0};
            struct in6_addr enb_ipv6 = {};
            convert_proto_ip_to_standard_ip_fmt(
                eps_bearer_ctxt.mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
                spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

            struct in_addr ue_ipv4 = {.s_addr = 0};
            struct in6_addr ue_ipv6 = {};
            convert_proto_ip_to_standard_ip_fmt(
                eps_bearer_ctxt.mutable_ue_ip_paa(), &ue_ipv4, &ue_ipv6, true);
#if !MME_UNIT_TEST  // skip tunnel deletion for unit tests
            if (gtp_tunnel_ops->del_tunnel(
                    enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
                    eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                    eps_bearer_ctxt.enb_teid_s1u(), nullptr) < 0) {
              OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                              "ERROR in deleting TUNNEL " TEID_FMT
                              " (eNB) <-> (SGW) " TEID_FMT "\n",
                              eps_bearer_ctxt.enb_teid_s1u(),
                              eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up());
            }
#endif
            eps_bearer_ctxt.set_num_sdf(0);
          }
          if (sgw_update_eps_bearer_entry(
                  sgw_context_p->mutable_pdn_connection(), ebi,
                  &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
            OAILOG_ERROR_UE(
                LOG_SPGW_APP, sgw_context_p->imsi64(),
                "Failed to update bearer context for bearer id :%u\n", ebi);
            OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
          }
        }
      }
      if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                      delete_session_req_pP->lbi,
                                      &eps_bearer_ctxt) ==
          magma::PROTO_MAP_OK) {
        eps_bearer_ctxt.set_num_sdf(0);

        sgi_delete_end_point_request.context_teid = delete_session_req_pP->teid;
        sgi_delete_end_point_request.sgw_S1u_teid =
            eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up();
        sgi_delete_end_point_request.eps_bearer_id = delete_session_req_pP->lbi;
        sgi_delete_end_point_request.pdn_type =
            (pdn_type_value_t)sgw_context_p->saved_message().pdn_type();

        sgi_delete_end_point_request.paa.pdn_type =
            (pdn_type_value_t)eps_bearer_ctxt.ue_ip_paa().pdn_type();
        convert_proto_ip_to_standard_ip_fmt(
            eps_bearer_ctxt.mutable_ue_ip_paa(),
            &sgi_delete_end_point_request.paa.ipv4_address,
            &sgi_delete_end_point_request.paa.ipv6_address, true);
        if (sgw_update_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                        delete_session_req_pP->lbi,
                                        &eps_bearer_ctxt) !=
            magma::PROTO_MAP_OK) {
          OAILOG_ERROR_UE(LOG_SPGW_APP, sgw_context_p->imsi64(),
                          "Failed to update bearer context for bearer id :%u\n",
                          delete_session_req_pP->lbi);
          OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
        }

        sgw_handle_sgi_endpoint_deleted(&sgi_delete_end_point_request, imsi64);
      } else {
        OAILOG_WARNING_UE(
            LOG_SPGW_APP, imsi64,
            "Can't find eps_bearer_entry for MME TEID " TEID_FMT " lbi %u\n",
            delete_session_req_pP->teid, delete_session_req_pP->lbi);
      }

      /*
       * Remove eps bearer context, S11 bearer context and s11 tunnel
       */
      sgw_cm_remove_bearer_context_information(delete_session_req_pP->teid,
                                               imsi64);
      increment_counter("spgw_delete_session", 1, 1, "result", "success");
    }

    delete_session_resp_p->trxn = delete_session_req_pP->trxn;
    delete_session_resp_p->peer_ip.s_addr =
        delete_session_req_pP->peer_ip.s_addr;

    delete_session_resp_p->lbi = delete_session_req_pP->lbi;

    message_p->ittiMsgHeader.imsi = imsi64;
    rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);

  } else {
    /*
     * Context not found... set the cause to CONTEXT_NOT_FOUND
     * * * * 3GPP TS 29.274 #7.2.10.1
     */

    if ((delete_session_req_pP->sender_fteid_for_cp.ipv4 == 0) &&
        (delete_session_req_pP->sender_fteid_for_cp.ipv6 == 0)) {
      delete_session_resp_p->teid = 0;
    } else {
      delete_session_resp_p->teid =
          delete_session_req_pP->sender_fteid_for_cp.teid;
    }

    delete_session_resp_p->cause.cause_value = CONTEXT_NOT_FOUND;
    delete_session_resp_p->trxn = delete_session_req_pP->trxn;
    delete_session_resp_p->peer_ip.s_addr =
        delete_session_req_pP->peer_ip.s_addr;
    delete_session_resp_p->lbi = delete_session_req_pP->lbi;

    message_p->ittiMsgHeader.imsi = imsi64;
    rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
    increment_counter("spgw_delete_session", 1, 2, "result", "failure", "cause",
                      "context_not_found");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rv);
  }

  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
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
void sgw_handle_release_access_bearers_request(
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Release Access Bearer Request Received in SGW\n");

  spgw_ue_context_t* ue_context_p = nullptr;
  gtpv2c_cause_value_t cause = CONTEXT_NOT_FOUND;
  magma::lte::oai::S11BearerContext* ctx_p = nullptr;

  map_uint64_spgw_ue_context_t* state_ue_map = get_spgw_ue_state();
  if (!state_ue_map) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get spgw_ue_state");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  state_ue_map->get(imsi64, &ue_context_p);

  if (ue_context_p) {
    sgw_s11_teid_t* s11_teid_p = nullptr;
    LIST_FOREACH(s11_teid_p, &ue_context_p->sgw_s11_teid_list, entries) {
      if (s11_teid_p) {
        ctx_p = sgw_cm_get_spgw_context(s11_teid_p->sgw_s11_teid);
        if (ctx_p) {
          sgw_process_release_access_bearer_request(
              LOG_SPGW_APP, imsi64, ctx_p->mutable_sgw_eps_bearer_context());
          cause = REQUEST_ACCEPTED;
        }
      }
    }
  }
  sgw_send_release_access_bearer_response(
      LOG_SPGW_APP, TASK_SPGW_APP, imsi64, cause, release_access_bearers_req_pP,
      ctx_p ? ctx_p->sgw_eps_bearer_context().mme_teid_s11() : 0);
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

//-------------------------------------------------------------------------
void handle_s5_create_session_response(
    spgw_state_t* state,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    s5_create_session_response_t session_resp) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  itti_s11_create_session_response_t* create_session_response_p = NULL;
  MessageDef* message_p = NULL;
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp = {};
  gtpv2c_cause_value_t cause = REQUEST_ACCEPTED;

  /* Since bearer context is not found, can not get mme_s11_teid, imsi64,
   * so Create Session Response will not be sent
   */
  if (!new_bearer_ctxt_info_p) {
    OAILOG_ERROR(
        LOG_SPGW_APP,
        "Failed to fetch sgw bearer context from sgw s11 teid: " TEID_FMT "\n",
        session_resp.context_teid);
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
      "Handle s5_create_session_response, for Context SGW S11 teid, " TEID_FMT
      "EPS bearer id %u\n",
      session_resp.context_teid, session_resp.eps_bearer_id);

  // PCO processing
  magma::lte::oai::Pco* pco_req =
      new_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_saved_message()
          ->mutable_pco();
  protocol_configuration_options_t pco_resp = {0};
  protocol_configuration_options_ids_t pco_ids;
  memset(&pco_ids, 0, sizeof pco_ids);

  if (pgw_process_pco_request(pco_req, &pco_resp, &pco_ids) != RETURNok) {
    OAILOG_ERROR_UE(LOG_SPGW_APP,
                    new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
                    "Error in processing PCO in create session request for "
                    "context_id: " TEID_FMT "\n",
                    session_resp.context_teid);
    session_resp.failure_cause = S5_OK;
    session_resp.status = SGI_STATUS_ERROR_FAILED_TO_PROCESS_PCO;
  }
  copy_protocol_configuration_options(&sgi_create_endpoint_resp.pco, &pco_resp);
  clear_protocol_configuration_options(&pco_resp);

  // Fill SGi create endpoint resp data
  sgi_create_endpoint_resp.status = session_resp.status;
  sgi_create_endpoint_resp.context_teid = session_resp.context_teid;
  sgi_create_endpoint_resp.eps_bearer_id = session_resp.eps_bearer_id;
  sgi_create_endpoint_resp.paa.pdn_type =
      (pdn_type_value_t)new_bearer_ctxt_info_p->sgw_eps_bearer_context()
          .saved_message()
          .pdn_type();

  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
      "Status of SGI_CREATE_ENDPOINT_RESPONSE within S5_CREATE_BEARER_RESPONSE "
      "is: %u\n",
      sgi_create_endpoint_resp.status);

  if (session_resp.failure_cause == S5_OK) {
    switch (sgi_create_endpoint_resp.status) {
      case SGI_STATUS_OK:
        // Send Create Session Response with ack
        sgw_handle_sgi_endpoint_created(
            state, &sgi_create_endpoint_resp,
            new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64());
        increment_counter("spgw_create_session", 1, 1, "result", "success");
        OAILOG_FUNC_OUT(LOG_SPGW_APP);

      case SGI_STATUS_ERROR_CONTEXT_NOT_FOUND:
        cause = CONTEXT_NOT_FOUND;
        increment_counter("spgw_create_session", 1, 1, "result", "failure",
                          "cause", "context_not_found");

        break;

      case SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED:
        cause = ALL_DYNAMIC_ADDRESSES_ARE_OCCUPIED;
        increment_counter("spgw_create_session", 1, 1, "result", "failure",
                          "cause", "resource_not_available");

        break;

      case SGI_STATUS_ERROR_SERVICE_NOT_SUPPORTED:
        cause = SERVICE_NOT_SUPPORTED;
        increment_counter("spgw_create_session", 1, 1, "result", "failure",
                          "cause", "pdn_type_ipv6_not_supported");

        break;
      case SGI_STATUS_ERROR_FAILED_TO_PROCESS_PCO:
        cause = REQUEST_REJECTED;
        increment_counter("spgw_create_session", 1, 1, "result", "failure",
                          "cause", "failed_to_process_pco_req");
        break;
      default:
        cause = REQUEST_REJECTED;  // Unspecified reason

        break;
    }
  } else if (session_resp.failure_cause == PCEF_FAILURE) {
    cause = SERVICE_DENIED;
  } else if (session_resp.failure_cause == IP_ALLOCATION_FAILURE) {
    cause = SYSTEM_FAILURE;
  }
  // Send Create Session Response with Nack
  message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_CREATE_SESSION_RESPONSE);
  if (!message_p) {
    OAILOG_ERROR_UE(LOG_SPGW_APP,
                    new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
                    "Message Create Session Response allocation failed\n");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  create_session_response_p = &message_p->ittiMsg.s11_create_session_response;
  memset(create_session_response_p, 0,
         sizeof(itti_s11_create_session_response_t));
  create_session_response_p->cause.cause_value = cause;
  create_session_response_p->bearer_contexts_marked_for_removal
      .bearer_contexts[0]
      .cause.cause_value = cause;
  create_session_response_p->bearer_contexts_marked_for_removal
      .num_bearer_context += 1;
  create_session_response_p->bearer_contexts_marked_for_removal
      .bearer_contexts[0]
      .eps_bearer_id = session_resp.eps_bearer_id;
  create_session_response_p->teid =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context().mme_teid_s11();
  memcpy(create_session_response_p->trxn,
         reinterpret_cast<const void*>(
             new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().c_str()),
         new_bearer_ctxt_info_p->sgw_eps_bearer_context().trxn().size());
  message_p->ittiMsgHeader.imsi =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64();
  OAILOG_DEBUG_UE(
      LOG_SPGW_APP, new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
      "Sending S11 Create Session Response to MME, MME S11 teid = %u\n",
      create_session_response_p->teid);

  send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);

  /* Remove the default bearer context entry already created as create session
   * response failure is received
   */
  OAILOG_INFO_UE(
      LOG_SPGW_APP, new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64(),
      "Deleted default bearer context with SGW C-plane TEID =" TEID_FMT
      "as create session response failure is received\n",
      new_bearer_ctxt_info_p->sgw_eps_bearer_context().mme_teid_s11());
  sgw_cm_remove_bearer_context_information(
      session_resp.context_teid,
      new_bearer_ctxt_info_p->sgw_eps_bearer_context().imsi64());

  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

/*
 * Handle Suspend Notification from MME, set the state of default bearer to
 * suspend and discard the DL data for this UE and delete the GTPv1-U tunnel
 * TODO for multiple PDN support, suspend all bearers and disard the DL data for
 * the UE
 */
status_code_e sgw_handle_suspend_notification(
    const itti_s11_suspend_notification_t* const suspend_notification_pP,
    imsi64_t imsi64) {
  itti_s11_suspend_acknowledge_t* suspend_acknowledge_p = NULL;
  MessageDef* message_p = NULL;
  status_code_e rv = RETURNok;

  OAILOG_FUNC_IN(LOG_SPGW_APP);
  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64, "Rx SUSPEND_NOTIFICATION, teid %u\n",
                  suspend_notification_pP->teid);

  message_p = itti_alloc_new_message(TASK_SPGW_APP, S11_SUSPEND_ACKNOWLEDGE);

  if (!message_p) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Unable to allocate itti message: S11_SUSPEND_ACKNOWLEDGE \n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  suspend_acknowledge_p = &message_p->ittiMsg.s11_suspend_acknowledge;
  memset((void*)suspend_acknowledge_p, 0,
         sizeof(itti_s11_suspend_acknowledge_t));

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  magma::lte::oai::S11BearerContext* ctx_p =
      sgw_cm_get_spgw_context(suspend_notification_pP->teid);
  if (ctx_p) {
    ctx_p->mutable_sgw_eps_bearer_context()
        ->mutable_pdn_connection()
        ->set_ue_suspended_for_ps_handover(true);
    suspend_acknowledge_p->cause.cause_value = REQUEST_ACCEPTED;
    suspend_acknowledge_p->teid =
        ctx_p->sgw_eps_bearer_context().mme_teid_s11();
    /*
     * TODO Need to discard the DL data
     * deleting the GTPV1-U tunnel in suspended mode
     * This tunnel will be added again when UE moves back to connected mode.
     */
    map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
    eps_bearer_map.map = ctx_p->mutable_sgw_eps_bearer_context()
                             ->mutable_pdn_connection()
                             ->mutable_eps_bearer_map();
    if (eps_bearer_map.get(suspend_notification_pP->lbi, &eps_bearer_ctxt) ==
        magma::PROTO_MAP_OK) {
      OAILOG_DEBUG_UE(
          LOG_SPGW_APP, imsi64,
          "Handle S11_SUSPEND_NOTIFICATION: Discard the Data received GTP-U "
          "Tunnel mapping in"
          "GTP-U Kernel module \n");
      // delete GTPv1-U tunnel
      struct in_addr ue_ipv4 = {.s_addr = 0};
      struct in6_addr ue_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(eps_bearer_ctxt.mutable_ue_ip_paa(),
                                          &ue_ipv4, &ue_ipv6, true);
#if !MME_UNIT_TEST
      if (gtp_tunnel_ops->discard_data_on_tunnel(
              ue_ipv4, &ue_ipv6, eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
              NULL) < 0) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "ERROR in Disabling DL data on TUNNEL\n");
      }
#endif
    } else {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "Bearer context not found \n");
    }
  } else {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Sending Suspend Acknowledge for sgw_s11_teid :%d for context not "
        "found "
        "\n",
        suspend_notification_pP->teid);
    suspend_acknowledge_p->cause.cause_value = CONTEXT_NOT_FOUND;
    suspend_acknowledge_p->teid = 0;
  }

  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Send Suspend acknowledge for teid :%d\n",
                 suspend_acknowledge_p->teid);
  message_p->ittiMsgHeader.imsi = imsi64;
  rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rv);
}

/*
 * Handle NW initiated Dedicated Bearer Activation Rsp from MME
 */

status_code_e sgw_handle_nw_initiated_actv_bearer_rsp(
    const itti_s11_nw_init_actv_bearer_rsp_t* const s11_actv_bearer_rsp,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint32_t msg_bearer_index = 0;
  status_code_e rc = RETURNerror;
  gtpv2c_cause_value_t cause = REQUEST_REJECTED;
  bearer_context_within_create_bearer_response_t bearer_context = {0};
  char policy_rule_name[POLICY_RULE_NAME_MAXLEN + 1];
  ebi_t default_bearer_id;

  bearer_context =
      s11_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index];
  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Received nw_initiated_bearer_actv_rsp from MME with EBI %u\n",
                 bearer_context.eps_bearer_id);

  magma::lte::oai::S11BearerContext* spgw_context =
      sgw_cm_get_spgw_context(s11_actv_bearer_rsp->sgw_s11_teid);
  if (!spgw_context) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Error in retrieving s_plus_p_gw context from sgw_s11_teid " TEID_FMT
        "\n",
        s11_actv_bearer_rsp->sgw_s11_teid);
    handle_failed_create_bearer_response(
        spgw_context->mutable_sgw_eps_bearer_context(),
        s11_actv_bearer_rsp->cause.cause_value, imsi64, &bearer_context,
        nullptr, LOG_SPGW_APP);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
  }

  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      spgw_context->mutable_sgw_eps_bearer_context();
  default_bearer_id = sgw_context_p->pdn_connection().default_bearer();

  //--------------------------------------
  // EPS bearer entry
  //--------------------------------------
  // TODO multiple bearers
  if (!(sgw_context_p->pending_procedures_size())) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "Failed to get create bearer procedure from temporary stored context, "
        "so did not create new EPS bearer entry for EBI %u\n",
        bearer_context.eps_bearer_id);
    handle_failed_create_bearer_response(
        sgw_context_p, s11_actv_bearer_rsp->cause.cause_value, imsi64,
        &bearer_context, nullptr, LOG_SPGW_APP);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
  }
  // If UE did not accept the request send reject to NW
  if (s11_actv_bearer_rsp->cause.cause_value != REQUEST_ACCEPTED) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Did not create new EPS bearer entry as "
                    "UE rejected the request for EBI %u\n",
                    bearer_context.eps_bearer_id);
    handle_failed_create_bearer_response(
        sgw_context_p, s11_actv_bearer_rsp->cause.cause_value, imsi64,
        &bearer_context, nullptr, LOG_SPGW_APP);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
  }
  magma::lte::oai::PgwCbrProcedure* pgw_ni_cbr_proc = nullptr;
  uint8_t num_of_bearers_deleted = 0;
  uint8_t num_of_pending_procedures = sgw_context_p->pending_procedures_size();
  for (uint8_t proc_index = 0;
       proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
    pgw_ni_cbr_proc = sgw_context_p->mutable_pending_procedures(proc_index);
    if (!pgw_ni_cbr_proc) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, sgw_context_p->imsi64(),
          "Pending procedure within sgw_context is null for proc_index:%u",
          proc_index);
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
    }
    if (pgw_ni_cbr_proc->type() ==
        PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST) {
      num_of_bearers_deleted = pgw_ni_cbr_proc->pending_eps_bearers_size();
      for (uint8_t bearer_index = 0;
           bearer_index < pgw_ni_cbr_proc->pending_eps_bearers_size();
           bearer_index++) {
        magma::lte::oai::SgwEpsBearerContext* bearer_context_proto =
            pgw_ni_cbr_proc->mutable_pending_eps_bearers(bearer_index);
        if (!bearer_context_proto) {
          OAILOG_ERROR_UE(LOG_SPGW_APP, sgw_context_p->imsi64(),
                          "Bearer context within pending procedure is null for "
                          "proc_index:%u bearer_index:%u",
                          proc_index, bearer_index);
          OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
        }

        if (bearer_context_proto->sgw_teid_s1u_s12_s4_up() ==
            bearer_context.s1u_sgw_fteid.teid) {
          bearer_context_proto->set_eps_bearer_id(bearer_context.eps_bearer_id);
          FTEID_T_2_PROTO_IP(&bearer_context.s1u_enb_fteid,
                             bearer_context_proto->mutable_enb_s1u_ip_addr());
          bearer_context_proto->set_enb_teid_s1u(
              bearer_context.s1u_enb_fteid.teid);
          if (sgw_cm_insert_eps_bearer_ctxt_in_collection(
                  sgw_context_p->mutable_pdn_connection(),
                  bearer_context_proto) != magma::PROTO_MAP_OK) {
            OAILOG_ERROR_UE(
                LOG_SPGW_APP, imsi64,
                "Failed to create new EPS bearer entry for bearer id: %u\n",
                bearer_context.eps_bearer_id);
            increment_counter("s11_actv_bearer_rsp", 1, 2, "result", "failure",
                              "cause", "internal_software_error");
          } else {
            OAILOG_INFO_UE(
                LOG_SPGW_APP, imsi64,
                "Successfully created new EPS bearer entry with EBI %u\n",
                bearer_context.eps_bearer_id);

            cause = REQUEST_ACCEPTED;
            snprintf(policy_rule_name, POLICY_RULE_NAME_MAXLEN + 1, "%s",
                     bearer_context_proto->policy_rule_name().c_str());
            // setup GTPv1-U tunnel for each packet filter
            // enb, UE and imsi are common across rules
            add_tunnel_helper(spgw_context, bearer_context_proto, imsi64);
          }
          --num_of_bearers_deleted;
        }
        if (num_of_bearers_deleted == 0) {
          pgw_ni_cbr_proc->clear_pending_eps_bearers();
          --num_of_pending_procedures;
        }
      }  // end of bearer index loop
    }
  }  // end of procedure index loop
  if (num_of_pending_procedures == 0) {
    sgw_context_p->clear_pending_procedures();
  }

  // Send ACTIVATE_DEDICATED_BEARER_RSP to PCRF
  rc = spgw_send_nw_init_activate_bearer_rsp(
      cause, imsi64, &bearer_context, default_bearer_id, policy_rule_name);
  if (rc != RETURNok) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to send ACTIVATE_DEDICATED_BEARER_RSP to PCRF\n");
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

/*
 * Handle NW-initiated dedicated bearer dectivation rsp from MME
 */

status_code_e sgw_handle_nw_initiated_deactv_bearer_rsp(
    spgw_state_t* spgw_state,
    const itti_s11_nw_init_deactv_bearer_rsp_t* const
        s11_pcrf_ded_bearer_deactv_rsp,
    imsi64_t imsi64) {
  status_code_e rc = RETURNok;
  uint32_t i = 0;
  uint32_t no_of_bearers = 0;
  ebi_t ebi = {0};
  itti_sgi_delete_end_point_request_t sgi_delete_end_point_request;

  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Received nw_initiated_deactv_bearer_rsp from MME\n");

  no_of_bearers =
      s11_pcrf_ded_bearer_deactv_rsp->bearer_contexts.num_bearer_context;
  //--------------------------------------
  // Get EPS bearer entry
  //--------------------------------------
  magma::lte::oai::S11BearerContext* spgw_ctxt =
      sgw_cm_get_spgw_context(s11_pcrf_ded_bearer_deactv_rsp->s_gw_teid_s11_s4);
  if (!spgw_ctxt) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "failed to get s_plus_p_gw_eps_bearer_context for teid" TEID_FMT,
        s11_pcrf_ded_bearer_deactv_rsp->s_gw_teid_s11_s4);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
  }
  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      spgw_ctxt->mutable_sgw_eps_bearer_context();
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  // Remove the default bearer entry
  if (s11_pcrf_ded_bearer_deactv_rsp->delete_default_bearer) {
    if (!s11_pcrf_ded_bearer_deactv_rsp->lbi) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64, "LBI received from MME is NULL\n");
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
    }
    // Delete all the dedicated bearers linked to this default bearer
    std::vector<uint32_t> remove_bearer_list = {};
    if (sgw_context_p->mutable_pdn_connection()->eps_bearer_map_size()) {
      map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
      eps_bearer_map.map =
          sgw_context_p->mutable_pdn_connection()->mutable_eps_bearer_map();
      for (auto itr = eps_bearer_map.map->begin();
           itr != eps_bearer_map.map->end(); itr++) {
        magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt = itr->second;

        if (eps_bearer_ctxt.eps_bearer_id() !=
            *s11_pcrf_ded_bearer_deactv_rsp->lbi) {
          struct in_addr enb = {.s_addr = 0};
          struct in6_addr enb_ipv6 = {};
          convert_proto_ip_to_standard_ip_fmt(
              eps_bearer_ctxt.mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
              spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

          struct in_addr ue_ipv4 = {.s_addr = 0};
          struct in6_addr ue_ipv6 = {};
          convert_proto_ip_to_standard_ip_fmt(
              eps_bearer_ctxt.mutable_ue_ip_paa(), &ue_ipv4, &ue_ipv6, true);
          bool is_ue_ipv6_pres = false;
          if (((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6 ||
                eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) &&
              !(eps_bearer_ctxt.ue_ip_paa().ipv6_addr().empty())) {
            is_ue_ipv6_pres = true;
          }

#if !MME_UNIT_TEST
          int rc = RETURNerror;
          if (is_ue_ipv6_pres) {
            rc = gtp_tunnel_ops->del_tunnel(
                enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
                eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                eps_bearer_ctxt.enb_teid_s1u(), nullptr);
          } else {
            rc = gtp_tunnel_ops->del_tunnel(
                enb, &enb_ipv6, ue_ipv4, nullptr,
                eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                eps_bearer_ctxt.enb_teid_s1u(), nullptr);
          }
          if (rc != RETURNok) {
            OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                            "ERROR in deleting TUNNEL " TEID_FMT
                            " (eNB) <-> (SGW) " TEID_FMT "\n",
                            eps_bearer_ctxt.enb_teid_s1u(),
                            eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up());
          } else {
            OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                           "Removed dedicated bearer context for (ebi = %d)\n",
                           ebi);
          }
#endif
          remove_bearer_list.push_back(eps_bearer_ctxt.eps_bearer_id());
        }
      }
      for (uint8_t i = 0; i < remove_bearer_list.size(); i++) {
        sgw_remove_eps_bearer_context(sgw_context_p->mutable_pdn_connection(),
                                      remove_bearer_list[i]);
      }
      remove_bearer_list.clear();
    }
    OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                   "Removed default bearer context for (ebi = %d)\n",
                   *s11_pcrf_ded_bearer_deactv_rsp->lbi);
    ebi = *s11_pcrf_ded_bearer_deactv_rsp->lbi;
    if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                    ebi,
                                    &eps_bearer_ctxt) != magma::PROTO_MAP_OK) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Failed to get eps_bearer_ctxt for bearer id:%u", ebi);
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
    }

    sgi_delete_end_point_request.context_teid =
        sgw_context_p->sgw_teid_s11_s4();
    sgi_delete_end_point_request.sgw_S1u_teid =
        eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up();
    sgi_delete_end_point_request.eps_bearer_id = ebi;
    sgi_delete_end_point_request.pdn_type =
        sgw_context_p->saved_message().pdn_type();

    sgi_delete_end_point_request.paa.pdn_type =
        (pdn_type_value_t)eps_bearer_ctxt.ue_ip_paa().pdn_type();
    convert_proto_ip_to_standard_ip_fmt(
        eps_bearer_ctxt.mutable_ue_ip_paa(),
        &sgi_delete_end_point_request.paa.ipv4_address,
        &sgi_delete_end_point_request.paa.ipv6_address, true);

    sgw_handle_sgi_endpoint_deleted(&sgi_delete_end_point_request, imsi64);

    sgw_cm_remove_bearer_context_information(
        s11_pcrf_ded_bearer_deactv_rsp->s_gw_teid_s11_s4, imsi64);
  } else {
    // Remove the dedicated bearer/s context
    for (i = 0; i < no_of_bearers; i++) {
      ebi = s11_pcrf_ded_bearer_deactv_rsp->bearer_contexts.bearer_contexts[i]
                .eps_bearer_id;
      if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                      ebi, &eps_bearer_ctxt) ==
          magma::PROTO_MAP_OK) {
        OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                       "Removed bearer context for (ebi = %u)\n", ebi);
        // Get all the DL flow rules for this dedicated bearer
        for (uint8_t itrn = 0;
             itrn < eps_bearer_ctxt.tft().number_of_packet_filters(); ++itrn) {
          // Prepare DL flow rule from stored packet filters
          struct ip_flow_dl dlflow = {0};
          struct in_addr ue_ipv4 = {.s_addr = 0};
          struct in6_addr ue_ipv6 = {};
          convert_proto_ip_to_standard_ip_fmt(
              eps_bearer_ctxt.mutable_ue_ip_paa(), &ue_ipv4, &ue_ipv6, true);
          bool is_ue_ipv6_pres = false;
          if (((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6 ||
                eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) &&
              !(eps_bearer_ctxt.ue_ip_paa().ipv6_addr().empty())) {
            is_ue_ipv6_pres = true;
          }

          magma::lte::oai::PacketFilter create_new_tft =
              eps_bearer_ctxt.tft().packet_filter_list().create_new_tft(itrn);
          create_new_tft_t new_tft = {0};
          proto_to_packet_filter(create_new_tft, &new_tft);
          if (is_ue_ipv6_pres) {
            generate_dl_flow(&new_tft.packetfiltercontents, ue_ipv4.s_addr,
                             &ue_ipv6, &dlflow);
          } else {
            generate_dl_flow(&new_tft.packetfiltercontents, ue_ipv4.s_addr,
                             nullptr, &dlflow);
          }

          struct in_addr enb = {.s_addr = 0};
          struct in6_addr enb_ipv6 = {};
          convert_proto_ip_to_standard_ip_fmt(
              eps_bearer_ctxt.mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
              spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

#if !MME_UNIT_TEST
          if (gtp_tunnel_ops->del_tunnel(
                  enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
                  eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                  eps_bearer_ctxt.enb_teid_s1u(), &dlflow) < 0) {
            OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                            "ERROR in deleting TUNNEL " TEID_FMT
                            " (eNB) <-> (SGW) " TEID_FMT "\n",
                            eps_bearer_ctxt.enb_teid_s1u(),
                            eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up());
          }
#endif
        }

        sgw_remove_eps_bearer_context(sgw_context_p->mutable_pdn_connection(),
                                      eps_bearer_ctxt.eps_bearer_id());
      }
    }
  }
  // Send DEACTIVATE_DEDICATED_BEARER_RSP to SPGW Service
  if (!s11_pcrf_ded_bearer_deactv_rsp->mme_initiated_local_deact) {
    spgw_handle_nw_init_deactivate_bearer_rsp(
        s11_pcrf_ded_bearer_deactv_rsp->cause, ebi);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

status_code_e sgw_handle_ip_allocation_rsp(
    spgw_state_t* spgw_state,
    const itti_ip_allocation_response_t* ip_allocation_rsp, imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  OAILOG_DEBUG_UE(LOG_SPGW_APP, imsi64,
                  "Received ip_allocation_rsp from gRPC task handler\n");

  magma::lte::oai::S11BearerContext* bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ip_allocation_rsp->context_teid);

  if (!bearer_ctxt_info_p) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                    "Failed to get SPGW UE context for teid" TEID_FMT,
                    ip_allocation_rsp->context_teid);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      bearer_ctxt_info_p->mutable_sgw_eps_bearer_context();
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctx;
  if (sgw_cm_get_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                  ip_allocation_rsp->eps_bearer_id,
                                  &eps_bearer_ctx) != magma::PROTO_MAP_OK) {
    OAILOG_ERROR_UE(LOG_SPGW_APP, sgw_context_p->imsi64(),
                    "Failed to get default bearer context for bearer id: %u\n",
                    ip_allocation_rsp->eps_bearer_id);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  std::string imsi = sgw_context_p->imsi();
  if (ip_allocation_rsp->status == SGI_STATUS_OK) {
    eps_bearer_ctx.mutable_ue_ip_paa()->set_pdn_type(
        ip_allocation_rsp->paa.pdn_type);
    // create session in PCEF
    s5_create_session_request_t session_req = {0};
    session_req.context_teid = ip_allocation_rsp->context_teid;
    session_req.eps_bearer_id = ip_allocation_rsp->eps_bearer_id;
    session_req.status = ip_allocation_rsp->status;
    struct pcef_create_session_data session_data = {};
    get_session_req_data(spgw_state, sgw_context_p->mutable_saved_message(),
                         &session_data);

    // Create session based IPv4, IPv6 or IPv4v6 PDN type
    if (ip_allocation_rsp->paa.pdn_type == IPv4) {
      char ip_str[INET_ADDRSTRLEN];
      inet_ntop(AF_INET, &(ip_allocation_rsp->paa.ipv4_address.s_addr), ip_str,
                INET_ADDRSTRLEN);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv4_addr(ip_str);
      pcef_create_session(imsi, ip_str, NULL, &session_data, session_req);
    } else if (ip_allocation_rsp->paa.pdn_type == IPv6) {
      char ip6_str[INET6_ADDRSTRLEN];
      inet_ntop(AF_INET6, &(ip_allocation_rsp->paa.ipv6_address), ip6_str,
                INET6_ADDRSTRLEN);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv6_addr(ip6_str);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv6_prefix_length(
          ip_allocation_rsp->paa.ipv6_prefix_length);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_vlan(ip_allocation_rsp->paa.vlan);
      pcef_create_session(imsi, NULL, ip6_str, &session_data, session_req);
    } else if (ip_allocation_rsp->paa.pdn_type == IPv4_AND_v6) {
      char ip4_str[INET_ADDRSTRLEN];
      inet_ntop(AF_INET, &(ip_allocation_rsp->paa.ipv4_address.s_addr), ip4_str,
                INET_ADDRSTRLEN);
      char ip6_str[INET6_ADDRSTRLEN];
      inet_ntop(AF_INET6, &(ip_allocation_rsp->paa.ipv6_address), ip6_str,
                INET6_ADDRSTRLEN);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv4_addr(ip4_str);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv6_addr(ip6_str);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_ipv6_prefix_length(
          ip_allocation_rsp->paa.ipv6_prefix_length);
      eps_bearer_ctx.mutable_ue_ip_paa()->set_vlan(ip_allocation_rsp->paa.vlan);
      pcef_create_session(imsi, ip4_str, ip6_str, &session_data, session_req);
    }
    if (sgw_update_eps_bearer_entry(sgw_context_p->mutable_pdn_connection(),
                                    ip_allocation_rsp->eps_bearer_id,
                                    &eps_bearer_ctx) != magma::PROTO_MAP_OK) {
      OAILOG_ERROR_UE(
          LOG_SPGW_APP, sgw_context_p->imsi64(),
          "Failed to update default bearer context for bearer id %u\n",
          ip_allocation_rsp->eps_bearer_id);
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
    }
  } else {
    if (ip_allocation_rsp->status == SGI_STATUS_ERROR_SYSTEM_FAILURE) {
      /*
       * This implies that UE session was not release properly.
       * Release the IP address so that subsequent attempt is successfull
       */
      // TODO - Release the GTP-tunnel corresponding to this IP address
      std::string apn = sgw_context_p->pdn_connection().apn_in_use();
      if (ip_allocation_rsp->paa.pdn_type == IPv4) {
        release_ipv4_address(imsi, apn, &ip_allocation_rsp->paa.ipv4_address);
      } else if (ip_allocation_rsp->paa.pdn_type == IPv6) {
        release_ipv6_address(imsi, apn, &ip_allocation_rsp->paa.ipv6_address);
      } else if (ip_allocation_rsp->paa.pdn_type == IPv4_AND_v6) {
        release_ipv4v6_address(imsi, apn, &ip_allocation_rsp->paa.ipv4_address,
                               &ip_allocation_rsp->paa.ipv6_address);
      }
    }

    // If we are here then the IP address allocation has failed
    s5_create_session_response_t s5_response;
    s5_response.eps_bearer_id = ip_allocation_rsp->eps_bearer_id;
    s5_response.context_teid = ip_allocation_rsp->context_teid;
    s5_response.failure_cause = IP_ALLOCATION_FAILURE;
    handle_s5_create_session_response(spgw_state, bearer_ctxt_info_p,
                                      s5_response);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

bool is_enb_ip_address_same(const fteid_t* fte_p, struct in_addr ipv4,
                            struct in6_addr ipv6) {
  bool rc = true;
  if ((fte_p)->ipv4) {
    if (ipv4.s_addr != (fte_p)->ipv4_address.s_addr) {
      rc = false;
    }
  } else if ((fte_p)->ipv6) {
    if (memcmp(&ipv6, &(fte_p)->ipv6_address, sizeof(struct in6_addr)) != 0) {
      rc = false;
    }
  } else {
    rc = true;
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

void handle_failed_s8_create_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_context,
    sgw_eps_bearer_ctxt_t* dedicated_bearer_ctxt_p, log_proto_t module) {
  OAILOG_FUNC_IN(module);
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = NULL;
  struct sgw_eps_bearer_entry_wrapper_s* sgw_eps_bearer_entry_p = NULL;
  char policy_rule_name[POLICY_RULE_NAME_MAXLEN + 1];
  ebi_t default_bearer_id = 0;

  if (sgw_context_p) {
    default_bearer_id = sgw_context_p->pdn_connection.default_bearer;
    pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(sgw_context_p);
    if (((pgw_ni_cbr_proc) &&
         (!LIST_EMPTY(pgw_ni_cbr_proc->pending_eps_bearers)))) {
      sgw_eps_bearer_entry_p = LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
      while (sgw_eps_bearer_entry_p) {
        if (bearer_context->s1u_sgw_fteid.teid ==
            sgw_eps_bearer_entry_p->sgw_eps_bearer_entry
                ->s_gw_teid_S1u_S12_S4_up) {
          if (module == LOG_SPGW_APP) {
            snprintf(
                policy_rule_name, POLICY_RULE_NAME_MAXLEN + 1, "%s",
                sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->policy_rule_name);
          }
          if (dedicated_bearer_ctxt_p) {
            memcpy(dedicated_bearer_ctxt_p,
                   sgw_eps_bearer_entry_p->sgw_eps_bearer_entry,
                   sizeof(sgw_eps_bearer_ctxt_t));
          }
          // Remove the temporary spgw entry
          LIST_REMOVE(sgw_eps_bearer_entry_p, entries);
          if (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry) {
            free_wrapper((void**)&sgw_eps_bearer_entry_p->sgw_eps_bearer_entry
                             ->pgw_cp_ip_port);
            free_cpp_wrapper(reinterpret_cast<void**>(
                &sgw_eps_bearer_entry_p->sgw_eps_bearer_entry));
          }
          free_cpp_wrapper(reinterpret_cast<void**>(&sgw_eps_bearer_entry_p));
          break;
        }
        sgw_eps_bearer_entry_p = LIST_NEXT(sgw_eps_bearer_entry_p, entries);
      }
      if (pgw_ni_cbr_proc &&
          (LIST_EMPTY(pgw_ni_cbr_proc->pending_eps_bearers))) {
        pgw_base_proc_t* base_proc1 =
            LIST_FIRST(sgw_context_p->pending_procedures);
        LIST_REMOVE(base_proc1, entries);
        free_cpp_wrapper(
            reinterpret_cast<void**>(&sgw_context_p->pending_procedures));
        free_cpp_wrapper(
            reinterpret_cast<void**>(&pgw_ni_cbr_proc->pending_eps_bearers));
        pgw_free_procedure_create_bearer((pgw_ni_cbr_proc_t**)&pgw_ni_cbr_proc);
      }
    }
  }
  if (module == LOG_SPGW_APP) {
    int rc = spgw_send_nw_init_activate_bearer_rsp(
        cause, imsi64, bearer_context, default_bearer_id, policy_rule_name);
    if (rc != RETURNok) {
      OAILOG_ERROR_UE(module, imsi64,
                      "Failed to send ACTIVATE_DEDICATED_BEARER_RSP to PCRF\n");
    }
  }
  OAILOG_FUNC_OUT(module);
}

void handle_failed_create_bearer_response(
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p,
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_context,
    magma::lte::oai::SgwEpsBearerContext* dedicated_bearer_ctxt_p,
    log_proto_t module) {
  OAILOG_FUNC_IN(module);
  char policy_rule_name[POLICY_RULE_NAME_MAXLEN + 1];
  ebi_t default_bearer_id = 0;

  if (sgw_context_p) {
    if (!(sgw_context_p->pending_procedures_size())) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Failed to get create bearer procedure from temporary "
                      "stored context for EBI %u\n",
                      bearer_context->eps_bearer_id);
      OAILOG_FUNC_OUT(module);
    }
    default_bearer_id = sgw_context_p->pdn_connection().default_bearer();

    uint8_t num_of_bearers_deleted = 0;
    uint8_t num_of_pending_procedures =
        sgw_context_p->pending_procedures_size();
    for (uint8_t proc_index = 0;
         proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
      magma::lte::oai::PgwCbrProcedure* pgw_ni_cbr_proc =
          sgw_context_p->mutable_pending_procedures(proc_index);
      if (!pgw_ni_cbr_proc) {
        OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                        "Failed to get create bearer procedure from temporary "
                        "stored context for EBI %u\n",
                        bearer_context->eps_bearer_id);
        OAILOG_FUNC_OUT(module);
      }
      if (pgw_ni_cbr_proc->type() ==
          PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST) {
        num_of_bearers_deleted = pgw_ni_cbr_proc->pending_eps_bearers_size();
        for (uint8_t bearer_index = 0;
             bearer_index < pgw_ni_cbr_proc->pending_eps_bearers_size();
             bearer_index++) {
          magma::lte::oai::SgwEpsBearerContext* stored_bearer_context_info =
              pgw_ni_cbr_proc->mutable_pending_eps_bearers(bearer_index);
          if (stored_bearer_context_info) {
            if (bearer_context->s1u_sgw_fteid.teid ==
                stored_bearer_context_info->sgw_teid_s1u_s12_s4_up()) {
              if (module == LOG_SPGW_APP) {
                snprintf(
                    policy_rule_name, POLICY_RULE_NAME_MAXLEN + 1, "%s",
                    stored_bearer_context_info->policy_rule_name().c_str());
              }
              if (dedicated_bearer_ctxt_p) {
                memcpy(dedicated_bearer_ctxt_p, stored_bearer_context_info,
                       sizeof(dedicated_bearer_ctxt_p));
              }
            }
            --num_of_bearers_deleted;
          }
          if (num_of_bearers_deleted == 0) {
            pgw_ni_cbr_proc->clear_pending_eps_bearers();
            --num_of_pending_procedures;
          }
        }
      }
    }
    if (num_of_pending_procedures == 0) {
      sgw_context_p->clear_pending_procedures();
    }
    if (module == LOG_SPGW_APP) {
      if (spgw_send_nw_init_activate_bearer_rsp(cause, imsi64, bearer_context,
                                                default_bearer_id,
                                                policy_rule_name) != RETURNok) {
        OAILOG_ERROR_UE(
            module, imsi64,
            "Failed to send ACTIVATE_DEDICATED_BEARER_RSP to PCRF\n");
      }
    }
  }
  OAILOG_FUNC_OUT(module);
}

// Fills up downlink (DL) flow match rule from packet filters of eps bearer
void generate_dl_flow(packet_filter_contents_t* packet_filter,
                      in_addr_t ipv4_s_addr, struct in6_addr* ue_ipv6,
                      struct ip_flow_dl* dlflow) {
  // Prepare DL flow rule
  // The TFTs are DL TFTs: UE is the destination/local,
  // PDN end point is the source/remote.

  // Adding UE to the rule is safe
  if (ipv4_s_addr && ue_ipv6) {
    /* In case of ipv4v6 since there is no other way to know if ipv4 or ipv6
     * address should be set, check the remote address flag and set the
     * ips accordingly
     */
    if ((TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG & packet_filter->flags) ==
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG) {
      dlflow->src_ip.s_addr =
          (((uint32_t)packet_filter->ipv4remoteaddr[0].addr) << 24) +
          (((uint32_t)packet_filter->ipv4remoteaddr[1].addr) << 16) +
          (((uint32_t)packet_filter->ipv4remoteaddr[2].addr) << 8) +
          (((uint32_t)packet_filter->ipv4remoteaddr[3].addr));

      dlflow->set_params |= SRC_IPV4;
      dlflow->dst_ip.s_addr = ipv4_s_addr;
      dlflow->set_params |= DST_IPV4;
    } else if ((TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG &
                packet_filter->flags) ==
               TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG) {
      struct in6_addr remoteaddr = {};
      for (uint8_t itr = 0; itr < 16; itr++) {
        remoteaddr.s6_addr[itr] = packet_filter->ipv6remoteaddr[itr].addr;
      }
      dlflow->src_ip6 = remoteaddr;
      dlflow->set_params |= SRC_IPV6;
      dlflow->dst_ip6 = *ue_ipv6;
      dlflow->set_params |= DST_IPV6;
    }
  } else if (ipv4_s_addr) {
    dlflow->dst_ip.s_addr = ipv4_s_addr;
    dlflow->set_params |= DST_IPV4;
    if ((TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG & packet_filter->flags) ==
        TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG) {
      dlflow->src_ip.s_addr =
          (((uint32_t)packet_filter->ipv4remoteaddr[0].addr) << 24) +
          (((uint32_t)packet_filter->ipv4remoteaddr[1].addr) << 16) +
          (((uint32_t)packet_filter->ipv4remoteaddr[2].addr) << 8) +
          (((uint32_t)packet_filter->ipv4remoteaddr[3].addr));
      dlflow->set_params |= SRC_IPV4;
    }
  } else if (ue_ipv6) {
    dlflow->dst_ip6 = *ue_ipv6;
    dlflow->set_params |= DST_IPV6;
    if ((TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG & packet_filter->flags) ==
        TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG) {
      struct in6_addr remoteaddr = {};
      for (uint8_t itr = 0; itr < 16; itr++) {
        remoteaddr.s6_addr[itr] = packet_filter->ipv6remoteaddr[itr].addr;
      }
      dlflow->src_ip6 = remoteaddr;
      dlflow->set_params |= SRC_IPV6;
    }
  }
  // Specify the next header
  dlflow->ip_proto = packet_filter->protocolidentifier_nextheader;
  // Match on proto if it is explicity specified to be
  // other than the dummy IP. When PCRF RAR message does not
  // define the protocol type, this field defaults to value 0.
  // OVS would still apply exact match on 0  if parameter is set,
  // although incoming packets will have a proper protocol number
  // in its header leading to no match.
  if (dlflow->ip_proto != IPPROTO_IP) {
    dlflow->set_params |= IP_PROTO;
  }

  // Process remote port if present
  if ((TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG & packet_filter->flags) ==
      TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG) {
    if (dlflow->ip_proto == IPPROTO_TCP) {
      dlflow->set_params |= TCP_SRC_PORT;
      dlflow->tcp_src_port = packet_filter->singleremoteport;
    } else if (dlflow->ip_proto == IPPROTO_UDP) {
      dlflow->set_params |= UDP_SRC_PORT;
      dlflow->udp_src_port = packet_filter->singleremoteport;
    }
  }

  // Process UE port if present
  if ((TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG & packet_filter->flags) ==
      TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG) {
    if (dlflow->ip_proto == IPPROTO_TCP) {
      dlflow->set_params |= TCP_DST_PORT;
      dlflow->tcp_dst_port = packet_filter->singlelocalport;
    } else if (dlflow->ip_proto == IPPROTO_UDP) {
      dlflow->set_params |= UDP_DST_PORT;
      dlflow->udp_dst_port = packet_filter->singlelocalport;
    }
  }
}

// Helper function to generate dl flows and add tunnel for ipv4/ipv6/ipv4v6
// bearers
static void add_tunnel_helper(
    magma::lte::oai::S11BearerContext* spgw_context,
    magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt_entry_p,
    imsi64_t imsi64) {
  uint32_t rc = RETURNerror;
  struct in_addr enb = {.s_addr = 0};
  const char* apn = reinterpret_cast<const char*>(
      spgw_context->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection()
          ->apn_in_use()
          .c_str());
  struct in6_addr enb_ipv6 = {};

  convert_proto_ip_to_standard_ip_fmt(
      eps_bearer_ctxt_entry_p->mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6,
      spgw_config.sgw_config.ipv6.s1_ipv6_enabled);

  struct in_addr ue_ipv4 = {.s_addr = 0};
  struct in6_addr ue_ipv6 = {};
  convert_proto_ip_to_standard_ip_fmt(
      eps_bearer_ctxt_entry_p->mutable_ue_ip_paa(), &ue_ipv4, &ue_ipv6, true);

  bool is_ue_ipv6_pres = false;
  if (((eps_bearer_ctxt_entry_p->ue_ip_paa().pdn_type() == IPv6 ||
        eps_bearer_ctxt_entry_p->ue_ip_paa().pdn_type() == IPv4_AND_v6)) &&
      !(eps_bearer_ctxt_entry_p->ue_ip_paa().ipv6_addr().empty())) {
    is_ue_ipv6_pres = true;
  }
  int vlan = eps_bearer_ctxt_entry_p->ue_ip_paa().vlan();
  Imsi_t imsi;
  memcpy(imsi.digit,
         spgw_context->mutable_sgw_eps_bearer_context()->imsi().c_str(),
         spgw_context->mutable_sgw_eps_bearer_context()->imsi().size());
  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64, "Number of packet filter rules: %d\n",
                 eps_bearer_ctxt_entry_p->tft().number_of_packet_filters());
  for (uint8_t i = 0;
       i < eps_bearer_ctxt_entry_p->tft().number_of_packet_filters(); ++i) {
    struct ip_flow_dl dlflow = {0};
    magma::lte::oai::PacketFilter create_new_tft =
        eps_bearer_ctxt_entry_p->tft().packet_filter_list().create_new_tft(i);
    create_new_tft_t new_tft = {0};
    proto_to_packet_filter(create_new_tft, &new_tft);
    if (is_ue_ipv6_pres) {
      generate_dl_flow(&new_tft.packetfiltercontents, ue_ipv4.s_addr, &ue_ipv6,
                       &dlflow);
    } else {
      generate_dl_flow(&new_tft.packetfiltercontents, ue_ipv4.s_addr, nullptr,
                       &dlflow);
    }

#if !MME_UNIT_TEST
    rc = gtpv1u_add_tunnel(ue_ipv4, &ue_ipv6, vlan, enb, &enb_ipv6,
                           eps_bearer_ctxt_entry_p->sgw_teid_s1u_s12_s4_up(),
                           eps_bearer_ctxt_entry_p->enb_teid_s1u(), imsi,
                           &dlflow,
                           eps_bearer_ctxt_entry_p->tft()
                               .packet_filter_list()
                               .create_new_tft(i)
                               .eval_precedence(),
                           apn);

    if (rc != RETURNok) {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "ERROR in setting up TUNNEL err=%d\n", rc);
    } else {
      OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                     "Successfully setup flow rule for EPS bearer id %u "
                     "tunnel " TEID_FMT " (eNB) <-> (SGW) " TEID_FMT "\n",
                     eps_bearer_ctxt_entry_p->eps_bearer_id(),
                     eps_bearer_ctxt_entry_p->enb_teid_s1u(),
                     eps_bearer_ctxt_entry_p->sgw_teid_s1u_s12_s4_up());
    }
#endif
  }
}

bool does_sgw_bearer_context_hold_valid_enb_ip(
    ip_address_t enb_ip_address_S1u) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  static struct in6_addr ipv6_address = {};
  switch (enb_ip_address_S1u.pdn_type) {
    case IPv4:
      if (enb_ip_address_S1u.address.ipv4_address.s_addr) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    case IPv4_AND_v6:
      if ((enb_ip_address_S1u.address.ipv4_address.s_addr) ||
          (memcmp(&ipv6_address, &(enb_ip_address_S1u.address.ipv6_address),
                  sizeof(struct in6_addr)) != 0)) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    case IPv6:
      if (memcmp(&ipv6_address, &(enb_ip_address_S1u.address.ipv6_address),
                 sizeof(struct in6_addr)) != 0) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    default:
      OAILOG_ERROR(LOG_SPGW_APP, "Invalid pdn-type \n");
      break;
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, false);
}

bool does_bearer_context_hold_valid_enb_ip(
    magma::lte::oai::IpTupple enb_ip_address_S1u) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  switch (enb_ip_address_S1u.pdn_type()) {
    case IPv4:
      if (enb_ip_address_S1u.ipv4_addr().size()) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    case IPv4_AND_v6:
      if ((enb_ip_address_S1u.ipv4_addr().size()) ||
          ((enb_ip_address_S1u.ipv6_addr().size()))) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    case IPv6:
      if (enb_ip_address_S1u.ipv6_addr().size()) {
        OAILOG_FUNC_RETURN(LOG_SPGW_APP, true);
      }
      break;
    default:
      OAILOG_ERROR(LOG_SPGW_APP, "Invalid pdn-type \n");
      break;
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, false);
}

void sgw_send_release_access_bearer_response(
    log_proto_t module, task_id_t origin_task_id, imsi64_t imsi64,
    gtpv2c_cause_value_t cause,
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    teid_t mme_teid_s11) {
  OAILOG_FUNC_IN(module);
  status_code_e rv = RETURNok;
  itti_s11_release_access_bearers_response_t* release_access_bearers_resp_p =
      NULL;
  MessageDef* message_p = itti_alloc_new_message(
      origin_task_id, S11_RELEASE_ACCESS_BEARERS_RESPONSE);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        module, imsi64,
        "Failed to allocate memory for Release Access Bearer Response \n");
    OAILOG_FUNC_OUT(module);
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  release_access_bearers_resp_p =
      &message_p->ittiMsg.s11_release_access_bearers_response;
  release_access_bearers_resp_p->cause.cause_value = cause;
  release_access_bearers_resp_p->teid = mme_teid_s11;
  release_access_bearers_resp_p->trxn = release_access_bearers_req_pP->trxn;

  if (module == LOG_SPGW_APP) {
    rv = send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME, message_p);
  } else if (module == LOG_SGW_S8) {
    rv = send_msg_to_task(&sgw_s8_task_zmq_ctx, TASK_MME, message_p);
  } else {
    OAILOG_ERROR_UE(module, imsi64, "Invalid module \n");
  }
  if (rv != RETURNok) {
    OAILOG_ERROR_UE(module, imsi64,
                    "Failed to send Release Access Bearer Response to MME\n");
    OAILOG_FUNC_OUT(module);
  }
  OAILOG_DEBUG_UE(module, imsi64,
                  "Release Access Bearer Response sent to MME\n");
  OAILOG_FUNC_OUT(module);
}

void sgw_process_release_access_bearer_request(
    log_proto_t module, imsi64_t imsi64,
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context) {
  OAILOG_FUNC_IN(module);
  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  /*
   * Release the tunnels so that in idle state, DL packets are not sent
   * towards eNB.
   * These tunnels will be added again when UE moves back to connected mode.
   */
  magma::lte::oai::SgwPdnConnection* pdn_connection_p =
      sgw_context->mutable_pdn_connection();
  if (pdn_connection_p->eps_bearer_map_size()) {
    map_uint32_spgw_eps_bearer_context_t eps_bearer_map;
    eps_bearer_map.map = pdn_connection_p->mutable_eps_bearer_map();
    for (auto itr = eps_bearer_map.map->begin();
         itr != eps_bearer_map.map->end(); itr++) {
      magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt = itr->second;

      struct in_addr enb = {.s_addr = 0};
      struct in6_addr enb_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(
          eps_bearer_ctxt.mutable_enb_s1u_ip_addr(), &enb, &enb_ipv6, true);

      struct in_addr ue_ipv4 = {.s_addr = 0};
      struct in6_addr ue_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(eps_bearer_ctxt.mutable_ue_ip_paa(),
                                          &ue_ipv4, &ue_ipv6, true);
      bool is_ue_ipv6_pres = false;
      if (((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6 ||
            eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) &&
          !(eps_bearer_ctxt.ue_ip_paa().ipv6_addr().empty())) {
        is_ue_ipv6_pres = true;
      }

      struct in_addr pgw = {.s_addr = 0};
      struct in6_addr pgw_ipv6 = {};
      convert_proto_ip_to_standard_ip_fmt(
          eps_bearer_ctxt.mutable_pgw_address_up(), &pgw, &pgw_ipv6, true);

      OAILOG_DEBUG_UE(module, imsi64,
                      "Deleting tunnel for bearer_id %u ue addr %x enb_ip %x "
                      "sgw_teid_s1u_s12_s4_up %x, enb_teid_S1u %x pgw_up_ip %x "
                      "pgw_up_teid " TEID_FMT
                      "s_gw_ip_address_S5_S8_up %s"
                      "s_gw_teid_S5_S8_up " TEID_FMT,
                      eps_bearer_ctxt.eps_bearer_id(), ue_ipv4.s_addr,
                      enb.s_addr, eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                      eps_bearer_ctxt.enb_teid_s1u(), pgw.s_addr,
                      eps_bearer_ctxt.pgw_teid_s5_s8_up(),
                      eps_bearer_ctxt.sgw_ip_address_s5_s8_up().c_str(),
                      eps_bearer_ctxt.sgw_teid_s5_s8_up());
#if !MME_UNIT_TEST  // skip tunnel deletion for unit tests
      int rv = RETURNok;
      if (module == LOG_SPGW_APP) {
        if (is_ue_ipv6_pres) {
          rv = gtp_tunnel_ops->del_tunnel(
              enb, &enb_ipv6, ue_ipv4, &ue_ipv6,
              eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
              eps_bearer_ctxt.enb_teid_s1u(), nullptr);
        } else {
          rv = gtp_tunnel_ops->del_tunnel(
              enb, &enb_ipv6, ue_ipv4, nullptr,
              eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
              eps_bearer_ctxt.enb_teid_s1u(), nullptr);
        }
      } else if (module == LOG_SGW_S8) {
        rv = gtpv1u_del_s8_tunnel(enb, &enb_ipv6, pgw, &pgw_ipv6, ue_ipv4,
                                  &ue_ipv6,
                                  eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up(),
                                  eps_bearer_ctxt.sgw_teid_s5_s8_up());
      }

      // TODO Need to add handling on failing to delete s1-u tunnel rules from
      // ovs flow table
      if (rv < 0) {
        OAILOG_ERROR_UE(module, imsi64,
                        "ERROR in deleting TUNNEL " TEID_FMT
                        " (eNB) <-> (SGW) " TEID_FMT "\n",
                        eps_bearer_ctxt.enb_teid_s1u(),
                        eps_bearer_ctxt.sgw_teid_s1u_s12_s4_up());
      }
      // Paging is performed without packet buffering
      Imsi_t imsi;
      memcpy(imsi.digit, sgw_context->imsi().c_str(),
             sgw_context->imsi().size());
      if (is_ue_ipv6_pres) {
        rv = gtp_tunnel_ops->add_paging_rule(imsi, ue_ipv4, &ue_ipv6);
      } else {
        rv = gtp_tunnel_ops->add_paging_rule(imsi, ue_ipv4, nullptr);
      }
      // Convert to string for logging
      char* ip_str = inet_ntoa(ue_ipv4);
      if (rv < 0) {
        OAILOG_ERROR_UE(module, imsi64,
                        "ERROR in setting paging rule for IP Addr: %s\n",
                        ip_str);
      }
      if ((eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv6) ||
          (eps_bearer_ctxt.ue_ip_paa().pdn_type() == IPv4_AND_v6)) {
        OAILOG_DEBUG(module, "Set the paging rule for IPv6 Addr: %s\n",
                     eps_bearer_ctxt.ue_ip_paa().ipv6_addr().c_str());
      } else {
        OAILOG_DEBUG(module, "Set the paging rule for IP Addr: %s\n", ip_str);
      }
#endif
      // release enb's UP teid and IP address
      eps_bearer_ctxt.clear_enb_s1u_ip_addr();
      eps_bearer_ctxt.set_enb_teid_s1u(INVALID_TEID);
      eps_bearer_map.update_val(eps_bearer_ctxt.eps_bearer_id(),
                                &eps_bearer_ctxt);
    }
  }
  OAILOG_FUNC_OUT(module);
}

// Generates random s11 control plane teid
static teid_t sgw_generate_new_s11_cp_teid(void) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  magma::lte::oai::S11BearerContext* s_plus_p_gw_eps_bearer_ctxt_info_p =
      nullptr;
  teid_t teid = INVALID_TEID;
  // note srand with seed is initialized at main
  do {
    teid = (teid_t)rand();
    s_plus_p_gw_eps_bearer_ctxt_info_p = sgw_cm_get_spgw_context(teid);
  } while (s_plus_p_gw_eps_bearer_ctxt_info_p);

  OAILOG_FUNC_RETURN(LOG_SGW_S8, teid);
}

// Handles delete bearer cmd from MME
void sgw_handle_delete_bearer_cmd(
    itti_s11_delete_bearer_command_t* s11_delete_bearer_command,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  OAILOG_INFO_UE(LOG_SPGW_APP, imsi64,
                 "Received s11_delete_bearer_command for teid" TEID_FMT "\n",
                 s11_delete_bearer_command->teid);
  magma::lte::oai::S11BearerContext* spgw_ctxt =
      sgw_cm_get_spgw_context(s11_delete_bearer_command->teid);
  if (!spgw_ctxt) {
    OAILOG_ERROR_UE(
        LOG_SPGW_APP, imsi64,
        "failed to get s_plus_pgw_eps bearer context for teid " TEID_FMT
        " while processing s11_delete_bearer_command\n",
        s11_delete_bearer_command->teid);
    OAILOG_FUNC_OUT(TASK_SPGW_APP);
  }

  const char* imsi = reinterpret_cast<const char*>(
      spgw_ctxt->sgw_eps_bearer_context().imsi().c_str());
  pcef_delete_dedicated_bearer(imsi, s11_delete_bearer_command->ebi_list);

  // Send itti_s11_nw_init_deactv_bearer_request_t to MME to delete the bearer/s
  spgw_build_and_send_s11_deactivate_bearer_req(
      imsi64, s11_delete_bearer_command->ebi_list.num_ebi,
      s11_delete_bearer_command->ebi_list.ebis, false,
      s11_delete_bearer_command->local_teid, LOG_SPGW_APP);

  OAILOG_FUNC_OUT(TASK_SPGW_APP);
}

void convert_proto_ip_to_standard_ip_fmt(magma::lte::oai::IpTupple* proto_ip,
                                         struct in_addr* ipv4,
                                         struct in6_addr* ipv6,
                                         bool ipv6_enabled) {
  inet_pton(AF_INET, proto_ip->ipv4_addr().c_str(), &ipv4->s_addr);
  if (ipv6_enabled &&
      (proto_ip->pdn_type() == IPv6 || proto_ip->pdn_type() == IPv4_AND_v6)) {
    inet_pton(AF_INET6, proto_ip->ipv6_addr().c_str(), ipv6);
  }
}

static void get_session_req_data(
    spgw_state_t* spgw_state,
    const magma::lte::oai::CreateSessionMessage* saved_req,
    struct pcef_create_session_data* data) {
  data->msisdn_len = saved_req->msisdn().size();
  memcpy(data->msisdn, saved_req->msisdn().c_str(), data->msisdn_len);
  data->imeisv_exists = saved_req->mei().size() > 0 ? true : false;
  if (data->imeisv_exists) {
    memcpy(data->imeisv, saved_req->mei().c_str(), saved_req->mei().size());
  }

  data->uli_exists = saved_req->uli().size() > 0 ? true : false;
  if (data->uli_exists) {
    memcpy(data->uli, saved_req->uli().c_str(), saved_req->uli().size());
  }

  data->mcc_mnc_len = saved_req->serving_network().mcc().size();
  memset(data->mcc_mnc, 0, 7);
  memcpy(data->mcc_mnc, saved_req->serving_network().mcc().c_str(),
         saved_req->serving_network().mcc().size());

  memcpy(&data->mcc_mnc[data->mcc_mnc_len],
         saved_req->serving_network().mnc().c_str(),
         saved_req->serving_network().mnc().size());
  data->mcc_mnc_len += saved_req->serving_network().mnc().size();
  get_imsi_plmn_from_session_req(saved_req->imsi(), data);

  memcpy(data->charging_characteristics.value,
         saved_req->charging_characteristics().c_str(),
         saved_req->charging_characteristics().size());
  data->charging_characteristics.length =
      saved_req->charging_characteristics().size();

  memset(data->apn, 0, APN_MAX_LENGTH + 1);
  memcpy(data->apn, saved_req->apn().c_str(), saved_req->apn().size());
  data->pdn_type = saved_req->pdn_type();

  inet_ntop(AF_INET, &spgw_state->sgw_ip_address_S1u_S12_S4_up, data->sgw_ip,
            INET_ADDRSTRLEN);

  // QoS Info
  data->ambr_dl = saved_req->ambr().br_dl();
  data->ambr_ul = saved_req->ambr().br_ul();
  auto qos = saved_req->bearer_contexts_to_be_created(0).bearer_level_qos();
  data->pl = qos.pl();
  data->pci = qos.pci();
  data->pvi = qos.pvi();
  data->qci = qos.qci();
}
