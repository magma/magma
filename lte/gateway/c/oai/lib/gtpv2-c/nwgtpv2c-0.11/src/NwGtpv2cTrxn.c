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
#include <stdbool.h>

#include "bstrlib.h"

#include "NwTypes.h"
#include "NwLog.h"
#include "NwUtils.h"
#include "NwGtpv2cLog.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cPrivate.h"
#include "NwGtpv2cTrxn.h"
#include "log.h"

/*--------------------------------------------------------------------------*
                   P R I V A T E  D E C L A R A T I O N S
  --------------------------------------------------------------------------*/

#ifdef __cplusplus
extern "C" {
#endif

static nw_gtpv2c_trxn_t* gpGtpv2cTrxnPool = NULL;

/*--------------------------------------------------------------------------*
                     P R I V A T E      F U N C T I O N S
  --------------------------------------------------------------------------*/

/*---------------------------------------------------------------------------
   Send msg retransmission to peer via data request to UDP Entity
  --------------------------------------------------------------------------*/

static nw_rc_t nwGtpv2cTrxnSendMsgRetransmission(nw_gtpv2c_trxn_t* thiz) {
  nw_rc_t rc;
  NW_ASSERT(thiz);
  NW_ASSERT(thiz->pMsg);
  rc = thiz->pStack->udp.udpDataReqCallback(
      thiz->pStack->udp.hUdp, thiz->pMsg->msgBuf, thiz->pMsg->msgLen,
      thiz->localPort, (struct sockaddr*) &thiz->peer_ip, thiz->peerPort);
  thiz->maxRetries--;
  return rc;
}

static nw_rc_t nwGtpv2cTrxnPeerRspWaitTimeout(void* arg) {
  nw_rc_t rc = NW_OK;
  nw_gtpv2c_trxn_t* thiz;
  nw_gtpv2c_stack_t* pStack;

  thiz   = ((nw_gtpv2c_trxn_t*) arg);
  pStack = thiz->pStack;
  NW_ASSERT(pStack);
  OAILOG_WARNING(
      LOG_GTPV2C, "T3 Response timer expired for transaction %p\n", thiz);
  thiz->hRspTmr = 0;

  if (thiz->trx_flags & INTERNAL_FLAG_TRIGGERED_ACK) {
    OAILOG_ERROR(
        LOG_GTPV2C,
        "Transaction transaction %p (seqNo=0x%x) was acknowledged. Removing "
        "for timeout. \n",
        thiz, thiz->seqNum);
    RB_REMOVE(
        NwGtpv2cOutstandingTxSeqNumTrxnMap, &(pStack->outstandingTxSeqNumMap),
        thiz);
    rc = nwGtpv2cTrxnDelete(&thiz);
    return rc;
  }

  if (thiz->maxRetries) {
    /** Check if a tunnel endpoint exists. */
    nw_gtpv2c_tunnel_t *pLocalTunnel = NULL, keyTunnel = {0};
    keyTunnel.teid = thiz->teidLocal;
    memcpy(
        (void*) &keyTunnel.ipAddrRemote, (void*) &thiz->peer_ip,
        sizeof(thiz->peer_ip));
    pLocalTunnel = RB_FIND(NwGtpv2cTunnelMap, &(pStack->tunnelMap), &keyTunnel);
    if (pLocalTunnel) {
      rc = nwGtpv2cTrxnSendMsgRetransmission(thiz);
      NW_ASSERT(NW_OK == rc);
      rc = nwGtpv2cStartTimer(
          thiz->pStack, thiz->t3Timer, 0, NW_GTPV2C_TMR_TYPE_ONE_SHOT,
          nwGtpv2cTrxnPeerRspWaitTimeout, thiz, &thiz->hRspTmr);
    } else {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Tunnel for local-TEID 0x%x is removed for request transaction %p "
          "(seqNo=0x%x)! Removing the trx and ignoring timeout. \n",
          thiz->teidLocal, thiz, thiz->seqNum);
      RB_REMOVE(
          NwGtpv2cOutstandingTxSeqNumTrxnMap, &(pStack->outstandingTxSeqNumMap),
          thiz);
      rc = nwGtpv2cTrxnDelete(&thiz);
    }
  } else {
    nw_gtpv2c_ulp_api_t ulpApi;
    memset(&ulpApi, 0, sizeof(nw_gtpv2c_ulp_api_t));

    ulpApi.hMsg    = 0;
    ulpApi.apiType = NW_GTPV2C_ULP_API_RSP_FAILURE_IND;
    ulpApi.u_api_info.rspFailureInfo.hUlpTrxn = thiz->hUlpTrxn;
    ulpApi.u_api_info.rspFailureInfo.noDelete = thiz->noDelete;
    ulpApi.u_api_info.rspFailureInfo.msgType =
        thiz->pMsg ? thiz->pMsg->msgType : 0;
    ulpApi.u_api_info.rspFailureInfo.hUlpTunnel =
        ((thiz->hTunnel) ? ((nw_gtpv2c_tunnel_t*) (thiz->hTunnel))->hUlpTunnel :
                           0);
    ulpApi.u_api_info.rspFailureInfo.teidLocal =
        (thiz->hTunnel) ? ((nw_gtpv2c_tunnel_t*) (thiz->hTunnel))->teid : 0;
    /** Set the flags. */
    ulpApi.u_api_info.rspFailureInfo.trx_flags = thiz->trx_flags;
    OAILOG_ERROR(LOG_GTPV2C, "N3 retries expired for transaction %p\n", thiz);
    RB_REMOVE(
        NwGtpv2cOutstandingTxSeqNumTrxnMap, &(pStack->outstandingTxSeqNumMap),
        thiz);
    rc = nwGtpv2cTrxnDelete(&thiz);
    rc = pStack->ulp.ulpReqCallback(pStack->ulp.hUlp, &ulpApi);
  }

  return rc;
}

static nw_rc_t nwGtpv2cTrxnDuplicateRequestWaitTimeout(void* arg) {
  nw_rc_t rc = NW_OK;
  nw_gtpv2c_trxn_t* thiz;
  nw_gtpv2c_stack_t* pStack;

  thiz = ((nw_gtpv2c_trxn_t*) arg);
  NW_ASSERT(thiz);
  pStack = thiz->pStack;
  NW_ASSERT(pStack);
  OAILOG_DEBUG(
      LOG_GTPV2C,
      "Duplicate request hold timer expired for transaction %p with seqNum "
      "%d\n",
      thiz, thiz->seqNum);
  thiz->hRspTmr = 0;
  RB_REMOVE(
      NwGtpv2cOutstandingRxSeqNumTrxnMap, &(pStack->outstandingRxSeqNumMap),
      thiz);
  rc = nwGtpv2cTrxnDelete(&thiz);
  NW_ASSERT(NW_OK == rc);
  return rc;
}

/**
   Start timer to wait for rsp of a req message

   @param[in] thiz : Pointer to transaction
   @param[in] timeoutCallbackFunc : Timeout handler callback function.
   @return NW_OK on success.
*/

nw_rc_t nwGtpv2cTrxnStartPeerRspWaitTimer(nw_gtpv2c_trxn_t* thiz) {
  nw_rc_t rc;

  rc = nwGtpv2cStartTimer(
      thiz->pStack, thiz->t3Timer, 0, NW_GTPV2C_TMR_TYPE_ONE_SHOT,
      nwGtpv2cTrxnPeerRspWaitTimeout, thiz, &thiz->hRspTmr);
  return rc;
}

/**
  Start timer to wait before pruginf a req tran for which response has been sent

  @param[in] thiz : Pointer to transaction
  @return NW_OK on success.
*/

nw_rc_t nwGtpv2cTrxnStartDulpicateRequestWaitTimer(nw_gtpv2c_trxn_t* thiz) {
  nw_rc_t rc;

  rc = nwGtpv2cStartTimer(
      thiz->pStack, thiz->t3Timer * thiz->maxRetries, 0,
      NW_GTPV2C_TMR_TYPE_ONE_SHOT, nwGtpv2cTrxnDuplicateRequestWaitTimeout,
      thiz, &thiz->hRspTmr);
  return rc;
}

/**
  Send timer stop request to TmrMgr Entity.

  @param[in] thiz : Pointer to transaction
  @return NW_OK on success.
*/

static nw_rc_t nwGtpv2cTrxnStopPeerRspTimer(nw_gtpv2c_trxn_t* thiz) {
  nw_rc_t rc;

  NW_ASSERT(thiz->pStack->tmrMgr.tmrStopCallback != NULL);
  rc            = nwGtpv2cStopTimer(thiz->pStack, thiz->hRspTmr);
  thiz->hRspTmr = 0;
  return rc;
}

/*--------------------------------------------------------------------------*
                        P U B L I C    F U N C T I O N S
  --------------------------------------------------------------------------*/

/**
   Constructor

   @param[in] thiz : Pointer to stack
   @param[out] ppTrxn : Pointer to pointer to Trxn object.
   @return NW_OK on success.
*/
nw_gtpv2c_trxn_t* nwGtpv2cTrxnNew(NW_IN nw_gtpv2c_stack_t* thiz) {
  nw_gtpv2c_trxn_t* pTrxn;

  if (gpGtpv2cTrxnPool) {
    pTrxn            = gpGtpv2cTrxnPool;
    gpGtpv2cTrxnPool = gpGtpv2cTrxnPool->next;
  } else {
    NW_GTPV2C_MALLOC(thiz, sizeof(nw_gtpv2c_trxn_t), pTrxn, nw_gtpv2c_trxn_t*);
  }

  if (pTrxn) {
    OAILOG_DEBUG(
        LOG_GTPV2C,
        "Created not trx without seqNum as transaction %p. Head %p, Next %p\n",
        pTrxn, gpGtpv2cTrxnPool,
        (gpGtpv2cTrxnPool) ? gpGtpv2cTrxnPool->next : NULL);

    pTrxn->pStack     = thiz;
    pTrxn->pMsg       = NULL;
    pTrxn->maxRetries = 2;
    pTrxn->t3Timer    = 10;
    pTrxn->seqNum     = thiz->seqNum;
    pTrxn->pt_trx     = false;

    thiz->seqNum++;  // Increment sequence number

    if (thiz->seqNum == 0x800000) thiz->seqNum = 0;
  }

  OAILOG_DEBUG(LOG_GTPV2C, "Created transaction %p\n", pTrxn);
  return pTrxn;
}

/**
   Overloaded Constructor

   @param[in] thiz : Pointer to stack.
   @param[in] seqNum : Sequence number for this transaction.
   @return Pointer to Trxn object.
*/
nw_gtpv2c_trxn_t* nwGtpv2cTrxnWithSeqNumNew(
    NW_IN nw_gtpv2c_stack_t* thiz, NW_IN uint32_t seqNum) {
  nw_gtpv2c_trxn_t* pTrxn;

  if (gpGtpv2cTrxnPool) {
    pTrxn            = gpGtpv2cTrxnPool;
    gpGtpv2cTrxnPool = gpGtpv2cTrxnPool->next;
  } else {
    NW_GTPV2C_MALLOC(thiz, sizeof(nw_gtpv2c_trxn_t), pTrxn, nw_gtpv2c_trxn_t*);
  }

  if (pTrxn) {
    OAILOG_DEBUG(
        LOG_GTPV2C,
        "Created new trx with seqNum %p as transaction %u. Head %p, Next %p\n",
        pTrxn, seqNum, gpGtpv2cTrxnPool,
        (gpGtpv2cTrxnPool) ? gpGtpv2cTrxnPool->next : NULL);

    pTrxn->pStack     = thiz;
    pTrxn->pMsg       = NULL;
    pTrxn->maxRetries = 2;
    pTrxn->t3Timer    = 2;
    pTrxn->seqNum     = seqNum;
    pTrxn->pMsg       = NULL;
    pTrxn->pt_trx     = false;
  }
  return pTrxn;
}

/**
   Another overloaded constructor. Create transaction as outstanding
   RX transaction for detecting duplicated requests.

   @param[in] thiz : Pointer to stack.
   @param[in] teidLocal : Trxn teid.
   @param[in] peerIp : Peer Ip address.
   @param[in] peerPort : Peer Ip port.
   @param[in] seqNum : Seq Number.
   @return NW_OK on success.
*/

nw_gtpv2c_trxn_t* nwGtpv2cTrxnOutstandingRxNew(
    NW_IN nw_gtpv2c_stack_t* thiz,
    __attribute__((unused)) NW_IN uint32_t teidLocal,
    NW_IN struct sockaddr* peerIp, NW_IN uint32_t peerPort,
    NW_IN uint32_t seqNum) {
  nw_rc_t rc;
  nw_gtpv2c_trxn_t *pTrxn, *pCollision;

  // todo: ipv6 for retransmission1

  if (gpGtpv2cTrxnPool) {
    pTrxn            = gpGtpv2cTrxnPool;
    gpGtpv2cTrxnPool = gpGtpv2cTrxnPool->next;
  } else {
    NW_GTPV2C_MALLOC(thiz, sizeof(nw_gtpv2c_trxn_t), pTrxn, nw_gtpv2c_trxn_t*);
  }

  if (pTrxn) {
    OAILOG_DEBUG(
        LOG_GTPV2C, "Received new Rx transaction %p, Head %p, Next %p\n", pTrxn,
        gpGtpv2cTrxnPool, (gpGtpv2cTrxnPool) ? gpGtpv2cTrxnPool->next : NULL);

    pTrxn->pStack     = thiz;
    pTrxn->maxRetries = 2;
    pTrxn->t3Timer    = 2;
    pTrxn->seqNum     = seqNum;
    memcpy(
        (void*) &pTrxn->peer_ip, peerIp,
        (peerIp->sa_family == AF_INET) ? sizeof(struct sockaddr_in) :
                                         sizeof(struct sockaddr_in6));
    pTrxn->peerPort = peerPort;
    pTrxn->pMsg     = NULL;
    pTrxn->hRspTmr  = 0;
    pTrxn->pt_trx   = false;
    pCollision      = RB_INSERT(
        NwGtpv2cOutstandingRxSeqNumTrxnMap, &(thiz->outstandingRxSeqNumMap),
        pTrxn);

    if (pCollision) {
      OAILOG_WARNING(
          LOG_GTPV2C,
          "Duplicate request message received for seq num 0x%x for trx (%p)!\n",
          (uint32_t) seqNum, pCollision);

      rc = nwGtpv2cTrxnDelete(&pTrxn);
      NW_ASSERT(NW_OK == rc);
      pTrxn = NULL;

      // Case of duplicate request message from peer. Retransmit response.

      if (pCollision->pMsg) {
        rc = pCollision->pStack->udp.udpDataReqCallback(
            pCollision->pStack->udp.hUdp, pCollision->pMsg->msgBuf,
            pCollision->pMsg->msgLen, pCollision->localPort,
            (struct sockaddr*) &pCollision->peer_ip, pCollision->peerPort);
      } else if (pCollision->pt_trx) {
        /** Transaction is PT, continue with processing. */
        OAILOG_DEBUG(
            LOG_GTPV2C,
            "Outstanding RX transaction (%p) with seqNum %d is set as "
            "passthrough, continuing with processing it. \n",
            pCollision, seqNum);
        /** Remove the newly created transaction. */
        pTrxn = pCollision;
      }
    }
  }

  if (pTrxn)
    OAILOG_DEBUG(
        LOG_GTPV2C, "Created outstanding RX transaction %p with seqNum %d \n",
        pTrxn, seqNum);

  return (pTrxn);
}

/**
   Destructor

   @param[out] pthiz : Pointer to pointer to Trxn object.
   @return NW_OK on success.
*/
nw_rc_t nwGtpv2cTrxnDelete(NW_INOUT nw_gtpv2c_trxn_t** pthiz) {
  nw_rc_t rc = NW_OK;
  nw_gtpv2c_stack_t* pStack;
  nw_gtpv2c_trxn_t* thiz = *pthiz;

  pStack = thiz->pStack;

  if (thiz->hRspTmr) {
    rc = nwGtpv2cTrxnStopPeerRspTimer(thiz);
    if (NW_OK != rc) {
      OAILOG_INFO(
          LOG_GTPV2C, "Stopping peer response timer for trxn %p failed!\n",
          thiz);
    }
  }

  if (thiz->pMsg) {
    rc = nwGtpv2cMsgDelete(
        (nw_gtpv2c_stack_handle_t) pStack, (nw_gtpv2c_msg_handle_t) thiz->pMsg);
    NW_ASSERT(NW_OK == rc);
  }

  OAILOG_DEBUG(
      LOG_GTPV2C,
      "Purging  transaction %p with seqNum %d. (before) Head %p, Next %p. \n",
      thiz, thiz->seqNum, gpGtpv2cTrxnPool,
      (gpGtpv2cTrxnPool) ? gpGtpv2cTrxnPool->next : 0);
  thiz->next       = gpGtpv2cTrxnPool;
  gpGtpv2cTrxnPool = thiz;
  *pthiz           = NULL;

  OAILOG_DEBUG(
      LOG_GTPV2C, "After purging  transaction %p, Head %p, Next %p\n", thiz,
      gpGtpv2cTrxnPool, (gpGtpv2cTrxnPool) ? gpGtpv2cTrxnPool->next : NULL);

  return rc;
}

#ifdef __cplusplus
}
#endif

/*--------------------------------------------------------------------------*
                            E N D   O F   F I L E
  --------------------------------------------------------------------------*/
