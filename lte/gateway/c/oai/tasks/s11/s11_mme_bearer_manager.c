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
#include "msc.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "msc.h"

#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgParser.h"

#include "s11_common.h"
#include "s11_mme_bearer_manager.h"
#include "s11_ie_formatter.h"

extern hash_table_ts_t *s11_mme_teid_2_gtv2c_teid_handle;

//------------------------------------------------------------------------------
int s11_mme_release_access_bearers_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_release_access_bearers_request_t *req_p)
{
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_rc_t rc;
  //uint8_t                                 restart_counter = 0;

  DevAssert(stack_p);
  DevAssert(req_p);
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;
  /*
   * Prepare a new Create Session Request msg
   */
  rc = nwGtpv2cMsgNew(
    *stack_p,
    true,
    NW_GTP_RELEASE_ACCESS_BEARERS_REQ,
    req_p->teid,
    0,
    &(ulp_req.hMsg));
  ulp_req.u_api_info.initialReqInfo.peerIp = req_p->peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal = req_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
    s11_mme_teid_2_gtv2c_teid_handle,
    (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
    (void **) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
      LOG_S11,
      "Could not get GTPv2-C hTunnel for local teid %X\n",
      ulp_req.u_api_info.initialReqInfo.teidLocal);
    return RETURNerror;
  }

  // TODO add node_type_t originating_node if ISR active
  rc = nwGtpv2cMsgAddIe(
    (ulp_req.hMsg),
    NW_GTPV2C_IE_NODE_TYPE,
    1,
    0,
    (uint8_t *) &req_p->originating_node);
  DevAssert(NW_OK == rc);

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  MSC_LOG_TX_MESSAGE(
    MSC_S11_MME,
    MSC_SGW,
    NULL,
    0,
    "0 RELEASE_ACCESS_BEARERS_REQUEST local S11 teid " TEID_FMT " ",
    req_p->local_teid);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_handle_release_access_bearer_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_release_access_bearers_response_t *resp_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p =
    itti_alloc_new_message(TASK_S11, S11_RELEASE_ACCESS_BEARERS_RESPONSE);
  resp_p = &message_p->ittiMsg.s11_release_access_bearers_response;

  resp_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);

  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_RELEASE_ACCESS_BEARERS_RSP,
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
  /*
   * Recovery IE
   */
  /*rc = nwGtpv2cMsgParserAddIe (pMsgParser, NW_GTPV2C_IE_RECOVERY, NW_GTPV2C_IE_INSTANCE_ZERO, NW_GTPV2C_IE_PRESENCE_OPTIONAL, s11_fteid_ie_get,
		  &resp_p->recovery);
  DevAssert (NW_OK == rc);*/

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
      "0 RELEASE_ACCESS_BEARERS_RESPONSE local S11 teid " TEID_FMT " ",
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
    "0 RELEASE_ACCESS_BEARERS_RESPONSE local S11 teid " TEID_FMT " cause %u",
    resp_p->teid,
    resp_p->cause);

  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);
  return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int s11_mme_modify_bearer_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_modify_bearer_request_t *req_p)
{
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_rc_t rc;
  //uint8_t                                 restart_counter = 0;

  DevAssert(stack_p);
  DevAssert(req_p);
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ;
  /*
   * Prepare a new Modify Bearer Request msg
   */
  rc = nwGtpv2cMsgNew(
    *stack_p, true, NW_GTP_MODIFY_BEARER_REQ, req_p->teid, 0, &(ulp_req.hMsg));
  ulp_req.u_api_info.initialReqInfo.peerIp = req_p->peer_ip;
  ulp_req.u_api_info.initialReqInfo.teidLocal = req_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
    s11_mme_teid_2_gtv2c_teid_handle,
    (hash_key_t) ulp_req.u_api_info.initialReqInfo.teidLocal,
    (void **) (uintptr_t) &ulp_req.u_api_info.initialReqInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
      LOG_S11,
      "Could not get GTPv2-C hTunnel for local teid %X\n",
      ulp_req.u_api_info.initialReqInfo.teidLocal);
    return RETURNerror;
  }

  /*
   * Sender F-TEID for Control Plane (MME S11)
   */
  rc = nwGtpv2cMsgAddIeFteid(
    (ulp_req.hMsg),
    NW_GTPV2C_IE_INSTANCE_ZERO,
    S11_MME_GTP_C,
    req_p->sender_fteid_for_cp.teid,
    req_p->sender_fteid_for_cp.ipv4 ? &req_p->sender_fteid_for_cp.ipv4_address :
                                      0,
    req_p->sender_fteid_for_cp.ipv6 ? &req_p->sender_fteid_for_cp.ipv6_address :
                                      NULL);

  for (int i = 0; i < req_p->bearer_contexts_to_be_modified.num_bearer_context;
       i++) {
    rc =
      gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_set(
        &(ulp_req.hMsg),
        &req_p->bearer_contexts_to_be_modified.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }

  MSC_LOG_TX_MESSAGE(
    MSC_S11_MME,
    MSC_SGW,
    NULL,
    0,
    "0 MODIFY_BEARER_REQUEST local S11 teid " TEID_FMT " num bearers ctx %u",
    req_p->local_teid,
    req_p->bearer_contexts_to_be_modified.num_bearer_context);

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_create_bearer_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_create_bearer_response_t *response_p)
{
  gtpv2c_cause_t cause;
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;

  DevAssert(stack_p);
  DevAssert(response_p);
  trxn = (nw_gtpv2c_trxn_handle_t) response_p->trxn;
  /*
   * Prepare a create bearer response to send to SGW.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc = nwGtpv2cMsgNew(
    *stack_p,
    true,
    NW_GTP_CREATE_BEARER_RSP,
    response_p->teid,
    0,
    &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  ulp_req.u_api_info.triggeredRspInfo.teidLocal = response_p->local_teid;

  hashtable_rc_t hash_rc = hashtable_ts_get(
    s11_mme_teid_2_gtv2c_teid_handle,
    (hash_key_t) response_p->local_teid,
    (void **) (uintptr_t) &ulp_req.u_api_info.triggeredRspInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
      LOG_S11,
      "Could not get GTPv2-C hTunnel for local teid %X\n",
      response_p->local_teid);
    return RETURNerror;
  }

  // TODO relay cause
  cause = response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);

  for (int i = 0; i < response_p->bearer_contexts.num_bearer_context; i++) {
    rc = gtpv2c_bearer_context_within_create_bearer_response_ie_set(
      &(ulp_req.hMsg), &response_p->bearer_contexts.bearer_contexts[i]);
    DevAssert(NW_OK == rc);
  }

  MSC_LOG_TX_MESSAGE(
    MSC_S11_MME,
    MSC_SGW,
    NULL,
    0,
    "0 CREATE_BEARER_RESPONSE S11 teid " TEID_FMT " num bearers ctx %u",
    response_p->teid,
    response_p->bearer_contexts.num_bearer_context);

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_mme_handle_modify_bearer_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_modify_bearer_response_t *resp_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_MODIFY_BEARER_RESPONSE);
  resp_p = &message_p->ittiMsg.s11_modify_bearer_response;

  resp_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);

  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_MODIFY_BEARER_RSP,
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
  /*
   * Recovery IE
   */
  /*rc = nwGtpv2cMsgParserAddIe (pMsgParser, NW_GTPV2C_IE_RECOVERY, NW_GTPV2C_IE_INSTANCE_ZERO, NW_GTPV2C_IE_PRESENCE_OPTIONAL, s11_fteid_ie_get,
		  &resp_p->recovery);
  DevAssert (NW_OK == rc);*/

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
      "0 MODIFY_BEARER_RESPONSE local S11 teid " TEID_FMT " ",
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
    "0 MODIFY_BEARER_RESPONSE local S11 teid " TEID_FMT " cause %u",
    resp_p->teid,
    resp_p->cause);
  rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
  DevAssert(NW_OK == rc);
  return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int s11_mme_handle_create_bearer_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_create_bearer_request_t *req_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_CREATE_BEARER_REQUEST);

  if (message_p) {
    req_p = &message_p->ittiMsg.s11_create_bearer_request;

    req_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
    req_p->trxn = (void *) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;

    /*
     * Create a new message parser
     */
    rc = nwGtpv2cMsgParserNew(
      *stack_p,
      NW_GTP_CREATE_BEARER_REQ,
      s11_ie_indication_generic,
      NULL,
      &pMsgParser);
    DevAssert(NW_OK == rc);

    rc = nwGtpv2cMsgParserAddIe(
      pMsgParser,
      NW_GTPV2C_IE_EBI,
      NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_MANDATORY,
      gtpv2c_ebi_ie_get,
      &req_p->linked_eps_bearer_id);
    DevAssert(NW_OK == rc);

    rc = nwGtpv2cMsgParserAddIe(
      pMsgParser,
      NW_GTPV2C_IE_PCO,
      NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_OPTIONAL,
      gtpv2c_pco_ie_get,
      &req_p->pco);
    DevAssert(NW_OK == rc);

    rc = nwGtpv2cMsgParserAddIe(
      pMsgParser,
      NW_GTPV2C_IE_BEARER_CONTEXT,
      NW_GTPV2C_IE_INSTANCE_ZERO,
      NW_GTPV2C_IE_PRESENCE_MANDATORY,
      gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_get,
      &req_p->bearer_contexts);
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
        "0 CREATE_BEARER_REQUEST local S11 teid " TEID_FMT " ",
        req_p->teid);
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
      "0 CREATE_BEARER_REQUEST local S11 teid " TEID_FMT " lebi %u",
      req_p->teid,
      req_p->linked_eps_bearer_id);
    rc = nwGtpv2cMsgParserDelete(*stack_p, pMsgParser);
    DevAssert(NW_OK == rc);
    rc = nwGtpv2cMsgDelete(*stack_p, (pUlpApi->hMsg));
    DevAssert(NW_OK == rc);
    return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  }
  return RETURNerror;
}
