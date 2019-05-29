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
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file s11_sgw_bearer_manager.c
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

#include "assertions.h"
#include "intertask_interface.h"
#include "queue.h"
#include "hashtable.h"
#include "NwLog.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgParser.h"
#include "sgw_ie_defs.h"
#include "s11_common.h"
#include "s11_sgw_bearer_manager.h"
#include "s11_ie_formatter.h"
#include "log.h"

extern hash_table_ts_t *s11_sgw_teid_2_gtv2c_teid_handle;

//------------------------------------------------------------------------------
int s11_sgw_handle_modify_bearer_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_modify_bearer_request_t *request_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_MODIFY_BEARER_REQUEST);
  request_p = &message_p->ittiMsg.s11_modify_bearer_request;
  request_p->trxn = (void *) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
  request_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_MODIFY_BEARER_REQ,
    s11_ie_indication_generic,
    NULL,
    &pMsgParser);
  DevAssert(NW_OK == rc);
  /*
   * Indication Flags IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_INDICATION,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_indication_flags_ie_get,
    &request_p->indication_flags);
  DevAssert(NW_OK == rc);
  /*
   * MME-FQ-CSID IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_FQ_CSID,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_fqcsid_ie_get,
    &request_p->mme_fq_csid);
  DevAssert(NW_OK == rc);
  /*
   * RAT Type IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_RAT_TYPE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_rat_type_ie_get,
    &request_p->rat_type);
  DevAssert(NW_OK == rc);
  /*
   * Delay Value IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_DELAY_VALUE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_delay_value_ie_get,
    &request_p->delay_dl_packet_notif_req);
  DevAssert(NW_OK == rc);
  /*
   * Bearer Context to be modified IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_BEARER_CONTEXT,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_get,
    &request_p->bearer_contexts_to_be_modified);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgParserRun(
    pMsgParser,
    pUlpApi->hMsg,
    &offendingIeType,
    &offendingIeInstance,
    &offendingIeLength);

  if (rc != NW_OK) {
    gtpv2c_cause_t cause;
    nw_gtpv2c_ulp_api_t ulp_req;

    memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
    memset(&cause, 0, sizeof(gtpv2c_cause_t));
    cause.offending_ie_type = offendingIeType;
    cause.offending_ie_length = offendingIeLength;
    cause.offending_ie_instance = offendingIeInstance;

    switch (rc) {
      case NW_GTPV2C_MANDATORY_IE_MISSING:
        OAILOG_DEBUG(
          LOG_S11,
          "Mandatory IE type '%u' of instance '%u' missing!\n",
          offendingIeType,
          offendingIeLength);
        cause.cause_value = NW_GTPV2C_CAUSE_MANDATORY_IE_MISSING;
        break;

      default:
        OAILOG_DEBUG(LOG_S11, "Unknown message parse error!\n");
        cause.cause_value = 0;
        break;
    }

    /*
     * Send Modify Bearer response with failure to Gtpv2c Stack Instance
     */
    ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
    ulp_req.u_api_info.triggeredRspInfo.hTrxn =
      pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
    rc = nwGtpv2cMsgNew(
      *stack_p,
      true,
      NW_GTP_MODIFY_BEARER_RSP,
      0,
      nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg),
      &(ulp_req.hMsg));
    gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
    OAILOG_DEBUG(
      LOG_S11,
      "Received NW_GTP_MODIFY_BEARER_REQ, Sending NW_GTP_MODIFY_BEARER_RSP!\n");
    rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
    DevAssert(NW_OK == rc);
    itti_free(ITTI_MSG_ORIGIN_ID(message_p), message_p);
    message_p = NULL;
    rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return NW_OK;
  }

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);
  return itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int s11_sgw_handle_modify_bearer_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_modify_bearer_response_t *response_p)
{
  gtpv2c_cause_t cause;
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;

  DevAssert(stack_p);
  DevAssert(response_p);
  trxn = (nw_gtpv2c_trxn_handle_t) response_p->trxn;
  /*
   * Prepare a modify bearer response to send to MME.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc = nwGtpv2cMsgNew(
    *stack_p, true, NW_GTP_MODIFY_BEARER_RSP, 0, 0, &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  rc = nwGtpv2cMsgSetTeid(ulp_req.hMsg, response_p->teid);
  DevAssert(NW_OK == rc);
  cause = response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_sgw_handle_release_access_bearers_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_release_access_bearers_request_t *request_p = NULL;
  MessageDef *message_p = NULL;
  nw_gtpv2c_msg_parser_t *pMsgParser = NULL;

  DevAssert(stack_p);
  message_p =
    itti_alloc_new_message(TASK_S11, S11_RELEASE_ACCESS_BEARERS_REQUEST);
  request_p = &message_p->ittiMsg.s11_release_access_bearers_request;

  request_p->trxn = (void *) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
  request_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_RELEASE_ACCESS_BEARERS_REQ,
    s11_ie_indication_generic,
    NULL,
    &pMsgParser);
  DevAssert(NW_OK == rc);

  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_NODE_TYPE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_node_type_ie_get,
    &request_p->originating_node);

  rc = nwGtpv2cMsgParserRun(
    pMsgParser,
    pUlpApi->hMsg,
    &offendingIeType,
    &offendingIeInstance,
    &offendingIeLength);

  if (rc != NW_OK) {
    gtpv2c_cause_t cause;
    nw_gtpv2c_ulp_api_t ulp_req;

    memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
    memset(&cause, 0, sizeof(gtpv2c_cause_t));
    cause.offending_ie_type = offendingIeType;
    cause.offending_ie_length = offendingIeLength;
    cause.offending_ie_instance = offendingIeInstance;

    switch (rc) {
      case NW_GTPV2C_MANDATORY_IE_MISSING:
        OAILOG_DEBUG(
          LOG_S11,
          "Mandatory IE type '%u' of instance '%u' missing!\n",
          offendingIeType,
          offendingIeLength);
        cause.cause_value = NW_GTPV2C_CAUSE_MANDATORY_IE_MISSING;
        break;

      default:
        OAILOG_DEBUG(LOG_S11, "Unknown message parse error!\n");
        cause.cause_value = 0;
        break;
    }

    /*
     * Send Release Access bearer response with failure to Gtpv2c Stack Instance
     */
    ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
    ulp_req.u_api_info.triggeredRspInfo.hTrxn =
      pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
    rc = nwGtpv2cMsgNew(
      *stack_p,
      true,
      NW_GTP_RELEASE_ACCESS_BEARERS_RSP,
      0,
      nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg),
      &(ulp_req.hMsg));
    gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
    OAILOG_DEBUG(
      LOG_S11,
      "Received NW_GTP_RELEASE_ACCESS_BEARERS_REQ, Sending "
      "NW_GTP_RELEASE_ACCESS_BEARERS_RSP!\n");
    rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
    DevAssert(NW_OK == rc);
    itti_free(ITTI_MSG_ORIGIN_ID(message_p), message_p);
    message_p = NULL;
    rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return RETURNok;
  }

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);

  return itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int s11_sgw_handle_release_access_bearers_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_release_access_bearers_response_t *response_p)
{
  gtpv2c_cause_t cause;
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;

  DevAssert(stack_p);
  DevAssert(response_p);
  trxn = (nw_gtpv2c_trxn_handle_t) response_p->trxn;
  /*
   * Prepare a release access bearer response to send to MME.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc = nwGtpv2cMsgNew(
    *stack_p, true, NW_GTP_RELEASE_ACCESS_BEARERS_RSP, 0, 0, &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  rc = nwGtpv2cMsgSetTeid(ulp_req.hMsg, response_p->teid);
  DevAssert(NW_OK == rc);
  //TODO relay cause
  cause = response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_sgw_handle_create_bearer_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_create_bearer_request_t *request_p)
{
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;

  DevAssert(stack_p);
  DevAssert(request_p);
  /*
   * Prepare a create bearer request to send to MME.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));

  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;
  rc = nwGtpv2cMsgNew(
    *stack_p,
    true,
    NW_GTP_CREATE_BEARER_REQ,
    request_p->teid,
    0,
    &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);

  ulp_req.u_api_info.initialReqInfo.peerIp = request_p->peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal = request_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
    s11_sgw_teid_2_gtv2c_teid_handle,
    (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
    (void **) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);
  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
      LOG_S11,
      "Could not get GTPv2-C hTunnel for local teid %X\n",
      ulp_req.u_api_info.initialReqInfo.teidLocal);
  }

  /*
   * Set the remote TEID
   */
  rc = nwGtpv2cMsgSetTeid(ulp_req.hMsg, request_p->teid);
  DevAssert(NW_OK == rc);

  // TODO   pti_t                      pti; ///< C: This IE shall be sent on the S5/S8 and S4/S11 interfaces

  gtpv2c_ebi_ie_set(
    &(ulp_req.hMsg), (unsigned) request_p->linked_eps_bearer_id);

  if (request_p->pco.num_protocol_or_container_id) {
    rc = gtpv2c_pco_ie_set(&(ulp_req.hMsg), &request_p->pco);
    DevAssert(NW_OK == rc);
  }

  for (int i = 0; i < request_p->bearer_contexts.num_bearer_context; i++) {
    rc =
      gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_set(
        &(ulp_req.hMsg), &request_p->bearer_contexts.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }

  // TODO pgw_fq_csid, sgw_fq_csid

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_sgw_handle_create_bearer_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_create_bearer_response_t *resp_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_CREATE_BEARER_RESPONSE);
  resp_p = &message_p->ittiMsg.s11_create_bearer_response;

  resp_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);

  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_CREATE_BEARER_RSP,
    s11_ie_indication_generic,
    NULL,
    &pMsgParser);
  DevAssert(NW_OK == rc);
  /*
   * Cause IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_CAUSE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_cause_ie_get,
    &resp_p->cause);
  DevAssert(NW_OK == rc);

  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_BEARER_CONTEXT,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_bearer_context_within_create_bearer_response_ie_get,
    &resp_p->bearer_contexts);
  DevAssert(NW_OK == rc);

  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_PCO,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_OPTIONAL,
    gtpv2c_pco_ie_get,
    &resp_p->pco);
  DevAssert(NW_OK == rc);

  /*
   * Run the parser
   */
  rc = nwGtpv2cMsgParserRun(
    pMsgParser,
    (pUlpApi->hMsg),
    &offendingIeType,
    &offendingIeInstance,
    &offendingIeLength);

  if (rc != NW_OK) {
    MSC_LOG_RX_DISCARDED_MESSAGE(
      MSC_S11_MME,
      MSC_SGW,
      NULL,
      0,
      "0 CREATE_BEARER_RESPONSE local S11 teid " TEID_FMT " ",
      resp_p->teid);
    /*
     * TODO: handle this case
     */
    itti_free(ITTI_MSG_ORIGIN_ID(message_p), message_p);
    message_p = NULL;
    rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return RETURNerror;
  }

  MSC_LOG_RX_MESSAGE(
    MSC_S11_MME,
    MSC_SGW,
    NULL,
    0,
    "0 CREATE_BEARER_RESPONSE local S11 teid " TEID_FMT " cause %u",
    resp_p->teid,
    resp_p->cause);

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);
  return itti_send_msg_to_task(TASK_SPGW_APP, INSTANCE_DEFAULT, message_p);
}
