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

#ifndef __NW_GTPV2C_H__
#define __NW_GTPV2C_H__

#include <sys/types.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include "NwTypes.h"
#include "NwError.h"
#include "sgw_ie_defs.h"
#include "udp_messages_types.h"

/** @mainpage

  @section intro Introduction

  nw-gtpv2c library is a free and open source control plane implementation of
  GPRS Tunneling protocol v2 also known as eGTPc as per 3GPP TS29274-930. The
  library is published under BSD three clause license.

  @section scope Scope

  The stack library also does basic tasks like packet/header validation,
  retransmission and message parsing.

  @section design Design Philosophy

  The stack is fully-asynchronous in design for compatibility with event loop
  mechanisms such as select, poll, etc. and can also be used for multi-threaded
  applications. It should compile on Linux, *BSD, Mac OS X, Solaris and Windows
  (cygwin).

  The stack is designed for high portability not only for the hardware and OS it
  will run on but also for the application software that uses it. The stack
  doesn't mandate conditions on the user application architecture or design. The
  stack relies on the user application for infrastructure utilities such as I/O,
  timers, logs and multithreading. This realized by using callback mechanisms
  and enables the stack library to seamlessly integrate without or very little
  changes to the existing application framework.

  The stack architecture builds upon following mentioned entities that are
  external to it.

  User Layer Protocol (ULP) Entity:
  This layer implements the intelligent logic for the application and sits on
  top of the stack.

  UDP Entity:
  This is the layer below the stack and is responsible for UDP I/O with the
  stack and network. It may or may not be housed in ULP.

  Timer Manager Entity:
  Timer Manager Entity provides the stack with infrastructure for timer CRUD
  operations.

  Log Manager Entity:
  Log Manager Entity provides the stack with callbacks for logging operations.
  It may or may not be housed in ULP.

  The application may implement all above entities as a single or multiple
  object.

  @section applications Applications and Usage

  Please refer sample applications under 'test-app' directory for usage
  examples.

 */

/**
 * @file NwGtpv2c.h
 * @author Amit Chawre
 * @brief
 *
 * This header file contains all required definitions and functions
 * prototypes for using nw-gtpv2c library.
 *
 **/

/*--------------------------------------------------------------------------*
 *            S T A C K    H A N D L E    D E F I N I T I O N S             *
 *--------------------------------------------------------------------------*/

typedef NwPtrT nw_gtpv2c_stack_handle_t; /**< Gtpv2c Stack Handle */
typedef NwPtrT nw_gtpv2c_ulp_handle_t; /**< Gtpv2c Stack Ulp Entity Handle   */
typedef NwPtrT nw_gtpv2c_udp_handle_t; /**< Gtpv2c Stack Udp Entity Handle   */
typedef NwPtrT
    nw_gtpv2c_timer_mgr_handle_t; /**< Gtpv2c Stack Timer Manager Handle  */
typedef NwPtrT
    nw_gtpv2c_mem_mgr_handle_t; /**< Gtpv2c Stack Memory Manager Handle */
typedef NwPtrT
    nw_gtpv2c_log_mgr_handle_t; /**< Gtpv2c Stack Log Manager Handle    */
typedef NwPtrT nw_gtpv2c_timer_handle_t;  /**< Gtpv2c Stack Timer Handle  */
typedef NwPtrT nw_gtpv2c_msg_handle_t;    /**< Gtpv2c Msg Handle    */
typedef NwPtrT nw_gtpv2c_trxn_handle_t;   /**< Gtpv2c Transaction Handle   */
typedef NwPtrT nw_gtpv2c_tunnel_handle_t; /**< Gtpv2c Ulp Tunnel Handle */
typedef NwPtrT
    nw_gtpv2c_ulp_trxn_handle_t; /**< Gtpv2c Ulp Transaction Handle      */
typedef NwPtrT nw_gtpv2c_ulp_tunnel_handle_t; /**< Gtpv2c Ulp Tunnel Handle */

typedef uint8_t nw_gtpv2c_msg_type_t; /**< Gtpv2c Msg Type                    */

typedef struct nw_gtpv2c_stack_config_s {
  uint16_t __tbd;
} nw_gtpv2c_stack_config_t;

/*--------------------------------------------------------------------------*
 *            S T A C K        A P I      D E F I N I T I O N S             *
 *--------------------------------------------------------------------------*/

#define NW_GTPV2C_ULP_API_FLAG_NONE (0x00 << 24)
#define NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL (0x01 << 24)
#define NW_GTPV2C_ULP_API_FLAG_DELETE_LOCAL_TUNNEL (0x02 << 24)
#define NW_GTPV2C_ULP_API_FLAG_IS_COMMAND_MESSAGE (0x03 << 24)

/*---------------------------------------------------------------------------
 * Gtpv2c Stack ULP API type definitions
 *--------------------------------------------------------------------------*/

/**
 * APIs types between ULP and Stack
 */

typedef enum nw_gtpv2c_ulp_api_type_e {
  /* APIs from ULP to stack */

  NW_GTPV2C_ULP_API_INITIAL_REQ = 0x00000000, /**< Send a initial message */
  NW_GTPV2C_ULP_API_TRIGGERED_REQ, /**< Send a triggered req message */
  NW_GTPV2C_ULP_API_TRIGGERED_RSP, /**< Send a triggered rsp message */
  NW_GTPV2C_ULP_API_TRIGGERED_ACK, /**< Send a triggered ack message */

  /* APIs from stack to ULP */

  NW_GTPV2C_ULP_API_INITIAL_REQ_IND,   /**< Receive a initial message from stack
                                        */
  NW_GTPV2C_ULP_API_TRIGGERED_RSP_IND, /**< Recieve a triggered rsp message from
                                          stack */
  NW_GTPV2C_ULP_API_TRIGGERED_REQ_IND, /**< Recieve a triggered req message from
                                          stack */
  NW_GTPV2C_ULP_API_TRIGGERED_ACK_IND, /**< Receive a triggered ACK from stack
                                        */
  NW_GTPV2C_ULP_API_RSP_FAILURE_IND,   /**< Rsp failure for gtpv2 message from
                                          stack   */

  /* Local tunnel management APIs from ULP to stack */

  NW_GTPV2C_ULP_CREATE_LOCAL_TUNNEL, /**< Create a local tunnel */
  NW_GTPV2C_ULP_DELETE_LOCAL_TUNNEL, /**< Delete a local tunnel */

  NW_GTPV2C_ULP_FIND_LOCAL_TUNNEL, /**< FIND a local tunnel */

  /* Do not add below this */
  NW_GTPV2C_ULP_API_END = 0xFFFFFFFF,

} nw_gtpv2c_ulp_api_type_t;

/**
 * Error information of incoming GTP messages
 */

typedef struct nw_gtpv2c_error_s {
  NW_IN uint8_t cause;
  NW_IN uint8_t flags;
  struct {
    NW_IN uint8_t type;
    NW_IN uint8_t instance;
  } offendingIe;
} nw_gtpv2c_error_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c initial message.
 */

typedef struct nw_gtpv2c_initial_req_info_s {
  NW_INOUT nw_gtpv2c_tunnel_handle_t
      hTunnel; /**< Tunnel handle over which the mesasge is to be sent.*/
  NW_IN uint16_t t3Timer;
  NW_IN uint16_t maxRetries;
  NW_IN nw_gtpv2c_ulp_trxn_handle_t
      hUlpTrxn; /**< Optional handle to be returned in rsp of this msg. */

  NW_IN struct sockaddr*
      edns_peer_ip; /**< Required only in case when hTunnel == 0            */
  NW_IN uint8_t
      internal_flags; /**< Required only in case when hTunnel == 0            */
  NW_IN uint32_t teidLocal; /**< Required only in case when hTunnel == 0 */
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t
      hUlpTunnel; /**< Required only in case when hTunnel == 0            */
  NW_IN bool noDelete; /**< Set if the tunnel shall not be deleted automatically
                          by the response. */
} nw_gtpv2c_initial_req_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c triggered request message.
 */

typedef struct nw_gtpv2c_triggered_req_info_s {
  NW_IN nw_gtpv2c_tunnel_handle_t
      hTunnel; /**< Tunnel handle over which the mesasge is to be sent */
  NW_IN nw_gtpv2c_trxn_handle_t hTrxn; /**< Request Trxn handle which to which
                                          triggered req is being sent */
  NW_IN uint16_t t3Timer;
  NW_IN uint16_t maxRetries;
  NW_IN nw_gtpv2c_ulp_trxn_handle_t
      hUlpTrxn; /**< Optional handle to be returned in rsp of this msg. */
  NW_IN struct sockaddr_in*
      peerIp; /**< Required only in case when hTunnel == 0            */
  NW_IN uint32_t teidLocal; /**< Required only in case when hTunnel == 0 */
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t
      hUlpTunnel; /**< Required only in case when hTunnel == 0            */

} nw_gtpv2c_triggered_req_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c triggered response message.
 */

typedef struct nw_gtpv2c_triggered_rsp_info_s {
  NW_IN nw_gtpv2c_trxn_handle_t hTrxn; /**< Request Trxn handle which to which
                                          triggered rsp is being sent */
  NW_IN uint32_t teidLocal;            /**< Required only if
                                          NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL is set
                                          to flags. */
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t
      hUlpTunnel; /**< Required only if
                     NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL
                     is set to flags. */
  NW_IN bool
      pt_trx; /**< Make the transaction passthrough, such that the message is
                 forwarded, if no msg is appended to the trx. */

  NW_OUT nw_gtpv2c_tunnel_handle_t
      hTunnel; /**< Returned only in case flags is set to
                NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL */
} nw_gtpv2c_triggered_rsp_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c triggered acknowledgement message.
 */

typedef struct nw_gtpv2c_triggered_ack_info_s {
  NW_IN nw_gtpv2c_trxn_handle_t hTrxn; /**< Request Trxn handle which to which
                                          triggered rsp is being sent */
  NW_IN uint32_t teidLocal;            /**< Required only if
                                          NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL is set
                                          to flags. */
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t
      hUlpTunnel; /**< Required only if
                     NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL
                     is set to flags. */

  NW_OUT nw_gtpv2c_tunnel_handle_t
      hTunnel; /**< Returned only in case flags is set to
                NW_GTPV2C_ULP_API_FLAG_CREATE_LOCAL_TUNNEL */
  NW_IN struct sockaddr* peerIp;
  NW_IN uint32_t peerPort;
  NW_IN uint32_t localPort;
} nw_gtpv2c_triggered_ack_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c initial message.
 */

typedef struct nw_gtpv2c_initial_req_ind_info_s {
  NW_IN nw_gtpv2c_error_t error;
  NW_IN nw_gtpv2c_trxn_handle_t hTrxn;
  NW_IN nw_gtpv2c_ulp_trxn_handle_t hUlpTrxn;
  NW_IN nw_gtpv2c_msg_type_t msgType;
  NW_IN struct sockaddr* peerIp;
  NW_IN uint32_t peerPort;
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
  NW_INOUT nw_gtpv2c_tunnel_handle_t hTunnel;
} nw_gtpv2c_initial_req_ind_info_t;

typedef nw_gtpv2c_initial_req_info_t NwGtpv2cFindInfoT;
typedef nw_gtpv2c_initial_req_ind_info_t nw_gtpv2c_triggered_ack_ind_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c triggered request message.
 */

typedef struct nw_gtpv2c_triggered_req_ind_info_s {
  NW_IN nw_gtpv2c_error_t error;
  NW_IN nw_gtpv2c_trxn_handle_t hTrxn;
  NW_IN nw_gtpv2c_ulp_trxn_handle_t hUlpTrxn;
  NW_IN nw_gtpv2c_msg_type_t msgType;
  NW_IN uint32_t seqNum;
  NW_IN uint32_t teidLocal;
  NW_IN uint32_t teidRemote;
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
} nw_gtpv2c_triggered_req_ind_info_t;

/**
 * API information elements between ULP and Stack for
 * sending a Gtpv2c triggered response message.
 */

typedef struct nw_gtpv2c_triggered_rsp_ind_info_s {
  NW_IN nw_gtpv2c_error_t error;
  NW_IN nw_gtpv2c_ulp_trxn_handle_t hUlpTrxn;
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
  NW_IN uint8_t trx_flags;
  NW_IN nw_gtpv2c_msg_type_t msgType;
  NW_IN bool noDelete;
  NW_IN struct sockaddr* peerIp;
  NW_IN uint32_t localPort;
  NW_IN uint32_t peerPort;
} nw_gtpv2c_triggered_rsp_ind_info_t;

/**
 * API information elements between ULP and Stack for
 * receving a path failure indication from stack.
 */

typedef struct nw_gtpv2c_rsp_failure_ind_info_s {
  NW_IN nw_gtpv2c_ulp_trxn_handle_t hUlpTrxn;
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
  NW_IN nw_gtpv2c_msg_type_t msgType;
  NW_IN bool noDelete;
  NW_IN uint8_t trx_flags;

  NW_IN uint32_t teidLocal;
} nw_gtpv2c_rsp_failure_ind_info_t;

/**
 * API information elements between ULP and Stack for
 * creating local tunnel.
 */

typedef struct nw_gtpv2c_create_local_tunnel_info_s {
  NW_OUT nw_gtpv2c_tunnel_handle_t hTunnel;
  NW_IN nw_gtpv2c_ulp_tunnel_handle_t hUlpTunnel;
  NW_IN uint32_t teidLocal;
  NW_IN struct sockaddr* peerIp;
} nw_gtpv2c_create_local_tunnel_info_t;

/**
 * API information elements between ULP and Stack for
 * deleting a local tunnel.
 */

typedef struct nw_gtpv2c_delete_local_tunnel_info_s {
  NW_IN nw_gtpv2c_tunnel_handle_t hTunnel;
} nw_gtpv2c_delete_local_tunnel_info_t;

/**
 * API container structure between ULP and Stack.
 */

typedef struct nw_gtpv2c_ulp_api_s {
  nw_gtpv2c_ulp_api_type_t
      apiType; /**< First bytes of this field is used as flag holder   */
  nw_gtpv2c_msg_handle_t hMsg; /**< Handle associated with this API */
  union {
    nw_gtpv2c_initial_req_info_t initialReqInfo;
    nw_gtpv2c_triggered_rsp_info_t triggeredRspInfo;
    nw_gtpv2c_triggered_req_info_t triggeredReqInfo;
    nw_gtpv2c_triggered_ack_info_t triggeredAckInfo;
    nw_gtpv2c_initial_req_ind_info_t initialReqIndInfo;
    nw_gtpv2c_triggered_ack_ind_info_t triggeredAckIndInfo;
    nw_gtpv2c_triggered_rsp_ind_info_t triggeredRspIndInfo;
    nw_gtpv2c_triggered_req_ind_info_t triggeredReqIndInfo;
    nw_gtpv2c_rsp_failure_ind_info_t rspFailureInfo;
    nw_gtpv2c_create_local_tunnel_info_t createLocalTunnelInfo;
    nw_gtpv2c_delete_local_tunnel_info_t deleteLocalTunnelInfo;

    // todo: remove this one
    NwGtpv2cFindInfoT findLocalTunnelInfo;
  } u_api_info;
} nw_gtpv2c_ulp_api_t;

/*--------------------------------------------------------------------------*
 *           S T A C K    E N T I T I E S    D E F I N I T I O N S          *
 *--------------------------------------------------------------------------*/

/**
 * Gtpv2c ULP entity definition
 */

typedef struct nw_gtpv2c_ulp_entity_s {
  nw_gtpv2c_ulp_handle_t hUlp;
  nw_rc_t (*ulpReqCallback)(
      NW_IN nw_gtpv2c_ulp_handle_t hUlp, NW_IN nw_gtpv2c_ulp_api_t* pUlpApi);
} nw_gtpv2c_ulp_entity_t;

/**
 * Gtpv2c UDP entity definition
 */

typedef struct nw_gtpv2c_udp_entity_s {
  nw_gtpv2c_udp_handle_t hUdp;
  uint16_t gtpv2cStandardPort;
  nw_rc_t (*udpDataReqCallback)(
      NW_IN nw_gtpv2c_udp_handle_t udpHandle, NW_IN uint8_t* dataBuf,
      NW_IN uint32_t dataSize, NW_IN uint16_t localPort,
      NW_IN struct sockaddr* peerIp, NW_IN uint16_t peerPort);
} nw_gtpv2c_udp_entity_t;

/**
 * Gtpv2c Memory Manager entity definition
 */

typedef struct nw_gtpv2c_mem_mgr_entity_s {
  nw_gtpv2c_mem_mgr_handle_t hMemMgr;
  void* (*memAlloc)(
      NW_IN nw_gtpv2c_mem_mgr_handle_t hMemMgr, NW_IN uint32_t memSize,
      NW_IN char* fileName, NW_IN uint32_t lineNumber);

  void (*memFree)(
      NW_IN nw_gtpv2c_mem_mgr_handle_t hMemMgr, NW_IN void* hMem,
      NW_IN char* fileName, NW_IN uint32_t lineNumber);
} nw_gtpv2c_mem_mgr_entity_t;

#define NW_GTPV2C_TMR_TYPE_ONE_SHOT (0)
#define NW_GTPV2C_TMR_TYPE_REPETITIVE (1)
/**
 * Gtpv2c Timer Manager entity definition
 */

typedef struct nw_gtpv2c_timer_mgr_entity_s {
  nw_gtpv2c_timer_mgr_handle_t tmrMgrHandle;
  nw_rc_t (*tmrStartCallback)(
      NW_IN nw_gtpv2c_timer_mgr_handle_t tmrMgrHandle,
      NW_IN uint32_t timeoutSec, NW_IN uint32_t timeoutUsec,
      NW_IN uint32_t tmrType, NW_IN void* tmrArg,
      NW_OUT nw_gtpv2c_timer_handle_t* tmrHandle);

  nw_rc_t (*tmrStopCallback)(
      NW_IN nw_gtpv2c_timer_mgr_handle_t tmrMgrHandle,
      NW_IN nw_gtpv2c_timer_handle_t tmrHandle);
} nw_gtpv2c_timer_mgr_entity_t;

/**
 * Gtpv2c Log manager entity definition
 */

typedef struct nw_gtpv2c_log_mgr_entity_s {
  nw_gtpv2c_log_mgr_handle_t logMgrHandle;
  nw_rc_t (*logReqCallback)(
      NW_IN nw_gtpv2c_log_mgr_handle_t logMgrHandle, NW_IN uint32_t logLevel,
      NW_IN char* filename, NW_IN uint32_t line, NW_IN char* logStr);
} nw_gtpv2c_log_mgr_entity_t;

/*--------------------------------------------------------------------------*
 *                     P U B L I C   F U N C T I O N S                      *
 *--------------------------------------------------------------------------*/

#ifdef __cplusplus
extern "C" {
#endif

/**
 Constructor. Initialize nw-gtpv2c stack instance.

 @param[in,out] phGtpcStackHandle : Pointer to stack instance handle
 */

nw_rc_t nwGtpv2cInitialize(
    NW_INOUT nw_gtpv2c_stack_handle_t* phGtpcStackHandle);

/**
 Destructor. Destroy nw-gtpv2c stack instance .

 @param[in] hGtpcStackHandle : Stack instance handle
 */

nw_rc_t nwGtpv2cFinalize(NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle);

/**
 Set Configuration for the nw-gtpv2c stack.

 @param[in,out] phGtpcStackHandle : Pointer to stack handle
 */

nw_rc_t nwGtpv2cConfigSet(
    NW_IN nw_gtpv2c_stack_handle_t* phGtpcStackHandle,
    NW_IN nw_gtpv2c_stack_config_t* pConfig);

//#define T3_TIMER  10

/**
 Get Configuration for the nw-gtpv2c stack.

 @param[in,out] phGtpcStackHandle : Pointer to stack handle
 */

nw_rc_t nwGtpv2cConfigGet(
    NW_IN nw_gtpv2c_stack_handle_t* phGtpcStackHandle,
    NW_OUT nw_gtpv2c_stack_config_t* pConfig);

/**
 Set ULP entity for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pUlpEntity : Pointer to ULP entity.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetUlpEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_ulp_entity_t* pUlpEntity);

/**
 Set UDP entity for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pUdpEntity : Pointer to UDP entity.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetUdpEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_udp_entity_t* pUdpEntity);

/**
 Set MemMgr entity for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pMemMgr : Pointer to Memory Manager.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetMemMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_mem_mgr_entity_t* pMemMgr);

/**
 Set TmrMgr entity for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pTmrMgr : Pointer to Timer Manager.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetTimerMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_timer_mgr_entity_t* pTmrMgr);

/**
 Set LogMgr entity for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pLogMgr : Pointer to Log Manager.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetLogMgrEntity(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_log_mgr_entity_t* pLogMgr);

/**
 Set log level for the stack.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] logLevel : Log Level.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cSetLogLevel(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint32_t logLevel);

/**
 Process Data Request from UDP entity.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] udpData : Pointer to received UDP data.
 @param[in] udpDataLen : Received data length.
 @param[in] localPort : Received on local port.
 @param[in] dstPort : Received on port.
 @param[in] from : Received from peer information.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cProcessUdpReq(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle, NW_IN uint8_t* udpData,
    NW_IN uint32_t udpDataLen, NW_IN uint16_t localPort,
    NW_IN uint16_t peerPort, NW_IN struct sockaddr* peerIp);

/**
 Process Request from ULP entity.

 @param[in] hGtpcStackHandle : Stack handle
 @param[in] pLogMgr : Pointer to Ulp Req.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cProcessUlpReq(
    NW_IN nw_gtpv2c_stack_handle_t hGtpcStackHandle,
    NW_IN nw_gtpv2c_ulp_api_t* ulpReq);

/**
 Process Timer timeout Request from Timer Manager

 @param[in] pLogMgr : Pointer timeout arguments.
 @return NW_OK on success.
 */

nw_rc_t nwGtpv2cProcessTimeout(NW_IN void* timeoutArg);

#ifdef __cplusplus
}
#endif

#endif /* __NW_GTPV2C_H__ */

/*--------------------------------------------------------------------------*
 *                      E N D     O F    F I L E                            *
 *--------------------------------------------------------------------------*/
