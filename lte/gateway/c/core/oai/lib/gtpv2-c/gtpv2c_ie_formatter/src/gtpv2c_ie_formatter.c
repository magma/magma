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

/*! \file s11_ie_formatter.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "../../gtpv2c_ie_formatter/shared/gtpv2c_ie_formatter.h"

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
#include "TrafficFlowAggregateDescription.h"
#include "security_types.h"
#include "common_types.h"

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_imsi_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  imsi_t* imsi    = (imsi_t*) arg;
  uint8_t decoded = 0;

  DevAssert(arg);
  DevAssert(ieLength <= IMSI_BCD_DIGITS_MAX);
  imsi->length = 0;
  while (decoded < ieLength) {
    uint8_t tmp                   = ieValue[decoded++];
    imsi->u.value[imsi->length++] = (tmp >> 4) | (tmp << 4);
  }
  for (int i = imsi->length; i < IMSI_BCD8_SIZE; i++) {
    imsi->u.value[i] = 0xff;
  }

  OAILOG_DEBUG(
      LOG_GTPV2C, "\t- IMSI (l=%d) %u%u%u%u%u%u%u%u%u%u%u%u%u%u%u\n",
      imsi->length, imsi->u.num.digit1, imsi->u.num.digit2, imsi->u.num.digit3,
      imsi->u.num.digit4, imsi->u.num.digit5, imsi->u.num.digit6,
      imsi->u.num.digit7, imsi->u.num.digit8, imsi->u.num.digit9,
      imsi->u.num.digit10, imsi->u.num.digit11, imsi->u.num.digit12,
      imsi->u.num.digit13, imsi->u.num.digit14, imsi->u.num.digit15);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_imsi_ie_set(nw_gtpv2c_msg_handle_t* msg, const imsi_t* imsi) {
  nw_rc_t rc;
  imsi_t imsi_nbo = {0};
  DevAssert(msg);
  DevAssert(imsi);

  for (int i = 0; i < IMSI_BCD8_SIZE; i++) {
    uint8_t tmp         = imsi->u.value[i];
    imsi_nbo.u.value[i] = (tmp >> 4) | (tmp << 4);
  }
  imsi_nbo.length = imsi->length;

  rc = nwGtpv2cMsgAddIe(
      *msg, NW_GTPV2C_IE_IMSI, (imsi_nbo.length + 1) / 2, 0,
      (uint8_t*) imsi_nbo.u.value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_cause_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  gtpv2c_cause_t* cause = (gtpv2c_cause_t*) arg;

  DevAssert(cause);
  cause->cause_value = ieValue[0];
  cause->cs          = ieValue[1] & 0x01;
  cause->bce         = (ieValue[1] & 0x02) >> 1;
  cause->pce         = (ieValue[1] & 0x04) >> 2;
  if (6 == ieLength) {
    cause->offending_ie_type   = ieValue[2];
    cause->offending_ie_length = ((uint16_t) ieValue[3]) << 8;
    cause->offending_ie_length |= ((uint16_t) ieValue[4]);
    cause->offending_ie_instance = ieValue[5] & 0x0F;
  }
  OAILOG_DEBUG(LOG_GTPV2C, "\t- Cause value %u\n", cause->cause_value);
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_cause_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const gtpv2c_cause_t* cause) {
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
    rc       = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_CAUSE, 6, 0, value);
  } else {
    rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_CAUSE, 2, 0, value);
  }

  DevAssert(NW_OK == rc);
  return rc == NW_OK ? 0 : -1;
}

//------------------------------------------------------------------------------
int gtpv2c_fteid_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const fteid_t* fteid, const uint8_t instance) {
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
    uint32_t hbo    = ntohl(fteid->ipv4_address.s_addr);
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
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  uint8_t offset = 0;
  fteid_t* fteid = (fteid_t*) arg;

  DevAssert(fteid);
  fteid->ipv4           = (ieValue[0] & 0x80) >> 7;
  fteid->ipv6           = (ieValue[0] & 0x40) >> 6;
  fteid->interface_type = ieValue[0] & 0x1F;
  OAILOG_DEBUG(LOG_GTPV2C, "\t- F-TEID type %d\n", fteid->interface_type);
  /*
   * Copy the TEID or GRE key
   */
  fteid->teid = ntoh_int32_buf(&ieValue[1]);
  OAILOG_DEBUG(LOG_GTPV2C, "\t- TEID/GRE    %08x\n", fteid->teid);

  if (fteid->ipv4 == 1) {
    /*
     * IPv4 present: copy the 4 bytes
     */
    uint32_t hbo = (((uint32_t) ieValue[5]) << 24) |
                   (((uint32_t) ieValue[6]) << 16) |
                   (((uint32_t) ieValue[7]) << 8) | (uint32_t) ieValue[8];
    fteid->ipv4_address.s_addr = htonl(hbo);
    offset                     = 4;
    OAILOG_DEBUG(
        LOG_GTPV2C, "\t- IPv4 addr   " IN_ADDR_FMT "\n",
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
        AF_INET6, (void*) &fteid->ipv6_address, ipv6_ascii, INET6_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_GTPV2C, "\t- IPv6 addr   %s\n", ipv6_ascii);
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_paa_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  uint8_t offset = 0;
  paa_t* paa     = (paa_t*) arg;

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
    inet_ntop(AF_INET, (void*) &paa->ipv4_address, ipv4, INET_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_S11, "\t- IPv4 addr %s\n", ipv4);
  }

  paa->pdn_type -= 1;
  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_paa_ie_set(nw_gtpv2c_msg_handle_t* msg, const paa_t* paa) {
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
  pdn_type     = paa->pdn_type + 1;
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
    uint32_t hbo   = ntohl(paa->ipv4_address.s_addr);
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
int gtpv2c_ebi_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const unsigned ebi, const uint8_t instance) {
  nw_rc_t rc;
  uint8_t value = 0;

  value = ebi & 0x0F;
  rc    = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_EBI, 1, instance, &value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ebi_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  uint8_t* ebi = (uint8_t*) arg;

  DevAssert(ebi);
  *ebi = ieValue[0] & 0x0F;
  OAILOG_DEBUG(LOG_S11, "\t- EBI %u\n", *ebi);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ambr_ie_set(nw_gtpv2c_msg_handle_t* msg, ambr_t* ambr) {
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(ambr);
  /*
   * MCC Decimal | MCC Hundreds
   */

  uint8_t ambr_br[8];
  uint8_t* p_ambr;
  p_ambr = ambr_br;

  memset(ambr_br, 0, 8);

  INT32_TO_BUFFER((ambr->br_ul / 1000), p_ambr);
  p_ambr += 4;

  INT32_TO_BUFFER(((ambr->br_dl / 1000)), p_ambr);
  // todo: byte order?

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_AMBR, 8, 0, ambr_br);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_ambr_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  ambr_t* ambr = (ambr_t*) arg;

  DevAssert(ambr);

  BUFFER_TO_INT32(ieValue, ambr->br_ul);
  BUFFER_TO_INT32((ieValue + 4), ambr->br_dl);
  ambr->br_dl *= 1000;
  ambr->br_ul *= 1000;
  OAILOG_DEBUG(LOG_S11, "\t- AMBR UL %" PRIu64 "\n", ambr->br_ul);
  OAILOG_DEBUG(LOG_S11, "\t- AMBR DL %" PRIu64 "\n", ambr->br_dl);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_qos_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  bearer_qos_t* bearer_qos = (bearer_qos_t*) arg;

  DevAssert(bearer_qos);

  if (22 <= ieLength) {
    int index       = 0;
    bearer_qos->pci = (ieValue[index] >> 6) & 0x01;
    bearer_qos->pci = (ieValue[index] >> 6) & 0x01;
    bearer_qos->pl  = (ieValue[index] >> 2) & 0x0F;
    bearer_qos->pvi = ieValue[index++] & 0x01;
    bearer_qos->qci = ieValue[index++];

    bearer_qos->mbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->mbr.br_ul |= (bit_rate_t) ieValue[index++];
    bearer_qos->mbr.br_ul *= 1000;

    bearer_qos->mbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->mbr.br_dl |= (bit_rate_t) ieValue[index++];
    bearer_qos->mbr.br_dl *= 1000;

    bearer_qos->gbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->gbr.br_ul |= (bit_rate_t) ieValue[index++];
    bearer_qos->gbr.br_ul *= 1000;

    bearer_qos->gbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    bearer_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    bearer_qos->gbr.br_dl |= (bit_rate_t) ieValue[index++];
    bearer_qos->gbr.br_dl *= 1000;

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
    nw_gtpv2c_msg_handle_t* msg, const bearer_qos_t* const bearer_qos) {
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
  value[index++] = ((bearer_qos->mbr.br_ul / 1000) & 0xFF00000000) >> 32;
  value[index++] = ((bearer_qos->mbr.br_ul / 1000) & 0x00FF000000) >> 24;
  value[index++] = ((bearer_qos->mbr.br_ul / 1000) & 0x0000FF0000) >> 16;
  value[index++] = ((bearer_qos->mbr.br_ul / 1000) & 0x000000FF00) >> 8;
  value[index++] = ((bearer_qos->mbr.br_ul / 1000) & 0x00000000FF);

  value[index++] = ((bearer_qos->mbr.br_dl / 1000) & 0xFF00000000) >> 32;
  value[index++] = ((bearer_qos->mbr.br_dl / 1000) & 0x00FF000000) >> 24;
  value[index++] = ((bearer_qos->mbr.br_dl / 1000) & 0x0000FF0000) >> 16;
  value[index++] = ((bearer_qos->mbr.br_dl / 1000) & 0x000000FF00) >> 8;
  value[index++] = ((bearer_qos->mbr.br_dl / 1000) & 0x00000000FF);

  value[index++] = ((bearer_qos->gbr.br_ul / 1000) & 0xFF00000000) >> 32;
  value[index++] = ((bearer_qos->gbr.br_ul / 1000) & 0x00FF000000) >> 24;
  value[index++] = ((bearer_qos->gbr.br_ul / 1000) & 0x0000FF0000) >> 16;
  value[index++] = ((bearer_qos->gbr.br_ul / 1000) & 0x000000FF00) >> 8;
  value[index++] = ((bearer_qos->gbr.br_ul / 1000) & 0x00000000FF);

  value[index++] = ((bearer_qos->gbr.br_dl / 1000) & 0xFF00000000) >> 32;
  value[index++] = ((bearer_qos->gbr.br_dl / 1000) & 0x00FF000000) >> 24;
  value[index++] = ((bearer_qos->gbr.br_dl / 1000) & 0x0000FF0000) >> 16;
  value[index++] = ((bearer_qos->gbr.br_dl / 1000) & 0x000000FF00) >> 8;
  value[index++] = ((bearer_qos->gbr.br_dl / 1000) & 0x00000000FF);

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_BEARER_LEVEL_QOS, 22, 0, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------

nw_rc_t gtpv2c_flow_qos_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  flow_qos_t* flow_qos = (flow_qos_t*) arg;

  DevAssert(flow_qos);

  if (21 <= ieLength) {
    int index     = 0;
    flow_qos->qci = ieValue[index++];

    flow_qos->mbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    flow_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    flow_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    flow_qos->mbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    flow_qos->mbr.br_ul |= (bit_rate_t) ieValue[index++];

    flow_qos->mbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    flow_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    flow_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    flow_qos->mbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    flow_qos->mbr.br_dl |= (bit_rate_t) ieValue[index++];

    flow_qos->gbr.br_ul = ((bit_rate_t) ieValue[index++]) << 32;
    flow_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 24);
    flow_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 16);
    flow_qos->gbr.br_ul |= (((bit_rate_t) ieValue[index++]) << 8);
    flow_qos->gbr.br_ul |= (bit_rate_t) ieValue[index++];

    flow_qos->gbr.br_dl = ((bit_rate_t) ieValue[index++]) << 32;
    flow_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 24);
    flow_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 16);
    flow_qos->gbr.br_dl |= (((bit_rate_t) ieValue[index++]) << 8);
    flow_qos->gbr.br_dl |= (bit_rate_t) ieValue[index++];

    if (21 < ieLength) {
      OAILOG_ERROR(LOG_S11, "TODO gtpv2c_flow_qos_ie_get() FlowQOS_t\n");
      return NW_GTPV2C_IE_INCORRECT;
    }
    return NW_OK;
  } else {
    OAILOG_ERROR(LOG_S11, "Bad IE length %" PRIu8 "\n", ieLength);
    return NW_GTPV2C_IE_INCORRECT;
  }
}

//------------------------------------------------------------------------------

int gtpv2c_flow_qos_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const flow_qos_t* flow_qos) {
  nw_rc_t rc;
  uint8_t value[21];
  int index = 0;

  DevAssert(msg);
  DevAssert(flow_qos);
  value[index++] = flow_qos->qci;
  /*
   * TODO: check endianness
   */

  value[index++] = (flow_qos->mbr.br_ul & 0xFF00000000) >> 32;
  value[index++] = (flow_qos->mbr.br_ul & 0x00FF000000) >> 24;
  value[index++] = (flow_qos->mbr.br_ul & 0x0000FF0000) >> 16;
  value[index++] = (flow_qos->mbr.br_ul & 0x000000FF00) >> 8;
  value[index++] = (flow_qos->mbr.br_ul & 0x00000000FF);

  value[index++] = (flow_qos->mbr.br_dl & 0xFF00000000) >> 32;
  value[index++] = (flow_qos->mbr.br_dl & 0x00FF000000) >> 24;
  value[index++] = (flow_qos->mbr.br_dl & 0x0000FF0000) >> 16;
  value[index++] = (flow_qos->mbr.br_dl & 0x000000FF00) >> 8;
  value[index++] = (flow_qos->mbr.br_dl & 0x00000000FF);

  value[index++] = (flow_qos->gbr.br_ul & 0xFF00000000) >> 32;
  value[index++] = (flow_qos->gbr.br_ul & 0x00FF000000) >> 24;
  value[index++] = (flow_qos->gbr.br_ul & 0x0000FF0000) >> 16;
  value[index++] = (flow_qos->gbr.br_ul & 0x000000FF00) >> 8;
  value[index++] = (flow_qos->gbr.br_ul & 0x00000000FF);

  value[index++] = (flow_qos->gbr.br_dl & 0xFF00000000) >> 32;
  value[index++] = (flow_qos->gbr.br_dl & 0x00FF000000) >> 24;
  value[index++] = (flow_qos->gbr.br_dl & 0x0000FF0000) >> 16;
  value[index++] = (flow_qos->gbr.br_dl & 0x000000FF00) >> 8;
  value[index++] = (flow_qos->gbr.br_dl & 0x00000000FF);

  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_FLOW_QOS, 21, 0, value);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_tft_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  uint8_t offset               = 0;
  traffic_flow_template_t* tft = (traffic_flow_template_t*) arg;

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
    nw_gtpv2c_msg_handle_t* msg, const traffic_flow_template_t* const tft) {
  uint8_t temp[TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH];
  uint8_t offset = 0;
  nw_rc_t rc     = NW_OK;

  DevAssert(tft);
  offset = encode_traffic_flow_template(
      tft, temp, TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH);
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_BEARER_TFT, offset, 0, temp);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_tad_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  uint8_t offset = 0;
  traffic_flow_aggregate_description_t* tad =
      (traffic_flow_aggregate_description_t*) arg;

  DevAssert(tad);
  offset = decode_traffic_flow_template(tad, ieValue, ieLength);
  if ((0 < offset) && (TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH >= offset))
    return NW_OK;
  else {
    OAILOG_ERROR(LOG_S11, "Incorrect TAD IE\n");
    return NW_GTPV2C_IE_INCORRECT;
  }
}

//------------------------------------------------------------------------------
int gtpv2c_tad_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const traffic_flow_aggregate_description_t* tad) {
  uint8_t temp[TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH];
  uint8_t offset = 0;
  nw_rc_t rc     = NW_OK;

  DevAssert(tad);
  offset = encode_traffic_flow_template(
      tad, temp, TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH);
  rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_TAD, offset, 0, temp);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_target_identification_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  target_identification_t* target_identification =
      (target_identification_t*) arg;

  DevAssert(target_identification);
  target_identification->target_type = ieValue[0];

  target_identification->mcc[1] = (ieValue[1] & 0xF0) >> 4;
  target_identification->mcc[0] = (ieValue[1] & 0x0F);
  target_identification->mcc[2] = (ieValue[2] & 0x0F);

  if ((ieValue[1] & 0xF0) == 0xF0) {
    /*
     * Two digits MNC
     */
    target_identification->mnc[0] = 0;
    target_identification->mnc[1] = (ieValue[3] & 0x0F);
    target_identification->mnc[2] = (ieValue[3] & 0xF0) >> 4;
  } else {
    target_identification->mnc[0] = (ieValue[3] & 0x0F);
    target_identification->mnc[1] = (ieValue[3] & 0xF0) >> 4;
    target_identification->mnc[2] = (ieValue[2] & 0xF0) >> 4;
  }

  switch (target_identification->target_type) {
    case TARGET_ID_RNC_ID: {
      target_identification->target_id.rnc_id.lac =
          (ieValue[4] << 8) | ieValue[5];
      target_identification->target_id.rnc_id.rac = ieValue[6];

      if (ieLength == 11) {
        /*
         * Extended RNC id
         */
        target_identification->target_id.rnc_id.id =
            (ieValue[7] << 24) | (ieValue[8] << 16);
        target_identification->target_id.rnc_id.xid =
            (ieValue[9] << 8) | (ieValue[10]);
      } else if (ieLength == 9) {
        /*
         * Normal RNC id
         */
        target_identification->target_id.rnc_id.id =
            (ieValue[7] << 8) | ieValue[8];
        target_identification->target_id.rnc_id.xid = 0;
      } else {
        /*
         * This case is not possible
         */
        return NW_GTPV2C_IE_INCORRECT;
      }

      OAILOG_DEBUG(
          LOG_S10, "\t\t- LAC 0x%04x\n",
          target_identification->target_id.rnc_id.lac);
      OAILOG_DEBUG(
          LOG_S10, "\t\t- RAC 0x%02x\n",
          target_identification->target_id.rnc_id.rac);
      OAILOG_DEBUG(
          LOG_S10, "\t\t- RNC -ID 0x%08x\n",
          target_identification->target_id.rnc_id.id);
      OAILOG_DEBUG(
          LOG_S10, "\t\t- RNC -XID 0x%08x\n",
          target_identification->target_id.rnc_id.xid);
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
          LOG_S10, "\t\t- ENB Id 0x%06x\n",
          target_identification->target_id.macro_enb_id.enb_id);
      OAILOG_DEBUG(
          LOG_S10, "\t\t- TAC    0x%04x\n",
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
          LOG_S10, "\t\t- ENB Id 0x%07x\n",
          target_identification->target_id.home_enb_id.enb_id);
      OAILOG_DEBUG(
          LOG_S10, "\t\t- TAC    0x%04x\n",
          target_identification->target_id.home_enb_id.tac);
    } break;

    default:
      return NW_GTPV2C_IE_INCORRECT;
  }

  return NW_OK;
}

//------------------------------------------------------------------------------
int gtpv2c_target_identification_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const target_identification_t* target_identification) {
  nw_rc_t rc;
  uint8_t target_id[3];
  // uint32_t                               macro_enb_id;
  // uint16_t                               tac;

  DevAssert(msg);
  DevAssert(target_identification);
  /*
   * MCC Decimal | MCC Hundreds
   */
  target_id[0] = ((target_identification->mcc[1] & 0x0F) << 4) |
                 (target_identification->mcc[0] & 0x0F);
  target_id[1] = target_identification->mcc[2] & 0x0F;

  if ((target_identification->mnc[0] & 0xF) == 0xF) {
    /*
     * Only two digits
     */
    target_id[1] |= 0xF0;
    target_id[2] = ((target_identification->mnc[2] & 0x0F) << 4) |
                   (target_identification->mnc[1] & 0x0F);
  } else {
    target_id[1] |= (target_identification->mnc[2] & 0x0F) << 4;
    target_id[2] = ((target_identification->mnc[1] & 0x0F) << 4) |
                   (target_identification->mnc[0] & 0x0F);
  }

  /** Build an array for the TargetIe payload. */

  uint8_t enbId[4];

  uint8_t targetIeBuf[10];
  uint8_t* pTargetIeBuf = targetIeBuf;

  uint8_t target_id_length = 0;

  memset(pTargetIeBuf, 0, 10);
  /** Target Type. */
  *pTargetIeBuf = target_identification->target_type;
  target_id_length++;
  /** Set the plmn. */
  memcpy((uint8_t*) (pTargetIeBuf + target_id_length), target_id, 3);
  target_id_length += 3;

  /** Check the eNB Type. */
  if (target_identification->target_type == TARGET_ID_HOME_ENB_ID) {
    /** Set the enb-id as home. No need to skipping places. */
    enbId[0] = target_identification->target_id.home_enb_id.enb_id >> 24 & 0x0F;
    enbId[1] = target_identification->target_id.home_enb_id.enb_id >> 16 & 0xFF;
    enbId[2] = target_identification->target_id.home_enb_id.enb_id >> 8 & 0xFF;
    enbId[3] = target_identification->target_id.home_enb_id.enb_id & 0xFF;
    memcpy(pTargetIeBuf + target_id_length, enbId, 4);
    target_id_length += 4;
    /** TAC. */
    *((uint16_t*) (pTargetIeBuf + target_id_length)) =
        htons(target_identification->target_id.home_enb_id.tac);
    target_id_length += 2;
  } else if (target_identification->target_type == TARGET_ID_MACRO_ENB_ID) {
    /** Here always assume macro, also for RNC. */

    /** Macro Enb Id. */
    enbId[0] = target_identification->target_id.home_enb_id.enb_id >> 16 & 0xFF;
    enbId[1] = target_identification->target_id.home_enb_id.enb_id >> 8 & 0xFF;
    enbId[2] = target_identification->target_id.home_enb_id.enb_id & 0xFF;
    memcpy(pTargetIeBuf + target_id_length, enbId, 3);
    target_id_length += 3;
    /** TAC. */
    *((uint16_t*) (pTargetIeBuf + target_id_length)) =
        htons(target_identification->target_id.macro_enb_id.tac);
    target_id_length += 2;
  } else if (target_identification->target_type == TARGET_ID_RNC_ID) {
    /** Encode the RNC ID type. */
    /** Here always assume macro, also for RNC. */
    /** Macro Enb Id. */
    // Skip 1 place --> todo: copy 20 bits
    target_id_length += 2;
    *((uint32_t*) (pTargetIeBuf + target_id_length)) =
        htons(target_identification->target_id.rnc_id.lac);
    target_id_length += 2;
    *((uint32_t*) (pTargetIeBuf + target_id_length)) =
        target_identification->target_id.rnc_id.rac;
    target_id_length++;
    /** RNC-ID. */
    *((uint16_t*) (pTargetIeBuf + target_id_length)) =
        htons(target_identification->target_id.rnc_id.id);
    target_id_length += 2;
    /** RNC-XID. */
    *((uint16_t*) (pTargetIeBuf + target_id_length)) =
        htonl(target_identification->target_id.rnc_id.xid);
    target_id_length += 2;
  } else {
    /** Received unhandled cell id. */
    return RETURNerror;
  }

  // todo: extended f-cause?!
  /** Reset the pointer. */

  rc = nwGtpv2cMsgAddIe(
      *msg, NW_GTPV2C_IE_TARGET_IDENTIFICATION, target_id_length,
      NW_GTPV2C_IE_INSTANCE_ZERO, pTargetIeBuf);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_context_created_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  bearer_contexts_created_t* bearer_contexts = (bearer_contexts_created_t*) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
      MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS >=
      bearer_contexts->num_bearer_context);
  bearer_context_created_t* bearer_context =
      &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint16_t read = 0;
  nw_rc_t rc;

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t* ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t*) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
            ie_p->t, ntohs(ie_p->l), ie_p->i,
            &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
            &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_FTEID:
        switch (ie_p->i) {
          case 0:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s1u_sgw_fteid);
            break;
          case 1:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s4u_sgw_fteid);
            break;
          case 2:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s5_s8_u_pgw_fteid);
            break;
          case 3:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s12_sgw_fteid);
            break;
          default:
            OAILOG_ERROR(
                LOG_S11, "Received unexpected IE %u instance %u\n", ie_p->t,
                ie_p->i);
            return NW_GTPV2C_IE_INCORRECT;
        }
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_CAUSE:
        rc = gtpv2c_cause_ie_get(
            ie_p->t, ntohs(ie_p->l), ie_p->i,
            &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
            &bearer_context->cause);
        break;

      case NW_GTPV2C_IE_BEARER_LEVEL_QOS:
        rc = gtpv2c_bearer_qos_ie_get(
            ie_p->t, ntohs(ie_p->l), ie_p->i,
            &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
            &bearer_context->bearer_level_qos);
        break;

      case NW_GTPV2C_IE_BEARER_FLAGS:
        OAILOG_WARNING(
            LOG_S11, "Received IE BEARER_FLAGS %u to implement\n", ie_p->t);
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
    nw_gtpv2c_msg_handle_t* msg, const bearer_context_created_t* const bearer) {
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer);
  /*
   * Start section for grouped IE: bearer context created
   */
  rc = nwGtpv2cMsgGroupedIeStart(
      *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(msg, bearer->eps_bearer_id, NW_GTPV2C_IE_INSTANCE_ZERO);
  rc = gtpv2c_cause_ie_set(msg, &bearer->cause);
  DevAssert(NW_OK == rc);
  // todo: data forwarding..
  /*
   * End section for grouped IE: bearer context created
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_bearer_context_to_be_created_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  bearer_contexts_to_be_created_t* bearer_contexts =
      (bearer_contexts_to_be_created_t*) arg;
  DevAssert(bearer_contexts);
  DevAssert(0 <= bearer_contexts->num_bearer_context);
  DevAssert(
      MSG_CREATE_SESSION_REQUEST_MAX_BEARER_CONTEXTS >=
      bearer_contexts->num_bearer_context);
  bearer_context_to_be_created_t* bearer_context =
      &bearer_contexts->bearer_contexts[bearer_contexts->num_bearer_context];
  uint16_t read = 0;
  nw_rc_t rc;

  while (ieLength > read) {
    nw_gtpv2c_ie_tlv_t* ie_p;

    ie_p = (nw_gtpv2c_ie_tlv_t*) &ieValue[read];

    switch (ie_p->t) {
      case NW_GTPV2C_IE_EBI:
        rc = gtpv2c_ebi_ie_get(
            ie_p->t, ntohs(ie_p->l), ie_p->i,
            &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
            &bearer_context->eps_bearer_id);
        DevAssert(NW_OK == rc);
        break;

      case NW_GTPV2C_IE_BEARER_LEVEL_QOS:
        rc = gtpv2c_bearer_qos_ie_get(
            ie_p->t, ntohs(ie_p->l), ie_p->i,
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
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s1u_enb_fteid);
            break;
          case 1:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s4u_sgsn_fteid);
            break;
          case 2:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s5_s8_u_sgw_fteid);
            break;
          case 3:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s5_s8_u_pgw_fteid);
            break;
          case 4:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s12_rnc_fteid);
            break;
          case 5:
            rc = gtpv2c_fteid_ie_get(
                ie_p->t, ntohs(ie_p->l), ie_p->i,
                &ieValue[read + sizeof(nw_gtpv2c_ie_tlv_t)],
                &bearer_context->s2b_u_epdg_fteid);
            break;
          default:
            OAILOG_ERROR(
                LOG_S11, "Received unexpected IE %u instance %u\n", ie_p->t,
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
int gtpv2c_bearer_context_to_be_created_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_created_t* bearer_context) {
  nw_rc_t rc;

  DevAssert(msg);
  DevAssert(bearer_context);
  /*
   * Start section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeStart(
      *msg, NW_GTPV2C_IE_BEARER_CONTEXT, NW_GTPV2C_IE_INSTANCE_ZERO);
  DevAssert(NW_OK == rc);
  gtpv2c_ebi_ie_set(
      msg, bearer_context->eps_bearer_id, NW_GTPV2C_IE_INSTANCE_ZERO);
  gtpv2c_bearer_qos_ie_set(msg, &bearer_context->bearer_level_qos);
  /** Set the TFT. */
  if (&bearer_context->tft) {
    gtpv2c_tft_ie_set(msg, &bearer_context->tft);
  }

  /*
   * End section for grouped IE: bearer context to create
   */
  rc = nwGtpv2cMsgGroupedIeEnd(*msg);
  DevAssert(NW_OK == rc);
  return RETURNok;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_selection_mode_ie_set(
    nw_gtpv2c_msg_handle_t* msg, SelectionMode_t* sm) {
  DevAssert(msg);
  DevAssert(sm);
  uint8_t val8 = *sm;
  nw_rc_t rc = nwGtpv2cMsgAddIe(*msg, NW_GTPV2C_IE_SELECTION_MODE, 1, 0, &val8);
  DevAssert(NW_OK == rc);
  return NW_OK;
}

//------------------------------------------------------------------------------
nw_rc_t gtpv2c_selection_mode_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  SelectionMode_t* sm = (SelectionMode_t*) arg;

  DevAssert(sm);
  *sm = (ieValue[0] & 0x03);
  OAILOG_DEBUG(LOG_GTPV2C, "\t- Selection Mode  %d\n", *sm);
  return NW_OK;
}
