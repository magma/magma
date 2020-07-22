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

static nw_rc_t nwGtpv2cCreateSessionRequestIeIndication(
    uint8_t ieType, uint8_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg) {
  NW_LOG(
      NW_LOG_LEVEL_DEBG,
      "Received IE Parse Indication for of type %u, length %u, instance %u!",
      ieType, ieLength, ieInstance);
  return NW_OK;
}

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

nw_rc_t nwGtpv2cUlpSenEchoRequestToPeer(
    NwGtpv2cNodeUlpT* thiz, uint32_t peerIp) {
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulpReq;

  /*
   * Send Message Request to Gtpv2c Stack Instance
   */
  ulpReq.apiType                           = NW_GTPV2C_ULP_API_INITIAL_REQ;
  ulpReq.apiInfo.initialReqInfo.hTunnel    = 0;
  ulpReq.apiInfo.initialReqInfo.hUlpTrxn   = (nw_gtpv2c_UlpTrxnHandleT) thiz;
  ulpReq.apiInfo.initialReqInfo.hUlpTunnel = (nw_gtpv2c_UlpTrxnHandleT) thiz;
  ulpReq.apiInfo.initialReqInfo.teidLocal  = (nw_gtpv2c_UlpTrxnHandleT) 0;
  ulpReq.apiInfo.initialReqInfo.peerIp     = htonl(peerIp);
  rc                                       = nwGtpv2cMsgNew(
      thiz->hGtpv2cStack, NW_FALSE, NW_GTP_ECHO_REQ, 0, 0, &(ulpReq.hMsg));
  NW_ASSERT(NW_OK == rc);
  rc = nwGtpv2cMsgAddIeTV1(
      (ulpReq.hMsg), NW_GTPV2C_IE_RECOVERY, 0, thiz->restartCounter);
  NW_ASSERT(NW_OK == rc);
  rc = nwGtpv2cProcessUlpReq(thiz->hGtpv2cStack, &ulpReq);
  NW_ASSERT(NW_OK == rc);
  return NW_OK;
}

nw_rc_t nwGtpv2cUlpCreateSessionRequestToPeer(NwGtpv2cNodeUlpT* thiz) {
  nw_rc_t rc;
  nw_gtpv2c_ulp_api_t ulpReq;

  /*
   * Send Message Request to Gtpv2c Stack Instance
   */
  ulpReq.apiType =
      (NW_GTPV2C_ULP_API_INITIAL_REQ |
       NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL);
  ulpReq.apiInfo.initialReqInfo.hUlpTrxn  = (nw_gtpv2c_UlpTrxnHandleT) thiz;
  ulpReq.apiInfo.initialReqInfo.teidLocal = (nw_gtpv2c_UlpTrxnHandleT) thiz;
  ulpReq.apiInfo.initialReqInfo.peerIp    = ntohl(inet_addr(thiz->peerIpStr));
  rc                                      = nwGtpv2cMsgNew(
      thiz->hGtpv2cStack, NW_TRUE, NW_GTP_CREATE_SESSION_REQ, 0, 0,
      &(ulpReq.hMsg));
  NW_ASSERT(NW_OK == rc);
  rc = nwGtpv2cMsgAddIeTV1(
      (ulpReq.hMsg), NW_GTPV2C_IE_RECOVERY, 0, thiz->restartCounter);
  NW_ASSERT(NW_OK == rc);
  rc = nwGtpv2cProcessUlpReq(thiz->hGtpv2cStack, &ulpReq);
  NW_ASSERT(NW_OK == rc);
  return NW_OK;
}

nw_rc_t nwGtpv2cUlpProcessStackReqCallback(
    nw_gtpv2c_UlpHandleT hUlp, nw_gtpv2c_ulp_api_t* pUlpApi) {
  nw_rc_t rc;
  NwGtpv2cNodeUlpT* thiz;

  NW_ASSERT(pUlpApi != NULL);
  thiz = (NwGtpv2cNodeUlpT*) hUlp;

  switch (pUlpApi->apiType) {
    case NW_GTPV2C_ULP_API_INITIAL_REQ_IND: {
      nw_gtpv2c_ulp_api_t ulpReq;
      nw_gtpv2c_msg_parser_t* pMsgParser;

      NW_LOG(
          NW_LOG_LEVEL_DEBG,
          "Received NW_GTPV2C_ULP_API_INITIAL_REQ_IND from gtpv2c stack! %X:%u",
          pUlpApi->apiInfo.initialReqIndInfo.peerIp,
          pUlpApi->apiInfo.initialReqIndInfo.peerPort);
      rc = nwGtpv2cMsgParserNew(
          thiz->hGtpv2cStack, NW_GTP_CREATE_SESSION_REQ,
          nwGtpv2cCreateSessionRequestIeIndication, NULL, &pMsgParser);
      NW_ASSERT(NW_OK == rc);

      if (pUlpApi->apiInfo.initialReqIndInfo.msgType ==
          NW_GTP_CREATE_SESSION_REQ) {
        struct {
          uint8_t causeValue;
          uint8_t spare : 5;
          uint8_t pce : 1;
          uint8_t bce : 1;
          uint8_t cs : 1;
        } cause;

        rc = nwGtpv2cMsgParserAddIe(
            pMsgParser, NW_GTPV2C_IE_RECOVERY, NW_GTPV2C_IE_INSTANCE_ZERO,
            NW_GTPV2C_IE_PRESENCE_MANDATORY,
            nwGtpv2cCreateSessionRequestIeIndication, NULL);
        NW_ASSERT(NW_OK == rc);
        uint8_t offendingIeType, OffendingIeLength;

        rc = nwGtpv2cMsgParserRun(
            pMsgParser, (pUlpApi->hMsg), &offendingIeType, &OffendingIeLength);

        if (rc != NW_OK) {
          switch (rc) {
            case NW_GTPV2C_MANDATORY_IE_MISSING:
              NW_LOG(
                  NW_LOG_LEVEL_ERRO,
                  "Mandatory IE type '%u' of instance '%u' missing!",
                  offendingIeType, OffendingIeLength);
              cause.causeValue = NW_GTPV2C_CAUSE_MANDATORY_IE_MISSING;
              break;

            default:
              NW_LOG(NW_LOG_LEVEL_ERRO, "Unknown message parse error!");
              cause.causeValue = 0;
              break;
          }
        } else {
          cause.causeValue = NW_GTPV2C_CAUSE_REQUEST_ACCEPTED;
        }

        /*
         * Send Message Request to Gtpv2c Stack Instance
         */
        ulpReq.apiType = NW_GTPV2C_ULP_API_TRIGGERED_RSP;
        ulpReq.apiInfo.triggeredRspInfo.hTrxn =
            pUlpApi->apiInfo.initialReqIndInfo.hTrxn;
        rc = nwGtpv2cMsgNew(
            thiz->hGtpv2cStack, NW_TRUE, NW_GTP_CREATE_SESSION_RSP, 0,
            nwGtpv2cMsgGetSeqNumber(pUlpApi->hMsg), &(ulpReq.hMsg));
        rc = nwGtpv2cMsgAddIe(
            (ulpReq.hMsg), NW_GTPV2C_IE_CAUSE, 2, 0, (uint8_t*) &cause);
        NW_ASSERT(NW_OK == rc);
        NW_LOG(
            NW_LOG_LEVEL_NOTI,
            "Received NW_GTP_CREATE_SESSION_REQ, Sending "
            "NW_GTP_CREATE_SESSION_RSP!");
        rc = nwGtpv2cProcessUlpReq(thiz->hGtpv2cStack, &ulpReq);
        NW_ASSERT(NW_OK == rc);
      }

      rc = nwGtpv2cMsgParserDelete(thiz->hGtpv2cStack, pMsgParser);
      NW_ASSERT(NW_OK == rc);
      rc = nwGtpv2cMsgDelete(thiz->hGtpv2cStack, (pUlpApi->hMsg));
      NW_ASSERT(NW_OK == rc);
    } break;

    case NW_GTPV2C_ULP_API_TRIGGERED_RSP_IND: {
      nw_gtpv2c_ulp_api_t ulpReq;
      nw_gtpv2c_msg_parser_t* pMsgParser;

      NW_LOG(
          NW_LOG_LEVEL_DEBG,
          "Received NW_GTPV2C_ULP_API_TRIGGERED_RSP_IND from gtpv2c stack!",
          pUlpApi->apiInfo.triggeredRspIndInfo,
          pUlpApi->apiInfo.triggeredRspIndInfo);
      rc = nwGtpv2cMsgParserNew(
          thiz->hGtpv2cStack, NW_GTP_CREATE_SESSION_REQ,
          nwGtpv2cCreateSessionRequestIeIndication, NULL, &pMsgParser);
      NW_ASSERT(NW_OK == rc);

      if (pUlpApi->apiInfo.triggeredRspIndInfo.msgType ==
          NW_GTP_CREATE_SESSION_RSP) {
        rc = nwGtpv2cMsgParserAddIe(
            pMsgParser, NW_GTPV2C_IE_CAUSE, NW_GTPV2C_IE_INSTANCE_ZERO,
            NW_GTPV2C_IE_PRESENCE_MANDATORY,
            nwGtpv2cCreateSessionRequestIeIndication, NULL);
        NW_ASSERT(NW_OK == rc);
        uint8_t offendingIeType, OffendingIeLength;

        rc = nwGtpv2cMsgParserRun(
            pMsgParser, (pUlpApi->hMsg), &offendingIeType, &OffendingIeLength);

        if (rc != NW_OK) {
          switch (rc) {
            case NW_GTPV2C_MANDATORY_IE_MISSING:
              NW_LOG(NW_LOG_LEVEL_ERRO, "Mandatory IE Missing!");
              break;

            default:
              NW_LOG(NW_LOG_LEVEL_ERRO, "Unknown message parse error!");
              break;
          }
        }

        /*
         * Send Message Request to Gtpv2c Stack Instance
         */
        NW_LOG(
            NW_LOG_LEVEL_NOTI,
            "Received NW_GTP_CREATE_SESSION_RSP, Sending "
            "NW_GTP_CREATE_SESSION_REQ!");
        nwGtpv2cUlpCreateSessionRequestToPeer(thiz);
      }

      rc = nwGtpv2cMsgParserDelete(thiz->hGtpv2cStack, pMsgParser);
      NW_ASSERT(NW_OK == rc);
      rc = nwGtpv2cMsgDelete(thiz->hGtpv2cStack, (pUlpApi->hMsg));
      NW_ASSERT(NW_OK == rc);
    } break;

    case NW_GTPV2C_ULP_API_RSP_FAILURE_IND: {
      NW_LOG(
          NW_LOG_LEVEL_ERRO,
          "Received NW_GTPV2C_ULP_API_RSP_FAILURE from gtpv2c stack for "
          "transaction '%x'!",
          thiz);
    } break;

    default:
      NW_LOG(NW_LOG_LEVEL_WARN, "Received undefined UlpApi from gtpv2c stack!");
  }

  return NW_OK;
}

#ifdef __cplusplus
}
#endif
