/*----------------------------------------------------------------------------*
 *                                                                            *
              M I N I M A L I S T I C     U L P     E N T I T Y
 *                                                                            *
                      Copyright (C) 2010 Amit Chawre.
 *                                                                            *
  ----------------------------------------------------------------------------*/

/**
   @file NwMiniUlpEntity.c
   @brief This file contains example of a minimalistic ULP entity.
*/

#include <stdio.h>
#include <string.h>
#include <assert.h>
#include <sys/time.h>
#include "NwEvt.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgParser.h"
#include "NwMiniLogMgrEntity.h"
#include "NwMiniUlpEntity.h"

#ifndef NW_ASSERT
#define NW_ASSERT assert
#endif

#ifdef __cplusplus
extern "C" {
#endif

static NwCharT* gLogLevelStr[] = {"EMER", "ALER", "CRIT", "ERRO",
                                  "WARN", "NOTI", "INFO", "DEBG"};

/*---------------------------------------------------------------------------
   Public Functions
  --------------------------------------------------------------------------*/

nw_rc_t nwGtpv2cUlpInit(
    NwGtpv2cNodeUlpT* thiz, nw_gtpv2c_StackHandleT hGtpv2cStack,
    char* peerIpStr) {
  nw_rc_t rc;

  thiz->hGtpv2cStack = hGtpv2cStack;
  strcpy(thiz->peerIpStr, peerIpStr);
  return NW_OK;
}

nw_rc_t nwGtpv2cUlpDestroy(NwGtpv2cNodeUlpT* thiz) {
  NW_ASSERT(thiz);
  memset(thiz, 0, sizeof(NwGtpv2cNodeUlpT));
  return NW_OK;
}

typedef struct NwGtpv2cPeerS {
  uint32_t ipv4Addr;
  uint32_t pingCount;
  uint32_t pingInterval;
  uint32_t t3Time;
  uint32_t n3Count;

  uint32_t sendTimeStamp;
  nw_gtpv2c_TunnelHandleT hTunnel;
} NwGtpv2cPeerT;

NwGtpv2cPeerT* nwGtpv2cUlpCreatePeerContext(
    NwGtpv2cNodeUlpT* thiz, uint32_t peerIp) {
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulpReq;
  NwGtpv2cPeerT* pPeer = (NwGtpv2cPeerT*) malloc(sizeof(NwGtpv2cPeerT));

  if (pPeer) {
    pPeer->ipv4Addr = peerIp;
    /*
     * Send Message Request to Gtpv2c Stack Instance
     */
    ulpReq.apiType = NW_GTPV2C_ULP_CREATE_LOCAL_TUNNEL;
    ulpReq.u_api_info.createLocalTunnelInfo.hTunnel = 0;
    ulpReq.u_api_info.createLocalTunnelInfo.hUlpTunnel =
        (nw_gtpv2c_ulp_trxn_handle_t) thiz;
    ulpReq.u_api_info.createLocalTunnelInfo.teidLocal =
        (nw_gtpv2c_ulp_trxn_handle_t) 0;
    ulpReq.u_api_info.createLocalTunnelInfo.peerIp = htonl(peerIp);
    rc = nwGtpv2cProcessUlpReq(thiz->hGtpv2cStack, &ulpReq);
    NW_ASSERT(NW_OK == rc);
    pPeer->hTunnel = ulpReq.u_api_info.createLocalTunnelInfo.hTunnel;
  }

  return pPeer;
}

nw_rc_t nwGtpv2cUlpSendEchoRequestToPeer(
    NwGtpv2cNodeUlpT* thiz, NwGtpv2cPeerT* pPeer) {
  nw_rc_t rc;
  struct timeval tv;
  nw_gtpv2c_ulp_api_t ulpReq;

  /*
   * Send Message Request to Gtpv2c Stack Instance
   */
  ulpReq.apiType                           = NW_GTPV2C_ULP_API_INITIAL_REQ;
  ulpReq.u_api_info.initialReqInfo.hTunnel = pPeer->hTunnel;
  ulpReq.u_api_info.initialReqInfo.hUlpTrxn =
      (nw_gtpv2c_ulp_trxn_handle_t) pPeer;
  ulpReq.u_api_info.initialReqInfo.hUlpTunnel =
      (nw_gtpv2c_ulp_tunnel_handle_t) pPeer;
  rc = nwGtpv2cMsgNew(
      thiz->hGtpv2cStack, NW_FALSE, NW_GTP_ECHO_REQ, 0, 0, &(ulpReq.hMsg));
  NW_ASSERT(NW_OK == rc);
  rc = nwGtpv2cMsgAddIeTV1(
      (ulpReq.hMsg), NW_GTPV2C_IE_RECOVERY, 0, thiz->restartCounter);
  NW_ASSERT(NW_OK == rc);
  NW_ASSERT(gettimeofday(&tv, NULL) == 0);
  pPeer->sendTimeStamp = (tv.tv_sec * 1000000) + tv.tv_usec;
  rc                   = nwGtpv2cProcessUlpReq(thiz->hGtpv2cStack, &ulpReq);
  NW_ASSERT(NW_OK == rc);
  return NW_OK;
}

nw_rc_t nwGtpv2cUlpPing(
    NwGtpv2cNodeUlpT* thiz, uint32_t peerIp, uint32_t pingCount,
    uint32_t pingInterval, uint32_t t3Time, uint32_t n3Count) {
  nw_rc_t rc;
  NwGtpv2cPeerT* pPeer;
  nw_gtpv2c_ulp_api_t ulpReq;

  pPeer               = nwGtpv2cUlpCreatePeerContext(thiz, peerIp);
  pPeer->pingCount    = pingCount;
  pPeer->pingInterval = pingInterval;
  pPeer->t3Time       = t3Time;
  pPeer->n3Count      = n3Count;
  /*
   * Send Echo Request to peer
   */
  rc = nwGtpv2cUlpSendEchoRequestToPeer(thiz, pPeer);
  return rc;
}

nw_rc_t nwGtpv2cUlpProcessStackReqCallback(
    nw_gtpv2c_UlpHandleT hUlp, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc;
  uint32_t seqNum;
  uint32_t len;
  uint32_t recvTimeStamp;
  struct timeval tv;
  NwGtpv2cPeerT* pPeer;
  NwGtpv2cNodeUlpT* thiz;

  NW_ASSERT(pUlpApi != NULL);
  thiz = (NwGtpv2cNodeUlpT*) hUlp;

  switch (pUlpApi->apiType) {
    case NW_GTPV2C_ULP_API_TRIGGERED_RSP_IND: {
      pPeer = (NwGtpv2cPeerT*) pUlpApi->u_api_info.triggeredRspIndInfo.hUlpTrxn;

      if (pUlpApi->u_api_info.triggeredRspIndInfo.msgType == NW_GTP_ECHO_RSP) {
        seqNum = nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg);
        len    = nwGtpv2cMsgGetLength(pUlpApi->hMsg);
        NW_ASSERT(gettimeofday(&tv, NULL) == 0);
        recvTimeStamp = (tv.tv_sec * 1000000) + tv.tv_usec;
        NW_LOG(
            NW_LOG_LEVEL_NOTI,
            "%u bytes of response from " NW_IPV4_ADDR
            ": gtp_seq=%u time=%2.2f ms",
            len, NW_IPV4_ADDR_FORMAT(pPeer->ipv4Addr), seqNum,
            (float) (recvTimeStamp - pPeer->sendTimeStamp) / 1000);

        if (pPeer->pingCount) {
          sleep(pPeer->pingInterval);
          rc = nwGtpv2cUlpSendEchoRequestToPeer(thiz, pPeer);

          if (pPeer->pingCount != 0xffffffff) pPeer->pingCount--;
        }
      }
    } break;

    case NW_GTPV2C_ULP_API_RSP_FAILURE_IND: {
      pPeer = (NwGtpv2cPeerT*) pUlpApi->u_api_info.rspFailureInfo.hUlpTrxn;
      NW_LOG(
          NW_LOG_LEVEL_DEBG, "No response from " NW_IPV4_ADDR " (2123)!",
          NW_IPV4_ADDR_FORMAT(pPeer->ipv4Addr));
      rc = nwGtpv2cUlpSendEchoRequestToPeer(thiz, pPeer);
    } break;

    default:
      NW_LOG(NW_LOG_LEVEL_WARN, "Received undefined UlpApi from gtpv2c stack!");
  }

  return NW_OK;
}

#ifdef __cplusplus
}
#endif
