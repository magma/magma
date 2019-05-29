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

/*! \file s11_sgw_session_manager.c
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
#include "s11_sgw_session_manager.h"
#include "s11_ie_formatter.h"
#include "log.h"

extern hash_table_ts_t *s11_sgw_teid_2_gtv2c_teid_handle;

//------------------------------------------------------------------------------
int s11_sgw_handle_create_session_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_create_session_request_t *create_session_request_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_CREATE_SESSION_REQUEST);
  create_session_request_p = &message_p->ittiMsg.s11_create_session_request;
  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_CREATE_SESSION_REQ,
    s11_ie_indication_generic,
    NULL,
    &pMsgParser);
  DevAssert(NW_OK == rc);
  /*
   * Imsi IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_IMSI,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_imsi_ie_get,
    &create_session_request_p->imsi);
  DevAssert(NW_OK == rc);
  /*
   * MSISDN IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_MSISDN,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_msisdn_ie_get,
    &create_session_request_p->msisdn);
  DevAssert(NW_OK == rc);
  /*
   * MEI IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_MEI,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_mei_ie_get,
    &create_session_request_p->mei);
  DevAssert(NW_OK == rc);
  /*
   * ULI IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_ULI,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_uli_ie_get,
    &create_session_request_p->uli);
  DevAssert(NW_OK == rc);
  /*
   * Serving Network IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_SERVING_NETWORK,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_serving_network_ie_get,
    &create_session_request_p->serving_network);
  DevAssert(NW_OK == rc);
  /*
   * RAT Type IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_RAT_TYPE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_rat_type_ie_get,
    &create_session_request_p->rat_type);
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
    &create_session_request_p->indication_flags);
  DevAssert(NW_OK == rc);
  /*
   * APN IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_APN,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_apn_ie_get,
    &create_session_request_p->apn);
  DevAssert(NW_OK == rc);
  /*
   * Selection Mode IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_SELECTION_MODE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    s11_ie_indication_generic,
    NULL);
  DevAssert(NW_OK == rc);
  /*
   * PDN Type IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_PDN_TYPE,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_pdn_type_ie_get,
    &create_session_request_p->pdn_type);
  DevAssert(NW_OK == rc);
  /*
   * PAA IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_PAA,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_paa_ie_get,
    &create_session_request_p->paa);
  DevAssert(NW_OK == rc);
  /*
   * Sender FTEID for CP IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_FTEID,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_fteid_ie_get,
    &create_session_request_p->sender_fteid_for_cp);
  DevAssert(NW_OK == rc);
  /*
   * PGW FTEID for CP IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_FTEID,
    NW_GTPV2C_IE_INSTANCE_ONE,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_fteid_ie_get,
    &create_session_request_p->pgw_address_for_cp);
  DevAssert(NW_OK == rc);
  /*
   * APN Restriction IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_APN_RESTRICTION,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    s11_ie_indication_generic,
    NULL);
  DevAssert(NW_OK == rc);
  /*
   * Bearer Context IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_BEARER_CONTEXT,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    gtpv2c_bearer_context_to_be_created_within_create_session_request_ie_get,
    &create_session_request_p->bearer_contexts_to_be_created);
  DevAssert(NW_OK == rc);

  /*
   * Protocol Configuration Options IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_PCO,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_pco_ie_get,
    &create_session_request_p->pco);
  DevAssert(NW_OK == rc);

  /*TODO rc = nwGtpv2cMsgParserAddIe (pMsgParser, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ONE, NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
           s11_bearer_context_to_be_removed_ie_get, &create_session_request_p->bearer_contexts_to_be_removed);
  DevAssert (NW_OK == rc);*/

  /*
   * AMBR IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_AMBR,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_ambr_ie_get,
    &create_session_request_p->ambr);
  DevAssert(NW_OK == rc);
  /*
   * Recovery IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_RECOVERY,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_MANDATORY,
    s11_ie_indication_generic,
    NULL);
  DevAssert(NW_OK == rc);
  create_session_request_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
  create_session_request_p->trxn =
    (void *) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
  create_session_request_p->peer_ip =
    pUlpApi->u_api_info.initialReqIndInfo.peerIp;
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
     * Send Create session response with failure to Gtpv2c Stack Instance
     */
    ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
    ulp_req.u_api_info.triggeredRspInfo.hTrxn =
      pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
    rc = nwGtpv2cMsgNew(
      *stack_p,
      true,
      NW_GTP_CREATE_SESSION_RSP,
      0,
      nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg),
      &(ulp_req.hMsg));
    gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
    OAILOG_DEBUG(
      LOG_S11,
      "Received NW_GTP_CREATE_SESSION_REQ, Sending "
      "NW_GTP_CREATE_SESSION_RSP!\n");
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
  return itti_send_msg_to_task(TASK_SPGW_APP, message_p);
}

//------------------------------------------------------------------------------
int s11_sgw_handle_create_session_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_create_session_response_t *create_session_response_p)
{
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;
  gtpv2c_cause_t cause;

  DevAssert(create_session_response_p);
  DevAssert(stack_p);
  trxn = (nw_gtpv2c_trxn_handle_t) create_session_response_p->trxn;
  DevAssert(trxn);
  /*
   * Create a tunnel for the GTPv2-C stack
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  ulp_req.apiType = NW_GTPV2C_ULP_CREATE_LOCAL_TUNNEL;
  ulp_req.u_api_info.createLocalTunnelInfo.teidLocal =
    create_session_response_p->s11_sgw_fteid.teid;
  ulp_req.u_api_info.createLocalTunnelInfo.peerIp.s_addr =
    create_session_response_p->peer_ip.s_addr;
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);

  hashtable_rc_t hash_rc = hashtable_ts_insert(
    s11_sgw_teid_2_gtv2c_teid_handle,
    (hash_key_t) create_session_response_p->s11_sgw_fteid.teid,
    (void *) ulp_req.u_api_info.createLocalTunnelInfo.hTunnel);

  if (HASH_TABLE_OK != hash_rc) {
    OAILOG_WARNING(
      LOG_S11,
      "Could not save GTPv2-C hTunnel %p for local teid %X\n",
      (void *) ulp_req.u_api_info.createLocalTunnelInfo.hTunnel,
      create_session_response_p->s11_sgw_fteid.teid);
    return RETURNerror;
  }
  /*
   * Prepare a create session response to send to MME.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc = nwGtpv2cMsgNew(
    *stack_p, true, NW_GTP_CREATE_SESSION_RSP, 0, 0, &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  rc = nwGtpv2cMsgSetTeid(ulp_req.hMsg, create_session_response_p->teid);
  DevAssert(NW_OK == rc);
  cause = create_session_response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);

  rc = nwGtpv2cMsgAddIeFteid(
    (ulp_req.hMsg),
    NW_GTPV2C_IE_INSTANCE_ZERO,
    S11_SGW_GTP_C,
    create_session_response_p->s11_sgw_fteid.teid,
    create_session_response_p->s11_sgw_fteid.ipv4 ?
      &create_session_response_p->s11_sgw_fteid.ipv4_address :
      0,
    create_session_response_p->s11_sgw_fteid.ipv6 ?
      &create_session_response_p->s11_sgw_fteid.ipv6_address :
      NULL);

  rc = nwGtpv2cMsgAddIeFteid(
    (ulp_req.hMsg),
    NW_GTPV2C_IE_INSTANCE_ONE,
    S5_S8_PGW_GTP_C,
    create_session_response_p->s11_sgw_fteid.teid,
    create_session_response_p->s11_sgw_fteid.ipv4 ?
      &create_session_response_p->s11_sgw_fteid.ipv4_address :
      0,
    create_session_response_p->s11_sgw_fteid.ipv6 ?
      &create_session_response_p->s11_sgw_fteid.ipv6_address :
      NULL);

  gtpv2c_paa_ie_set(&(ulp_req.hMsg), &create_session_response_p->paa);
  /*
   * Put 0 for now i.e. no existing context or restriction
   */
  gtpv2c_apn_restriction_ie_set(&(ulp_req.hMsg), 0);
  gtpv2c_pco_ie_set(&(ulp_req.hMsg), &create_session_response_p->pco);

  for (int i = 0;
       i <
       create_session_response_p->bearer_contexts_created.num_bearer_context;
       i++) {
    gtpv2c_bearer_context_created_ie_set(
      &(ulp_req.hMsg),
      &create_session_response_p->bearer_contexts_created.bearer_contexts[i]);
  }

  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);

  return RETURNok;
}

//------------------------------------------------------------------------------
int s11_sgw_handle_delete_session_request(
  nw_gtpv2c_stack_handle_t *stack_p,
  nw_gtpv2c_ulp_api_t *pUlpApi)
{
  nw_rc_t rc = NW_OK;
  uint8_t offendingIeType, offendingIeInstance;
  uint16_t offendingIeLength;
  itti_s11_delete_session_request_t *delete_session_request_p;
  MessageDef *message_p;
  nw_gtpv2c_msg_parser_t *pMsgParser;

  DevAssert(stack_p);
  message_p = itti_alloc_new_message(TASK_S11, S11_DELETE_SESSION_REQUEST);
  delete_session_request_p = &message_p->ittiMsg.s11_delete_session_request;
  /*
   * Create a new message parser
   */
  rc = nwGtpv2cMsgParserNew(
    *stack_p,
    NW_GTP_DELETE_SESSION_REQ,
    s11_ie_indication_generic,
    NULL,
    &pMsgParser);
  DevAssert(NW_OK == rc);
  /*
   * MME FTEID for CP IE
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_FTEID,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_OPTIONAL,
    gtpv2c_fteid_ie_get,
    &delete_session_request_p->sender_fteid_for_cp);
  DevAssert(NW_OK == rc);
  /*
   * Linked EPS Bearer Id IE
   * * * * This information element shall not be present for TAU/RAU/Handover with
   * * * * S-GW relocation procedures.
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_EBI,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_OPTIONAL,
    gtpv2c_ebi_ie_get,
    &delete_session_request_p->lbi);
  DevAssert(NW_OK == rc);
  /*
   * Indication Flags IE
   * * * * For a Delete Session Request on S11 interface,
   * * * * only the Operation Indication flag might be present.
   */
  rc = nwGtpv2cMsgParserAddIe(
    pMsgParser,
    NW_GTPV2C_IE_INDICATION,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    NW_GTPV2C_IE_PRESENCE_CONDITIONAL,
    gtpv2c_indication_flags_ie_get,
    &delete_session_request_p->indication_flags);
  DevAssert(NW_OK == rc);
  delete_session_request_p->teid = nwGtpv2cMsgGetTeid(pUlpApi->hMsg);
  delete_session_request_p->trxn =
    (void *) pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
  delete_session_request_p->peer_ip =
    pUlpApi->u_api_info.initialReqIndInfo.peerIp;
  rc = nwGtpv2cMsgParserRun(
    pMsgParser,
    pUlpApi->hMsg,
    &offendingIeType,
    &offendingIeInstance,
    &offendingIeLength);

  if (rc != NW_OK) {
    nw_gtpv2c_ulp_api_t ulp_req;
    gtpv2c_cause_t cause = {0};

    memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
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
     * Send Create session response with failure to Gtpv2c Stack Instance
     */
    ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
    ulp_req.u_api_info.triggeredRspInfo.hTrxn =
      pUlpApi->u_api_info.initialReqIndInfo.hTrxn;
    rc = nwGtpv2cMsgNew(
      *stack_p,
      true,
      NW_GTP_DELETE_SESSION_RSP,
      0,
      nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg),
      &(ulp_req.hMsg));
    /*
     * Adding the cause
     */
    gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
    OAILOG_DEBUG(
      LOG_S11,
      "Received NW_GTP_DELETE_SESSION_REQ, Sending "
      "NW_GTP_DELETE_SESSION_RSP!\n");
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
  return itti_send_msg_to_task(TASK_SPGW_APP, message_p);
}

//------------------------------------------------------------------------------
int s11_sgw_handle_delete_session_response(
  nw_gtpv2c_stack_handle_t *stack_p,
  itti_s11_delete_session_response_t *delete_session_response_p)
{
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulp_req;
  nw_gtpv2c_trxn_handle_t trxn;
  gtpv2c_cause_t cause;

  DevAssert(delete_session_response_p);
  DevAssert(stack_p);
  trxn = (nw_gtpv2c_trxn_handle_t) delete_session_response_p->trxn;
  DevAssert(trxn);
  /*
   * Prepare a delete session response to send to MME.
   */
  memset(&ulp_req, 0, sizeof(nw_gtpv2c_ulp_api_t));
  memset(&cause, 0, sizeof(gtpv2c_cause_t));
  ulp_req.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
  ulp_req.u_api_info.triggeredRspInfo.hTrxn = trxn;
  rc = nwGtpv2cMsgNew(
    *stack_p, true, NW_GTP_DELETE_SESSION_RSP, 0, 0, &(ulp_req.hMsg));
  DevAssert(NW_OK == rc);
  /*
   * Set the remote TEID
   */
  rc = nwGtpv2cMsgSetTeid(ulp_req.hMsg, delete_session_response_p->teid);
  DevAssert(NW_OK == rc);
  cause = delete_session_response_p->cause;
  gtpv2c_cause_ie_set(&(ulp_req.hMsg), &cause);
  rc = nwGtpv2cProcessUlpReq(*stack_p, &ulp_req);
  DevAssert(NW_OK == rc);
  return RETURNok;
}
