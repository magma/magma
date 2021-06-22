/*----------------------------------------------------------------------------*
 *                                                                            *
 *                              n w - g t p v 2 c                             *
 *    G P R S   T u n n e l i n g    P r o t o c o l   v 2 c    S t a c k     *
 *                                                                            *
 *                                                                            *
 * Copyright (c) 2010-2011 Amit Chawre                                        *
 * All rights reserved.                                                       *
 *                                                                            *
 * Redistribution and use in source and binary forms, with or without         *
 * modification, are permitted provided that the following conditions         *
 * are met:                                                                   *
 *                                                                            *
 * 1. Redistributions of source code must retain the above copyright          *
 *    notice, this list of conditions and the following disclaimer.           *
 * 2. Redistributions in binary form must reproduce the above copyright       *
 *    notice, this list of conditions and the following disclaimer in the     *
 *    documentation and/or other materials provided with the distribution.    *
 * 3. The name of the author may not be used to endorse or promote products   *
 *    derived from this software without specific prior written permission.   *
 *                                                                            *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR       *
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES  *
 * OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.    *
 * IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT,           *
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT   *
 * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,  *
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY      *
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT        *
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF   *
 * THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.          *
 *----------------------------------------------------------------------------*/

#ifndef __NW_GTPV2C_PRIVATE_H__
#define __NW_GTPV2C_PRIVATE_H__

#include <sys/time.h>

#include "assertions.h"
#include "tree.h"
#include "queue.h"

#include "NwTypes.h"
#include "NwError.h"
#include "NwGtpv2c.h"
#include "NwGtpv2cIe.h"
#include "NwGtpv2cMsg.h"
#include "NwGtpv2cMsgIeParseInfo.h"
#include "NwGtpv2cTunnel.h"

/**
 * @file NwGtpv2cPrivate.h
 * @brief This header file contains nw-gtpv2c private definitions not to be
 * exposed to user application.
 */

#ifdef __cplusplus
extern "C" {
#endif

#define NW_GTPV2C_MALLOC(_stack, _size, _mem, _type)                           \
  do {                                                                         \
    if (((nw_gtpv2c_stack_t*) (_stack))->memMgr.memAlloc &&                    \
        ((nw_gtpv2c_stack_t*) (_stack))->memMgr.memFree) {                     \
      _mem = (_type)((nw_gtpv2c_stack_t*) (_stack))                            \
                 ->memMgr.memAlloc(                                            \
                     ((nw_gtpv2c_stack_t*) (_stack))->memMgr.hMemMgr, _size,   \
                     __FILE__, __LINE__);                                      \
      Fatal("Do not use this Mem manager");                                    \
    } else {                                                                   \
      _mem = (_type) malloc(_size);                                            \
    }                                                                          \
  } while (0)

#define NW_GTPV2C_FREE(_stack, _mem)                                           \
  do {                                                                         \
    if (((nw_gtpv2c_stack_t*) (_stack))->memMgr.memAlloc &&                    \
        ((nw_gtpv2c_stack_t*) (_stack))->memMgr.memFree) {                     \
      ((nw_gtpv2c_stack_t*) (_stack))                                          \
          ->memMgr.memFree(                                                    \
              ((nw_gtpv2c_stack_t*) (_stack))->memMgr.hMemMgr, _mem, __FILE__, \
              __LINE__);                                                       \
      Fatal("Do not use this Mem manager");                                    \
    } else {                                                                   \
      free((void*) _mem);                                                      \
    }                                                                          \
  } while (0)

/*--------------------------------------------------------------------------*
 *  G T P V 2 C   S T A C K   O B J E C T   T Y P E    D E F I N I T I O N  *
 *--------------------------------------------------------------------------*/

/**
 * gtpv2c stack class definition
 */

typedef struct nw_gtpv2c_stack_s {
  uint32_t id;
  nw_gtpv2c_ulp_entity_t ulp;
  nw_gtpv2c_udp_entity_t udp;
  nw_gtpv2c_mem_mgr_entity_t memMgr;
  nw_gtpv2c_timer_mgr_entity_t tmrMgr;
  nw_gtpv2c_log_mgr_entity_t logMgr;

  uint32_t seqNum;
  uint32_t logLevel;
  uint32_t restartCounter;

  nw_gtpv2c_msg_ie_parse_info_t* pGtpv2cMsgIeParseInfo[NW_GTP_MSG_END];
  struct nw_gtpv2c_timeout_info_s* activeTimerInfo;

  RB_HEAD(NwGtpv2cTunnelMap, nw_gtpv2c_tunnel_s) tunnelMap;
  RB_HEAD(NwGtpv2cOutstandingTxSeqNumTrxnMap, nw_gtpv2c_trxn_s)
  outstandingTxSeqNumMap;
  RB_HEAD(NwGtpv2cOutstandingRxSeqNumTrxnMap, nw_gtpv2c_trxn_s)
  outstandingRxSeqNumMap;
  RB_HEAD(NwGtpv2cActiveTimerList, nw_gtpv2c_timeout_info_s) activeTimerList;
  NwPtrT hTmrMinHeap;
} nw_gtpv2c_stack_t;

/*--------------------------------------------------------------------------*
 * Timeout Info Type Definition
 *--------------------------------------------------------------------------*/

/**
 * gtpv2c timeout info
 */

typedef struct nw_gtpv2c_timeout_info_s {
  nw_gtpv2c_stack_handle_t hStack;
  struct timeval tvTimeout;
  uint32_t tmrType;
  void* timeoutArg;
  nw_rc_t (*timeoutCallbackFunc)(void*);
  nw_gtpv2c_timer_handle_t hTimer;
  RB_ENTRY(nw_gtpv2c_timeout_info_s)
  activeTimerListRbtNode; /**< RB Tree Data Structure Node        */
  uint32_t timerMinHeapIndex;
  struct nw_gtpv2c_timeout_info_s* next;
} nw_gtpv2c_timeout_info_t;

/*---------------------------------------------------------------------------
 * GTPv2c Message Container Definition
 *--------------------------------------------------------------------------*/

#define NW_GTPV2C_MAX_MSG_LEN                                                  \
  (4096) /**< Maximum supported gtpv2c packet length including header */

/**
 * NwGtpv2cMsgT holds gtpv2c messages to/from the peer.
 */
typedef struct nw_gtpv2c_msg_s {
  uint8_t version;
  uint8_t teidPresent;
  uint8_t msgType;
  uint16_t msgLen;
  uint32_t teid;
  uint32_t seqNum;
  uint8_t* pMsgStart;

#define NW_GTPV2C_MAX_GROUPED_IE_DEPTH (2)
  struct {
    nw_gtpv2c_ie_tlv_t* pIe[NW_GTPV2C_MAX_GROUPED_IE_DEPTH];
    uint8_t top;
  } groupedIeEncodeStack;

  bool isIeValid[NW_GTPV2C_IE_TYPE_MAXIMUM][NW_GTPV2C_IE_INSTANCE_MAXIMUM];
  uint8_t* pIe[NW_GTPV2C_IE_TYPE_MAXIMUM][NW_GTPV2C_IE_INSTANCE_MAXIMUM];
  uint8_t msgBuf[NW_GTPV2C_MAX_MSG_LEN];
  nw_gtpv2c_stack_handle_t hStack;
  struct nw_gtpv2c_msg_s* next;
} nw_gtpv2c_msg_t;

/**
 * Transaction structure
 */

typedef struct nw_gtpv2c_trxn_s {
  uint32_t seqNum;
  uint32_t teidLocal;
  union {
    struct sockaddr_in addrv4;
    struct sockaddr_in6 addrv6;
  } peer_ip;
  uint32_t localPort;
  uint32_t peerPort;
  uint32_t noDelete;
  uint8_t t3Timer;
  uint8_t maxRetries;
  nw_gtpv2c_msg_t* pMsg;
  bool pt_trx; /**< Make the transaction passthrough, such that the message is
                  forwarded, if no msg is appended to the trx. */

  nw_gtpv2c_stack_t* pStack;
  nw_gtpv2c_timer_handle_t hRspTmr;  /**< Handle to reponse timer            */
  nw_gtpv2c_tunnel_handle_t hTunnel; /**< Handle to local tunnel context     */
  nw_gtpv2c_ulp_trxn_handle_t hUlpTrxn; /**< Handle to ULP tunnel context */
  uint8_t trx_flags; /**< Flags in the trx to be signalized back. */
  RB_ENTRY(nw_gtpv2c_trxn_s)
  outstandingTxSeqNumMapRbtNode; /**< RB Tree Data Structure Node        */
  RB_ENTRY(nw_gtpv2c_trxn_s)
  outstandingRxSeqNumMapRbtNode; /**< RB Tree Data Structure Node        */
  struct nw_gtpv2c_trxn_s* next;
} nw_gtpv2c_trxn_t;

/**
 *  GTPv2c Path Context
 */

typedef struct NwGtpv2cPathS {
  uint32_t hUlpPath; /**< Handle to ULP path contect         */
  uint32_t ipv4Address;
  uint32_t restartCounter;
  uint16_t t3ResponseTimout;
  uint16_t n3RequestCount;
  nw_gtpv2c_timer_handle_t
      hKeepAliveTmr; /**< Handle to path keep alive echo timer */
  RB_ENTRY(NwGtpv2cPathS) pathMapRbtNode;
} NwGtpv2cPathT;

RB_PROTOTYPE(
    NwGtpv2cTunnelMap, nw_gtpv2c_tunnel_s, tunnelMapRbtNode,
    nwGtpv2cCompareTunnel)
RB_PROTOTYPE(
    NwGtpv2cOutstandingTxSeqNumTrxnMap, nw_gtpv2c_trxn_s,
    outstandingTxSeqNumMapRbtNode, nwGtpv2cCompareSeqNum)
RB_PROTOTYPE(
    NwGtpv2cOutstandingRxSeqNumTrxnMap, nw_gtpv2c_trxn_s,
    outstandingRxSeqNumMapRbtNode, nwGtpv2cCompareSeqNum)
RB_PROTOTYPE(
    NwGtpv2cActiveTimerList, nw_gtpv2c_timeout_info_s, activeTimerListRbtNode,
    nwGtpv2cCompareOutstandingTxRexmitTime)

/**
 * Start Timer with ULP Timer Manager
 */

nw_rc_t nwGtpv2cStartTimer(
    nw_gtpv2c_stack_t* thiz, uint32_t timeoutSec, uint32_t timeoutUsec,
    uint32_t tmrType, nw_rc_t (*timeoutCallbackFunc)(void*),
    void* timeoutCallbackArg, nw_gtpv2c_timer_handle_t* phTimer);

/**
 * Stop Timer with ULP Timer Manager
 */

nw_rc_t nwGtpv2cStopTimer(
    nw_gtpv2c_stack_t* thiz, nw_gtpv2c_timer_handle_t hTimer);

#ifdef __cplusplus
}
#endif

#endif /* __NW_GTPV2C_PRIVATE_H__ */
/*--------------------------------------------------------------------------*
 *                      E N D     O F    F I L E                            *
 *--------------------------------------------------------------------------*/
