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
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <stdbool.h>

#include "bstrlib.h"

#include "NwTypes.h"
#include "NwUtils.h"
#include "NwError.h"
#include "NwGtpv2cPrivate.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cTrxn.h"
#include "NwGtpv2cLog.h"
#include "dynamic_memory_check.h"
#include "gcc_diag.h"
#include "log.h"

#ifdef _NWGTPV2C_HAVE_TIMERADD
#define NW_GTPV2C_TIMER_ADD(tvp, uvp, vvp) timeradd((tvp), (uvp), (vvp))
#define NW_GTPV2C_TIMER_SUB(tvp, uvp, vvp) timersub((tvp), (uvp), (vvp))
#define NW_GTPV2C_TIMER_CMP_P(a, b, CMP) timercmp(a, b, CMP)
#else

#define NW_GTPV2C_TIMER_ADD(tvp, uvp, vvp)                                     \
  do {                                                                         \
    (vvp)->tv_sec  = (tvp)->tv_sec + (uvp)->tv_sec;                            \
    (vvp)->tv_usec = (tvp)->tv_usec + (uvp)->tv_usec;                          \
    if ((vvp)->tv_usec >= 1000000) {                                           \
      (vvp)->tv_sec++;                                                         \
      (vvp)->tv_usec -= 1000000;                                               \
    }                                                                          \
  } while (0)

#define NW_GTPV2C_TIMER_SUB(tvp, uvp, vvp)                                     \
  do {                                                                         \
    (vvp)->tv_sec  = (tvp)->tv_sec - (uvp)->tv_sec;                            \
    (vvp)->tv_usec = (tvp)->tv_usec - (uvp)->tv_usec;                          \
    if ((vvp)->tv_usec < 0) {                                                  \
      (vvp)->tv_sec--;                                                         \
      (vvp)->tv_usec += 1000000;                                               \
    }                                                                          \
  } while (0)

#define NW_GTPV2C_TIMER_CMP_P(a, b, CMP)                                       \
  (((a)->tv_sec == (b)->tv_sec) ? ((a)->tv_usec CMP(b)->tv_usec) :             \
                                  ((a)->tv_sec CMP(b)->tv_sec))

#endif

#define NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(__thiz, __msgType)                    \
  do {                                                                         \
    __thiz->pGtpv2cMsgIeParseInfo[__msgType] = nwGtpv2cMsgIeParseInfoNew(      \
        (nw_gtpv2c_stack_handle_t) __thiz, __msgType);                         \
  } while (0)

#define NW_GTPV2C_UDP_PORT (2123)

#ifdef __cplusplus
extern "C" {
#endif

static nw_gtpv2c_timeout_info_t* gpGtpv2cTimeoutInfoPool = NULL;

typedef struct {
  int currSize;
  int maxSize;
  nw_gtpv2c_timeout_info_t** pHeap;
} NwGtpv2cTmrMinHeapT;

#define NW_HEAP_PARENT_INDEX(__child) (((__child) -1) / 2)

NwGtpv2cTmrMinHeapT* nwGtpv2cTmrMinHeapNew(int maxSize) {
  NwGtpv2cTmrMinHeapT* thiz =
      (NwGtpv2cTmrMinHeapT*) malloc(sizeof(NwGtpv2cTmrMinHeapT));

  if (thiz) {
    thiz->currSize = 0;
    thiz->maxSize  = maxSize;
    thiz->pHeap    = (nw_gtpv2c_timeout_info_t**) malloc(
        maxSize * sizeof(nw_gtpv2c_timeout_info_t*));
  }

  return thiz;
}

void nwGtpv2cTmrMinHeapDelete(NwGtpv2cTmrMinHeapT* thiz) {
  free_wrapper((void**) &thiz->pHeap);
  free_wrapper((void**) &thiz);
}

static nw_rc_t nwGtpv2cTmrMinHeapInsert(
    NwGtpv2cTmrMinHeapT* thiz, nw_gtpv2c_timeout_info_t* pTimerEvent) {
  int holeIndex = thiz->currSize++;

  NW_ASSERT(thiz->currSize < thiz->maxSize);

  while ((holeIndex > 0) &&
         NW_GTPV2C_TIMER_CMP_P(
             &(thiz->pHeap[NW_HEAP_PARENT_INDEX(holeIndex)])->tvTimeout,
             &(pTimerEvent->tvTimeout), >)) {
    thiz->pHeap[holeIndex] = thiz->pHeap[NW_HEAP_PARENT_INDEX(holeIndex)];
    thiz->pHeap[holeIndex]->timerMinHeapIndex = holeIndex;
    holeIndex                                 = NW_HEAP_PARENT_INDEX(holeIndex);
  }

  thiz->pHeap[holeIndex]         = pTimerEvent;
  pTimerEvent->timerMinHeapIndex = holeIndex;
  return holeIndex;
}

#define NW_MIN_HEAP_INDEX_INVALID (0xFFFFFFFF)
static nw_rc_t nwGtpv2cTmrMinHeapRemove(
    NwGtpv2cTmrMinHeapT* thiz, int minHeapIndex) {
  nw_gtpv2c_timeout_info_t* pTimerEvent = NULL;
  int holeIndex                         = minHeapIndex;
  int minChild = 0, maxChild = 0;

  if (minHeapIndex == NW_MIN_HEAP_INDEX_INVALID) return NW_FAILURE;

  if (minHeapIndex < thiz->currSize) {
    thiz->pHeap[minHeapIndex]->timerMinHeapIndex = NW_MIN_HEAP_INDEX_INVALID;
    thiz->currSize--;
    pTimerEvent = thiz->pHeap[thiz->currSize];
    holeIndex   = minHeapIndex;
    minChild    = (2 * holeIndex) + 1;
    maxChild    = minChild + 1;

    while ((maxChild) <= thiz->currSize) {
      if (NW_GTPV2C_TIMER_CMP_P(
              &(thiz->pHeap[minChild]->tvTimeout),
              &(thiz->pHeap[maxChild]->tvTimeout), >))
        minChild = maxChild;

      if (NW_GTPV2C_TIMER_CMP_P(
              &(pTimerEvent->tvTimeout), &(thiz->pHeap[minChild]->tvTimeout),
              <)) {
        break;
      }

      thiz->pHeap[holeIndex]                    = thiz->pHeap[minChild];
      thiz->pHeap[holeIndex]->timerMinHeapIndex = holeIndex;
      holeIndex                                 = minChild;
      minChild                                  = (2 * holeIndex) + 1;
      maxChild                                  = minChild + 1;
    }

    while ((holeIndex > 0) &&
           NW_GTPV2C_TIMER_CMP_P(
               &((thiz->pHeap[NW_HEAP_PARENT_INDEX(holeIndex)])->tvTimeout),
               &(pTimerEvent->tvTimeout), >)) {
      thiz->pHeap[holeIndex] = thiz->pHeap[NW_HEAP_PARENT_INDEX(holeIndex)];
      thiz->pHeap[holeIndex]->timerMinHeapIndex = holeIndex;
      holeIndex = NW_HEAP_PARENT_INDEX(holeIndex);
    }

    if (holeIndex < thiz->currSize) {
      thiz->pHeap[holeIndex]         = pTimerEvent;
      pTimerEvent->timerMinHeapIndex = holeIndex;
    }

    thiz->pHeap[thiz->currSize] = NULL;
    return NW_OK;
  }

  return NW_FAILURE;
}

static nw_gtpv2c_timeout_info_t* nwGtpv2cTmrMinHeapPeek(
    NwGtpv2cTmrMinHeapT* thiz) {
  if (thiz->currSize) {
    return thiz->pHeap[0];
  }

  return NULL;
}
/*--------------------------------------------------------------------------*
                      P R I V A T E    F U N C T I O N S
  --------------------------------------------------------------------------*/

static void nwGtpv2cDisplayBanner(nw_gtpv2c_stack_t* thiz) {
#if DISPLAY_LICENCE_INFO
  OAILOG_INFO(
      LOG_GTPV2C,
      " *----------------------------------------------------------------------"
      "------*\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                             n w - g t p v 2 c                        "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *    G P R S    T u n n e l i n g    P r o t o c o l   v 2 c   S t a c "
      "k     *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * Copyright (c) 2010-2011 Amit Chawre                                  "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * All rights reserved.                                                 "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * Redistribution and use in source and binary forms, with or without   "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * modification, are permitted provided that the following conditions   "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * are met:                                                             "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * 1. Redistributions of source code must retain the above copyright    "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *    notice, this list of conditions and the following disclaimer.     "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * 2. Redistributions in binary form must reproduce the above copyright "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *    notice, this list of conditions and the following disclaimer in "
      "the     *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *    documentation and/or other materials provided with the "
      "distribution.    *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * 3. The name of the author may not be used to endorse or promote "
      "products   *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *    derived from this software without specific prior written "
      "permission.   *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *                                                                      "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED "
      "WARRANTIES  *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE "
      "DISCLAIMED.    *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,     "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, "
      "BUT   *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF "
      "USE,  *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON "
      "ANY      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT  "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE "
      "USE OF   *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.    "
      "      *\n");
  OAILOG_INFO(
      LOG_GTPV2C,
      " *----------------------------------------------------------------------"
      "------*\n\n");
#endif
}

/*---------------------------------------------------------------------------
   Tunnel RBTree Search Data Structure
  --------------------------------------------------------------------------*/

/**
  Comparator funtion for comparing two sequence number transactions.

  @param[in] a: Pointer to session a.
  @param[in] b: Pointer to session b.
  @return  An integer greater than, equal to or less than zero according to
  whether the object pointed to by a is greater than, equal to or less than the
  object pointed to by b.
*/

static inline int32_t nwGtpv2cCompareTunnel(
    struct nw_gtpv2c_tunnel_s* a, struct nw_gtpv2c_tunnel_s* b) {
  if (a->teid > b->teid) return 1;

  if (a->teid < b->teid) return -1;

  /** Compare the sa_family. */
  if (((struct sockaddr*) &a->ipAddrRemote)->sa_family >
      ((struct sockaddr*) &b->ipAddrRemote)->sa_family)
    return 1;

  if (((struct sockaddr*) &a->ipAddrRemote)->sa_family <
      ((struct sockaddr*) &b->ipAddrRemote)->sa_family)
    return -1;

  /** compare the address. */
  if (((struct sockaddr*) &a->ipAddrRemote)->sa_family == AF_INET) {
    if (((struct sockaddr_in*) &a->ipAddrRemote)->sin_addr.s_addr >
        ((struct sockaddr_in*) &b->ipAddrRemote)->sin_addr.s_addr)
      return 1;

    if (((struct sockaddr_in*) &a->ipAddrRemote)->sin_addr.s_addr <
        ((struct sockaddr_in*) &b->ipAddrRemote)->sin_addr.s_addr)
      return -1;
  } else {
    DevAssert(((struct sockaddr*) &a->ipAddrRemote)->sa_family == AF_INET6);
    /** Should return 1 if a is bigger. */
    return memcmp(
        ((struct sockaddr_in6*) &a->ipAddrRemote)->sin6_addr.s6_addr,
        ((struct sockaddr_in6*) &b->ipAddrRemote)->sin6_addr.s6_addr, 16);
  }

  return 0;
}

RB_GENERATE(
    NwGtpv2cTunnelMap, nw_gtpv2c_tunnel_s, tunnelMapRbtNode,
    nwGtpv2cCompareTunnel)

/*---------------------------------------------------------------------------
   Transaction RBTree Search Data Structure
  --------------------------------------------------------------------------*/
/**
  Comparator funtion for comparing two outstancing TX transactions.

  @param[in] a: Pointer to session a.
  @param[in] b: Pointer to session b.
  @return  An integer greater than, equal to or less than zero according to
  whether the object pointed to by a is greater than, equal to or less than the
  object pointed to by b.
*/
static inline int32_t nwGtpv2cCompareOutstandingTxSeqNumTrxn(
    struct nw_gtpv2c_trxn_s* a, struct nw_gtpv2c_trxn_s* b) {
  if (a->seqNum > b->seqNum) return 1;

  if (a->seqNum < b->seqNum) return -1;

  /** Compare the sa_family. */
  if (((struct sockaddr*) &a->peer_ip)->sa_family >
      ((struct sockaddr*) &b->peer_ip)->sa_family)
    return 1;

  if (((struct sockaddr*) &a->peer_ip)->sa_family <
      ((struct sockaddr*) &b->peer_ip)->sa_family)
    return -1;

  /** compare the address. */
  if (((struct sockaddr*) &a->peer_ip)->sa_family == AF_INET) {
    if (((struct sockaddr_in*) &a->peer_ip)->sin_addr.s_addr >
        ((struct sockaddr_in*) &b->peer_ip)->sin_addr.s_addr)
      return 1;

    if (((struct sockaddr_in*) &a->peer_ip)->sin_addr.s_addr <
        ((struct sockaddr_in*) &b->peer_ip)->sin_addr.s_addr)
      return -1;
  } else {
    DevAssert(((struct sockaddr*) &a->peer_ip)->sa_family == AF_INET6);
    /** Should return 1 if a is bigger. */
    return memcmp(
        ((struct sockaddr_in6*) &a->peer_ip)->sin6_addr.s6_addr,
        ((struct sockaddr_in6*) &b->peer_ip)->sin6_addr.s6_addr, 16);
  }

  return 0;
}

RB_GENERATE(
    NwGtpv2cOutstandingTxSeqNumTrxnMap, nw_gtpv2c_trxn_s,
    outstandingTxSeqNumMapRbtNode, nwGtpv2cCompareOutstandingTxSeqNumTrxn)

/**
  Comparator funtion for comparing outstanding RX transactions.

  @param[in] a: Pointer to session a.
  @param[in] b: Pointer to session b.
  @return  An integer greater than, equal to or less than zero according to
  whether the object pointed to by a is greater than, equal to or less than the
  object pointed to by b.
*/
static inline int32_t nwGtpv2cCompareOutstandingRxSeqNumTrxn(
    struct nw_gtpv2c_trxn_s* a, struct nw_gtpv2c_trxn_s* b) {
  if (a->seqNum > b->seqNum) return 1;

  if (a->seqNum < b->seqNum) return -1;
  /** Compare the sa_family. */
  if (((struct sockaddr*) &a->peer_ip)->sa_family >
      ((struct sockaddr*) &b->peer_ip)->sa_family)
    return 1;

  if (((struct sockaddr*) &a->peer_ip)->sa_family <
      ((struct sockaddr*) &b->peer_ip)->sa_family)
    return -1;

  /** compare the address. */
  if (((struct sockaddr*) &a->peer_ip)->sa_family == AF_INET) {
    if (((struct sockaddr_in*) &a->peer_ip)->sin_addr.s_addr >
        ((struct sockaddr_in*) &b->peer_ip)->sin_addr.s_addr)
      return 1;

    if (((struct sockaddr_in*) &a->peer_ip)->sin_addr.s_addr <
        ((struct sockaddr_in*) &b->peer_ip)->sin_addr.s_addr)
      return -1;
  } else {
    DevAssert(((struct sockaddr*) &a->peer_ip)->sa_family == AF_INET6);
    /** Should return 1 if a is bigger. */
    return memcmp(
        ((struct sockaddr_in6*) &a->peer_ip)->sin6_addr.s6_addr,
        ((struct sockaddr_in6*) &b->peer_ip)->sin6_addr.s6_addr, 16);
  }

  if (a->peerPort > b->peerPort) return 1;

  if (a->peerPort < b->peerPort) return -1;

  return 0;
}

RB_GENERATE(
    NwGtpv2cOutstandingRxSeqNumTrxnMap, nw_gtpv2c_trxn_s,
    outstandingRxSeqNumMapRbtNode, nwGtpv2cCompareOutstandingRxSeqNumTrxn)

/*---------------------------------------------------------------------------
   Timer RB-tree data structure.
  --------------------------------------------------------------------------*/
/**
  Comparator funtion for comparing two outstancing TX transactions.

  @param[in] a: Pointer to session a.
  @param[in] b: Pointer to session b.
  @return  An integer greater than, equal to or less than zero according to
  whether the object pointed to by a is greater than, equal to or less than the
  object pointed to by b.
*/
static inline int32_t nwGtpv2cCompareOutstandingTxRexmitTime(
    struct nw_gtpv2c_timeout_info_s* a, struct nw_gtpv2c_timeout_info_s* b) {
  if (NW_GTPV2C_TIMER_CMP_P(&a->tvTimeout, &b->tvTimeout, >)) return 1;

  if (NW_GTPV2C_TIMER_CMP_P(&a->tvTimeout, &b->tvTimeout, <)) return -1;

  return 0;
}

RB_GENERATE(
    NwGtpv2cActiveTimerList, nw_gtpv2c_timeout_info_s, activeTimerListRbtNode,
    nwGtpv2cCompareOutstandingTxRexmitTime)

/**
   Send msg to peer via data request to UDP Entity

   @param[in] thiz : Pointer to stack.
   @param[in] peerIp : Peer Ip address.
   @param[in] peerPort : Local Ip port to send the message from.
   @param[in] peerPort : Peer Ip port.
   @param[in] pMsg : Message to be sent.
   @return NW_OK on success.
*/
static nw_rc_t nwGtpv2cCreateAndSendMsg(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t seqNum,
    NW_IN uint32_t localPort, NW_IN struct sockaddr* peerIp,
    NW_IN uint32_t peerPort, NW_IN nw_gtpv2c_msg_t* pMsg) {
  nw_rc_t rc      = NW_FAILURE;
  uint8_t* msgHdr = NULL;

  NW_ASSERT(thiz);
  NW_ASSERT(pMsg);
  msgHdr = pMsg->msgBuf;
  // Set flags in header
  *(msgHdr++) = (pMsg->version << 5) | (pMsg->teidPresent << 3);
  // Set msg type in header
  *(msgHdr++) = (pMsg->msgType);
  // Set msg length in header
  *((uint16_t*) msgHdr) = htons(pMsg->msgLen - 4);
  msgHdr += 2;
  // Set TEID, if present in header
  if (pMsg->teidPresent) {
    *((uint32_t*) msgHdr) = htonl(pMsg->teid);
    msgHdr += 4;
  }
  // Set seq num in header
  *((uint32_t*) msgHdr) = htonl(seqNum << 8);
  // Call UDP data request callback
  NW_ASSERT(thiz->udp.udpDataReqCallback != NULL);
  rc = thiz->udp.udpDataReqCallback(
      thiz->udp.hUdp, pMsg->msgBuf, pMsg->msgLen, localPort, peerIp, peerPort);
  NW_ASSERT(NW_OK == rc);
  return rc;
}

/**
  Send an Version Not Supported message

  @param[in] thiz : Stack pointer
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cSendVersionNotSupportedInd(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN struct sockaddr* peerIp,
    NW_IN uint32_t peerPort, NW_IN uint32_t seqNum) {
  nw_rc_t rc                  = NW_FAILURE;
  nw_gtpv2c_msg_handle_t hMsg = 0;

  rc = nwGtpv2cMsgNew(
      (nw_gtpv2c_stack_handle_t) thiz, false, NW_GTP_VERSION_NOT_SUPPORTED_IND,
      0x00, seqNum, (&hMsg));
  NW_ASSERT(NW_OK == rc);
  // todo:   OAILOG_NOTICE (LOG_GTPV2C,  "Sending Version Not Supported
  // Indication message to %x:%x with seq %u\n", peerIp->s_addr, peerPort,
  // seqNum);
  rc = nwGtpv2cCreateAndSendMsg(
      thiz, seqNum, NW_GTPV2C_UDP_PORT, peerIp, peerPort,
      (nw_gtpv2c_msg_t*) hMsg);
  rc = nwGtpv2cMsgDelete((nw_gtpv2c_stack_handle_t) thiz, hMsg);
  NW_ASSERT(NW_OK == rc);
  return rc;
}

/**
  Create a local tunnel.

  @param[in] thiz : Stack pointer
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cCreateLocalTunnel(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t teid,
    NW_IN struct sockaddr* fa, NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel,
    NW_OUT nw_gtpv2c_tunnel_handle_t* phTunnel) {
  nw_rc_t rc                  = NW_FAILURE;
  nw_gtpv2c_tunnel_t *pTunnel = NULL, *pCollision = NULL;

  char ip[INET_ADDRSTRLEN];

  if (fa->sa_family == AF_INET) {
    inet_ntop(
        AF_INET, &((struct sockaddr_in*) fa)->sin_addr, ip, INET_ADDRSTRLEN);
    OAILOG_DEBUG(
        LOG_GTPV2C, "Creating local tunnel with teid '0x%x' and peer IPv4 %s\n",
        teid, ip);
  } else {
    inet_ntop(
        AF_INET6, &((struct sockaddr_in6*) fa)->sin6_addr, ip,
        INET6_ADDRSTRLEN);
    OAILOG_DEBUG(
        LOG_GTPV2C, "Creating local tunnel with teid '0x%x' and peer IPv6 %s\n",
        teid, ip);
  }

  OAILOG_FUNC_IN(LOG_GTPV2C);
  pTunnel = nwGtpv2cTunnelNew(thiz, teid, fa, hUlpTunnel);

  if (pTunnel) {
    pCollision = RB_INSERT(NwGtpv2cTunnelMap, &(thiz->tunnelMap), pTunnel);

    if (pCollision) {
      rc = nwGtpv2cTunnelDelete(thiz, pTunnel);
      NW_ASSERT(NW_OK == rc);
      *phTunnel = (nw_gtpv2c_tunnel_handle_t) 0;
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Local tunnel creation failed for teid '0x%x' and peer IP %s. Tunnel "
          "already exists!\n",
          teid, ip);
      OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_FAILURE);
    }
  } else {
    rc = NW_FAILURE;
  }

  *phTunnel = (nw_gtpv2c_tunnel_handle_t) pTunnel;
  OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
}

/**
  Delete a local tunnel.

  @param[in] thiz : Stack pointer
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cDeleteLocalTunnel(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_OUT nw_gtpv2c_tunnel_handle_t hTunnel) {
  nw_rc_t rc                  = NW_FAILURE;
  nw_gtpv2c_tunnel_t* pTunnel = (nw_gtpv2c_tunnel_t*) hTunnel;
  char ip[INET6_ADDRSTRLEN];

  OAILOG_FUNC_IN(LOG_GTPV2C);

  pTunnel = RB_REMOVE(
      NwGtpv2cTunnelMap, &(thiz->tunnelMap), (nw_gtpv2c_tunnel_t*) hTunnel);
  NW_ASSERT(pTunnel == (nw_gtpv2c_tunnel_t*) hTunnel);

  inet_ntop(
      ((struct sockaddr*) &pTunnel->ipAddrRemote)->sa_family,
      (void*) &pTunnel->ipAddrRemote, ip,
      ((struct sockaddr*) &pTunnel->ipAddrRemote)->sa_family == AF_INET ?
          INET_ADDRSTRLEN :
          INET6_ADDRSTRLEN);
  OAILOG_DEBUG(
      LOG_GTPV2C, "Deleting local tunnel with teid '0x%x' and peer IP %s\n",
      pTunnel->teid, ip);
  rc = nwGtpv2cTunnelDelete(thiz, pTunnel);
  NW_ASSERT(NW_OK == rc);

  OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
}

/*---------------------------------------------------------------------------
   ULP API Processing Functions
  --------------------------------------------------------------------------*/

/**
  Process NW_GTPV2C_ULP_API_INITIAL_REQ Request from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpInitialReq(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  nw_rc_t rc                       = NW_FAILURE;
  nw_gtpv2c_trxn_t* pTrxn          = NULL;
  nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};

  OAILOG_FUNC_IN(LOG_GTPV2C);

  // Create New Transaction
  pTrxn = nwGtpv2cTrxnNew(thiz);

  // For MME this is the only place to create S11 tunnel!
  if (pTrxn) {
    if (!pUlpReq->u_api_info.initialReqInfo.hTunnel) {
      /** Check if a tunnel already exists depending on the flag. */
      keyTunnel.teid = pUlpReq->u_api_info.initialReqInfo.teidLocal;

      memcpy(
          ((struct sockaddr*) &keyTunnel.ipAddrRemote),
          pUlpReq->u_api_info.initialReqInfo.edns_peer_ip,
          pUlpReq->u_api_info.initialReqInfo.edns_peer_ip->sa_family ==
                  AF_INET ?
              sizeof(struct sockaddr_in) :
              sizeof(struct sockaddr_in6));

      pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(thiz->tunnelMap), &keyTunnel);
      if (!pLocalTunnel) {
        pLocalTunnel = RB_MIN(NwGtpv2cTunnelMap, &(thiz->tunnelMap));
        OAILOG_WARNING(
            LOG_GTPV2C,
            "Request message received on non-existent teid 0x%x received! "
            "Creating new tunnel.\n",
            ntohl(pUlpReq->u_api_info.initialReqInfo.teidLocal));
        rc = nwGtpv2cCreateLocalTunnel(
            thiz, pUlpReq->u_api_info.initialReqInfo.teidLocal,
            pUlpReq->u_api_info.initialReqInfo.edns_peer_ip,
            pUlpReq->u_api_info.initialReqInfo.hUlpTunnel,
            &pUlpReq->u_api_info.initialReqInfo.hTunnel);
        NW_ASSERT(NW_OK == rc);
      } else {
        pUlpReq->u_api_info.initialReqInfo.hTunnel =
            (nw_gtpv2c_tunnel_handle_t) pLocalTunnel;
      }
    }
    pTrxn->pMsg     = (nw_gtpv2c_msg_t*) pUlpReq->hMsg;
    pTrxn->hTunnel  = pUlpReq->u_api_info.initialReqInfo.hTunnel;
    pTrxn->hUlpTrxn = pUlpReq->u_api_info.initialReqInfo.hUlpTrxn;
    /** This will stay. */

    memcpy(
        (void*) &pTrxn->peer_ip,
        pUlpReq->u_api_info.initialReqInfo.edns_peer_ip,
        (pUlpReq->u_api_info.initialReqInfo.edns_peer_ip->sa_family ==
         AF_INET) ?
            sizeof(struct sockaddr_in) :
            sizeof(struct sockaddr_in6));

    pTrxn->peerPort =
        NW_GTPV2C_UDP_PORT; /**< Initial Requests always to 2123. */
    /* No Delete. */
    pTrxn->teidLocal = pUlpReq->u_api_info.initialReqInfo.teidLocal;
    pTrxn->noDelete  = pUlpReq->u_api_info.initialReqInfo.noDelete;
    pTrxn->trx_flags = pUlpReq->u_api_info.initialReqInfo.internal_flags;
    pTrxn->localPort = 0; /**< Set the local port to 0 (initialize it). */
    if (pUlpReq->apiType & NW_GTPV2C_ULP_API_FLAG_IS_COMMAND_MESSAGE) {
      pTrxn->seqNum |= 0x00800000UL;
    }

    char peer_ip[INET_ADDRSTRLEN];
    inet_ntop(
        AF_INET, (void*) &pTrxn->peer_ip.addrv4.sin_addr, peer_ip,
        INET_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_GTPV2C, "peer IP information %s\n", peer_ip);

    rc = nwGtpv2cCreateAndSendMsg(
        thiz, pTrxn->seqNum, 0, (struct sockaddr*) &pTrxn->peer_ip,
        pTrxn->peerPort,
        pTrxn->pMsg); /**< Send it from the socket with the high port. */

    if (NW_OK == rc) {
      rc = nwGtpv2cTrxnStartPeerRspWaitTimer(pTrxn);  // Start guard timer
      NW_ASSERT(NW_OK == rc);

      // Insert into search tree

      pTrxn = RB_INSERT(
          NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
          pTrxn);
      NW_ASSERT(pTrxn == NULL);
    } else {
      rc = nwGtpv2cTrxnDelete(&pTrxn);
      NW_ASSERT(NW_OK == rc);
    }
  } else {
    rc = NW_FAILURE;
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Process NW_GTPV2C_ULP_API_TRIGGERED_REQ Request from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpTriggeredReq(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  nw_rc_t rc                 = NW_FAILURE;
  nw_gtpv2c_trxn_t* pTrxn    = NULL;
  nw_gtpv2c_trxn_t* pReqTrxn = NULL;

  OAILOG_FUNC_IN(LOG_GTPV2C);

  // Create New Transaction
  pTrxn = nwGtpv2cTrxnWithSeqNumNew(
      thiz, (((nw_gtpv2c_msg_t*) (pUlpReq->hMsg))->seqNum));

  if (pTrxn) {
    pReqTrxn = (nw_gtpv2c_trxn_t*) pUlpReq->u_api_info.triggeredReqInfo.hTrxn;
    pTrxn->hUlpTrxn = pUlpReq->u_api_info.triggeredReqInfo.hUlpTrxn;
    memcpy(
        (void*) &pTrxn->peer_ip,
        pUlpReq->u_api_info.initialReqInfo.edns_peer_ip,
        (((struct sockaddr*) &pReqTrxn->peer_ip)->sa_family == AF_INET) ?
            sizeof(struct sockaddr_in) :
            sizeof(struct sockaddr_in6));

    pTrxn->peerPort = pReqTrxn->peerPort;
    pTrxn->pMsg     = (nw_gtpv2c_msg_t*) pUlpReq->hMsg;
    rc              = nwGtpv2cCreateAndSendMsg(
        thiz, pTrxn->seqNum, NW_GTPV2C_UDP_PORT,
        (struct sockaddr*) &pTrxn->peer_ip, pTrxn->peerPort, pTrxn->pMsg);

    if (NW_OK == rc) {
      // Start guard timer
      rc = nwGtpv2cTrxnStartPeerRspWaitTimer(pTrxn);
      NW_ASSERT(NW_OK == rc);

      // Insert into search tree

      RB_INSERT(
          NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
          pTrxn);

      if (!pUlpReq->u_api_info.triggeredReqInfo.hTunnel) {
        rc = nwGtpv2cCreateLocalTunnel(
            thiz, pUlpReq->u_api_info.triggeredReqInfo.teidLocal,
            (struct sockaddr*) &pReqTrxn->peer_ip,
            pUlpReq->u_api_info.triggeredReqInfo.hUlpTunnel,
            &pUlpReq->u_api_info.triggeredReqInfo.hTunnel);
      }
    } else {
      rc = nwGtpv2cTrxnDelete(&pTrxn);
      NW_ASSERT(NW_OK == rc);
    }
  } else {
    rc = NW_FAILURE;
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Process NW_GTPV2C_ULP_API_TRIGGERED_RSP Request from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpTriggeredRsp(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpRsp) {
  nw_rc_t rc                 = NW_FAILURE;
  nw_gtpv2c_trxn_t* pReqTrxn = NULL;

  OAILOG_FUNC_IN(LOG_GTPV2C);
  pReqTrxn = (nw_gtpv2c_trxn_t*) pUlpRsp->u_api_info.triggeredRspInfo.hTrxn;
  NW_ASSERT(pReqTrxn != NULL);

  if (((nw_gtpv2c_msg_t*) pUlpRsp->hMsg)->seqNum == 0)
    ((nw_gtpv2c_msg_t*) pUlpRsp->hMsg)->seqNum = pReqTrxn->seqNum;

  OAILOG_DEBUG(
      LOG_GTPV2C, "Sending response message over seq '0x%x'\n",
      pReqTrxn->seqNum);
  rc = nwGtpv2cCreateAndSendMsg(
      thiz, pReqTrxn->seqNum, pReqTrxn->localPort,
      (struct sockaddr*) &pReqTrxn->peer_ip, pReqTrxn->peerPort,
      (nw_gtpv2c_msg_t*) pUlpRsp->hMsg);
  /** Depending on the cause type, add it or not. */
  if (pUlpRsp->u_api_info.triggeredRspInfo.pt_trx) {
    OAILOG_DEBUG(
        LOG_GTPV2C,
        "Making transaction with seq '0x%x' (%p) passtrough due temporary "
        "reject. Not continuing with message. \n",
        pReqTrxn->seqNum, pReqTrxn);
    /** We wan't this message to be able to trigger something, so we remove the
     * message, too. */
    if (pReqTrxn->pMsg) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Transaction (%d) with seq '%p'contained a previous message. "
          "Discarding first. \n",
          pReqTrxn->seqNum, pReqTrxn);
      rc = nwGtpv2cMsgDelete(
          (nw_gtpv2c_stack_handle_t) thiz,
          (nw_gtpv2c_msg_handle_t) pReqTrxn->pMsg);
      NW_ASSERT(NW_OK == rc);
    }
    /** Set the transaction as pt. */
    pReqTrxn->pt_trx = true;
    OAILOG_DEBUG(
        LOG_GTPV2C,
        "Removing the message (%p) for passthrough trx with seqNo %d. \n",
        (void*) pUlpRsp->hMsg, pReqTrxn->seqNum);
    rc = nwGtpv2cMsgDelete((nw_gtpv2c_stack_handle_t) thiz, pUlpRsp->hMsg);
    OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
  } else {
    /** Check if there was a message, remove it first. */
    if (pReqTrxn->pMsg) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Transaction (%p) with seq '0x%x' contained a previous message. "
          "Discarding first. \n",
          pReqTrxn, pReqTrxn->seqNum);
      rc = nwGtpv2cMsgDelete(
          (nw_gtpv2c_stack_handle_t) thiz,
          (nw_gtpv2c_msg_handle_t) pReqTrxn->pMsg);
      NW_ASSERT(NW_OK == rc);
    }
    pReqTrxn->pMsg = (nw_gtpv2c_msg_t*) pUlpRsp->hMsg;
  }
  rc = nwGtpv2cTrxnStartDulpicateRequestWaitTimer(pReqTrxn);

  /** Creating a local tunnel if flag is set. */
  if ((pUlpRsp->apiType & 0xFF000000) ==
      NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL) {
    /** Check if there is a local tunnel already existing, if not create a local
     * S10 tunnel. */
    nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};
    keyTunnel.teid = pUlpRsp->u_api_info.triggeredRspInfo.teidLocal;
    memcpy(
        ((struct sockaddr*) &keyTunnel.ipAddrRemote),
        ((struct sockaddr*) &pReqTrxn->peer_ip),
        ((struct sockaddr*) &pReqTrxn->peer_ip)->sa_family == AF_INET ?
            sizeof(struct sockaddr_in) :
            sizeof(struct sockaddr_in6));

    pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(thiz->tunnelMap), &keyTunnel);
    char ip[INET6_ADDRSTRLEN];
    inet_ntop(
        AF_INET, (void*) &pReqTrxn->peer_ip, ip,
        ((struct sockaddr*) &pReqTrxn->peer_ip)->sa_family == AF_INET ?
            INET_ADDRSTRLEN :
            INET6_ADDRSTRLEN);
    if (!pLocalTunnel) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Triggered response not containing a tunnel. Creating one for "
          "local_teid 0x%x and peer %s!\n",
          pUlpRsp->u_api_info.triggeredRspInfo.teidLocal, ip);
      rc = nwGtpv2cCreateLocalTunnel(
          thiz, pUlpRsp->u_api_info.triggeredRspInfo.teidLocal,
          (struct sockaddr*) &pReqTrxn->peer_ip,
          pUlpRsp->u_api_info.triggeredRspInfo.hUlpTunnel,
          &pUlpRsp->u_api_info.triggeredRspInfo.hTunnel);
      NW_ASSERT(NW_OK == rc);
    } else {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Triggered response already containing a tunnel. Not creating a new "
          "one for local_teid 0x%x and peer %s.\n",
          (pUlpRsp->u_api_info.triggeredRspInfo.teidLocal), ip);
    }
  }
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Process NW_GTPV2C_ULP_API_TRIGGERED_ACK from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
 */
static nw_rc_t nwGtpv2cHandleUlpTriggeredAck(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpAck) {
  OAILOG_FUNC_IN(LOG_GTPV2C);

  nw_gtpv2c_trxn_t *pAckTrxn = NULL, keyTrxn;
  nw_rc_t rc                 = NW_FAILURE;

  /** Try to find the initial request transaction for the triggered ACK. */

  OAILOG_FUNC_IN(LOG_GTPV2C);
  ;

  keyTrxn.seqNum = ((nw_gtpv2c_msg_t*) pUlpAck->hMsg)->seqNum;
  memcpy(
      (void*) &keyTrxn.peer_ip, pUlpAck->u_api_info.triggeredAckInfo.peerIp,
      (((struct sockaddr*) pUlpAck->u_api_info.triggeredAckInfo.peerIp)
           ->sa_family == AF_INET) ?
          sizeof(struct sockaddr_in) :
          sizeof(struct sockaddr_in6));

  /** A transaction of the initial request (cmd) for the triggered request
   * should exist. */
  pAckTrxn = RB_FIND(
      NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
      &keyTrxn);

  if (pAckTrxn) {
    OAILOG_INFO(
        LOG_GTPV2C,
        "Found an initial request transaction for the triggered ACK. Appending "
        "the ACK and keeping the transaction for a while.\n");
    /** Remove the original message and append the current message. */
    DevAssert(!pAckTrxn->pMsg);

    /** Append the triggered ACK to the message. */
    pAckTrxn->pMsg         = (nw_gtpv2c_msg_t*) pUlpAck->hMsg;
    pAckTrxn->pMsg->seqNum = pAckTrxn->seqNum;

  } else {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Triggered ACK message without a matching outstanding request (req) "
        "received! Discarding.\n");
    if (pUlpAck->hMsg) {
      rc = nwGtpv2cMsgDelete((nw_gtpv2c_stack_handle_t) thiz, pUlpAck->hMsg);
      NW_ASSERT(NW_OK == rc);
    }
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  OAILOG_DEBUG(
      LOG_GTPV2C, "Sending a triggered ACK message over seq '0x%x'\n",
      pAckTrxn->seqNum);
  rc = nwGtpv2cCreateAndSendMsg(
      thiz, pAckTrxn->seqNum, pUlpAck->u_api_info.triggeredAckInfo.localPort,
      pUlpAck->u_api_info.triggeredAckInfo.peerIp,
      pUlpAck->u_api_info.triggeredAckInfo.peerPort,
      (nw_gtpv2c_msg_t*) pUlpAck->hMsg);

  /** Remove the tunnel. */
  rc = nwGtpv2cDeleteLocalTunnel(
      thiz, pUlpAck->u_api_info.triggeredAckInfo.hTunnel);

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Process NW_GTPV2C_ULP_CREATE_LOCAL_TUNNEL Request from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpCreateLocalTunnel(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  OAILOG_FUNC_IN(LOG_GTPV2C);
  nw_rc_t rc                  = NW_FAILURE;
  nw_gtpv2c_tunnel_t *pTunnel = NULL, *pCollision = NULL;
  char ipv4[INET_ADDRSTRLEN];

  inet_ntop(
      AF_INET, (void*) pUlpReq->u_api_info.createLocalTunnelInfo.peerIp, ipv4,
      INET_ADDRSTRLEN);
  OAILOG_DEBUG(
      LOG_GTPV2C, "Creating local tunnel with teid '0x%x' and peer IP %s\n",
      pUlpReq->u_api_info.createLocalTunnelInfo.teidLocal, ipv4);
  pTunnel = nwGtpv2cTunnelNew(
      thiz, pUlpReq->u_api_info.createLocalTunnelInfo.teidLocal,
      pUlpReq->u_api_info.createLocalTunnelInfo.peerIp,
      pUlpReq->u_api_info.triggeredRspInfo.hUlpTunnel);
  NW_ASSERT(pTunnel);
  pCollision = RB_INSERT(NwGtpv2cTunnelMap, &(thiz->tunnelMap), pTunnel);

  if (pCollision) {
    rc = nwGtpv2cTunnelDelete(thiz, pTunnel);
    NW_ASSERT(NW_OK == rc);
    pUlpReq->u_api_info.createLocalTunnelInfo.hTunnel =
        (nw_gtpv2c_tunnel_handle_t) 0;
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_FAILURE);
  }

  pUlpReq->u_api_info.createLocalTunnelInfo.hTunnel =
      (nw_gtpv2c_tunnel_handle_t) pTunnel;
  OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
}

/**
  Process NW_GTPV2C_ULP_DELETE_LOCAL_TUNNEL Request from ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @param[in] pUlpReq : Pointer to Ulp Req.
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpDeleteLocalTunnel(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  nw_rc_t rc = NW_FAILURE;

  OAILOG_FUNC_IN(LOG_GTPV2C);
  rc = nwGtpv2cDeleteLocalTunnel(
      thiz, pUlpReq->u_api_info.deleteLocalTunnelInfo.hTunnel);
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Send GTPv2c Initial Request Message Indication to ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cSendInitialReqIndToUlp(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_error_t* pError,
    NW_IN nw_gtpv2c_trxn_t* pTrxn, NW_IN uint32_t hUlpTunnel,
    NW_IN uint32_t msgType, NW_IN struct sockaddr* peerIp,
    NW_IN uint16_t peerPort, NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_rc_t rc = NW_FAILURE;
  nw_gtpv2c_ulp_api_t ulpApi;

  OAILOG_FUNC_IN(LOG_GTPV2C);
  ulpApi.hMsg    = hMsg;
  ulpApi.apiType = NW_GTPV2C_ULP_API_INITIAL_REQ_IND;
  ulpApi.u_api_info.initialReqIndInfo.msgType = msgType;
  ulpApi.u_api_info.initialReqIndInfo.hTrxn   = (nw_gtpv2c_trxn_handle_t) pTrxn;
  ulpApi.u_api_info.initialReqIndInfo.hUlpTunnel = hUlpTunnel;
  ulpApi.u_api_info.initialReqIndInfo.peerIp     = peerIp;
  ulpApi.u_api_info.initialReqIndInfo.peerPort   = peerPort;
  ulpApi.u_api_info.initialReqIndInfo.error      = *pError;
  rc = thiz->ulp.ulpReqCallback(thiz->ulp.hUlp, &ulpApi);
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Send GTPv2c Triggered Request Message Indication to ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cSendTriggeredReqIndToUlp(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_error_t* pError,
    NW_IN nw_gtpv2c_trxn_t* pTrxn, NW_IN uint32_t hUlpTunnel,
    NW_IN uint32_t msgType, NW_IN struct sockaddr* peerIp,
    NW_IN uint16_t peerPort, NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_rc_t rc = NW_FAILURE;
  nw_gtpv2c_ulp_api_t ulpApi;

  OAILOG_FUNC_IN(LOG_GTPV2C);
  ulpApi.hMsg    = hMsg;
  ulpApi.apiType = NW_GTPV2C_ULP_API_TRIGGERED_REQ_IND;

  ulpApi.u_api_info.triggeredReqIndInfo.msgType = msgType;
  ulpApi.u_api_info.triggeredReqIndInfo.hTrxn = (nw_gtpv2c_trxn_handle_t) pTrxn;
  ulpApi.u_api_info.triggeredReqIndInfo.hUlpTunnel = hUlpTunnel;
  ulpApi.u_api_info.triggeredReqIndInfo.error      = *pError;
  rc = thiz->ulp.ulpReqCallback(thiz->ulp.hUlp, &ulpApi);
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Send GTPv2c Triggered Response Indication to ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cSendTriggeredRspIndToUlp(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_error_t* pError,
    NW_IN uint32_t hUlpTrxn, NW_IN uint8_t* trxFlags_p,
    NW_IN uint16_t localPort, NW_IN uint16_t peerPort,
    NW_IN struct sockaddr* peerIp, NW_IN uint32_t hUlpTunnel,
    NW_IN uint32_t msgType, NW_IN bool noDelete,
    NW_IN nw_gtpv2c_msg_handle_t hMsg) {
  nw_rc_t rc = NW_FAILURE;
  nw_gtpv2c_ulp_api_t ulpApi;

  OAILOG_FUNC_IN(LOG_GTPV2C);
  ulpApi.hMsg    = hMsg;
  ulpApi.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP_IND;
  ulpApi.u_api_info.triggeredRspIndInfo.msgType    = msgType;
  ulpApi.u_api_info.triggeredRspIndInfo.hUlpTrxn   = hUlpTrxn;
  ulpApi.u_api_info.triggeredRspIndInfo.localPort  = localPort;
  ulpApi.u_api_info.triggeredRspIndInfo.peerPort   = peerPort;
  ulpApi.u_api_info.triggeredRspIndInfo.peerIp     = peerIp;
  ulpApi.u_api_info.triggeredRspIndInfo.hUlpTunnel = hUlpTunnel;
  ulpApi.u_api_info.triggeredRspIndInfo.error      = *pError;
  ulpApi.u_api_info.triggeredRspIndInfo.trx_flags  = *trxFlags_p;
  ulpApi.u_api_info.triggeredRspIndInfo.noDelete   = noDelete;
  rc          = thiz->ulp.ulpReqCallback(thiz->ulp.hUlp, &ulpApi);
  *trxFlags_p = ulpApi.u_api_info.triggeredRspIndInfo.trx_flags;
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
  Send GTPv2c Initial Request Message Indication to ULP entity.

  @param[in] hGtpcStackHandle : Stack handle
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleUlpFindLocalTunnel(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  // nw_rc_t                                   rc = NW_FAILURE;

  OAILOG_FUNC_IN(LOG_GTPV2C);

  /** Check if a tunnel already exists depending on the flag. */
  nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};
  keyTunnel.teid = pUlpReq->u_api_info.findLocalTunnelInfo.teidLocal;
  memcpy(
      (void*) &keyTunnel.ipAddrRemote,
      pUlpReq->u_api_info.findLocalTunnelInfo.edns_peer_ip,
      (((struct sockaddr*) &keyTunnel.ipAddrRemote)->sa_family == AF_INET) ?
          sizeof(struct sockaddr_in) :
          sizeof(struct sockaddr_in6));
  pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(thiz->tunnelMap), &keyTunnel);
  pUlpReq->u_api_info.findLocalTunnelInfo.hTunnel =
      (nw_gtpv2c_tunnel_handle_t) pLocalTunnel;

  if (pLocalTunnel) {
    //    todo:  OAILOG_DEBUG (LOG_GTPV2C, "FOUND local tunnel with teid '0x%x'
    //    and peer IP 0x%x\n", keyTunnel.teid, keyTunnel.ipv4AddrRemote);
  } else {
    // todo:      OAILOG_DEBUG (LOG_GTPV2C, "DID NOT FOUND local tunnel with
    // teid '0x%x' and peer IP 0x%x\n", keyTunnel.teid,
    // keyTunnel.ipv4AddrRemote);
  }
  return RETURNok;
}

/**
  Handle Echo Request from Peer Entity.

  @param[in] thiz : Stack context
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleEchoReq(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t msgType,
    NW_IN uint8_t* msgBuf, NW_IN uint32_t msgBufLen, NW_IN uint16_t peerPort,
    NW_IN struct sockaddr* peerIp) {
  nw_rc_t rc                  = NW_FAILURE;
  uint32_t seqNum             = 0;
  nw_gtpv2c_msg_handle_t hMsg = 0;

  seqNum = ntohl(*((uint32_t*) (msgBuf + (((*msgBuf) & 0x08) ? 8 : 4)))) >> 8;

  // Send Echo Response
  rc = nwGtpv2cMsgNew(
      (nw_gtpv2c_stack_handle_t) thiz, false, /* TEID present flag    */
      NW_GTP_ECHO_RSP,                        /* Msg Type             */
      0x00,                                   /* TEID                 */
      seqNum,                                 /* Seq Number           */
      (&hMsg));
  NW_ASSERT(NW_OK == rc);
  rc =
      nwGtpv2cMsgAddIeTV1(hMsg, NW_GTPV2C_IE_RECOVERY, 0, thiz->restartCounter);
  char ipv4[INET_ADDRSTRLEN];
  inet_ntop(AF_INET, (void*) peerIp, ipv4, INET_ADDRSTRLEN);
  OAILOG_ERROR(
      LOG_GTPV2C, "Sending NW_GTP_ECHO_RSP message to %s:%u with seq %u\n",
      ipv4, peerPort, (seqNum));
  rc = nwGtpv2cCreateAndSendMsg(
      thiz, (seqNum), NW_GTPV2C_UDP_PORT, peerIp, peerPort,
      (nw_gtpv2c_msg_t*) hMsg);
  rc = nwGtpv2cMsgDelete((nw_gtpv2c_stack_handle_t) thiz, hMsg);
  NW_ASSERT(NW_OK == rc);
  return rc;
}

/**
  Handle Initial Request from Peer Entity.

  @param[in] thiz : Stack context
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleInitialReq(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t msgType,
    NW_IN uint8_t* msgBuf, NW_IN uint32_t msgBufLen, NW_IN uint16_t peerPort,
    NW_IN struct sockaddr* peerIp) {
  nw_rc_t rc                       = NW_FAILURE;
  uint32_t seqNum                  = 0;
  uint32_t teidLocal               = 0;
  nw_gtpv2c_trxn_t* pTrxn          = NULL;
  nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};
  nw_gtpv2c_msg_handle_t hMsg              = 0;
  nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel = 0;
  nw_gtpv2c_error_t error                  = {0};
  char ip[INET6_ADDRSTRLEN];

  teidLocal = *((uint32_t*) (msgBuf + 4));
  inet_ntop(
      AF_INET, (void*) peerIp, ip,
      peerIp->sa_family == AF_INET ? INET_ADDRSTRLEN : INET6_ADDRSTRLEN);

  if (teidLocal) {
    keyTunnel.teid = ntohl(teidLocal);
    memcpy(
        (void*) &keyTunnel.ipAddrRemote, peerIp,
        (peerIp->sa_family == AF_INET) ? sizeof(struct sockaddr_in) :
                                         sizeof(struct sockaddr_in6));
    pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(thiz->tunnelMap), &keyTunnel);

    if (!pLocalTunnel) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Request message received on non-existent teid 0x%x from peer %s "
          "received! Discarding.\n",
          ntohl(teidLocal), ip);
      return NW_OK;
    }

    hUlpTunnel = pLocalTunnel->hUlpTunnel;
  } else {
    hUlpTunnel = 0;
  }

  seqNum = ntohl(*((uint32_t*) (msgBuf + (((*msgBuf) & 0x08) ? 8 : 4)))) >> 8;
  OAILOG_DEBUG(
      LOG_GTPV2C,
      "RECEIVED GTPV2c initial request message of type %d, length %d and "
      "seqNum 0x%x.\n",
      msgType, msgBufLen, seqNum);

  pTrxn = nwGtpv2cTrxnOutstandingRxNew(
      thiz, ntohl(teidLocal), peerIp, peerPort, (seqNum));

  if (pTrxn) {
    pTrxn->localPort = thiz->udp.gtpv2cStandardPort;
    rc               = nwGtpv2cMsgFromBufferNew(
        (nw_gtpv2c_stack_handle_t) thiz, msgBuf, msgBufLen, &(hMsg));
    NW_ASSERT(thiz->pGtpv2cMsgIeParseInfo[msgType]);
    rc = nwGtpv2cMsgIeParse(thiz->pGtpv2cMsgIeParseInfo[msgType], hMsg, &error);

    if (rc != NW_OK) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Malformed request message received on TEID %u from peer %s. "
          "Notifying ULP.\n",
          ntohl(teidLocal), ip);
    }

    rc = nwGtpv2cSendInitialReqIndToUlp(
        thiz, &error, pTrxn, hUlpTunnel, msgType, peerIp, peerPort, hMsg);
  }

  return NW_OK;
}

/**
  Handle Triggered Request from Peer Entity.

  @param[in] thiz : Stack context
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleTriggeredReq(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t msgType,
    NW_IN uint8_t* msgBuf, NW_IN uint32_t msgBufLen, NW_IN uint16_t localPort,
    NW_IN uint16_t peerPort, NW_IN struct sockaddr* peerIp) {
  nw_rc_t rc                       = NW_FAILURE;
  nw_gtpv2c_trxn_t *pTrxn          = NULL, keyTrxn;
  nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};
  uint32_t teidLocal                       = 0;
  nw_gtpv2c_msg_handle_t hMsg              = 0;
  nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel = 0;
  nw_gtpv2c_error_t error                  = {0};

  // bool                                       noDelete = false ;
  keyTrxn.seqNum =
      ntohl(*((uint32_t*) (msgBuf + (((*msgBuf) & 0x08) ? 8 : 4)))) >> 8;
  ;
  memcpy(
      (void*) &keyTrxn.peer_ip, peerIp,
      (peerIp->sa_family == AF_INET) ? sizeof(struct sockaddr_in) :
                                       sizeof(struct sockaddr_in6));

  OAILOG_DEBUG(
      LOG_GTPV2C,
      "RECEIVED GTPV2c triggered request message of type %d, length %d and "
      "seqNum %x.\n",
      msgType, msgBufLen, keyTrxn.seqNum);

  /** A transaction of the initial request (cmd) for the triggered request
   * should exist. */
  pTrxn = RB_FIND(
      NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
      &keyTrxn);

  if (pTrxn) {
    /**
     * We remove the transaction of the initial request and create a new
     * transaction the the received triggered request.
     */
    RB_REMOVE(
        NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
        pTrxn);
    rc = nwGtpv2cTrxnDelete(&pTrxn);
    NW_ASSERT(NW_OK == rc);
  } else {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Triggered request message without a matching outstanding request "
        "(cmd) received! Discarding.\n");
    rc = NW_OK;
  }
  NW_ASSERT(msgBuf && msgBufLen);

  /** Process it like an initial request. */
  if (teidLocal) {
    keyTunnel.teid = ntohl(teidLocal);
    memcpy(
        (void*) &keyTunnel.ipAddrRemote, peerIp,
        (peerIp->sa_family == AF_INET) ? sizeof(struct sockaddr_in) :
                                         sizeof(struct sockaddr_in6));
    pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(thiz->tunnelMap), &keyTunnel);

    if (!pLocalTunnel) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Request message received on non-existent teid 0x%x from peer %p "
          "received! Discarding.\n",
          ntohl(teidLocal), peerIp);
      return NW_OK;
    }
    hUlpTunnel = pLocalTunnel->hUlpTunnel;
  } else {
    hUlpTunnel = 0;
  }

  /** Create a new transaction for the same transaction id. */
  pTrxn = nwGtpv2cTrxnOutstandingRxNew(
      thiz, ntohl(teidLocal), peerIp, peerPort, (keyTrxn.seqNum));

  if (pTrxn) {
    pTrxn->localPort = localPort;
    rc               = nwGtpv2cMsgFromBufferNew(
        (nw_gtpv2c_stack_handle_t) thiz, msgBuf, msgBufLen, &(hMsg));
    NW_ASSERT(thiz->pGtpv2cMsgIeParseInfo[msgType]);
    rc = nwGtpv2cMsgIeParse(thiz->pGtpv2cMsgIeParseInfo[msgType], hMsg, &error);

    if (rc != NW_OK) {
      char ipv4[INET_ADDRSTRLEN];
      inet_ntop(AF_INET, (void*) peerIp, ipv4, INET_ADDRSTRLEN);
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Malformed triggered request message received on TEID %u from peer "
          "%s. Notifying ULP.\n",
          ntohl(teidLocal), ipv4);
    }

    rc = nwGtpv2cSendTriggeredReqIndToUlp(
        thiz, &error, pTrxn, hUlpTunnel, msgType, peerIp, peerPort, hMsg);
  }
  return rc;
}

/**
  Handle Triggered Response from Peer Entity.

  @param[in] thiz : Stack context
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cHandleTriggeredRsp(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t msgType,
    NW_IN uint8_t* msgBuf, NW_IN uint32_t msgBufLen, NW_IN uint16_t localPort,
    NW_IN uint16_t peerPort, NW_IN struct sockaddr* peerIp, NW_IN bool remove) {
  nw_rc_t rc                  = NW_FAILURE;
  nw_gtpv2c_trxn_t *pTrxn     = NULL, keyTrxn;
  nw_gtpv2c_msg_handle_t hMsg = 0;
  nw_gtpv2c_error_t error     = {0};

  bool noDelete = false;
  keyTrxn.seqNum =
      ntohl(*((uint32_t*) (msgBuf + (((*msgBuf) & 0x08) ? 8 : 4)))) >> 8;
  ;
  memcpy(
      (void*) &keyTrxn.peer_ip, peerIp,
      (peerIp->sa_family == AF_INET) ? sizeof(struct sockaddr_in) :
                                       sizeof(struct sockaddr_in6));

  OAILOG_DEBUG(
      LOG_GTPV2C,
      "RECEIVED GTPV2c  response message of type %d, length %d and seqNum "
      "%x.\n",
      msgType, msgBufLen, keyTrxn.seqNum);

  pTrxn = RB_FIND(
      NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
      &keyTrxn);
  uint8_t trx_flags = 0;
  if (pTrxn) {
    uint32_t hUlpTunnel;

    noDelete = pTrxn->noDelete;
    hUlpTunnel =
        (pTrxn->hTunnel ? ((nw_gtpv2c_tunnel_t*) (pTrxn->hTunnel))->hUlpTunnel :
                          0);

    if (remove) {
      /**
       * Remove all except create session request.
       * The remove flag might be changed for CSR.
       */
      trx_flags = pTrxn->trx_flags;
      OAILOG_DEBUG(
          LOG_GTPV2C,
          "Not removing the initial request transaction for message type %d, "
          "seqNo %x (altough remove flag set). \n",
          msgType, keyTrxn.seqNum);

    } else {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Not removing the initial request transaction for message type %d, "
          "seqNo %x. \n",
          msgType, keyTrxn.seqNum);
      /** This transaction should be taken care automatically. No additional
       * information is necessary. */
      rc = nwGtpv2cMsgDelete(
          (nw_gtpv2c_stack_handle_t) thiz,
          (nw_gtpv2c_msg_handle_t) pTrxn->pMsg);
      NW_ASSERT(NW_OK == rc);
      pTrxn->pMsg = NULL;
      NW_ASSERT(NW_OK == rc);
      /** Mark the transaction as ACKed. */
      pTrxn->trx_flags |= INTERNAL_FLAG_TRIGGERED_ACK;
      trx_flags = pTrxn->trx_flags;
    }

    NW_ASSERT(msgBuf && msgBufLen);
    rc = nwGtpv2cMsgFromBufferNew(
        (nw_gtpv2c_stack_handle_t) thiz, msgBuf, msgBufLen, &(hMsg));
    NW_ASSERT(thiz->pGtpv2cMsgIeParseInfo[msgType]);
    rc = nwGtpv2cMsgIeParse(thiz->pGtpv2cMsgIeParseInfo[msgType], hMsg, &error);

    /** We just forward the remaining message. */
    if (rc != NW_OK) {
      char ip[INET6_ADDRSTRLEN];
      inet_ntop(
          peerIp->sa_family, (void*) peerIp, ip,
          peerIp->sa_family == AF_INET ? INET_ADDRSTRLEN : INET_ADDRSTRLEN);
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Malformed message received on TEID %u from peer %s. Notifying "
          "ULP.\n",
          ntohl((*((uint32_t*) (msgBuf + 4)))), ip);
    }
    rc = nwGtpv2cSendTriggeredRspIndToUlp(
        thiz, &error, keyTrxn.seqNum, &trx_flags, localPort, peerPort, peerIp,
        hUlpTunnel, msgType, noDelete, hMsg);
    if (remove && !(trx_flags & INTERNAL_LATE_RESPONS_IND)) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Removing the initial request transaction for message type %d, seqNo "
          "%x in conclusion (not late response). \n",
          msgType, keyTrxn.seqNum);
      /** Remove the transaction. */
      RB_REMOVE(
          NwGtpv2cOutstandingTxSeqNumTrxnMap, &(thiz->outstandingTxSeqNumMap),
          pTrxn);
      rc = nwGtpv2cTrxnDelete(&pTrxn);
      NW_ASSERT(NW_OK == rc);
      remove = false;
    } else {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Not removing the initial request transaction for message type %d, "
          "seqNo %x since it was a late response. \n",
          msgType, keyTrxn.seqNum);
      /** Remove the decoded message. */
      rc = NW_OK;
    }
  } else {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Response message without a matching outstanding request received! "
        "Discarding.\n");
    rc = NW_OK;
  }

  return rc;
}

/*--------------------------------------------------------------------------*
                       P U B L I C   F U N C T I O N S
  --------------------------------------------------------------------------*/

/**
   Constructor
*/

nw_rc_t nwGtpv2cInitialize(
    NW_INOUT nw_gtpv2c_stack_handle_t* hGtpcStackHandle) {
  nw_rc_t rc              = NW_OK;
  nw_gtpv2c_stack_t* thiz = NULL;

  thiz = (nw_gtpv2c_stack_t*) malloc(sizeof(nw_gtpv2c_stack_t));
  memset(thiz, 0, sizeof(nw_gtpv2c_stack_t));

  if (thiz) {
    OAI_GCC_DIAG_OFF("-Wpointer-to-int-cast");
    thiz->id     = (uint32_t) thiz;
    thiz->seqNum = ((uint32_t) thiz) & 0x0000FFFF;
    OAI_GCC_DIAG_ON("-Wpointer-to-int-cast");
    RB_INIT(&(thiz->tunnelMap));
    RB_INIT(&(thiz->outstandingTxSeqNumMap));
    RB_INIT(&(thiz->outstandingRxSeqNumMap));
    RB_INIT(&(thiz->activeTimerList));
    OAI_GCC_DIAG_OFF("-Wpointer-to-int-cast");
    thiz->hTmrMinHeap = (NwPtrT) nwGtpv2cTmrMinHeapNew(10000);
    OAI_GCC_DIAG_ON("-Wpointer-to-int-cast");
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_ECHO_RSP);

    // For S11 interface
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CREATE_SESSION_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CREATE_SESSION_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_SESSION_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_SESSION_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_MODIFY_BEARER_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_MODIFY_BEARER_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CREATE_BEARER_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CREATE_BEARER_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_UPDATE_BEARER_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_UPDATE_BEARER_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_BEARER_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_BEARER_RSP);
    /** Congestion related command. */
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_BEARER_CMD);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DELETE_BEARER_FAILURE_IND);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_BEARER_RESOURCE_CMD);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_BEARER_RESOURCE_FAILURE_IND);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_RELEASE_ACCESS_BEARERS_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_RELEASE_ACCESS_BEARERS_RSP);
    /** Paging related GTPv2c signaling. */
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_DOWNLINK_DATA_NOTIFICATION);

    // For S10 interface
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_FORWARD_RELOCATION_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_FORWARD_RELOCATION_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_FORWARD_ACCESS_CONTEXT_NTF);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_FORWARD_ACCESS_CONTEXT_ACK);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(
        thiz, NW_GTP_FORWARD_RELOCATION_COMPLETE_NTF);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(
        thiz, NW_GTP_FORWARD_RELOCATION_COMPLETE_ACK);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CONTEXT_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CONTEXT_RSP);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_CONTEXT_ACK);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_RELOCATION_CANCEL_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_RELOCATION_CANCEL_RSP);

    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_IDENTIFICATION_REQ);
    NW_GTPV2C_INIT_MSG_IE_PARSE_INFO(thiz, NW_GTP_IDENTIFICATION_RSP);
    nwGtpv2cDisplayBanner(thiz);
  } else {
    rc = NW_FAILURE;
  }
  *hGtpcStackHandle = (nw_gtpv2c_stack_handle_t) thiz;
  return rc;
}

/**
   Destructor
*/

nw_rc_t nwGtpv2cFinalize(NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle) {
  if (!hGtpcStackHandle) return NW_FAILURE;

  nwGtpv2cMsgIeParseInfoDelete(((nw_gtpv2c_stack_t*) hGtpcStackHandle)
                                   ->pGtpv2cMsgIeParseInfo[NW_GTP_ECHO_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(((nw_gtpv2c_stack_t*) hGtpcStackHandle)
                                   ->pGtpv2cMsgIeParseInfo[NW_GTP_ECHO_RSP]);

  // For S11 interface
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_CREATE_SESSION_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_CREATE_SESSION_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_SESSION_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_SESSION_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_MODIFY_BEARER_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_CREATE_BEARER_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_CREATE_BEARER_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_UPDATE_BEARER_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_UPDATE_BEARER_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_BEARER_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_BEARER_RSP]);

  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_BEARER_RESOURCE_CMD]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_BEARER_RESOURCE_FAILURE_IND]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_BEARER_CMD]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DELETE_BEARER_FAILURE_IND]);

  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_RELEASE_ACCESS_BEARERS_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_RELEASE_ACCESS_BEARERS_RSP]);

  /** Paging related GTPv2c signaling. */
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DOWNLINK_DATA_NOTIFICATION]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_DOWNLINK_DATA_NOTIFICATION_ACK]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo
              [NW_GTP_DOWNLINK_DATA_NOTIFICATION_FAILURE_IND]);

  // For S10 interface
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_RELOCATION_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_RELOCATION_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_ACCESS_CONTEXT_NTF]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_ACCESS_CONTEXT_ACK]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_RELOCATION_COMPLETE_NTF]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_FORWARD_RELOCATION_COMPLETE_ACK]);
  nwGtpv2cMsgIeParseInfoDelete(((nw_gtpv2c_stack_t*) hGtpcStackHandle)
                                   ->pGtpv2cMsgIeParseInfo[NW_GTP_CONTEXT_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(((nw_gtpv2c_stack_t*) hGtpcStackHandle)
                                   ->pGtpv2cMsgIeParseInfo[NW_GTP_CONTEXT_RSP]);
  nwGtpv2cMsgIeParseInfoDelete(((nw_gtpv2c_stack_t*) hGtpcStackHandle)
                                   ->pGtpv2cMsgIeParseInfo[NW_GTP_CONTEXT_ACK]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_RELOCATION_CANCEL_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_RELOCATION_CANCEL_RSP]);

  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_IDENTIFICATION_REQ]);
  nwGtpv2cMsgIeParseInfoDelete(
      ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->pGtpv2cMsgIeParseInfo[NW_GTP_IDENTIFICATION_RSP]);

  OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
  nwGtpv2cTmrMinHeapDelete(
      (NwGtpv2cTmrMinHeapT*) ((nw_gtpv2c_stack_t*) hGtpcStackHandle)
          ->hTmrMinHeap);
  OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");
  free_wrapper((void**) &hGtpcStackHandle);
  return NW_OK;
}

/**
   Set ULP entity
*/

nw_rc_t nwGtpv2cSetUlpEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_ulp_entity_t* pUlpEntity) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  if (!pUlpEntity) return NW_FAILURE;

  thiz->ulp = *(pUlpEntity);
  return NW_OK;
}

/**
   Set UDP entity
*/
nw_rc_t nwGtpv2cSetUdpEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_udp_entity_t* pUdpEntity) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  if (!pUdpEntity) return NW_FAILURE;

  thiz->udp = *(pUdpEntity);
  return NW_OK;
}

/**
   Set MEM MGR entity
*/
nw_rc_t nwGtpv2cSetMemMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_mem_mgr_entity_t* pMemMgrEntity) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  if (!pMemMgrEntity) return NW_FAILURE;

  thiz->memMgr = *(pMemMgrEntity);
  return NW_OK;
}

/**
   Set TMR MGR entity
*/

nw_rc_t nwGtpv2cSetTimerMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_timer_mgr_entity_t* pTmrMgrEntity) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  if (!pTmrMgrEntity) return NW_FAILURE;

  thiz->tmrMgr = *(pTmrMgrEntity);
  return NW_OK;
}

/**
   Set LOG MGR entity
*/

nw_rc_t nwGtpv2cSetLogMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_log_mgr_entity_t* pLogMgrEntity) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  if (!pLogMgrEntity) return NW_FAILURE;

  thiz->logMgr = *(pLogMgrEntity);
  return NW_OK;
}

/**
  Set log level for the stack.
*/

nw_rc_t nwGtpv2cSetLogLevel(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint32_t logLevel) {
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  thiz->logLevel = logLevel;
  return NW_OK;
}

/**
   Process Request from Udp Layer
*/

nw_rc_t nwGtpv2cProcessUdpReq(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint8_t* udpData,
    NW_IN uint32_t udpDataLen, NW_IN uint16_t localPort,
    NW_IN uint16_t peerPort, NW_IN struct sockaddr* peerIp) {
  nw_rc_t rc              = NW_FAILURE;
  nw_gtpv2c_stack_t* thiz = NULL;
  uint16_t msgType        = 0;

  thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;
  NW_ASSERT(thiz);
  OAILOG_FUNC_IN(LOG_GTPV2C);

  if (udpDataLen < NW_GTPV2C_MINIMUM_HEADER_SIZE) {
    /*
     * TS 29.274 Section 7.7.3:
     * If a GTP entity receives a message, which is too short to
     * contain the respective GTPv2 header, the GTP-PDU shall be
     * silently discarded
     */
    OAILOG_WARNING(LOG_GTPV2C, "Received message too small! Discarding.\n");
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  if ((ntohs(*((uint16_t*) ((uint8_t*) udpData + 2))) /* Length */
       + ((*((uint8_t*) (udpData)) & 0x08) ?
              4 :
              0) /* Extra Header length if TEID present */) > udpDataLen) {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Received message with erroneous length of %u against expected length "
        "of %u! Discarding\n",
        udpDataLen,
        ntohs(*((uint16_t*) ((uint8_t*) udpData + 2))) +
            ((*((uint8_t*) (udpData)) & 0x08) ? 4 : 0));
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  if (((*((uint8_t*) (udpData)) & 0xE0) >> 5) != NW_GTP_VERSION) {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Received unsupported GTP version '%u' message! Discarding.\n",
        ((*((uint8_t*) (udpData)) & 0xE0) >> 5));

    // Send Version Not Supported Message to peer
    rc = nwGtpv2cSendVersionNotSupportedInd(
        thiz, peerIp, peerPort,
        *((uint32_t*) (udpData + ((*((uint8_t*) (udpData)) & 0x08) ? 8 : 4))) /* Seq Num */);
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  msgType = *((uint8_t*) (udpData + 1));

  switch (msgType) {
    case NW_GTP_ECHO_REQ: {
      rc = nwGtpv2cHandleEchoReq(
          thiz, msgType, udpData, udpDataLen, peerPort, peerIp);
    } break;

    /** Definitive Initial Requests. */
    case NW_GTP_CREATE_SESSION_REQ:
    case NW_GTP_MODIFY_BEARER_REQ:
    case NW_GTP_DELETE_SESSION_REQ:
    case NW_GTP_RELEASE_ACCESS_BEARERS_REQ:
    case NW_GTP_CREATE_INDIRECT_DATA_FORWARDING_TUNNEL_REQ:
    case NW_GTP_DELETE_INDIRECT_DATA_FORWARDING_TUNNEL_REQ:
    case NW_GTP_REMOTE_UE_REPORT_NOTIFICATION:
      /** Handover Related Messages. */
    case NW_GTP_FORWARD_RELOCATION_REQ:
    case NW_GTP_FORWARD_RELOCATION_COMPLETE_NTF:
    case NW_GTP_RELOCATION_CANCEL_REQ:
    case NW_GTP_CONTEXT_REQ:
      /** S11: Paging. */
    case NW_GTP_DOWNLINK_DATA_NOTIFICATION:
      rc = nwGtpv2cHandleInitialReq(
          thiz, msgType, udpData, udpDataLen, peerPort, peerIp);
      break;

    case NW_GTP_FORWARD_ACCESS_CONTEXT_NTF:
      rc = nwGtpv2cHandleInitialReq(
          thiz, msgType, udpData, udpDataLen, peerPort, peerIp);
      break;

    /** May be initial request or triggered requests. */
    case NW_GTP_CREATE_BEARER_REQ:
    case NW_GTP_UPDATE_BEARER_REQ:
    case NW_GTP_DELETE_BEARER_REQ: {
      /** Check the received port, if it is an Initial Request, a Triggered
       * Request or a Triggered Response. */
      if (localPort == thiz->udp.gtpv2cStandardPort) {
        /** Message received on standard port, checking for Initial Requests. */
        rc = nwGtpv2cHandleInitialReq(
            thiz, msgType, udpData, udpDataLen, peerPort, peerIp);
        break;
      } else {
        /** Message received on high port, checking for triggered requests and
         * responses. */
        rc = nwGtpv2cHandleTriggeredReq(
            thiz, msgType, udpData, udpDataLen, localPort, peerPort, peerIp);
        break;
      }
    } break;

    case NW_GTP_ECHO_RSP:
    case NW_GTP_CREATE_SESSION_RSP:
    case NW_GTP_MODIFY_BEARER_RSP:
    case NW_GTP_DELETE_SESSION_RSP:
    case NW_GTP_CREATE_BEARER_RSP:
    case NW_GTP_UPDATE_BEARER_RSP:
    case NW_GTP_DELETE_BEARER_RSP:
    case NW_GTP_DELETE_BEARER_FAILURE_IND:
    case NW_GTP_GTP_REMOTE_UE_REPORT_ACK:
    case NW_GTP_RELEASE_ACCESS_BEARERS_RSP:
    case NW_GTP_CREATE_INDIRECT_DATA_FORWARDING_TUNNEL_RSP:
    /** Handover Related Messages. */
    case NW_GTP_FORWARD_RELOCATION_RSP:
    case NW_GTP_FORWARD_ACCESS_CONTEXT_ACK:
    case NW_GTP_DELETE_INDIRECT_DATA_FORWARDING_TUNNEL_RSP:
    case NW_GTP_FORWARD_RELOCATION_COMPLETE_ACK:
    case NW_GTP_RELOCATION_CANCEL_RSP:
      rc = nwGtpv2cHandleTriggeredRsp(
          thiz, msgType, udpData, udpDataLen, localPort, peerPort, peerIp,
          true); /**< We will check inside, if the received response is to be
                    acked. */
      break;
    case NW_GTP_CONTEXT_RSP:
      rc = nwGtpv2cHandleTriggeredRsp(
          thiz, msgType, udpData, udpDataLen, localPort, peerPort, peerIp,
          false); /**< We will check inside, if the received response is to be
                     acked. */
      break;
    case NW_GTP_CONTEXT_ACK:
      /** Ignore the received Ctx ACK (no transaction). */
      // todo: handle eventually..
      rc = NW_OK;
      break;

    default: {
      /*
       * TS 29.274 Section 7.7.4:
       * If a GTP entity receives a message with an unknown Message Type
       * value, it shall silently discard the message.
       */
      OAILOG_WARNING(
          LOG_GTPV2C, "Received unknown message type %u from UDP! Ignoring.\n",
          msgType);
      rc = NW_OK;
    }
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

//   Process Request from Upper Layer

nw_rc_t nwGtpv2cProcessUlpReq(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_ulp_api_t* pUlpReq) {
  nw_rc_t rc              = NW_FAILURE;
  nw_gtpv2c_stack_t* thiz = (nw_gtpv2c_stack_t*) hGtpcStackHandle;

  NW_ASSERT(thiz);
  NW_ASSERT(pUlpReq != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);

  switch (pUlpReq->apiType & 0x00FFFFFFL) {
    case NW_GTPV2C_ULP_API_INITIAL_REQ: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received initial request from ulp\n");
      rc = nwGtpv2cHandleUlpInitialReq(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_API_TRIGGERED_REQ: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received triggered request from ulp\n");
      rc = nwGtpv2cHandleUlpTriggeredReq(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_API_TRIGGERED_RSP: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received triggered response from ulp\n");
      rc = nwGtpv2cHandleUlpTriggeredRsp(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_API_TRIGGERED_ACK: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received triggered acknowledgement from ulp\n");
      rc = nwGtpv2cHandleUlpTriggeredAck(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_CREATE_LOCAL_TUNNEL: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received create local tunnel from ulp\n");
      rc = nwGtpv2cHandleUlpCreateLocalTunnel(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_DELETE_LOCAL_TUNNEL: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received delete local tunnel from ulp\n");
      rc = nwGtpv2cHandleUlpDeleteLocalTunnel(thiz, pUlpReq);
    } break;

    case NW_GTPV2C_ULP_FIND_LOCAL_TUNNEL: {
      OAILOG_DEBUG(LOG_GTPV2C, "Received find local tunnel from ulp\n");
      rc = nwGtpv2cHandleUlpFindLocalTunnel(thiz, pUlpReq);
    } break;

    default: {
      OAILOG_WARNING(
          LOG_GTPV2C, "Received unhandled API 0x%x from ULP! Ignoring.\n",
          pUlpReq->apiType);
      rc = NW_FAILURE;
    } break;
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
   Process Timer timeout Request from Timer ULP Manager
*/

nw_rc_t nwGtpv2cProcessTimeoutOld(void* arg) {
  nw_rc_t rc                                 = NW_FAILURE;
  nw_gtpv2c_stack_t* thiz                    = NULL;
  nw_gtpv2c_timeout_info_t* timeoutInfo      = (nw_gtpv2c_timeout_info_t*) arg;
  nw_gtpv2c_timeout_info_t* pNextTimeoutInfo = NULL;
  struct timeval tv                          = {0};

  NW_ASSERT(timeoutInfo != NULL);
  thiz =
      (nw_gtpv2c_stack_t*) (((nw_gtpv2c_timeout_info_t*) timeoutInfo)->hStack);
  NW_ASSERT(thiz != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);

  if (thiz->activeTimerInfo == timeoutInfo) {
    thiz->activeTimerInfo = NULL;
    RB_REMOVE(NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);
    timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = timeoutInfo;
    rc = ((timeoutInfo)->timeoutCallbackFunc)(timeoutInfo->timeoutArg);
  } else {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Received timeout event from ULP for non-existent timeoutInfo 0x%p and "
        "activeTimer 0x%p!\n",
        timeoutInfo, thiz->activeTimerInfo);
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  NW_ASSERT(gettimeofday(&tv, NULL) == 0);

  for ((timeoutInfo) =
           RB_MIN(NwGtpv2cActiveTimerList, &(thiz->activeTimerList));
       (timeoutInfo) != NULL;) {
    if (NW_GTPV2C_TIMER_CMP_P(&timeoutInfo->tvTimeout, &tv, >)) break;

    pNextTimeoutInfo =
        RB_NEXT(NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);
    RB_REMOVE(NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);
    timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = timeoutInfo;
    rc          = ((timeoutInfo)->timeoutCallbackFunc)(timeoutInfo->timeoutArg);
    timeoutInfo = pNextTimeoutInfo;
  }

  // activeTimerInfo may be reset by the timeoutCallbackFunc call above
  if (thiz->activeTimerInfo == NULL) {
    timeoutInfo = RB_MIN(NwGtpv2cActiveTimerList, &(thiz->activeTimerList));

    if (timeoutInfo) {
      NW_GTPV2C_TIMER_SUB(&timeoutInfo->tvTimeout, &tv, &tv);
      rc = thiz->tmrMgr.tmrStartCallback(
          thiz->tmrMgr.tmrMgrHandle, tv.tv_sec, tv.tv_usec,
          timeoutInfo->tmrType, (void*) timeoutInfo, &timeoutInfo->hTimer);
      NW_ASSERT(NW_OK == rc);
      thiz->activeTimerInfo = timeoutInfo;
    }
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

nw_rc_t nwGtpv2cProcessTimeout(void* arg) {
  nw_rc_t rc                            = NW_FAILURE;
  nw_gtpv2c_stack_t* thiz               = NULL;
  nw_gtpv2c_timeout_info_t* timeoutInfo = (nw_gtpv2c_timeout_info_t*) arg;
  struct timeval tv                     = {0};

  NW_ASSERT(timeoutInfo != NULL);
  thiz = (nw_gtpv2c_stack_t*) (timeoutInfo->hStack);
  NW_ASSERT(thiz != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);

  if (thiz->activeTimerInfo == timeoutInfo) {
    thiz->activeTimerInfo = NULL;
    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    rc = nwGtpv2cTmrMinHeapRemove(
        (NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap,
        timeoutInfo->timerMinHeapIndex);
    OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");
    timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = timeoutInfo;
    rc = ((timeoutInfo)->timeoutCallbackFunc)(timeoutInfo->timeoutArg);
  } else {
    OAILOG_WARNING(
        LOG_GTPV2C,
        "Received timeout event from ULP for "
        "non-existent timeoutInfo 0x%p and activeTimer 0x%p!\n",
        timeoutInfo, thiz->activeTimerInfo);
    OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
  }

  NW_ASSERT(gettimeofday(&tv, NULL) == 0);
  OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
  timeoutInfo =
      nwGtpv2cTmrMinHeapPeek((NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap);
  OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");

  while ((timeoutInfo) != NULL) {
    if (NW_GTPV2C_TIMER_CMP_P(&timeoutInfo->tvTimeout, &tv, >)) break;

    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    rc = nwGtpv2cTmrMinHeapRemove(
        (NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap,
        timeoutInfo->timerMinHeapIndex);
    OAI_GCC_DIAG_ON(int - to - pointer - cast);
    timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = timeoutInfo;
    rc = ((timeoutInfo)->timeoutCallbackFunc)(timeoutInfo->timeoutArg);
    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    timeoutInfo =
        nwGtpv2cTmrMinHeapPeek((NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap);
    OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");
  }

  // activeTimerInfo may be reset by the timeoutCallbackFunc call above
  if (thiz->activeTimerInfo == NULL) {
    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    timeoutInfo =
        nwGtpv2cTmrMinHeapPeek((NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap);
    OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");

    if (timeoutInfo) {
      NW_GTPV2C_TIMER_SUB(&timeoutInfo->tvTimeout, &tv, &tv);
      rc = thiz->tmrMgr.tmrStartCallback(
          thiz->tmrMgr.tmrMgrHandle, tv.tv_sec, tv.tv_usec,
          timeoutInfo->tmrType, (void*) timeoutInfo, &timeoutInfo->hTimer);
      NW_ASSERT(NW_OK == rc);
      thiz->activeTimerInfo = timeoutInfo;
    }
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
   Start Timer with ULP Timer Manager
*/

nw_rc_t nwGtpv2cStartTimer(
    nw_gtpv2c_stack_t* thiz, uint32_t timeoutSec, uint32_t timeoutUsec,
    uint32_t tmrType, nw_rc_t (*timeoutCallbackFunc)(void*),
    void* timeoutCallbackArg, nw_gtpv2c_timer_handle_t* phTimer) {
  nw_rc_t rc                            = NW_OK;
  struct timeval tv                     = {0};
  nw_gtpv2c_timeout_info_t* timeoutInfo = NULL;

  OAILOG_FUNC_IN(LOG_GTPV2C);

  if (gpGtpv2cTimeoutInfoPool) {
    timeoutInfo             = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = gpGtpv2cTimeoutInfoPool->next;
  } else {
    NW_GTPV2C_MALLOC(
        thiz, sizeof(nw_gtpv2c_timeout_info_t), timeoutInfo,
        nw_gtpv2c_timeout_info_t*);
  }

  if (timeoutInfo) {
    timeoutInfo->tmrType             = tmrType;
    timeoutInfo->timeoutArg          = timeoutCallbackArg;
    timeoutInfo->timeoutCallbackFunc = timeoutCallbackFunc;
    timeoutInfo->hStack              = (nw_gtpv2c_stack_handle_t) thiz;
    NW_ASSERT(gettimeofday(&tv, NULL) == 0);
    NW_ASSERT(gettimeofday(&timeoutInfo->tvTimeout, NULL) == 0);
    timeoutInfo->tvTimeout.tv_sec  = timeoutSec;
    timeoutInfo->tvTimeout.tv_usec = timeoutUsec;
    NW_GTPV2C_TIMER_ADD(&tv, &timeoutInfo->tvTimeout, &timeoutInfo->tvTimeout);
    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    rc = nwGtpv2cTmrMinHeapInsert(
        (NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap, timeoutInfo);
    OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");
#if 0

      do {
        collision = RB_INSERT (NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);

        if (!collision)
          break;

        OAILOG_WARNING (LOG_GTPV2C,  "timer collision!\n");
        timeoutInfo->tvTimeout.tv_usec++;       /* HACK: In case there is a collision, schedule this event 1 usec later */

        if (timeoutInfo->tvTimeout.tv_usec > (999999 /*1000000 - 1 */ )) {
          timeoutInfo->tvTimeout.tv_usec = 0;
          timeoutInfo->tvTimeout.tv_sec++;
        }
      } while (1);

#endif

    if (thiz->activeTimerInfo) {
      if (NW_GTPV2C_TIMER_CMP_P(
              &(thiz->activeTimerInfo->tvTimeout), &(timeoutInfo->tvTimeout),
              >)) {
        //          OAILOG_DEBUG (LOG_GTPV2C, "Stopping active timer 0x%"
        //          PRIxPTR " for info 0x%p!\n", thiz->activeTimerInfo->hTimer,
        //          thiz->activeTimerInfo);
        rc = thiz->tmrMgr.tmrStopCallback(
            thiz->tmrMgr.tmrMgrHandle, thiz->activeTimerInfo->hTimer);
        NW_ASSERT(NW_OK == rc);
      } else {
        OAILOG_DEBUG(
            LOG_GTPV2C, "Already Started timer 0x%" PRIxPTR " for info 0x%p!\n",
            thiz->activeTimerInfo->hTimer, thiz->activeTimerInfo);
        *phTimer = (nw_gtpv2c_timer_handle_t) timeoutInfo;
        OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
      }
    }

    rc = thiz->tmrMgr.tmrStartCallback(
        thiz->tmrMgr.tmrMgrHandle, timeoutSec, timeoutUsec, tmrType,
        (void*) timeoutInfo, &timeoutInfo->hTimer);
    OAILOG_DEBUG(
        LOG_GTPV2C, "Started timer 0x%" PRIxPTR " for info 0x%p!\n",
        timeoutInfo->hTimer, timeoutInfo);
    NW_ASSERT(NW_OK == rc);
    thiz->activeTimerInfo = timeoutInfo;
  }

  *phTimer = (nw_gtpv2c_timer_handle_t) timeoutInfo;
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

nw_rc_t nwGtpv2cStartTimerOld(
    nw_gtpv2c_stack_t* thiz, uint32_t timeoutSec, uint32_t timeoutUsec,
    uint32_t tmrType, nw_rc_t (*timeoutCallbackFunc)(void*),
    void* timeoutCallbackArg, nw_gtpv2c_timer_handle_t* phTimer) {
  nw_rc_t rc = NW_OK;
  struct timeval tv;
  nw_gtpv2c_timeout_info_t* timeoutInfo;
  nw_gtpv2c_timeout_info_t* collision;

  NW_ASSERT(thiz != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);

  if (gpGtpv2cTimeoutInfoPool) {
    timeoutInfo             = gpGtpv2cTimeoutInfoPool;
    gpGtpv2cTimeoutInfoPool = gpGtpv2cTimeoutInfoPool->next;
  } else {
    NW_GTPV2C_MALLOC(
        thiz, sizeof(nw_gtpv2c_timeout_info_t), timeoutInfo,
        nw_gtpv2c_timeout_info_t*);
  }

  if (timeoutInfo) {
    timeoutInfo->tmrType             = tmrType;
    timeoutInfo->timeoutArg          = timeoutCallbackArg;
    timeoutInfo->timeoutCallbackFunc = timeoutCallbackFunc;
    timeoutInfo->hStack              = (nw_gtpv2c_stack_handle_t) thiz;
    NW_ASSERT(gettimeofday(&tv, NULL) == 0);
    NW_ASSERT(gettimeofday(&timeoutInfo->tvTimeout, NULL) == 0);
    timeoutInfo->tvTimeout.tv_sec  = timeoutSec;
    timeoutInfo->tvTimeout.tv_usec = timeoutUsec;
    NW_GTPV2C_TIMER_ADD(&tv, &timeoutInfo->tvTimeout, &timeoutInfo->tvTimeout);

    do {
      collision = RB_INSERT(
          NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);

      if (!collision) break;

      OAILOG_WARNING(LOG_GTPV2C, "timer collision!\n");
      timeoutInfo->tvTimeout.tv_usec++; /* HACK: In case there is a collision,
                                           schedule this event 1 usec later */

      if (timeoutInfo->tvTimeout.tv_usec > (999999 /*1000000 - 1 */)) {
        timeoutInfo->tvTimeout.tv_usec = 0;
        timeoutInfo->tvTimeout.tv_sec++;
      }
    } while (1);

    if (thiz->activeTimerInfo) {
      if (NW_GTPV2C_TIMER_CMP_P(
              &(thiz->activeTimerInfo->tvTimeout), &(timeoutInfo->tvTimeout),
              >)) {
        OAILOG_DEBUG(
            LOG_GTPV2C, "Stopping active timer 0x%" PRIxPTR " for info 0x%p!\n",
            thiz->activeTimerInfo->hTimer, thiz->activeTimerInfo);
        rc = thiz->tmrMgr.tmrStopCallback(
            thiz->tmrMgr.tmrMgrHandle, thiz->activeTimerInfo->hTimer);
        NW_ASSERT(NW_OK == rc);
      } else {
        OAILOG_DEBUG(
            LOG_GTPV2C, "Already Started timer 0x%" PRIxPTR " for info 0x%p!\n",
            thiz->activeTimerInfo->hTimer, thiz->activeTimerInfo);
        *phTimer = (nw_gtpv2c_timer_handle_t) timeoutInfo;
        OAILOG_FUNC_RETURN(LOG_GTPV2C, NW_OK);
      }
    }

    rc = thiz->tmrMgr.tmrStartCallback(
        thiz->tmrMgr.tmrMgrHandle, timeoutSec, timeoutUsec, tmrType,
        (void*) timeoutInfo, &timeoutInfo->hTimer);
    OAILOG_DEBUG(
        LOG_GTPV2C, "Started timer 0x%" PRIxPTR " for info 0x%p!\n",
        timeoutInfo->hTimer, timeoutInfo);
    NW_ASSERT(NW_OK == rc);
    thiz->activeTimerInfo = timeoutInfo;
  }

  *phTimer = (nw_gtpv2c_timer_handle_t) timeoutInfo;
  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

/**
   Stop Timer with ULP Timer Manager
*/
nw_rc_t nwGtpv2cStopTimer(
    nw_gtpv2c_stack_t* thiz, nw_gtpv2c_timer_handle_t hTimer) {
  nw_rc_t rc = NW_OK;
  struct timeval tv;
  nw_gtpv2c_timeout_info_t* timeoutInfo;

  NW_ASSERT(thiz != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);
  timeoutInfo = (nw_gtpv2c_timeout_info_t*) hTimer;
  OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
  rc = nwGtpv2cTmrMinHeapRemove(
      (NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap, timeoutInfo->timerMinHeapIndex);
  OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");
  timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
  gpGtpv2cTimeoutInfoPool = timeoutInfo;
  //    OAILOG_DEBUG (LOG_GTPV2C, "Stopping active timer 0x%" PRIxPTR " for info
  //    0x%p!\n", timeoutInfo->hTimer, timeoutInfo);

  if (thiz->activeTimerInfo == timeoutInfo) {
    OAILOG_DEBUG(
        LOG_GTPV2C, "Stopping active timer 0x%" PRIxPTR " for info 0x%p!\n",
        timeoutInfo->hTimer, timeoutInfo);
    rc = thiz->tmrMgr.tmrStopCallback(
        thiz->tmrMgr.tmrMgrHandle, timeoutInfo->hTimer);
    thiz->activeTimerInfo = NULL;
    if (NW_OK != rc) {
      OAILOG_ERROR(
          LOG_GTPV2C,
          "Stopping active timer 0x%" PRIxPTR " for info 0x%p failed!\n",
          timeoutInfo->hTimer, timeoutInfo);
    } else
      OAILOG_INFO(
          LOG_GTPV2C, "Stopped active timer 0x%" PRIxPTR " for info 0x%p!\n",
          timeoutInfo->hTimer, timeoutInfo);
    OAI_GCC_DIAG_OFF("-Wint-to-pointer-cast");
    timeoutInfo =
        nwGtpv2cTmrMinHeapPeek((NwGtpv2cTmrMinHeapT*) thiz->hTmrMinHeap);
    OAI_GCC_DIAG_ON("-Wint-to-pointer-cast");

    if (timeoutInfo) {
      NW_ASSERT(gettimeofday(&tv, NULL) == 0);

      if (NW_GTPV2C_TIMER_CMP_P(&timeoutInfo->tvTimeout, &tv, <)) {
        thiz->activeTimerInfo = timeoutInfo;
        rc                    = nwGtpv2cProcessTimeout(timeoutInfo);
        NW_ASSERT(NW_OK == rc);
      } else {
        NW_GTPV2C_TIMER_SUB(&timeoutInfo->tvTimeout, &tv, &tv);
        rc = thiz->tmrMgr.tmrStartCallback(
            thiz->tmrMgr.tmrMgrHandle, tv.tv_sec, tv.tv_usec,
            timeoutInfo->tmrType, (void*) timeoutInfo, &timeoutInfo->hTimer);
        NW_ASSERT(NW_OK == rc);
        OAILOG_DEBUG(
            LOG_GTPV2C, "Started timer 0x%" PRIxPTR " for info 0x%p!\n",
            timeoutInfo->hTimer, timeoutInfo);
        thiz->activeTimerInfo = timeoutInfo;
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

nw_rc_t nwGtpv2cStopTimerOld(
    nw_gtpv2c_stack_t* thiz, nw_gtpv2c_timer_handle_t hTimer) {
  nw_rc_t rc = NW_OK;
  struct timeval tv;
  nw_gtpv2c_timeout_info_t* timeoutInfo;

  NW_ASSERT(thiz != NULL);
  OAILOG_FUNC_IN(LOG_GTPV2C);
  timeoutInfo = (nw_gtpv2c_timeout_info_t*) hTimer;
  RB_REMOVE(NwGtpv2cActiveTimerList, &(thiz->activeTimerList), timeoutInfo);
  timeoutInfo->next       = gpGtpv2cTimeoutInfoPool;
  gpGtpv2cTimeoutInfoPool = timeoutInfo;
  OAILOG_DEBUG(
      LOG_GTPV2C, "Stopping active timer 0x%" PRIxPTR " for info 0x%p!\n",
      timeoutInfo->hTimer, timeoutInfo);

  if (thiz->activeTimerInfo == timeoutInfo) {
    OAILOG_DEBUG(
        LOG_GTPV2C, "Stopping active timer 0x%" PRIxPTR " for info 0x%p!\n",
        timeoutInfo->hTimer, timeoutInfo);
    rc = thiz->tmrMgr.tmrStopCallback(
        thiz->tmrMgr.tmrMgrHandle, timeoutInfo->hTimer);
    thiz->activeTimerInfo = NULL;
    NW_ASSERT(NW_OK == rc);
    timeoutInfo = RB_MIN(NwGtpv2cActiveTimerList, &(thiz->activeTimerList));

    if (timeoutInfo) {
      NW_ASSERT(gettimeofday(&tv, NULL) == 0);

      if (NW_GTPV2C_TIMER_CMP_P(&timeoutInfo->tvTimeout, &tv, <)) {
        thiz->activeTimerInfo = timeoutInfo;
        rc                    = nwGtpv2cProcessTimeout(timeoutInfo);
        NW_ASSERT(NW_OK == rc);
      } else {
        NW_GTPV2C_TIMER_SUB(&timeoutInfo->tvTimeout, &tv, &tv);
        rc = thiz->tmrMgr.tmrStartCallback(
            thiz->tmrMgr.tmrMgrHandle, tv.tv_sec, tv.tv_usec,
            timeoutInfo->tmrType, (void*) timeoutInfo, &timeoutInfo->hTimer);
        NW_ASSERT(NW_OK == rc);
        OAILOG_DEBUG(
            LOG_GTPV2C, "Started timer 0x%" PRIxPTR " for info 0x%p!\n",
            timeoutInfo->hTimer, timeoutInfo);
        thiz->activeTimerInfo = timeoutInfo;
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_GTPV2C, rc);
}

#ifdef __cplusplus
}
#endif

/*--------------------------------------------------------------------------*
                        E N D     O F    F I L E
  --------------------------------------------------------------------------*/
