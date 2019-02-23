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

/*! \file s11_ie_formatter.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <inttypes.h>
#include <pthread.h>

#include "bstrlib.h"

#include "dynamic_memory_check.h"
#include "common_defs.h"
#include "gcc_diag.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "3gpp_33.401.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_24.007.h"
#include "3gpp_29.274.h"
#include "3gpp_36.413.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgParser.h"
#include "s11_common.h"
#include "security_types.h"
#include "common_types.h"
#include "sgw_ie_defs.h"
#include "PdnType.h"
#include "s11_ie_formatter.h"

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_imsi_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  imsi_t *imsi = (imsi_t *) arg;
  uint8_t decoded = 0;

  DevAssert(arg);
  DevAssert(ieLength <= IMSI_BCD_DIGITS_MAX);
  imsi->length = 0;
  while (decoded < ieLength) {
    uint8_t tmp = ieValue[decoded++];
    imsi->u.value[imsi->length++] = (tmp >> 4) | (tmp << 4);
  }

  for (int i = imsi->length; i < IMSI_BCD8_SIZE; i++) {
    imsi->u.value[i] = 0xff;
  }

  OAILOG_DEBUG(
    LOG_S11,
    "\t- IMSI (l=%d) %u%u%u%u%u%u%u%u%u%u%u%u%u%u%u\n",
    imsi->length,
    imsi->u.num.digit1,
    imsi->u.num.digit2,
    imsi->u.num.digit3,
    imsi->u.num.digit4,
    imsi->u.num.digit5,
    imsi->u.num.digit6,
    imsi->u.num.digit7,
    imsi->u.num.digit8,
    imsi->u.num.digit9,
    imsi->u.num.digit10,
    imsi->u.num.digit11,
    imsi->u.num.digit12,
    imsi->u.num.digit13,
    imsi->u.num.digit14,
    imsi->u.num.digit15);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_imsi_ie_set(nw_gtpv2c_msg_handle_t *msg, const imsi_t *imsi)
{
  nw_rc_t rc;
  imsi_t imsi_nbo = {0};

  DevAssert(msg);
  DevAssert(imsi);
  memcpy(&imsi_nbo, imsi, sizeof(imsi_nbo));
  for (int i = 0; i < IMSI_BCD8_SIZE; i++) {
    uint8_t tmp = imsi_nbo.u.value[i];
    imsi_nbo.u.value[i] = (tmp >> 4) | (tmp << 4);
  }

  rc = nwGtpv2cMsgAddIe(
    *msg, NW_GTPV2C_IE_IMSI, imsi_nbo.length, 0, (uint8_t *) imsi_nbo.u.value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_msisdn_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  Msisdn_t *msisdn;
  uint8_t i;
  uint8_t mask = 0x0F;
  uint8_t msisdn_length = 2 * ieLength;

  DevAssert(arg);
  msisdn = (Msisdn_t *) arg;

  for (i = 0; i < ieLength * 2; i++) {
    if (mask == 0x0F) {
      msisdn->digit[i] = (ieValue[i / 2] & (mask));
    } else {
      msisdn->digit[i] = (ieValue[i / 2] & (mask)) >> 4;
    }

    msisdn->digit[i] += '0';
    mask = ~mask;
  }

  if (msisdn->digit[msisdn_length - 1] == (0x0f + '0')) {
    msisdn->digit[msisdn_length - 1] = 0;
    msisdn_length--;
  }

  msisdn->length = msisdn_length;
  OAILOG_DEBUG(LOG_S11, "\t- MSISDN length %d\n", msisdn->length);
  OAILOG_DEBUG(
    LOG_S11, "\t-        value  %*s\n", msisdn->length, (char *) msisdn->digit);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_mei_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  Mei_t *mei = (Mei_t *) arg;

  DevAssert(mei);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_node_type_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  node_type_t *node_type = (node_type_t *) arg;

  DevAssert(node_type);

  if (*ieValue == 0) {
    *node_type = NODE_TYPE_MME;
  } else if (*ieValue == 1) {
    *node_type = NODE_TYPE_SGSN;
  } else {
    OAILOG_ERROR(
      LOG_S11, "Received unknown value for Node Type: %u\n", *ieValue);
    return NW_GTPV2C_IE_INCORRECT;
  }

  OAILOG_DEBUG(LOG_S11, "\t- Node type %u\n", *node_type);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_node_type_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const node_type_t *node_type)
{
  nw_rc_t rc;
  uint8_t value;

  DevAssert(node_type);
  DevAssert(msg);

  switch (*node_type) {
    case NODE_TYPE_MME: value = 0; break;

    case NODE_TYPE_SGSN: value = 1; break;

    default:
      OAILOG_ERROR(LOG_S11, "Invalid Node type received: %d\n", *node_type);
      return RETURNerror;
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_NODE_TYPE, 1, 0, (uint8_t *) &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_pdn_type_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  pdn_type_t *pdn_type = (pdn_type_t *) arg;

  DevAssert(pdn_type);

  if (*ieValue == 1) {
    /*
     * Only IPv4
     */
    *pdn_type = IPv4;
  } else if (*ieValue == 2) {
    /*
     * Only IPv6
     */
    *pdn_type = IPv6;
  } else if (*ieValue == 3) {
    /*
     * IPv4 and/or IPv6
     */
    *pdn_type = IPv4_AND_v6;
  } else {
    OAILOG_ERROR(
      LOG_S11, "Received unknown value for PDN Type: %u\n", *ieValue);
    return NW_GTPV2C_IE_INCORRECT;
  }

  OAILOG_DEBUG(LOG_S11, "\t- PDN type %u\n", *pdn_type);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_pdn_type_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const pdn_type_t *pdn_type)
{
  nw_rc_t rc;
  uint8_t value;

  DevAssert(pdn_type);
  DevAssert(msg);

  switch (*pdn_type) {
    case IPv4: value = 1; break;

    case IPv6: value = 2; break;

    case IPv4_AND_v6:
    case IPv4_OR_v6: value = 3; break;

    default:
      OAILOG_ERROR(LOG_S11, "Invalid PDN type received: %d\n", *pdn_type);
      return RETURNerror;
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_PDN_TYPE, 1, 0, (uint8_t *) &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_rat_type_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  rat_type_t *rat_type = (rat_type_t *) arg;

  DevAssert(rat_type);

  switch (*ieValue) {
    case 1: *rat_type = RAT_UTRAN; break;

    case 2: *rat_type = RAT_GERAN; break;

    case 3: *rat_type = RAT_WLAN; break;

    case 4: *rat_type = RAT_GAN; break;

    case 5: *rat_type = RAT_HSPA_EVOLUTION; break;

    case 6: *rat_type = RAT_EUTRAN; break;

    default:
      OAILOG_ERROR(
        LOG_S11,
        "Can't map GTP RAT type %u to EPC definition\n"
        "\tCheck TS.29.274 #8.17 for possible values\n",
        *ieValue);
      return NW_GTPV2C_IE_INCORRECT;
  }

  OAILOG_DEBUG(LOG_S11, "\t- RAT type (%d): %d\n", *ieValue, *rat_type);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_rat_type_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const rat_type_t *rat_type)
{
  nw_rc_t rc;
  uint8_t value;

  DevAssert(rat_type);
  DevAssert(msg);

  switch (*rat_type) {
    case RAT_UTRAN: value = 1; break;

    case RAT_GERAN: value = 2; break;

    case RAT_WLAN: value = 3; break;

    case RAT_GAN: value = 4; break;

    case RAT_HSPA_EVOLUTION: value = 5; break;

    case RAT_EUTRAN: value = 6; break;

    default:
      OAILOG_ERROR(
        LOG_S11,
        "Can't map RAT type %d to GTP RAT type\n"
        "\tCheck TS.29.274 #8.17 for possible values\n",
        *rat_type);
      return RETURNerror;
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_RAT_TYPE, 1, 0, (uint8_t *) &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int gtpv2c_ebi_ie_set(nw_gtpv2c_msg_handle_t *msg, const unsigned ebi)
{
  nw_rc_t rc;
  uint8_t value = 0;

  value = ebi & 0x0F;
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_EBI, 1, 0, &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ebi_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t *ebi = (uint8_t *) arg;

  DevAssert(ebi);
  *ebi = ieValue[0] & 0x0F;
  OAILOG_DEBUG(LOG_S11, "\t- EBI %u\n", *ebi);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ebi_ie_get_list(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  ebi_list_t *ebi_list = (ebi_list_t *) arg;
  DevAssert(ebi_list);
  DevAssert(RELEASE_ACCESS_BEARER_MAX_BEARERS > ebi_list->num_ebi);
  uint8_t *ebi = (uint8_t *) &ebi_list->ebis[ebi_list->num_ebi];

  DevAssert(ebi);
  *ebi = ieValue[0] & 0x0F;
  OAILOG_DEBUG(LOG_S11, "\t- EBI %u\n", *ebi);
  ebi_list->num_ebi += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_cause_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  gtpv2c_cause_t *cause = (gtpv2c_cause_t *) arg;

  DevAssert(cause);
  cause->cause_value = ieValue[0];
  cause->cs = ieValue[1] & 0x01;
  cause->bce = (ieValue[1] & 0x02) >> 1;
  cause->pce = (ieValue[1] & 0x04) >> 2;
  if (6 == ieLength) {
    cause->offending_ie_type = ieValue[2];
    cause->offending_ie_length = ((uint16_t) ieValue[3]) << 8;
    cause->offending_ie_length |= ((uint16_t) ieValue[4]);
    cause->offending_ie_instance = ieValue[5] & 0x0F;
  }
  OAILOG_DEBUG(LOG_S11, "\t- Cause value %u\n", cause->cause_value);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_cause_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const gtpv2c_cause_t *cause)
{
  nw_rc_t rc;
  uint8_t value[6];

  DevAssert(msg);
  DevAssert(cause);
  value[0] = cause->cause_value;
  value[1] =
    ((cause->pce & 0x1) << 2) | ((cause->bce & 0x1) << 1) | (cause->cs & 0x1);

  if (cause->offending_ie_type) {
    value[2] = cause->offending_ie_type;
    value[3] = (cause->offending_ie_length & 0xFF00) >> 8;
    value[4] = cause->offending_ie_length & 0x00FF;
    value[5] = cause->offending_ie_instance & 0x0F;
    rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_CAUSE, 6, 0, value);
  } else {
    rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_CAUSE, 2, 0, value);
  }

  DevAssert(NW_OK == rc);
  return rc == NW_OK ? 0 : -1;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_to_create_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_to_create_t *bearer_to_create)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_to_create);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer_to_create->eps_bearer_id);
  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t
gtpv2c_bearer_context_to_be_created_within_create_session_request_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_contexts_to_be_created_t *bearer_contexts =
    (bearer_contexts_to_be_created_t *) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
    MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS >=
    bearer_contexts->num_bearer_context);
  bearer_context_to_be_created_t *bearer_context =
    &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint8_t read = 0;
  nw_rc_t rc;

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t *ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t *) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_BEARER_LEVEL_QOS:
        rc = gtpv2c_bearer_qos_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->bearer_level_qos);
        break;

      case NW_GTPV2C_IE_BEARER_TFT:
        OAILOG_ERROR(LOG_S11, "Received IE %u to implement\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
        break;

      case NW_GTPV2C_IE_FTEID:
        switch (ie_p->i) {
          case 0:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s1u_enb_fteid);
            break;
          case 1:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s4u_sgsn_fteid);
            break;
          case 2:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s5_s8_u_sgw_fteid);
            break;
          case 3:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s5_s8_u_pgw_fteid);
            break;
          case 4:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s12_rnc_fteid);
            break;
          case 5:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s2b_u_epdg_fteid);
            break;
          default:
            OAILOG_ERROR(
              LOG_S11,
              "Received unexpected IE %u instance %u\n",
              ie_p->t,
              ie_p->i);
            return NW_GTPV2C_IE_INCORRECT;
        }
        DevAssert(NW_OK == rc);
        break;

      default:
        OAILOG_ERROR(LOG_S11, "Received unexpected IE %u\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
    }

    read += (ntohs(ie_p->l) + sizeof(nw_gtpv2c_ie_tlv_t));
  }
  bearer_contexts->num_bearer_context += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_to_be_created_within_create_session_request_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_context_to_be_created_t *bearer_context)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_context);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer_context->eps_bearer_id);
  gtpv2c_bearer_qos_ie_set(msg, &bearer_context->bearer_level_qos);
  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_contexts_within_create_bearer_request_t *bearer_contexts =
    (bearer_contexts_within_create_bearer_request_t *) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
    MSG_CREATE_BEARER_REQUEST_MAX_BEARER_CONTEXTS >=
    bearer_contexts->num_bearer_context);
  bearer_context_within_create_bearer_request_t *bearer_context =
    &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint8_t read = 0;
  nw_rc_t rc;

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t *ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t *) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_BEARER_LEVEL_QOS:
        rc = gtpv2c_bearer_qos_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->bearer_level_qos);
        break;

      case NW_GTPV2C_IE_BEARER_TFT:
        rc = gtpv2c_tft_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->tft);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_PCO:
        rc = gtpv2c_pco_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->pco);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_FTEID:
        switch (ie_p->i) {
          case 0:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s1u_sgw_fteid);
            break;
          case 1:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s5_s8_u_pgw_fteid);
            break;
          case 2:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s12_sgw_fteid);
            break;
          case 3:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s4_u_sgw_fteid);
            break;
          case 4:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s2b_u_pgw_fteid);
            break;
          default:
            OAILOG_ERROR(
              LOG_S11,
              "Received unexpected IE %u instance %u\n",
              ie_p->t,
              ie_p->i);
            return NW_GTPV2C_IE_INCORRECT;
        }
        DevAssert(NW_OK == rc);
        break;

      default:
        OAILOG_ERROR(LOG_S11, "Received unexpected IE %u\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
    }

    read += (ntohs(ie_p->l) + sizeof(nw_gtpv2c_ie_tlv_t));
  }
  bearer_contexts->num_bearer_context += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_context_within_create_bearer_request_t *bearer_context)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_context);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer_context->eps_bearer_id);
  if (bearer_context->pco.num_protocol_or_container_id) {
    gtpv2c_pco_ie_set(msg, &bearer_context->pco);
  }
  if (bearer_context->s1u_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s1u_sgw_fteid, 0);
  }
  if (bearer_context->s5_s8_u_pgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s5_s8_u_pgw_fteid, 1);
  }
  if (bearer_context->s12_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s12_sgw_fteid, 2);
  }
  if (bearer_context->s4_u_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s4_u_sgw_fteid, 3);
  }
  if (bearer_context->s2b_u_pgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s2b_u_pgw_fteid, 4);
  }
  gtpv2c_bearer_qos_ie_set(msg, &bearer_context->bearer_level_qos);
  gtpv2c_tft_ie_set(msg, &bearer_context->tft);

  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_within_create_bearer_response_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_context_within_create_bearer_response_t *bearer_context)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_context);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer_context->eps_bearer_id);
  gtpv2c_cause_ie_set(msg, &bearer_context->cause);
  gtpv2c_fteid_ie_set(msg, &bearer_context->s1u_enb_fteid, 0);
  gtpv2c_fteid_ie_set(msg, &bearer_context->s1u_sgw_fteid, 1);
  if (bearer_context->s5_s8_u_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s5_s8_u_sgw_fteid, 2);
  }
  if (bearer_context->s5_s8_u_pgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s5_s8_u_pgw_fteid, 3);
  }
  if (bearer_context->s12_rnc_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s12_rnc_fteid, 4);
  }
  if (bearer_context->s12_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s12_sgw_fteid, 5);
  }
  if (bearer_context->s4_u_sgsn_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s4_u_sgsn_fteid, 6);
  }
  if (bearer_context->s4_u_sgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s4_u_sgw_fteid, 7);
  }
  if (bearer_context->s2b_u_epdg_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s2b_u_epdg_fteid, 8);
  }
  if (bearer_context->s2b_u_pgw_fteid.teid) {
    gtpv2c_fteid_ie_set(msg, &bearer_context->s2b_u_pgw_fteid, 9);
  }
  if (bearer_context->pco.num_protocol_or_container_id) {
    gtpv2c_pco_ie_set(msg, &bearer_context->pco);
  }

  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_context_within_create_bearer_response_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_contexts_within_create_bearer_response_t *bearer_contexts =
    (bearer_contexts_within_create_bearer_response_t *) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
    MSG_MODIFY_BEARER_REQUEST_MAX_BEARER_CONTEXTS >=
    bearer_contexts->num_bearer_context);
  bearer_context_within_create_bearer_response_t *bearer_context =
    &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint8_t read = 0;
  nw_rc_t rc;

  DevAssert(bearer_context);

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t *ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t *) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_CAUSE:
        rc = gtpv2c_cause_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->cause);
        break;

      case NW_GTPV2C_IE_PCO:
        rc = gtpv2c_pco_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->pco);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_FTEID:
        switch (ie_p->i) {
          case 0:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s1u_enb_fteid);
            break;
          case 1:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s1u_sgw_fteid);
            break;
          case 2:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s5_s8_u_sgw_fteid);
            break;
          case 3:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s5_s8_u_pgw_fteid);
            break;
          case 4:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s12_rnc_fteid);
            break;
          case 5:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s12_sgw_fteid);
            break;
          case 6:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s4_u_sgsn_fteid);
            break;
          case 7:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s4_u_sgw_fteid);
            break;
          case 8:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s2b_u_epdg_fteid);
            break;
          case 9:
            rc = gtpv2c_fteid_ie_get(
              ie_p->t,
              ntohs(ie_p->l),
              ie_p->i,
              &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
              &bearer_context->s2b_u_pgw_fteid);
            break;

          default:
            OAILOG_ERROR(
              LOG_S11, "Received unexpected instance %u for fteid\n", ie_p->i);
        }
        break;

      default:
        OAILOG_ERROR(LOG_S11, "Received unexpected IE %u\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
    }

    read += (ntohs(ie_p->l) + sizeof(nw_gtpv2c_ie_tlv_t));
  }
  bearer_contexts->num_bearer_context += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_context_to_be_modified_t *bearer_context)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_context);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer_context->eps_bearer_id);
  gtpv2c_fteid_ie_set(msg, &bearer_context->s1_eNB_fteid, 0);
  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t
gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_contexts_to_be_modified_t *bearer_contexts =
    (bearer_contexts_to_be_modified_t *) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
    MSG_MODIFY_BEARER_REQUEST_MAX_BEARER_CONTEXTS >=
    bearer_contexts->num_bearer_context);
  bearer_context_to_be_modified_t *bearer_context =
    &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint8_t read = 0;
  nw_rc_t rc;

  DevAssert(bearer_context);

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t *ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t *) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_FTEID:
        rc = gtpv2c_fteid_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->s1_eNB_fteid);
        break;

      default:
        OAILOG_ERROR(LOG_S11, "Received unexpected IE %u\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
    }

    read += (ntohs(ie_p->l) + sizeof(nw_gtpv2c_ie_tlv_t));
  }
  bearer_contexts->num_bearer_context += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_context_created_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_contexts_created_t *bearer_contexts =
    (bearer_contexts_created_t *) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
    MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS >=
    bearer_contexts->num_bearer_context);
  bearer_context_created_t *bearer_context =
    &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint8_t read = 0;
  nw_rc_t rc;

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t *ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t *) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_FTEID:
        rc = gtpv2c_fteid_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->s1u_sgw_fteid);
        break;

      case NW_GTPV2C_IE_CAUSE:
        rc = gtpv2c_cause_ie_get(
          ie_p->t,
          ntohs(ie_p->l),
          ie_p->i,
          &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
          &bearer_context->cause);
        break;

      default:
        OAILOG_ERROR(LOG_S11, "Received unexpected IE %u\n", ie_p->t);
        return NW_GTPV2C_IE_INCORRECT;
    }

    read += (ntohs(ie_p->l) + sizeof(nw_gtpv2c_ie_tlv_t));
  }
  bearer_contexts->num_bearer_context += 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_context_created_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_context_created_t const *bearer)
{
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer);
  /*
   * Start section for grouped IE: bearer context created
   */
  rc = nwGtpv2cMsgGroupedIeStart(
    *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer->eps_bearer_id);
  rc = gtpv2c_cause_ie_set(msg, &bearer->cause);
  DevAssert(NW_OK == rc);
  rc = nwGtpv2cMsgAddIeFteid(
    *msg,
    NW_GTPV2C_IE_INSTANCE_ZERO,
    bearer->s1u_sgw_fteid.interface_type,
    bearer->s1u_sgw_fteid.teid,
    bearer->s1u_sgw_fteid.ipv4 ? &bearer->s1u_sgw_fteid.ipv4_address : 0,
    bearer->s1u_sgw_fteid.ipv6 ? &bearer->s1u_sgw_fteid.ipv6_address : NULL);
  DevAssert(NW_OK == rc);
  /*
   * End section for grouped IE: bearer context created
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_apn_restriction_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  APNRestriction_t *apn_restriction = (APNRestriction_t *) arg;
  OAILOG_DEBUG(LOG_S11, "\t- APN restriction 0x%02x\n", *apn_restriction);
  return NW_OK;
}

//------------------------------------------------------------------------------
/* This IE shall be included in the E-UTRAN initial attach,
   PDP Context Activation and UE Requested PDN connectivity procedures.
   This IE denotes the most stringent restriction as required
   by any already active bearer context. If there are no already active bearer
   contexts, this value is set to the least restrictive type.
*/
int gtpv2c_apn_restriction_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const uint8_t apn_restriction)
{
  nw_rc_t rc;

  DevAssert(msg);
  rc = nwGtpv2cMsgAddIe(
    *msg, NW_GTPV2C_IE_APN_RESTRICTION, 1, 0, (uint8_t *) &apn_restriction);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_serving_network_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  ServingNetwork_t *serving_net = (ServingNetwork_t *) arg;

  DevAssert(serving_net);
  serving_net->mcc[1] = (ieValue[0] & 0xF0) >> 4;
  serving_net->mcc[0] = (ieValue[0] & 0x0F);
  serving_net->mcc[2] = (ieValue[1] & 0x0F);

  if ((ieValue[1] & 0xF0) == 0xF0) {
    /*
     * Two digits MNC
     */
    serving_net->mnc[0] = 0;
    serving_net->mnc[1] = (ieValue[2] & 0x0F);
    serving_net->mnc[2] = (ieValue[2] & 0xF0) >> 4;
  } else {
    serving_net->mnc[0] = (ieValue[2] & 0x0F);
    serving_net->mnc[1] = (ieValue[2] & 0xF0) >> 4;
    serving_net->mnc[2] = (ieValue[1] & 0xF0) >> 4;
  }

  OAILOG_DEBUG(
    LOG_S11,
    "\t- Serving network %d.%d\n",
    serving_net->mcc[0] * 100 + serving_net->mcc[1] * 10 + serving_net->mcc[2],
    serving_net->mnc[0] * 100 + serving_net->mnc[1] * 10 + serving_net->mnc[2]);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_serving_network_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const ServingNetwork_t *serving_network)
{
  nw_rc_t rc;
  uint8_t value[3];

  DevAssert(msg);
  DevAssert(serving_network);
  /*
   * MCC Decimal | MCC Hundreds
   */
  value[0] =
    ((serving_network->mcc[1] & 0x0F) << 4) | (serving_network->mcc[0] & 0x0F);
  value[1] = serving_network->mcc[2] & 0x0F;

  if ((serving_network->mnc[0] & 0xF) == 0xF) {
    /*
     * Only two digits
     */
    value[1] |= 0xF0;
    value[2] = ((serving_network->mnc[2] & 0x0F) << 4) |
               (serving_network->mnc[1] & 0x0F);
  } else {
    value[1] |= (serving_network->mnc[2] & 0x0F) << 4;
    value[2] = ((serving_network->mnc[1] & 0x0F) << 4) |
               (serving_network->mnc[0] & 0x0F);
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_SERVING_NETWORK, 3, 0, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
int gtpv2c_fteid_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const fteid_t *fteid,
  const uint8_t instance)
{
  nw_rc_t rc;
  uint8_t value[25];

  DevAssert(msg);
  DevAssert(fteid);
  /*
   * MCC Decimal | MCC Hundreds
   */
  value[0] =
    (fteid->ipv4 << 7) | (fteid->ipv6 << 6) | (fteid->interface_type & 0x3F);
  value[1] = (fteid->teid >> 24);
  value[2] = (fteid->teid >> 16) & 0xFF;
  value[3] = (fteid->teid >> 8) & 0xFF;
  value[4] = (fteid->teid >> 0) & 0xFF;

  int offset = 5;
  if (fteid->ipv4 == 1) {
    uint32_t hbo = ntohl(fteid->ipv4_address.s_addr);
    value[offset++] = (uint8_t)(hbo >> 24);
    value[offset++] = (uint8_t)(hbo >> 16);
    value[offset++] = (uint8_t)(hbo >> 8);
    value[offset++] = (uint8_t) hbo;
  }
  if (fteid->ipv6 == 1) {
    /*
     * IPv6 present: copy the 16 bytes
     */
    memcpy(&value[offset], fteid->ipv6_address.__in6_u.__u6_addr8, 16);
    offset += 16;
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_FTEID, offset, instance, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_fteid_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t offset = 0;
  fteid_t *fteid = (fteid_t *) arg;

  DevAssert(fteid);
  fteid->ipv4 = (ieValue[0] & 0x80) >> 7;
  fteid->ipv6 = (ieValue[0] & 0x40) >> 6;
  fteid->interface_type = ieValue[0] & 0x1F;
  OAILOG_DEBUG(LOG_S11, "\t- F-TEID type %d\n", fteid->interface_type);
  /*
   * Copy the TEID or GRE key
   */
  fteid->teid = ntoh_int32_buf(&ieValue[1]);
  OAILOG_DEBUG(LOG_S11, "\t- TEID/GRE    %08x\n", fteid->teid);

  if (fteid->ipv4 == 1) {
    /*
     * IPv4 present: copy the 4 bytes
     */
    uint32_t hbo = (((uint32_t) ieValue[5]) << 24) |
                   (((uint32_t) ieValue[6]) << 16) |
                   (((uint32_t) ieValue[7]) << 8) | (uint32_t) ieValue[8];
    fteid->ipv4_address.s_addr = htonl(hbo);
    offset = 4;
    OAILOG_DEBUG(
      LOG_S11,
      "\t- IPv4 addr   " IN_ADDR_FMT "\n",
      PRI_IN_ADDR(fteid->ipv4_address));
  }

  if (fteid->ipv6 == 1) {
    char ipv6_ascii[INET6_ADDRSTRLEN];

    /*
     * IPv6 present: copy the 16 bytes
     * * * * WARNING: if Ipv4 is present, 4 bytes of offset should be applied
     */
    memcpy(fteid->ipv6_address.__in6_u.__u6_addr8, &ieValue[5 + offset], 16);
    inet_ntop(
      AF_INET6, (void *) &fteid->ipv6_address, ipv6_ascii, INET6_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_S11, "\t- IPv6 addr   %s\n", ipv6_ascii);
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_pco_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t offset = 0;
  protocol_configuration_options_t *pco =
    (protocol_configuration_options_t *) arg;

  DevAssert(pco);
  offset = decode_protocol_configuration_options(pco, ieValue, ieLength);
  if ((0 < offset) && (PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH >= offset))
    return NW_OK;
  else
    return NW_GTPV2C_IE_INCORRECT;
}

//------------------------------------------------------------------------------
int gtpv2c_pco_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const protocol_configuration_options_t *pco)
{
  uint8_t temp[PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH];
  uint8_t offset = 0;
  nw_rc_t rc = NW_OK;

  DevAssert(pco);
  offset = encode_protocol_configuration_options(
    pco, temp, PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH);
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_PCO, offset, 0, temp);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_tft_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t offset = 0;
  traffic_flow_template_t *tft = (traffic_flow_template_t *) arg;

  DevAssert(tft);
  offset = decode_traffic_flow_template(tft, ieValue, ieLength);
  if ((0 < offset) && (TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH >= offset))
    return NW_OK;
  else {
    OAILOG_ERROR(LOG_S11, "Incorrect TFT IE\n");
    return NW_GTPV2C_IE_INCORRECT;
  }
}

//------------------------------------------------------------------------------
int gtpv2c_tft_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const traffic_flow_template_t *tft)
{
  uint8_t temp[TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH];
  uint8_t offset = 0;
  nw_rc_t rc = NW_OK;

  DevAssert(tft);
  offset = encode_traffic_flow_template(
    tft, temp, TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH);
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_BEARER_TFT, offset, 0, temp);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_paa_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t offset = 0;
  paa_t *paa = (paa_t *) arg;

  DevAssert(paa);
  paa->pdn_type = ieValue[0] & 0x07;
  OAILOG_DEBUG(LOG_S11, "\t- PAA type  %d\n", paa->pdn_type);

  if (paa->pdn_type & 0x2) {
    char ipv6_ascii[INET6_ADDRSTRLEN];

    /*
     * IPv6 present: copy the 16 bytes
     * * * * WARNING: if both ipv4 and ipv6 are present,
     * * * *          17 bytes of offset should be applied for ipv4
     * * * * NOTE: an ipv6 prefix length is prepend
     * * * * NOTE: in Rel.8 the prefix length has a default value of /64
     */
    paa->ipv6_prefix_length = ieValue[1];
    memcpy(paa->ipv6_address.__in6_u.__u6_addr8, &ieValue[2], 16);
    inet_ntop(AF_INET6, &paa->ipv6_address, ipv6_ascii, INET6_ADDRSTRLEN);
    OAILOG_DEBUG(
      LOG_S11, "\t- IPv6 addr %s/%u\n", ipv6_ascii, paa->ipv6_prefix_length);
  }

  if (paa->pdn_type == 3) {
    offset = 17;
  }

  if (paa->pdn_type & 0x1) {
    uint32_t ip = (((uint32_t) ieValue[1 + offset]) << 24) |
                  (((uint32_t) ieValue[2 + offset]) << 16) |
                  (((uint32_t) ieValue[3 + offset]) << 8) |
                  ((uint32_t) ieValue[4 + offset]);

    paa->ipv4_address.s_addr = htonl(ip);
    char ipv4[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, (void *) &paa->ipv4_address, ipv4, INET_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_S11, "\t- IPv4 addr %s\n", ipv4);
  }

  paa->pdn_type -= 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_paa_ie_set(nw_gtpv2c_msg_handle_t *msg, const paa_t *paa)
{
  /*
   * ipv4 address = 4 + ipv6 address = 16 + ipv6 prefix length = 1
   * * * * + pdn_type = 1
   * * * * = maximum of 22 bytes
   */
  uint8_t temp[22];
  uint8_t pdn_type;
  uint8_t offset = 0;
  nw_rc_t rc;

  DevAssert(paa);
  pdn_type = paa->pdn_type + 1;
  temp[offset] = pdn_type;
  offset++;

  if (pdn_type & 0x2) {
    /*
     * If ipv6 or ipv4v6 present
     */
    temp[1] = paa->ipv6_prefix_length;
    memcpy(&temp[2], paa->ipv6_address.__in6_u.__u6_addr8, 16);
    offset += 17;
  }

  if (pdn_type & 0x1) {
    uint32_t hbo = ntohl(paa->ipv4_address.s_addr);
    temp[offset++] = (uint8_t)(hbo >> 24);
    temp[offset++] = (uint8_t)(hbo >> 16);
    temp[offset++] = (uint8_t)(hbo >> 8);
    temp[offset++] = (uint8_t) hbo;
  }

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_PAA, offset, 0, temp);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
/* The encoding of the APN shall follow the Name Syntax defined in RFC 2181,
   RFC 1035 and RFC 1123. The APN consists of one or more labels. Each label
   is coded as a one octet length field followed by that number of octets
   coded as 8 bit ASCII characters. Following RFC 1035 the labels shall consist
   only of the alphabetic characters (A-Z and a-z), digits (0-9)
   and the hyphen (-).
*/
nw_rc_t gtpv2c_apn_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  uint8_t read = 1;
  uint8_t word_length;
  char *apn = (char *) arg;

  DevAssert(apn);
  DevCheck(
    ieLength <= ACCESS_POINT_NAME_MAX_LENGTH,
    ieLength,
    ACCESS_POINT_NAME_MAX_LENGTH,
    0);
  word_length = ieValue[0];

  while (read < ieLength) {
    if (word_length > 0) {
      apn[read - 1] = ieValue[read];
      word_length--;
    } else {
      /*
       * This is not an alphanumeric character
       */
      word_length = ieValue[read];
      /*
       * Replace the length attribute by '.'
       */
      apn[read - 1] = '.';
    }

    read++;
  }

  apn[read - 1] = '\0';
  OAILOG_DEBUG(LOG_S11, "\t- APN %s\n", apn);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_apn_ie_set(nw_gtpv2c_msg_handle_t *msg, const char *apn)
{
  nw_rc_t rc;
  uint8_t *value;
  uint8_t apn_length;
  uint8_t offset = 0;
  uint8_t *last_size;
  uint8_t word_length = 0;

  DevAssert(apn);
  DevAssert(msg);
  apn_length = strlen(apn);
  value = calloc(apn_length + 1, sizeof(uint8_t));
  last_size = &value[0];

  while (apn[offset]) {
    /*
     * We replace the . by the length of the word
     */
    if (apn[offset] == '.') {
      *last_size = word_length;
      word_length = 0;
      last_size = &value[offset + 1];
    } else {
      word_length++;
      value[offset + 1] = apn[offset];
    }

    offset++;
  }

  *last_size = word_length;
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_APN, apn_length + 1, 0, value);
  DevAssert(NW_OK == rc);
  free_wrapper((void **) &value);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ambr_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  ambr_t *ambr = (ambr_t *) arg;

  DevAssert(ambr);
  ambr->br_ul = ntoh_int32_buf(&ieValue[0]);
  ambr->br_dl = ntoh_int32_buf(&ieValue[4]);
  OAILOG_DEBUG(LOG_S11, "\t- AMBR UL %" PRIu64 "\n", ambr->br_ul);
  OAILOG_DEBUG(LOG_S11, "\t- AMBR DL %" PRIu64 "\n", ambr->br_dl);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_uli_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  Uli_t *uli = (Uli_t *) arg;

  DevAssert(uli);
  uli->present = ieValue[0];

  if (uli->present & ULI_CGI) {
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_qos_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_qos_t *bearer_qos = (bearer_qos_t *) arg;

  DevAssert(bearer_qos);

  if (22 <= ieLength) {
    int index = 0;
    bearer_qos->pci = (ieValue[index] >> 6) & 0x01;
    bearer_qos->pci = (ieValue[index] >> 6) & 0x01;
    bearer_qos->pl = (ieValue[index] >> 2) & 0x0F;
    bearer_qos->pvi = ieValue[index++] & 0x01;
    bearer_qos->qci = ieValue[index++];

    bearer_qos->mbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->mbr.br_ul |= (bit_rate_t) ieValue[index++];

    bearer_qos->mbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->mbr.br_dl |= (bit_rate_t) ieValue[index++];

    bearer_qos->gbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->gbr.br_ul |= (bit_rate_t) ieValue[index++];

    bearer_qos->gbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->gbr.br_dl |= (bit_rate_t) ieValue[index++];

    if (22 < ieLength) {
      OAILOG_ERROR(LOG_S11, "TODO gtpv2c_bearer_qos_ie_get() BearerQOS_t\n");
      return NW_GTPV2C_IE_INCORRECT;
    }
    return NW_OK;
  } else {
    OAILOG_ERROR(LOG_S11, "Bad IE length %" PRIu8 "\n", ieLength);
    return NW_GTPV2C_IE_INCORRECT;
  }
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_qos_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_qos_t *bearer_qos)
{
  nw_rc_t rc;
  uint8_t value[22];
  int index = 0;

  DevAssert(msg);
  DevAssert(bearer_qos);
  value[index++] =
    (bearer_qos->pci << 6) | (bearer_qos->pl << 2) | (bearer_qos->pvi);
  value[index++] = bearer_qos->qci;
  /*
   * TODO: check endianness
   */
  value[index++] = (bearer_qos->mbr.br_ul & 0xFF00000000) >> 32;
  value[index++] = (bearer_qos->mbr.br_ul & 0x00FF000000) >> 24;
  value[index++] = (bearer_qos->mbr.br_ul & 0x0000FF0000) >> 16;
  value[index++] = (bearer_qos->mbr.br_ul & 0x000000FF00) >> 8;
  value[index++] = (bearer_qos->mbr.br_ul & 0x00000000FF);

  value[index++] = (bearer_qos->mbr.br_dl & 0xFF00000000) >> 32;
  value[index++] = (bearer_qos->mbr.br_dl & 0x00FF000000) >> 24;
  value[index++] = (bearer_qos->mbr.br_dl & 0x0000FF0000) >> 16;
  value[index++] = (bearer_qos->mbr.br_dl & 0x000000FF00) >> 8;
  value[index++] = (bearer_qos->mbr.br_dl & 0x00000000FF);

  value[index++] = (bearer_qos->gbr.br_ul & 0xFF00000000) >> 32;
  value[index++] = (bearer_qos->gbr.br_ul & 0x00FF000000) >> 24;
  value[index++] = (bearer_qos->gbr.br_ul & 0x0000FF0000) >> 16;
  value[index++] = (bearer_qos->gbr.br_ul & 0x000000FF00) >> 8;
  value[index++] = (bearer_qos->gbr.br_ul & 0x00000000FF);

  value[index++] = (bearer_qos->gbr.br_dl & 0xFF00000000) >> 32;
  value[index++] = (bearer_qos->gbr.br_dl & 0x00FF000000) >> 24;
  value[index++] = (bearer_qos->gbr.br_dl & 0x0000FF0000) >> 16;
  value[index++] = (bearer_qos->gbr.br_dl & 0x000000FF00) >> 8;
  value[index++] = (bearer_qos->gbr.br_dl & 0x00000000FF);
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_BEARER_LEVEL_QOS, 22, 0, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ip_address_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  gtp_ip_address_t *ip_address = (gtp_ip_address_t *) arg;

  DevAssert(ip_address);

  if (ieLength == 4) {
    /*
     * This is an IPv4 Address
     */
    ip_address->present = GTP_IP_ADDR_v4;
    memcpy(ip_address->address.v4, ieValue, 4);
  } else if (ieLength == 16) {
    /*
     * This is an IPv6 Address
     */
    ip_address->present = GTP_IP_ADDR_v6;
    memcpy(ip_address->address.v6, ieValue, 16);
  } else {
    /*
     * Length doesn't lie in possible values
     */
    return NW_GTPV2C_IE_INCORRECT;
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_ip_address_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const gtp_ip_address_t *ip_address)
{
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_delay_value_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  DelayValue_t *delay_value = (DelayValue_t *) arg;

  DevAssert(arg);

  if (ieLength != 1) {
    return NW_GTPV2C_IE_INCORRECT;
  }

  *delay_value = ieValue[0];
  OAILOG_DEBUG(LOG_S11, "\t - Delay Value %u\n", *delay_value);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_delay_value_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const DelayValue_t *delay_value)
{
  uint8_t value;
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(delay_value);
  value = *delay_value;
  rc =
    nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_DELAY_VALUE, 1, 0, (uint8_t *) &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ue_time_zone_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  UETimeZone_t *ue_time_zone = (UETimeZone_t *) arg;

  DevAssert(ue_time_zone);

  if (ieLength != 2) {
    return NW_GTPV2C_IE_INCORRECT;
  }

  ue_time_zone->time_zone = ieValue[0];
  ue_time_zone->daylight_saving_time = ieValue[1] & 0x03;
  OAILOG_DEBUG(LOG_S11, "\t - Time Zone    %u\n", ue_time_zone->time_zone);
  OAILOG_DEBUG(
    LOG_S11, "\t - Daylight SVT %u\n", ue_time_zone->daylight_saving_time);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_ue_time_zone_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const UETimeZone_t *ue_time_zone)
{
  uint8_t value[2];
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(ue_time_zone);
  value[0] = ue_time_zone->time_zone;
  value[1] = ue_time_zone->daylight_saving_time;
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_UE_TIME_ZONE, 2, 0, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_target_identification_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  target_identification_t *target_identification =
    (target_identification_t *) arg;

  DevAssert(target_identification);
  target_identification->target_type = ieValue[0];

  switch (target_identification->target_type) {
    case TARGET_ID_RNC_ID: {
      target_identification->target_id.rnc_id.lac =
        (ieValue[4] << 8) | ieValue[5];
      target_identification->target_id.rnc_id.rac = ieValue[6];

      if (ieLength == 11) {
        /*
         * Extended RNC id
         */
        target_identification->target_id.rnc_id.rnc_id =
          (ieValue[7] << 24) | (ieValue[8] << 16) | (ieValue[9] << 8) |
          (ieValue[10]);
      } else if (ieLength == 9) {
        /*
         * Normal RNC id
         */
        target_identification->target_id.rnc_id.rnc_id =
          (ieValue[7] << 8) | ieValue[8];
      } else {
        /*
         * This case is not possible
         */
        return NW_GTPV2C_IE_INCORRECT;
      }

      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- LAC 0x%04x\n",
        target_identification->target_id.rnc_id.lac);
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- RAC 0x%02x\n",
        target_identification->target_id.rnc_id.rac);
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- RNC 0x%08x\n",
        target_identification->target_id.rnc_id.rnc_id);
    } break;

    case TARGET_ID_MACRO_ENB_ID: {
      if (ieLength != 9) {
        return NW_GTPV2C_IE_INCORRECT;
      }

      target_identification->target_id.macro_enb_id.enb_id =
        ((ieValue[4] & 0x0F) << 16) | (ieValue[5] << 8) | ieValue[6];
      target_identification->target_id.macro_enb_id.tac =
        (ieValue[7] << 8) | ieValue[8];
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- ENB Id 0x%06x\n",
        target_identification->target_id.macro_enb_id.enb_id);
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- TAC    0x%04x\n",
        target_identification->target_id.macro_enb_id.tac);
    } break;

    case TARGET_ID_HOME_ENB_ID: {
      if (ieLength != 10) {
        return NW_GTPV2C_IE_INCORRECT;
      }

      target_identification->target_id.home_enb_id.enb_id =
        ((ieValue[4] & 0x0F) << 14) | (ieValue[5] << 16) | (ieValue[6] << 8) |
        ieValue[7];
      target_identification->target_id.home_enb_id.tac =
        (ieValue[8] << 8) | ieValue[9];
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- ENB Id 0x%07x\n",
        target_identification->target_id.home_enb_id.enb_id);
      OAILOG_DEBUG(
        LOG_S11,
        "\t\t- TAC    0x%04x\n",
        target_identification->target_id.home_enb_id.tac);
    } break;

    default: return NW_GTPV2C_IE_INCORRECT;
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_flags_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  bearer_flags_t *bearer_flags = (bearer_flags_t *) arg;

  DevAssert(arg);

  if (ieLength != 1) {
    return NW_GTPV2C_IE_INCORRECT;
  }

  bearer_flags->ppc = ieValue[0] & 0x01;
  bearer_flags->vb = ieValue[0] & 0x02;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_bearer_flags_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const bearer_flags_t *bearer_flags)
{
  nw_rc_t rc;
  uint8_t value;

  DevAssert(msg);
  DevAssert(bearer_flags);
  value = (bearer_flags->vb << 1) | bearer_flags->ppc;
  rc =
    nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_BEARER_FLAGS, 1, 0, (uint8_t *) &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_indication_flags_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  indication_flags_t *indication_flags = (indication_flags_t *) arg;

  DevAssert(indication_flags);

  if (2 <= ieLength) { // think about more than 3 later
    indication_flags->daf = (ieValue[0] >> DAF_FLAG_BIT_POS) & 0x01;
    indication_flags->dtf = (ieValue[0] >> DTF_FLAG_BIT_POS) & 0x01;
    indication_flags->hi = (ieValue[0] >> HI_FLAG_BIT_POS) & 0x01;
    indication_flags->dfi = (ieValue[0] >> DFI_FLAG_BIT_POS) & 0x01;
    indication_flags->oi = (ieValue[0] >> OI_FLAG_BIT_POS) & 0x01;
    indication_flags->isrsi = (ieValue[0] >> ISRSI_FLAG_BIT_POS) & 0x01;
    indication_flags->israi = (ieValue[0] >> ISRAI_FLAG_BIT_POS) & 0x01;
    indication_flags->sgwci = (ieValue[0] >> SGWCI_FLAG_BIT_POS) & 0x01;

    indication_flags->sqci = (ieValue[1] >> SQSI_FLAG_BIT_POS) & 0x01;
    indication_flags->uimsi = (ieValue[1] >> UIMSI_FLAG_BIT_POS) & 0x01;
    indication_flags->cfsi = (ieValue[1] >> CFSI_FLAG_BIT_POS) & 0x01;
    indication_flags->crsi = (ieValue[1] >> CRSI_FLAG_BIT_POS) & 0x01;
    indication_flags->p = (ieValue[1] >> P_FLAG_BIT_POS) & 0x01;
    indication_flags->pt = (ieValue[1] >> PT_FLAG_BIT_POS) & 0x01;
    indication_flags->si = (ieValue[1] >> SI_FLAG_BIT_POS) & 0x01;
    indication_flags->msv = (ieValue[1] >> MSV_FLAG_BIT_POS) & 0x01;

    if (2 == ieLength) {
      return NW_OK;
    }
    if (3 == ieLength) {
      indication_flags->spare1 = 0;
      indication_flags->spare2 = 0;
      indication_flags->spare3 = 0;
      indication_flags->s6af = (ieValue[2] >> S6AF_FLAG_BIT_POS) & 0x01;
      indication_flags->s4af = (ieValue[2] >> S4AF_FLAG_BIT_POS) & 0x01;
      indication_flags->mbmdt = (ieValue[2] >> MBMDT_FLAG_BIT_POS) & 0x01;
      indication_flags->israu = (ieValue[2] >> ISRAU_FLAG_BIT_POS) & 0x01;
      indication_flags->ccrsi = (ieValue[2] >> CRSI_FLAG_BIT_POS) & 0x01;
      return NW_OK;
    }
  }
  return NW_GTPV2C_IE_INCORRECT;
}

//------------------------------------------------------------------------------
int gtpv2c_indication_flags_ie_set(
  nw_gtpv2c_msg_handle_t *msg,
  const indication_flags_t *indication_flags)
{
  nw_rc_t rc;
  uint8_t value[3];

  DevAssert(msg);
  DevAssert(indication_flags);
  value[0] = (indication_flags->daf << DAF_FLAG_BIT_POS) |
             (indication_flags->dtf << DTF_FLAG_BIT_POS) |
             (indication_flags->hi << HI_FLAG_BIT_POS) |
             (indication_flags->dfi << DFI_FLAG_BIT_POS) |
             (indication_flags->oi << OI_FLAG_BIT_POS) |
             (indication_flags->isrsi << ISRSI_FLAG_BIT_POS) |
             (indication_flags->israi << ISRAI_FLAG_BIT_POS) |
             (indication_flags->sgwci << SGWCI_FLAG_BIT_POS);

  value[1] = (indication_flags->sqci << SQSI_FLAG_BIT_POS) |
             (indication_flags->uimsi << UIMSI_FLAG_BIT_POS) |
             (indication_flags->cfsi << CFSI_FLAG_BIT_POS) |
             (indication_flags->crsi << CRSI_FLAG_BIT_POS) |
             (indication_flags->p << P_FLAG_BIT_POS) |
             (indication_flags->pt << PT_FLAG_BIT_POS) |
             (indication_flags->si << SI_FLAG_BIT_POS) |
             (indication_flags->msv << MSV_FLAG_BIT_POS);

  value[2] = (indication_flags->s6af << S6AF_FLAG_BIT_POS) |
             (indication_flags->s4af << S4AF_FLAG_BIT_POS) |
             (indication_flags->mbmdt << MBMDT_FLAG_BIT_POS) |
             (indication_flags->israu << ISRAU_FLAG_BIT_POS) |
             (indication_flags->ccrsi << CCRSI_FLAG_BIT_POS);

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_INDICATION, 3, 0, (uint8_t *) value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_fqcsid_ie_get(
  uint8_t ieType,
  uint8_t ieLength,
  uint8_t ieInstance,
  uint8_t *ieValue,
  void *arg)
{
  FQ_CSID_t *fq_csid = (FQ_CSID_t *) arg;

  DevAssert(fq_csid);
  fq_csid->node_id_type = (ieValue[0] & 0xF0) >> 4;
  OAILOG_DEBUG(LOG_S11, "\t- FQ-CSID type %u\n", fq_csid->node_id_type);

  /*
   * NOTE: Values of Number of CSID other than 1 are only employed in the
   * * * * Delete PDN Connection Set Request and Response.
   */
  if ((ieValue[0] & 0x0F) != 1) {
    return NW_GTPV2C_IE_INCORRECT;
  }

  switch (fq_csid->node_id_type) {
    case GLOBAL_UNICAST_IPv4: {
      char ipv4[INET_ADDRSTRLEN];
      if (ieLength != 7) {
        return NW_GTPV2C_IE_INCORRECT;
      }

      int addr = (ieValue[1] << 24) | (ieValue[2] << 16) | (ieValue[3] << 8) |
                 (ieValue[4]);
      fq_csid->node_id.unicast_ipv4.s_addr = addr;
      fq_csid->csid = (ieValue[5] << 8) | ieValue[6];
      inet_ntop(
        AF_INET,
        (void *) &fq_csid->node_id.unicast_ipv4,
        ipv4,
        INET_ADDRSTRLEN);
      OAILOG_DEBUG(LOG_S11, "\t- v4 address [%s]\n", ipv4);
    } break;

    case GLOBAL_UNICAST_IPv6: {
      char ipv6[INET6_ADDRSTRLEN];

      if (ieLength != 19) {
        return NW_GTPV2C_IE_INCORRECT;
      }

      memcpy(fq_csid->node_id.unicast_ipv6.__in6_u.__u6_addr8, &ieValue[1], 16);
      fq_csid->csid = (ieValue[17] << 8) | ieValue[18];
      /*
       * Convert the ipv6 to printable string
       */
      inet_ntop(
        AF_INET6,
        (void *) &fq_csid->node_id.unicast_ipv6,
        ipv6,
        INET6_ADDRSTRLEN);
      OAILOG_DEBUG(LOG_S11, "\t- v6 address [%s]\n", ipv6);
    } break;

    default: return NW_GTPV2C_IE_INCORRECT;
  }

  OAILOG_DEBUG(LOG_S11, "\t- CSID 0x%04x\n", fq_csid->csid);
  return NW_OK;
}
