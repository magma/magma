/*----------------------------------------------------------------------------*
 *                                                                            *
                                n w - g t p v 2 c
      G P R S   T u n n e l i n g    P r o t o c o l   v 2 c    S t a c k
 *                                                                            *
 *                                                                            *
   Copyright (c) 2010-2011 Amit Chawre
   All rights reserved.
 *                                                                            *
   Redistribution and use in source and binary forms, with or without
   modification, are permitted provided that the following conditions
   are met:
 *                                                                            *
   1. Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
   2. Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
   3. The name of the author may not be used to endorse or promote products
      derived from this software without specific prior written permission.
 *                                                                            *
   THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
   IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
   OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
   IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,
   INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
   NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
   DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
   THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
   (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
   THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
  ----------------------------------------------------------------------------*/

#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <inttypes.h>
#include <stdbool.h>

#include "bstrlib.h"

#include "NwTypes.h"
#include "NwLog.h"
#include "NwUtils.h"
#include "NwGtpv2cLog.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cPrivate.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "log.h"

#ifdef __cplusplus
extern "C" {
#endif

/*----------------------------------------------------------------------------*
                       P R I V A T E     F U N C T I O N S
  ----------------------------------------------------------------------------*/

static nw_gtpv2c_msg_t* gpGtpv2cMsgPool = NULL;

/*----------------------------------------------------------------------------*
                         P U B L I C   F U N C T I O N S
  ----------------------------------------------------------------------------*/

nw_rc_t nwGtpv2cMsgNew(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint8_t teidPresent,
    NW_IN uint8_t msgType, NW_IN uint32_t teid, NW_IN uint32_t seqNum,
    NW_OUT nw_gtpv2c_msg_handle_t* phMsg) {
  nw_gtpv2c_stack_t* pStack = (nw_gtpv2c_stack_t*) hGtpcStackHandle;
  nw_gtpv2c_msg_t* pMsg;
  NW_ASSERT(pStack);

  if (gpGtpv2cMsgPool) {
    pMsg            = gpGtpv2cMsgPool;
    gpGtpv2cMsgPool = gpGtpv2cMsgPool->next;
  } else {
    NW_GTPV2C_MALLOC(pStack, sizeof(nw_gtpv2c_msg_t), pMsg, nw_gtpv2c_msg_t*);
    OAILOG_DEBUG(LOG_GTPV2C, "ALLOCATED NEW MESSAGE %p!\n", pMsg);
  }

  if (pMsg) {
    pMsg->version     = NW_GTP_VERSION;
    pMsg->teidPresent = teidPresent;
    pMsg->msgType     = msgType;
    pMsg->teid        = teid;
    pMsg->seqNum      = seqNum;
    pMsg->msgLen = (NW_GTPV2C_EPC_SPECIFIC_HEADER_SIZE - (teidPresent ? 0 : 4));
    pMsg->groupedIeEncodeStack.top = 0;
    pMsg->hStack                   = hGtpcStackHandle;
    *phMsg                         = (nw_gtpv2c_msg_handle_t) pMsg;
    OAILOG_DEBUG(LOG_GTPV2C, "Created message %p!\n", pMsg);
    return NW_OK;
  }

  return NW_FAILURE;
}

nw_rc_t nwGtpv2cMsgFromBufferNew(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint8_t* pBuf,
    NW_IN uint32_t bufLen, NW_OUT nw_gtpv2c_msg_handle_t* phMsg) {
  nw_gtpv2c_stack_t* pStack = (nw_gtpv2c_stack_t*) hGtpcStackHandle;
  nw_gtpv2c_msg_t* pMsg;

  NW_ASSERT(pStack);

  if (gpGtpv2cMsgPool) {
    pMsg            = gpGtpv2cMsgPool;
    gpGtpv2cMsgPool = gpGtpv2cMsgPool->next;
  } else {
    NW_GTPV2C_MALLOC(pStack, sizeof(nw_gtpv2c_msg_t), pMsg, nw_gtpv2c_msg_t*);
  }

  if (pMsg) {
    *phMsg = (nw_gtpv2c_msg_handle_t) pMsg;
    memcpy(pMsg->msgBuf, pBuf, bufLen);
    pMsg->msgLen      = bufLen;
    pMsg->version     = ((*pBuf) & 0xE0) >> 5;
    pMsg->teidPresent = ((*pBuf) & 0x08) >> 3;
    pBuf++;
    pMsg->msgType = *(pBuf);
    pBuf += 3;

    if (pMsg->teidPresent) {
      pMsg->teid = ntohl(*((uint32_t*) (pBuf)));
      pBuf += 4;
    }

    memcpy(((uint8_t*) &pMsg->seqNum) + 1, pBuf, 3);
    pMsg->seqNum = ntohl(pMsg->seqNum);
    pMsg->hStack = hGtpcStackHandle;
    OAILOG_DEBUG(LOG_GTPV2C, "Created message %p!\n", pMsg);
    return NW_OK;
  }

  return NW_FAILURE;
}

nw_rc_t nwGtpv2cMsgDelete(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  OAILOG_DEBUG(LOG_GTPV2C, "Purging message 0x%" PRIxPTR "!\n", hMsg);
  ((nw_gtpv2c_msg_t*) hMsg)->next = gpGtpv2cMsgPool;
  gpGtpv2cMsgPool                 = (nw_gtpv2c_msg_t*) hMsg;
  OAILOG_DEBUG(
      LOG_GTPV2C, "Message pool %p! Next element %p !\n", gpGtpv2cMsgPool,
      gpGtpv2cMsgPool->next);

  return NW_OK;
}

/**
   Set TEID for gtpv2c message.

   @param[in] hMsg : Message handle.
   @param[in] teid: TEID value.
*/

nw_rc_t nwGtpv2cMsgSetTeid(NW_IN nw_gtpv2c_msg_handle_t hMsg, uint32_t teid) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  thiz->teid = teid;
  return NW_OK;
}

/**
   Set TEID present flag for gtpv2c message.

   @param[in] hMsg : Message handle.
   @param[in] teidPesent: Flag boolean value.
*/

nw_rc_t nwGtpv2cMsgSetTeidPresent(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, bool teidPresent) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  thiz->teidPresent = teidPresent;
  return NW_OK;
}

/**
   Set sequence for gtpv2c message.

   @param[in] hMsg : Message handle.
   @param[in] seqNum: Flag boolean value.
*/

nw_rc_t nwGtpv2cMsgSetSeqNumber(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, uint32_t seqNum) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  thiz->seqNum = seqNum;
  return NW_OK;
}

/**
   Get TEID present for gtpv2c message.

   @param[in] hMsg : Message handle.
*/

uint32_t nwGtpv2cMsgGetTeid(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  return (thiz->teid);
}

/**
   Get TEID present for gtpv2c message.

   @param[in] hMsg : Message handle.
*/

bool nwGtpv2cMsgGetTeidPresent(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  return (thiz->teidPresent);
}

/**
   Get sequence number for gtpv2c message.

   @param[in] hMsg : Message handle.
*/

uint32_t nwGtpv2cMsgGetSeqNumber(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  return (thiz->seqNum);
}

/**
   Get msg type for gtpv2c message.

   @param[in] hMsg : Message handle.
*/

uint32_t nwGtpv2cMsgGetMsgType(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  return (thiz->msgType);
}

/**
   Get msg type for gtpv2c message.

   @param[in] hMsg : Message handle.
*/

uint32_t nwGtpv2cMsgGetLength(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  return (thiz->msgLen);
}

nw_rc_t nwGtpv2cMsgAddIeTV1(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_IN uint8_t value) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv1_t* pIe;

  pIe    = (nw_gtpv2c_ie_tv1_t*) (pMsg->msgBuf + pMsg->msgLen);
  pIe->t = type;
  pIe->l = htons(0x0001);
  pIe->i = instance & 0x00ff;
  pIe->v = value;
  pMsg->msgLen += sizeof(nw_gtpv2c_ie_tv1_t);
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgAddIeTV2(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_IN uint16_t value) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv2_t* pIe;

  pIe    = (nw_gtpv2c_ie_tv2_t*) (pMsg->msgBuf + pMsg->msgLen);
  pIe->t = type;
  pIe->l = htons(0x0002);
  pIe->i = instance & 0x00ff;
  pIe->v = htons(value);
  pMsg->msgLen += sizeof(nw_gtpv2c_ie_tv2_t);
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgAddIeTV4(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_IN uint32_t value) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv4_t* pIe;

  pIe    = (nw_gtpv2c_ie_tv4_t*) (pMsg->msgBuf + pMsg->msgLen);
  pIe->t = type;
  pIe->l = htons(0x0004);
  pIe->i = instance & 0x00ff;
  pIe->v = htonl(value);
  pMsg->msgLen += sizeof(nw_gtpv2c_ie_tv4_t);
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgAddIe(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint16_t length, NW_IN uint8_t instance, NW_IN uint8_t* pVal) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  pIe    = (nw_gtpv2c_ie_tlv_t*) (pMsg->msgBuf + pMsg->msgLen);
  pIe->t = type;
  pIe->l = htons(length);
  pIe->i = instance & 0x00ff;
  memcpy(((uint8_t*) pIe) + 4, pVal, length);
  pMsg->msgLen += (4 + length);
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgGroupedIeStart(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  pIe    = (nw_gtpv2c_ie_tlv_t*) (pMsg->msgBuf + pMsg->msgLen);
  pIe->t = type;
  pIe->i = instance & 0x00ff;
  pMsg->msgLen += (4);
  pIe->l = (pMsg->msgLen);
  NW_ASSERT(pMsg->groupedIeEncodeStack.top < NW_GTPV2C_MAX_GROUPED_IE_DEPTH);
  pMsg->groupedIeEncodeStack.pIe[pMsg->groupedIeEncodeStack.top] = pIe;
  pMsg->groupedIeEncodeStack.top++;
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgGroupedIeEnd(NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  NW_ASSERT(pMsg->groupedIeEncodeStack.top > 0);
  pMsg->groupedIeEncodeStack.top--;
  pIe    = pMsg->groupedIeEncodeStack.pIe[pMsg->groupedIeEncodeStack.top];
  pIe->l = htons(pMsg->msgLen - pIe->l);
  return NW_OK;
}

nw_rc_t nwGtpv2cMsgAddIeCause(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_IN uint8_t causeValue, NW_IN uint8_t bitFlags,
    NW_IN uint8_t offendingIeType, NW_IN uint8_t offendingIeInstance) {
  uint8_t causeBuf[8];

  causeBuf[0] = causeValue;
  causeBuf[1] = bitFlags;

  if (offendingIeType) {
    causeBuf[2] = offendingIeType;
    causeBuf[3] = 0;
    causeBuf[4] = 0;
    causeBuf[5] = (offendingIeInstance & 0x0f);
  }

  return (nwGtpv2cMsgAddIe(
      hMsg, NW_GTPV2C_IE_CAUSE, (offendingIeType ? 6 : 2), instance, causeBuf));
}

nw_rc_t nwGtpv2cMsgAddIeFteid(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_IN uint8_t ifType, NW_IN const uint32_t teidOrGreKey,
    NW_IN const struct in_addr* const ipv4Addr,
    NW_IN const struct in6_addr* const pIpv6Addr) {
  uint8_t fteidBuf[32];
  uint8_t* pFteidBuf = fteidBuf;

  fteidBuf[0] = (ifType & 0x1F);
  pFteidBuf++;
  *((uint32_t*) (pFteidBuf)) = htonl((teidOrGreKey));
  pFteidBuf += 4;

  if (ipv4Addr) {
    fteidBuf[0] |= (0x01 << 7);
    *((uint32_t*) (pFteidBuf)) = ipv4Addr->s_addr;
    pFteidBuf += 4;
  }

  if (pIpv6Addr) {
    fteidBuf[0] |= (0x01 << 6);
    memcpy((pFteidBuf), pIpv6Addr->__in6_u.__u6_addr8, 16);
    pFteidBuf += 16;
  }

  return (nwGtpv2cMsgAddIe(
      hMsg, NW_GTPV2C_IE_FTEID, (pFteidBuf - fteidBuf), instance, fteidBuf));
}

nw_rc_t nwGtpv2cMsgAddIeFCause(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_IN uint8_t fcauseType, NW_IN uint8_t fcauseValue) {
  uint8_t fcauseBuf[8];

  fcauseBuf[0] = fcauseType;
  fcauseBuf[1] = fcauseValue;
  // todo: extended f-cause?!
  return (nwGtpv2cMsgAddIe(hMsg, NW_GTPV2C_IE_F_CAUSE, 2, instance, fcauseBuf));
}

nw_rc_t nwGtpv2cMsgAddIeFContainer(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_IN uint8_t* container_value, NW_IN uint32_t container_data_size,
    NW_IN uint8_t container_type) {
  uint8_t fContainerBuf[container_data_size + 1];

  fContainerBuf[0] = container_type;
  memcpy(&fContainerBuf[1], container_value, container_data_size);
  // todo: extended f-cause?!
  return (nwGtpv2cMsgAddIe(
      hMsg, NW_GTPV2C_IE_F_CONTAINER, container_data_size + 1, instance,
      fContainerBuf));
}

nw_rc_t nwGtpv2cMsgAddIeCompleteRequestMessage(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_IN uint8_t* request_value, NW_IN uint32_t request_size,
    NW_IN uint8_t request_type) {
  uint8_t requestBuf[request_size];

  requestBuf[0] = request_type;
  memcpy(&requestBuf[1], request_value, request_size - 1);
  // todo: extended f-cause?!
  return (nwGtpv2cMsgAddIe(
      hMsg, NW_GTPV2C_IE_COMPLETE_REQUEST_MESSAGE, request_size, instance,
      requestBuf));
}

bool nwGtpv2cMsgIsIePresent(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;

  if ((nw_gtpv2c_ie_tv1_t*) thiz->pIe[type][instance]) return true;

  return false;
}

nw_rc_t nwGtpv2cMsgGetIeTV1(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_OUT uint8_t* pVal) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv1_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tv1_t*) thiz->pIe[type][instance];

    if (ntohs(pIe->l) != 0x01) return NW_GTPV2C_IE_INCORRECT;

    if (pVal) *pVal = pIe->v;

    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeTV2(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_OUT uint16_t* pVal) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv2_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tv2_t*) thiz->pIe[type][instance];

    if (ntohs(pIe->l) != 0x02) return NW_GTPV2C_IE_INCORRECT;

    if (pVal) *pVal = ntohs(pIe->v);

    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeTV4(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_OUT uint32_t* pVal) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv4_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tv4_t*) thiz->pIe[type][instance];

    if (ntohs(pIe->l) != 0x04) return NW_GTPV2C_IE_INCORRECT;

    if (pVal) *pVal = ntohl(pIe->v);

    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeTV8(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_OUT uint64_t* pVal) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tv8_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tv8_t*) thiz->pIe[type][instance];

    if (ntohs(pIe->l) != 0x08) return NW_GTPV2C_IE_INCORRECT;

    if (pVal) *pVal = NW_NTOHLL((pIe->v));

    return NW_OK;
  }

  OAILOG_ERROR(
      LOG_GTPV2C, "Cannot retrieve IE of type %u instance %u !\n", type,
      instance);
  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeTlv(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_IN uint16_t maxLen, NW_OUT uint8_t* pVal,
    NW_OUT uint16_t* pLen) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tlv_t*) thiz->pIe[type][instance];

    if (ntohs(pIe->l) <= maxLen) {
      if (pVal) memcpy(pVal, ((uint8_t*) pIe) + 4, ntohs(pIe->l));

      if (pLen) *pLen = ntohs(pIe->l);

      return NW_OK;
    }
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeTlvP(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t type,
    NW_IN uint8_t instance, NW_OUT uint8_t** ppVal, NW_OUT uint16_t* pLen) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[type][instance]) {
    pIe = (nw_gtpv2c_ie_tlv_t*) thiz->pIe[type][instance];

    if (ppVal) *ppVal = ((uint8_t*) pIe) + 4;

    if (pLen) *pLen = ntohs(pIe->l);

    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeCause(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_OUT uint8_t* causeValue, NW_OUT uint8_t* flags,
    NW_OUT uint8_t* offendingIeType, NW_OUT uint8_t* offendingIeInstance) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[NW_GTPV2C_IE_CAUSE][instance]) {
    pIe         = (nw_gtpv2c_ie_tlv_t*) thiz->pIe[NW_GTPV2C_IE_CAUSE][instance];
    *causeValue = *((uint8_t*) (((uint8_t*) pIe) + 4));
    *flags      = *((uint8_t*) (((uint8_t*) pIe) + 5));

    if (pIe->l == 6) {
      *offendingIeType = *((uint8_t*) (((uint8_t*) pIe) + 6));
      *offendingIeType = *((uint8_t*) (((uint8_t*) pIe) + 8));
    }

    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgGetIeFteid(
    NW_IN nw_gtpv2c_msg_handle_t hMsg, NW_IN uint8_t instance,
    NW_OUT uint8_t* ifType, NW_OUT uint32_t* teidOrGreKey,
    NW_OUT struct in_addr* ipv4Addr, NW_OUT struct in6_addr* pIpv6Addr) {
  nw_gtpv2c_msg_t* thiz = (nw_gtpv2c_msg_t*) hMsg;
  nw_gtpv2c_ie_tlv_t* pIe;

  NW_ASSERT(instance <= NW_GTPV2C_IE_INSTANCE_MAXIMUM);

  if (thiz->isIeValid[NW_GTPV2C_IE_FTEID][instance]) {
    pIe = (nw_gtpv2c_ie_tlv_t*) thiz->pIe[NW_GTPV2C_IE_FTEID][instance];
    uint8_t flags;
    uint8_t* pIeValue = ((uint8_t*) pIe) + 4;

    flags   = (*pIeValue) & 0xE0;
    *ifType = (*pIeValue) & 0x1F;
    pIeValue += 1;
    *teidOrGreKey = ntohl(*((uint32_t*) (pIeValue)));
    pIeValue += 4;

    if (flags & 0x80) {
      ipv4Addr->s_addr = (*((uint32_t*) (pIeValue)));
      pIeValue += 4;
    }
    return NW_OK;
  }

  return NW_GTPV2C_IE_MISSING;
}

nw_rc_t nwGtpv2cMsgHexDump(nw_gtpv2c_msg_handle_t hMsg, FILE* fp) {
  nw_gtpv2c_msg_t* pMsg = (nw_gtpv2c_msg_t*) hMsg;
  uint8_t* data         = pMsg->msgBuf;
  uint32_t size         = pMsg->msgLen;
  unsigned char* p      = (unsigned char*) data;
  unsigned char c;
  int n;
  char bytestr[4]          = {0};
  char addrstr[10]         = {0};
  char hexstr[16 * 3 + 5]  = {0};
  char charstr[16 * 1 + 5] = {0};
  fprintf((FILE*) fp, "\n");

  for (n = 1; n <= size; n++) {
    if (n % 16 == 1) {
      // store address for this line
      snprintf(
          addrstr, sizeof(addrstr), "%.4lx",
          ((unsigned long) p - (unsigned long) data));
    }

    c = *p;

    if (isalnum(c) == 0) {
      c = '.';
    }

    // store hex str (for left side)
    snprintf(bytestr, sizeof(bytestr), "%02X ", *p);
    strncat(hexstr, bytestr, sizeof(hexstr) - strlen(hexstr) - 1);
    /*
     * store char str (for right side)
     */
    snprintf(bytestr, sizeof(bytestr), "%c", c);
    strncat(charstr, bytestr, sizeof(charstr) - strlen(charstr) - 1);

    if (n % 16 == 0) {
      // line completed
      fprintf((FILE*) fp, "[%4.4s]   %-50.50s  %s\n", addrstr, hexstr, charstr);
      hexstr[0]  = 0;
      charstr[0] = 0;
    } else if (n % 8 == 0) {
      // half line: add whitespaces
      strncat(hexstr, "  ", sizeof(hexstr) - strlen(hexstr) - 1);
      strncat(charstr, " ", sizeof(charstr) - strlen(charstr) - 1);
    }

    p++; /* next byte */
  }

  if (strlen(hexstr) > 0) {
    // print rest of buffer if not empty
    fprintf((FILE*) fp, "[%4.4s]   %-50.50s  %s\n", addrstr, hexstr, charstr);
  }

  fprintf((FILE*) fp, "\n");
  return NW_OK;
}

#ifdef __cplusplus
}
#endif

/*--------------------------------------------------------------------------*
                            E N D   O F   F I L E
  --------------------------------------------------------------------------*/
