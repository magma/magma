/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file pgw_handlers.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define S5_HANDLERS_C

#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "log.h"
#include "spgw_config.h"
#include "pgw_pco.h"
#include "dynamic_memory_check.h"
#include "MobilityClientAPI.h"
#include "pgw_handlers.h"
#include "sgw_handlers.h"
#include "pcef_handlers.h"
#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_23.401.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "common_types.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "ip_forward_messages_types.h"
#include "itti_types.h"
#include "pgw_config.h"
#include "s11_messages_types.h"
#include "service303.h"
#include "sgw_context_manager.h"
#include "sgw_ie_defs.h"
#include "pgw_procedures.h"
#include "spgw_types.h"

extern spgw_config_t spgw_config;
extern void print_bearer_ids_helper(const ebi_t*, uint32_t);

static int _spgw_build_and_send_s11_create_bearer_request(
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p,
  const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
  spgw_state_t* spgw_state,
  teid_t s1_u_sgw_fteid);

static int _create_temporary_dedicated_bearer_context(
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p,
  const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
  spgw_state_t* spgw_state,
  teid_t s1_u_sgw_fteid);

static void _delete_temporary_dedicated_bearer_context(
  teid_t s1_u_sgw_fteid,
  ebi_t lbi,
  s_plus_p_gw_eps_bearer_context_information_t* spgw_context_p);

static int32_t _spgw_build_and_send_s11_deactivate_bearer_req(
  imsi64_t imsi64,
  uint8_t no_of_bearers_to_be_deact,
  ebi_t* ebi_to_be_deactivated,
  bool delete_default_bearer,
  teid_t mme_teid_S11);

static void _spgw_handle_s5_response_with_error(
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  itti_sgi_create_end_point_response_t* sgi_create_endpoint_resp,
  s5_create_session_response_t* s5_response);

//--------------------------------------------------------------------------------

void handle_s5_create_session_request(
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t *new_bearer_ctxt_info_p,
  teid_t context_teid,
  ebi_t eps_bearer_id)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp = {0};
  s5_create_session_response_t s5_response = {0};
  struct in_addr inaddr;
  char *imsi = NULL;
  char *apn = NULL;

  OAILOG_DEBUG(
    LOG_SPGW_APP,
    "Handle s5_create_session_request, for context sgw s11 teid, " TEID_FMT
    "EPS bearer id %u\n",
    context_teid,
    eps_bearer_id);

  if (!new_bearer_ctxt_info_p) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to fetch sgw bearer context from the received context "
      "teid" TEID_FMT "\n",
      context_teid);
    sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_CONTEXT_NOT_FOUND;
    _spgw_handle_s5_response_with_error(
      spgw_state,
      new_bearer_ctxt_info_p,
      context_teid,
      eps_bearer_id,
      &sgi_create_endpoint_resp,
      &s5_response);
  }

  // PCO processing
  protocol_configuration_options_t* pco_req =
    &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
       .pco;
  protocol_configuration_options_t pco_resp = {0};
  protocol_configuration_options_ids_t pco_ids;
  memset(&pco_ids, 0, sizeof pco_ids);

  if (pgw_process_pco_request(pco_req, &pco_resp, &pco_ids) != RETURNok) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Error in processing PCO in create session request for "
      "context_id: " TEID_FMT "\n",
      context_teid);
    sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_FAILED_TO_PROCESS_PCO;
    _spgw_handle_s5_response_with_error(
      spgw_state,
      new_bearer_ctxt_info_p,
      context_teid,
      eps_bearer_id,
      &sgi_create_endpoint_resp,
      &s5_response);
  }
  copy_protocol_configuration_options(&sgi_create_endpoint_resp.pco, &pco_resp);
  clear_protocol_configuration_options(&pco_resp);

  // IP forward will forward packets to this teid
  sgi_create_endpoint_resp.context_teid = context_teid;
  sgi_create_endpoint_resp.eps_bearer_id = eps_bearer_id;
  sgi_create_endpoint_resp.paa.pdn_type =
    new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message
      .pdn_type;

  imsi =
    (char*)
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.imsi.digit;

  apn = (char*) new_bearer_ctxt_info_p->sgw_eps_bearer_context_information
          .pdn_connection.apn_in_use;

  switch (sgi_create_endpoint_resp.paa.pdn_type) {
    case IPv4:
      // Use NAS by default if no preference is set.
      //
      // For context, the protocol configuration options (PCO) section of
      // packet from the UE is optional, which means that it is perfectly
      // valid UE to send no PCO preferences at all. The previous logic only
      // allocates an IPv4 address if the UE has explicitly set the PCO
      // parameter for allocating IPv4 via NAS signaling (as opposed to via
      // DHCPv4). This means that, in the absence of either parameter being,
      // set the does not know what to do, so we need a default option as well.
      //
      // Since we only support the NAS signaling option right now, we will
      // default to using NAS signaling UNLESS we see a preference for DHCPv4.
      // This means that all IPv4 addresses are now allocated via NAS signaling
      // unless specified otherwise.
      //
      // In the long run, we will want to evolve the logic to use whatever
      // information we have to choose the ``best" allocation method. This means
      // adding new bitfields to pco_ids in pgw_pco.h, setting them in pgw_pco.c
      // and using them here in conditional logic. We will also want to
      // implement different logic between the PDN types.
      if (!pco_ids.ci_ipv4_address_allocation_via_dhcpv4) {
        pgw_handle_allocate_ipv4_address(
          imsi,
          apn,
          &inaddr,
          sgi_create_endpoint_resp,
          "ipv4",
          context_teid,
          eps_bearer_id,
          spgw_state,
          new_bearer_ctxt_info_p,
          s5_response);
      }
      break;

    case IPv6:
      increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", "ipv6", "result", "failure");
      OAILOG_ERROR(LOG_SPGW_APP, "IPV6 PDN type NOT Supported\n");
      sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_SERVICE_NOT_SUPPORTED;
      break;

    case IPv4_AND_v6:
      pgw_handle_allocate_ipv4_address(
        imsi,
        apn,
        &inaddr,
        sgi_create_endpoint_resp,
        "ipv4v6",
        context_teid,
        eps_bearer_id,
        spgw_state,
        new_bearer_ctxt_info_p,
        s5_response);
      break;

    default:
      AssertFatal(
        0, "BAD paa.pdn_type %d", sgi_create_endpoint_resp.paa.pdn_type);
      break;
  }
}

void _spgw_handle_s5_response_with_error(
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  itti_sgi_create_end_point_response_t* sgi_create_endpoint_resp,
  s5_create_session_response_t* s5_response)
{
  s5_response->context_teid = context_teid;
  s5_response->eps_bearer_id = eps_bearer_id;
  s5_response->sgi_create_endpoint_resp = (*sgi_create_endpoint_resp);
  s5_response->failure_cause = S5_OK;

  OAILOG_DEBUG(
    LOG_SPGW_APP,
    "Sending S5 Create Session Response to SGW: with context teid, " TEID_FMT
    "EPS Bearer Id = %u\n",
    s5_response->context_teid,
    s5_response->eps_bearer_id);
  handle_s5_create_session_response(
    spgw_state, new_bearer_ctxt_info_p, (*s5_response));
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

/*
 * Handle NW initiated Dedicated Bearer Activation from SPGW service
 */
int spgw_handle_nw_initiated_bearer_actv_req(
  spgw_state_t* spgw_state,
  const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
  imsi64_t imsi64,
  gtpv2c_cause_value_t* failed_cause)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint32_t i = 0;
  int rc = RETURNok;
  hash_table_ts_t* hashtblP = NULL;
  uint32_t num_elements = 0;
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p = NULL;
  hash_node_t* node = NULL;
  bool is_imsi_found = false;
  bool is_lbi_found = false;

  OAILOG_INFO(
    LOG_SPGW_APP,
    "Received Create Bearer Req from PCRF with lbi:%d IMSI\n" IMSI_64_FMT,
    bearer_req_p->lbi,
    imsi64);

  // TODO: Revisit this if UE context struct manages multiple PDN connections
  hashtblP = get_spgw_ue_state();
  if (!hashtblP) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "No s11_bearer_context_information hash table found \n");
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  /* On reception of Dedicated Bearer Activation Request from PCRF,
   * SPGW shall identify whether valid PDN session exists for the UE
   * using IMSI and LBI, for which Dedicated Bearer Activation is requested.
   */
  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size) &&
         (!is_lbi_found)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    while (node) {
      num_elements++;
      hashtable_ts_get(
        hashtblP, (const hash_key_t) node->key, (void **) &spgw_ctxt_p);
      if (spgw_ctxt_p != NULL) {
        if (!strncmp(
          (const char *)
            spgw_ctxt_p->sgw_eps_bearer_context_information.imsi.digit,
          (const char *) bearer_req_p->imsi,
          strlen((const char *) bearer_req_p->imsi))) {
          is_imsi_found = true;
          if (
            spgw_ctxt_p->sgw_eps_bearer_context_information.pdn_connection
              .default_bearer == bearer_req_p->lbi) {
            is_lbi_found = true;
            break;
          }
        }
      }
      node = node->next;
    }
    i++;
  }

  if ((!is_imsi_found) || (!is_lbi_found)) {
    OAILOG_INFO(
      LOG_SPGW_APP,
      "is_imsi_found (%d), is_lbi_found (%d)\n",
      is_imsi_found, is_lbi_found);
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Sending dedicated_bearer_actv_rsp with REQUEST_REJECTED cause to NW\n");
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(spgw_state);
  // Create temporary dedicated bearer context
  rc = _create_temporary_dedicated_bearer_context(
    spgw_ctxt_p, bearer_req_p, spgw_state, s1_u_sgw_fteid);
  if (rc != RETURNok) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to create temporary dedicated bearer context for lbi: %u \n ",
      bearer_req_p->lbi);
    *failed_cause = REQUEST_REJECTED;
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  // Build and send ITTI message, s11_create_bearer_request to MME APP
  rc = _spgw_build_and_send_s11_create_bearer_request(
    spgw_ctxt_p, bearer_req_p, spgw_state, s1_u_sgw_fteid);
  if (rc != RETURNok) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to build and send S11 Create Bearer Request for lbi :%u \n",
      bearer_req_p->lbi);

    *failed_cause = REQUEST_REJECTED;
    _delete_temporary_dedicated_bearer_context(
      s1_u_sgw_fteid, bearer_req_p->lbi, spgw_ctxt_p);
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

//------------------------------------------------------------------------------
int32_t spgw_handle_nw_initiated_bearer_deactv_req(
  spgw_state_t* spgw_state,
  const itti_gx_nw_init_deactv_bearer_request_t* const bearer_req_p,
  imsi64_t imsi64)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  int32_t rc = RETURNok;
  hash_table_ts_t* hashtblP = NULL;
  uint32_t num_elements = 0;
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p = NULL;
  hash_node_t *node = NULL;
  bool is_lbi_found = false;
  bool is_imsi_found = false;
  bool is_ebi_found = false;
  ebi_t ebi_to_be_deactivated[BEARERS_PER_UE] = {0};
  uint32_t no_of_bearers_to_be_deact = 0;
  uint32_t no_of_bearers_rej = 0;
  ebi_t invalid_bearer_id[BEARERS_PER_UE] = {0};

  OAILOG_INFO(
    LOG_SPGW_APP,
    "Received nw_initiated_deactv_bearer_req from SPGW service \n");
  print_bearer_ids_helper(bearer_req_p->ebi, bearer_req_p->no_of_bearers);

  hashtblP = get_spgw_ue_state();
  if (hashtblP == NULL) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "No s11_bearer_context_information hash table is found\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  // Check if valid LBI and EBI recvd
  /* For multi PDN, same IMSI can have multiple sessions, which means there
   * will be multiple entries for different sessions with the same IMSI. Hence
   * even though IMSI is found search the entire list for the LBI
   */
  uint32_t i = 0;
  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size) &&
         (!is_lbi_found)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
      spgw_ctxt_p = node->data;
      num_elements++;
      if (spgw_ctxt_p != NULL) {
        if (!strcmp(
              (const char*)
                spgw_ctxt_p->sgw_eps_bearer_context_information.imsi.digit,
              (const char*) bearer_req_p->imsi)) {
          is_imsi_found = true;
          if (
            (bearer_req_p->lbi != 0) &&
            (bearer_req_p->lbi ==
             spgw_ctxt_p->sgw_eps_bearer_context_information.pdn_connection
               .default_bearer)) {
            is_lbi_found = true;
            // Check if the received EBI is valid
            for (uint32_t itrn = 0; itrn < bearer_req_p->no_of_bearers;
                 itrn++) {
              if (sgw_cm_get_eps_bearer_entry(
                    &spgw_ctxt_p->sgw_eps_bearer_context_information
                       .pdn_connection,
                    bearer_req_p->ebi[itrn])) {
                is_ebi_found = true;
                ebi_to_be_deactivated[no_of_bearers_to_be_deact] =
                  bearer_req_p->ebi[itrn];
                no_of_bearers_to_be_deact++;
              } else {
                invalid_bearer_id[no_of_bearers_rej] = bearer_req_p->ebi[itrn];
                no_of_bearers_rej++;
              }
            }
          }
        }
      }
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    i++;
  }

  /* Send reject to NW if we did not find ebi/lbi/imsi.
   * Also in case of multiple bearers, if some EBIs are valid and some are not,
   * send reject to those for which we did not find EBI.
   * Proceed with deactivation by sending s5_nw_init_deactv_bearer_request to
   * SGW for valid EBIs
   */
  if ((!is_ebi_found) || (!is_lbi_found) || (!is_imsi_found) ||
    (no_of_bearers_rej > 0)) {
    OAILOG_INFO(
      LOG_SPGW_APP,
      "is_imsi_found (%d), is_lbi_found (%d), is_ebi_found (%d) \n",
      is_imsi_found,
      is_lbi_found,
      is_ebi_found);
    OAILOG_ERROR(
      LOG_SPGW_APP, "Sending dedicated bearer deactivation reject to NW\n");
    print_bearer_ids_helper(invalid_bearer_id, no_of_bearers_rej);
    // TODO-Uncomment once implemented at PCRF
    /* rc = send_dedicated_bearer_deactv_rsp(invalid_bearer_id,
         REQUEST_REJECTED);*/
  }

  if (no_of_bearers_to_be_deact > 0) {
    bool delete_default_bearer =
      (bearer_req_p->lbi == bearer_req_p->ebi[0]) ? true : false;
    rc = _spgw_build_and_send_s11_deactivate_bearer_req(
      imsi64,
      no_of_bearers_to_be_deact,
      ebi_to_be_deactivated,
      delete_default_bearer,
      spgw_ctxt_p->sgw_eps_bearer_context_information.mme_teid_S11);
  }
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

// Send ITTI message,S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST to mme_app
static int32_t _spgw_build_and_send_s11_deactivate_bearer_req(
  imsi64_t imsi64,
  uint8_t no_of_bearers_to_be_deact,
  ebi_t* ebi_to_be_deactivated,
  bool delete_default_bearer,
  teid_t mme_teid_S11)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  MessageDef* message_p = itti_alloc_new_message(
    TASK_SPGW_APP, S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "itti_alloc_new_message failed for nw_initiated_deactv_bearer_req\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  itti_s11_nw_init_deactv_bearer_request_t* s11_bearer_deactv_request =
    &message_p->ittiMsg.s11_nw_init_deactv_bearer_request;
  memset(
    s11_bearer_deactv_request,
    0,
    sizeof(itti_s11_nw_init_deactv_bearer_request_t));

  s11_bearer_deactv_request->s11_mme_teid = mme_teid_S11;
  /* If default bearer has to be deleted then the EBI list in the received
   * pgw_nw_init_deactv_bearer_request message contains a single entry at 0th
   * index and LBI == bearer_req_p->ebi[0]
   */
  s11_bearer_deactv_request->delete_default_bearer = delete_default_bearer;
  s11_bearer_deactv_request->no_of_bearers = no_of_bearers_to_be_deact;

  memcpy(
    s11_bearer_deactv_request->ebi,
    ebi_to_be_deactivated,
    (sizeof(ebi_t) * no_of_bearers_to_be_deact));
  print_bearer_ids_helper(
    s11_bearer_deactv_request->ebi, s11_bearer_deactv_request->no_of_bearers);

  message_p->ittiMsgHeader.imsi = imsi64;
  OAILOG_INFO(
    LOG_SPGW_APP,
    "Sending nw_initiated_deactv_bearer_req to mme_app "
    "with delete_default_bearer flag set to %d\n",
    s11_bearer_deactv_request->delete_default_bearer);
  int rc = itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

//------------------------------------------------------------------------------
int spgw_send_nw_init_activate_bearer_rsp(
  gtpv2c_cause_value_t cause,
  imsi64_t imsi64,
  uint8_t eps_bearer_id)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  uint32_t rc = RETURNok;

  OAILOG_INFO(
    LOG_SPGW_APP,
    "To be implemented: Sending Create Bearer Rsp to PCRF with EBI %d with "
    "cause :%d \n",
    eps_bearer_id,
    cause);
  // Send Create Bearer Rsp to PCRF
  // TODO-Uncomment once implemented at PCRF
  /* rc = send_dedicated_bearer_actv_rsp(act_ded_bearer_rsp->ebi,
       act_ded_bearer_rsp->cause);*/
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

//------------------------------------------------------------------------------
uint32_t spgw_handle_nw_init_deactivate_bearer_rsp(
  gtpv2c_cause_t cause,
  ebi_t lbi)
{
  uint32_t rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SPGW_APP);

  OAILOG_INFO(
    LOG_SPGW_APP,
    "To be implemented: Sending Delete Bearer Rsp to PCRF with LBI %u with "
    "cause :%d\n",
    lbi,
    cause.cause_value);
  // Send Delete Bearer Rsp to PCRF
  // TODO-Uncomment once implemented at PCRF
  // rc = send_dedicated_bearer_deactv_rsp(lbi, cause);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

// Build and send ITTI message, s11_create_bearer_request to MME APP
static int _spgw_build_and_send_s11_create_bearer_request(
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p,
  const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
  spgw_state_t* spgw_state,
  teid_t s1_u_sgw_fteid)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  MessageDef* message_p = NULL;
  int rc = RETURNerror;

  message_p = itti_alloc_new_message(
    TASK_SPGW_APP, S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST);
  if (!message_p) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to allocate message_p for"
      "S11_NW_INITIATED_BEARER_ACTV_REQUEST\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
  }

  itti_s11_nw_init_actv_bearer_request_t* s11_actv_bearer_request =
    &message_p->ittiMsg.s11_nw_init_actv_bearer_request;
  memset(
    s11_actv_bearer_request, 0, sizeof(itti_s11_nw_init_actv_bearer_request_t));
  // Context TEID
  s11_actv_bearer_request->s11_mme_teid =
    spgw_ctxt_p->sgw_eps_bearer_context_information.mme_teid_S11;
  // LBI
  s11_actv_bearer_request->lbi = bearer_req_p->lbi;
  // UL TFT to be sent to UE
  memcpy(
    &s11_actv_bearer_request->tft,
    &bearer_req_p->ul_tft,
    sizeof(traffic_flow_template_t));
  // QoS
  memcpy(
    &s11_actv_bearer_request->eps_bearer_qos,
    &bearer_req_p->eps_bearer_qos,
    sizeof(bearer_qos_t));
  // S1U SGW F-TEID
  s11_actv_bearer_request->s1_u_sgw_fteid.teid = s1_u_sgw_fteid;
  s11_actv_bearer_request->s1_u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
  // Set IPv4 address type bit
  s11_actv_bearer_request->s1_u_sgw_fteid.ipv4 = true;

  // TODO - IPv6 address
  s11_actv_bearer_request->s1_u_sgw_fteid.ipv4_address.s_addr =
    spgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr;
  message_p->ittiMsgHeader.imsi =
    spgw_ctxt_p->sgw_eps_bearer_context_information.imsi64;
  OAILOG_INFO(
    LOG_SPGW_APP,
    "Sending S11 Create Bearer Request to MME_APP for LBI %d IMSI " IMSI_64_FMT,
    bearer_req_p->lbi,
    message_p->ittiMsgHeader.imsi);
  rc = itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, rc);
}

// Create temporary dedicated bearer context
static int _create_temporary_dedicated_bearer_context(
  s_plus_p_gw_eps_bearer_context_information_t* spgw_ctxt_p,
  const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
  spgw_state_t* spgw_state,
  teid_t s1_u_sgw_fteid)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  sgw_eps_bearer_ctxt_t* eps_bearer_ctxt_p =
    calloc(1, sizeof(sgw_eps_bearer_ctxt_t));

  if (!eps_bearer_ctxt_p) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "Failed to allocate memory for eps_bearer_ctxt_p\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  // Copy PAA from default bearer cntxt
  sgw_eps_bearer_ctxt_t* default_eps_bearer_entry_p =
    sgw_cm_get_eps_bearer_entry(
      &spgw_ctxt_p->sgw_eps_bearer_context_information.pdn_connection,
      spgw_ctxt_p->sgw_eps_bearer_context_information.pdn_connection
        .default_bearer);

  if (!default_eps_bearer_entry_p) {
    OAILOG_ERROR(LOG_SPGW_APP, "Failed to get default bearer context\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }

  eps_bearer_ctxt_p->eps_bearer_id = 0;
  eps_bearer_ctxt_p->paa = default_eps_bearer_entry_p->paa;
  // SGW FTEID
  eps_bearer_ctxt_p->s_gw_teid_S1u_S12_S4_up = s1_u_sgw_fteid;

  eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.pdn_type = IPv4;
  eps_bearer_ctxt_p->s_gw_ip_address_S1u_S12_S4_up.address.ipv4_address.s_addr =
    spgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr;
  // DL TFT
  memcpy(
    &eps_bearer_ctxt_p->tft,
    &bearer_req_p->dl_tft,
    sizeof(traffic_flow_template_t));
  // QoS
  memcpy(
    &eps_bearer_ctxt_p->eps_bearer_qos,
    &bearer_req_p->eps_bearer_qos,
    sizeof(bearer_qos_t));

  OAILOG_INFO(
    LOG_SPGW_APP,
    "Number of DL packet filter rules: %d\n",
    eps_bearer_ctxt_p->tft.numberofpacketfilters);

  // Create temporary spgw bearer context entry
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
    pgw_get_procedure_create_bearer(spgw_ctxt_p);
  if (!pgw_ni_cbr_proc) {
    OAILOG_DEBUG(
      LOG_SPGW_APP, "Creating a new temporary eps bearer context entry\n");
    pgw_ni_cbr_proc = pgw_create_procedure_create_bearer(spgw_ctxt_p);
    if (!pgw_ni_cbr_proc) {
      OAILOG_ERROR(
        LOG_SPGW_APP, "Failed to create temporary eps bearer context entry\n");
      OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
    }
  }
  struct sgw_eps_bearer_entry_wrapper_s* sgw_eps_bearer_entry_p =
    calloc(1, sizeof(*sgw_eps_bearer_entry_p));
  if (!sgw_eps_bearer_entry_p) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "Failed to allocate memory for sgw_eps_bearer_entry_p\n");
    OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNerror);
  }
  sgw_eps_bearer_entry_p->sgw_eps_bearer_entry = eps_bearer_ctxt_p;
  LIST_INSERT_HEAD(
    (pgw_ni_cbr_proc->pending_eps_bearers), sgw_eps_bearer_entry_p, entries);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

// Deletes temporary dedicated bearer context
static void _delete_temporary_dedicated_bearer_context(
  teid_t s1_u_sgw_fteid,
  ebi_t lbi,
  s_plus_p_gw_eps_bearer_context_information_t* spgw_context_p)
{
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc = NULL;
  struct sgw_eps_bearer_entry_wrapper_s* spgw_eps_bearer_entry_p = NULL;
  pgw_ni_cbr_proc = pgw_get_procedure_create_bearer(spgw_context_p);
  if (!pgw_ni_cbr_proc) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to get Create bearer procedure from temporary stored contexts "
      "for lbi :%u \n",
      lbi);
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  OAILOG_INFO(
    LOG_SPGW_APP, "Delete temporary bearer context for lbi :%u \n", lbi);
  spgw_eps_bearer_entry_p = LIST_FIRST(pgw_ni_cbr_proc->pending_eps_bearers);
  while (spgw_eps_bearer_entry_p) {
    if (
      s1_u_sgw_fteid ==
      spgw_eps_bearer_entry_p->sgw_eps_bearer_entry->s_gw_teid_S1u_S12_S4_up) {
      // Remove the temporary spgw entry
      LIST_REMOVE(spgw_eps_bearer_entry_p, entries);
      if (spgw_eps_bearer_entry_p->sgw_eps_bearer_entry) {
        free_wrapper((void**) &spgw_eps_bearer_entry_p->sgw_eps_bearer_entry);
      }
      free_wrapper((void**) &spgw_eps_bearer_entry_p);
      break;
    }
    spgw_eps_bearer_entry_p = LIST_NEXT(spgw_eps_bearer_entry_p, entries);
  }
  if (LIST_EMPTY(pgw_ni_cbr_proc->pending_eps_bearers)) {
    pgw_free_procedure_create_bearer((pgw_ni_cbr_proc_t**) &pgw_ni_cbr_proc);
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}
