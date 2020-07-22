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

/*! \file s11_mme_bearer_manager.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>

#include "bstrlib.h"

#include "hashtable.h"
#include "obj_hashtable.h"
#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"

#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgParser.h"

#include "s11_common.h"
#include "s11_mme_bearer_manager.h"
#include "gtpv2c_ie_formatter.h"
#include "s11_ie_formatter.h"

extern hash_table_ts_t* s11_mme_teid_2_gtv2c_teid_handle;

//------------------------------------------------------------------------------
int s11_mme_release_access_bearers_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_release_access_bearers_request_t* req_p) {
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_rc_t rc;

  DevAssert(stack_p);
  DevAssert(req_p);
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;

  // Prepare a new Create Session Request msg
  rc = nwGtpv2cMsgNew(
      *stack_p, true, NW_GTP_RELEASE_ACCESS_BEARERS_REQ, req_p->teid, 0,
      &(ulp_req.hMsg));
  ulp_req.u_api_info.initialReqInfo.edns_peer_ip =
      (struct sockaddr*) &req_p->edns_peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal = req_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
      s11_mme_teid_2_gtv2c_teid_handle,
      (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
      (void**) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
        LOG_S11, "Could not get GTPv2-C hTunnel for local teid %X\n",
        ulp_req.u_api_info.initialReqInfo.teidLocal);
    rc = nwGtpv2cMsgDelete(*stack_p, (ulp_req.hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  // TODO add node_type_t originating_node if ISR active
  rc = nwGtpv2cMsgAddIe(
      (ulp_req.hMsg), NW_GTPV2C_IE_NODE_TYPE, 1, 0,
      (uint8_t*) &req_p->originating_node);
  DevAssert(NW_OK == rc);

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_downlink_data_notification_acknowledge(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_downlink_data_notification_acknowledge_t* ack_p) {
  nw_gtpv2c_ulp_api_t ulp_ack;
  gtpv2c_cause_t cause;
  nw_rc_t rc;
  nw_gtpv2c_trxn_handle_t trxn;

  DevAssert(stack_p);
  DevAssert(ack_p);
  memset(&ulp_ack, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_ack.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;

  trxn = (nw_gtpv2c_trxn_handle_t) ack_p->trxn;

  memset(
      &cause, 0,
      sizeof(
          gtpv2c_cause_t));  // Prepare a create bearer response to send to SGW.
  ulp_ack.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc                                        = nwGtpv2cMsgNew(
      *stack_p, true, NW_GTP_DOWNLINK_DATA_NOTIFICATION_ACK, ack_p->teid, 0,
      &(ulp_ack.hMsg));
  DevAssert(NW_OK == rc);

  ulp_ack.u_api_info.triggeredRspInfo.teidLocal =
      ack_p->local_teid;  // Set the remote TEID

  hashtable_rc_t hash_rc = hashtable_ts_get(
      s11_mme_teid_2_gtv2c_teid_handle, (hash_key_t) ack_p->local_teid,
      (void**) (uintptr_t) &ulp_ack.u_api_info.triggeredRspInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
        LOG_S11, "Could not get GTPv2-C hTunnel for local teid %X\n",
        ack_p->local_teid);
    rc = nwGtpv2cMsgDelete(*stack_p, (ulp_ack.hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  // TODO relay cause
  cause = ack_p->cause;
  gtpv2c_cause_ie_set(&(ulp_ack.hMsg), &cause);
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_ack);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_handle_release_access_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_release_access_bearers_response_t* resp_p;
  MessageDef* message_p;
  nw_gtpv2c_msg_parser_t* pMsgParser;

  DevAssert(stack_p);
  message_p =
      itti_alloc_new_message(TASK_S11, S11_RELEASE_ACCESS_BEARERS_RESPONSE);
  resp_p = &message_p->ittiMsg.s11_release_access_bearers_response;

  resp_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);

  // Create a new message parser
  rc = nwGtpv2cMsgParserNew(
      *stack_p, NW_GTP_RELEASE_ACCESS_BEARERS_RSP, s11_ie_indication_generic,
      NULL, &pMsgParser);
  DevAssert(NW_OK == rc);

  // Cause IE
  rc = nwGtpv2cMsgParserAddIe(
      pMsgParser, NW_GTPV2C_IE_CAUSE, NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_MANDATORY, gtpv2c_cause_ie_get, &resp_p->cause);
  DevAssert(NW_OK == rc);

  // Run the parser
  rc = nwGtpv2cMsgParserRun(
      pMsgParser, (pUlpApi->hMsg), &offendingIeType, &offendingIeInstance,
      &offendingIeLength);

  if (rc != NW_OK) {
    // TODO: handle this case
    free(message_p);
    message_p = NULL;
    rc        = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);
  return send_msg_to_task(&s11_task_zmq_ctx, TASK_MME_APP, message_p);
}

//------------------------------------------------------------------------------
int s11_mme_modify_bearer_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_modify_bearer_request_t* req_p) {
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_rc_t rc;

  DevAssert(stack_p);
  DevAssert(req_p);
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;

  // Prepare a new Modify Bearer Request msg
  rc = nwGtpv2cMsgNew(
      *stack_p, true, NW_GTP_MODIFY_BEARER_REQ, req_p->teid, 0,
      &(ulp_req.hMsg));
  ulp_req.u_api_info.initialReqInfo.edns_peer_ip =
      (struct sockaddr*) &req_p->edns_peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal      = req_p->local_teid;
  ulp_req.u_api_info.initialReqInfo.internal_flags = req_p->internal_flags;

  hashtable_rc_t hash_rc = hashtable_ts_get(
      s11_mme_teid_2_gtv2c_teid_handle,
      (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
      (void**) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
        LOG_S11, "Could not get GTPv2-C hTunnel for local teid %X\n",
        ulp_req.u_api_info.initialReqInfo.teidLocal);
    rc = nwGtpv2cMsgDelete(*stack_p, (ulp_req.hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  // Sender F-TEID for Control Plane (MME S11)
  rc = nwGtpv2cMsgAddIeFteid(
      (ulp_req.hMsg), NW_GTPV2C_IE_INSTANCE_ZERO, S11_MME_GTP_C,
      req_p->sender_fteid_for_cp.teid,
      req_p->sender_fteid_for_cp.ipv4 ?
          &req_p->sender_fteid_for_cp.ipv4_address :
          0,
      req_p->sender_fteid_for_cp.ipv6 ?
          &req_p->sender_fteid_for_cp.ipv6_address :
          NULL);

  for (int i = 0; i < req_p->bearer_contexts_to_be_modified.num_bearer_context;
       i++) {
    rc =
        gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_set(
            &(ulp_req.hMsg),
            &req_p->bearer_contexts_to_be_modified.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }

  for (int i = 0; i < req_p->bearer_contexts_to_be_removed.num_bearer_context;
       i++) {
    rc =
        gtpv2c_bearer_context_to_be_removed_within_modify_bearer_request_ie_set(
            &(ulp_req.hMsg),
            &req_p->bearer_contexts_to_be_removed.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_handle_modify_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_modify_bearer_response_t* resp_p;
  MessageDef* message_p;
  nw_gtpv2c_msg_parser_t* pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_MODIFY_BEARER_RESPONSE);
  resp_p    = &message_p->ittiMsg.s11_modify_bearer_response;

  resp_p->teid           = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
  resp_p->internal_flags = pUlpApi->u_api_info.triggeredRspIndInfo.trx_flags;

  // Create a new message parser
  rc = nwGtpv2cMsgParserNew(
      *stack_p, NW_GTP_MODIFY_BEARER_RSP, s11_ie_indication_generic, NULL,
      &pMsgParser);
  DevAssert(NW_OK == rc);

  // Cause IE
  rc = nwGtpv2cMsgParserAddIe(
      pMsgParser, NW_GTPV2C_IE_CAUSE, NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_MANDATORY, gtpv2c_cause_ie_get, &resp_p->cause);
  DevAssert(NW_OK == rc);

  // Bearer Contexts Created IE
  rc = nwGtpv2cMsgParserAddIe(
      pMsgParser, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_CONDITIONAL, gtpv2c_bearer_context_modified_ie_get,
      &resp_p->bearer_contexts_modified);
  DevAssert(NW_OK == rc);

  /*
   * Bearer Contexts Marked For Removal IE.
   * todo: we only process everything we marked locally. Currently disregardign
   * this element
   */
  rc = nwGtpv2cMsgParserAddIe(
      pMsgParser, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ONE,
      NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
      gtpv2c_bearer_context_marked_for_removal_ie_get,
      &resp_p->bearer_contexts_marked_for_removal);

  // Run the parser
  rc = nwGtpv2cMsgParserRun(
      pMsgParser, (pUlpApi->hMsg), &offendingIeType, &offendingIeInstance,
      &offendingIeLength);

  if (rc != NW_OK) {
    // TODO: handle this case
    free(message_p);
    message_p = NULL;
    rc        = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);

  /** Check the cause. */
  if (resp_p->cause.cause_value == LATE_OVERLAPPING_REQUEST) {
    pUlpApi->u_api_info.triggeredRspIndInfo.trx_flags |=
        LATE_OVERLAPPING_REQUEST;
    OAILOG_WARNING(
        LOG_S11,
        "Received a late overlapping request (MBR). Not forwarding message to "
        "MME_APP layer. \n");
    free(message_p);
    message_p = NULL;
    return RETURNok;
  }

  return send_msg_to_task(&s11_task_zmq_ctx, TASK_MME_APP, message_p);
}

//------------------------------------------------------------------------------
int s11_mme_delete_bearer_command(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_delete_bearer_command_t* cmd_p) {
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_rc_t rc;

  DevAssert(stack_p);
  DevAssert(cmd_p);
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;
  ulp_req.apiType |= NW_GTPV2C_ULP_API_FLAG_IS_COMMAND_MESSAGE;

  // Prepare a new Delete Session Request msg
  rc = nwGtpv2cMsgNew(
      *stack_p, true, NW_GTP_DELETE_BEARER_CMD, cmd_p->teid, 0,
      &(ulp_req.hMsg));
  ulp_req.u_api_info.initialReqInfo.edns_peer_ip =
      (struct sockaddr*) &cmd_p->edns_peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal = cmd_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
      s11_mme_teid_2_gtv2c_teid_handle,
      (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
      (void**) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
        LOG_S11, "Could not get GTPv2-C hTunnel for local teid %X\n",
        ulp_req.u_api_info.initialReqInfo.teidLocal);
    rc = nwGtpv2cMsgDelete(*stack_p, (ulp_req.hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  // Add bearer contexts to be removed.
  for (int num_ebi = 0; num_ebi < cmd_p->ebi_list.num_ebi; num_ebi++) {
    rc = gtpv2c_bearer_context_ebi_only_ie_set(
        &(ulp_req.hMsg), cmd_p->ebi_list.ebis[num_ebi]);
    DevAssert(NW_OK == rc);
  }

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_handle_create_bearer_request(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_create_bearer_request_t* req_p;
  MessageDef* message_p;
  nw_gtpv2c_msg_parser_t* pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_CREATE_BEARER_REQUEST);

  if (message_p) {
    req_p = &message_p->ittiMsg.s11_create_bearer_request;

    req_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
    req_p->trxn = (void*) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;

    // Create a new message parser
    rc = nwGtpv2cMsgParserNew(
        *stack_p, NW_GTP_CREATE_BEARER_REQ, s11_ie_indication_generic, NULL,
        &pMsgParser);
    DevAssert(NW_OK == rc);

    rc = nwGtpv2cMsgParserAddIe(
        pMsgParser, NW_GTPV2C_IE_EBI, NW_GTPV2C_IE_INSTANCE_ZERO,
        NW_GTPV2C_IE_PRESENCE_MANDATORY, gtpv2c_ebi_ie_get,
        &req_p->linked_eps_bearer_id);
    DevAssert(NW_OK == rc);

    rc = nwGtpv2cMsgParserAddIe(
        pMsgParser, NW_GTPV2C_IE_PCO, NW_GTPV2C_IE_INSTANCE_ZERO,
        NW_GTPV2C_IE_PRESENCE_OPTIONAL, gtpv2c_pco_ie_get, &req_p->pco);
    DevAssert(NW_OK == rc);

    DevAssert(!&req_p->bearer_contexts);

    rc = nwGtpv2cMsgParserAddIe(
        pMsgParser, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO,
        NW_GTPV2C_IE_PRESENCE_MANDATORY,
        gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_get,
        &req_p->bearer_contexts);
    DevAssert(NW_OK == rc);

    /** Add the PTI to inform to UEs. */
    rc = nwGtpv2cMsgParserAddIe(
        pMsgParser, NW_GTPV2C_IE_PROCEDURE_TRANSACTION_ID,
        NW_GTPV2C_IE_INSTANCE_ZERO, NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
        gtpv2c_pti_ie_get, &req_p->pti);
    DevAssert(NW_OK == rc);

    // Run the parser
    rc = nwGtpv2cMsgParserRun(
        pMsgParser, (pUlpApi->hMsg), &offendingIeType, &offendingIeInstance,
        &offendingIeLength);

    if (rc != NW_OK) {
      // TODO: handle this case
      free(message_p);
      message_p = NULL;
      rc        = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
      DevAssert(NW_OK == rc);
      rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
      DevAssert(NW_OK == rc);
      return RETURNerror;
    }

    rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return send_msg_to_task(&s11_task_zmq_ctx, TASK_MME_APP, message_p);
  }
  return RETURNerror;
}

//------------------------------------------------------------------------------
int s11_mme_create_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_create_bearer_response_t* response_p) {
  gtpv2c_cause_t cause;
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;

  DevAssert(stack_p);
  DevAssert(response_p);
  trxn = (nw_gtpv2c_trxn_handle_t) response_p->trxn;

  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));  // Prepare a create bearer
                                                     // response to send to SGW.
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType                           = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc                                        = nwGtpv2cMsgNew(
      *stack_p, true, NW_GTP_CREATE_BEARER_RSP, response_p->teid, 0,
      &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  ulp_req.u_api_info.triggeredRspInfo.teidLocal = response_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
      s11_mme_teid_2_gtv2c_teid_handle, (hash_key_t) response_p->local_teid,
      (void**) (uintptr_t) &ulp_req.u_api_info.triggeredRspInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
        LOG_S11, "Could not get GTPv2-C hTunnel for local teid %X\n",
        response_p->local_teid);
    rc = nwGtpv2cMsgDelete(*stack_p, (ulp_req.hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  // TODO relay cause
  cause = response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
  if (cause.cause_value == TEMP_REJECT_HO_IN_PROGRESS)
    ulp_req.u_api_info.triggeredRspInfo.pt_trx =
        true; /**< Using boolean, such that not to add any dependencies in
                 NwGtpv2c.h etc.. */

  for (int i = 0; i < response_p->bearer_contexts.num_bearer_context; i++) {
    rc = gtpv2c_bearer_context_within_create_bearer_response_ie_set(
        &(ulp_req.hMsg), &response_p->bearer_contexts.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
/* @brief Handle Downlink Data Notification received from source MME. */
int s11_mme_handle_downlink_data_notification(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc = NW_OK;
  itti_s11_downlink_data_notification_t* notif_p;
  MessageDef* message_p;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S10, S11_DOWNLINK_DATA_NOTIFICATION);
  notif_p   = &message_p->ittiMsg.s11_downlink_data_notification;
  memset(notif_p, 0, sizeof(*notif_p));
  notif_p->teid = nwGtpv2cMsgGetTeid(
      pUlpApi->hMsg); /**< When the message is sent, this is the field,
where the MME_APP sets the destination TEID. In this case, at reception and
decoding, it is the local TEID, used to find the MME_APP ue_context. */
  notif_p->trxn = (void*) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;

  /** Message will not be removed as part of the transaction. */
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);

  return send_msg_to_task(&s11_task_zmq_ctx, TASK_MME_APP, message_p);
}
