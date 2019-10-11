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

#include <stdint.h>
#include <string.h>

#include "s1ap_common.h"
#include "s1ap_ies_defs.h"
#include "s1ap_mme_encoder.h"
#include "assertions.h"
#include "log.h"
#include "S1AP-PDU.h"
#include "S1ap-Criticality.h"
#include "S1ap-DownlinkNASTransport.h"
#include "S1ap-E-RABSetupRequest.h"
#include "S1ap-InitialContextSetupRequest.h"
#include "S1ap-Paging.h"
#include "S1ap-ProcedureCode.h"
#include "S1ap-ResetAcknowledge.h"
#include "S1ap-S1SetupFailure.h"
#include "S1ap-S1SetupResponse.h"
#include "S1ap-UEContextModificationRequest.h"
#include "S1ap-UEContextReleaseCommand.h"
#include "S1ap-E-RABReleaseCommand.h"

static inline int s1ap_mme_encode_initial_context_setup_request(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
static inline int s1ap_mme_encode_s1setupresponse(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
static inline int s1ap_mme_encode_s1setupfailure(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
static inline int s1ap_mme_encode_ue_context_release_command(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
static inline int s1ap_mme_encode_downlink_nas_transport(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_e_rab_setup(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_initiating(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
static inline int s1ap_mme_encode_successfull_outcome(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *len);
static inline int s1ap_mme_encode_unsuccessfull_outcome(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *len);

static inline int s1ap_mme_encode_resetack(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_e_rab_release_command(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);
//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_paging(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_ue_context_modification_request(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_mme_configuration_transfer(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_pathswitchreqack(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_pathswitchreqfailure(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length);

static inline int s1ap_mme_encode_initial_context_setup_request(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_InitialContextSetupRequest_t initialContextSetupRequest;
  S1ap_InitialContextSetupRequest_t *initialContextSetupRequest_p =
    &initialContextSetupRequest;

  memset(
    initialContextSetupRequest_p, 0, sizeof(S1ap_InitialContextSetupRequest_t));

  if (
    s1ap_encode_s1ap_initialcontextsetuprequesties(
      initialContextSetupRequest_p,
      &message_p->msg.s1ap_InitialContextSetupRequestIEs) < 0) {
    return -1;
  }

  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_InitialContextSetup,
    S1ap_Criticality_reject,
    &asn_DEF_S1ap_InitialContextSetupRequest,
    initialContextSetupRequest_p);
}

//------------------------------------------------------------------------------
int s1ap_mme_encode_pdu(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  DevAssert(message_p != NULL);
  DevAssert(buffer != NULL);
  DevAssert(length != NULL);

  switch (message_p->direction) {
    case S1AP_PDU_PR_initiatingMessage:
      return s1ap_mme_encode_initiating(message_p, buffer, length);

    case S1AP_PDU_PR_successfulOutcome:
      return s1ap_mme_encode_successfull_outcome(message_p, buffer, length);

    case S1AP_PDU_PR_unsuccessfulOutcome:
      return s1ap_mme_encode_unsuccessfull_outcome(message_p, buffer, length);

    default:
      OAILOG_DEBUG(
        LOG_S1AP,
        "Unknown message outcome (%d) or not implemented",
        (int) message_p->direction);
      break;
  }

  return -1;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_initiating(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  switch (message_p->procedureCode) {
    case S1ap_ProcedureCode_id_downlinkNASTransport:
      return s1ap_mme_encode_downlink_nas_transport(message_p, buffer, length);

    case S1ap_ProcedureCode_id_InitialContextSetup:
      return s1ap_mme_encode_initial_context_setup_request(
        message_p, buffer, length);

    case S1ap_ProcedureCode_id_UEContextRelease:
      return s1ap_mme_encode_ue_context_release_command(
        message_p, buffer, length);

    case S1ap_ProcedureCode_id_Paging:
      return s1ap_mme_encode_paging(message_p, buffer, length);

    case S1ap_ProcedureCode_id_UEContextModification:
      return s1ap_mme_encode_ue_context_modification_request(
        message_p, buffer, length);

    case S1ap_ProcedureCode_id_E_RABSetup:
      return s1ap_mme_encode_e_rab_setup(message_p, buffer, length);

    case S1ap_ProcedureCode_id_E_RABRelease:
      return s1ap_mme_encode_e_rab_release_command(message_p, buffer, length);

    case S1ap_ProcedureCode_id_MMEConfigurationTransfer:
      return s1ap_mme_encode_mme_configuration_transfer(
        message_p, buffer, length);

    default:
      OAILOG_DEBUG(
        LOG_S1AP,
        "Unknown procedure ID (%d) for initiating message_p\n",
        (int) message_p->procedureCode);
      break;
  }

  return -1;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_successfull_outcome(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  switch (message_p->procedureCode) {
    case S1ap_ProcedureCode_id_S1Setup:
      return s1ap_mme_encode_s1setupresponse(message_p, buffer, length);
    case S1ap_ProcedureCode_id_Reset:
      return s1ap_mme_encode_resetack(message_p, buffer, length);
    case S1ap_ProcedureCode_id_PathSwitchRequest:
      return s1ap_mme_encode_pathswitchreqack(message_p, buffer, length);

    default:
      OAILOG_DEBUG(
        LOG_S1AP,
        "Unknown procedure ID (%d) for successfull outcome message\n",
        (int) message_p->procedureCode);
      break;
  }

  return -1;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_unsuccessfull_outcome(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  switch (message_p->procedureCode) {
    case S1ap_ProcedureCode_id_S1Setup:
      return s1ap_mme_encode_s1setupfailure(message_p, buffer, length);
    case S1ap_ProcedureCode_id_PathSwitchRequest:
      return s1ap_mme_encode_pathswitchreqfailure(message_p, buffer, length);

    default:
      OAILOG_DEBUG(
        LOG_S1AP,
        "Unknown procedure ID (%d) for unsuccessfull outcome message\n",
        (int) message_p->procedureCode);
      break;
  }

  return -1;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_s1setupresponse(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_S1SetupResponse_t s1SetupResponse;
  S1ap_S1SetupResponse_t *s1SetupResponse_p = &s1SetupResponse;

  memset(s1SetupResponse_p, 0, sizeof(S1ap_S1SetupResponse_t));

  if (
    s1ap_encode_s1ap_s1setupresponseies(
      s1SetupResponse_p, &message_p->msg.s1ap_S1SetupResponseIEs) < 0) {
    return -1;
  }

  return s1ap_generate_successfull_outcome(
    buffer,
    length,
    S1ap_ProcedureCode_id_S1Setup,
    message_p->criticality,
    &asn_DEF_S1ap_S1SetupResponse,
    s1SetupResponse_p);
}

static inline int s1ap_mme_encode_resetack(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_ResetAcknowledge_t s1ResetAck;
  S1ap_ResetAcknowledge_t *s1ResetAck_p = &s1ResetAck;

  memset(s1ResetAck_p, 0, sizeof(S1ap_ResetAcknowledge_t));

  if (
    s1ap_encode_s1ap_resetacknowledgeies(
      s1ResetAck_p, &message_p->msg.s1ap_ResetAcknowledgeIEs) < 0) {
    return -1;
  }

  return s1ap_generate_successfull_outcome(
    buffer,
    length,
    S1ap_ProcedureCode_id_Reset,
    message_p->criticality,
    &asn_DEF_S1ap_ResetAcknowledge,
    s1ResetAck_p);
}

static inline int s1ap_mme_encode_pathswitchreqack(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_PathSwitchRequestAcknowledge_t s1PathSwitchRequestAck;
  S1ap_PathSwitchRequestAcknowledge_t *s1PathSwitchRequestAck_p =
    &s1PathSwitchRequestAck;

  memset(s1PathSwitchRequestAck_p, 0,
    sizeof(S1ap_PathSwitchRequestAcknowledge_t));

  if (
    s1ap_encode_s1ap_pathswitchrequestacknowledgeies(
      s1PathSwitchRequestAck_p,
      &message_p->msg.s1ap_PathSwitchRequestAcknowledgeIEs) < 0) {
    return -1;
  }

  return s1ap_generate_successfull_outcome(
    buffer,
    length,
    S1ap_ProcedureCode_id_PathSwitchRequest,
    message_p->criticality,
    &asn_DEF_S1ap_PathSwitchRequestAcknowledge,
    s1PathSwitchRequestAck_p);
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_s1setupfailure(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_S1SetupFailure_t s1SetupFailure;
  S1ap_S1SetupFailure_t *s1SetupFailure_p = &s1SetupFailure;

  memset(s1SetupFailure_p, 0, sizeof(S1ap_S1SetupFailure_t));

  if (
    s1ap_encode_s1ap_s1setupfailureies(
      s1SetupFailure_p, &message_p->msg.s1ap_S1SetupFailureIEs) < 0) {
    return -1;
  }

  return s1ap_generate_unsuccessfull_outcome(
    buffer,
    length,
    S1ap_ProcedureCode_id_S1Setup,
    message_p->criticality,
    &asn_DEF_S1ap_S1SetupFailure,
    s1SetupFailure_p);
}

static inline int s1ap_mme_encode_pathswitchreqfailure(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_PathSwitchRequestFailure_t s1PathSwitchRequestFailure;
  S1ap_PathSwitchRequestFailure_t *s1ap_PathSwitchRequestFailure_p =
    &s1PathSwitchRequestFailure;

  memset(s1ap_PathSwitchRequestFailure_p, 0,
        sizeof(S1ap_PathSwitchRequestFailure_t));

  if (
    s1ap_encode_s1ap_pathswitchrequestfailureies(
      s1ap_PathSwitchRequestFailure_p,
      &message_p->msg.s1ap_PathSwitchRequestFailureIEs) < 0) {
    return -1;
  }

  return s1ap_generate_unsuccessfull_outcome(
    buffer,
    length,
    S1ap_ProcedureCode_id_PathSwitchRequest,
    message_p->criticality,
    &asn_DEF_S1ap_PathSwitchRequestFailure,
    s1ap_PathSwitchRequestFailure_p);
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_downlink_nas_transport(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_DownlinkNASTransport_t downlinkNasTransport;
  S1ap_DownlinkNASTransport_t *downlinkNasTransport_p = &downlinkNasTransport;

  memset(downlinkNasTransport_p, 0, sizeof(S1ap_DownlinkNASTransport_t));

  /*
   * Convert IE structure into asn1 message_p
   */
  if (
    s1ap_encode_s1ap_downlinknastransporties(
      downlinkNasTransport_p, &message_p->msg.s1ap_DownlinkNASTransportIEs) <
    0) {
    return -1;
  }

  /*
   * Generate Initiating message_p for the list of IEs
   */
  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_downlinkNASTransport,
    S1ap_Criticality_ignore,
    &asn_DEF_S1ap_DownlinkNASTransport,
    downlinkNasTransport_p);
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_ue_context_release_command(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_UEContextReleaseCommand_t ueContextReleaseCommand;
  S1ap_UEContextReleaseCommand_t *ueContextReleaseCommand_p =
    &ueContextReleaseCommand;

  memset(ueContextReleaseCommand_p, 0, sizeof(S1ap_UEContextReleaseCommand_t));

  /*
   * Convert IE structure into asn1 message_p
   */
  if (
    s1ap_encode_s1ap_uecontextreleasecommandies(
      ueContextReleaseCommand_p,
      &message_p->msg.s1ap_UEContextReleaseCommandIEs) < 0) {
    return -1;
  }

  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_UEContextRelease,
    S1ap_Criticality_reject,
    &asn_DEF_S1ap_UEContextReleaseCommand,
    ueContextReleaseCommand_p);
}

static inline int s1ap_mme_encode_paging(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_Paging_t paging;
  S1ap_Paging_t *paging_p = &paging;
  memset(paging_p, 0, sizeof(S1ap_Paging_t));
  if (
    s1ap_encode_s1ap_pagingies(paging_p, &message_p->msg.s1ap_PagingIEs) < 0) {
    return -1;
  }
  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_Paging,
    S1ap_Criticality_ignore,
    &asn_DEF_S1ap_Paging,
    paging_p);
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_e_rab_setup(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_E_RABSetupRequest_t e_rab_setup;
  S1ap_E_RABSetupRequest_t *e_rab_setup_p = &e_rab_setup;

  memset(e_rab_setup_p, 0, sizeof(S1ap_E_RABSetupRequest_t));

  /*
   * Convert IE structure into asn1 message_p
   */
  if (
    s1ap_encode_s1ap_e_rabsetuprequesties(
      e_rab_setup_p, &message_p->msg.s1ap_E_RABSetupRequestIEs) < 0) {
    return -1;
  }
  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_E_RABSetup,
    message_p->criticality,
    &asn_DEF_S1ap_E_RABSetupRequest,
    e_rab_setup_p);
}

static inline int s1ap_mme_encode_ue_context_modification_request(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_UEContextModificationRequest_t ueContextModificationRequest;
  S1ap_UEContextModificationRequest_t *ueContextModificationRequest_p =
    &ueContextModificationRequest;

  memset(
    ueContextModificationRequest_p,
    0,
    sizeof(S1ap_UEContextModificationRequest_t));

  if (
    s1ap_encode_s1ap_uecontextmodificationrequesties(
      ueContextModificationRequest_p,
      &message_p->msg.s1ap_UEContextModificationRequestIEs) < 0) {
    return -1;
  }

  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_UEContextModification,
    S1ap_Criticality_reject,
    &asn_DEF_S1ap_UEContextModificationRequest,
    ueContextModificationRequest_p);
}

static inline int s1ap_mme_encode_e_rab_release_command(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{
  S1ap_E_RABReleaseCommand_t e_rab_rel_cmd;
  S1ap_E_RABReleaseCommand_t *e_rab_rel_cmd_p = &e_rab_rel_cmd;

  memset(e_rab_rel_cmd_p, 0, sizeof(S1ap_E_RABReleaseCommand_t));

  /*
   * Convert IE structure into asn1 message_p
   */
  if (
    s1ap_encode_s1ap_e_rabreleasecommandies(
      e_rab_rel_cmd_p, &message_p->msg.s1ap_E_RABReleaseCommandIEs) < 0) {
    return -1;
  }
  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_E_RABRelease,
    message_p->criticality,
    &asn_DEF_S1ap_E_RABSetupRequest,
    e_rab_rel_cmd_p);
}

static inline int s1ap_mme_encode_mme_configuration_transfer(
  s1ap_message *message_p,
  uint8_t **buffer,
  uint32_t *length)
{

  S1ap_MMEConfigurationTransfer_t mmeConfigurationTransfer;
  S1ap_MMEConfigurationTransfer_t *mmeConfigurationTransfer_p =
    &mmeConfigurationTransfer;

  memset(
    mmeConfigurationTransfer_p,
    0,
    sizeof(S1ap_MMEConfigurationTransfer_t));

  if (
    s1ap_encode_s1ap_mmeconfigurationtransferies(
      mmeConfigurationTransfer_p,
      &message_p->msg.s1ap_MMEConfigurationTransferIEs) < 0) {
    return -1;
  }

  return s1ap_generate_initiating_message(
    buffer,
    length,
    S1ap_ProcedureCode_id_MMEConfigurationTransfer,
    S1ap_Criticality_ignore,
    &asn_DEF_S1ap_MMEConfigurationTransfer,
    mmeConfigurationTransfer_p);
}
